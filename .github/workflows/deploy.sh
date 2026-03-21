#!/bin/bash
set -e

sudo docker load < agent-care-bot.tar.gz
sudo docker stop agent-care-bot || true
sudo docker rm agent-care-bot || true
sudo docker run -d \
  --name agent-care-bot \
  --restart always \
  --env TG_BOT_TOKEN="${TG_BOT_TOKEN}" \
  --env DATABASE_URL="${DATABASE_URL}" \
  agent-care-bot
