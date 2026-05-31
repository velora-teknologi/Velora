Local Development Environment
===========================

Use this guide to set up the development environment for Velora locally.

Quick start (dry-run):

```bash
./scripts/setup-dev.sh
```

To actually perform the steps (creates `.env` and starts services):

```bash
./scripts/setup-dev.sh --apply
```

What the script does (when `--apply`):
- Ensures `backend/.env` exists (copies from `backend/.env.example` if missing).
- Ensures `frontend/.env.example` exists.
- Runs `docker compose up -d --build` to start services (Postgres, Redis, NATS, Backend, Frontend).
- Runs `backend/scripts/init-db.sh` and `backend/scripts/seed-db.sh` to initialize and seed the database.

Notes:
- You should have Docker and Docker Compose installed locally.
- If you already run services separately, skip `--apply` or stop compose before applying.
