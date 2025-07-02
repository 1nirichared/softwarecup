# DeepSeek AI 本地部署 - 快速开始

## 🚀 快速部署

### Windows 用户

1. **运行部署脚本**
   ```cmd
   scripts\deploy_deepseek.bat
   ```

2. **进入部署目录**
   ```cmd
   cd deepseek_deploy
   ```

3. **启动服务**
   ```cmd
   start.bat
   ```

4. **测试服务**
   ```cmd
   test.bat
   ```

### Linux/macOS 用户

1. **给脚本添加执行权限**
   ```bash
   chmod +x scripts/deploy_deepseek.sh
   ```

2. **运行部署脚本**
   ```bash
   ./scripts/deploy_deepseek.sh
   ```

3. **进入部署目录**
   ```bash
   cd deepseek_deploy
   ```

4. **启动服务**
   ```bash
   ./start.sh
   ```

5. **测试服务**
   ```bash
   ./test.sh
   ```

## 📋 系统要求

### 最低配置
- **GPU**: NVIDIA RTX 3090 (24GB VRAM)
- **内存**: 32GB RAM
- **存储**: 100GB 可用空间
- **Docker**: 20.10+

### 推荐配置
- **GPU**: NVIDIA RTX 4090 或 A100
- **内存**: 64GB RAM
- **存储**: 200GB SSD

## 🔧 服务地址

部署成功后，您可以通过以下地址访问：

- **Ollama API**: http://localhost:11434
- **Web UI**: http://localhost:3000

## 📖 详细文档

更多详细信息请查看：
- [完整部署指南](docs/DEEPSEEK_DEPLOYMENT.md)
- [API 文档](docs/API.md)

## 🆘 常见问题

### 1. Docker 未安装
- Windows: 下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- Linux: `sudo apt install docker.io docker-compose`

### 2. GPU 内存不足
```bash
# 使用更小的模型
ollama pull deepseek-coder:1.3b
```

### 3. 服务启动失败
```bash
# 检查端口占用
netstat -an | findstr 11434  # Windows
netstat -tulpn | grep 11434  # Linux
```

### 4. 模型下载失败
```bash
# 检查网络连接
ping ollama.ai
```

## 📞 技术支持

如果遇到问题，请：
1. 查看 [故障排除指南](docs/DEEPSEEK_DEPLOYMENT.md#故障排除)
2. 检查服务日志：`docker-compose logs -f`
3. 确保系统满足最低要求

---

**注意**: 首次部署可能需要较长时间下载模型文件，请耐心等待。 