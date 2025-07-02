# DeepSeek API 配置指南

## 概述

本指南将帮助您配置 DeepSeek API 服务，将其集成到教学培训平台中。DeepSeek 提供了强大的代码生成和编程辅助能力。

## 1. 获取 API 密钥

### 注册账号
1. 访问 [DeepSeek Platform](https://platform.deepseek.com/)
2. 点击 "Sign Up" 注册账号
3. 验证邮箱地址

### 获取 API Key
1. 登录后进入 [API Keys](https://platform.deepseek.com/api_keys) 页面
2. 点击 "Create API Key"
3. 输入密钥名称（如：teaching-platform）
4. 复制生成的 API Key（注意保存，只显示一次）

## 2. 配置项目

### 更新配置文件

编辑 `backend/config/config.yaml`：

```yaml
# AI服务配置
ai:
  # DeepSeek API配置
  provider: "deepseek"  # deepseek, openai
  api_key: "sk-your-deepseek-api-key-here"
  base_url: "https://api.deepseek.com/v1"
  model: "deepseek-coder"  # deepseek-coder, deepseek-chat
  max_tokens: 2000
  temperature: 0.7
  timeout: 60
  
  # OpenAI配置（备用）
  openai_api_key: "sk-your-openai-api-key-here"
  openai_base_url: "https://api.openai.com/v1"
  openai_model: "gpt-3.5-turbo"
```

### 环境变量配置（推荐）

创建 `.env` 文件：

```bash
# DeepSeek API
DEEPSEEK_API_KEY=sk-your-deepseek-api-key-here
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-coder

# OpenAI API (备用)
OPENAI_API_KEY=sk-your-openai-api-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-3.5-turbo
```

## 3. 可用模型

### DeepSeek 模型列表

| 模型名称 | 类型 | 特点 | 适用场景 |
|---------|------|------|----------|
| `deepseek-coder` | 代码生成 | 专精代码生成和编程 | 编程教学、代码审查 |
| `deepseek-chat` | 通用对话 | 通用对话能力 | 教学问答、内容生成 |
| `deepseek-coder-33b-instruct` | 代码生成 | 更大模型，更高质量 | 复杂编程任务 |

### 模型选择建议

- **编程教学**: 使用 `deepseek-coder`
- **通用教学**: 使用 `deepseek-chat`
- **高级编程**: 使用 `deepseek-coder-33b-instruct`

## 4. API 使用示例

### 基础文本生成

```bash
curl -X POST https://api.deepseek.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "deepseek-coder",
    "messages": [
      {"role": "user", "content": "Write a Python function to calculate fibonacci numbers"}
    ],
    "max_tokens": 2000,
    "temperature": 0.7
  }'
```

### 编程教学场景

```bash
curl -X POST https://api.deepseek.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-your-api-key" \
  -d '{
    "model": "deepseek-coder",
    "messages": [
      {"role": "system", "content": "You are a programming teacher helping students learn Python."},
      {"role": "user", "content": "Explain how to use list comprehensions in Python with examples"}
    ],
    "max_tokens": 2000,
    "temperature": 0.7
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
| 401 | 认证失败 | 检查 API Key 是否正确 |
| 429 | 请求频率限制 | 降低请求频率或升级套餐 |
| 500 | 服务器错误 | 稍后重试 |
| 400 | 请求参数错误 | 检查请求格式 |

### 错误处理示例

```go
func (s *AIService) handleAPIError(err error) error {
	if strings.Contains(err.Error(), "401") {
		return fmt.Errorf("API认证失败，请检查API Key")
	}
	if strings.Contains(err.Error(), "429") {
		return fmt.Errorf("请求频率过高，请稍后重试")
	}
	if strings.Contains(err.Error(), "500") {
		return fmt.Errorf("服务器错误，请稍后重试")
	}
	return err
}
```

## 7. 成本控制

### 计费方式

DeepSeek API 按 token 计费：
- 输入 token：$0.0001 / 1K tokens
- 输出 token：$0.0002 / 1K tokens

### 优化建议

1. **设置合理的 max_tokens**
   ```yaml
   max_tokens: 1000  # 根据实际需求设置
   ```

2. **使用系统提示词**
   ```go
   messages := []DeepSeekMessage{
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
	log.Printf("AI Request - Duration: %v, Tokens: %d", duration, len(prompt)+len(response))
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
	metrics.RecordAPILatency("deepseek", duration)
	metrics.RecordAPISuccess("deepseek", err == nil)
}
```

## 9. 安全考虑

### API Key 安全

1. **环境变量存储**
   ```bash
   export DEEPSEEK_API_KEY="sk-your-key"
   ```

2. **配置文件权限**
   ```bash
   chmod 600 config/config.yaml
   ```

3. **密钥轮换**
   - 定期更换 API Key
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
func TestDeepSeekService(t *testing.T) {
	service := NewDeepSeekService()
	
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

- [ ] 获取 DeepSeek API Key
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
curl -H "Authorization: Bearer sk-your-key" \
     https://api.deepseek.com/v1/models
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
     -X POST https://api.deepseek.com/v1/chat/completions \
     -H "Authorization: Bearer sk-your-key" \
     -d '{"model":"deepseek-coder","messages":[{"role":"user","content":"test"}]}'
```

## 总结

通过以上配置，您已经成功将 DeepSeek API 集成到教学培训平台中。DeepSeek 强大的代码生成能力将为您的教学平台提供优质的 AI 辅助功能。

记住：
- 妥善保管 API Key
- 监控使用成本
- 定期测试服务
- 关注 API 更新 