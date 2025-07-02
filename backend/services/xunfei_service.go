package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type XunfeiX1Config struct {
	AppID     string
	APIKey    string
	APISecret string
}

type XunfeiX1Service struct {
	config XunfeiX1Config
}

func NewXunfeiX1Service(appID, apiKey, apiSecret string) *XunfeiX1Service {
	return &XunfeiX1Service{
		config: XunfeiX1Config{
			AppID:     appID,
			APIKey:    apiKey,
			APISecret: apiSecret,
		},
	}
}

// 生成X1 WebSocket鉴权URL
func (s *XunfeiX1Service) genWsAuthUrl() (string, error) {
	host := "spark-api.xf-yun.com"
	path := "/v1/x1"
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\nGET %s HTTP/1.1", host, date, path)
	h := hmac.New(sha256.New, []byte(s.config.APISecret))
	h.Write([]byte(signatureOrigin))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationOrigin := fmt.Sprintf(
		"api_key=\"%s\", algorithm=\"hmac-sha256\", headers=\"host date request-line\", signature=\"%s\"",
		s.config.APIKey, signature,
	)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))
	wsUrl := fmt.Sprintf(
		"wss://%s%s?authorization=%s&date=%s&host=%s",
		host, path,
		url.QueryEscape(authorization),
		url.QueryEscape(date),
		host,
	)
	return wsUrl, nil
}

// Chat: 发送用户消息，返回AI回复（流式，返回完整内容）
func (s *XunfeiX1Service) Chat(userMsg string) (string, error) {
	wsUrl, err := s.genWsAuthUrl()
	if err != nil {
		return "", err
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	req := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": s.config.AppID,
			"uid":    "test_user",
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": []map[string]interface{}{
					{
						"role":    "user",
						"content": userMsg,
					},
				},
			},
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":            "x1",
				"max_tokens":        4096,
				"temperature":       0.5,
				"presence_penalty":  1,
				"frequency_penalty": 0.02,
				"top_k":             5,
			},
		},
	}
	data, _ := json.Marshal(req)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return "", err
	}

	var fullReply string
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// 解析响应体
		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}
		// 取出payload.choices.text内容
		payload, ok := resp["payload"].(map[string]interface{})
		if !ok {
			continue
		}
		choices, ok := payload["choices"].(map[string]interface{})
		if !ok {
			continue
		}
		textArr, ok := choices["text"].([]interface{})
		if !ok {
			continue
		}
		for _, t := range textArr {
			item, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			if content, ok := item["content"].(string); ok {
				fullReply += content
			}
		}
		// 判断是否最后一包
		header, ok := resp["header"].(map[string]interface{})
		if ok {
			if status, ok := header["status"].(float64); ok && int(status) == 2 {
				break
			}
		}
	}
	return fullReply, nil
}

// ChatStream: 发送用户消息，流式回调AI回复内容
func (s *XunfeiX1Service) ChatStream(userMsg string, onChunk func(string) error) error {
	wsUrl, err := s.genWsAuthUrl()
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": s.config.AppID,
			"uid":    "test_user",
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": []map[string]interface{}{
					{
						"role":    "user",
						"content": userMsg,
					},
				},
			},
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":            "x1",
				"max_tokens":        4096,
				"temperature":       0.5,
				"presence_penalty":  1,
				"frequency_penalty": 0.02,
				"top_k":             5,
			},
		},
	}
	data, _ := json.Marshal(req)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}
		payload, ok := resp["payload"].(map[string]interface{})
		if !ok {
			continue
		}
		choices, ok := payload["choices"].(map[string]interface{})
		if !ok {
			continue
		}
		textArr, ok := choices["text"].([]interface{})
		if !ok {
			continue
		}
		for _, t := range textArr {
			item, ok := t.(map[string]interface{})
			if !ok {
				continue
			}
			if content, ok := item["content"].(string); ok {
				err := onChunk(content)
				if err != nil {
					return err
				}
			}
		}
		header, ok := resp["header"].(map[string]interface{})
		if ok {
			if status, ok := header["status"].(float64); ok && int(status) == 2 {
				return nil
			}
		}
	}
	return nil
}
