# Deployment Guide

This document outlines the automated deployment process for the Agent Care bot. The deployment is handled via GitHub Actions, targeting a Google Cloud Platform (GCP) Compute Engine virtual machine (VM).

## 1. Deployment Trigger

The deployment workflow is automatically triggered on every `push` event to the `main` branch of this repository.

## 2. GitHub Actions Workflow (`.github/workflows/deploy.yml`)

The `deploy.yml` workflow orchestrates the entire deployment process. It runs on `ubuntu-latest` and performs the following steps:

### Steps:

1.  **Checkout Code**: Clones the repository to the GitHub Actions runner.
    *   `uses: actions/checkout@v3`
2.  **Authenticate with GCP**: Authenticates with Google Cloud Platform using a Service Account Key stored as a GitHub Secret.
    *   `uses: google-github-actions/auth@v1`
    *   `with: credentials_json: ${{ secrets.GCP_SA_KEY }}`
3.  **Set up Cloud SDK**: Installs and configures the Google Cloud SDK on the runner.
    *   `uses: google-github-actions/setup-gcloud@v1`
4.  **Build Docker Image**: Builds the Docker image for the bot using the `Dockerfile` in the repository. The image is tagged `agent-care-bot`.
    *   `run: docker build -t agent-care-bot .`
5.  **Save Docker Image**: Saves the built Docker image as a gzipped tar archive (`agent-care-bot.tar.gz`).
    *   `run: docker save agent-care-bot | gzip > agent-care-bot.tar.gz`
6.  **Copy Image to VM**: Securely copies the Docker image archive to the target GCP VM using `gcloud compute scp`. The VM name and zone are retrieved from GitHub Secrets.
    *   `gcloud compute scp agent-care-bot.tar.gz ${{ secrets.GCP_VM_NAME }}:~ --zone ${{ secrets.GCP_VM_ZONE }}`
7.  **Create `.env` file**: Creates an `.env` file on the GitHub Actions runner with sensitive environment variables (Telegram Bot Token, Database URL) retrieved from GitHub Secrets. This file is then given read-only permissions for the owner.
    *   `TG_BOT_TOKEN=${{ secrets.TG_BOT_TOKEN }}`
    *   `DATABASE_URL=${{ secrets.DATABASE_URL }}`
8.  **Copy `.env` file to VM**: Copies the generated `.env` file to the target GCP VM.
    *   `gcloud compute scp .env ${{ secrets.GCP_VM_NAME }}:~ --zone ${{ secrets.GCP_VM_ZONE }}`
9.  **Copy Deploy Script to VM**: Copies the local deployment script (`.github/workflows/deploy.sh`) to the target GCP VM.
    *   `gcloud compute scp .github/workflows/deploy.sh ${{ secrets.GCP_VM_NAME }}:~ --zone ${{ secrets.GCP_VM_ZONE }}`
10. **Deploy on VM**: Connects to the GCP VM via SSH and executes the copied `deploy.sh` script.
    *   `gcloud compute ssh ${{ secrets.GCP_VM_NAME }} --zone ${{ secrets.GCP_VM_ZONE }} --command "bash ~/deploy.sh"`

## 3. Deployment Script (`~/deploy.sh` on VM)

The `deploy.sh` script is executed directly on the target GCP VM to manage the Docker container.

### Steps:

1.  **Load Docker Image**: Loads the previously copied `agent-care-bot.tar.gz` into the Docker daemon on the VM.
    *   `sudo docker load < agent-care-bot.tar.gz`
2.  **Stop Existing Container**: Stops any running container named `agent-care-bot`. The `|| true` ensures the script doesn't fail if the container isn't running.
    *   `sudo docker stop agent-care-bot || true`
3.  **Remove Existing Container**: Removes any stopped container named `agent-care-bot`. The `|| true` ensures the script doesn't fail if the container doesn't exist.
    *   `sudo docker rm agent-care-bot || true`
4.  **Run New Docker Container**: Starts a new Docker container named `agent-care-bot` in detached mode (`-d`). It's configured to `restart always` and uses the `.env` file for environment variables.
    *   `sudo docker run -d --name agent-care-bot --restart always --env-file ~/.env agent-care-bot`
5.  **Clean up `.env` file**: Removes the `.env` file from the VM's home directory after use.
    *   `rm -f ~/.env`
6.  **Clean up Docker Image Archive**: Removes the `agent-care-bot.tar.gz` file from the VM's home directory.
    *   `rm -f ~/agent-care-bot.tar.gz`
7.  **Prune Docker Images**: Removes unused Docker images from the VM to save disk space.
    *   `sudo docker image prune -f`

## 4. Required GitHub Secrets

For this deployment workflow to function correctly, the following GitHub Secrets must be configured in your repository settings:

*   `GCP_SA_KEY`: The JSON key for a Google Cloud Service Account with permissions to create and manage VM instances, and transfer files (e.g., Compute Instance Admin, Service Account User).
*   `GCP_VM_NAME`: The name of your target GCP Compute Engine VM instance.
*   `GCP_VM_ZONE`: The GCP zone where your VM instance is located (e.g., `us-central1-a`).
*   `TG_BOT_TOKEN`: Your Telegram bot token.
*   `DATABASE_URL`: The connection string for your PostgreSQL database (e.g., Supabase).

## 5. Manual Setup on GCP VM (Prerequisites)

Before the automated deployment can succeed, ensure your GCP Compute Engine VM instance is set up with:

*   **Docker**: Docker must be installed and running on the VM.
*   **Google Cloud SDK**: Basic GCP SDK tools should be available, particularly `gcloud`.
*   **SSH Access**: Ensure that the GitHub Actions runner (via the service account) has SSH access to the VM.
*   **`.env` file permissions**: The user running the Docker container on the VM should have read access to `~/.env`.

This guide should provide a comprehensive overview of the deployment process.
