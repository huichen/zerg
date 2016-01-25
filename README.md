# zerg
基于docker的分布式爬虫服务

![](https://raw.github.com/huichen/zerg/master/doc/zerg.png)

### 特性

* 多机多 IP，充分利用 IP 资源
* 服务自动发现和注册（基于 etcd 和 registrator）
* 负载均衡 + 一致性哈希
* 服务端客户端通信基于 gRPC
* 可设置抓取超时，页面重新抓取频率等
* 基于 docker volume 和文件系统的持久化存储

### 如何部署

#### 第一步：配置 etcd

我把 etcd 容器化了，并开发了脚本使得部署 etcd cluster 非常容易，见 [github.com/huichen/etcd_docker](https://github.com/huichen/etcd_docker)

为了容灾，请至少在三台服务器上运行 etcd 实例。为了方便调用，你可以固定 etc endpoint 的端口号，并在所有机器上手工添加 etcd host 的hostname。

#### 第二步：启动 registrator 服务发现程序

你需要在集群的每一台服务器上都运行 registrator，这使得我们可以自动发现和注册 dist_crawl 服务

```
docker run -d --name=registrator --net=host --volume=/var/run/docker.sock:/tmp/docker.sock \
  gliderlabs/registrator etcd://<etcd 接入点的 ip:port>/services
```

请把上面的 etcd 接入点换成你的 etcd 地址。

#### 第三部：部署 zerg 服务

进入 service_container 子目录，然后运行

```
./build_docker_image.sh
```

这会生成 unmerged/zerg 容器。然后在集群的每台服务器上启动容器：

```
docker run -d -P -v /opt/zerg_cache:/cache unmerged/zerg
```

registrator 会自动注册这些服务到 etcd。如果单机有多个IP，你可以单机启动多个容器，并在 -P 中分别指定IP。

抓取的页面内容会通过 docker volume 存储在 /opt/zerg_cache 目录下。

#### 第四部：调用样例代码

进入 examples 目录，运行

```
go run zerg_crawl.go --endpoints http://<你的 etcd host:ip> --url http://taobao.com
```

#### 可选步骤

1、重新生成 protobuf service

```
protoc protos/crawl.proto --go_out=plugins=grpc:protos -I protos/
```

2、本地测试

启动本地服务。进入 service_container 目录，然后运行

```
go run service.go
```

然后进入 examples 目录，运行

```
go run single_machine_crawl.go --url http://taobao.com
```
