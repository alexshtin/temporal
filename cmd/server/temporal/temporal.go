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
	"log"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/urfave/cli"

	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/log/loggerimpl"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/service/config"
	"go.temporal.io/server/tools/cassandra"
	"go.temporal.io/server/tools/sql"
)

// validServices is the list of all valid temporal services
var validServices = []string{primitives.FrontendService, primitives.HistoryService, primitives.MatchingService, primitives.WorkerService}

func loadConfig(c *cli.Context) *config.Config {
	env := getEnvironment(c)
	zone := getZone(c)
	configDir := getConfigDir(c)

	log.Printf("Loading config; env=%v,zone=%v,configDir=%v\n", env, zone, configDir)

	var cfg config.Config
	err := config.Load(env, configDir, zone, &cfg)
	if err != nil {
		log.Fatal("Config file corrupted.", err)
	}
	if cfg.Log.Level == "debug" {
		log.Printf("config=\n%v\n", cfg.String())
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}
	// cassandra schema version validation
	if err := cassandra.VerifyCompatibleVersion(cfg.Persistence); err != nil {
		log.Fatalf("cassandra schema version compatibility check failed: %v", err)
	}
	// sql schema version validation
	if err := sql.VerifyCompatibleVersion(cfg.Persistence); err != nil {
		log.Fatalf("sql schema version compatibility check failed: %v", err)
	}

	if err := cfg.Global.PProf.NewInitializer(loggerimpl.NewLogger(cfg.Log.NewZapLogger())).Start(); err != nil {
		log.Fatalf("fail to start PProf: %v", err)
	}

	return &cfg
}

func getEnvironment(c *cli.Context) string {
	return strings.TrimSpace(c.GlobalString("env"))
}

func getZone(c *cli.Context) string {
	return strings.TrimSpace(c.GlobalString("zone"))
}

// getServices parses the services arg from cli
// and returns a list of services to start
func getServices(c *cli.Context) []string {

	val := strings.TrimSpace(c.String("services"))
	tokens := strings.Split(val, ",")

	if len(tokens) == 0 {
		log.Fatal("list of services is empty")
	}

	for _, t := range tokens {
		if !isValidService(t) {
			log.Fatalf("invalid service `%v` in service list [%v]", t, val)
		}
	}

	return tokens
}

func isValidService(in string) bool {
	for _, s := range validServices {
		if s == in {
			return true
		}
	}
	return false
}

func getConfigDir(c *cli.Context) string {
	return path.Join(getRootDir(c), c.GlobalString("config"))
}

func getRootDir(c *cli.Context) string {
	dirpath := c.GlobalString("root")
	if len(dirpath) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed, err=%v", err)
		}
		return cwd
	}
	return dirpath
}

// BuildCLI is the main entry point for the temporal server
func BuildCLI() *cli.App {

	app := cli.NewApp()
	app.Name = "temporal"
	app.Usage = "Temporal server"
	app.Version = headers.ServerVersion

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "root, r",
			Value:  ".",
			Usage:  "root directory of execution environment",
			EnvVar: config.EnvKeyRoot,
		},
		cli.StringFlag{
			Name:   "config, c",
			Value:  "config",
			Usage:  "config dir path relative to root",
			EnvVar: config.EnvKeyConfigDir,
		},
		cli.StringFlag{
			Name:   "env, e",
			Value:  "development",
			Usage:  "runtime environment",
			EnvVar: config.EnvKeyEnvironment,
		},
		cli.StringFlag{
			Name:   "zone, az",
			Value:  "",
			Usage:  "availability zone",
			EnvVar: config.EnvKeyAvailabilityZone,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{""},
			Usage:   "start temporal server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "services, s",
					Value: strings.Join(validServices, ","),
					Usage: "list of services to start",
				},
			},
			Action: func(c *cli.Context) error {
				cfg := loadConfig(c)
				services := getServices(c)

				var servers []*server
				sigtermCh := make(chan os.Signal, 1)
				signal.Notify(sigtermCh, os.Interrupt, syscall.SIGTERM)
				for _, svc := range services {
					s := newServer(
						ForService(svc),
						WithConfig(cfg),
					)
					servers = append(servers, s)
					s.Start()
				}

				sigterm := <-sigtermCh
				log.Printf("Received %v signal, initiating shutdown.\n", sigterm)
				for _, s := range servers {
					s.Stop()
				}
				return cli.NewExitError("All services are stopped.", 0)
			},
		},
	}
	return app
}
