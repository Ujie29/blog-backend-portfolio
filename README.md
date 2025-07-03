# 🛠 blog-backend-portfolio

This project is the backend service for my personal blog system, built with Golang, implementing API architecture for both frontend and admin interfaces.

---

## 📌 Project Features

- Uses **Golang full development**
- Clear separation of frontend and backend architecture:
  - `admin/`: Backend API and management logic
  - `member/`: Frontend API (for user operations)
- Supports modular expansion, can add new sites or services at any time
- Adopts RESTful API design
- Uses middleware to abstract shared logic
- Deployment method: **Push Dockerfile to GCP, automatically built by GCP Cloud Build to create image and deploy to production environment**

---

## 🧱 Tech Stack

- Go 1.20+
- Uses bin as main framework
- ORM adopts bun
- Docker / Docker Compose
- Google Cloud Build / GCP Deploy
- RESTful API architecture design
- Modular design, well-structured file layering

---

## 📁 Project Architecture Overview

```bash
blog-backend/
├── api/                  # Frontend and backend API projects (admin/member)
├── common/               # Shared logic (middleware, entity, utils)
├── infra/                # Deployment related files (Cloud Build YAML, etc.)
├── docker-compose.yml    # For local building
├── Makefile              # Build / test shortcut commands
├── go.mod                # Go module configuration
└── README.md
```

---

## 🚀 Deployment Commands (GCP Cloud Build)

Each deployment uses the current date and Git commit hash to generate TAG, automatically submitting for build.

### 🧩 Frontend Services
```bash
gcloud builds submit \
  --config=infra/member/cloudbuild-apigw.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)

gcloud builds submit \
  --config=infra/member/cloudbuild-post.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)

gcloud builds submit \
  --config=infra/member/cloudbuild-category.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)
```

### 🛠 Backend Services
```bash
gcloud builds submit \
  --config=infra/admin/cloudbuild-apigw.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)

gcloud builds submit \
  --config=infra/admin/cloudbuild-post.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)

gcloud builds submit \
  --config=infra/admin/cloudbuild-category.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)
```

### ⏱ Batch Jobs
```bash
gcloud builds submit \
  --config=infra/batch/cloudbuild.yaml \
  --substitutions=_TAG=$(date +%Y%m%d)-$(git rev-parse --short=10 HEAD)
```

---

## 🔧 Environment Variables Configuration

```bash
# 🌐 CORS Configuration
CORS_ALLOW_ORIGINS=*
# Recommended for actual deployment:
# CORS_ALLOW_ORIGINS=http://localhost:5173,https://your-app.com

# 🧩 Frontend Services
POST_member_SERVICE=http://localhost:8081
CATEGORY_member_SERVICE=http://localhost:8082

# 🛠 Backend Services
POST_admin_SERVICE=http://localhost:8181
CATEGORY_admin_SERVICE=http://localhost:8182

# 🗄 Database
POSTGRES_URL=postgres://postgres:0000@localhost:5432/blog

# 🌱 Environment
ENV=local

# ☁️ R2 Object Storage (Cloudflare R2)
R2_ENDPOINT=https://xxx.r2.cloudflarestorage.com
R2_ACCESS_KEY=xxx
R2_SECRET_KEY=xxx
R2_PUBLIC_BASE_URL=https://pub-xxx.r2.dev
```