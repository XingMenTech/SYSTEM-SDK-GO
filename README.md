# SYSTEM-SDK-GO

Go语言实现的系统SDK，提供HTTP和gRPC两种通信协议的客户端支持，包含AES加密解密功能。

## 功能特性

- 🚀 **双协议支持**：同时支持HTTP REST API和gRPC协议调用
- 🔒 **数据安全**：内置AES加密解密功能，保障数据传输安全
- 📦 **易于集成**：简洁的API设计，快速集成到现有项目中
- 🧪 **完整测试**：包含完整的单元测试用例
- 📝 **类型安全**：基于Protocol Buffer定义的数据结构

## 目录结构

```
SYSTEM-SDK-GO/
├── pb/                 # Protocol Buffer相关文件
│   ├── rpcpb.proto     # gRPC服务定义
│   ├── rpcpb.pb.go     # 生成的Protocol Buffer代码
│   ├── rpcpb_grpc.pb.go # 生成的gRPC客户端代码
│   └── shell.sh        # protobuf编译脚本
├── aes.go             # AES加密解密实现
├── client_http.go     # HTTP客户端实现
├── client_grpc.go     # gRPC客户端实现
├── models.go          # 数据模型定义
├── http_api_test.go   # HTTP API测试用例
├── grpc_api_test.go   # gRPC API测试用例
├── go.mod             # Go模块定义
└── README.md          # 项目说明文档
```

## 安装使用

### 前置要求

- Go版本 >= 1.22.10
- gRPC相关依赖（如使用gRPC功能）

### 安装

```bash
go get github.com/your-org/SYSTEM-SDK-GO
```

或者在项目中直接导入：

```go
import "github.com/your-org/SYSTEM-SDK-GO"
```

## 快速开始

### HTTP客户端使用示例

```go
package main

import (
    "fmt"
    "log"
    sdk "github.com/your-org/SYSTEM-SDK-GO"
)

func main() {
    // 创建HTTP客户端
    client := sdk.NewHTTPClient("http://localhost:8080")
    
    // 发送请求
    response, err := client.SendRequest("POST", "/api/data", []byte(`{"message":"hello"}`))
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %s\n", string(response))
}
```

### gRPC客户端使用示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    sdk "github.com/your-org/SYSTEM-SDK-GO"
)

func main() {
    // 创建gRPC客户端
    client, err := sdk.NewGRPCClient("localhost:50051")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 调用gRPC方法
    ctx := context.Background()
    response, err := client.CallMethod(ctx, &sdk.Request{
        Data: []byte("hello world"),
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %s\n", string(response.Data))
}
```

### AES加密解密示例

```go
package main

import (
    "fmt"
    "log"
    sdk "github.com/your-org/SYSTEM-SDK-GO"
)

func main() {
    // 创建AES实例
    aes, err := sdk.NewAES([]byte("your-32-byte-key-here"))
    if err != nil {
        log.Fatal(err)
    }
    
    // 加密数据
    plaintext := []byte("敏感数据")
    ciphertext, err := aes.Encrypt(plaintext)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("加密后: %x\n", ciphertext)
    
    // 解密数据
    decrypted, err := aes.Decrypt(ciphertext)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("解密后: %s\n", string(decrypted))
}
```

## API参考

### HTTP客户端

#### `NewHTTPClient(baseURL string) *HTTPClient`
创建新的HTTP客户端实例

参数：
- `baseURL`: 服务基础URL

返回：
- `*HTTPClient`: HTTP客户端实例

#### `SendRequest(method, path string, body []byte) ([]byte, error)`
发送HTTP请求

参数：
- `method`: HTTP方法 (GET, POST, PUT, DELETE等)
- `path`: 请求路径
- `body`: 请求体数据

返回：
- `[]byte`: 响应数据
- `error`: 错误信息

### gRPC客户端

#### `NewGRPCClient(address string) (*GRPCClient, error)`
创建新的gRPC客户端实例

参数：
- `address`: gRPC服务地址

返回：
- `*GRPCClient`: gRPC客户端实例
- `error`: 错误信息

#### `CallMethod(ctx context.Context, req *Request) (*Response, error)`
调用gRPC方法

参数：
- `ctx`: 上下文
- `req`: 请求数据

返回：
- `*Response`: 响应数据
- `error`: 错误信息

### AES加密

#### `NewAES(key []byte) (*AES, error)`
创建新的AES实例

参数：
- `key`: 32字节密钥

返回：
- `*AES`: AES实例
- `error`: 错误信息

#### `Encrypt(plaintext []byte) ([]byte, error)`
加密数据

参数：
- `plaintext`: 明文数据

返回：
- `[]byte`: 密文数据
- `error`: 错误信息

#### `Decrypt(ciphertext []byte) ([]byte, error)`
解密数据

参数：
- `ciphertext`: 密文数据

返回：
- `[]byte`: 明文数据
- `error`: 错误信息

## 测试

运行所有测试：

```bash
go test -v ./...
```

运行特定测试：

```bash
# 运行HTTP API测试
go test -v -run TestHTTP

# 运行gRPC API测试  
go test -v -run TestGRPC

# 运行AES测试
go test -v -run TestAES
```

## 开发指南

### Protocol Buffer编译

如果修改了`.proto`文件，需要重新编译：

```bash
cd pb
./shell.sh
```

### 依赖管理

```bash
# 添加新依赖
go get github.com/some/package

# 清理未使用的依赖
go mod tidy

# 更新所有依赖到最新版本
go get -u ./...
```

## 配置说明

### 环境变量

```bash
# HTTP服务地址
HTTP_SERVICE_URL=http://localhost:8080

# gRPC服务地址  
GRPC_SERVICE_ADDR=localhost:50051

# AES密钥（32字节）
AES_KEY=your-32-byte-secret-key-here
```

## 错误处理

SDK遵循Go的标准错误处理模式，所有可能出错的方法都会返回`error`类型：

```go
response, err := client.SendRequest("GET", "/api/data", nil)
if err != nil {
    // 处理错误
    if errors.Is(err, sdk.ErrConnectionFailed) {
        // 连接失败处理
    } else if errors.Is(err, sdk.ErrTimeout) {
        // 超时处理
    }
    return err
}
```

## 性能优化建议

1. **连接复用**：HTTP客户端会自动复用连接，建议在应用生命周期内复用客户端实例
2. **超时设置**：为gRPC调用设置合适的超时时间
3. **批量操作**：对于大量数据操作，考虑使用批量接口
4. **缓存策略**：对频繁访问但不常变化的数据实施缓存

## 安全注意事项

1. **密钥管理**：AES密钥应该安全存储，避免硬编码在源码中
2. **传输安全**：生产环境建议使用HTTPS/gRPC TLS
3. **输入验证**：对外部输入进行适当验证和清理
4. **日志记录**：避免在日志中记录敏感信息

## 版本兼容性

- v1.0.0: 初始版本，支持基本HTTP/gRPC功能
- v1.1.0: 添加AES加密功能
- v1.2.0: 改进错误处理和测试覆盖

## 贡献指南

欢迎提交Issue和Pull Request：

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情

## 技术支持

如有问题，请通过以下方式联系：

- 提交GitHub Issue
- 发送邮件至：support@your-org.com
- 查看官方文档

---
*最后更新：2026年3月*