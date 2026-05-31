CI/CD setup

This repository includes a GitHub Actions workflow at `.github/workflows/ci.yml`.

What it does:
- Runs Go unit tests and `golangci-lint` on the backend.
- Runs frontend `npm test` and builds the Vite app.
- Optionally builds and pushes Docker images for `backend` and `frontend` when `DOCKERHUB_USERNAME` and `DOCKERHUB_PASSWORD` secrets are set.

How to enable Docker image publish:
1. Add the following repository secrets in GitHub:
   - `DOCKERHUB_USERNAME`
   - `DOCKERHUB_PASSWORD`

Workflow triggers:
- `push` to `main`
- `pull_request` targeting `main`

To run locally, you can execute the tests manually:

Backend:
```bash
cd backend
go test ./...
```

Frontend:
```bash
cd frontend
npm ci
npm test
```

CD / Deployment
----------------

This repository includes a GitHub Actions CD workflow at `.github/workflows/cd.yml`.

What it does:
- On `push` to `main`, the workflow uses an SSH key to connect to your deployment host and runs a deployment script that does a `git reset --hard origin/main` and then runs `docker compose up -d --build`.

Required repository secrets:
- `DEPLOY_SSH_PRIVATE_KEY` - private SSH key that can log in as `DEPLOY_USER` on the server.
- `DEPLOY_HOST` - host or IP of the deployment server.
- `DEPLOY_USER` - SSH user on the deployment server.
- `DEPLOY_PATH` - path on the server where the repository is checked out (e.g., `/srv/velora`).

Prerequisites on server:
- The repository should already be cloned at `DEPLOY_PATH` and configured to track `origin/main`.
- Docker and Docker Compose must be installed and the `docker` user available to `DEPLOY_USER`.
- Add the public key corresponding to `DEPLOY_SSH_PRIVATE_KEY` to `~DEPLOY_USER/.ssh/authorized_keys`.

Recommended remote setup steps (on the server):

```bash
# create deploy user (if needed)
sudo useradd -m -s /bin/bash deploy
sudo mkdir -p /home/deploy/.ssh
sudo chown deploy:deploy /home/deploy/.ssh

# add your public key to authorized_keys
echo "<your-public-key>" | sudo tee -a /home/deploy/.ssh/authorized_keys
sudo chown deploy:deploy /home/deploy/.ssh/authorized_keys

# clone repository
sudo -u deploy git clone https://github.com/<owner>/<repo>.git /srv/velora
cd /srv/velora
sudo -u deploy git checkout main

# ensure docker and docker-compose installed
sudo apt update && sudo apt install -y docker.io docker-compose
sudo usermod -aG docker deploy
```

After setting secrets in GitHub and preparing the server, pushes to `main` will trigger an automatic deployment.

