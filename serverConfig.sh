#!/bin/bash

echo "Updating package lists..."
apt update

echo "Installing Docker and dependencies..."
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "Starting Docker service..."
systemctl start docker
systemctl enable docker

echo "Starting Docker Compose..."
docker compose up -d
