# 讯飞星火 API 配置指南

## 概述

本指南将帮助您配置讯飞星火 API 服务，将其集成到教学培训平台中。讯飞星火是科大讯飞推出的大语言模型，具有优秀的中文理解和生成能力。

## 1. 获取 API 凭据

### 注册账号
1. 访问 [讯飞开放平台](https://www.xfyun.cn/)
2. 注册并登录账号
3. 完成实名认证

### 创建应用
1. 进入 [控制台](https://console.xfyun.cn/)
2. 创建新应用
3. 选择"星火认知大模型"服务
4. 获取以下凭据：
   - **AppID**: 应用ID
   - **APISecret**: API密钥
   - **APIKey**: API Key

## 2. 配置项目

### 更新配置文件

编辑 `backend/config/config.yaml`：

```yaml
# AI服务配置
ai:
  # 主要AI提供商
  provider: "xunfei"  # xunfei, deepseek, openai
  
  # DeepSeek API配置
  deepseek_api_key: "your-deepseek-api-key"
  deepseek_base_url: "https://api.deepseek.com/v1"
  deepseek_model: "deepseek-coder"
  
  # OpenAI配置（备用）
  openai_api_key: "your-openai-api-key"
  openai_base_url: "https://api.openai.com/v1"
  openai_model: "gpt-3.5-turbo"
  
  # 通用配置
  max_tokens: 2000
  temperature: 0.7
  timeout: 60

# 讯飞星火配置
xunfei:
  app_id: "your-app-id"
  api_secret: "your-api-secret"
  api_key: "your-api-key"
  base_url: "https://spark-api.xf-yun.com/v3.1/chat"
  model: "spark-v3.1"
  max_tokens: 2000
  timeout: 60
```

### 环境变量配置（推荐）

创建 `.env` 文件：

```bash
# 讯飞星火 API
XUNFEI_APP_ID=your-app-id
XUNFEI_API_SECRET=your-api-secret
XUNFEI_API_KEY=your-api-key
XUNFEI_BASE_URL=https://spark-api.xf-yun.com/v3.1/chat
XUNFEI_MODEL=spark-v3.1

# DeepSeek API (备用)
DEEPSEEK_API_KEY=your-deepseek-api-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-coder

# OpenAI API (备用)
OPENAI_API_KEY=your-openai-api-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-3.5-turbo
```

## 3. 可用模型

### 讯飞星火模型列表

| 模型名称 | 版本 | 特点 | 适用场景 |
|---------|------|------|----------|
| `spark-v1.5` | V1.5 | 基础版本 | 通用对话 |
| `spark-v2.0` | V2.0 | 增强版本 | 复杂任务 |
| `spark-v3.0` | V3.0 | 最新版本 | 高级应用 |
| `spark-v3.1` | V3.1 | 最新版本 | 推荐使用 |

### 模型选择建议

- **教学问答**: 使用 `spark-v3.1`
- **代码生成**: 使用 `spark-v3.1`
- **内容创作**: 使用 `spark-v3.1`

## 4. API 使用示例

### 基础聊天

```bash
curl -X POST "https://spark-api.xf-yun.com/v3.1/chat?authorization=xxx&date=xxx&host=spark-api.xf-yun.com" \
  -H "Content-Type: application/json" \
  -d '{
    "header": {
      "app_id": "your-app-id",
      "uid": "12345"
    },
    "parameter": {
      "chat": {
        "domain": "general",
        "temperature": 0.7,
        "max_tokens": 2000
      }
    },
    "payload": {
      "message": {
        "text": [
          {"role": "user", "content": "你好，请介绍一下你自己"}
        ]
      }
    }
  }'
```

### 编程教学场景

```bash
curl -X POST "https://spark-api.xf-yun.com/v3.1/chat?authorization=xxx&date=xxx&host=spark-api.xf-yun.com" \
  -H "Content-Type: application/json" \
  -d '{
    "header": {
      "app_id": "your-app-id",
      "uid": "12345"
    },
    "parameter": {
      "chat": {
        "domain": "general",
        "temperature": 0.7,
        "max_tokens": 2000
      }
    },
    "payload": {
      "message": {
        "text": [
          {"role": "user", "content": "请用Python写一个简单的计算器函数，并解释代码逻辑"}
        ]
      }
    }
  }'
```

## 5. 集成到教学平台

### 课程计划生成

```go
func (s *AIService) GenerateLessonPlan(subject, topic, level string) (string, error) {
	prompt := fmt.Sprintf(`作为一位%s老师，请为%s水平的学生制定一个关于"%s"的详细课程计划。
包括：
1. 教学目标
2. 课程大纲
3. 教学方法
4. 评估方式
5. 课后作业`, subject, level, topic)
	
	return s.chatCompletion(prompt)
}
```

### 习题生成

```go
func (s *AIService) GenerateExercises(subject, topic, difficulty string, count int) (string, error) {
	prompt := fmt.Sprintf(`请为%s主题"%s"生成%d道%s难度的练习题。
每道题包括：
1. 题目描述
2. 参考答案
3. 解题思路
4. 知识点说明`, subject, topic, count, difficulty)
	
	return s.chatCompletion(prompt)
}
```

### 代码评估

```go
func (s *AIService) EvaluateCode(code, language, requirements string) (string, error) {
	prompt := fmt.Sprintf(`请评估以下%s代码：
代码：
%s

要求：
%s

请从以下方面进行评估：
1. 功能正确性
2. 代码质量
3. 性能优化
4. 最佳实践
5. 改进建议`, language, code, requirements)
	
	return s.chatCompletion(prompt)
}
```

## 6. 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 10001 | 参数错误 | 检查请求参数格式 |
| 10002 | 鉴权失败 | 检查API凭据 |
| 10003 | 余额不足 | 充值或检查套餐 |
| 10004 | 请求频率限制 | 降低请求频率 |
| 10005 | 服务异常 | 稍后重试 |

### 错误处理示例

```go
func (s *AIService) handleAPIError(err error) error {
	if strings.Contains(err.Error(), "10001") {
		return fmt.Errorf("请求参数错误，请检查输入")
	}
	if strings.Contains(err.Error(), "10002") {
		return fmt.Errorf("API鉴权失败，请检查凭据")
	}
	if strings.Contains(err.Error(), "10003") {
		return fmt.Errorf("账户余额不足，请充值")
	}
	if strings.Contains(err.Error(), "10004") {
		return fmt.Errorf("请求频率过高，请稍后重试")
	}
	return err
}
```

## 7. 成本控制

### 计费方式

讯飞星火 API 按 token 计费：
- 输入 token：0.1元 / 1K tokens
- 输出 token：0.1元 / 1K tokens

### 优化建议

1. **设置合理的 max_tokens**
   ```yaml
   max_tokens: 1000  # 根据实际需求设置
   ```

2. **使用系统提示词**
   ```go
   messages := []XunfeiText{
       {Role: "system", Content: "你是一个编程老师，请简洁回答"},
       {Role: "user", Content: prompt},
   }
   ```

3. **缓存常用响应**
   ```go
   // 使用 Redis 缓存常见问题的回答
   cacheKey := fmt.Sprintf("ai:response:%s", hash(prompt))
   if cached, err := redis.Get(cacheKey); err == nil {
       return cached, nil
   }
   ```

## 8. 监控和日志

### 请求监控

```go
func (s *AIService) logRequest(prompt string, response string, duration time.Duration) {
	log.Printf("AI Request - Provider: xunfei, Duration: %v, Tokens: %d", duration, len(prompt)+len(response))
}
```

### 性能监控

```go
func (s *AIService) monitorPerformance() {
	// 记录响应时间
	start := time.Now()
	response, err := s.chatCompletion(prompt)
	duration := time.Since(start)
	
	// 记录指标
	metrics.RecordAPILatency("xunfei", duration)
	metrics.RecordAPISuccess("xunfei", err == nil)
}
```

## 9. 安全考虑

### API 凭据安全

1. **环境变量存储**
   ```bash
   export XUNFEI_APP_ID="your-app-id"
   export XUNFEI_API_SECRET="your-api-secret"
   export XUNFEI_API_KEY="your-api-key"
   ```

2. **配置文件权限**
   ```bash
   chmod 600 config/config.yaml
   ```

3. **凭据轮换**
   - 定期更换 API 凭据
   - 监控异常使用

### 输入验证

```go
func (s *AIService) validateInput(prompt string) error {
	if len(prompt) > 10000 {
		return fmt.Errorf("输入内容过长")
	}
	if strings.Contains(prompt, "恶意内容") {
		return fmt.Errorf("输入内容包含敏感信息")
	}
	return nil
}
```

## 10. 测试

### 单元测试

```go
func TestXunfeiService(t *testing.T) {
	service := NewXunfeiService()
	
	response, err := service.Generate("Hello, how are you?")
	if err != nil {
		t.Fatalf("API调用失败: %v", err)
	}
	
	if response == "" {
		t.Error("响应为空")
	}
}
```

### 集成测试

```go
func TestAIIntegration(t *testing.T) {
	service := NewAIService()
	
	// 测试课程计划生成
	plan, err := service.GenerateLessonPlan("Python", "函数", "初级")
	if err != nil {
		t.Fatalf("课程计划生成失败: %v", err)
	}
	
	// 测试习题生成
	exercises, err := service.GenerateExercises("Python", "函数", "中等", 3)
	if err != nil {
		t.Fatalf("习题生成失败: %v", err)
	}
	
	t.Logf("课程计划: %s", plan)
	t.Logf("习题: %s", exercises)
}
```

## 11. 部署检查清单

- [ ] 获取讯飞星火 API 凭据
- [ ] 更新配置文件
- [ ] 设置环境变量
- [ ] 测试 API 连接
- [ ] 配置错误处理
- [ ] 设置监控日志
- [ ] 配置成本控制
- [ ] 进行安全测试
- [ ] 运行集成测试

## 12. 故障排除

### API 连接问题

```bash
# 测试连接
python scripts/test_xunfei_api.py "your-app-id" "your-api-secret" "your-api-key"
```

### 配置问题

```bash
# 检查配置
go run main.go --config-check
```

### 性能问题

```bash
# 监控请求
curl -w "@curl-format.txt" -o /dev/null -s \
     -X POST "https://spark-api.xf-yun.com/v3.1/chat?authorization=xxx&date=xxx&host=spark-api.xf-yun.com" \
     -H "Content-Type: application/json" \
     -d '{"header":{"app_id":"test"},"parameter":{"chat":{"domain":"general"}},"payload":{"message":{"text":[{"role":"user","content":"test"}]}}}'
```

## 13. 快速配置

### 使用提供的凭据

如果您有现成的凭据，可以直接运行配置脚本：

```bash
# Linux/macOS
chmod +x scripts/setup_xunfei_api.sh
./scripts/setup_xunfei_api.sh

# Windows
scripts\setup_xunfei_api.bat
```

脚本会自动：
1. 使用提供的凭据配置项目
2. 更新配置文件
3. 创建环境变量文件
4. 测试 API 连接
5. 提供下一步指导

## 总结

通过以上配置，您已经成功将讯飞星火 API 集成到教学培训平台中。讯飞星火优秀的中文理解和生成能力将为您的教学平台提供优质的 AI 辅助功能。

记住：
- 妥善保管 API 凭据
- 监控使用成本
- 定期测试服务
- 关注 API 更新 