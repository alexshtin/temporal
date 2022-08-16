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
	"time"

	"github.com/pborman/uuid"
	commonpb "go.temporal.io/api/common/v1"
	failurepb "go.temporal.io/api/failure/v1"
	updatepb "go.temporal.io/api/update/v1"

	"go.temporal.io/server/common"
)

type (
	Update struct {
		id                           string
		requestID                    string
		requestTime                  time.Time
		transient                    bool
		scheduledWorkflowTaskEventID int64
		update                       *updatepb.WorkflowUpdate
		resultCh                     chan *Result
		unblockCh                    chan struct{}
	}

	Result struct {
		// one of:
		Success *commonpb.Payloads
		Failure *failurepb.Failure
	}
)

func newUpdate(update *updatepb.WorkflowUpdate, requestID string, requestTime time.Time, isTransient bool) *Update {
	return &Update{
		id:                           uuid.New(),
		requestID:                    requestID,
		requestTime:                  requestTime,
		transient:                    isTransient,
		scheduledWorkflowTaskEventID: common.EmptyEventID,
		update:                       update,
		resultCh:                     make(chan *Result),
		unblockCh:                    make(chan struct{}),
	}
}

func (u *Update) Update() *updatepb.WorkflowUpdate {
	return u.update
}

func (u *Update) ResultCh() <-chan *Result {
	return u.resultCh
}

func (u *Update) UnblockCh() <-chan struct{} {
	return u.unblockCh
}

func (u *Update) SendResult(success *commonpb.Payloads, failure *failurepb.Failure) {
	// TODO: validate that one and only one is not nil
	r := &Result{
		Success: success,
		Failure: failure,
	}

	u.resultCh <- r
}

func (u *Update) ID() string {
	return u.id
}

func (u *Update) RequestID() string {
	return u.requestID
}

func (u *Update) RequestTime() *time.Time {
	return &u.requestTime
}

func (u *Update) Transient() bool {
	return u.transient
}

func (u *Update) SetScheduledWorkflowTaskEventID(scheduledWorkflowTaskEventID int64) {
	u.scheduledWorkflowTaskEventID = scheduledWorkflowTaskEventID
}
func (u *Update) ScheduledWorkflowTaskEventID() int64 {
	return u.scheduledWorkflowTaskEventID
}
