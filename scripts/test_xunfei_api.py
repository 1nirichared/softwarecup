#!/usr/bin/env python3
"""
讯飞星火 API 测试脚本
用于验证讯飞星火 API 配置是否正确
"""

import os
import sys
import json
import time
import hmac
import hashlib
import base64
import requests
from datetime import datetime
from urllib.parse import urlencode
from typing import Dict, Any

class XunfeiAPITester:
    def __init__(self, app_id: str = None, api_secret: str = None, api_key: str = None):
        self.app_id = app_id or os.getenv("XUNFEI_APP_ID", "04107cca")
        self.api_secret = api_secret or os.getenv("XUNFEI_API_SECRET", "NmYyYjc2OTk1Yjc4ZGMwZDhkYWM5YjBj")
        self.api_key = api_key or os.getenv("XUNFEI_API_KEY", "1a6c4989345073f44140f47aa57e5ae8")
        self.base_url = "https://spark-api.xf-yun.com/v3.1/chat"
    
    def create_url(self) -> str:
        """生成鉴权URL"""
        # 生成RFC1123格式的时间戳
        now = datetime.now()
        date = now.strftime('%a, %d %b %Y %H:%M:%S GMT')
        
        # 拼接字符串
        signature_origin = "host: spark-api.xf-yun.com\n"
        signature_origin += "date: " + date + "\n"
        signature_origin += "GET /v3.1/chat HTTP/1.1"
        
        # 使用APISecret对signature_origin进行hmac-sha256加密
        signature_sha = hmac.new(
            self.api_secret.encode('utf-8'),
            signature_origin.encode('utf-8'),
            digestmod=hashlib.sha256
        ).digest()
        
        signature_sha_base64 = base64.b64encode(signature_sha).decode()
        authorization_origin = f'api_key="{self.api_key}", algorithm="hmac-sha256", headers="host date request-line", signature="{signature_sha_base64}"'
        authorization = base64.b64encode(authorization_origin.encode('utf-8')).decode()
        
        # 构建URL
        url = f"{self.base_url}?authorization={authorization}&date={date}&host=spark-api.xf-yun.com"
        return url
    
    def test_connection(self) -> bool:
        """测试 API 连接"""
        try:
            url = self.create_url()
            data = {
                "header": {
                    "app_id": self.app_id,
                    "uid": "12345"
                },
                "parameter": {
                    "chat": {
                        "domain": "general",
                        "temperature": 0.7,
                        "max_tokens": 100
                    }
                },
                "payload": {
                    "message": {
                        "text": [
                            {"role": "user", "content": "Hello"}
                        ]
                    }
                }
            }
            
            response = requests.post(url, json=data, timeout=30)
            
            if response.status_code == 200:
                result = response.json()
                if result.get("header", {}).get("code") == 0:
                    print("✅ API 连接成功")
                    return True
                else:
                    print(f"❌ API 返回错误: {result.get('header', {}).get('message')}")
                    return False
            else:
                print(f"❌ API 连接失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 连接错误: {e}")
            return False
    
    def test_chat_completion(self) -> bool:
        """测试聊天完成功能"""
        try:
            url = self.create_url()
            data = {
                "header": {
                    "app_id": self.app_id,
                    "uid": "12345"
                },
                "parameter": {
                    "chat": {
                        "domain": "general",
                        "temperature": 0.7,
                        "max_tokens": 200
                    }
                },
                "payload": {
                    "message": {
                        "text": [
                            {"role": "user", "content": "请简单介绍一下你自己"}
                        ]
                    }
                }
            }
            
            response = requests.post(url, json=data, timeout=30)
            
            if response.status_code == 200:
                result = response.json()
                if result.get("header", {}).get("code") == 0:
                    content = result.get("payload", {}).get("choices", {}).get("text", [{}])[0].get("content", "")
                    print(f"✅ 聊天完成测试成功")
                    print(f"响应内容: {content}")
                    return True
                else:
                    print(f"❌ 聊天完成测试失败: {result.get('header', {}).get('message')}")
                    return False
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
            url = self.create_url()
            data = {
                "header": {
                    "app_id": self.app_id,
                    "uid": "12345"
                },
                "parameter": {
                    "chat": {
                        "domain": "general",
                        "temperature": 0.7,
                        "max_tokens": 500
                    }
                },
                "payload": {
                    "message": {
                        "text": [
                            {"role": "user", "content": "请用Python写一个简单的计算器函数"}
                        ]
                    }
                }
            }
            
            response = requests.post(url, json=data, timeout=30)
            
            if response.status_code == 200:
                result = response.json()
                if result.get("header", {}).get("code") == 0:
                    content = result.get("payload", {}).get("choices", {}).get("text", [{}])[0].get("content", "")
                    print(f"✅ 代码生成测试成功")
                    print(f"生成的代码:\n{content}")
                    return True
                else:
                    print(f"❌ 代码生成测试失败: {result.get('header', {}).get('message')}")
                    return False
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
            url = self.create_url()
            data = {
                "header": {
                    "app_id": self.app_id,
                    "uid": "12345"
                },
                "parameter": {
                    "chat": {
                        "domain": "general",
                        "temperature": 0.7,
                        "max_tokens": 300
                    }
                },
                "payload": {
                    "message": {
                        "text": [
                            {"role": "user", "content": "请解释什么是函数，并给出一个简单的例子"}
                        ]
                    }
                }
            }
            
            response = requests.post(url, json=data, timeout=30)
            
            if response.status_code == 200:
                result = response.json()
                if result.get("header", {}).get("code") == 0:
                    content = result.get("payload", {}).get("choices", {}).get("text", [{}])[0].get("content", "")
                    print(f"✅ 教学场景测试成功")
                    print(f"教学回答:\n{content}")
                    return True
                else:
                    print(f"❌ 教学场景测试失败: {result.get('header', {}).get('message')}")
                    return False
            else:
                print(f"❌ 教学场景测试失败: {response.status_code}")
                print(f"错误信息: {response.text}")
                return False
        except Exception as e:
            print(f"❌ 教学场景测试错误: {e}")
            return False
    
    def test_stream_chat(self) -> bool:
        """测试流式聊天"""
        try:
            url = self.create_url()
            data = {
                "header": {
                    "app_id": self.app_id,
                    "uid": "12345"
                },
                "parameter": {
                    "chat": {
                        "domain": "general",
                        "temperature": 0.7,
                        "max_tokens": 200
                    }
                },
                "payload": {
                    "message": {
                        "text": [
                            {"role": "user", "content": "请用一句话介绍Python"}
                        ]
                    }
                }
            }
            
            response = requests.post(url, json=data, timeout=30, stream=True)
            
            if response.status_code == 200:
                print(f"✅ 流式聊天测试成功")
                print("响应内容:")
                for line in response.iter_lines():
                    if line:
                        line_str = line.decode('utf-8')
                        if line_str.startswith('data: '):
                            data_str = line_str[6:]
                            if data_str == '[DONE]':
                                break
                            try:
                                data_json = json.loads(data_str)
                                content = data_json.get("payload", {}).get("choices", {}).get("text", [{}])[0].get("content", "")
                                if content:
                                    print(content, end='', flush=True)
                            except json.JSONDecodeError:
                                continue
                print()  # 换行
                return True
            else:
                print(f"❌ 流式聊天测试失败: {response.status_code}")
                return False
        except Exception as e:
            print(f"❌ 流式聊天测试错误: {e}")
            return False
    
    def run_all_tests(self) -> bool:
        """运行所有测试"""
        print("🚀 开始讯飞星火 API 测试...")
        print("=" * 50)
        
        if not all([self.app_id, self.api_secret, self.api_key]):
            print("❌ 缺少必要的配置参数")
            print("请设置环境变量或传入参数：")
            print("  - XUNFEI_APP_ID")
            print("  - XUNFEI_API_SECRET")
            print("  - XUNFEI_API_KEY")
            return False
        
        tests = [
            ("API 连接测试", self.test_connection),
            ("聊天完成测试", self.test_chat_completion),
            ("代码生成测试", self.test_code_generation),
            ("教学场景测试", self.test_teaching_scenario),
            ("流式聊天测试", self.test_stream_chat),
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
            print("🎉 所有测试通过！讯飞星火 API 配置正确")
            return True
        else:
            print("⚠️  部分测试失败，请检查配置")
            return False

def main():
    """主函数"""
    print("讯飞星火 API 测试工具")
    print("=" * 50)
    
    # 检查命令行参数
    app_id = None
    api_secret = None
    api_key = None
    
    if len(sys.argv) > 3:
        app_id = sys.argv[1]
        api_secret = sys.argv[2]
        api_key = sys.argv[3]
    
    # 创建测试器
    tester = XunfeiAPITester(app_id, api_secret, api_key)
    
    # 运行测试
    success = tester.run_all_tests()
    
    if success:
        print("\n✅ 配置验证完成，可以开始使用讯飞星火 API")
        sys.exit(0)
    else:
        print("\n❌ 配置验证失败，请检查 API 配置和网络连接")
        sys.exit(1)

if __name__ == "__main__":
    main() 