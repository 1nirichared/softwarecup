import requests
import json
from typing import List, Dict, Any

class OllamaService:
    def __init__(self, base_url: str = "http://localhost:11434"):
        self.base_url = base_url
        self.model = "deepseek-coder:6.7b"
    
    def generate(self, prompt: str, max_tokens: int = 2000, temperature: float = 0.7) -> str:
        """生成文本"""
        url = f"{self.base_url}/api/generate"
        
        data = {
            "model": self.model,
            "prompt": prompt,
            "stream": False,
            "options": {
                "num_predict": max_tokens,
                "temperature": temperature,
                "top_p": 0.9,
                "top_k": 40
            }
        }
        
        try:
            response = requests.post(url, json=data, timeout=60)
            response.raise_for_status()
            result = response.json()
            return result.get("response", "")
        except Exception as e:
            print(f"Ollama API调用失败: {e}")
            return ""
    
    def chat(self, messages: List[Dict[str, str]], max_tokens: int = 2000, temperature: float = 0.7) -> str:
        """聊天对话"""
        url = f"{self.base_url}/api/chat"
        
        data = {
            "model": self.model,
            "messages": messages,
            "stream": False,
            "options": {
                "num_predict": max_tokens,
                "temperature": temperature,
                "top_p": 0.9,
                "top_k": 40
            }
        }
        
        try:
            response = requests.post(url, json=data, timeout=60)
            response.raise_for_status()
            result = response.json()
            return result.get("message", {}).get("content", "")
        except Exception as e:
            print(f"Ollama Chat API调用失败: {e}")
            return ""
    
    def list_models(self) -> List[str]:
        """列出可用模型"""
        url = f"{self.base_url}/api/tags"
        
        try:
            response = requests.get(url)
            response.raise_for_status()
            result = response.json()
            return [model["name"] for model in result.get("models", [])]
        except Exception as e:
            print(f"获取模型列表失败: {e}")
            return []
    
    def health_check(self) -> bool:
        """健康检查"""
        try:
            response = requests.get(f"{self.base_url}/api/tags")
            return response.status_code == 200
        except:
            return False

# 使用示例
if __name__ == "__main__":
    ollama = OllamaService()
    
    # 检查服务状态
    if ollama.health_check():
        print("Ollama服务运行正常")
        
        # 测试生成
        prompt = "请用Python写一个简单的计算器"
        response = ollama.generate(prompt)
        print(f"生成结果: {response}")
    else:
        print("Ollama服务未运行") 