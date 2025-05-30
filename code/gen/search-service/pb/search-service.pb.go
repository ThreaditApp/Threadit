// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.4
// source: search-service.proto

package pb

import (
	pb "gen/models/pb"
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

type SearchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Query         string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Offset        *int32                 `protobuf:"varint,2,opt,name=offset,proto3,oneof" json:"offset,omitempty"`
	Limit         *int32                 `protobuf:"varint,3,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	mi := &file_search_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_search_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_search_service_proto_rawDescGZIP(), []int{0}
}

func (x *SearchRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchRequest) GetOffset() int32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

func (x *SearchRequest) GetLimit() int32 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

type GlobalSearchResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	ThreadResults    []*pb.Thread           `protobuf:"bytes,1,rep,name=thread_results,json=threadResults,proto3" json:"thread_results,omitempty"`
	CommunityResults []*pb.Community        `protobuf:"bytes,2,rep,name=community_results,json=communityResults,proto3" json:"community_results,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *GlobalSearchResponse) Reset() {
	*x = GlobalSearchResponse{}
	mi := &file_search_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GlobalSearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalSearchResponse) ProtoMessage() {}

func (x *GlobalSearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalSearchResponse.ProtoReflect.Descriptor instead.
func (*GlobalSearchResponse) Descriptor() ([]byte, []int) {
	return file_search_service_proto_rawDescGZIP(), []int{1}
}

func (x *GlobalSearchResponse) GetThreadResults() []*pb.Thread {
	if x != nil {
		return x.ThreadResults
	}
	return nil
}

func (x *GlobalSearchResponse) GetCommunityResults() []*pb.Community {
	if x != nil {
		return x.CommunityResults
	}
	return nil
}

type CommunitySearchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Results       []*pb.Community        `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommunitySearchResponse) Reset() {
	*x = CommunitySearchResponse{}
	mi := &file_search_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommunitySearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommunitySearchResponse) ProtoMessage() {}

func (x *CommunitySearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommunitySearchResponse.ProtoReflect.Descriptor instead.
func (*CommunitySearchResponse) Descriptor() ([]byte, []int) {
	return file_search_service_proto_rawDescGZIP(), []int{2}
}

func (x *CommunitySearchResponse) GetResults() []*pb.Community {
	if x != nil {
		return x.Results
	}
	return nil
}

type ThreadSearchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Results       []*pb.Thread           `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ThreadSearchResponse) Reset() {
	*x = ThreadSearchResponse{}
	mi := &file_search_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ThreadSearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ThreadSearchResponse) ProtoMessage() {}

func (x *ThreadSearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ThreadSearchResponse.ProtoReflect.Descriptor instead.
func (*ThreadSearchResponse) Descriptor() ([]byte, []int) {
	return file_search_service_proto_rawDescGZIP(), []int{3}
}

func (x *ThreadSearchResponse) GetResults() []*pb.Thread {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_search_service_proto protoreflect.FileDescriptor

const file_search_service_proto_rawDesc = "" +
	"\n" +
	"\x14search-service.proto\x12\x06search\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/api/annotations.proto\x1a\fmodels.proto\"r\n" +
	"\rSearchRequest\x12\x14\n" +
	"\x05query\x18\x01 \x01(\tR\x05query\x12\x1b\n" +
	"\x06offset\x18\x02 \x01(\x05H\x00R\x06offset\x88\x01\x01\x12\x19\n" +
	"\x05limit\x18\x03 \x01(\x05H\x01R\x05limit\x88\x01\x01B\t\n" +
	"\a_offsetB\b\n" +
	"\x06_limit\"\x8d\x01\n" +
	"\x14GlobalSearchResponse\x125\n" +
	"\x0ethread_results\x18\x01 \x03(\v2\x0e.models.ThreadR\rthreadResults\x12>\n" +
	"\x11community_results\x18\x02 \x03(\v2\x11.models.CommunityR\x10communityResults\"F\n" +
	"\x17CommunitySearchResponse\x12+\n" +
	"\aresults\x18\x01 \x03(\v2\x11.models.CommunityR\aresults\"@\n" +
	"\x14ThreadSearchResponse\x12(\n" +
	"\aresults\x18\x01 \x03(\v2\x0e.models.ThreadR\aresults2\xe7\x02\n" +
	"\rSearchService\x12=\n" +
	"\vCheckHealth\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\x12T\n" +
	"\fGlobalSearch\x12\x15.search.SearchRequest\x1a\x1c.search.GlobalSearchResponse\"\x0f\x82\xd3\xe4\x93\x02\t\x12\a/search\x12d\n" +
	"\x0fCommunitySearch\x12\x15.search.SearchRequest\x1a\x1f.search.CommunitySearchResponse\"\x19\x82\xd3\xe4\x93\x02\x13\x12\x11/search/community\x12[\n" +
	"\fThreadSearch\x12\x15.search.SearchRequest\x1a\x1c.search.ThreadSearchResponse\"\x16\x82\xd3\xe4\x93\x02\x10\x12\x0e/search/threadB\x1aZ\x18gen/search-service/pb;pbb\x06proto3"

var (
	file_search_service_proto_rawDescOnce sync.Once
	file_search_service_proto_rawDescData []byte
)

func file_search_service_proto_rawDescGZIP() []byte {
	file_search_service_proto_rawDescOnce.Do(func() {
		file_search_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_search_service_proto_rawDesc), len(file_search_service_proto_rawDesc)))
	})
	return file_search_service_proto_rawDescData
}

var file_search_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_search_service_proto_goTypes = []any{
	(*SearchRequest)(nil),           // 0: search.SearchRequest
	(*GlobalSearchResponse)(nil),    // 1: search.GlobalSearchResponse
	(*CommunitySearchResponse)(nil), // 2: search.CommunitySearchResponse
	(*ThreadSearchResponse)(nil),    // 3: search.ThreadSearchResponse
	(*pb.Thread)(nil),               // 4: models.Thread
	(*pb.Community)(nil),            // 5: models.Community
	(*emptypb.Empty)(nil),           // 6: google.protobuf.Empty
}
var file_search_service_proto_depIdxs = []int32{
	4, // 0: search.GlobalSearchResponse.thread_results:type_name -> models.Thread
	5, // 1: search.GlobalSearchResponse.community_results:type_name -> models.Community
	5, // 2: search.CommunitySearchResponse.results:type_name -> models.Community
	4, // 3: search.ThreadSearchResponse.results:type_name -> models.Thread
	6, // 4: search.SearchService.CheckHealth:input_type -> google.protobuf.Empty
	0, // 5: search.SearchService.GlobalSearch:input_type -> search.SearchRequest
	0, // 6: search.SearchService.CommunitySearch:input_type -> search.SearchRequest
	0, // 7: search.SearchService.ThreadSearch:input_type -> search.SearchRequest
	6, // 8: search.SearchService.CheckHealth:output_type -> google.protobuf.Empty
	1, // 9: search.SearchService.GlobalSearch:output_type -> search.GlobalSearchResponse
	2, // 10: search.SearchService.CommunitySearch:output_type -> search.CommunitySearchResponse
	3, // 11: search.SearchService.ThreadSearch:output_type -> search.ThreadSearchResponse
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_search_service_proto_init() }
func file_search_service_proto_init() {
	if File_search_service_proto != nil {
		return
	}
	file_search_service_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_search_service_proto_rawDesc), len(file_search_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_search_service_proto_goTypes,
		DependencyIndexes: file_search_service_proto_depIdxs,
		MessageInfos:      file_search_service_proto_msgTypes,
	}.Build()
	File_search_service_proto = out.File
	file_search_service_proto_goTypes = nil
	file_search_service_proto_depIdxs = nil
}
