#!/usr/bin/env bash
set -euo pipefail

# setup-dev.sh - Helper to prepare local development environment
# Usage:
#   ./scripts/setup-dev.sh       # show plan (dry-run)
#   ./scripts/setup-dev.sh --apply  # perform actions

DRY_RUN=true
if [[ "${1:-}" == "--apply" ]]; then
  DRY_RUN=false
fi

echo "[setup-dev] Starting (dry-run=${DRY_RUN})"

plan() {
  echo "Plan:"
  echo " - Ensure backend/.env exists (copy from backend/.env.example if missing)"
  echo " - Ensure frontend/.env.example exists"
  echo " - Start infrastructure via docker compose"
  echo " - (Optional) Initialize DB and seed using backend/scripts"
}

ensure_backend_env() {
  if [[ ! -f backend/.env ]]; then
    echo "backend/.env not found. Will copy backend/.env.example -> backend/.env"
    if [[ "$DRY_RUN" == false ]]; then
      cp backend/.env.example backend/.env
      echo "Copied backend/.env.example to backend/.env"
    fi
  else
    echo "backend/.env already exists"
  fi
}

ensure_frontend_env_example() {
  if [[ ! -f frontend/.env.example ]]; then
    echo "Creating frontend/.env.example"
    if [[ "$DRY_RUN" == false ]]; then
      cat > frontend/.env.example <<EOF
VITE_API_URL=http://localhost:8080/api/v1
EOF
      echo "Created frontend/.env.example"
    fi
  else
    echo "frontend/.env.example already exists"
  fi
}

start_compose() {
  echo "Starting docker compose services (postgres, redis, nats, backend, frontend)"
  if [[ "$DRY_RUN" == false ]]; then
    docker compose up -d --build
  fi
}

init_db() {
  echo "Initializing database (runs migrations and seeds)"
  if [[ "$DRY_RUN" == false ]]; then
    # Pass through .env values if present
    (cd backend && DB_HOST=${DB_HOST:-localhost} DB_PORT=${DB_PORT:-5432} DB_USER=${DB_USER:-postgres} DB_PASSWORD=${DB_PASSWORD:-} DB_NAME=${DB_NAME:-velora_dev} bash scripts/init-db.sh)
    (cd backend && DB_HOST=${DB_HOST:-localhost} DB_PORT=${DB_PORT:-5432} DB_USER=${DB_USER:-postgres} DB_PASSWORD=${DB_PASSWORD:-} DB_NAME=${DB_NAME:-velora_dev} bash scripts/seed-db.sh)
  fi
}

plan
ensure_backend_env
ensure_frontend_env_example

if [[ "$DRY_RUN" == false ]]; then
  start_compose
  init_db
  echo "Setup complete. Frontend: http://localhost:3000  Backend: http://localhost:8080"
else
  echo "Dry run complete. Rerun with --apply to perform actions."
fi
