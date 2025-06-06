// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: vote-service.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VoteThreadRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ThreadId      string                 `protobuf:"bytes,1,opt,name=thread_id,json=threadId,proto3" json:"thread_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteThreadRequest) Reset() {
	*x = VoteThreadRequest{}
	mi := &file_vote_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteThreadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteThreadRequest) ProtoMessage() {}

func (x *VoteThreadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vote_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteThreadRequest.ProtoReflect.Descriptor instead.
func (*VoteThreadRequest) Descriptor() ([]byte, []int) {
	return file_vote_service_proto_rawDescGZIP(), []int{0}
}

func (x *VoteThreadRequest) GetThreadId() string {
	if x != nil {
		return x.ThreadId
	}
	return ""
}

type VoteCommentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CommentId     string                 `protobuf:"bytes,1,opt,name=comment_id,json=commentId,proto3" json:"comment_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteCommentRequest) Reset() {
	*x = VoteCommentRequest{}
	mi := &file_vote_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteCommentRequest) ProtoMessage() {}

func (x *VoteCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_vote_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteCommentRequest.ProtoReflect.Descriptor instead.
func (*VoteCommentRequest) Descriptor() ([]byte, []int) {
	return file_vote_service_proto_rawDescGZIP(), []int{1}
}

func (x *VoteCommentRequest) GetCommentId() string {
	if x != nil {
		return x.CommentId
	}
	return ""
}

var File_vote_service_proto protoreflect.FileDescriptor

const file_vote_service_proto_rawDesc = "" +
	"\n" +
	"\x12vote-service.proto\x12\x04vote\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/api/annotations.proto\"0\n" +
	"\x11VoteThreadRequest\x12\x1b\n" +
	"\tthread_id\x18\x01 \x01(\tR\bthreadId\"3\n" +
	"\x12VoteCommentRequest\x12\x1d\n" +
	"\n" +
	"comment_id\x18\x01 \x01(\tR\tcommentId2\x84\x04\n" +
	"\vVoteService\x12=\n" +
	"\vCheckHealth\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12h\n" +
	"\fUpvoteThread\x12\x17.vote.VoteThreadRequest\x1a\x16.google.protobuf.Empty\"'\x82\xd3\xe4\x93\x02!:\x01*\"\x1c/votes/thread/{thread_id}/up\x12l\n" +
	"\x0eDownvoteThread\x12\x17.vote.VoteThreadRequest\x1a\x16.google.protobuf.Empty\")\x82\xd3\xe4\x93\x02#:\x01*\"\x1e/votes/thread/{thread_id}/down\x12l\n" +
	"\rUpvoteComment\x12\x18.vote.VoteCommentRequest\x1a\x16.google.protobuf.Empty\")\x82\xd3\xe4\x93\x02#:\x01*\"\x1e/votes/comment/{comment_id}/up\x12p\n" +
	"\x0fDownvoteComment\x12\x18.vote.VoteCommentRequest\x1a\x16.google.protobuf.Empty\"+\x82\xd3\xe4\x93\x02%:\x01*\" /votes/comment/{comment_id}/downB\x18Z\x16gen/vote-service/pb;pbb\x06proto3"

var (
	file_vote_service_proto_rawDescOnce sync.Once
	file_vote_service_proto_rawDescData []byte
)

func file_vote_service_proto_rawDescGZIP() []byte {
	file_vote_service_proto_rawDescOnce.Do(func() {
		file_vote_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_vote_service_proto_rawDesc), len(file_vote_service_proto_rawDesc)))
	})
	return file_vote_service_proto_rawDescData
}

var file_vote_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_vote_service_proto_goTypes = []any{
	(*VoteThreadRequest)(nil),  // 0: vote.VoteThreadRequest
	(*VoteCommentRequest)(nil), // 1: vote.VoteCommentRequest
	(*emptypb.Empty)(nil),      // 2: google.protobuf.Empty
}
var file_vote_service_proto_depIdxs = []int32{
	2, // 0: vote.VoteService.CheckHealth:input_type -> google.protobuf.Empty
	0, // 1: vote.VoteService.UpvoteThread:input_type -> vote.VoteThreadRequest
	0, // 2: vote.VoteService.DownvoteThread:input_type -> vote.VoteThreadRequest
	1, // 3: vote.VoteService.UpvoteComment:input_type -> vote.VoteCommentRequest
	1, // 4: vote.VoteService.DownvoteComment:input_type -> vote.VoteCommentRequest
	2, // 5: vote.VoteService.CheckHealth:output_type -> google.protobuf.Empty
	2, // 6: vote.VoteService.UpvoteThread:output_type -> google.protobuf.Empty
	2, // 7: vote.VoteService.DownvoteThread:output_type -> google.protobuf.Empty
	2, // 8: vote.VoteService.UpvoteComment:output_type -> google.protobuf.Empty
	2, // 9: vote.VoteService.DownvoteComment:output_type -> google.protobuf.Empty
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_vote_service_proto_init() }
func file_vote_service_proto_init() {
	if File_vote_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_vote_service_proto_rawDesc), len(file_vote_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_vote_service_proto_goTypes,
		DependencyIndexes: file_vote_service_proto_depIdxs,
		MessageInfos:      file_vote_service_proto_msgTypes,
	}.Build()
	File_vote_service_proto = out.File
	file_vote_service_proto_goTypes = nil
	file_vote_service_proto_depIdxs = nil
}
