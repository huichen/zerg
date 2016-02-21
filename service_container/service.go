package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	pb "github.com/huichen/zerg/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	address = flag.String("address", ":50051", "服务器地址")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("无法绑定地址: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCrawlServer(s, &server{})
	s.Serve(lis)
}

type server struct {
	client *http.Client
}

func (s *server) Crawl(ctx context.Context, in *pb.CrawlRequest) (*pb.CrawlResponse, error) {
	return s.internalCrawl(in)
}

func (s *server) internalCrawl(in *pb.CrawlRequest) (*pb.CrawlResponse, error) {
	response := pb.CrawlResponse{}

	// 获取 http 连接
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	if in.Timeout > 0 {
		client.Timeout = time.Millisecond * time.Duration(in.Timeout)
	}

	// 根据不同的 method 类型，分别调用不同 HTTP 方法
	var resp *http.Response
	var err error
	if in.Method == pb.Method_GET {
		resp, err = client.Get(in.Url)
	} else if in.Method == pb.Method_HEAD {
		resp, err = client.Head(in.Url)
	} else if in.Method == pb.Method_POST {
		buff := bytes.NewBufferString(in.PostBody)
		resp, err = client.Post(in.Url, in.BodyType, buff)
	}
	if err != nil {
		return nil, err
	}

	// 只有当 method 不为 HEAD 时才读取页面内容
	var body []byte
	if in.Method != pb.Method_HEAD {
		// 读取页面内容
		var err error
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		// 充填 response
		if !in.OnlyReturnMetadata {
			response.Content = string(body)
		}
	}

	// 充填 metadata
	response.Metadata = &pb.Metadata{}
	response.Metadata.Length = uint32(len(body))
	for key, vs := range resp.Header {
		for _, v := range vs {
			response.Metadata.Header = append(response.Metadata.Header, &pb.KV{
				Key:   key,
				Value: v,
			})
		}
	}
	response.Metadata.Status = resp.Status
	response.Metadata.StatusCode = int32(resp.StatusCode)

	return &response, nil
}
