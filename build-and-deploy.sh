#!/bin/bash
# Build and deploy CleanCast for UNRAID

set -e

echo "Building Docker image..."
docker build -t cleancast:custom .

echo "Stopping existing container if running..."
docker stop CleanCast 2>/dev/null || true
docker rm CleanCast 2>/dev/null || true

echo "Starting CleanCast container..."
docker run \
  -d \
  --name='CleanCast' \
  --net='bridge' \
  --pids-limit 2048 \
  -e TZ="America/New_York" \
  -e HOST_OS="Unraid" \
  -e HOST_HOSTNAME="UNRAID" \
  -e HOST_CONTAINERNAME="CleanCast" \
  -e 'GOOGLE_API_KEY'='AIzaSyDZuJay1GyPMwA8z2mnENflB8IwFjpSwVk' \
  -e 'TOKEN'='A0VgrfV8VkfnEARmChkcwAUDPEF1TOB' \
  -e 'SPONSORBLOCK_CATEGORIES'='sponsor,selfpromo,intro,outro' \
  -e 'MIN_DURATION'='3m' \
  -l net.unraid.docker.managed=dockerman \
  -l net.unraid.docker.webui='http://[IP]:[PORT:8080]/' \
  -l net.unraid.docker.icon='https://icon-library.com/images/podcast-icon-png/podcast-icon-png-0.jpg' \
  -p '8416:8080/tcp' \
  -v '/mnt/zcache/appdata/CleanCast':'/config':'rw' \
  -v '/mnt/user/Media/clean-cast-audio-downloads':'/config/audio':'rw' \
  cleancast:custom

echo "Container started! Check status with: docker logs CleanCast"
