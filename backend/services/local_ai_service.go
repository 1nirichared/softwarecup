package services

import (
	"strings"
)

// LocalAIService 本地AI服务，提供基本的AI回复功能
type LocalAIService struct{}

// NewLocalAIService 创建本地AI服务实例
func NewLocalAIService() *LocalAIService {
	return &LocalAIService{}
}

// Generate 生成文本回复
func (s *LocalAIService) Generate(prompt string) (string, error) {
	return s.generateResponse(prompt), nil
}

// Chat 聊天对话
func (s *LocalAIService) Chat(messages []interface{}) (string, error) {
	if len(messages) == 0 {
		return "你好！我是本地AI助手，很高兴为你服务。", nil
	}

	// 获取最后一条用户消息
	lastMessage := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if msg, ok := messages[i].(map[string]interface{}); ok {
			if role, exists := msg["role"].(string); exists && role == "user" {
				if content, exists := msg["content"].(string); exists {
					lastMessage = content
					break
				}
			}
		}
	}

	return s.generateResponse(lastMessage), nil
}

// generateResponse 根据用户输入生成回复
func (s *LocalAIService) generateResponse(userInput string) string {
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	// 根据关键词匹配回复
	switch {
	case strings.Contains(userInput, "你好") || strings.Contains(userInput, "hello"):
		return "你好！我是你的AI学习助手，很高兴为你服务。我可以帮助你解答学习问题、生成练习题、提供学习建议等。"

	case strings.Contains(userInput, "介绍") || strings.Contains(userInput, "自己"):
		return "我是专门为教学平台设计的AI助手，具备以下功能：\n1. 智能问答 - 解答学习问题\n2. 练习生成 - 根据课程内容生成练习题\n3. 学习建议 - 提供个性化学习建议\n4. 代码评估 - 评估编程作业\n\n请告诉我你需要什么帮助！"

	case strings.Contains(userInput, "javascript") || strings.Contains(userInput, "js"):
		return "JavaScript是一门强大的编程语言，主要用于网页开发。\n\n主要特点：\n• 动态类型语言\n• 支持函数式编程\n• 丰富的生态系统\n• 前后端都可以使用\n\n你想了解JavaScript的哪个方面？比如变量、函数、对象、异步编程等。"

	case strings.Contains(userInput, "python"):
		return "Python是一门优秀的编程语言，特别适合初学者。\n\n主要特点：\n• 语法简洁易读\n• 丰富的标准库\n• 强大的第三方库\n• 广泛应用于数据科学、AI、Web开发\n\n你想学习Python的哪个部分？"

	case strings.Contains(userInput, "算法") || strings.Contains(userInput, "数据结构"):
		return "算法和数据结构是编程的基础。\n\n常见数据结构：\n• 数组、链表、栈、队列\n• 树、图、哈希表\n\n常见算法：\n• 排序算法（冒泡、快速、归并）\n• 搜索算法（二分、深度优先、广度优先）\n• 动态规划、贪心算法\n\n你想了解哪个具体的算法或数据结构？"

	case strings.Contains(userInput, "练习") || strings.Contains(userInput, "题目"):
		return "我可以为你生成练习题！请告诉我：\n1. 你想练习什么主题？（如JavaScript、Python、算法等）\n2. 题目类型？（选择题、填空题、编程题等）\n3. 难度级别？（简单、中等、困难）\n\n我会根据你的需求生成合适的练习题。"

	case strings.Contains(userInput, "学习") || strings.Contains(userInput, "建议"):
		return "学习建议：\n\n1. 制定学习计划 - 设定明确的目标和时间表\n2. 实践为主 - 多写代码，多做练习\n3. 循序渐进 - 从基础开始，逐步深入\n4. 及时复习 - 定期回顾已学内容\n5. 项目驱动 - 通过实际项目巩固知识\n\n你目前在学习什么？我可以提供更具体的建议。"

	case strings.Contains(userInput, "谢谢") || strings.Contains(userInput, "感谢"):
		return "不客气！很高兴能帮助你。如果还有其他问题，随时可以问我。祝你学习愉快！"

	case strings.Contains(userInput, "再见") || strings.Contains(userInput, "拜拜"):
		return "再见！如果学习过程中遇到问题，随时回来找我。祝你学习进步！"

	default:
		// 根据输入内容提供更有针对性的回复
		return s.generateContextualResponse(userInput)
	}
}

// generateContextualResponse 根据上下文生成更有针对性的回复
func (s *LocalAIService) generateContextualResponse(userInput string) string {
	userInput = strings.ToLower(userInput)

	// 编程语言相关
	if strings.Contains(userInput, "golang") || strings.Contains(userInput, "go语言") {
		return "Go语言（Golang）是由Google开发的开源编程语言。\n\n主要特点：\n• 简洁的语法设计\n• 强大的并发支持（goroutines和channels）\n• 高效的垃圾回收\n• 优秀的性能表现\n• 内置工具链\n\nGo特别适合：\n• 网络服务开发\n• 微服务架构\n• 云原生应用\n• 系统编程\n\n你想了解Go的哪个方面？比如语法、并发、Web开发等。"
	}

	if strings.Contains(userInput, "变量") || strings.Contains(userInput, "var") {
		return "变量是存储数据的容器。\n\nJavaScript变量声明：\n```javascript\nvar name = 'John';        // 旧方式\nlet age = 25;             // 块级作用域\nconst PI = 3.14;          // 常量\n```\n\nGo变量声明：\n```go\nvar name string = \"John\"  // 显式类型\nage := 25                 // 类型推断\nconst PI = 3.14           // 常量\n```\n\n你想了解哪种语言的变量？"
	}

	if strings.Contains(userInput, "函数") || strings.Contains(userInput, "function") {
		return "函数是可重用的代码块。\n\nJavaScript函数：\n```javascript\nfunction greet(name) {\n    return `Hello, ${name}!`;\n}\n\n// 箭头函数\nconst add = (a, b) => a + b;\n```\n\nGo函数：\n```go\nfunc greet(name string) string {\n    return \"Hello, \" + name + \"!\"\n}\n\n// 匿名函数\nadd := func(a, b int) int {\n    return a + b\n}\n```\n\n你想了解函数的哪个方面？"
	}

	// 技术问题
	if strings.Contains(userInput, "如何") || strings.Contains(userInput, "怎么") {
		return "解决技术问题的一般步骤：\n\n1. **明确问题** - 清楚定义要解决的问题\n2. **分析原因** - 找出问题的根本原因\n3. **查找资料** - 搜索相关文档和解决方案\n4. **制定方案** - 设计解决步骤\n5. **实施验证** - 执行并测试结果\n6. **总结经验** - 记录解决方案和注意事项\n\n你能具体描述一下遇到的问题吗？这样我可以提供更精确的帮助。"
	}

	if strings.Contains(userInput, "为什么") {
		return "理解原理很重要！分析原因时可以考虑：\n\n1. **基本原理** - 涉及的核心概念和机制\n2. **影响因素** - 可能影响结果的各种因素\n3. **应用场景** - 在什么情况下会出现这种现象\n4. **最佳实践** - 如何避免或优化\n\n你能提供更多上下文信息吗？这样我可以给出更准确的解释。"
	}

	// 学习建议
	if strings.Contains(userInput, "学习") || strings.Contains(userInput, "建议") {
		return "学习编程的建议：\n\n1. **打好基础** - 掌握基本语法和概念\n2. **多写代码** - 实践是最好的老师\n3. **项目驱动** - 通过实际项目巩固知识\n4. **阅读源码** - 学习优秀代码的写法\n5. **参与社区** - 与他人交流学习\n6. **持续学习** - 技术更新很快，保持学习\n\n你目前在学习什么技术？我可以提供更具体的建议。"
	}

	// 简单问候或确认
	if len(userInput) <= 10 {
		return "我理解你的意思。如果你有具体的学习问题或需要帮助，请告诉我，我会尽力为你提供有用的信息和建议。"
	}

	// 默认回复 - 更个性化
	return "我理解你的问题。作为AI学习助手，我可以帮助你：\n\n• 解答编程相关问题\n• 解释技术概念\n• 提供学习建议\n• 生成练习题\n• 代码审查和建议\n\n请告诉我你具体想了解什么，我会为你提供有针对性的帮助。"
}

// HealthCheck 健康检查
func (s *LocalAIService) HealthCheck() bool {
	return true
}

// GetModelInfo 获取模型信息
func (s *LocalAIService) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"provider": "local",
		"model":    "local-ai-assistant",
		"version":  "1.0.0",
		"features": []string{"chat", "qa", "exercise_generation"},
	}
}
