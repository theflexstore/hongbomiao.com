// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: api/graphql_server/proto/torchserve/v1/inference.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PredictionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of model.
	ModelName string `protobuf:"bytes,1,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"` //required
	// Version of model to run prediction on.
	ModelVersion string `protobuf:"bytes,2,opt,name=model_version,json=modelVersion,proto3" json:"model_version,omitempty"` //optional
	// input data for model prediction
	Input map[string][]byte `protobuf:"bytes,3,rep,name=input,proto3" json:"input,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` //required
}

func (x *PredictionsRequest) Reset() {
	*x = PredictionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictionsRequest) ProtoMessage() {}

func (x *PredictionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictionsRequest.ProtoReflect.Descriptor instead.
func (*PredictionsRequest) Descriptor() ([]byte, []int) {
	return file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescGZIP(), []int{0}
}

func (x *PredictionsRequest) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *PredictionsRequest) GetModelVersion() string {
	if x != nil {
		return x.ModelVersion
	}
	return ""
}

func (x *PredictionsRequest) GetInput() map[string][]byte {
	if x != nil {
		return x.Input
	}
	return nil
}

type PredictionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TorchServe health
	Prediction []byte `protobuf:"bytes,1,opt,name=prediction,proto3" json:"prediction,omitempty"`
}

func (x *PredictionResponse) Reset() {
	*x = PredictionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PredictionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictionResponse) ProtoMessage() {}

func (x *PredictionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictionResponse.ProtoReflect.Descriptor instead.
func (*PredictionResponse) Descriptor() ([]byte, []int) {
	return file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescGZIP(), []int{1}
}

func (x *PredictionResponse) GetPrediction() []byte {
	if x != nil {
		return x.Prediction
	}
	return nil
}

type TorchServeHealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TorchServe health
	Health string `protobuf:"bytes,1,opt,name=health,proto3" json:"health,omitempty"`
}

func (x *TorchServeHealthResponse) Reset() {
	*x = TorchServeHealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorchServeHealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorchServeHealthResponse) ProtoMessage() {}

func (x *TorchServeHealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorchServeHealthResponse.ProtoReflect.Descriptor instead.
func (*TorchServeHealthResponse) Descriptor() ([]byte, []int) {
	return file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescGZIP(), []int{2}
}

func (x *TorchServeHealthResponse) GetHealth() string {
	if x != nil {
		return x.Health
	}
	return ""
}

var File_api_graphql_server_proto_torchserve_v1_inference_proto protoreflect.FileDescriptor

var file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDesc = []byte{
	0x0a, 0x36, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x5f, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x6f, 0x72, 0x63, 0x68,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x20, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x79,
	0x74, 0x6f, 0x72, 0x63, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe9, 0x01, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x55, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x3f, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x79, 0x74, 0x6f, 0x72, 0x63, 0x68, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x69, 0x6e, 0x66, 0x65, 0x72,
	0x65, 0x6e, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x34, 0x0a, 0x12, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70,
	0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x32, 0x0a, 0x18, 0x54, 0x6f, 0x72,
	0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x65, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x32, 0xf1, 0x01,
	0x0a, 0x14, 0x49, 0x6e, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x41, 0x50, 0x49, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x3a, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x79, 0x74,
	0x6f, 0x72, 0x63, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x69, 0x6e, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x54, 0x6f, 0x72, 0x63, 0x68, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x7b, 0x0a, 0x0b, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x34, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x79, 0x74, 0x6f, 0x72, 0x63,
	0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x69, 0x6e, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x6f, 0x72, 0x67, 0x2e,
	0x70, 0x79, 0x74, 0x6f, 0x72, 0x63, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x69, 0x6e, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x2a, 0x5a, 0x28, 0x2e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x71, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x74, 0x6f, 0x72, 0x63, 0x68, 0x73, 0x65, 0x72, 0x76, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescOnce sync.Once
	file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescData = file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDesc
)

func file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescGZIP() []byte {
	file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescOnce.Do(func() {
		file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescData)
	})
	return file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDescData
}

var file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_graphql_server_proto_torchserve_v1_inference_proto_goTypes = []interface{}{
	(*PredictionsRequest)(nil),       // 0: org.pytorch.serve.grpc.inference.PredictionsRequest
	(*PredictionResponse)(nil),       // 1: org.pytorch.serve.grpc.inference.PredictionResponse
	(*TorchServeHealthResponse)(nil), // 2: org.pytorch.serve.grpc.inference.TorchServeHealthResponse
	nil,                              // 3: org.pytorch.serve.grpc.inference.PredictionsRequest.InputEntry
	(*emptypb.Empty)(nil),            // 4: google.protobuf.Empty
}
var file_api_graphql_server_proto_torchserve_v1_inference_proto_depIdxs = []int32{
	3, // 0: org.pytorch.serve.grpc.inference.PredictionsRequest.input:type_name -> org.pytorch.serve.grpc.inference.PredictionsRequest.InputEntry
	4, // 1: org.pytorch.serve.grpc.inference.InferenceAPIsService.Ping:input_type -> google.protobuf.Empty
	0, // 2: org.pytorch.serve.grpc.inference.InferenceAPIsService.Predictions:input_type -> org.pytorch.serve.grpc.inference.PredictionsRequest
	2, // 3: org.pytorch.serve.grpc.inference.InferenceAPIsService.Ping:output_type -> org.pytorch.serve.grpc.inference.TorchServeHealthResponse
	1, // 4: org.pytorch.serve.grpc.inference.InferenceAPIsService.Predictions:output_type -> org.pytorch.serve.grpc.inference.PredictionResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_graphql_server_proto_torchserve_v1_inference_proto_init() }
func file_api_graphql_server_proto_torchserve_v1_inference_proto_init() {
	if File_api_graphql_server_proto_torchserve_v1_inference_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictionsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PredictionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorchServeHealthResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_graphql_server_proto_torchserve_v1_inference_proto_goTypes,
		DependencyIndexes: file_api_graphql_server_proto_torchserve_v1_inference_proto_depIdxs,
		MessageInfos:      file_api_graphql_server_proto_torchserve_v1_inference_proto_msgTypes,
	}.Build()
	File_api_graphql_server_proto_torchserve_v1_inference_proto = out.File
	file_api_graphql_server_proto_torchserve_v1_inference_proto_rawDesc = nil
	file_api_graphql_server_proto_torchserve_v1_inference_proto_goTypes = nil
	file_api_graphql_server_proto_torchserve_v1_inference_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// InferenceAPIsServiceClient is the client API for InferenceAPIsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InferenceAPIsServiceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TorchServeHealthResponse, error)
	// Predictions entry point to get inference using default model version.
	Predictions(ctx context.Context, in *PredictionsRequest, opts ...grpc.CallOption) (*PredictionResponse, error)
}

type inferenceAPIsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInferenceAPIsServiceClient(cc grpc.ClientConnInterface) InferenceAPIsServiceClient {
	return &inferenceAPIsServiceClient{cc}
}

func (c *inferenceAPIsServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TorchServeHealthResponse, error) {
	out := new(TorchServeHealthResponse)
	err := c.cc.Invoke(ctx, "/org.pytorch.serve.grpc.inference.InferenceAPIsService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inferenceAPIsServiceClient) Predictions(ctx context.Context, in *PredictionsRequest, opts ...grpc.CallOption) (*PredictionResponse, error) {
	out := new(PredictionResponse)
	err := c.cc.Invoke(ctx, "/org.pytorch.serve.grpc.inference.InferenceAPIsService/Predictions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InferenceAPIsServiceServer is the server API for InferenceAPIsService service.
type InferenceAPIsServiceServer interface {
	Ping(context.Context, *emptypb.Empty) (*TorchServeHealthResponse, error)
	// Predictions entry point to get inference using default model version.
	Predictions(context.Context, *PredictionsRequest) (*PredictionResponse, error)
}

// UnimplementedInferenceAPIsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedInferenceAPIsServiceServer struct {
}

func (*UnimplementedInferenceAPIsServiceServer) Ping(context.Context, *emptypb.Empty) (*TorchServeHealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedInferenceAPIsServiceServer) Predictions(context.Context, *PredictionsRequest) (*PredictionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Predictions not implemented")
}

func RegisterInferenceAPIsServiceServer(s *grpc.Server, srv InferenceAPIsServiceServer) {
	s.RegisterService(&_InferenceAPIsService_serviceDesc, srv)
}

func _InferenceAPIsService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InferenceAPIsServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/org.pytorch.serve.grpc.inference.InferenceAPIsService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InferenceAPIsServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _InferenceAPIsService_Predictions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InferenceAPIsServiceServer).Predictions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/org.pytorch.serve.grpc.inference.InferenceAPIsService/Predictions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InferenceAPIsServiceServer).Predictions(ctx, req.(*PredictionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _InferenceAPIsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "org.pytorch.serve.grpc.inference.InferenceAPIsService",
	HandlerType: (*InferenceAPIsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _InferenceAPIsService_Ping_Handler,
		},
		{
			MethodName: "Predictions",
			Handler:    _InferenceAPIsService_Predictions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/graphql_server/proto/torchserve/v1/inference.proto",
}
