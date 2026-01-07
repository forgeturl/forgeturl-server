<p align="center">
  <h1 align="center">🔖 ForgetURL Server</h1>
  <p align="center">
    <strong>Minimalist Bookmark Management Service - Make Link Collection Simple</strong>
  </p>
  <p align="center">
    <a href="#features">Features</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#api-documentation">API Docs</a> •
    <a href="#deployment">Deployment</a>
  </p>
  <p align="center">
    <a href="./README_zh.md">🇨🇳 中文文档</a>
  </p>
</p>

---

## 📖 Introduction

ForgetURL is a modern bookmark management platform that helps users easily save, organize, and share web links. With a clean, elegant interface and powerful backend service, link collection becomes effortless.

**Why ForgetURL?**

- 🎯 **Minimalist Design** - Focus on core features, no bloat
- 🔐 **Secure & Reliable** - Support multiple OAuth providers
- 🔗 **Flexible Sharing** - Multi-level permission control for different sharing scenarios
- 📦 **Seamless Migration** - Import/export bookmarks for easy data migration

## ✨ Features

### 🔐 User Authentication
- Support Google, GitHub and other OAuth providers
- Secure token-based authentication
- User profile management

### 📄 Page Management
- **Create Pages** - Quickly create bookmark collection pages
- **Edit Pages** - Real-time editing of title, description, and link collections
- **Delete Pages** - Safely remove unwanted pages
- **Sort Pages** - Customize page order

### 📁 Link Organization
- **Link Collections** - Group related links together
- **Tagging System** - Add tags to links for easy filtering
- **Sub-links** - Attach related sub-links to main links

### 🔗 Permission Sharing
| Link Type | Prefix | Permission |
|-----------|--------|------------|
| Read-only | `R` | View only, no editing |
| Editable | `E` | View and edit content |
| Admin | `A` | Full control permissions |

### 📥 Import/Export
- Import bookmarks from browsers
- Export to universal formats

## 🛠 Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.23 |
| **Web Framework** | Gin |
| **ORM** | GORM + Gen |
| **Database** | MySQL |
| **Cache** | Redis |
| **API Spec** | Protocol Buffers / gRPC |
| **Container** | Docker |
| **Auth** | Goth (OAuth) |

## 📁 Project Structure

```
forgeturl-server/
├── app/
│   ├── api/                    # API Layer
│   │   ├── proto/              # Protobuf definitions
│   │   │   ├── space.proto     # Space management API
│   │   │   ├── login.proto     # Authentication API
│   │   │   └── dumplinks.proto # Import/Export API
│   │   ├── space/              # Generated space service code
│   │   ├── login/              # Generated login service code
│   │   └── docs/               # Swagger documentation
│   ├── cmd/                    # CLI entry points
│   ├── conf/                   # Configuration files
│   ├── dal/                    # Data Access Layer
│   │   ├── model/              # Data models
│   │   └── query/              # GORM Gen queries
│   ├── pkg/                    # Shared packages
│   │   ├── connector-google/   # Google OAuth connector
│   │   ├── core/               # Core utilities
│   │   ├── lcache/             # Local cache
│   │   ├── maths/              # Math utilities (ID generation)
│   │   └── middleware/         # Middlewares
│   ├── route/                  # Router configuration
│   ├── main.go                 # Main entry point
│   └── go.mod                  # Go module definition
├── tests/                      # Test files
├── Dockerfile                  # Docker build file
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- MySQL 5.7+
- Redis 6.0+

### Local Development

```bash
# 1. Clone the repository
git clone https://github.com/your-username/forgeturl.git
cd forgeturl/forgeturl-server

# 2. Install dependencies
cd app
go mod download

# 3. Configure environment
cp conf/local.toml.example conf/local.toml
# Edit conf/local.toml to configure database and Redis connection

# 4. Start the server
go run main.go api start
```

Server runs at `http://127.0.0.1:80` by default.

### Docker Deployment

```bash
# Build image
docker build -t forgeturl-server .

# Run container
docker run -d -p 80:80 forgeturl-server
```

## 📚 API Documentation

### Authentication

All endpoints requiring authentication need the Token in request header:

```http
X-Token: your_access_token
```

### Core Endpoints

#### Space Management

| Endpoint | Method | Description | Auth Required |
|----------|--------|-------------|---------------|
| `/space/getUserInfo` | POST | Get user info | Optional |
| `/space/getMySpace` | POST | Get my space | ✅ |
| `/space/getPage` | POST | Get page details | Optional |
| `/space/createPage` | POST | Create page | ✅ |
| `/space/updatePage` | POST | Update page | ✅ |
| `/space/deletePage` | POST | Delete page | ✅ |
| `/space/savePageIds` | POST | Save page order | ✅ |
| `/space/addPageLink` | POST | Generate share link | ✅ |
| `/space/removePageLink` | POST | Remove share link | ✅ |

#### Authentication

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/auth/{provider}` | GET | OAuth authorization redirect |
| `/callback/{provider}` | GET | OAuth callback handler |
| `/login/logout` | POST | Logout |

#### Import/Export

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/dumplinks/importBookmarks` | POST | Import bookmarks |
| `/dumplinks/exportBookmarks` | POST | Export bookmarks |

### Request Examples

<details>
<summary><b>Get User Info</b></summary>

```bash
curl 'http://127.0.0.1:80/space/getUserInfo' \
  -d '{"uid": 1}' \
  -H 'Content-Type: application/json'
```
</details>

<details>
<summary><b>Get My Space</b></summary>

```bash
curl 'http://127.0.0.1:80/space/getMySpace' \
  -d '{}' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```
</details>

<details>
<summary><b>Create Page</b></summary>

```bash
curl 'http://127.0.0.1:80/space/createPage' \
  -d '{
    "title": "My Bookmarks",
    "brief": "Frequently used links",
    "collections": [
      {
        "title": "Dev Tools",
        "links": [
          {
            "title": "GitHub",
            "url": "https://github.com",
            "tags": ["code", "opensource"]
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
<summary><b>Update Page</b></summary>

```bash
curl 'http://127.0.0.1:80/space/updatePage' \
  -d '{
    "page_id": "O3sFmpq",
    "title": "Updated Title",
    "brief": "Updated description",
    "collections": [...],
    "version": 0,
    "mask": 7
  }' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```

**Update Mask Values:**
- `0x01` (1): Update title
- `0x02` (2): Update description
- `0x04` (4): Update collections
- `0x07` (7): Update all fields
</details>

<details>
<summary><b>Generate Share Link</b></summary>

```bash
curl 'http://127.0.0.1:80/space/addPageLink' \
  -d '{
    "page_id": "O3sFmpq",
    "page_type": "readonly"
  }' \
  -H 'Content-Type: application/json' \
  -H 'X-Token: your_token'
```

**page_type Options:**
- `readonly` - Read-only link
- `edit` - Editable link
- `admin` - Admin link
</details>

## 🌐 Environment Configuration

The project supports three environments:

| Environment | Config File | API Address |
|-------------|-------------|-------------|
| Local | `local.toml` | `http://127.0.0.1:80` |
| Test | `test.toml` | `https://test-api.brightguo.com` |
| Production | `onl.toml` | `https://api.brightguo.com` |

## 🔧 Development Guide

### Code Generation

```bash
# Generate API code (from proto files)
./genapi.sh

# Generate GORM models and queries
cd dal && ./gensql.sh
```

### Running Tests

```bash
cd tests
pip install -r requirements.txt
python run_tests.py
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Issues and Pull Requests are welcome!

---

<p align="center">
  Made with ❤️ by ForgetURL Team
</p>
