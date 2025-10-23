# AI-Powered Ticketing System

Cloud-native backend for support ticketing (Go, Gin, Postgres, Kafka, Gemini, Redis, Notifications).

## Setup & Run Locally
- Clone: `git clone yourrepo`
- Deps: `go mod tidy`
- Docker: `cd docker && docker-compose up -d`
- Run: `go run cmd/[name]-service/main.go` for each.

## Docker
`cd docker && docker-compose up --build -d`

## K8s (Minikube)
`minikube start`
`kubectl apply -f k8s/`
`minikube tunnel`

## CI/CD
GitHub Actions: Tests/builds/deploys on push/PR.

## Demo
[Video link or screenshot]