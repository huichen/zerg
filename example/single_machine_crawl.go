package main

import (
	"flag"
	pb "github.com/huichen/zerg/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var (
	address = flag.String("address", ":50051", "服务器地址")
	url     = flag.String("url", "", "URL")
	ttl     = flag.Int64("ttl", 0, "重新抓取TTL")
	method  = flag.String("method", "GET", "HTTP 请求类型：GET HEAD POST POSTFORM")
)

func main() {
	flag.Parse()

	// 得到 CrawlClient
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewCrawlClient(conn)

	log.Printf("开始抓取")
	request := pb.CrawlRequest{
		Url:        *url,
		Timeout:    10000,
		RecrawlTtl: *ttl,
		Method:     pb.Method(pb.Method_value[*method]),
	}
	response, err := client.Crawl(context.Background(), &request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("抓取完毕")
	log.Printf("%+v", response.Metadata)
	log.Printf("%d", len(response.Content))
}
