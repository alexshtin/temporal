// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package temporal

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pborman/uuid"
	sdkclient "go.temporal.io/sdk/client"
	"go.uber.org/zap"

	"go.temporal.io/server/api/persistenceblobs/v1"
	"go.temporal.io/server/common"
	"go.temporal.io/server/common/archiver"
	"go.temporal.io/server/common/archiver/provider"
	"go.temporal.io/server/common/authorization"
	"go.temporal.io/server/common/cluster"
	"go.temporal.io/server/common/elasticsearch"
	l "go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/loggerimpl"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/messaging"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/persistence"
	persistenceClient "go.temporal.io/server/common/persistence/client"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/resource"
	"go.temporal.io/server/common/rpc"
	"go.temporal.io/server/common/rpc/encryption"
	"go.temporal.io/server/common/service/config"
	"go.temporal.io/server/common/service/config/ringpop"
	"go.temporal.io/server/common/service/dynamicconfig"
	"go.temporal.io/server/service/frontend"
	"go.temporal.io/server/service/history"
	"go.temporal.io/server/service/matching"
	"go.temporal.io/server/service/worker"
)

type (
	Server struct {
		so                *serverOptions
		services          map[string]common.Daemon
		serviceStoppedChs map[string]chan struct{}
		stoppedCh         chan struct{}
		logger            l.Logger
	}
)

// Services is the list of all valid temporal services
var (
	Services = []string{
		primitives.FrontendService,
		primitives.HistoryService,
		primitives.MatchingService,
		primitives.WorkerService,
	}
)

// NewServer returns a new instance of server that serves one or many services.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		so:                newServerOptions(opts),
		services:          make(map[string]common.Daemon),
		serviceStoppedChs: make(map[string]chan struct{}),
	}
	return s
}

// Start temporal server.
func (s *Server) Start() error {
	err := s.so.validate()
	if err != nil {
		return err
	}

	s.stoppedCh = make(chan struct{})

	zapLogger := s.so.config.Log.NewZapLogger()
	s.logger = loggerimpl.NewLogger(zapLogger)

	if err := s.so.config.Global.PProf.NewInitializer(s.logger).Start(); err != nil {
		log.Fatalf("fail to start PProf: %v", err)
	}

	err = ringpop.ValidateRingpopConfig(&s.so.config.Global.Membership)
	if err != nil {
		log.Fatalf("Ringpop config validation error - %v", err)
	}

	tlsFactory, err := encryption.NewTLSConfigProviderFromConfig(s.so.config.Global.TLS)
	if err != nil {
		log.Fatalf("error initializing TLS provider: %v", err)
	}

	dynamicConfig, err := dynamicconfig.NewFileBasedClient(&s.so.config.DynamicConfigClient, s.logger, s.stoppedCh)
	if err != nil {
		log.Printf("error creating file based dynamic config client, use no-op config client instead. error: %v", err)
		dynamicConfig = dynamicconfig.NewNopClient()
	}
	dc := dynamicconfig.NewCollection(dynamicConfig, s.logger)

	// This call performs a config check against the configured persistence store for immutable cluster metadata.
	// If there is a mismatch, the persisted values take precedence and will be written over in the config objects.
	// This is to keep this check hidden from independent downstream daemons and keep this in a single place.
	immutableClusterMetadataInitialization(s.logger, dc, &s.so.config.Persistence, s.so.config.ClusterMetadata)

	clusterMetadata := cluster.NewMetadata(
		s.logger,
		dc.GetBoolProperty(dynamicconfig.EnableGlobalNamespace, s.so.config.ClusterMetadata.EnableGlobalNamespace),
		s.so.config.ClusterMetadata.FailoverVersionIncrement,
		s.so.config.ClusterMetadata.MasterClusterName,
		s.so.config.ClusterMetadata.CurrentClusterName,
		s.so.config.ClusterMetadata.ClusterInformation,
		s.so.config.ClusterMetadata.ReplicationConsumer,
	)

	for _, svcName := range s.so.serviceNames {
		params := resource.BootstrapParams{}
		params.Name = svcName
		params.Logger = s.logger
		params.PersistenceConfig = s.so.config.Persistence
		params.DynamicConfig = dynamicConfig

		svcCfg := s.so.config.Services[svcName]
		rpcFactory := rpc.NewFactory(&svcCfg.RPC, svcName, s.logger, tlsFactory)
		params.RPCFactory = rpcFactory

		// Ringpop uses a different port to register handlers, this map is needed to resolve
		// services to correct addresses used by clients through ServiceResolver lookup API
		servicePortMap := make(map[string]int)
		for svcName, svcCfg := range s.so.config.Services {
			servicePortMap[svcName] = svcCfg.RPC.GRPCPort
		}

		params.MembershipFactoryInitializer =
			func(persistenceBean persistenceClient.Bean, logger l.Logger) (resource.MembershipMonitorFactory, error) {
				return ringpop.NewRingpopFactory(
					&s.so.config.Global.Membership,
					rpcFactory.GetRingpopChannel(),
					svcName,
					servicePortMap,
					logger,
					persistenceBean.GetClusterMetadataManager(),
				)
			}

		params.DCRedirectionPolicy = s.so.config.DCRedirectionPolicy
		metricsScope := svcCfg.Metrics.NewScope(s.logger)
		params.MetricsScope = metricsScope
		metricsClient := metrics.NewClient(metricsScope, metrics.GetMetricsServiceIdx(svcName, s.logger))
		params.MetricsClient = metricsClient
		params.ClusterMetadata = clusterMetadata

		options, err := tlsFactory.GetFrontendClientConfig()
		if err != nil {
			log.Fatalf("unable to load frontend tls configuration: %v", err)
		}

		params.PublicClient, err = sdkclient.NewClient(sdkclient.Options{
			HostPort:     s.so.config.PublicClient.HostPort,
			Namespace:    common.SystemLocalNamespace,
			MetricsScope: metricsScope,
			Logger:       l.NewZapAdapter(zapLogger),
			ConnectionOptions: sdkclient.ConnectionOptions{
				TLS:                options,
				DisableHealthCheck: true,
			},
		})
		if err != nil {
			log.Fatalf("failed to create public client: %v", err)
		}

		advancedVisMode := dc.GetStringProperty(
			dynamicconfig.AdvancedVisibilityWritingMode,
			common.GetDefaultAdvancedVisibilityWritingMode(s.so.config.Persistence.IsAdvancedVisibilityConfigExist()),
		)()
		isAdvancedVisEnabled := advancedVisMode != common.AdvancedVisibilityWritingModeOff
		if clusterMetadata.IsGlobalNamespaceEnabled() {
			params.MessagingClient = messaging.NewKafkaClient(&s.so.config.Kafka, metricsClient, zap.NewNop(), s.logger, metricsScope, true, isAdvancedVisEnabled)
		} else if isAdvancedVisEnabled {
			params.MessagingClient = messaging.NewKafkaClient(&s.so.config.Kafka, metricsClient, zap.NewNop(), s.logger, metricsScope, false, isAdvancedVisEnabled)
		} else {
			params.MessagingClient = nil
		}

		if isAdvancedVisEnabled {
			// verify config of advanced visibility store
			advancedVisStoreKey := s.so.config.Persistence.AdvancedVisibilityStore
			advancedVisStore, ok := s.so.config.Persistence.DataStores[advancedVisStoreKey]
			if !ok {
				log.Fatalf("not able to find advanced visibility store in config: %v", advancedVisStoreKey)
			}

			esClient, err := elasticsearch.NewClient(advancedVisStore.ElasticSearch)
			if err != nil {
				log.Fatalf("error creating elastic search client: %v", err)
			}
			params.ESConfig = advancedVisStore.ElasticSearch
			params.ESClient = esClient

			// verify index name
			indexName, ok := advancedVisStore.ElasticSearch.Indices[common.VisibilityAppName]
			if !ok || len(indexName) == 0 {
				log.Fatalf("elastic search config missing visibility index")
			}
		}

		params.ArchivalMetadata = archiver.NewArchivalMetadata(
			dc,
			s.so.config.Archival.History.State,
			s.so.config.Archival.History.EnableRead,
			s.so.config.Archival.Visibility.State,
			s.so.config.Archival.Visibility.EnableRead,
			&s.so.config.NamespaceDefaults.Archival,
		)

		params.ArchiverProvider = provider.NewArchiverProvider(s.so.config.Archival.History.Provider, s.so.config.Archival.Visibility.Provider)
		params.PersistenceConfig.TransactionSizeLimit = dc.GetIntProperty(dynamicconfig.TransactionSizeLimit, common.DefaultTransactionSizeLimit)
		params.Authorizer = authorization.NewNopAuthorizer()

		var svc common.Daemon
		switch svcName {
		case primitives.FrontendService:
			svc, err = frontend.NewService(&params)
		case primitives.HistoryService:
			svc, err = history.NewService(&params)
		case primitives.MatchingService:
			svc, err = matching.NewService(&params)
		case primitives.WorkerService:
			svc, err = worker.NewService(&params)
		default:
			return fmt.Errorf("uknown service %q", svcName)
		}
		if err != nil {
			s.logger.Fatal("Fail to start service", tag.Service(svcName), tag.Error(err))
		}

		s.services[svcName] = svc
		s.serviceStoppedChs[svcName] = make(chan struct{})

		go func(svc common.Daemon, svcStoppedCh chan<- struct{}) {
			// Start is blocked until Stop() is called.
			svc.Start()
			close(svcStoppedCh)
		}(svc, s.serviceStoppedChs[svcName])

	}

	if s.so.blockingStart {
		// If s.so.interruptCh is nil this will wait forever.
		interruptSignal := <-s.so.interruptCh
		log.Printf("Received %v signal, stopping the server.\n", interruptSignal)
		s.Stop()
	}

	return nil
}

// Stops the server.
func (s *Server) Stop() {
	var wg sync.WaitGroup
	wg.Add(len(s.services))
	close(s.stoppedCh)

	for svcName, svc := range s.services {
		go func(svc common.Daemon, svcName string, svcStoppedCh <-chan struct{}) {
			svc.Stop()
			select {
			case <-svcStoppedCh:
			case <-time.After(time.Minute):
				s.logger.Error("Timed out (1 minute) waiting for service to stop.", tag.Service(svcName))
			}
			wg.Done()
		}(svc, svcName, s.serviceStoppedChs[svcName])
	}
	wg.Wait()
	s.logger.Info("All services are stopped.")
}

func immutableClusterMetadataInitialization(
	logger l.Logger,
	dc *dynamicconfig.Collection,
	persistenceConfig *config.Persistence,
	clusterMetadata *config.ClusterMetadata) {

	logger = logger.WithTags(tag.ComponentMetadataInitializer)
	factory := persistenceClient.NewFactory(
		persistenceConfig,
		dc.GetIntProperty(dynamicconfig.HistoryPersistenceMaxQPS, 3000),
		nil,
		clusterMetadata.CurrentClusterName,
		nil,
		logger,
	)

	clusterMetadataManager, err := factory.NewClusterMetadataManager()
	if err != nil {
		log.Fatalf("Error initializing cluster metadata manager: %v", err)
	}
	defer clusterMetadataManager.Close()

	applied, err := clusterMetadataManager.SaveClusterMetadata(
		&persistence.SaveClusterMetadataRequest{
			ClusterMetadata: persistenceblobs.ClusterMetadata{
				HistoryShardCount: persistenceConfig.NumHistoryShards,
				ClusterName:       clusterMetadata.CurrentClusterName,
				ClusterId:         uuid.New(),
			}})
	if err != nil {
		log.Fatalf("Error while saving cluster metadata: %v", err)
	}
	if applied {
		logger.Info("Successfully saved cluster metadata.")
	} else {
		resp, err := clusterMetadataManager.GetClusterMetadata()
		if err != nil {
			log.Fatalf("Error while fetching cluster metadata: %v", err)
		}
		if clusterMetadata.CurrentClusterName != resp.ClusterName {
			logImmutableMismatch(logger,
				"ClusterMetadata.CurrentClusterName",
				clusterMetadata.CurrentClusterName,
				resp.ClusterName)

			clusterMetadata.CurrentClusterName = resp.ClusterName
		}

		var persistedShardCount = resp.HistoryShardCount
		if persistenceConfig.NumHistoryShards != persistedShardCount {
			logImmutableMismatch(logger,
				"Persistence.NumHistoryShards",
				persistenceConfig.NumHistoryShards,
				persistedShardCount)

			persistenceConfig.NumHistoryShards = persistedShardCount
		}
	}

	metadataManager, err := factory.NewMetadataManager()
	if err != nil {
		log.Fatalf("Error initializing metadata manager: %v", err)
	}
	defer metadataManager.Close()
	if err := metadataManager.InitializeSystemNamespaces(clusterMetadata.CurrentClusterName); err != nil {
		log.Fatalf("failed to register system namespace: %v", err)
	}
}

func logImmutableMismatch(l l.Logger, key string, ignored interface{}, value interface{}) {
	l.Error(
		"Supplied configuration key/value mismatches persisted ImmutableClusterMetadata."+
			"Continuing with the persisted value as this value cannot be changed once initialized.",
		tag.Key(key),
		tag.IgnoredValue(ignored),
		tag.Value(value))
}