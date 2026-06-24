# Deployment Guide — agent-care-tg

This document describes how the Telegram bot is built and deployed via GitHub Actions to a single GCP Compute Engine VM, running two isolated environments (`prod` and `staging`) as separate Docker containers.

## Architecture overview

- **One VM**, two Docker containers — no separate VM per environment.
- Each environment gets its own:
  - Docker image tag: `agent-care-bot-prod` / `agent-care-bot-staging`
  - Container name: `agent-care-bot-prod` / `agent-care-bot-staging`
  - Remote folder on the VM: `~/agent-care-bot-prod/` / `~/agent-care-bot-staging/`
  - Telegram bot token (separate bot created via @BotFather)
  - Database (separate Supabase project/database — see below)
- No port mapping is needed on either container — the bot uses **long polling** (`tg.LongPoller`), so it reaches out to Telegram rather than receiving a webhook. No exposed port means no port-collision risk between the two containers.

## Branch → environment mapping

| Branch    | Resolves to | GitHub Environment |
|-----------|-------------|---------------------|
| `main`    | `prod`      | `prod`              |
| `staging` | `staging`   | `staging`           |


## Database — Supabase

Using **Supabase** for both environments, with **two separate Supabase projects** (one for prod, one for staging) so test data never touches production data.

Use the **direct connection string** (port `5432`), not the pooled/transaction-mode connection (port `6543`). The direct connection is the right choice here because the bot is a long-running Go process holding a persistent connection — the pooler (PgBouncer) is intended for short-lived/serverless connections and can cause issues with long-lived sessions or certain Postgres features. You'll find the direct connection string under Supabase → Project Settings → Database → Connection string → URI.

```
DATABASE_URL=postgresql://postgres:[PASSWORD]@[HOST]:5432/postgres
```

Set a different one of these per GitHub Environment (see below).

## GitHub Environments & secrets

Go to **Settings → Environments** and create two environments: `prod` and `staging`.

### Environment-scoped secrets (different value per environment)

Add these inside **each** environment's "Environment secrets" section:

| Secret           | prod value                          | staging value                          |
|-------------------|--------------------------------------|------------------------------------------|
| `TG_BOT_TOKEN`    | Real bot's token (from @BotFather)  | Separate test bot's token               |
| `DATABASE_URL`    | Prod Supabase direct connection string | Staging Supabase direct connection string |

⚠️ **Important gotcha:** if `TG_BOT_TOKEN`/`DATABASE_URL` ever existed as plain repo-level secrets, GitHub will silently fall back to those for any environment that doesn't define its own override. Make sure both secrets are explicitly set in **both** environments, then delete any old repo-level versions of these two so there's no fallback path left.

### Repo-level secrets (shared, same value for both environments)

Add these under **Settings → Secrets and variables → Actions** (not inside an Environment), since they're identical regardless of which environment is deploying — same VM, same service account either way:

| Secret          | Description                                      |
|------------------|---------------------------------------------------|
| `GCP_SA_KEY`     | Full JSON key of the GCP service account used for deploys |
| `GCP_VM_NAME`    | Name of the Compute Engine VM                     |
| `GCP_VM_ZONE`    | Zone the VM lives in                               |

### Optional: protect `prod`

On the `prod` environment page → **Deployment protection rules** → enable **Required reviewers** and add yourself. This pauses any deploy to `main` until you manually approve it in the Actions tab — a cheap safety net against an accidental bad push to prod.

## Workflow file — `.github/workflows/deploy.yml`

```yaml
name: Deployment
on:
  push:
    branches:
    - main
    - staging
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: ${{ github.ref == 'refs/heads/main' && 'prod' || 'staging' }}
    env:
      ENV_NAME: ${{ github.ref == 'refs/heads/main' && 'prod' || 'staging' }}
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      - name: Authenticate with GCP
        uses: google-github-actions/auth@v1
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}
      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v1

      - name: Build Docker image
        run: docker build -t agent-care-bot-${{ env.ENV_NAME }} .
      - name: Save Docker image
        run: docker save agent-care-bot-${{ env.ENV_NAME }} | gzip > agent-care-bot-${{ env.ENV_NAME }}.tar.gz
      - name: Copy image to VM
        run: |
          gcloud compute ssh ${{ secrets.GCP_VM_NAME }} \
            --zone ${{ secrets.GCP_VM_ZONE }} \
            --command "mkdir -p ~/agent-care-bot-${{ env.ENV_NAME }}"
          gcloud compute scp agent-care-bot-${{ env.ENV_NAME }}.tar.gz \
            ${{ secrets.GCP_VM_NAME }}:~/agent-care-bot-${{ env.ENV_NAME }}/ \
            --zone ${{ secrets.GCP_VM_ZONE }}

      - name: Create .env file
        run: |
          cat <<EOF > .env
          TG_BOT_TOKEN=${{ secrets.TG_BOT_TOKEN }}
          DATABASE_URL=${{ secrets.DATABASE_URL }}
          EOF
          chmod 600 .env
      - name: Copy .env file to VM
        run: |
          gcloud compute scp .env \
            ${{ secrets.GCP_VM_NAME }}:~/agent-care-bot-${{ env.ENV_NAME }}/ \
            --zone ${{ secrets.GCP_VM_ZONE }}
      - name: Copy deploy script to VM
        run: |
          gcloud compute scp .github/workflows/deploy.sh \
            ${{ secrets.GCP_VM_NAME }}:~/agent-care-bot-${{ env.ENV_NAME }}/ \
            --zone ${{ secrets.GCP_VM_ZONE }}
      - name: Deploy on VM
        run: |
          gcloud compute ssh ${{ secrets.GCP_VM_NAME }} \
            --zone ${{ secrets.GCP_VM_ZONE }} \
            --command "cd ~/agent-care-bot-${{ env.ENV_NAME }} && bash deploy.sh ${{ env.ENV_NAME }}"
```

## Deploy script — `.github/workflows/deploy.sh`

```bash
#!/bin/bash
set -e
ENV_NAME=$1
sudo docker load < agent-care-bot-${ENV_NAME}.tar.gz
sudo docker stop agent-care-bot-${ENV_NAME} || true
sudo docker rm agent-care-bot-${ENV_NAME} || true
sudo docker run -d \
  --name agent-care-bot-${ENV_NAME} \
  --restart always \
  --env-file .env \
  agent-care-bot-${ENV_NAME}
rm -f .env
rm -f agent-care-bot-${ENV_NAME}.tar.gz
sudo docker image prune -f
```

## One-time setup checklist

- [ ] Create a separate Telegram bot via @BotFather for staging
- [ ] Create a separate Supabase project for staging; grab its **direct** connection string
- [ ] Create `prod` and `staging` GitHub Environments
- [ ] Add `TG_BOT_TOKEN` and `DATABASE_URL` to **both** environments with the correct per-environment values
- [ ] Confirm `GCP_SA_KEY`, `GCP_VM_NAME`, `GCP_VM_ZONE` exist as repo-level secrets
- [ ] (Optional) Enable required reviewers on the `prod` environment
- [ ] Push to `staging` branch and watch the first run closely

## Known gaps / future improvements

Not addressed yet — flagged for later, not blocking initial use:

- **No health check after deploy** — if the new binary crashes immediately, `systemctl`/`docker run -d` won't surface that; the workflow goes green but the bot could be down. Consider checking `docker ps` for the container shortly after starting it and failing the job if it's not running.
- **No rollback mechanism** — the previous image/tarball is overwritten with no backup. Consider keeping the prior image tagged as `-previous` before replacing it, so reverting is a quick manual step.
- **`docker image prune -f` is environment-agnostic** — it prunes all dangling images on the VM, not just the one just deployed. Rare edge case, but worth knowing if a prod rollback image goes missing unexpectedly.
