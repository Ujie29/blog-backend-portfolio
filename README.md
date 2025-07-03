# 🛠 blog-backend-portfolio

This project is the backend system of my personal blog platform (showcase version), built entirely with Golang. It includes both frontend-facing (`member`) and admin-facing (`admin`) APIs. The project structure is modular, allowing future expansion with additional services or sites.

This repository is a **sanitized version for portfolio presentation**, with all environment variables and secrets removed. However, the **core logic and architecture remain fully intact**, making it a solid example of my backend development skills.

---

## 📌 Features

- Fully developed in **Golang**
- Clearly separated frontend/backend API structure:
  - `admin/`: backend API for admin management
  - `member/`: frontend API for user interactions
- Modular and scalable architecture (easy to add new sites/services)
- Follows RESTful API design principles
- Uses middleware to abstract shared logic
- Deployable via GCP Cloud Build (`cloudbuild.yaml`)

---

## 🧱 Tech Stack

- Go 1.20+
- Docker / Docker Compose
- Google Cloud Build / GCP Deploy
- RESTful API Architecture
- Modular and well-organized folder structure

---

## 📁 Project Structure

```bash
blog-backend/
├── api/                  # Frontend/Admin APIs (admin/member)
├── common/               # Shared logic (middleware, entity, utils)
├── infra/                # Deployment configs (Cloud Build YAML)
├── docker-compose.yml    # For local development
├── Makefile              # Compile / Test shortcuts
├── go.mod                # Go module config
└── README.md
```

---

## 🚀 Deployment Commands (GCP Cloud Build)

Each deployment generates a unique TAG using the current date and Git commit hash, and then triggers a build.

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

### 🛠 Admin Services
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

## 🌐 CORS Configuration

```env
CORS_ALLOW_ORIGINS=*  
# For production use:
# CORS_ALLOW_ORIGINS=http://localhost:5173,https://your-app.com
```

---

## 🧩 Frontend Service URLs

```env
POST_member_SERVICE=http://localhost:8081  
CATEGORY_member_SERVICE=http://localhost:8082
```

---

## 🛠 Admin Service URLs

```env
POST_admin_SERVICE=http://localhost:8181  
CATEGORY_admin_SERVICE=http://localhost:8182
```

---

## 🗄 Database Configuration

```env
POSTGRES_URL=postgres://postgres:0000@localhost:5432/blog
```

---

## 🌱 Environment

```env
ENV=local
```

---

## ☁️ R2 Object Storage (Cloudflare R2)

```env
R2_ENDPOINT=https://xxx.r2.cloudflarestorage.com  
R2_ACCESS_KEY=xxx  
R2_SECRET_KEY=xxx  
R2_PUBLIC_BASE_URL=https://pub-xxx.r2.dev
```

---

> ⚠️ Note: Comments in the source code are written in Traditional Chinese, as this project was originally built for personal use. The code structure and logic are universal and easy to follow.