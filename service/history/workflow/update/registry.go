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

package update

import (
	"sync"
	"time"

	failurepb "go.temporal.io/api/failure/v1"
	updatepb "go.temporal.io/api/update/v1"
)

type (
	Registry interface {
		Pending(id string) *Update
		Transient() *Update
		Add(update *updatepb.WorkflowUpdate, requestID string, requestTime time.Time) (*Update, RemoveFunc, *Update)
		Len() int

		Clear() // TODO: not sure if this is needed
	}

	RemoveFunc func()

	registryImpl struct {
		sync.RWMutex
		updates map[string]*Update
	}
)

func NewRegistry() *registryImpl {
	return &registryImpl{
		updates: make(map[string]*Update),
	}
}

func (r *registryImpl) Pending(id string) *Update {
	r.RLock()
	defer r.RUnlock()
	return r.updates[id]
}

func (r *registryImpl) Transient() *Update {
	r.RLock()
	defer r.RUnlock()

	return r.transientNoLock()
}

func (r *registryImpl) transientNoLock() *Update {
	// Update can be transient only when it is a single update in the registry.
	if len(r.updates) != 1 {
		return nil
	}

	var trUpdate *Update
	for _, u := range r.updates {
		trUpdate = u
	}
	if !trUpdate.transient {
		return nil
	}
	return trUpdate
}

func (r *registryImpl) Add(update *updatepb.WorkflowUpdate, requestID string, requestTime time.Time) (*Update, RemoveFunc, *Update) {
	// requireInline := false // TODO: should be parameter

	r.Lock()
	defer r.Unlock()
	isTransient := len(r.updates) == 0 // && we want to skip persistence
	u := newUpdate(update, requestID, requestTime, isTransient)
	currentTransient := r.transientNoLock()
	if currentTransient != nil {
		currentTransient.transient = false
		// if requireInline {
		// 	// TODO: block new add until current transient is completed
		// 	// select {
		// 	// case <-r.transient.C():
		// 	//
		// 	// }
		// }
	}

	r.updates[u.id] = u
	return u, func() { r.remove(u.id) }, currentTransient
}

func (r *registryImpl) Len() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.updates)
}

func (r *registryImpl) Clear() {
	r.Lock()
	defer r.Unlock()
	for _, u := range r.updates {
		u.SendResult(nil, r.clearFailure())
	}
	r.updates = make(map[string]*Update)
}

func (r *registryImpl) clearFailure() *failurepb.Failure {
	return &failurepb.Failure{
		Message: "update cleared, please retry",
		FailureInfo: &failurepb.Failure_ServerFailureInfo{
			ServerFailureInfo: &failurepb.ServerFailureInfo{
				NonRetryable: true,
			},
		},
	}
}

func (r *registryImpl) remove(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.updates, id)
}
