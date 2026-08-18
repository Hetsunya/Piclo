# Piclo - Image Management Service

Piclo is a modern image management service built with Go and React. It provides efficient storage, retrieval, and management of images using PostgreSQL and MinIO.

## Features

- Image upload and storage
- Database integration with PostgreSQL
- Cloud storage with MinIO
- Background cleanup with TTL worker
- RESTful API endpoints
- Responsive web interface

## Technologies Used

### Backend
- Go 1.21+
- PostgreSQL
- MinIO
- Gorilla Mux for routing
- pgx for database connection
- GORM for ORM (if used)

### Frontend
- React 18
- TypeScript
- Vite
- Tailwind CSS

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker and Docker Compose (for development)
- PostgreSQL
- MinIO

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Piclo
```

2. Install backend dependencies:
```bash
go mod tidy
```

3. Install frontend dependencies:
```bash
npm install
```

4. Set up environment variables in `.env` file:
```bash
PORT=9090
DATABASE_URL=postgres://piclo:piclo@localhost:5433/piclo?sslmode=disable
PUBLIC_URL=http://localhost:9090
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=piclo-images
```

### Running the Application

#### Development Mode
```bash
# Start backend server
go run cmd/server/main.go

# Start frontend development server
npm run dev
```

#### Docker Compose
```bash
docker-compose up
```

## API Endpoints

- `POST /images` - Upload image
- `GET /images/:id` - Get image by ID
- `DELETE /images/:id` - Delete image
- `GET /images` - List all images

## Project Structure

```
.
├── cmd/                 # Main application entry points
│   └── server/          # Server main package
├── internal/            # Internal packages
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP handlers and routes
│   ├── model/           # Data models
│   ├── repository/      # Database operations
│   ├── service/         # Business logic
│   ├── storage/         # Storage operations (MinIO)
│   └── worker/          # Background workers
├── front/               # Frontend application
│   ├── src/             # Source code
│   ├── public/          # Static assets
│   └── ...
├── migrations/          # Database migration files
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker build file
└── README.md            # This file
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License.