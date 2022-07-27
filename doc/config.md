## 配置文件说明

```
# 日志配置
Log:

  # 日志级别：info、warning、error
  Level: "info"
  
  # 是否在控制台输出日志：0、1
  ExportConsole: 1
  
  # 输出日志到日志文件
  ExportFile:

    # 日志文件路径
    Path: "./run.log"


# 云主机 Frontend 服务配置
Frontend:

  # 监听 HTTP 端口，用户访问的端口，可以直接设置为 80、443，如果前面还有 Nginx，则换其他端口
  HttpPort: 8112
  
  # HTTP 请求的超时时间（秒）
  HttpTimeout: 30
  
  # SSL 证书（可选）
  HttpsCertFile: "....../*.pem"
  
  # SSL 私钥（可选）
  HttpsKeyFile: "....../*.key"
  
  # 监听 Websocket 端口，Proxy 将通过该端口和 Frontend 建立长连接
  WebsocketPort: 8113
  
  # Websocket 的路劲
  WebsocketPath: "/proxy"
  
  # Websocket 连接的密钥，防止被恶意连接，请自定改为复杂密钥
  SecretKey: "123456"


# 内网主机 Proxy 服务配置
Proxy:

  # Frontend 服务的连接地址
  FrontendUrl: "ws://127.0.0.1:8113/proxy"
  
  # Backend 服务的 IP 或域名
  BackendHost: "http://127.0.0.1:8114"

  # 和 Frontend 建立 Websocket 连接的数量，测试证明多个连接比单个连接有更好的并发性能
  WebsocketNum: 3
  
  # 和 Frontend 建立 Websocket 连接的密钥
  SecretKey: "123456"


# 调试时，可以启动一个模拟的 Backend 服务
Backend:
  HttpPort: 8114
```


