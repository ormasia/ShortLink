// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: proto/shortlinkpb/shortlink.proto

package shortlinkpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// 请求生成短链接
type ShortenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginalUrl   string                 `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenRequest) Reset() {
	*x = ShortenRequest{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenRequest) ProtoMessage() {}

func (x *ShortenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenRequest.ProtoReflect.Descriptor instead.
func (*ShortenRequest) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{0}
}

func (x *ShortenRequest) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *ShortenRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ShortenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortenResponse) Reset() {
	*x = ShortenResponse{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenResponse) ProtoMessage() {}

func (x *ShortenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenResponse.ProtoReflect.Descriptor instead.
func (*ShortenResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{1}
}

func (x *ShortenResponse) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

// 请求解析短链接
type ResolveRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResolveRequest) Reset() {
	*x = ResolveRequest{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResolveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveRequest) ProtoMessage() {}

func (x *ResolveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveRequest.ProtoReflect.Descriptor instead.
func (*ResolveRequest) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{2}
}

func (x *ResolveRequest) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

type ResolveResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginalUrl   string                 `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResolveResponse) Reset() {
	*x = ResolveResponse{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResolveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveResponse) ProtoMessage() {}

func (x *ResolveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveResponse.ProtoReflect.Descriptor instead.
func (*ResolveResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{3}
}

func (x *ResolveResponse) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type TopRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         int64                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TopRequest) Reset() {
	*x = TopRequest{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopRequest) ProtoMessage() {}

func (x *TopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopRequest.ProtoReflect.Descriptor instead.
func (*TopRequest) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{4}
}

func (x *TopRequest) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ShortLinkItem struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShortUrl      string                 `protobuf:"bytes,1,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	Clicks        float64                `protobuf:"fixed64,2,opt,name=clicks,proto3" json:"clicks,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShortLinkItem) Reset() {
	*x = ShortLinkItem{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShortLinkItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortLinkItem) ProtoMessage() {}

func (x *ShortLinkItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortLinkItem.ProtoReflect.Descriptor instead.
func (*ShortLinkItem) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{5}
}

func (x *ShortLinkItem) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *ShortLinkItem) GetClicks() float64 {
	if x != nil {
		return x.Clicks
	}
	return 0
}

type TopResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Top           []*ShortLinkItem       `protobuf:"bytes,1,rep,name=top,proto3" json:"top,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TopResponse) Reset() {
	*x = TopResponse{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopResponse) ProtoMessage() {}

func (x *TopResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopResponse.ProtoReflect.Descriptor instead.
func (*TopResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{6}
}

func (x *TopResponse) GetTop() []*ShortLinkItem {
	if x != nil {
		return x.Top
	}
	return nil
}

// 批量生成短链接的请求
type BatchShortenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 需要转换的原始长URL列表
	OriginalUrls []string `protobuf:"bytes,1,rep,name=original_urls,json=originalUrls,proto3" json:"original_urls,omitempty"`
	UserId       string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// 并发处理的数量，默认为10
	Concurrency   int32 `protobuf:"varint,3,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchShortenRequest) Reset() {
	*x = BatchShortenRequest{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchShortenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenRequest) ProtoMessage() {}

func (x *BatchShortenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenRequest.ProtoReflect.Descriptor instead.
func (*BatchShortenRequest) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{7}
}

func (x *BatchShortenRequest) GetOriginalUrls() []string {
	if x != nil {
		return x.OriginalUrls
	}
	return nil
}

func (x *BatchShortenRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *BatchShortenRequest) GetConcurrency() int32 {
	if x != nil {
		return x.Concurrency
	}
	return 0
}

// 批量生成短链接的单个结果
type BatchShortenResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 原始URL
	OriginalUrl string `protobuf:"bytes,1,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
	// 生成的短URL
	ShortUrl string `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
	// 错误信息，如果有的话
	Error         string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchShortenResult) Reset() {
	*x = BatchShortenResult{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchShortenResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenResult) ProtoMessage() {}

func (x *BatchShortenResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenResult.ProtoReflect.Descriptor instead.
func (*BatchShortenResult) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{8}
}

func (x *BatchShortenResult) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

func (x *BatchShortenResult) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

func (x *BatchShortenResult) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// 批量生成短链接的响应
type BatchShortenResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// 批量生成结果列表
	Results []*BatchShortenResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	// 总请求数量
	TotalCount int32 `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	// 成功生成的数量
	SuccessCount int32 `protobuf:"varint,3,opt,name=success_count,json=successCount,proto3" json:"success_count,omitempty"`
	// 处理耗时
	ElapsedTime   string `protobuf:"bytes,4,opt,name=elapsed_time,json=elapsedTime,proto3" json:"elapsed_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchShortenResponse) Reset() {
	*x = BatchShortenResponse{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchShortenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchShortenResponse) ProtoMessage() {}

func (x *BatchShortenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchShortenResponse.ProtoReflect.Descriptor instead.
func (*BatchShortenResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{9}
}

func (x *BatchShortenResponse) GetResults() []*BatchShortenResult {
	if x != nil {
		return x.Results
	}
	return nil
}

func (x *BatchShortenResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *BatchShortenResponse) GetSuccessCount() int32 {
	if x != nil {
		return x.SuccessCount
	}
	return 0
}

func (x *BatchShortenResponse) GetElapsedTime() string {
	if x != nil {
		return x.ElapsedTime
	}
	return ""
}

// 删除用户短链接的请求
type DeleteUserURLsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserURLsRequest) Reset() {
	*x = DeleteUserURLsRequest{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserURLsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserURLsRequest) ProtoMessage() {}

func (x *DeleteUserURLsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserURLsRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserURLsRequest) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteUserURLsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// 删除用户短链接的响应
type DeleteUserURLsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DeletedCount  int32                  `protobuf:"varint,1,opt,name=deleted_count,json=deletedCount,proto3" json:"deleted_count,omitempty"` // 删除的短链接数量
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserURLsResponse) Reset() {
	*x = DeleteUserURLsResponse{}
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserURLsResponse) ProtoMessage() {}

func (x *DeleteUserURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortlinkpb_shortlink_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserURLsResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserURLsResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortlinkpb_shortlink_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteUserURLsResponse) GetDeletedCount() int32 {
	if x != nil {
		return x.DeletedCount
	}
	return 0
}

var File_proto_shortlinkpb_shortlink_proto protoreflect.FileDescriptor

const file_proto_shortlinkpb_shortlink_proto_rawDesc = "" +
	"\n" +
	"!proto/shortlinkpb/shortlink.proto\x12\tshortlink\"L\n" +
	"\x0eShortenRequest\x12!\n" +
	"\foriginal_url\x18\x01 \x01(\tR\voriginalUrl\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\".\n" +
	"\x0fShortenResponse\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\"-\n" +
	"\x0eResolveRequest\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\"4\n" +
	"\x0fResolveResponse\x12!\n" +
	"\foriginal_url\x18\x01 \x01(\tR\voriginalUrl\"\"\n" +
	"\n" +
	"TopRequest\x12\x14\n" +
	"\x05count\x18\x01 \x01(\x03R\x05count\"D\n" +
	"\rShortLinkItem\x12\x1b\n" +
	"\tshort_url\x18\x01 \x01(\tR\bshortUrl\x12\x16\n" +
	"\x06clicks\x18\x02 \x01(\x01R\x06clicks\"9\n" +
	"\vTopResponse\x12*\n" +
	"\x03top\x18\x01 \x03(\v2\x18.shortlink.ShortLinkItemR\x03top\"u\n" +
	"\x13BatchShortenRequest\x12#\n" +
	"\roriginal_urls\x18\x01 \x03(\tR\foriginalUrls\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\x12 \n" +
	"\vconcurrency\x18\x03 \x01(\x05R\vconcurrency\"j\n" +
	"\x12BatchShortenResult\x12!\n" +
	"\foriginal_url\x18\x01 \x01(\tR\voriginalUrl\x12\x1b\n" +
	"\tshort_url\x18\x02 \x01(\tR\bshortUrl\x12\x14\n" +
	"\x05error\x18\x03 \x01(\tR\x05error\"\xb8\x01\n" +
	"\x14BatchShortenResponse\x127\n" +
	"\aresults\x18\x01 \x03(\v2\x1d.shortlink.BatchShortenResultR\aresults\x12\x1f\n" +
	"\vtotal_count\x18\x02 \x01(\x05R\n" +
	"totalCount\x12#\n" +
	"\rsuccess_count\x18\x03 \x01(\x05R\fsuccessCount\x12!\n" +
	"\felapsed_time\x18\x04 \x01(\tR\velapsedTime\"0\n" +
	"\x15DeleteUserURLsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"=\n" +
	"\x16DeleteUserURLsResponse\x12#\n" +
	"\rdeleted_count\x18\x01 \x01(\x05R\fdeletedCount2\x85\x03\n" +
	"\x10ShortlinkService\x12C\n" +
	"\n" +
	"ShortenURL\x12\x19.shortlink.ShortenRequest\x1a\x1a.shortlink.ShortenResponse\x12B\n" +
	"\tRedierect\x12\x19.shortlink.ResolveRequest\x1a\x1a.shortlink.ResolveResponse\x12<\n" +
	"\vGetTopLinks\x12\x15.shortlink.TopRequest\x1a\x16.shortlink.TopResponse\x12S\n" +
	"\x10BatchShortenURLs\x12\x1e.shortlink.BatchShortenRequest\x1a\x1f.shortlink.BatchShortenResponse\x12U\n" +
	"\x0eDeleteUserURLs\x12 .shortlink.DeleteUserURLsRequest\x1a!.shortlink.DeleteUserURLsResponseB\x15Z\x13./proto/shortlinkpbb\x06proto3"

var (
	file_proto_shortlinkpb_shortlink_proto_rawDescOnce sync.Once
	file_proto_shortlinkpb_shortlink_proto_rawDescData []byte
)

func file_proto_shortlinkpb_shortlink_proto_rawDescGZIP() []byte {
	file_proto_shortlinkpb_shortlink_proto_rawDescOnce.Do(func() {
		file_proto_shortlinkpb_shortlink_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_shortlinkpb_shortlink_proto_rawDesc), len(file_proto_shortlinkpb_shortlink_proto_rawDesc)))
	})
	return file_proto_shortlinkpb_shortlink_proto_rawDescData
}

var file_proto_shortlinkpb_shortlink_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_shortlinkpb_shortlink_proto_goTypes = []any{
	(*ShortenRequest)(nil),         // 0: shortlink.ShortenRequest
	(*ShortenResponse)(nil),        // 1: shortlink.ShortenResponse
	(*ResolveRequest)(nil),         // 2: shortlink.ResolveRequest
	(*ResolveResponse)(nil),        // 3: shortlink.ResolveResponse
	(*TopRequest)(nil),             // 4: shortlink.TopRequest
	(*ShortLinkItem)(nil),          // 5: shortlink.ShortLinkItem
	(*TopResponse)(nil),            // 6: shortlink.TopResponse
	(*BatchShortenRequest)(nil),    // 7: shortlink.BatchShortenRequest
	(*BatchShortenResult)(nil),     // 8: shortlink.BatchShortenResult
	(*BatchShortenResponse)(nil),   // 9: shortlink.BatchShortenResponse
	(*DeleteUserURLsRequest)(nil),  // 10: shortlink.DeleteUserURLsRequest
	(*DeleteUserURLsResponse)(nil), // 11: shortlink.DeleteUserURLsResponse
}
var file_proto_shortlinkpb_shortlink_proto_depIdxs = []int32{
	5,  // 0: shortlink.TopResponse.top:type_name -> shortlink.ShortLinkItem
	8,  // 1: shortlink.BatchShortenResponse.results:type_name -> shortlink.BatchShortenResult
	0,  // 2: shortlink.ShortlinkService.ShortenURL:input_type -> shortlink.ShortenRequest
	2,  // 3: shortlink.ShortlinkService.Redierect:input_type -> shortlink.ResolveRequest
	4,  // 4: shortlink.ShortlinkService.GetTopLinks:input_type -> shortlink.TopRequest
	7,  // 5: shortlink.ShortlinkService.BatchShortenURLs:input_type -> shortlink.BatchShortenRequest
	10, // 6: shortlink.ShortlinkService.DeleteUserURLs:input_type -> shortlink.DeleteUserURLsRequest
	1,  // 7: shortlink.ShortlinkService.ShortenURL:output_type -> shortlink.ShortenResponse
	3,  // 8: shortlink.ShortlinkService.Redierect:output_type -> shortlink.ResolveResponse
	6,  // 9: shortlink.ShortlinkService.GetTopLinks:output_type -> shortlink.TopResponse
	9,  // 10: shortlink.ShortlinkService.BatchShortenURLs:output_type -> shortlink.BatchShortenResponse
	11, // 11: shortlink.ShortlinkService.DeleteUserURLs:output_type -> shortlink.DeleteUserURLsResponse
	7,  // [7:12] is the sub-list for method output_type
	2,  // [2:7] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_proto_shortlinkpb_shortlink_proto_init() }
func file_proto_shortlinkpb_shortlink_proto_init() {
	if File_proto_shortlinkpb_shortlink_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_shortlinkpb_shortlink_proto_rawDesc), len(file_proto_shortlinkpb_shortlink_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_shortlinkpb_shortlink_proto_goTypes,
		DependencyIndexes: file_proto_shortlinkpb_shortlink_proto_depIdxs,
		MessageInfos:      file_proto_shortlinkpb_shortlink_proto_msgTypes,
	}.Build()
	File_proto_shortlinkpb_shortlink_proto = out.File
	file_proto_shortlinkpb_shortlink_proto_goTypes = nil
	file_proto_shortlinkpb_shortlink_proto_depIdxs = nil
}
