#!/bin/bash
set -e

IMAGE_NAME="sateninsync/incentive_calculator"
TAG="latest"

echo "🚀 Building Docker image..."
docker build -t $IMAGE_NAME:$TAG .

echo "🔍 Testing container locally..."
CONTAINER_NAME="myapp_test_container"

# Remove any old container with the same name
docker rm -f $CONTAINER_NAME >/dev/null 2>&1 || true

docker run -d --name $CONTAINER_NAME -p 8080:8080 $IMAGE_NAME:$TAG \
  --env prod \
  --port 8080 \
  --base-config-path ./resources

echo "⏳ Waiting for app to start..."
sleep 3

# Basic health check — you can customize this
if docker logs $CONTAINER_NAME 2>&1 | grep -qi "error"; then
  echo "❌ App reported an error in logs. Aborting push."
  docker rm -f $CONTAINER_NAME
  exit 1
fi

echo "🟢 Container started successfully. Stopping test container..."
docker rm -f $CONTAINER_NAME

echo "📤 Logging into Docker Hub..."
docker login

echo "⬆️ Pushing image to Docker Hub..."
docker push $IMAGE_NAME:$TAG

echo "✅ Done! Image pushed: $IMAGE_NAME:$TAG"
