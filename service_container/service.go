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
	"net/url"
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

	// 仅当 method 为 GET 或者 HEAD 时才读取文件缓存
	var cacheFilename, metadataFilename string
	if in.Method == pb.Method_GET || in.Method == pb.Method_HEAD {
		// 生成缓存文件名
		hasher := sha1.New()
		hasher.Write([]byte(in.Url))
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		cacheFilename = fmt.Sprintf("%s/page-%s", *pageCacheDir, sha)
		metadataFilename = fmt.Sprintf("%s/metadata-%s", *pageCacheDir, sha)

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
					if time.Now().UnixNano()-response.Metadata.LastCrawlTimestamp <
						in.RecrawlTtl*int64(time.Millisecond) {
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
	} else if in.Method == pb.Method_POSTFORM {
		formKeys := url.Values{}
		for _, kv := range in.FormValues {
			formKeys.Add(kv.Key, kv.Value)
		}
		resp, err = client.PostForm(in.Url, formKeys)
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
		response.IsFreshCrawl = true
	}

	// 充填 metadata
	response.Metadata = &pb.Metadata{}
	response.Metadata.LastCrawlTimestamp = time.Now().UnixNano()
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

	// 仅当 method 为 GET 才将页面内容写入文件
	if in.Method == pb.Method_GET {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		writer.Write(body)
		writer.Close()
		if err := ioutil.WriteFile(cacheFilename, buf.Bytes(), 0644); err != nil {
			return &response, err
		}
	}

	// 仅当 method 为 GET 或者 HEAD 时才将 metadata 写入文件
	if in.Method == pb.Method_GET || in.Method == pb.Method_HEAD {
		serializedMetadata, err := proto.Marshal(response.Metadata)
		if err != nil {
			return &response, err
		}
		if err := ioutil.WriteFile(metadataFilename, serializedMetadata, 0644); err != nil {
			return &response, err
		}
	}

	return &response, nil
}
