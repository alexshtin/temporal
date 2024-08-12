// The MIT License
//
// Copyright (c) 2020 Temporal Technologies, Inc.
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// plugins:
// 	protoc-gen-go
// 	protoc
// source: temporal/server/api/enums/v1/replication.proto

package enums

import (
	reflect "reflect"
	"strconv"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReplicationTaskType int32

const (
	REPLICATION_TASK_TYPE_UNSPECIFIED                      ReplicationTaskType = 0
	REPLICATION_TASK_TYPE_NAMESPACE_TASK                   ReplicationTaskType = 1
	REPLICATION_TASK_TYPE_HISTORY_TASK                     ReplicationTaskType = 2
	REPLICATION_TASK_TYPE_SYNC_SHARD_STATUS_TASK           ReplicationTaskType = 3
	REPLICATION_TASK_TYPE_SYNC_ACTIVITY_TASK               ReplicationTaskType = 4
	REPLICATION_TASK_TYPE_HISTORY_METADATA_TASK            ReplicationTaskType = 5
	REPLICATION_TASK_TYPE_HISTORY_V2_TASK                  ReplicationTaskType = 6
	REPLICATION_TASK_TYPE_SYNC_WORKFLOW_STATE_TASK         ReplicationTaskType = 7
	REPLICATION_TASK_TYPE_TASK_QUEUE_USER_DATA             ReplicationTaskType = 8
	REPLICATION_TASK_TYPE_SYNC_HSM_TASK                    ReplicationTaskType = 9
	REPLICATION_TASK_TYPE_BACKFILL_HISTORY_TASK            ReplicationTaskType = 10
	REPLICATION_TASK_TYPE_VERIFY_VERSIONED_TRANSITION_TASK ReplicationTaskType = 11
)

// Enum value maps for ReplicationTaskType.
var (
	ReplicationTaskType_name = map[int32]string{
		0:  "REPLICATION_TASK_TYPE_UNSPECIFIED",
		1:  "REPLICATION_TASK_TYPE_NAMESPACE_TASK",
		2:  "REPLICATION_TASK_TYPE_HISTORY_TASK",
		3:  "REPLICATION_TASK_TYPE_SYNC_SHARD_STATUS_TASK",
		4:  "REPLICATION_TASK_TYPE_SYNC_ACTIVITY_TASK",
		5:  "REPLICATION_TASK_TYPE_HISTORY_METADATA_TASK",
		6:  "REPLICATION_TASK_TYPE_HISTORY_V2_TASK",
		7:  "REPLICATION_TASK_TYPE_SYNC_WORKFLOW_STATE_TASK",
		8:  "REPLICATION_TASK_TYPE_TASK_QUEUE_USER_DATA",
		9:  "REPLICATION_TASK_TYPE_SYNC_HSM_TASK",
		10: "REPLICATION_TASK_TYPE_BACKFILL_HISTORY_TASK",
		11: "REPLICATION_TASK_TYPE_VERIFY_VERSIONED_TRANSITION_TASK",
	}
	ReplicationTaskType_value = map[string]int32{
		"REPLICATION_TASK_TYPE_UNSPECIFIED":                      0,
		"REPLICATION_TASK_TYPE_NAMESPACE_TASK":                   1,
		"REPLICATION_TASK_TYPE_HISTORY_TASK":                     2,
		"REPLICATION_TASK_TYPE_SYNC_SHARD_STATUS_TASK":           3,
		"REPLICATION_TASK_TYPE_SYNC_ACTIVITY_TASK":               4,
		"REPLICATION_TASK_TYPE_HISTORY_METADATA_TASK":            5,
		"REPLICATION_TASK_TYPE_HISTORY_V2_TASK":                  6,
		"REPLICATION_TASK_TYPE_SYNC_WORKFLOW_STATE_TASK":         7,
		"REPLICATION_TASK_TYPE_TASK_QUEUE_USER_DATA":             8,
		"REPLICATION_TASK_TYPE_SYNC_HSM_TASK":                    9,
		"REPLICATION_TASK_TYPE_BACKFILL_HISTORY_TASK":            10,
		"REPLICATION_TASK_TYPE_VERIFY_VERSIONED_TRANSITION_TASK": 11,
	}
)

func (x ReplicationTaskType) Enum() *ReplicationTaskType {
	p := new(ReplicationTaskType)
	*p = x
	return p
}

func (x ReplicationTaskType) String() string {
	switch x {
	case REPLICATION_TASK_TYPE_UNSPECIFIED:
		return "Unspecified"
	case REPLICATION_TASK_TYPE_NAMESPACE_TASK:
		return "NamespaceTask"
	case REPLICATION_TASK_TYPE_HISTORY_TASK:
		return "HistoryTask"
	case REPLICATION_TASK_TYPE_SYNC_SHARD_STATUS_TASK:
		return "SyncShardStatusTask"
	case REPLICATION_TASK_TYPE_SYNC_ACTIVITY_TASK:
		return "SyncActivityTask"
	case REPLICATION_TASK_TYPE_HISTORY_METADATA_TASK:
		return "HistoryMetadataTask"
	case REPLICATION_TASK_TYPE_HISTORY_V2_TASK:
		return "HistoryV2Task"
	case REPLICATION_TASK_TYPE_SYNC_WORKFLOW_STATE_TASK:

		// Deprecated: Use ReplicationTaskType.Descriptor instead.
		return "SyncWorkflowStateTask"
	case REPLICATION_TASK_TYPE_TASK_QUEUE_USER_DATA:
		return "TaskQueueUserData"
	case REPLICATION_TASK_TYPE_SYNC_HSM_TASK:
		return "SyncHsmTask"
	case REPLICATION_TASK_TYPE_BACKFILL_HISTORY_TASK:
		return "BackfillHistoryTask"
	case REPLICATION_TASK_TYPE_VERIFY_VERSIONED_TRANSITION_TASK:
		return "VerifyVersionedTransitionTask"
	default:
		return strconv.Itoa(int(x))
	}

}

func (ReplicationTaskType) Descriptor() protoreflect.EnumDescriptor {
	return file_temporal_server_api_enums_v1_replication_proto_enumTypes[0].Descriptor()
}

func (ReplicationTaskType) Type() protoreflect.EnumType {
	return &file_temporal_server_api_enums_v1_replication_proto_enumTypes[0]
}

func (x ReplicationTaskType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

func (ReplicationTaskType) EnumDescriptor() ([]byte, []int) {
	return file_temporal_server_api_enums_v1_replication_proto_rawDescGZIP(), []int{0}
}

type NamespaceOperation int32

const (
	NAMESPACE_OPERATION_UNSPECIFIED NamespaceOperation = 0
	NAMESPACE_OPERATION_CREATE      NamespaceOperation = 1
	NAMESPACE_OPERATION_UPDATE      NamespaceOperation = 2
)

// Enum value maps for NamespaceOperation.
var (
	NamespaceOperation_name = map[int32]string{
		0: "NAMESPACE_OPERATION_UNSPECIFIED",
		1: "NAMESPACE_OPERATION_CREATE",
		2: "NAMESPACE_OPERATION_UPDATE",
	}
	NamespaceOperation_value = map[string]int32{
		"NAMESPACE_OPERATION_UNSPECIFIED": 0,
		"NAMESPACE_OPERATION_CREATE":      1,
		"NAMESPACE_OPERATION_UPDATE":      2,
	}
)

func (x NamespaceOperation) Enum() *NamespaceOperation {
	p := new(NamespaceOperation)
	*p = x
	return p
}

func (x NamespaceOperation) String() string {
	switch x {
	case NAMESPACE_OPERATION_UNSPECIFIED:
		return "Unspecified"
	case NAMESPACE_OPERATION_CREATE:
		return "Create"
	case NAMESPACE_OPERATION_UPDATE:
		return "Update"
	default:
		return strconv.Itoa(int(x))
	}

}

func (NamespaceOperation) Descriptor() protoreflect.EnumDescriptor {
	return file_temporal_server_api_enums_v1_replication_proto_enumTypes[1].Descriptor()
}

func (NamespaceOperation) Type() protoreflect.EnumType {
	return &file_temporal_server_api_enums_v1_replication_proto_enumTypes[1]
}

func (x NamespaceOperation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NamespaceOperation.Descriptor instead.
func (NamespaceOperation) EnumDescriptor() ([]byte, []int) {
	return file_temporal_server_api_enums_v1_replication_proto_rawDescGZIP(), []int{1}
}

type ReplicationFlowControlCommand int32

const (
	REPLICATION_FLOW_CONTROL_COMMAND_UNSPECIFIED ReplicationFlowControlCommand = 0
	REPLICATION_FLOW_CONTROL_COMMAND_RESUME      ReplicationFlowControlCommand = 1
	REPLICATION_FLOW_CONTROL_COMMAND_PAUSE       ReplicationFlowControlCommand = 2
)

// Enum value maps for ReplicationFlowControlCommand.
var (
	ReplicationFlowControlCommand_name = map[int32]string{
		0: "REPLICATION_FLOW_CONTROL_COMMAND_UNSPECIFIED",
		1: "REPLICATION_FLOW_CONTROL_COMMAND_RESUME",
		2: "REPLICATION_FLOW_CONTROL_COMMAND_PAUSE",
	}
	ReplicationFlowControlCommand_value = map[string]int32{
		"REPLICATION_FLOW_CONTROL_COMMAND_UNSPECIFIED": 0,
		"REPLICATION_FLOW_CONTROL_COMMAND_RESUME":      1,
		"REPLICATION_FLOW_CONTROL_COMMAND_PAUSE":       2,
	}
)

func (x ReplicationFlowControlCommand) Enum() *ReplicationFlowControlCommand {
	p := new(ReplicationFlowControlCommand)
	*p = x
	return p
}

func (x ReplicationFlowControlCommand) String() string {
	switch x {
	case REPLICATION_FLOW_CONTROL_COMMAND_UNSPECIFIED:
		return "Unspecified"
	case REPLICATION_FLOW_CONTROL_COMMAND_RESUME:
		return "Resume"
	case REPLICATION_FLOW_CONTROL_COMMAND_PAUSE:
		return "Pause"
	default:
		return strconv.Itoa(int(x))
	}

}

func (ReplicationFlowControlCommand) Descriptor() protoreflect.EnumDescriptor {
	return file_temporal_server_api_enums_v1_replication_proto_enumTypes[2].Descriptor()
}

func (ReplicationFlowControlCommand) Type() protoreflect.EnumType {
	return &file_temporal_server_api_enums_v1_replication_proto_enumTypes[2]
}

func (x ReplicationFlowControlCommand) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReplicationFlowControlCommand.Descriptor instead.
func (ReplicationFlowControlCommand) EnumDescriptor() ([]byte, []int) {
	return file_temporal_server_api_enums_v1_replication_proto_rawDescGZIP(), []int{2}
}

var File_temporal_server_api_enums_v1_replication_proto protoreflect.FileDescriptor

var file_temporal_server_api_enums_v1_replication_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x72,
	0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x1c, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x76, 0x31, 0x2a, 0xc4,
	0x04, 0x0a, 0x13, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x61,
	0x73, 0x6b, 0x54, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x21, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x28, 0x0a,
	0x24, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53,
	0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45,
	0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x01, 0x12, 0x26, 0x0a, 0x22, 0x52, 0x45, 0x50, 0x4c, 0x49,
	0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x48, 0x49, 0x53, 0x54, 0x4f, 0x52, 0x59, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x02, 0x12,
	0x30, 0x0a, 0x2c, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f, 0x53, 0x48,
	0x41, 0x52, 0x44, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10,
	0x03, 0x12, 0x2c, 0x0a, 0x28, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f,
	0x41, 0x43, 0x54, 0x49, 0x56, 0x49, 0x54, 0x59, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x04, 0x12,
	0x2f, 0x0a, 0x2b, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x49, 0x53, 0x54, 0x4f, 0x52, 0x59,
	0x5f, 0x4d, 0x45, 0x54, 0x41, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x05,
	0x12, 0x29, 0x0a, 0x25, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x49, 0x53, 0x54, 0x4f, 0x52,
	0x59, 0x5f, 0x56, 0x32, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x06, 0x12, 0x32, 0x0a, 0x2e, 0x52,
	0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c,
	0x4f, 0x57, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x07, 0x12,
	0x2e, 0x0a, 0x2a, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x51, 0x55,
	0x45, 0x55, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x10, 0x08, 0x12,
	0x27, 0x0a, 0x23, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f, 0x48, 0x53,
	0x4d, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x09, 0x12, 0x2f, 0x0a, 0x2b, 0x52, 0x45, 0x50, 0x4c,
	0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x46, 0x49, 0x4c, 0x4c, 0x5f, 0x48, 0x49, 0x53, 0x54, 0x4f,
	0x52, 0x59, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x0a, 0x12, 0x3a, 0x0a, 0x36, 0x52, 0x45, 0x50,
	0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f,
	0x4e, 0x45, 0x44, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54,
	0x41, 0x53, 0x4b, 0x10, 0x0b, 0x2a, 0x79, 0x0a, 0x12, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x1f, 0x4e,
	0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x1e, 0x0a, 0x1a, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x4f, 0x50,
	0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x01,
	0x12, 0x1e, 0x0a, 0x1a, 0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x4f, 0x50,
	0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x02,
	0x2a, 0xaa, 0x01, 0x0a, 0x1d, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x6c, 0x6f, 0x77, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x12, 0x30, 0x0a, 0x2c, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x4f, 0x4c, 0x5f, 0x43,
	0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x2b, 0x0a, 0x27, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x4f, 0x4c,
	0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4d, 0x45, 0x10,
	0x01, 0x12, 0x2a, 0x0a, 0x26, 0x52, 0x45, 0x50, 0x4c, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x4f, 0x4c, 0x5f, 0x43, 0x4f,
	0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x5f, 0x50, 0x41, 0x55, 0x53, 0x45, 0x10, 0x02, 0x42, 0x2a, 0x5a,
	0x28, 0x67, 0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x69, 0x6f, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x73,
	0x2f, 0x76, 0x31, 0x3b, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_temporal_server_api_enums_v1_replication_proto_rawDescOnce sync.Once
	file_temporal_server_api_enums_v1_replication_proto_rawDescData = file_temporal_server_api_enums_v1_replication_proto_rawDesc
)

func file_temporal_server_api_enums_v1_replication_proto_rawDescGZIP() []byte {
	file_temporal_server_api_enums_v1_replication_proto_rawDescOnce.Do(func() {
		file_temporal_server_api_enums_v1_replication_proto_rawDescData = protoimpl.X.CompressGZIP(file_temporal_server_api_enums_v1_replication_proto_rawDescData)
	})
	return file_temporal_server_api_enums_v1_replication_proto_rawDescData
}

var file_temporal_server_api_enums_v1_replication_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_temporal_server_api_enums_v1_replication_proto_goTypes = []interface{}{
	(ReplicationTaskType)(0),           // 0: temporal.server.api.enums.v1.ReplicationTaskType
	(NamespaceOperation)(0),            // 1: temporal.server.api.enums.v1.NamespaceOperation
	(ReplicationFlowControlCommand)(0), // 2: temporal.server.api.enums.v1.ReplicationFlowControlCommand
}
var file_temporal_server_api_enums_v1_replication_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_server_api_enums_v1_replication_proto_init() }
func file_temporal_server_api_enums_v1_replication_proto_init() {
	if File_temporal_server_api_enums_v1_replication_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_temporal_server_api_enums_v1_replication_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_server_api_enums_v1_replication_proto_goTypes,
		DependencyIndexes: file_temporal_server_api_enums_v1_replication_proto_depIdxs,
		EnumInfos:         file_temporal_server_api_enums_v1_replication_proto_enumTypes,
	}.Build()
	File_temporal_server_api_enums_v1_replication_proto = out.File
	file_temporal_server_api_enums_v1_replication_proto_rawDesc = nil
	file_temporal_server_api_enums_v1_replication_proto_goTypes = nil
	file_temporal_server_api_enums_v1_replication_proto_depIdxs = nil
}
