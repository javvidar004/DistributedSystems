FROM golang:1.26.6-alpine3.24

WORKDIR /app

# Install dependencies
RUN go mod init myapp
RUN go get github.com/gin-gonic/gin
RUN	go get github.com/lib/pq
#COPY go.mod .
#COPY go.sum .
RUN go mod download

COPY main.go .


CMD [ "go", "run", "main.go" ]