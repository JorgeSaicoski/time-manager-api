FROM docker.io/library/golang:1.23.2-alpine

RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    bash \
    curl

WORKDIR /app

COPY go.mod go.sum ./

RUN go version

RUN go mod edit -go=1.23.2

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
