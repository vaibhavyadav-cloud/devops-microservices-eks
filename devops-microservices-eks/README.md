# DevOps Microservices — EKS Deep Dive

Learning project: 4 polyglot microservices + Postgres, containerized, deployed to AWS EKS
with a production-style CI/CD pipeline and (deliberately) broad coverage of Kubernetes
resource types.

## Services

| Service | Stack | Port | Responsibility |
|---|---|---|---|
| `auth-service` | Node.js + Express | 3000 | Register/login, JWT issuance |
| `catalog-service` | Python + FastAPI | 8000 | Product CRUD |
| `order-service` | Go | 8080 | Create/list orders |
| `notification-service` | Java + Spring Boot | 8081 | Simulated notification sending |

Each service owns its own Postgres database (database-per-service pattern):
`authdb`, `catalogdb`, `orderdb`. Notification service currently has no DB.

## Running each service locally (without Docker, for a quick sanity check)

**auth-service**
```
cd services/auth-service
cp .env.example .env
npm install
npm start
```

**catalog-service**
```
cd services/catalog-service
cp .env.example .env
python3 -m venv venv && source venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
```

**order-service**
```
cd services/order-service
cp .env.example .env
go mod tidy
go run main.go
```

**notification-service**
```
cd services/notification-service
mvn spring-boot:run
```

You'll need a local Postgres running with `authdb`, `catalogdb`, and `orderdb`
databases created for the services to connect successfully — each service also
starts up fine without one and will just log a DB warning (health check still works).

## Next steps (not done yet — intentionally left for you)

1. Write a `Dockerfile` for each service (this is on you, per your call — good learning).
2. `docker-compose.yml` at the root to run all 4 + Postgres together.
3. Push everything to a new GitHub repo.
4. Launch an EC2 instance as a dev environment, run the whole stack there via
   docker-compose.
5. GitHub Actions CI: build + push images.
6. EKS cluster (via `eksctl`), then Kubernetes manifests — Deployments, Services,
   Ingress, ConfigMaps, Secrets, StatefulSet for Postgres, HPA, Jobs/CronJobs, RBAC,
   NetworkPolicy, PDB, probes, etc.
7. Helm charts to templatize everything for dev/prod.

## Notes

- Health checks (`/health` on every service) are already in place — these map
  directly to K8s liveness/readiness probes later.
- DB config is read from env vars everywhere, so switching from `.env` files to
  K8s ConfigMaps/Secrets later is a drop-in change, no code edits needed.
