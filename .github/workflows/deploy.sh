#!/bin/bash
set -e

sudo docker load < agent-care-bot.tar.gz
sudo docker stop agent-care-bot || true
sudo docker rm agent-care-bot || true
sudo docker run -d \
  --name agent-care-bot \
  --restart always \
  --env-file ~/.env \
  agent-care-bot

rm -f ~/.env
rm -f ~/agent-care-bot.tar.gz
sudo docker image prune -f
