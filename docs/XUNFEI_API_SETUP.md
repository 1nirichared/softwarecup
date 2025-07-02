# 讯飞星火X1 API 配置指南

## 📋 概述

讯飞星火X1是科大讯飞推出的大语言模型，具有优秀的中文理解和生成能力。本平台集成讯飞星火X1 API，支持流式对话功能，提供类似ChatGPT的实时对话体验。

## 🚀 主要特性

- **流式对话**：支持SSE (Server-Sent Events) 流式通信
- **实时回复**：AI回复实时显示，无需等待完整响应
- **中文优化**：专精中文理解和生成
- **教学场景**：适合教育领域的问答和辅导
- **成本合理**：相比其他AI服务，成本更加合理

## 🔧 配置步骤

### 1. 注册讯飞开放平台

1. 访问 [讯飞开放平台](https://www.xfyun.cn/)
2. 注册账号并完成实名认证
3. 进入控制台，创建新应用

### 2. 创建应用

1. 在控制台点击"创建应用"
2. 选择"星火认知大模型"
3. 填写应用信息：
   - 应用名称：智能教学平台
   - 应用描述：教育AI助手
   - 应用类型：Web应用

### 3. 获取密钥信息

创建应用后，在应用详情页面获取以下信息：
- **AppID**：应用唯一标识
- **APIKey**：API访问密钥
- **APISecret**：API访问密钥（用于签名）

### 4. 配置后端

#### 修改配置文件
编辑 `backend/config/config.yaml`：

```yaml
ai:
  provider: "xunfei"  # 设置为讯飞服务
  max_tokens: 2048
  temperature: 0.7

xunfei:
  app_id: "your_app_id_here"
  api_key: "your_api_key_here"
  api_secret: "your_api_secret_here"
  host: "spark-api.xf-yun.com"
  path: "/v3.1/chat"
```

#### 环境变量配置（可选）
```bash
export XUNFEI_APP_ID="your_app_id_here"
export XUNFEI_API_KEY="your_api_key_here"
export XUNFEI_API_SECRET="your_api_secret_here"
```

### 5. 测试配置

#### 使用测试脚本
```bash
cd scripts
python test_xunfei_api.py
```

#### 手动测试
```bash
curl -X POST http://localhost:3002/api/v1/chat/messages \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "content": "你好，请介绍一下自己",
    "type": "text"
  }'
```

## 🔐 签名算法

讯飞API使用HMAC-SHA256签名算法进行身份验证。

### 签名生成步骤

1. **获取RFC1123格式的时间戳**
```go
timeStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
```

2. **拼接签名字符串**
```go
signatureOrigin := "host: " + host + "\n"
signatureOrigin += "date: " + timeStr + "\n"
signatureOrigin += "GET " + path + " HTTP/1.1"
```

3. **使用APISecret进行HMAC-SHA256签名**
```go
h := hmac.New(sha256.New, []byte(apiSecret))
h.Write([]byte(signatureOrigin))
signatureSha := base64.StdEncoding.EncodeToString(h.Sum(nil))
```

4. **生成Authorization头**
```go
authorizationOrigin := `api_key="` + apiKey + `", algorithm="hmac-sha256", headers="host date request-line", signature="` + signatureSha + `"`
authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))
```

### 完整签名示例
```go
func generateSignature(host, path, apiKey, apiSecret string) (string, string, string) {
    // 获取RFC1123格式的时间戳
    timeStr := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
    
    // 拼接签名字符串
    signatureOrigin := "host: " + host + "\n"
    signatureOrigin += "date: " + timeStr + "\n"
    signatureOrigin += "GET " + path + " HTTP/1.1"
    
    // 使用APISecret进行HMAC-SHA256签名
    h := hmac.New(sha256.New, []byte(apiSecret))
    h.Write([]byte(signatureOrigin))
    signatureSha := base64.StdEncoding.EncodeToString(h.Sum(nil))
    
    // 生成Authorization头
    authorizationOrigin := `api_key="` + apiKey + `", algorithm="hmac-sha256", headers="host date request-line", signature="` + signatureSha + `"`
    authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))
    
    return authorization, signatureSha, timeStr
}
```

## 🌐 WebSocket连接

### 连接URL
```
wss://spark-api.xf-yun.com/v3.1/chat
```

### 请求头
```
Host: spark-api.xf-yun.com
Date: Mon, 02 Jan 2024 10:00:00 GMT
Authorization: base64_encoded_authorization
```

### 消息格式

#### 发送消息
```json
{
  "header": {
    "app_id": "your_app_id",
    "uid": "user_unique_id"
  },
  "parameter": {
    "chat": {
      "domain": "general",
      "temperature": 0.7,
      "max_tokens": 2048
    }
  },
  "payload": {
    "message": {
      "text": [
        {
          "role": "user",
          "content": "用户消息内容"
        }
      ]
    }
  }
}
```

#### 接收消息
```json
{
  "header": {
    "code": 0,
    "message": "Success",
    "sid": "session_id"
  },
  "payload": {
    "choices": {
      "status": 2,
      "seq": 0,
      "text": [
        {
          "content": "AI回复内容",
          "role": "assistant"
        }
      ]
    },
    "usage": {
      "text": {
        "question_tokens": 10,
        "prompt_tokens": 10,
        "completion_tokens": 20,
        "total_tokens": 30
      }
    }
  }
}
```

## 🔄 流式对话实现

### 后端实现

#### 1. WebSocket连接管理
```go
type XunfeiService struct {
    appID     string
    apiKey    string
    apiSecret string
    host      string
    path      string
}

func (s *XunfeiService) ChatStream(ctx context.Context, messages []Message, callback func(string)) error {
    // 生成签名
    authorization, _, date := generateSignature(s.host, s.path, s.apiKey, s.apiSecret)
    
    // 建立WebSocket连接
    u := url.URL{Scheme: "wss", Host: s.host, Path: s.path}
    conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{
        "Authorization": {authorization},
        "Date":         {date},
    })
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // 发送消息
    request := buildRequest(s.appID, messages)
    if err := conn.WriteJSON(request); err != nil {
        return err
    }
    
    // 接收流式回复
    for {
        var response Response
        if err := conn.ReadJSON(&response); err != nil {
            return err
        }
        
        // 处理回复
        if response.Header.Code != 0 {
            return fmt.Errorf("API错误: %s", response.Header.Message)
        }
        
        // 调用回调函数
        if len(response.Payload.Choices.Text) > 0 {
            content := response.Payload.Choices.Text[0].Content
            callback(content)
        }
        
        // 检查是否完成
        if response.Payload.Choices.Status == 2 {
            break
        }
    }
    
    return nil
}
```

#### 2. SSE流式接口
```go
func (h *ChatHandler) StreamAIChat(c *gin.Context) {
    // 获取参数
    sessionID := c.Query("session_id")
    token := c.Query("token")
    
    // 验证token
    userID, err := h.validateToken(token)
    if err != nil {
        c.JSON(401, gin.H{"error": "未授权"})
        return
    }
    
    // 设置SSE头
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("Access-Control-Allow-Origin", "*")
    
    // 获取会话消息
    messages, err := h.getSessionMessages(sessionID, userID)
    if err != nil {
        c.SSEvent("error", gin.H{"message": "获取会话失败"})
        return
    }
    
    // 调用AI服务
    err = h.aiService.ChatStream(c.Request.Context(), messages, func(content string) {
        c.SSEvent("content", gin.H{"content": content})
        c.Writer.Flush()
    })
    
    if err != nil {
        c.SSEvent("error", gin.H{"message": "AI服务错误"})
    } else {
        c.SSEvent("end", gin.H{"message": "回复完成"})
    }
}
```

### 前端实现

#### 1. EventSource连接
```javascript
function streamAIChat(sessionId, token) {
    const eventSource = new EventSource(
        `/api/v1/chat/stream?session_id=${sessionId}&token=${token}`
    );
    
    eventSource.onmessage = function(event) {
        const data = JSON.parse(event.data);
        
        switch(data.type) {
            case 'content':
                // 更新UI显示内容
                appendMessage(data.content);
                break;
            case 'end':
                // 完成处理
                eventSource.close();
                break;
            case 'error':
                // 错误处理
                console.error('错误:', data.message);
                eventSource.close();
                break;
        }
    };
    
    eventSource.onerror = function(error) {
        console.error('SSE连接错误:', error);
        eventSource.close();
    };
}
```

#### 2. 发送消息
```javascript
async function sendMessage(sessionId, content) {
    try {
        const response = await fetch('/api/v1/chat/messages', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                session_id: sessionId,
                content: content,
                type: 'text'
            })
        });
        
        if (response.ok) {
            // 开始流式接收回复
            streamAIChat(sessionId, token);
        }
    } catch (error) {
        console.error('发送消息失败:', error);
    }
}
```

## 📊 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 10000 | 输入参数错误 | 检查请求参数格式 |
| 10001 | 签名验证失败 | 检查签名算法和密钥 |
| 10002 | 时间戳过期 | 检查系统时间 |
| 10003 | 应用不存在 | 检查AppID是否正确 |
| 10004 | 应用未授权 | 检查应用权限设置 |
| 10005 | 服务调用失败 | 检查网络连接 |

### 错误处理示例
```go
func handleXunfeiError(code int, message string) error {
    switch code {
    case 10000:
        return fmt.Errorf("参数错误: %s", message)
    case 10001:
        return fmt.Errorf("签名验证失败: %s", message)
    case 10002:
        return fmt.Errorf("时间戳过期: %s", message)
    case 10003:
        return fmt.Errorf("应用不存在: %s", message)
    case 10004:
        return fmt.Errorf("应用未授权: %s", message)
    case 10005:
        return fmt.Errorf("服务调用失败: %s", message)
    default:
        return fmt.Errorf("未知错误(%d): %s", code, message)
    }
}
```

## 🔧 性能优化

### 1. 连接池管理
```go
type XunfeiConnectionPool struct {
    connections chan *websocket.Conn
    maxConn     int
    mu          sync.Mutex
}

func (p *XunfeiConnectionPool) GetConnection() (*websocket.Conn, error) {
    select {
    case conn := <-p.connections:
        return conn, nil
    default:
        // 创建新连接
        return p.createConnection()
    }
}

func (p *XunfeiConnectionPool) ReturnConnection(conn *websocket.Conn) {
    select {
    case p.connections <- conn:
        // 连接已归还
    default:
        // 连接池已满，关闭连接
        conn.Close()
    }
}
```

### 2. 消息缓存
```go
type MessageCache struct {
    cache map[string][]Message
    mu    sync.RWMutex
}

func (c *MessageCache) Get(sessionID string) []Message {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.cache[sessionID]
}

func (c *MessageCache) Set(sessionID string, messages []Message) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[sessionID] = messages
}
```

## 📈 监控和日志

### 1. 请求日志
```go
func (s *XunfeiService) logRequest(messages []Message, responseTime time.Duration) {
    log.Printf("讯飞API请求 - 消息数: %d, 响应时间: %v", 
        len(messages), responseTime)
}
```

### 2. 错误监控
```go
func (s *XunfeiService) monitorError(err error) {
    // 记录错误日志
    log.Printf("讯飞API错误: %v", err)
    
    // 发送错误告警
    if s.errorCount > 10 {
        // 发送告警通知
        s.sendAlert("讯飞API错误率过高")
    }
}
```

## 🔒 安全考虑

### 1. 密钥管理
- 不要在代码中硬编码密钥
- 使用环境变量或配置文件
- 定期轮换密钥

### 2. 请求限制
```go
type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func (r *RateLimiter) Allow(userID string) bool {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    now := time.Now()
    if requests, exists := r.requests[userID]; exists {
        // 清理过期的请求记录
        valid := make([]time.Time, 0)
        for _, req := range requests {
            if now.Sub(req) < r.window {
                valid = append(valid, req)
            }
        }
        r.requests[userID] = valid
        
        if len(valid) >= r.limit {
            return false
        }
    }
    
    r.requests[userID] = append(r.requests[userID], now)
    return true
}
```

## 📚 相关文档

- [讯飞开放平台文档](https://www.xfyun.cn/doc/spark/Web.html)
- [API文档](../API.md) - 后端API接口说明
- [部署指南](../DEPLOYMENT.md) - 部署配置说明
- [项目总结](../PROJECT_SUMMARY.md) - 项目功能总结

## 🤝 技术支持

如果在配置过程中遇到问题，请：

1. 检查讯飞开放平台的应用配置
2. 验证签名算法实现
3. 查看网络连接状态
4. 检查错误日志
5. 参考讯飞官方文档 