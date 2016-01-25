package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/huichen/zerg/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	pageCacheDir = flag.String("page_cache_dir", "cache", "页面缓存目录")
	address      = flag.String("address", ":50051", "服务器地址")
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

	// 文件名
	hasher := sha1.New()
	hasher.Write([]byte(in.Url))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	cacheFilename := fmt.Sprintf("%s/page-%s", *pageCacheDir, sha)
	metadataFilename := fmt.Sprintf("%s/metadata-%s", *pageCacheDir, sha)

	// 检查页面的 metadata 和 page 是否已经被抓取
	// 如果有任何异常或者页面过旧，则重新抓取
	if _, err := os.Stat(metadataFilename); err == nil {
		serializedMetadata, err := ioutil.ReadFile(metadataFilename)
		if err == nil {
			response.Metadata = &pb.Metadata{}
			err = proto.Unmarshal(serializedMetadata, response.Metadata)
			if err == nil {
				if in.OnlyReturnMetadata {
					return &response, nil
				}
				if time.Now().UnixNano()-response.Metadata.LastCrawlTimestamp < in.CrawlFrequency*int64(time.Millisecond) {
					if _, err := os.Stat(cacheFilename); err == nil {
						content, err := ioutil.ReadFile(cacheFilename)
						if err == nil {
							contentReader := bytes.NewReader(content)
							reader, err := gzip.NewReader(contentReader)
							if err == nil {
								if b, err := ioutil.ReadAll(reader); err == nil {
									response.Content = string(b)
									return &response, nil
								}
								reader.Close()
							}
						}
					}
				}
			}
		}
	}

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
	resp, err := client.Get(in.Url)
	if err != nil {
		return nil, err
	}

	// 读取页面内容
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// 充填 response
	if !in.OnlyReturnMetadata {
		response.Content = string(body)
	}
	response.IsFreshCrawl = true

	// 充填 metadata
	response.Metadata = &pb.Metadata{}
	response.Metadata.LastCrawlTimestamp = time.Now().UnixNano()
	response.Metadata.Length = uint32(len(body))

	// 将 page content 写入文件
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	writer.Write(body)
	writer.Close()
	if err := ioutil.WriteFile(cacheFilename, buf.Bytes(), 0644); err != nil {
		return &response, err
	}

	// 将 metadata 写入文件
	serializedMetadata, err := proto.Marshal(response.Metadata)
	if err != nil {
		return &response, err
	}
	if err := ioutil.WriteFile(metadataFilename, serializedMetadata, 0644); err != nil {
		return &response, err
	}

	return &response, nil
}
