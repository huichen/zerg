// Code generated by protoc-gen-go.
// source: crawl.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	crawl.proto

It has these top-level messages:
	CrawlRequest
	KV
	CrawlResponse
	Metadata
*/
package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Method int32

const (
	Method_GET  Method = 0
	Method_HEAD Method = 1
	Method_POST Method = 2
)

var Method_name = map[int32]string{
	0: "GET",
	1: "HEAD",
	2: "POST",
}
var Method_value = map[string]int32{
	"GET":  0,
	"HEAD": 1,
	"POST": 2,
}

func (x Method) String() string {
	return proto.EnumName(Method_name, int32(x))
}
func (Method) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CrawlRequest struct {
	// 以 http:// 或者 https:// 开头的网址
	Url string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	// 抓取超时限制，单位毫秒，设为 0 时无超时
	Timeout int64 `protobuf:"varint,2,opt,name=timeout" json:"timeout,omitempty"`
	// 是否仅返回 metadata 而忽略 content
	OnlyReturnMetadata bool `protobuf:"varint,5,opt,name=only_return_metadata,json=onlyReturnMetadata" json:"only_return_metadata,omitempty"`
	// 请求的自定义 header
	Header []*KV `protobuf:"bytes,6,rep,name=header" json:"header,omitempty"`
	// 请求方法
	Method Method `protobuf:"varint,7,opt,name=method,enum=protos.Method" json:"method,omitempty"`
	// POST body，仅当请求类型为 POST 时有效
	PostBody string `protobuf:"bytes,8,opt,name=post_body,json=postBody" json:"post_body,omitempty"`
	BodyType string `protobuf:"bytes,9,opt,name=body_type,json=bodyType" json:"body_type,omitempty"`
}

func (m *CrawlRequest) Reset()                    { *m = CrawlRequest{} }
func (m *CrawlRequest) String() string            { return proto.CompactTextString(m) }
func (*CrawlRequest) ProtoMessage()               {}
func (*CrawlRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CrawlRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *CrawlRequest) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *CrawlRequest) GetOnlyReturnMetadata() bool {
	if m != nil {
		return m.OnlyReturnMetadata
	}
	return false
}

func (m *CrawlRequest) GetHeader() []*KV {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *CrawlRequest) GetMethod() Method {
	if m != nil {
		return m.Method
	}
	return Method_GET
}

func (m *CrawlRequest) GetPostBody() string {
	if m != nil {
		return m.PostBody
	}
	return ""
}

func (m *CrawlRequest) GetBodyType() string {
	if m != nil {
		return m.BodyType
	}
	return ""
}

type KV struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *KV) Reset()                    { *m = KV{} }
func (m *KV) String() string            { return proto.CompactTextString(m) }
func (*KV) ProtoMessage()               {}
func (*KV) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *KV) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KV) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type CrawlResponse struct {
	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	Content  string    `protobuf:"bytes,2,opt,name=content" json:"content,omitempty"`
}

func (m *CrawlResponse) Reset()                    { *m = CrawlResponse{} }
func (m *CrawlResponse) String() string            { return proto.CompactTextString(m) }
func (*CrawlResponse) ProtoMessage()               {}
func (*CrawlResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CrawlResponse) GetMetadata() *Metadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *CrawlResponse) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type Metadata struct {
	Length     uint32 `protobuf:"varint,1,opt,name=length" json:"length,omitempty"`
	Header     []*KV  `protobuf:"bytes,3,rep,name=header" json:"header,omitempty"`
	Status     string `protobuf:"bytes,4,opt,name=status" json:"status,omitempty"`
	StatusCode int32  `protobuf:"varint,5,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
}

func (m *Metadata) Reset()                    { *m = Metadata{} }
func (m *Metadata) String() string            { return proto.CompactTextString(m) }
func (*Metadata) ProtoMessage()               {}
func (*Metadata) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Metadata) GetLength() uint32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *Metadata) GetHeader() []*KV {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Metadata) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Metadata) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func init() {
	proto.RegisterType((*CrawlRequest)(nil), "protos.CrawlRequest")
	proto.RegisterType((*KV)(nil), "protos.KV")
	proto.RegisterType((*CrawlResponse)(nil), "protos.CrawlResponse")
	proto.RegisterType((*Metadata)(nil), "protos.Metadata")
	proto.RegisterEnum("protos.Method", Method_name, Method_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Crawl service

type CrawlClient interface {
	Crawl(ctx context.Context, in *CrawlRequest, opts ...grpc.CallOption) (*CrawlResponse, error)
}

type crawlClient struct {
	cc *grpc.ClientConn
}

func NewCrawlClient(cc *grpc.ClientConn) CrawlClient {
	return &crawlClient{cc}
}

func (c *crawlClient) Crawl(ctx context.Context, in *CrawlRequest, opts ...grpc.CallOption) (*CrawlResponse, error) {
	out := new(CrawlResponse)
	err := grpc.Invoke(ctx, "/protos.Crawl/Crawl", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Crawl service

type CrawlServer interface {
	Crawl(context.Context, *CrawlRequest) (*CrawlResponse, error)
}

func RegisterCrawlServer(s *grpc.Server, srv CrawlServer) {
	s.RegisterService(&_Crawl_serviceDesc, srv)
}

func _Crawl_Crawl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CrawlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrawlServer).Crawl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Crawl/Crawl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrawlServer).Crawl(ctx, req.(*CrawlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Crawl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Crawl",
	HandlerType: (*CrawlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Crawl",
			Handler:    _Crawl_Crawl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "crawl.proto",
}

func init() { proto.RegisterFile("crawl.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 409 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x92, 0xd1, 0x6a, 0xd4, 0x40,
	0x14, 0x86, 0x9b, 0xa4, 0x99, 0x4d, 0x4e, 0x6c, 0x09, 0xc3, 0x5a, 0x06, 0xbd, 0x30, 0x04, 0x94,
	0x20, 0xa5, 0x48, 0x04, 0x6f, 0xc5, 0xd6, 0x62, 0xa1, 0x16, 0x65, 0x5c, 0x76, 0x2f, 0x43, 0x76,
	0x73, 0x70, 0x65, 0xb3, 0x99, 0x98, 0x99, 0x28, 0xb9, 0xf2, 0x5d, 0x7c, 0x44, 0x9f, 0x40, 0x66,
	0x26, 0x59, 0x16, 0xc1, 0xab, 0x9c, 0xff, 0xff, 0x27, 0x93, 0x73, 0xbe, 0x1c, 0x88, 0x36, 0x5d,
	0xf9, 0xb3, 0xbe, 0x6a, 0x3b, 0xa1, 0x04, 0x25, 0xe6, 0x21, 0xd3, 0x3f, 0x0e, 0x3c, 0xba, 0xd1,
	0x3e, 0xc7, 0xef, 0x3d, 0x4a, 0x45, 0x63, 0xf0, 0xfa, 0xae, 0x66, 0x4e, 0xe2, 0x64, 0x21, 0xd7,
	0x25, 0x65, 0x30, 0x53, 0xdf, 0xf6, 0x28, 0x7a, 0xc5, 0xdc, 0xc4, 0xc9, 0x3c, 0x3e, 0x49, 0xfa,
	0x0a, 0xe6, 0xa2, 0xa9, 0x87, 0xa2, 0x43, 0xd5, 0x77, 0x4d, 0xb1, 0x47, 0x55, 0x56, 0xa5, 0x2a,
	0x99, 0x9f, 0x38, 0x59, 0xc0, 0xa9, 0xce, 0xb8, 0x89, 0x1e, 0xc6, 0x84, 0xa6, 0x40, 0xb6, 0x58,
	0x56, 0xd8, 0x31, 0x92, 0x78, 0x59, 0x94, 0x83, 0x6d, 0x47, 0x5e, 0xdd, 0x2f, 0xf9, 0x98, 0xd0,
	0x17, 0x40, 0xf6, 0xa8, 0xb6, 0xa2, 0x62, 0xb3, 0xc4, 0xc9, 0xce, 0xf3, 0xf3, 0xe9, 0xcc, 0x83,
	0x71, 0xf9, 0x98, 0xd2, 0xa7, 0x10, 0xb6, 0x42, 0xaa, 0x62, 0x2d, 0xaa, 0x81, 0x05, 0xa6, 0xdf,
	0x40, 0x1b, 0xd7, 0xa2, 0x1a, 0x74, 0xa8, 0xfd, 0x42, 0x0d, 0x2d, 0xb2, 0xd0, 0x86, 0xda, 0x58,
	0x0c, 0x2d, 0xa6, 0x97, 0xe0, 0xde, 0x2f, 0xf5, 0xa4, 0x3b, 0x1c, 0xa6, 0x49, 0x77, 0x38, 0xd0,
	0x39, 0xf8, 0x3f, 0xca, 0xba, 0x47, 0x33, 0x67, 0xc8, 0xad, 0x48, 0x57, 0x70, 0x36, 0x12, 0x92,
	0xad, 0x68, 0x24, 0xd2, 0x4b, 0x08, 0x0e, 0xa3, 0xea, 0xb7, 0xa3, 0x3c, 0x3e, 0x6a, 0xd1, 0xf8,
	0xfc, 0x70, 0x42, 0xe3, 0xdb, 0x88, 0x46, 0x61, 0xa3, 0xc6, 0x6b, 0x27, 0x99, 0xfe, 0x82, 0xe0,
	0x00, 0xe6, 0x02, 0x48, 0x8d, 0xcd, 0x57, 0xb5, 0x35, 0x37, 0x9e, 0xf1, 0x51, 0x1d, 0x01, 0xf3,
	0xfe, 0x0b, 0xec, 0x02, 0x88, 0x54, 0xa5, 0xea, 0x25, 0x3b, 0x35, 0x1f, 0x18, 0x15, 0x7d, 0x06,
	0x91, 0xad, 0x8a, 0x8d, 0xa8, 0xd0, 0xfc, 0x15, 0x9f, 0x83, 0xb5, 0x6e, 0x44, 0x85, 0x2f, 0x9f,
	0x03, 0xb1, 0x4c, 0xe9, 0x0c, 0xbc, 0x0f, 0xb7, 0x8b, 0xf8, 0x84, 0x06, 0x70, 0x7a, 0x77, 0xfb,
	0xee, 0x7d, 0xec, 0xe8, 0xea, 0xf3, 0xa7, 0x2f, 0x8b, 0xd8, 0xcd, 0xdf, 0x82, 0x6f, 0x00, 0xd0,
	0x37, 0x53, 0x31, 0x9f, 0xba, 0x38, 0x5e, 0x9d, 0x27, 0x8f, 0xff, 0x71, 0x2d, 0xae, 0xf4, 0xe4,
	0x9a, 0xfc, 0x76, 0xbd, 0xbb, 0x8f, 0xab, 0xb5, 0x5d, 0xba, 0xd7, 0x7f, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xda, 0x2c, 0xf5, 0x3c, 0x8a, 0x02, 0x00, 0x00,
}
