# WebHole 内网穿透代理

## 基本概念

**Client**：<br>
用户端。

**Frontend**：<br>
作为 Websocket 服务端、HTTP 前端服务端，直接接受用户端的 HTTP 请求，但是自身并不需要处理任何业务逻辑，而是将请求转发至 Proxy，并等待 Proxy 的响应，然后响应给用户端。Frontend 是拥有公网 IP，或许还解析了域名的云主机。

**Proxy**：<br>
作为 Websocket 客户端、HTTP 代理，主动与 Frontend 建立 Websocket 长连接，接受 Frontend 转发过来的请求，然后将请求代理到真正的业务后端 Backend，等待 Backend 的响应再返回给 Frontend。Proxy 是公司内网主机，或者直接就是家里的 PC，没有固定公网 IP。

**Backend**：<br>
HTTP 后端服务端。Backend 可以是任何环境的 HTTP 服务器，可以是本地 PC、内网主机、其他公网主机、或者是一个负载均衡，只要是 Proxy 能够访问到的 HTTP 服务都可以。

## 代理流程

<img src="https://github.com/njutsiang/web-hole/raw/main/doc/process.png">

## 能解决什么问题？

1、可以将你的网站部署在公司内网，或者你的 PC 上，并且实现高效低延迟的访问，只需要一台低配的云主机，无需再购买云数据库、云存储等等，能节省不少费用。

2、极大地方便开发阶段的调试，例如你正在开发一个有支付功能的新项目，支付结果需要通过公网回调到你的项目，但是你的项目还非常简陋，甚至可能会报错，需要调试，无法部署到公网环境。这种情况就是 WebHole 能发挥作用的时候了，将公网请求代理到你正在开发的本地项目中，实现在线 debug。

3、如果你的家里有部署 NAS，也或许只需要一台旧电脑，给它插上大硬盘，就可以将你收藏多年的“学习资料”分享给你的水友们了。

## TODO LIST

- [x] 解决 Websocket 并发写的问题
- [x] 解决等待响应的 ChanMap 并发读写的问题
- [x] Frontend 支持多个 Proxy 服务
- [x] 优化日志组件、日志级别
- [x] Frontend 支持 https
- [x] 解决发送心跳失败和普通消息存在并发写的问题
- [ ] 支持 301、302 跳转的代理
- [ ] Proxy 支持多进程多连接
- [ ] 完善使用说明文档

## 压力测试

在相同条件下，对未经过 WebHole 代理和经过 WebHole 代理分别做了压力测试，测试参数如下：

```
# -z 60s，指定持续 60s
# -c 4，指定 4 个 Worker 并发执行请求
# -q 1000，指定每个 Worker 的请求速率（QPS）为 1000/s
hey -z 60s -c 4 -q 1000 -m GET http://127.0.0.1:8112
```

未经过 WebHole 代理的压力测试结果：

```
Summary:
  Total:        60.0017 secs
  Slowest:      0.0077 secs
  Fastest:      0.0000 secs
  Average:      0.0003 secs
  Requests/sec: 3105.0432

  Total data:   2794620 bytes
  Size/request: 15 bytes

Response time histogram:
  0.000 [1]     |
  0.001 [183847]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [1865]  |
  0.002 [358]   |
  0.003 [121]   |
  0.004 [55]    |
  0.005 [35]    |
  0.005 [14]    |
  0.006 [0]     |
  0.007 [1]     |
  0.008 [11]    |

Latency distribution:
  10% in 0.0002 secs
  25% in 0.0002 secs
  50% in 0.0003 secs
  75% in 0.0003 secs
  90% in 0.0005 secs
  95% in 0.0006 secs
  99% in 0.0009 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0004 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0040 secs
  resp wait:    0.0002 secs, 0.0000 secs, 0.0077 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0076 secs

Status code distribution:
  [200] 186308 responses
```

经过 WebHole 代理的压力测试结果：

```
Summary:
  Total:        60.0022 secs
  Slowest:      0.0228 secs
  Fastest:      0.0002 secs
  Average:      0.0018 secs
  Requests/sec: 1688.1889

  Total data:   1519425 bytes
  Size/request: 15 bytes

Response time histogram:
  0.000 [1]     |
  0.002 [88670] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [6198]  |■■■
  0.007 [3580]  |■■
  0.009 [2006]  |■
  0.011 [546]   |
  0.014 [185]   |
  0.016 [68]    |
  0.018 [21]    |
  0.021 [13]    |
  0.023 [7]     |

Latency distribution:
  10% in 0.0008 secs
  25% in 0.0010 secs
  50% in 0.0012 secs
  75% in 0.0017 secs
  90% in 0.0032 secs
  95% in 0.0055 secs
  99% in 0.0089 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0008 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0061 secs
  resp wait:    0.0017 secs, 0.0002 secs, 0.0227 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0075 secs

Status code distribution:
  [200] 101295 responses
```
