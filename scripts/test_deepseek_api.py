#!/usr/bin/env python3
"""
DeepSeek API 测试脚本
用于验证 DeepSeek API 配置是否正确
"""

import os
import sys
import json
import requests
from typing import Dict, Any

class DeepSeekAPITester:
    def __init__(self, api_key: str = None, base_url: str = "https://api.deepseek.com/v1"):
        self.api_key = api_key or os.getenv("DEEPSEEK_API_KEY")
        self.base_url = base_url
        self.headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.api_key}"
        }
    
    def test_connection(self) -> bool:
        """测试 API 连接"""
        try:
            response = requests.get(f"{self.base_url}/models", headers=self.headers)
            if response.status_code == 200:
                print("✅ API 连接成功")
                return True
            else:
                print(f"❌ API 连接失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 连接错误: {e}")
            return False
    
    def test_chat_completion(self, model: str = "deepseek-coder") -> bool:
        """测试聊天完成功能"""
        try:
            data = {
                "model": model,
                "messages": [
                    {"role": "user", "content": "Hello, please respond with 'API test successful'"}
                ],
                "max_tokens": 100,
                "temperature": 0.7
            }
            
            response = requests.post(
                f"{self.base_url}/chat/completions",
                headers=self.headers,
                json=data
            )
            
            if response.status_code == 200:
                result = response.json()
                content = result["choices"][0]["message"]["content"]
                print(f"✅ 聊天完成测试成功")
                print(f"响应内容: {content}")
                return True
            else:
                print(f"❌ 聊天完成测试失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 聊天完成测试错误: {e}")
            return False
    
    def test_code_generation(self) -> bool:
        """测试代码生成功能"""
        try:
            data = {
                "model": "deepseek-coder",
                "messages": [
                    {"role": "user", "content": "Write a Python function to calculate fibonacci numbers"}
                ],
                "max_tokens": 500,
                "temperature": 0.7
            }
            
            response = requests.post(
                f"{self.base_url}/chat/completions",
                headers=self.headers,
                json=data
            )
            
            if response.status_code == 200:
                result = response.json()
                content = result["choices"][0]["message"]["content"]
                print(f"✅ 代码生成测试成功")
                print(f"生成的代码:\n{content}")
                return True
            else:
                print(f"❌ 代码生成测试失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 代码生成测试错误: {e}")
            return False
    
    def test_teaching_scenario(self) -> bool:
        """测试教学场景"""
        try:
            data = {
                "model": "deepseek-coder",
                "messages": [
                    {
                        "role": "system", 
                        "content": "你是一个编程老师，请用简洁明了的方式回答学生问题"
                    },
                    {
                        "role": "user", 
                        "content": "请解释什么是函数，并给出一个简单的例子"
                    }
                ],
                "max_tokens": 300,
                "temperature": 0.7
            }
            
            response = requests.post(
                f"{self.base_url}/chat/completions",
                headers=self.headers,
                json=data
            )
            
            if response.status_code == 200:
                result = response.json()
                content = result["choices"][0]["message"]["content"]
                print(f"✅ 教学场景测试成功")
                print(f"教学回答:\n{content}")
                return True
            else:
                print(f"❌ 教学场景测试失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 教学场景测试错误: {e}")
            return False
    
    def get_available_models(self) -> list:
        """获取可用模型列表"""
        try:
            response = requests.get(f"{self.base_url}/models", headers=self.headers)
            if response.status_code == 200:
                result = response.json()
                models = [model["id"] for model in result.get("data", [])]
                print(f"✅ 可用模型: {models}")
                return models
            else:
                print(f"❌ 获取模型列表失败: {response.status_code}")
                return []
        except Exception as e:
            print(f"❌ 获取模型列表错误: {e}")
            return []
    
    def run_all_tests(self) -> bool:
        """运行所有测试"""
        print("🚀 开始 DeepSeek API 测试...")
        print("=" * 50)
        
        if not self.api_key:
            print("❌ 未设置 API Key")
            print("请设置环境变量 DEEPSEEK_API_KEY 或在初始化时传入")
            return False
        
        tests = [
            ("API 连接测试", self.test_connection),
            ("获取模型列表", self.get_available_models),
            ("聊天完成测试", self.test_chat_completion),
            ("代码生成测试", self.test_code_generation),
            ("教学场景测试", self.test_teaching_scenario),
        ]
        
        passed = 0
        total = len(tests)
        
        for test_name, test_func in tests:
            print(f"\n📋 {test_name}")
            print("-" * 30)
            try:
                if test_func():
                    passed += 1
            except Exception as e:
                print(f"❌ {test_name} 异常: {e}")
        
        print("\n" + "=" * 50)
        print(f"📊 测试结果: {passed}/{total} 通过")
        
        if passed == total:
            print("🎉 所有测试通过！DeepSeek API 配置正确")
            return True
        else:
            print("⚠️  部分测试失败，请检查配置")
            return False

def main():
    """主函数"""
    print("DeepSeek API 测试工具")
    print("=" * 50)
    
    # 检查命令行参数
    api_key = None
    if len(sys.argv) > 1:
        api_key = sys.argv[1]
    
    # 创建测试器
    tester = DeepSeekAPITester(api_key)
    
    # 运行测试
    success = tester.run_all_tests()
    
    if success:
        print("\n✅ 配置验证完成，可以开始使用 DeepSeek API")
        sys.exit(0)
    else:
        print("\n❌ 配置验证失败，请检查 API Key 和网络连接")
        sys.exit(1)

if __name__ == "__main__":
    main() 