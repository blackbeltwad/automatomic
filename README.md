# Automatomic

An Internal Developer Platform (IDP) / CI-CD engine. Push to GitHub, and
Automatomic parses the pipeline, runs each step in an isolated container,
streams build logs live to a dashboard, gates the build on security
findings, and deploys to Kubernetes — including running its own CI/CD on
its own codebase.

## Status

In active development. Target completion: end of August 2026.

- [x] Project scoped and architecture defined
- [ ] Core control plane (in progress)
- [ ] Kubernetes execution + webhooks
- [ ] Observability + security gate
- [ ] Live AWS demo + write-up

## Why this exists

A platform engineering project built to demonstrate real production
patterns: a Go backend, container/Kubernetes orchestration, a real AWS
deployment, observability, and security gating — deployed deliberately as
a short-lived, cost-conscious demo rather than a service left running
indefinitely.

## Tech stack

- **Backend:** Go — control plane API, job workers, pipeline parser
- **Frontend:** Next.js + TypeScript, TailwindCSS
- **Data:** PostgreSQL, Redis (job queue)
- **Execution:** Docker (local) → Kubernetes (local via `kind`, then AWS EKS)
- **Infra:** Terraform, AWS (EKS, RDS, ElastiCache, S3, IAM/IRSA, ALB)
- **Observability:** OpenTelemetry, Prometheus, Grafana
- **Security:** Automated SAST scanning, HMAC-verified GitHub webhooks
- **Testing:** Go tests with Testcontainers, frontend E2E tests

## Architecture

```
GitHub push/PR
   -> Webhook (HMAC-verified) -> Go API
   -> Pipeline spec parsed (YAML)
   -> Job queued in Redis
   -> Worker schedules a Kubernetes Job/Pod
   -> Logs streamed live to the dashboard
   -> Logs archived to object storage on completion
   -> Security gate can fail the build
   -> Metrics/traces exported to Prometheus/Grafana
```

## What's in this project

**Core platform (the foundation, built first and most polished):**
GitHub OAuth login, a real pipeline parser, an async job queue with a
proper state machine, isolated container execution per build step, and a
live dashboard streaming build output in real time. This is the part that
has to feel solid — everything else builds on it.

**Kubernetes execution + this repo's own CI:** builds run as real
Kubernetes Jobs rather than local containers, and this repository runs its
own automated tests and checks on every push.

**Observability and security:** metrics and traces exported to a Grafana
dashboard, and an automated check that fails a build if it detects a
leaked credential or a known-vulnerable dependency. This layer covers the
essentials rather than every possible signal — the goal is a working,
legible example of the pattern, not exhaustive coverage.

**Cloud deployment:** the full system deployed to real AWS infrastructure
(Kubernetes, managed database, managed cache, object storage) for a
recorded demo, provisioned with infrastructure-as-code rather than clicked
together by hand, then torn down — the goal is proving it runs in a real
cloud environment, not keeping it live indefinitely.

**Demo & write-up:** a recorded walkthrough of the full pipeline running
end to end on real infrastructure, along with architecture notes and
screenshots, will be linked here once complete.

## Explicitly out of scope for this version

- A native low-level execution runner (a possible future, separate project)
- Physical/hardware status display
- High-availability infrastructure (multi-AZ, redundant networking) — this
  is a demo deployment, not a production SLA

## Local development

```bash
# prereqs: Go, Node, Docker, kind, kubectl, Terraform

docker compose up -d          # local Postgres, Redis, object storage
kind create cluster --name automatomic-dev

cd cmd/api && go run .        # backend
cd web && npm install && npm run dev   # frontend

go test ./...                 # tests
```

