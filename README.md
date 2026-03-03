# RunFlow

Modern Runbook Automation Platform

## Features

- YAML-based runbooks
- Variable substitution: `{{VAR_NAME}}`
- Multi-cloud support (AWS, Azure, GCP, Kubernetes)
- Infrastructure as Code (Terraform, Ansible)
- Docker-based isolated execution
- Real-time execution logs

## Quick Start

```bash
# Clone
git clone https://github.com/dablon/runflow.git
cd runflow

# Update .env with your secrets
cp .env.example .env

# Start
docker-compose up -d
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/runbooks | Create runbook |
| GET | /api/v1/runbooks | List runbooks |
| GET | /api/v1/runbooks/:id | Get runbook |
| PUT | /api/v1/runbooks/:id | Update runbook |
| DELETE | /api/v1/runbooks/:id | Delete runbook |
| POST | /api/v1/runbooks/:id/execute | Execute runbook |
| GET | /api/v1/executions/:id | Get execution status |
| GET | /api/v1/executions/:id/logs | Get execution logs |

## Example Runbook

```yaml
name: Deploy Application
version: "1.0"
description: Deploy app to Kubernetes
environment: production

variables:
  APP_NAME: myapp
  NAMESPACE: production

steps:
  - name: Check prerequisites
    command: kubectl version --client
    timeout: 30

  - name: Deploy
    command: kubectl apply -f deploy.yaml
    provider: kubectl
    timeout: 120

on_failure:
  - name: Rollback
    command: kubectl rollout undo deployment/{{APP_NAME}}
    provider: kubectl
```

## Architecture

```
┌─────────────┐
│   REST API  │  (Go + Gin)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Executor  │  (Commands)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Runner    │  (Docker)
└─────────────┘
```

## Development

```bash
# Run tests
go test -v ./...

# Build
go build -o runflow ./cmd/api
```

## License

MIT
