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
	freq    = flag.Int64("freq", 0, "抓取频率")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewZergClient(conn)

	request := pb.CrawlRequest{Url: *url, Timeout: 10000, CrawlFrequency: *freq}
	log.Printf("开始抓取")
	for i := 0; i < 10; i++ {
		response, err := client.Crawl(context.Background(), &request)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%+v", response.Metadata)
		log.Printf("%d", len(response.Content))
	}
}
