FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . . 

EXPOSE 3701

CMD [ "go", "run", "cmd/main.go" ]
