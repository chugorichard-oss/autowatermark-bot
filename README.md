# 🐯 AutoWatermark Bot
### Hybrid Cloud Microservice Architecture — GCP + HomeLab + N8N + Go

> Real-time automated image watermarking via Telegram.  
> Built with a production-grade hybrid architecture:  
> GCP Always Free + Reverse SSH Tunnel + self-hosted N8N + custom Go microservice.

---

## 🏗️ Architecture

📱 Mobile (photo)
↓
🤖 Telegram Bot API
↓
☁️  GCP VM e2-micro — Caddy (SSL/HTTPS termination)
↓
🔐 Reverse SSH Tunnel (NAT traversal — no open ports)
↓
🏠 HomeLab — N8N (Docker)
├── Download photo binary (best resolution — photo[2])
├── Load watermark logo from disk
├── Merge binaries by position
└── POST multipart/form-data → Go microservice
↓
⚡ Go Microservice (Fiber + image/draw)
├── Lanczos3 resize (25% of photo width)
├── RGBA compositing (bottom-right, 40px margin)
└── JPEG encoding (quality 90)
↓
🤖 Telegram → returns watermarked photo to user

---

## ⚡ Why Go instead of Python?

First version used Python (Flask + Pillow).  
Rewrote in Go (Fiber + image/draw) for:

| | Python | Go |
|---|---|---|
| Processing speed | ~200ms | ~20ms |
| RAM usage | ~120MB | ~15MB |
| Concurrency | Limited (GIL) | Native goroutines |
| Deploy | Runtime required | Single binary |

---

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| Orchestration | N8N self-hosted (Docker) | Visual workflow + webhook handling |
| Image Processing | Go + Fiber + Lanczos3 | Fast stateless microservice |
| Public Gateway | GCP e2-micro + Caddy | Free SSL/HTTPS termination |
| Tunnel | Reverse SSH (NAT traversal) | Secure homelab exposure |
| Messaging | Telegram Bot API | User interface |
| Runtime | Docker + Linux HomeLab | Full control, zero cloud cost |

---

## 📁 Project Structure
autowatermark-bot/
├── README.md
├── workflow/
│   └── autowatermark-workflow.json    # N8N workflow (import-ready)
├── watermark-service/
│   ├── main.go                        # Go microservice
│   └── Dockerfile                     # Multi-stage build
└── infrastructure/
└── docker-compose.yml             # N8N + service orchestration

---

## 🚀 Run it yourself

### Prerequisites
- Docker + Docker Compose
- GCP account (Always Free e2-micro)
- Telegram Bot Token (from @BotFather)
- Go 1.22+

### 1. Clone the repo
```bash
git clone https://github.com/chugorichard-oss/autowatermark-bot
cd autowatermark-bot
```

### 2. Start services
```bash
docker-compose -f infrastructure/docker-compose.yml up -d
```

### 3. Import N8N workflow
N8N → Settings → Import Workflow
→ workflow/autowatermark-workflow.json
→ Add your Telegram Bot credentials
→ Update HTTP Request URL to your local service IP
→ Activate workflow

### 4. Add your watermark logo
```bash
# Copy your logo to N8N files directory
cp your-logo.png /home/node/.n8n-files/logotigre.png
```

### 5. GCP tunnel setup
```bash
# On GCP VM — reverse SSH tunnel to expose local N8N
ssh -R 5678:localhost:5678 user@your-homelab-ip

# Caddyfile on GCP VM
your-domain.com {
    reverse_proxy localhost:5678
}
```

---

## 💡 Key Engineering Decisions

**Custom Go microservice over commercial tools**
Heavy Java-based solutions crashed the homelab RAM.
Rewrote as a lightweight stateless Go service — 
single binary, ~15MB RAM, handles concurrent requests natively.

**Reverse SSH Tunnel for NAT traversal**
Exposes local N8N to the internet without opening router ports.
GCP Always Free VM acts as the public gateway with Caddy handling SSL.

**N8N parallel branches**
Photo download and watermark logo load run in parallel branches,
then merge by position — minimizing total workflow execution time.

---

## 📸 Demo

[![LinkedIn Demo](https://img.shields.io/badge/LinkedIn-Watch%20Demo%20Video-blue?logo=linkedin)](https://www.linkedin.com/in/hugo-ramirez-cloud)

---

## 👤 Author

**Hugo Ramírez** — Infrastructure & AI Automation Engineer  
La Paz, Bolivia (LATAM) · Available for remote roles  

[![LinkedIn](https://img.shields.io/badge/LinkedIn-hugo--ramirez--cloud-blue?logo=linkedin)](https://linkedin.com/in/hugo-ramirez-cloud)
[![GitHub](https://img.shields.io/badge/GitHub-chugorichard--oss-black?logo=github)](https://github.com/chugorichard-oss)
