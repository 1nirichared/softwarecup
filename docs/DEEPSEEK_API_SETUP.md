# DeepSeek API 配置指南

## 📋 概述

DeepSeek是专注于代码生成和编程教学的大语言模型，具有优秀的代码理解和生成能力。本平台集成DeepSeek API作为备用AI服务，特别适合编程教学场景。

## 🚀 主要特性

- **代码生成**：优秀的代码生成和补全能力
- **编程教学**：专精编程语言教学和解释
- **多语言支持**：支持Python、Java、C++、JavaScript等主流编程语言
- **成本优势**：相比其他AI服务，成本更加合理
- **稳定可靠**：API服务稳定，响应速度快

## 🔧 配置步骤

### 1. 注册DeepSeek平台

1. 访问 [DeepSeek官网](https://www.deepseek.com/)
2. 注册账号并完成邮箱验证
3. 进入API控制台，获取API密钥

### 2. 获取API密钥

1. 登录DeepSeek控制台
2. 进入"API Keys"页面
3. 点击"Create API Key"创建新的API密钥
4. 复制并保存API密钥（注意：密钥只显示一次）

### 3. 配置后端

#### 修改配置文件
编辑 `backend/config/config.yaml`：

```yaml
ai:
  provider: "deepseek"  # 设置为DeepSeek服务
  max_tokens: 2048
  temperature: 0.7

deepseek:
  api_key: "your_deepseek_api_key_here"
  base_url: "https://api.deepseek.com/v1"
  model: "deepseek-chat"
  timeout: 60
```

#### 环境变量配置（可选）
```bash
export DEEPSEEK_API_KEY="your_deepseek_api_key_here"
export DEEPSEEK_BASE_URL="https://api.deepseek.com/v1"
export DEEPSEEK_MODEL="deepseek-chat"
```

### 4. 测试配置

#### 使用测试脚本
```bash
cd scripts
python test_deepseek_api.py
```

#### 手动测试
```bash
curl -X POST http://localhost:3002/api/v1/chat/messages \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "content": "请用Python写一个简单的计算器函数",
    "type": "text"
  }'
```

## 🌐 API接口

### 基础聊天接口
```
POST https://api.deepseek.com/v1/chat/completions
```

### 请求头
```
Authorization: Bearer your_api_key
Content-Type: application/json
```

### 请求格式
```json
{
  "model": "deepseek-chat",
  "messages": [
    {
      "role": "user",
      "content": "用户消息内容"
    }
  ],
  "max_tokens": 2048,
  "temperature": 0.7,
  "stream": false
}
```

### 响应格式
```json
{
  "id": "chatcmpl-123",
  "object": "chat.completion",
  "created": 1677652288,
  "model": "deepseek-chat",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "AI回复内容"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

## 🔄 流式对话实现

### 后端实现

#### 1. HTTP流式请求
```go
type DeepSeekService struct {
    apiKey   string
    baseURL  string
    model    string
    client   *http.Client
}

func (s *DeepSeekService) ChatStream(ctx context.Context, messages []Message, callback func(string)) error {
    // 构建请求
    request := DeepSeekRequest{
        Model:       s.model,
        Messages:    s.convertMessages(messages),
        MaxTokens:   2048,
        Temperature: 0.7,
        Stream:      true,
    }
    
    // 发送请求
    reqBody, _ := json.Marshal(request)
    req, _ := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
    req.Header.Set("Authorization", "Bearer "+s.apiKey)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := s.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // 处理流式响应
    reader := bufio.NewReader(resp.Body)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        
        // 解析SSE格式
        if strings.HasPrefix(line, "data: ") {
            data := strings.TrimPrefix(line, "data: ")
            if data == "[DONE]" {
                break
            }
            
            var streamResp DeepSeekStreamResponse
            if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
                continue
            }
            
            // 调用回调函数
            if len(streamResp.Choices) > 0 && len(streamResp.Choices[0].Delta.Content) > 0 {
                callback(streamResp.Choices[0].Delta.Content)
            }
        }
    }
    
    return nil
}
```

#### 2. 消息格式转换
```go
func (s *DeepSeekService) convertMessages(messages []Message) []DeepSeekMessage {
    result := make([]DeepSeekMessage, len(messages))
    for i, msg := range messages {
        result[i] = DeepSeekMessage{
            Role:    msg.Role,
            Content: msg.Content,
        }
    }
    return result
}
```

### 前端实现

#### 1. 流式接收
```javascript
async function streamDeepSeekChat(messages) {
    try {
        const response = await fetch('/api/v1/chat/stream', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                messages: messages,
                provider: 'deepseek'
            })
        });
        
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            
            const chunk = decoder.decode(value);
            const lines = chunk.split('\n');
            
            for (const line of lines) {
                if (line.startsWith('data: ')) {
                    const data = JSON.parse(line.slice(6));
                    if (data.content) {
                        // 更新UI显示内容
                        appendMessage(data.content);
                    }
                }
            }
        }
    } catch (error) {
        console.error('流式对话错误:', error);
    }
}
```

## 📊 编程教学场景

### 1. 代码生成
```go
func (s *DeepSeekService) GenerateCode(language, description string) (string, error) {
    prompt := fmt.Sprintf(`请用%s语言编写代码：%s

要求：
1. 代码要简洁清晰
2. 添加必要的注释
3. 遵循最佳实践
4. 提供使用示例`, language, description)
    
    return s.chatCompletion(prompt)
}
```

### 2. 代码解释
```go
func (s *DeepSeekService) ExplainCode(code, language string) (string, error) {
    prompt := fmt.Sprintf(`请解释以下%s代码：

```%s
%s
```

请从以下方面进行解释：
1. 代码功能
2. 关键语法
3. 执行流程
4. 可能的改进`, language, language, code)
    
    return s.chatCompletion(prompt)
}
```

### 3. 代码调试
```go
func (s *DeepSeekService) DebugCode(code, language, error string) (string, error) {
    prompt := fmt.Sprintf(`请帮助调试以下%s代码：

```%s
%s
```

错误信息：
%s

请：
1. 分析错误原因
2. 提供修复方案
3. 解释修复原理`, language, language, code, error)
    
    return s.chatCompletion(prompt)
}
```

### 4. 编程练习生成
```go
func (s *DeepSeekService) GenerateExercise(language, topic, difficulty string) (string, error) {
    prompt := fmt.Sprintf(`请为%s语言生成一道关于"%s"的%s难度编程练习。

要求：
1. 题目描述清晰
2. 提供输入输出示例
3. 包含解题思路
4. 提供参考答案
5. 说明涉及的知识点`, language, topic, difficulty)
    
    return s.chatCompletion(prompt)
}
```

## 📈 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求格式和参数 |
| 401 | 认证失败 | 检查API密钥是否正确 |
| 403 | 权限不足 | 检查API密钥权限 |
| 429 | 请求频率限制 | 降低请求频率 |
| 500 | 服务器错误 | 稍后重试 |

### 错误处理示例
```go
func handleDeepSeekError(resp *http.Response) error {
    if resp.StatusCode != 200 {
        var errorResp DeepSeekErrorResponse
        if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
            return fmt.Errorf("HTTP %d: 未知错误", resp.StatusCode)
        }
        
        switch resp.StatusCode {
        case 400:
            return fmt.Errorf("请求参数错误: %s", errorResp.Error.Message)
        case 401:
            return fmt.Errorf("认证失败: %s", errorResp.Error.Message)
        case 403:
            return fmt.Errorf("权限不足: %s", errorResp.Error.Message)
        case 429:
            return fmt.Errorf("请求频率限制: %s", errorResp.Error.Message)
        case 500:
            return fmt.Errorf("服务器错误: %s", errorResp.Error.Message)
        default:
            return fmt.Errorf("API错误(%d): %s", resp.StatusCode, errorResp.Error.Message)
        }
    }
    return nil
}
```

## 🔧 性能优化

### 1. 连接池管理
```go
type DeepSeekClient struct {
    client  *http.Client
    apiKey  string
    baseURL string
}

func NewDeepSeekClient(apiKey, baseURL string) *DeepSeekClient {
    return &DeepSeekClient{
        client: &http.Client{
            Timeout: 60 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        apiKey:  apiKey,
        baseURL: baseURL,
    }
}
```

### 2. 请求缓存
```go
type DeepSeekCache struct {
    cache map[string]string
    mu    sync.RWMutex
    ttl   time.Duration
}

func (c *DeepSeekCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if value, exists := c.cache[key]; exists {
        return value, true
    }
    return "", false
}

func (c *DeepSeekCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = value
    
    // 设置过期时间
    time.AfterFunc(c.ttl, func() {
        c.mu.Lock()
        delete(c.cache, key)
        c.mu.Unlock()
    })
}
```

### 3. 并发控制
```go
type DeepSeekLimiter struct {
    semaphore chan struct{}
    timeout   time.Duration
}

func NewDeepSeekLimiter(maxConcurrent int, timeout time.Duration) *DeepSeekLimiter {
    return &DeepSeekLimiter{
        semaphore: make(chan struct{}, maxConcurrent),
        timeout:   timeout,
    }
}

func (l *DeepSeekLimiter) Do(fn func() error) error {
    select {
    case l.semaphore <- struct{}{}:
        defer func() { <-l.semaphore }()
        return fn()
    case <-time.After(l.timeout):
        return fmt.Errorf("请求超时")
    }
}
```

## 📊 监控和日志

### 1. 请求监控
```go
type DeepSeekMetrics struct {
    requestCount   int64
    errorCount     int64
    totalLatency   time.Duration
    mu             sync.Mutex
}

func (m *DeepSeekMetrics) RecordRequest(duration time.Duration, err error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    atomic.AddInt64(&m.requestCount, 1)
    m.totalLatency += duration
    
    if err != nil {
        atomic.AddInt64(&m.errorCount, 1)
    }
}

func (m *DeepSeekMetrics) GetStats() map[string]interface{} {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    avgLatency := time.Duration(0)
    if m.requestCount > 0 {
        avgLatency = m.totalLatency / time.Duration(m.requestCount)
    }
    
    return map[string]interface{}{
        "total_requests":   m.requestCount,
        "error_count":      m.errorCount,
        "error_rate":       float64(m.errorCount) / float64(m.requestCount),
        "avg_latency":      avgLatency,
    }
}
```

### 2. 日志记录
```go
func (s *DeepSeekService) logRequest(messages []Message, response string, duration time.Duration) {
    log.Printf("DeepSeek API请求 - 消息数: %d, 响应长度: %d, 耗时: %v", 
        len(messages), len(response), duration)
}
```

## 🔒 安全考虑

### 1. API密钥管理
```go
type DeepSeekConfig struct {
    APIKey  string `json:"api_key" validate:"required"`
    BaseURL string `json:"base_url" validate:"required,url"`
    Model   string `json:"model" validate:"required"`
}

func LoadDeepSeekConfig() (*DeepSeekConfig, error) {
    // 优先从环境变量读取
    apiKey := os.Getenv("DEEPSEEK_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("DEEPSEEK_API_KEY环境变量未设置")
    }
    
    return &DeepSeekConfig{
        APIKey:  apiKey,
        BaseURL: getEnvOrDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
        Model:   getEnvOrDefault("DEEPSEEK_MODEL", "deepseek-chat"),
    }, nil
}
```

### 2. 输入验证
```go
func (s *DeepSeekService) validateInput(content string) error {
    if len(content) > 10000 {
        return fmt.Errorf("输入内容过长，最大支持10000字符")
    }
    
    // 检查敏感内容
    sensitiveWords := []string{"恶意内容", "攻击代码"}
    for _, word := range sensitiveWords {
        if strings.Contains(content, word) {
            return fmt.Errorf("输入内容包含敏感信息")
        }
    }
    
    return nil
}
```

### 3. 请求限制
```go
type DeepSeekRateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func (r *DeepSeekRateLimiter) Allow(userID string) bool {
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

- [DeepSeek官方文档](https://platform.deepseek.com/docs)
- [API文档](../API.md) - 后端API接口说明
- [部署指南](../DEPLOYMENT.md) - 部署配置说明
- [项目总结](../PROJECT_SUMMARY.md) - 项目功能总结

## 🤝 技术支持

如果在配置过程中遇到问题，请：

1. 检查DeepSeek API密钥是否正确
2. 验证网络连接状态
3. 查看API响应错误信息
4. 检查请求格式和参数
5. 参考DeepSeek官方文档 