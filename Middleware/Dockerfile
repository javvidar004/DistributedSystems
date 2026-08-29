FROM golang:1.26.6-alpine3.24

WORKDIR /app

# Install dependencies
COPY go.mod .
#COPY go.sum .
RUN go mod download

COPY . .


CMD [ "go", "run", "main.go" ]