# Golang Devops
Implement production grade CI/CD pipeline.

# To-Do List

## 🟢 Application Layer
- [x] Finish Go app
- [x] Setup database
- [x] Setup Monitoring (Prometheus)
- [ ] Setup structured logging (zerolog / slog)
- [ ] Setup distributed tracing (OpenTelemetry + Jaeger)
- [ ] Write unit tests (>80% coverage)
- [ ] Write integration tests (testcontainers-go)
- [ ] Setup graceful shutdown (signal handling + context cancellation)
- [ ] Add health check endpoints (/healthz, /readyz)
- [ ] Add input validation (go-playground/validator)
- [ ] Add rate limiting middleware
- [ ] API versioning (/api/v1/...)

## 🟢 Database Layer
- [x] Setup database (PostgreSQL)
- [ ] Setup database migrations (golang-migrate)
- [ ] Setup connection pooling (pgxpool)
- [ ] Add database seeding for local dev
- [ ] Setup read replica routing (write vs read separation)
- [ ] Database backup strategy (pg_dump + S3)

## 🐳 Containerization
- [x] Setup Docker container
- [ ] Multi-stage Dockerfile (builder → minimal final image)
- [ ] Setup Docker Compose for local dev (app + db + grafana + prometheus)
- [ ] Add .dockerignore
- [ ] Scan image for vulnerabilities (Trivy / Docker Scout)

## ☁️ Cloud Infrastructure (AWS + Terraform)
- [ ] Setup Cloud Infra using Terraform in AWS
- [ ] Remote Terraform state (S3 backend + DynamoDB lock)
- [ ] Separate Terraform workspaces (dev / staging / prod)
- [ ] Setup VPC (public + private subnets)
- [ ] Setup ECS Fargate or EKS cluster
- [ ] Setup RDS PostgreSQL (Multi-AZ for prod)
- [ ] Setup ElastiCache Redis
- [ ] Setup S3 buckets (app assets, backups, logs)
- [ ] Setup IAM roles with least-privilege policy
- [ ] Setup AWS Secrets Manager for credentials
- [ ] Setup CloudWatch log groups

## 🔀 Networking & Routing
- [ ] Setup Traefik
- [ ] Configure TLS termination (Let's Encrypt / ACM)
- [ ] Configure middlewares (rate limit, auth, cors)
- [ ] Setup domain + Route53 DNS records
- [ ] Configure Traefik dashboard with auth

## 🚀 CI/CD Pipeline
- [ ] Setup GitHub Actions workflows
- [ ] CI — lint (golangci-lint)
- [ ] CI — run tests + coverage report
- [ ] CI — build and push Docker image to ECR
- [ ] CI — Terraform plan on PR
- [ ] CD — Terraform apply on merge to main
- [ ] CD — rolling deploy to ECS/EKS
- [ ] Setup environment promotion (dev → staging → prod)
- [ ] Setup semantic versioning + release tags

## 📊 Observability Stack
- [x] Setup Prometheus metrics
- [ ] Setup Grafana dashboards
- [ ] Setup Grafana alerting rules (PagerDuty / Slack)
- [ ] Setup Loki for log aggregation
- [ ] Setup Jaeger / Tempo for distributed tracing
- [ ] Build the three pillars dashboard (metrics + logs + traces in one view)
- [ ] Setup uptime monitoring (Grafana Synthetic Monitoring / Checkly)

## 🔐 Security
- [ ] Secrets never in code or Docker image (use Secrets Manager)
- [ ] Enable AWS GuardDuty
- [ ] Setup SAST scanning in CI (gosec)
- [ ] Setup dependency vulnerability scan (govulncheck)
- [ ] C