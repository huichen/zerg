package main

import (
	"flag"
	pb "github.com/huichen/zerg/protos"
	"github.com/huichen/zerg/zerg_client"
	"log"
)

var (
	url         = flag.String("url", "", "URL")
	ttl         = flag.Int64("ttl", 0, "重新抓取TTL")
	endPoints   = flag.String("endpoints", "", "半角逗号分隔的 etcd 接入点列表，每个接入点地址以 http:// 开始")
	serviceName = flag.String("service_name", "/services/zerg", "zerg 服务名")
	method      = flag.String("method", "GET", "HTTP 请求类型：GET HEAD POST POSTFORM")
)

func main() {
	flag.Parse()

	// 创建新 ZergClient
	zc, err := zerg_client.NewZergClient(*endPoints, *serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer zc.Close()

	// 调用 zerg 服务
	request := pb.CrawlRequest{
		Url:        *url,
		Timeout:    10000, // 超时 10 秒
		RecrawlTtl: *ttl,
		Method:     pb.Method(pb.Method_value[*method]),
	}
	response, err := zc.Crawl(&request)
	if err != nil {
		// 处理异常
		log.Fatal(err)
	}

	// 处理返回结果
	log.Printf("metadata = %+v", response.Metadata)
	log.Printf("page content length = %d", len(response.Content))
}
