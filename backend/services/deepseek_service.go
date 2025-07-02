package services

import (
	"backend/config"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DeepSeekService struct {
	config config.AIConfig
	client *http.Client
}

type DeepSeekChatRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
	Stream      bool              `json:"stream,omitempty"`
	Tools       []DeepSeekTool    `json:"tools,omitempty"`
	ToolChoice  interface{}       `json:"tool_choice,omitempty"`
}

type DeepSeekMessage struct {
	Role      string             `json:"role"`
	Content   string             `json:"content"`
	Name      string             `json:"name,omitempty"`
	ToolCalls []DeepSeekToolCall `json:"tool_calls,omitempty"`
}

type DeepSeekToolCall struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Function DeepSeekFunctionCall `json:"function"`
}

type DeepSeekFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type DeepSeekTool struct {
	Type     string           `json:"type"`
	Function DeepSeekFunction `json:"function"`
}

type DeepSeekFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type DeepSeekChatResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []DeepSeekChoice `json:"choices"`
	Usage   DeepSeekUsage    `json:"usage"`
}

type DeepSeekChoice struct {
	Index        int             `json:"index"`
	Message      DeepSeekMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

type DeepSeekUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewDeepSeekService() *DeepSeekService {
	cfg := config.GlobalConfig.AI
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	return &DeepSeekService{
		config: cfg,
		client: client,
	}
}

// 生成文本
func (s *DeepSeekService) Generate(prompt string) (string, error) {
	messages := []DeepSeekMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return s.Chat(messages)
}

// 聊天对话
func (s *DeepSeekService) Chat(messages []DeepSeekMessage) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", s.config.DeepSeekBaseURL)

	request := DeepSeekChatRequest{
		Model:       s.config.DeepSeekModel,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		TopP:        0.9,
		Stream:      false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.DeepSeekAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API请求失败: %s, %s", resp.Status, string(body))
	}

	var response DeepSeekChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return response.Choices[0].Message.Content, nil
}

// 流式聊天
func (s *DeepSeekService) ChatStream(messages []DeepSeekMessage, callback func(string) error) error {
	url := fmt.Sprintf("%s/chat/completions", s.config.DeepSeekBaseURL)

	request := DeepSeekChatRequest{
		Model:       s.config.DeepSeekModel,
		Messages:    messages,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		TopP:        0.9,
		Stream:      true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.DeepSeekAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API请求失败: %s, %s", resp.Status, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取流数据失败: %v", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var streamResponse DeepSeekChatResponse
			if err := json.Unmarshal([]byte(data), &streamResponse); err != nil {
				continue
			}

			if len(streamResponse.Choices) > 0 {
				content := streamResponse.Choices[0].Message.Content
				if content != "" {
					if err := callback(content); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// 健康检查
func (s *DeepSeekService) HealthCheck() bool {
	url := fmt.Sprintf("%s/models", s.config.DeepSeekBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+s.config.DeepSeekAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// 获取可用模型列表
func (s *DeepSeekService) ListModels() ([]string, error) {
	url := fmt.Sprintf("%s/models", s.config.DeepSeekBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.DeepSeekAPIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取模型列表失败: %s", resp.Status)
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析模型列表失败: %v", err)
	}

	models := make([]string, len(result.Data))
	for i, model := range result.Data {
		models[i] = model.ID
	}

	return models, nil
}
