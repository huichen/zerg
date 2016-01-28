package zerg_client

import (
	"errors"
	"github.com/huichen/consistent_service"
	pb "github.com/huichen/zerg/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strings"
)

type ZergClient struct {
	endPoints         []string
	serviceName       string
	consistentService consistent_service.ConsistentService
	clients           map[string]pb.CrawlClient
	conns             map[string]*grpc.ClientConn
	initialized       bool
}

// endPoints: 逗号分隔的 etcd 接入点列表，每个接入点以 http:// 开始
func NewZergClient(endpoints string, servicename string) (*ZergClient, error) {
	// 解析 endPoints
	ep := strings.Split(endpoints, ",")
	if len(ep) == 0 {
		return nil, errors.New("无法解析endpoints")
	}

	zc := &ZergClient{
		endPoints:   ep,
		serviceName: servicename,
	}
	zc.clients = make(map[string]pb.CrawlClient)
	zc.conns = make(map[string]*grpc.ClientConn)
	err := zc.consistentService.Connect(servicename, ep)
	if err != nil {
		return nil, err
	}
	zc.initialized = true
	return zc, nil
}

func (zc *ZergClient) Crawl(in *pb.CrawlRequest, opts ...grpc.CallOption) (*pb.CrawlResponse, error) {
	// 检查是否已经初始化
	if !zc.initialized {
		return nil, errors.New("DistCrawlClient 没有初始化")
	}

	node, err := zc.consistentService.GetNode(in.Url)
	if err != nil {
		return nil, err
	}

	if _, ok := zc.conns[node]; !ok {
		conn, err := grpc.Dial(node, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		zc.conns[node] = conn
		client := pb.NewCrawlClient(conn)
		zc.clients[node] = client
	}

	return zc.clients[node].Crawl(context.Background(), in, opts...)
}

func (zc *ZergClient) Close() {
	for _, v := range zc.conns {
		v.Close()
	}
	zc.initialized = false
}
