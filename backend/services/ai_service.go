package services

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type AIService struct {
	client   *openai.Client
	config   config.AIConfig
	deepseek *DeepSeekService
	Xunfei   *XunfeiX1Service
	local    *LocalAIService
	provider string
}

func NewAIService() *AIService {
	cfg := config.GlobalConfig.AI

	// 创建 OpenAI 客户端（备用）
	openaiClient := openai.NewClient(cfg.OpenAIAPIKey)

	// 创建 DeepSeek 服务
	deepseekService := NewDeepSeekService()

	// 创建讯飞星火X1服务
	var xunfei *XunfeiX1Service
	if cfg.Provider == "xunfei" {
		xcfg := config.GlobalConfig.Xunfei
		xunfei = NewXunfeiX1Service(xcfg.AppID, xcfg.APIKey, xcfg.APISecret)
	}

	// 创建本地AI服务
	localService := NewLocalAIService()

	return &AIService{
		client:   openaiClient,
		config:   cfg,
		deepseek: deepseekService,
		Xunfei:   xunfei,
		local:    localService,
		provider: cfg.Provider,
	}
}

// 生成备课内容
func (s *AIService) GenerateLessonPlan(course *models.Course, chapter *models.Chapter, knowledge []models.Knowledge) (string, error) {
	prompt := fmt.Sprintf(`
请为以下课程生成详细的备课内容：

课程名称：%s
章节标题：%s
章节描述：%s

知识点：
%s

请生成包含以下内容的备课方案：
1. 教学目标
2. 教学重点难点
3. 教学内容安排
4. 教学方法
5. 时间分配
6. 课后作业建议

请用中文回答，内容要详细且实用。
`, course.Name, chapter.Title, chapter.Description, s.formatKnowledge(knowledge))

	return s.chatCompletion(prompt)
}

// 生成练习题
func (s *AIService) GenerateExercises(course *models.Course, chapter *models.Chapter, questionType models.QuestionType, count int) ([]models.Question, error) {
	prompt := fmt.Sprintf(`
请为以下课程内容生成%d道%s题：

课程名称：%s
章节标题：%s
章节内容：%s

要求：
1. 题目要符合教学内容
2. 难度适中
3. 包含答案和解析
4. 如果是编程题，请提供完整的代码示例

请以JSON格式返回，格式如下：
[
  {
    "title": "题目内容",
    "content": "题目详细描述",
    "options": "选项内容（如果是选择题）",
    "answer": "正确答案",
    "analysis": "解析",
    "score": 分值,
    "difficulty": 难度等级(1-5)
  }
]
`, count, s.getQuestionTypeName(questionType), course.Name, chapter.Title, chapter.Content)

	response, err := s.chatCompletion(prompt)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var questions []models.Question
	err = json.Unmarshal([]byte(response), &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

// 智能问答
func (s *AIService) ChatWithAI(session *models.ChatSession, message string, knowledgeBase []models.KnowledgeBase) (string, error) {
	// 构建上下文
	context := s.buildContext(session, knowledgeBase)

	prompt := fmt.Sprintf(`
你是一个专业的教学助手，请基于以下知识库内容回答学生的问题：

知识库内容：
%s

学生问题：%s

请提供准确、详细的回答，如果涉及编程，请提供代码示例。
`, context, message)

	return s.chatCompletion(prompt)
}

// 评估学生答案
func (s *AIService) EvaluateAnswer(question *models.Question, studentAnswer string) (int, bool, string, error) {
	prompt := fmt.Sprintf(`
请评估以下学生答案：

题目：%s
题目类型：%s
正确答案：%s
学生答案：%s

请评估：
1. 答案是否正确（true/false）
2. 得分（0-满分）
3. 详细反馈和建议

请以JSON格式返回：
{
  "is_correct": true/false,
  "score": 得分,
  "feedback": "详细反馈"
}
`, question.Title, question.Type, question.Answer, studentAnswer)

	response, err := s.chatCompletion(prompt)
	if err != nil {
		return 0, false, "", err
	}

	// 解析评估结果
	var result struct {
		IsCorrect bool   `json:"is_correct"`
		Score     int    `json:"score"`
		Feedback  string `json:"feedback"`
	}

	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		return 0, false, "", err
	}

	return result.Score, result.IsCorrect, result.Feedback, nil
}

// 生成学习建议
func (s *AIService) GenerateLearningAdvice(userID uint, progress []models.LearningProgress, answers []models.StudentAnswer) (string, error) {
	prompt := fmt.Sprintf(`
基于学生的学习数据，生成个性化学习建议：

学习进度：%s
答题情况：%s

请分析学生的学习情况，提供：
1. 薄弱知识点
2. 学习建议
3. 推荐练习
4. 学习方法指导
`, s.formatProgress(progress), s.formatAnswers(answers))

	return s.chatCompletion(prompt)
}

// 私有方法
func (s *AIService) chatCompletion(prompt string) (string, error) {
	// 根据配置选择AI提供商
	switch s.provider {
	case "xunfei":
		return s.Xunfei.Chat(prompt)
	case "deepseek":
		return s.deepseek.Generate(prompt)
	case "local":
		return s.local.Chat([]interface{}{
			map[string]interface{}{"role": "user", "content": prompt},
		})
	case "openai":
		fallthrough
	default:
		// 使用OpenAI服务
		resp, err := s.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       s.config.OpenAIModel,
				Messages:    []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleUser, Content: prompt}},
				MaxTokens:   s.config.MaxTokens,
				Temperature: float32(s.config.Temperature),
			},
		)

		if err != nil {
			return "", err
		}

		return resp.Choices[0].Message.Content, nil
	}
}

func (s *AIService) formatKnowledge(knowledge []models.Knowledge) string {
	var result strings.Builder
	for i, k := range knowledge {
		result.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, k.Title, k.Content))
	}
	return result.String()
}

func (s *AIService) getQuestionTypeName(qType models.QuestionType) string {
	switch qType {
	case models.QuestionTypeSingle:
		return "单选题"
	case models.QuestionTypeMultiple:
		return "多选题"
	case models.QuestionTypeFill:
		return "填空题"
	case models.QuestionTypeEssay:
		return "简答题"
	case models.QuestionTypeCode:
		return "编程题"
	default:
		return "题目"
	}
}

func (s *AIService) buildContext(session *models.ChatSession, knowledgeBase []models.KnowledgeBase) string {
	var context strings.Builder

	if session.CourseID != nil && session.Course != nil {
		context.WriteString(fmt.Sprintf("当前课程：%s\n", session.Course.Name))
	}

	if session.ChapterID != nil && session.Chapter != nil {
		context.WriteString(fmt.Sprintf("当前章节：%s\n", session.Chapter.Title))
	}

	context.WriteString("相关知识：\n")
	for _, kb := range knowledgeBase {
		context.WriteString(fmt.Sprintf("- %s: %s\n", kb.Name, kb.Content))
	}

	return context.String()
}

func (s *AIService) formatProgress(progress []models.LearningProgress) string {
	var result strings.Builder
	for _, p := range progress {
		result.WriteString(fmt.Sprintf("课程进度：%.1f%%, 学习时长：%d分钟\n",
			p.Progress, p.TimeSpent))
	}
	return result.String()
}

func (s *AIService) formatAnswers(answers []models.StudentAnswer) string {
	var result strings.Builder
	for _, a := range answers {
		result.WriteString(fmt.Sprintf("得分：%d, 正确：%v\n",
			a.Score, a.IsCorrect))
	}
	return result.String()
}
