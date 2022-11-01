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

package host

import (
	"strconv"
	"time"

	"github.com/pborman/uuid"
	commandpb "go.temporal.io/api/command/v1"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	failurepb "go.temporal.io/api/failure/v1"
	historypb "go.temporal.io/api/history/v1"
	interactionpb "go.temporal.io/api/interaction/v1"
	taskqueuepb "go.temporal.io/api/taskqueue/v1"
	"go.temporal.io/api/workflowservice/v1"

	"go.temporal.io/server/common/payloads"
	"go.temporal.io/server/common/primitives/timestamp"
)

func (s *integrationSuite) TestUpdateWorkflow_ExistingWorkflowTask_AcceptComplete() {
	id := "integration-update-workflow-test-1"
	wt := "integration-update-workflow-test-1-type"
	tq := "integration-update-workflow-test-1-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:    uuid.New(),
		Namespace:    s.namespace,
		WorkflowId:   id,
		WorkflowType: workflowType,
		TaskQueue:    taskQueue,
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted`, history)
			// interaction handler will add commands.
			return nil, nil
		case 2:
			s.EqualHistory(`
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowUpdateCompleted
  7 WorkflowTaskScheduled
  8 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		interaction := task.Interactions[0]

		s.Equal(payloads.EncodeString("update args"), interaction.GetInput().GetArgs())
		s.Equal("update_handler", interaction.GetInput().GetName())

		return []*commandpb.Command{
			{
				CommandType: enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_AcceptWorkflowUpdateCommandAttributes{AcceptWorkflowUpdateCommandAttributes: &commandpb.AcceptWorkflowUpdateCommandAttributes{
					Meta:  interaction.GetMeta(),
					Input: interaction.GetInput(),
				}},
			},
			{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
					Meta: interaction.GetMeta(),
					Output: &interactionpb.Output{
						Result: &interactionpb.Output_Success{
							Success: payloads.EncodeString("update success"),
						}}}},
			},
		}, nil
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateResultCh := make(chan UpdateResult)
	updateWorkflowFn := func() {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString("update args"),
			},
		})
		updateResultCh <- UpdateResult{Response: updateResponse, Err: err1}
	}
	go updateWorkflowFn()
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Process update in workflow.
	_, updateResp, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	updateResult := <-updateResultCh
	s.NoError(updateResult.Err)
	s.EqualValues(payloads.EncodeString("update success"), updateResult.Response.GetOutput().GetSuccess())
	s.EqualValues(3, updateResp.ResetHistoryEventId)

	// Complete workflow.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(8, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(2, wtHandlerCalls)
	s.Equal(1, interHandlerCalls)

	events := s.getHistory(s.namespace, we)

	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowUpdateCompleted
  7 WorkflowTaskScheduled
  8 WorkflowTaskStarted
  9 WorkflowTaskCompleted
 10 WorkflowExecutionCompleted`, events)
}

func (s *integrationSuite) TestUpdateWorkflow_NewWorkflowTask_AcceptComplete() {
	id := "integration-update-workflow-test-2"
	wt := "integration-update-workflow-test-2-type"
	tq := "integration-update-workflow-test-2-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:    uuid.New(),
		Namespace:    s.namespace,
		WorkflowId:   id,
		WorkflowType: workflowType,
		TaskQueue:    taskQueue,
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_SCHEDULE_ACTIVITY_TASK,
				Attributes: &commandpb.Command_ScheduleActivityTaskCommandAttributes{ScheduleActivityTaskCommandAttributes: &commandpb.ScheduleActivityTaskCommandAttributes{
					ActivityId:             strconv.Itoa(1),
					ActivityType:           &commonpb.ActivityType{Name: "activity_type_1"},
					TaskQueue:              &taskqueuepb.TaskQueue{Name: tq},
					ScheduleToCloseTimeout: timestamp.DurationPtr(10 * time.Hour),
				}},
			}}, nil
		case 2:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 ActivityTaskScheduled
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted`, history)
			// interaction handler will add commands.
			return nil, nil
		case 3:
			s.EqualHistory(`
  8 WorkflowTaskCompleted
  9 WorkflowUpdateAccepted
 10 WorkflowUpdateCompleted
 11 WorkflowTaskScheduled
 12 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		interaction := task.Interactions[0]

		s.Equal(payloads.EncodeString("update args"), interaction.GetInput().GetArgs())
		s.Equal("update_handler", interaction.GetInput().GetName())

		return []*commandpb.Command{
			{
				CommandType: enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_AcceptWorkflowUpdateCommandAttributes{AcceptWorkflowUpdateCommandAttributes: &commandpb.AcceptWorkflowUpdateCommandAttributes{
					Meta:  interaction.GetMeta(),
					Input: interaction.GetInput(),
				}},
			},
			{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
					Meta: interaction.GetMeta(),
					Output: &interactionpb.Output{
						Result: &interactionpb.Output_Success{
							Success: payloads.EncodeString("update success"),
						}}}},
			},
		}, nil
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	// Start activity using existing workflow task.
	_, err = poller.PollAndProcessWorkflowTask(true, false)
	s.NoError(err)

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateResultCh := make(chan UpdateResult)
	updateWorkflowFn := func() {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString("update args"),
			},
		})
		updateResultCh <- UpdateResult{Response: updateResponse, Err: err1}
	}
	go updateWorkflowFn()
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Process update in workflow.
	_, updateResp, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	updateResult := <-updateResultCh
	s.NoError(updateResult.Err)
	s.EqualValues(payloads.EncodeString("update success"), updateResult.Response.GetOutput().GetSuccess())
	s.EqualValues(7, updateResp.ResetHistoryEventId)

	// Complete workflow.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(12, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(3, wtHandlerCalls)
	s.Equal(1, interHandlerCalls)

	events := s.getHistory(s.namespace, we)

	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 ActivityTaskScheduled
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted
  8 WorkflowTaskCompleted
  9 WorkflowUpdateAccepted
 10 WorkflowUpdateCompleted
 11 WorkflowTaskScheduled
 12 WorkflowTaskStarted
 13 WorkflowTaskCompleted
 14 WorkflowExecutionCompleted`, events)
}

func (s *integrationSuite) TestUpdateWorkflow_ExistingWorkflowTask_Reject() {
	id := "integration-update-workflow-test-3"
	wt := "integration-update-workflow-test-3-type"
	tq := "integration-update-workflow-test-3-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:           uuid.New(),
		Namespace:           s.namespace,
		WorkflowId:          id,
		WorkflowType:        workflowType,
		TaskQueue:           taskQueue,
		WorkflowTaskTimeout: timestamp.DurationPtr(10 * time.Minute),
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted`, history)
			// interaction handler will add commands.
			return nil, nil
		case 2:
			s.EqualHistory(`
  4 WorkflowTaskCompleted
  5 WorkflowTaskScheduled
  6 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		interaction := task.Interactions[0]

		s.Equal(payloads.EncodeString("update args"), interaction.GetInput().GetArgs())
		s.Equal("update_handler", interaction.GetInput().GetName())

		return []*commandpb.Command{
			{
				CommandType: enumspb.COMMAND_TYPE_REJECT_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_RejectWorkflowUpdateCommandAttributes{RejectWorkflowUpdateCommandAttributes: &commandpb.RejectWorkflowUpdateCommandAttributes{
					Meta: interaction.GetMeta(),
					Failure: &failurepb.Failure{
						Message: "update rejected",
						// TODO: use specific failure type???
						FailureInfo: &failurepb.Failure_ApplicationFailureInfo{ApplicationFailureInfo: &failurepb.ApplicationFailureInfo{}},
					}},
				},
			}}, nil
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateResultCh := make(chan UpdateResult)
	updateWorkflowFn := func() {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString("update args"),
			},
		})
		updateResultCh <- UpdateResult{Response: updateResponse, Err: err1}
	}
	go updateWorkflowFn()
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Process update in workflow.
	_, updateResp, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	updateResult := <-updateResultCh
	s.NoError(updateResult.Err)
	s.Equal("update rejected", updateResult.Response.GetOutput().GetFailure().GetMessage())
	s.EqualValues(3, updateResp.ResetHistoryEventId)

	// Complete workflow.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(6, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(2, wtHandlerCalls)
	s.Equal(1, interHandlerCalls)

	events := s.getHistory(s.namespace, we)

	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowTaskScheduled
  6 WorkflowTaskStarted
  7 WorkflowTaskCompleted
  8 WorkflowExecutionCompleted`, events)
}

func (s *integrationSuite) TestUpdateWorkflow_NewWorkflowTask_Reject() {
	id := "integration-update-workflow-test-4"
	wt := "integration-update-workflow-test-4-type"
	tq := "integration-update-workflow-test-4-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:           uuid.New(),
		Namespace:           s.namespace,
		WorkflowId:          id,
		WorkflowType:        workflowType,
		TaskQueue:           taskQueue,
		WorkflowTaskTimeout: timestamp.DurationPtr(10 * time.Minute),
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_SCHEDULE_ACTIVITY_TASK,
				Attributes: &commandpb.Command_ScheduleActivityTaskCommandAttributes{ScheduleActivityTaskCommandAttributes: &commandpb.ScheduleActivityTaskCommandAttributes{
					ActivityId:             strconv.Itoa(1),
					ActivityType:           &commonpb.ActivityType{Name: "activity_type_1"},
					TaskQueue:              &taskqueuepb.TaskQueue{Name: tq},
					ScheduleToCloseTimeout: timestamp.DurationPtr(10 * time.Hour),
				}},
			}}, nil
		case 2:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 ActivityTaskScheduled
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted`, history)
			// interaction handler will add commands.
			return nil, nil
		case 3:
			s.EqualHistory(`
  4 WorkflowTaskCompleted
  5 ActivityTaskScheduled
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		interaction := task.Interactions[0]

		s.Equal(payloads.EncodeString("update args"), interaction.GetInput().GetArgs())
		s.Equal("update_handler", interaction.GetInput().GetName())

		return []*commandpb.Command{
			{
				CommandType: enumspb.COMMAND_TYPE_REJECT_WORKFLOW_UPDATE,
				Attributes: &commandpb.Command_RejectWorkflowUpdateCommandAttributes{RejectWorkflowUpdateCommandAttributes: &commandpb.RejectWorkflowUpdateCommandAttributes{
					Meta: interaction.GetMeta(),
					Failure: &failurepb.Failure{
						Message: "update rejected",
						// TODO: use specific failure type???
						FailureInfo: &failurepb.Failure_ApplicationFailureInfo{ApplicationFailureInfo: &failurepb.ApplicationFailureInfo{}},
					}},
				},
			}}, nil
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	// Start activity using existing workflow task.
	_, err = poller.PollAndProcessWorkflowTask(true, false)
	s.NoError(err)

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateResultCh := make(chan UpdateResult)
	updateWorkflowFn := func() {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString("update args"),
			},
		})
		updateResultCh <- UpdateResult{Response: updateResponse, Err: err1}
	}
	go updateWorkflowFn()
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Process update in workflow.
	_, updateResp, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	updateResult := <-updateResultCh
	s.NoError(updateResult.Err)
	s.Equal("update rejected", updateResult.Response.GetOutput().GetFailure().GetMessage())
	s.EqualValues(3, updateResp.ResetHistoryEventId)

	// Complete workflow.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(7, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(3, wtHandlerCalls)
	s.Equal(1, interHandlerCalls)

	events := s.getHistory(s.namespace, we)
	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 ActivityTaskScheduled
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted
  8 WorkflowTaskCompleted
  9 WorkflowExecutionCompleted`, events)
}

func (s *integrationSuite) TestUpdateWorkflow_ExistingWorkflowTask_1stAccept_2ndAccept_2ndComplete_1stComplete() {
	id := "integration-update-workflow-test-5"
	wt := "integration-update-workflow-test-5-type"
	tq := "integration-update-workflow-test-5-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:           uuid.New(),
		Namespace:           s.namespace,
		WorkflowId:          id,
		WorkflowType:        workflowType,
		TaskQueue:           taskQueue,
		WorkflowTaskTimeout: timestamp.DurationPtr(10 * time.Minute),
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	var interaction1Meta, interaction2Meta *interactionpb.Meta

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted`, history)
			// interaction handler accepts 1st update.
			return nil, nil
		case 2:
			// Complete WT with empty command list.
			s.EqualHistory(`
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted`, history)
			return nil, nil
		case 3:
			s.EqualHistory(`
  8 WorkflowTaskCompleted
  9 WorkflowTaskScheduled
 10 WorkflowTaskStarted`, history)
			// interaction handler accepts 2nd update.
			return nil, nil
		case 4:
			s.EqualHistory(`
 11 WorkflowTaskCompleted
 12 WorkflowUpdateAccepted
 13 WorkflowTaskScheduled
 14 WorkflowTaskStarted`, history)
			// complete 2nd update.
			s.NotNil(interaction2Meta)
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
						Meta: interaction2Meta,
						Output: &interactionpb.Output{
							Result: &interactionpb.Output_Success{
								Success: payloads.EncodeString("update2 success"),
							}}}},
				},
			}, nil
		case 5:
			s.EqualHistory(`
 15 WorkflowTaskCompleted
 16 WorkflowUpdateCompleted
 17 WorkflowTaskScheduled
 18 WorkflowTaskStarted`, history)
			s.NotNil(interaction1Meta)
			// complete 1st update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
						Meta: interaction1Meta,
						Output: &interactionpb.Output{
							Result: &interactionpb.Output_Success{
								Success: payloads.EncodeString("update1 success"),
							}}}},
				},
			}, nil

		case 6:
			s.EqualHistory(`
 19 WorkflowTaskCompleted
 20 WorkflowUpdateCompleted
 21 WorkflowTaskScheduled
 22 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		switch interHandlerCalls {
		case 1:
			s.Len(task.Interactions, 1)
			interaction := task.Interactions[0]
			interaction1Meta = interaction.GetMeta()
			// accept 1st update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_AcceptWorkflowUpdateCommandAttributes{AcceptWorkflowUpdateCommandAttributes: &commandpb.AcceptWorkflowUpdateCommandAttributes{
						Meta:  interaction.GetMeta(),
						Input: interaction.GetInput(),
					}},
				},
			}, nil
		case 2:
			s.Len(task.Interactions, 1)
			interaction := task.Interactions[0]
			interaction2Meta = interaction.GetMeta()
			// accept 2nd update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_AcceptWorkflowUpdateCommandAttributes{AcceptWorkflowUpdateCommandAttributes: &commandpb.AcceptWorkflowUpdateCommandAttributes{
						Meta:  interaction.GetMeta(),
						Input: interaction.GetInput(),
					}},
				},
			}, nil
		default:
			s.Fail("interHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateWorkflowFn := func(ch chan<- UpdateResult, updateArgs string) {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString(updateArgs),
			},
		})
		ch <- UpdateResult{Response: updateResponse, Err: err1}
	}

	updateResultCh1 := make(chan UpdateResult)
	go updateWorkflowFn(updateResultCh1, "update1 args")
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Accept update1 in WT1.
	_, updateAcceptResp1, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	s.NotNil(updateAcceptResp1)
	s.EqualValues(3, updateAcceptResp1.ResetHistoryEventId)

	// Start update after WT2 has started.
	updateResultCh2 := make(chan UpdateResult)
	go updateWorkflowFn(updateResultCh2, "update2 args")
	time.Sleep(500 * time.Millisecond) // This is to make sure that update2 gets to the server before WT2 is completed and WT3 is started.

	// WT2 from updateAcceptResp1.GetWorkflowTask() doesn't have 2nd update because WT2 has started before update was received.
	// Complete WT2 with empty commands list.
	completeEmptyWTResp, err := poller.HandlePartialWorkflowTask(updateAcceptResp1.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeEmptyWTResp)
	s.EqualValues(7, completeEmptyWTResp.ResetHistoryEventId)

	// WT3 from completeEmptyWTResp.GetWorkflowTask() has 2nd update.
	// Accept update2 in WT3.
	updateAcceptResp2, err := poller.HandlePartialWorkflowTask(completeEmptyWTResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(updateAcceptResp2)
	s.EqualValues(10, updateAcceptResp2.ResetHistoryEventId)

	// Complete update2 in WT4.
	updateCompleteResp2, err := poller.HandlePartialWorkflowTask(updateAcceptResp2.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(updateCompleteResp2)
	updateResult2 := <-updateResultCh2
	s.NoError(updateResult2.Err)
	s.EqualValues(payloads.EncodeString("update2 success"), updateResult2.Response.GetOutput().GetSuccess())
	s.EqualValues(14, updateCompleteResp2.ResetHistoryEventId)

	// Complete update1 in WT5.
	updateCompleteResp1, err := poller.HandlePartialWorkflowTask(updateCompleteResp2.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(updateCompleteResp1)
	updateResult1 := <-updateResultCh1
	s.NoError(updateResult1.Err)
	s.EqualValues(payloads.EncodeString("update1 success"), updateResult1.Response.GetOutput().GetSuccess())
	s.EqualValues(18, updateCompleteResp1.ResetHistoryEventId)

	// Complete WT6.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateCompleteResp1.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(22, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(6, wtHandlerCalls)
	s.Equal(2, interHandlerCalls)

	events := s.getHistory(s.namespace, we)

	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted
  8 WorkflowTaskCompleted
  9 WorkflowTaskScheduled
 10 WorkflowTaskStarted
 11 WorkflowTaskCompleted
 12 WorkflowUpdateAccepted
 13 WorkflowTaskScheduled
 14 WorkflowTaskStarted
 15 WorkflowTaskCompleted
 16 WorkflowUpdateCompleted
 17 WorkflowTaskScheduled
 18 WorkflowTaskStarted
 19 WorkflowTaskCompleted
 20 WorkflowUpdateCompleted
 21 WorkflowTaskScheduled
 22 WorkflowTaskStarted
 23 WorkflowTaskCompleted
 24 WorkflowExecutionCompleted`, events)
}

func (s *integrationSuite) TestUpdateWorkflow_ExistingWorkflowTask_1stAccept_2ndReject_1stComplete() {
	id := "integration-update-workflow-test-6"
	wt := "integration-update-workflow-test-6-type"
	tq := "integration-update-workflow-test-6-task-queue"

	workflowType := &commonpb.WorkflowType{Name: wt}
	taskQueue := &taskqueuepb.TaskQueue{Name: tq}

	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:           uuid.New(),
		Namespace:           s.namespace,
		WorkflowId:          id,
		WorkflowType:        workflowType,
		TaskQueue:           taskQueue,
		WorkflowTaskTimeout: timestamp.DurationPtr(10 * time.Minute),
	}

	startResp, err := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err)

	we := &commonpb.WorkflowExecution{
		WorkflowId: id,
		RunId:      startResp.GetRunId(),
	}

	var interaction1Meta *interactionpb.Meta

	wtHandlerCalls := 0
	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		wtHandlerCalls++
		switch wtHandlerCalls {
		case 1:
			s.EqualHistory(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted`, history)
			// interaction handler accepts 1st update.
			return nil, nil
		case 2:
			// Complete WT with empty command list.
			s.EqualHistory(`
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted`, history)
			// Return non-empty command list to indicate that WT is not heartbeat.
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_SCHEDULE_ACTIVITY_TASK,
				Attributes: &commandpb.Command_ScheduleActivityTaskCommandAttributes{ScheduleActivityTaskCommandAttributes: &commandpb.ScheduleActivityTaskCommandAttributes{
					ActivityId:             strconv.Itoa(1),
					ActivityType:           &commonpb.ActivityType{Name: "activity_type_1"},
					TaskQueue:              &taskqueuepb.TaskQueue{Name: tq},
					ScheduleToCloseTimeout: timestamp.DurationPtr(10 * time.Hour),
				}},
			}}, nil
		case 3:
			s.EqualHistory(`
  8 WorkflowTaskCompleted
  9 ActivityTaskScheduled
 10 WorkflowTaskScheduled
 11 WorkflowTaskStarted`, history)
			// interaction handler rejects 2nd update.
			return nil, nil
		case 4:
			// WT3 events were transient.
			s.EqualHistory(`
  8 WorkflowTaskCompleted
  9 ActivityTaskScheduled
 10 WorkflowTaskScheduled
 11 WorkflowTaskStarted`, history)
			s.NotNil(interaction1Meta)
			// complete 1st update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
						Meta: interaction1Meta,
						Output: &interactionpb.Output{
							Result: &interactionpb.Output_Success{
								Success: payloads.EncodeString("update1 success"),
							}}}},
				},
			}, nil

		case 5:
			s.EqualHistory(`
 12 WorkflowTaskCompleted
 13 WorkflowUpdateCompleted
 14 WorkflowTaskScheduled
 15 WorkflowTaskStarted`, history)
			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
				Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
					Result: payloads.EncodeString("done"),
				}},
			}}, nil
		default:
			s.Fail("wtHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	interHandlerCalls := 0
	interHandler := func(task *workflowservice.PollWorkflowTaskQueueResponse) ([]*commandpb.Command, error) {
		interHandlerCalls++
		switch interHandlerCalls {
		case 1:
			s.Len(task.Interactions, 1)
			interaction := task.Interactions[0]
			interaction1Meta = interaction.GetMeta()
			// accept 1st update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_ACCEPT_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_AcceptWorkflowUpdateCommandAttributes{AcceptWorkflowUpdateCommandAttributes: &commandpb.AcceptWorkflowUpdateCommandAttributes{
						Meta:  interaction.GetMeta(),
						Input: interaction.GetInput(),
					}},
				},
			}, nil
		case 2:
			s.Len(task.Interactions, 1)
			interaction := task.Interactions[0]
			// reject 2nd update.
			return []*commandpb.Command{
				{
					CommandType: enumspb.COMMAND_TYPE_REJECT_WORKFLOW_UPDATE,
					Attributes: &commandpb.Command_RejectWorkflowUpdateCommandAttributes{RejectWorkflowUpdateCommandAttributes: &commandpb.RejectWorkflowUpdateCommandAttributes{
						Meta: interaction.GetMeta(),
						Failure: &failurepb.Failure{
							Message: "update2 rejected",
							// TODO: use specific failure type???
							FailureInfo: &failurepb.Failure_ApplicationFailureInfo{ApplicationFailureInfo: &failurepb.ApplicationFailureInfo{}},
						}}},
				},
			}, nil
		default:
			s.Fail("interHandler shouldn't be called %d times", wtHandlerCalls)
			return nil, nil
		}
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		WorkflowTaskHandler: wtHandler,
		InteractionHandler:  interHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateWorkflowFn := func(ch chan<- UpdateResult, updateArgs string) {
		updateResponse, err1 := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: we,
			Input: &interactionpb.Input{
				Name: "update_handler",
				Args: payloads.EncodeString(updateArgs),
			},
		})
		ch <- UpdateResult{Response: updateResponse, Err: err1}
	}

	updateResultCh1 := make(chan UpdateResult)
	go updateWorkflowFn(updateResultCh1, "update1 args")
	time.Sleep(500 * time.Millisecond) // This is to make sure that update gets to the server.

	// Accept update1 in WT1.
	_, updateAcceptResp1, err := poller.PollAndProcessWorkflowTaskWithAttemptAndRetryAndForceNewWorkflowTask(false, false, false, false, 1, 5, true, nil)
	s.NoError(err)
	s.NotNil(updateAcceptResp1)
	s.EqualValues(3, updateAcceptResp1.ResetHistoryEventId)

	// Start update after WT2 has started.
	updateResultCh2 := make(chan UpdateResult)
	go updateWorkflowFn(updateResultCh2, "update2 args")
	time.Sleep(500 * time.Millisecond) // This is to make sure that update2 gets to the server before WT2 is completed and WT3 is started.

	// WT2 from updateAcceptResp1.GetWorkflowTask() doesn't have 2nd update because WT2 has started before update was received.
	// Complete WT2 with empty commands list.
	completeEmptyWTResp, err := poller.HandlePartialWorkflowTask(updateAcceptResp1.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeEmptyWTResp)
	s.EqualValues(7, completeEmptyWTResp.ResetHistoryEventId)

	// WT3 from completeEmptyWTResp.GetWorkflowTask() has 2nd update.
	// Reject update2 in WT3.
	updateRejectResp2, err := poller.HandlePartialWorkflowTask(completeEmptyWTResp.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(updateRejectResp2)
	updateResult2 := <-updateResultCh2
	s.NoError(updateResult2.Err)
	s.Equal("update2 rejected", updateResult2.Response.GetOutput().GetFailure().GetMessage())
	s.EqualValues(7, updateRejectResp2.ResetHistoryEventId)

	// Complete update1 in WT4.
	updateCompleteResp1, err := poller.HandlePartialWorkflowTask(updateRejectResp2.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(updateCompleteResp1)
	updateResult1 := <-updateResultCh1
	s.NoError(updateResult1.Err)
	s.EqualValues(payloads.EncodeString("update1 success"), updateResult1.Response.GetOutput().GetSuccess())
	s.EqualValues(11, updateCompleteResp1.ResetHistoryEventId)

	// Complete WT5.
	completeWorkflowResp, err := poller.HandlePartialWorkflowTask(updateCompleteResp1.GetWorkflowTask())
	s.NoError(err)
	s.NotNil(completeWorkflowResp)
	s.Nil(completeWorkflowResp.GetWorkflowTask())
	s.EqualValues(15, completeWorkflowResp.ResetHistoryEventId)

	s.Equal(5, wtHandlerCalls)
	s.Equal(2, interHandlerCalls)

	events := s.getHistory(s.namespace, we)

	s.EqualHistoryEvents(`
  1 WorkflowExecutionStarted
  2 WorkflowTaskScheduled
  3 WorkflowTaskStarted
  4 WorkflowTaskCompleted
  5 WorkflowUpdateAccepted
  6 WorkflowTaskScheduled
  7 WorkflowTaskStarted
  8 WorkflowTaskCompleted
  9 ActivityTaskScheduled
 10 WorkflowTaskScheduled
 11 WorkflowTaskStarted
 12 WorkflowTaskCompleted
 13 WorkflowUpdateCompleted
 14 WorkflowTaskScheduled
 15 WorkflowTaskStarted
 16 WorkflowTaskCompleted
 17 WorkflowExecutionCompleted`, events)
}

// // Send another signal without RunID
// signalName := "another signal"
// signalInput := payloads.EncodeString("another signal input")
// _, err = s.engine.SignalWorkflowExecution(NewContext(), &workflowservice.SignalWorkflowExecutionRequest{
// 	Namespace: s.namespace,
// 	WorkflowExecution: &commonpb.WorkflowExecution{
// 		WorkflowId: id,
// 	},
// 	SignalName: signalName,
// 	Input:      signalInput,
// 	Identity:   identity,
// })
// s.NoError(err)
//
// // Process signal in workflow
// _, err = poller.PollAndProcessWorkflowTask(true, false)
// s.Logger.Info("PollAndProcessWorkflowTask", tag.Error(err))
// s.NoError(err)
//
// s.False(workflowComplete)
// s.True(updateEvent != nil)
// s.Equal(signalName, updateEvent.GetWorkflowExecutionSignaledEventAttributes().SignalName)
// s.Equal(signalInput, updateEvent.GetWorkflowExecutionSignaledEventAttributes().Input)
// s.Equal(identity, updateEvent.GetWorkflowExecutionSignaledEventAttributes().Identity)
