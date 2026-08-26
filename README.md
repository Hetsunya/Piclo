
# Piclo

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

A modern, high-performance image management service built with Go and React. It provides efficient storage, retrieval, and automatic lifecycle management of images using PostgreSQL and MinIO.

![Screenshot](link-to-screenshot)

---

## ⚡️ Quick Start

To run the entire stack (Backend, Frontend, PostgreSQL, MinIO), a single command is sufficient. No local installation of Go or Node.js is required for basic usage.

```bash
docker-compose up -d
```

Once running, the services will be available at:
- **Frontend (UI)**: http://localhost:5173
- **Backend API**: http://localhost:9090
- **MinIO Console**: http://localhost:9001 (default credentials: `minioadmin` / `minioadmin`)

---

## ✨ Features

- **Reliable Storage**: Seamless integration with MinIO for highly available object storage.
- **Automated Cleanup (TTL)**: A background worker automatically purges expired images and their database records based on a configurable Time-To-Live.
- **Strict Validation**: Server-side validation of MIME types and file sizes before processing and storage.
- **Modern UI**: Intuitive Drag & Drop image upload interface built with React 18, TypeScript, and Vite.
- **Schema Management**: Database migrations are handled reliably via `golang-migrate`.

---

## 🔌 API

Core service endpoints (base path `/api/v1`):

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/images` | Upload a new image (`multipart/form-data`). |
| `GET` | `/api/v1/images` | Retrieve a paginated list of all uploaded images with metadata. |
| `GET` | `/api/v1/images/{id}` | Retrieve metadata for a specific image by its unique ID. |
| `DELETE` | `/api/v1/images/{id}` | Delete an image from both the object storage and the database. |

> **Note:** Full interactive API documentation is available via Swagger/OpenAPI at `/swagger/index.html` (when the documentation middleware is enabled).

---

## 🛠️ Local Development (Advanced)

For local development without Docker, ensure you have **Go 1.21+** and **Node.js 18+** installed, along with local instances of PostgreSQL and MinIO.

1. Install dependencies:
   ```bash
   go mod tidy
   cd front && npm install && cd ..
   ```

2. Configure environment variables: Create a `.env` file in the root directory with your local PostgreSQL and MinIO connection details (ports, credentials, bucket name).

3. Apply database migrations using `golang-migrate`.

4. Start the development servers in separate terminals:
   ```bash
   # Terminal 1: Backend
   go run cmd/server/main.go

   # Terminal 2: Frontend
   cd front && npm run dev
   ```

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
