#!/usr/bin/env python3
"""
测试AI服务是否正常工作
"""

import requests
import json
import time

def test_ai_service():
    """测试AI服务"""
    print("🧪 测试AI服务...")
    
    # 测试配置
    base_url = "http://localhost:8080/api/v1"
    
    # 1. 先注册一个用户
    print("1. 注册用户...")
    register_data = {
        "username": "testuser",
        "password": "123456",
        "email": "test@example.com",
        "real_name": "测试用户",
        "role": "student"
    }
    
    try:
        response = requests.post(f"{base_url}/auth/register", json=register_data)
        print(f"注册响应: {response.status_code}")
        if response.status_code == 200:
            print("✅ 用户注册成功")
        else:
            print(f"❌ 用户注册失败: {response.text}")
    except Exception as e:
        print(f"❌ 注册请求失败: {e}")
        return
    
    # 2. 登录获取token
    print("\n2. 用户登录...")
    login_data = {
        "username": "testuser",
        "password": "123456"
    }
    
    try:
        response = requests.post(f"{base_url}/auth/login", json=login_data)
        if response.status_code == 200:
            token = response.json()["data"]["token"]
            print("✅ 登录成功")
        else:
            print(f"❌ 登录失败: {response.text}")
            return
    except Exception as e:
        print(f"❌ 登录请求失败: {e}")
        return
    
    # 3. 创建聊天会话
    print("\n3. 创建聊天会话...")
    headers = {"Authorization": f"Bearer {token}"}
    session_data = {
        "title": "AI测试会话",
        "type": "general"
    }
    
    try:
        response = requests.post(f"{base_url}/chat/sessions", json=session_data, headers=headers)
        if response.status_code == 200:
            session_id = response.json()["data"]["id"]
            print(f"✅ 会话创建成功，ID: {session_id}")
        else:
            print(f"❌ 会话创建失败: {response.text}")
            return
    except Exception as e:
        print(f"❌ 创建会话请求失败: {e}")
        return
    
    # 4. 发送消息测试AI
    print("\n4. 发送消息测试AI...")
    message_data = {
        "content": "你好，请简单介绍一下你自己"
    }
    
    try:
        response = requests.post(f"{base_url}/chat/sessions/{session_id}/messages", 
                               json=message_data, headers=headers)
        print(f"发送消息响应: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            print("✅ AI回复成功")
            print(f"用户消息: {result['data']['user_message']['content']}")
            print(f"AI回复: {result['data']['ai_message']['content']}")
        else:
            print(f"❌ AI回复失败: {response.text}")
    except Exception as e:
        print(f"❌ 发送消息请求失败: {e}")
    
    print("\n🎉 AI服务测试完成！")

if __name__ == "__main__":
    test_ai_service() 