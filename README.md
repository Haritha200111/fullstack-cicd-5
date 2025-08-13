# fullstack-prisma-ci-migrate

Minimal full-stack project (Angular frontend + Go backend + Prisma migrations) preconfigured for:

- **Local development** using docker-compose (Postgres + backend + frontend)
- **CI-based migrations**: Prisma migrations are run in **GitHub Actions (CI)** before building and deploying images
- **CI/CD**: GitHub Actions builds images and pushes to **Docker Hub**, then **SSHs into EC2** to deploy (pull & run the images)
- **Database**: Use AWS RDS PostgreSQL in production (DATABASE_URL passed as GitHub Secret)

## Quick local start
1. Ensure Docker is installed.
2. From repo root:
```bash
docker compose up --build
```
3. Frontend: http://localhost
   Backend health: http://localhost:8080/health

## CI/CD (what happens on push to main)
1. GitHub Actions runs Prisma migrations against `DATABASE_URL` secret.
2. If migrations succeed, Actions builds backend/frontend Docker images, pushes to Docker Hub.
3. Actions SSHes to EC2, pulls images, and runs containers with `DATABASE_URL` env pointing to RDS.

## Required GitHub Secrets
- `DOCKERHUB_USERNAME` — Docker Hub username
- `DOCKERHUB_TOKEN` — Docker Hub access token / password
- `EC2_HOST` — EC2 public IP or DNS
- `EC2_USER` — SSH user (e.g., ubuntu)
- `EC2_SSH_KEY` — Private key contents (.pem) for SSH (paste full file)
- `DATABASE_URL` — RDS connection string (used during migration and passed to container)

## Notes
- Migrations run only in CI (before images are built) — containers **do not** run migrations at startup.
- For local migrations, run `docker compose run backend npx prisma migrate deploy --schema=prisma/schema.prisma` or use `make migrate`.
