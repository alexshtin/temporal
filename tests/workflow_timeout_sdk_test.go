// The MIT License
//
// Copyright (c) 2024 Temporal Technologies Inc.  All rights reserved.
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

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	enumspb "go.temporal.io/api/enums/v1"

	sdkclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"go.temporal.io/server/common/namespace"
	"go.temporal.io/server/common/testing/testvars"
	"go.temporal.io/server/tests/testcore"
)

type WorkflowTimeoutSdkSuite struct {
	testcore.ClientFunctionalSuite
}

func TestWorkflowTimeoutSdkSuite(t *testing.T) {
	s := new(WorkflowTimeoutSdkSuite)
	suite.Run(t, s)
}

func (s *WorkflowTimeoutSdkSuite) TestWorkflowTimeout() {

	tv := testvars.New(s.T()).WithTaskQueue(s.TaskQueue()).WithNamespaceName(namespace.Name(s.Namespace()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	workflowFn := func(ctx workflow.Context) error {
		s.NoError(workflow.Await(ctx, func() bool { return false }))
		return unreachableErr
	}
	s.Worker().RegisterWorkflow(workflowFn)

	firstRun, err := s.SdkClient().ExecuteWorkflow(ctx, sdkclient.StartWorkflowOptions{
		ID:                 tv.WorkflowID(),
		TaskQueue:          tv.TaskQueue().Name,
		WorkflowRunTimeout: 2 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
		},
	}, workflowFn)
	s.NoError(err)

	err = firstRun.GetWithOptions(ctx, nil, sdkclient.WorkflowRunGetOptions{DisableFollowingRuns: true})
	var canErr *workflow.ContinueAsNewError
	s.ErrorAs(err, &canErr)

	secondRun := s.SdkClient().GetWorkflow(ctx, tv.WorkflowID(), "")
	err = secondRun.Get(ctx, nil)
	var weErr *temporal.WorkflowExecutionError
	s.ErrorAs(err, &weErr)
	var timeoutErr *temporal.TimeoutError
	s.ErrorAs(weErr.Unwrap(), &timeoutErr)
	s.Equal(enumspb.TIMEOUT_TYPE_START_TO_CLOSE, timeoutErr.TimeoutType())

	s.NotEmpty(firstRun.GetID())
	s.NotEmpty(secondRun.GetID())
	s.NotEqual(firstRun.GetRunID(), secondRun.GetRunID())

	h1 := s.GetHistory(s.Namespace(), tv.WithRunID(firstRun.GetRunID()).WorkflowExecution())
	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowExecutionTimedOut`, h1)

	h2 := s.GetHistory(s.Namespace(), tv.WithRunID(secondRun.GetRunID()).WorkflowExecution())
	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowExecutionTimedOut`, h2)
}
