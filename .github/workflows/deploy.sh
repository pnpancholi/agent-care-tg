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
