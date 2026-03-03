# RunFlow

Modern Runbook Automation Platform

## Features

- YAML-based runbooks
- Multi-cloud support (AWS, Azure, GCP, Kubernetes)
- Infrastructure as Code (Terraform, Ansible)
- Docker-based isolated execution
- Real-time execution logs

## Quick Start

```bash
docker-compose up -d
```

## API Endpoints

```
POST   /api/v1/runbooks          - Create runbook
GET    /api/v1/runbooks          - List runbooks
POST   /api/v1/runbooks/:id/execute  - Execute runbook
GET    /api/v1/executions/:id   - Get execution status
```

## Example Runbook

```yaml
name: Deploy App
variables:
  APP_NAME: myapp
steps:
  - name: Deploy
    command: kubectl apply -f deploy.yaml
    provider: kubectl
```

## License

MIT
