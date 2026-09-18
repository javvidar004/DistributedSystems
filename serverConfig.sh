#!/bin/bash

echo "Updating package lists..."
apt-get update

echo "Installing Docker and dependencies..."
apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

apt-get update
apt-get install docker-ce docker-ce-cli containerd.io


echo "Starting Docker service..."
systemctl start docker
systemctl enable docker


# Solicitar IP de la máquina al usuario
read -r -p "Ingresa la IP del servidor: " SERVER_IP

# Validar que no se ingrese una IP en blanco
while [ -z "$SERVER_IP" ]; do
  echo "La IP no puede estar vacía."
  read -r -p "Ingresa la IP del servidor: " SERVER_IP
done

# Generar archivo .env
echo "API_HOST=http://${SERVER_IP}:8080" > .env
echo "Archivo .env creado correctamente:"
cat .env

echo "Starting Docker Compose..."
docker compose up -d
