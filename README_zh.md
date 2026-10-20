<p align="center">
  <h1 align="center">🔖 ForgetURL Server</h1>
  <p align="center">
    <strong>极简书签管理服务 - 让链接收藏变得简单</strong>
  </p>
  <p align="center">
    <a href="https://github.com/forgeturl/forgeturl-server/actions/workflows/deploy.yml">
      <img src="https://github.com/forgeturl/forgeturl-server/actions/workflows/deploy.yml/badge.svg" alt="Deploy Status">
    </a>
    <a href="https://github.com/forgeturl/forgeturl-server/releases">
      <img src="https://img.shields.io/github/v/release/forgeturl/forgeturl-server" alt="Release">
    </a>
    <a href="https://github.com/forgeturl/forgeturl-server/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/forgeturl/forgeturl-server" alt="License">
    </a>
    <img src="https://img.shields.io/badge/Go-1.23-00ADD8?logo=go" alt="Go Version">
    <a href="https://github.com/forgeturl/forgeturl-server/stargazers">
      <img src="https://img.shields.io/github/stars/forgeturl/forgeturl-server?style=social" alt="Stars">
    </a>
  </p>
  <p align="center">
    <a href="#功能特性">功能特性</a> •
    <a href="#快速开始">快速开始</a> •
    <a href="#api-文档">API 文档</a> •
    <a href="#部署指南">部署指南</a>
  </p>
  <p align="center">
    <a href="./README.md">🇬🇧 English</a>
  </p>
</p>

---

## 📖 项目简介

ForgetURL 是一个现代化的书签管理平台，帮助用户轻松保存、整理和分享网页链接。通过简洁优雅的界面和强大的后端服务，让链接收藏不再繁琐。

**为什么选择 ForgetURL？**

- 🎯 **极简设计** - 专注于核心功能，拒绝臃肿
- 🔐 **安全可靠** - 支持多种第三方 OAuth 登录
- 🔗 **灵活分享** - 多级权限控制，满足不同分享场景
- 📦 **无缝迁移** - 支持书签导入/导出，轻松迁移数据

## ✨ 功能特性

### 🔐 用户认证
- 支持 Google、GitHub 等第三方 OAuth 登录
- 安全的 Token 认证机制
- 用户信息管理

### 📄 页面管理
- **创建页面** - 快速创建书签收藏页
- **编辑页面** - 支持标题、描述、链接集合的实时编辑
- **删除页面** - 安全删除不需要的页面
- **排序页面** - 自定义页面顺序

### 📁 链接组织
- **链接集合** - 将相关链接分组管理
- **标签系统** - 为链接添加标签，便于筛选
- **子链接** - 支持主链接下挂载相关子链接

### 🔗 权限分享
| 链接类型 | 前缀 | 权限说明 |
|---------|------|---------|
| 只读链接 | `R` | 仅可查看，无法编辑 |
| 编辑链接 | `E` | 可查看和编辑内容 |
| 超级权限 | `A` | 完全控制权限 |

### 📥 数据导入/导出
- 支持从浏览器导入书签
- 支持导出为通用格式

## 🛠 技术栈

| 类别 | 技术 |
|-----|------|
| **语言** | Go 1.23 |
| **Web 框架** | Gin |
| **ORM** | GORM + Gen |
| **数据库** | MySQL |
| **缓存** | Redis |
| **API 规范** | Protocol Buffers / gRPC |
| **容器化** | Docker |
| **认证** | Goth (OAuth) |

## 📁 项目结构

```
forgeturl-server/
├── app/
│   ├── api/                    # API 层
│   │   ├── proto/              # Protobuf 定义文件
│   │   │   ├── space.proto     # 空间管理 API
│   │   │   ├── login.proto     # 登录认证 API
│   │   │   └── dumplinks.proto # 数据导入导出 API
│   │   ├── space/              # 生成的空间服务代码
│   │   ├── login/              # 生成的登录服务代码
│   │   └── docs/               # Swagger 文档
│   ├── cmd/                    # 命令行入口
│   ├── conf/                   # 配置文件
│   ├── dal/                    # 数据访问层
│   │   ├── model/              # 数据模型
│   │   └── query/              # GORM Gen 查询
│   ├── pkg/                    # 公共包
│   │   ├── connector-google/   # Google OAuth 连接器
│   │   ├── core/               # 核心工具
│   │   ├── lcache/             # 本地缓存
│   │   ├── maths/              # 数学工具（ID生成）
│   │   └── middleware/         # 中间件
│   ├── route/                  # 路由配置
│   ├── main.go                 # 主入口
│   └── go.mod                  # Go 模块定义
├── tests/                      # 测试文件
├── Dockerfile                  # Docker 构建文件
└── README.md
```

## 🚀 快速开始

### 环境要求

- Go 1.23+
- MySQL 5.7+
- Redis 6.0+

### 本地开发

```bash
# 1. 克隆项目
git clone https://github.com/your-username/forgeturl.git
cd forgeturl/forgeturl-server

# 2. 安装依赖
cd app
go mod download

# 3. 配置环境
cp conf/local.toml.example conf/local.toml
# 编辑 conf/local.toml 配置数据库和 Redis 连接

# 4. 启动服务
go run main.go api start
```

服务默认运行在 `http://127.0.0.1:80`

### Docker 部署

```bash
# 构建镜像
docker build -t forgeturl-server .

# 运行容器
docker run -d -p 80:80 forgeturl-server
```

## 📚 API 文档

### 认证机制

所有需要登录态的接口需要在请求头中携带 Token：

```http
X-Token: your_access_token
```

### 核心接口

#### 空间管理

| 接口 | 方法 | 说明 | 登录态 |
|-----|------|-----|-------|
| `/space/getUserInfo` | POST | 获取用户信息 | 可选 |
| `/space/getMySpace` | POST | 获取我的空间 | ✅ |
| `/space/getPage` | POST | 获取页面详情 | 可选 |
| `/space/createPage` | POST | 创建页面 | ✅ |
| `/space/updatePage` | POST | 更新页面 | ✅ |
| `/space/deletePage` | POST | 删除页面 | ✅ |
| `/space/savePageIds` | POST | 保存页面排序 | ✅ |
| `/space/addPageLink` | POST | 生成分享链接 | ✅ |
| `/space/removePageLink` | POST | 移除分享链接 | ✅ |

#### 登录认证

| 接口 | 方法 | 说明 |
|-----|------|-----|
| `/auth/{provider}` | GET | OAuth 授权跳转 |
| `/callback/{provider}` | GET | OAuth 回调处理 |
| `/login/logout` | POST | 登出 |

#### 数据导入导出

| 接口 | 方法 | 说明 |
|-----|------|-----|
| `/dumplinks/importBookmarks` | POST | 导入书签 |
| `/dumplinks/exportBookmarks` | POST | 导出书签 |

### 请求示例

<details>
<summary><b>获取用户信息</b></summary>

```bash
curl 'http://127.0.0.1:80/space/getUserInfo' \
  -d '{"uid": 1}' \
  -H 'Content-Type: application/json'
```
</details>

<details>
<summary><b>获取我的空间</b></summary>

```bash
curl 'http://127.0.0.1:80/space/getMySpace' \
  -d '{}' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```
</details>

<details>
<summary><b>创建页面</b></summary>

```bash
curl 'http://127.0.0.1:80/space/createPage' \
  -d '{
    "title": "我的书签",
    "brief": "常用链接收藏",
    "collections": [
      {
        "title": "开发工具",
        "links": [
          {
            "title": "GitHub",
            "url": "https://github.com",
            "tags": ["代码", "开源"]
          }
        ]
      }
    ]
  }' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```
</details>

<details>
<summary><b>更新页面</b></summary>

```bash
curl 'http://127.0.0.1:80/space/updatePage' \
  -d '{
    "page_id": "O3sFmpq",
    "title": "更新后的标题",
    "brief": "更新后的描述",
    "collections": [...],
    "version": 0,
    "mask": 7
  }' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```

**更新掩码 (mask) 说明：**
- `0x01` (1): 更新标题
- `0x02` (2): 更新描述
- `0x04` (4): 更新集合
- `0x07` (7): 更新所有字段
</details>

<details>
<summary><b>生成分享链接</b></summary>

```bash
curl 'http://127.0.0.1:80/space/addPageLink' \
  -d '{
    "page_id": "O3sFmpq",
    "page_type": "readonly"
  }' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```

**page_type 选项：**
- `readonly` - 只读链接
- `edit` - 可编辑链接
- `admin` - 超级权限链接
</details>

## 🌐 环境配置

项目支持三个运行环境：

| 环境 | 配置文件 | API 地址 |
|-----|---------|---------|
| 本地开发 | `local.toml` | `http://127.0.0.1:80` |
| 测试环境 | `test.toml` | `https://test-api.brightguo.com` |
| 生产环境 | `onl.toml` | `https://api.brightguo.com` |

## 🔧 开发指南

### 代码生成

```bash
# 生成 API 代码（从 proto 文件）
./genapi.sh

# 生成 GORM 模型和查询
cd dal && ./gensql.sh
```

### 运行测试

```bash
cd tests
pip install -r requirements.txt
python run_tests.py
```

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

<p align="center">
  Made with ❤️ by ForgetURL Team
</p>

### AVM 热点通知查看凭证

微信服务号桥接 `/login/connector/avm/wechat-mp/send` 支持可选 `share_token`（需同时提供 `recommendation_id`）。桥接只把凭证编码到配置的官网 `/hotspots` 链接，不接受调用方指定跳转域名。凭证由 AVM 签发并验证，只授权免登录读取对应的单条推荐。旧调用不传此字段时仍生成原登录链接。上线时先部署桥接，再部署 AVM API 和官网前端；历史通知不会自动补上凭证。
