package main

import (
	"flag"
	pb "github.com/huichen/zerg/protos"
	"github.com/huichen/zerg/zerg_client"
	"golang.org/x/net/context"
	"log"
)

var (
	address     = flag.String("address", ":50051", "服务器地址")
	url         = flag.String("url", "", "URL")
	freq        = flag.Int64("freq", 0, "抓取频率")
	endPoints   = flag.String("endpoints", "", "半角逗号分隔的 etcd 接入点列表，每个接入点地址以 http:// 开始")
	serviceName = flag.String("service_name", "/services/zerg", "zerg 服务名")
)

func main() {
	flag.Parse()

	if *url == "" {
		log.Fatal("--url 参数不能为空")
	}

	zc, err := zerg_client.NewZergClient(*endPoints, *serviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer zc.Close()

	request := pb.CrawlRequest{Url: *url, Timeout: 10000, CrawlFrequency: *freq}
	log.Printf("开始抓取")
	for i := 0; i < 10; i++ {
		// 调用 client.Crawl 前必须先调用 Get 命令获取 client，client 通过 url 的一致性哈希进行分配
		client, err := zc.Get(*url)
		if err != nil {
			log.Fatal(err)
		}

		response, err := client.Crawl(context.Background(), &request)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", response.Metadata)
		log.Printf("%d", len(response.Content))
	}
}
