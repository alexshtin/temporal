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
	historypb "go.temporal.io/api/history/v1"
	taskqueuepb "go.temporal.io/api/taskqueue/v1"
	updatepb "go.temporal.io/api/update/v1"
	"go.temporal.io/api/workflowservice/v1"

	"go.temporal.io/server/common/payloads"
	"go.temporal.io/server/common/primitives/timestamp"
)

func (s *integrationSuite) TestUpdateWorkflow() {
	id := "integration-signal-workflow-test"
	wt := "integration-signal-workflow-test-type"
	tl := "integration-signal-workflow-test-taskqueue"
	identity := "worker1"
	activityName := "activity_type1"

	workflowType := &commonpb.WorkflowType{Name: wt}

	taskQueue := &taskqueuepb.TaskQueue{Name: tl}

	// header := &commonpb.Header{
	// 	Fields: map[string]*commonpb.Payload{"update header key": payload.EncodeString("update header value")},
	// }

	// Start workflow execution
	request := &workflowservice.StartWorkflowExecutionRequest{
		RequestId:    uuid.New(),
		Namespace:    s.namespace,
		WorkflowId:   id,
		WorkflowType: workflowType,
		TaskQueue:    taskQueue,
		Input:        nil,
		Identity:     identity,
	}

	we, err0 := s.engine.StartWorkflowExecution(NewContext(), request)
	s.NoError(err0)

	// workflow logic
	workflowComplete := false
	activityScheduled := false
	var updateEvent *historypb.HistoryEvent

	wtHandler := func(execution *commonpb.WorkflowExecution, wt *commonpb.WorkflowType, previousStartedEventID, startedEventID int64, history *historypb.History) ([]*commandpb.Command, error) {
		if !activityScheduled {
			activityScheduled = true

			return []*commandpb.Command{{
				CommandType: enumspb.COMMAND_TYPE_SCHEDULE_ACTIVITY_TASK,
				Attributes: &commandpb.Command_ScheduleActivityTaskCommandAttributes{ScheduleActivityTaskCommandAttributes: &commandpb.ScheduleActivityTaskCommandAttributes{
					ActivityId:             strconv.Itoa(1),
					ActivityType:           &commonpb.ActivityType{Name: activityName},
					TaskQueue:              &taskqueuepb.TaskQueue{Name: tl},
					Input:                  nil,
					ScheduleToCloseTimeout: timestamp.DurationPtr(10 * time.Hour),
				}},
			}}, nil
		} else if previousStartedEventID > 0 {
			for _, event := range history.Events[previousStartedEventID:] {
				if event.GetEventType() == enumspb.EVENT_TYPE_WORKFLOW_UPDATE_REQUESTED {
					updateEvent = event
					return []*commandpb.Command{{
						CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_UPDATE,
						Attributes: &commandpb.Command_CompleteWorkflowUpdateCommandAttributes{CompleteWorkflowUpdateCommandAttributes: &commandpb.CompleteWorkflowUpdateCommandAttributes{
							UpdateId:             updateEvent.GetWorkflowUpdateRequestedEventAttributes().GetUpdateId(),
							DurabilityPreference: enumspb.WORKFLOW_UPDATE_DURABILITY_PREFERENCE_BYPASS,
							Result: &commandpb.CompleteWorkflowUpdateCommandAttributes_Success{
								Success: payloads.EncodeString("update success"),
							},
						}},
					}}, nil
				}
			}
		}

		workflowComplete = true
		return []*commandpb.Command{{
			CommandType: enumspb.COMMAND_TYPE_COMPLETE_WORKFLOW_EXECUTION,
			Attributes: &commandpb.Command_CompleteWorkflowExecutionCommandAttributes{CompleteWorkflowExecutionCommandAttributes: &commandpb.CompleteWorkflowExecutionCommandAttributes{
				Result: payloads.EncodeString("workflow result"),
			}},
		}}, nil
	}

	// activity handler
	atHandler := func(execution *commonpb.WorkflowExecution, activityType *commonpb.ActivityType, activityID string, input *commonpb.Payloads, taskToken []byte) (*commonpb.Payloads, bool, error) {
		return payloads.EncodeString("activity result"), false, nil
	}

	poller := &TaskPoller{
		Engine:              s.engine,
		Namespace:           s.namespace,
		TaskQueue:           taskQueue,
		Identity:            identity,
		WorkflowTaskHandler: wtHandler,
		ActivityTaskHandler: atHandler,
		Logger:              s.Logger,
		T:                   s.T(),
	}

	// Make first command to schedule activity
	_, err := poller.PollAndProcessWorkflowTask(false, false)
	s.NoError(err)

	// Send first signal using RunID
	type UpdateResult struct {
		Response *workflowservice.UpdateWorkflowResponse
		Err      error
	}
	updateResultCh := make(chan UpdateResult)
	updateWorkflowFn := func(client workflowservice.WorkflowServiceClient, queryType string) {
		updateResponse, err := s.engine.UpdateWorkflow(NewContext(), &workflowservice.UpdateWorkflowRequest{
			RequestId:         "",
			ResultAccessStyle: enumspb.WORKFLOW_UPDATE_RESULT_ACCESS_STYLE_REQUIRE_INLINE,
			Namespace:         s.namespace,
			WorkflowExecution: &commonpb.WorkflowExecution{
				WorkflowId: id,
				RunId:      we.RunId,
			},
			FirstExecutionRunId: "",
			Update: &updatepb.WorkflowUpdate{
				// Header: header,
				Name: "update_handler",
				Args: payloads.EncodeString("update args"),
			},
		})
		s.NoError(err)
		updateResultCh <- UpdateResult{Response: updateResponse, Err: err}
	}

	go updateWorkflowFn(s.engine, "RunID")

	// Process update in workflow
	_, err = poller.PollAndProcessWorkflowTask(true, false)
	s.NoError(err)
	updateResult := <-updateResultCh

	s.False(workflowComplete)
	s.True(updateEvent != nil)
	s.Equal("update_handler", updateEvent.GetWorkflowUpdateRequestedEventAttributes().GetUpdate().GetName())
	s.EqualValues(payloads.EncodeString("update args"), updateEvent.GetWorkflowUpdateRequestedEventAttributes().GetUpdate().GetArgs())
	// s.Equal(header, updateEvent.GetWorkflowUpdateRequestedEventAttributes().GetHeader())
	s.EqualValues(payloads.EncodeString("update success"), updateResult.Response.GetSuccess())

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
}
