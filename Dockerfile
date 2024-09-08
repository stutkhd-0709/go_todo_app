FROM golang:1.23.0-bullseye as deploy-builer

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# デプロイ用のコンテナ
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builer /app/app .

CMD ["./app"]

# 開発用のホットリロードコンテナ
FROM golang:1.23.0 as dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD ["air"]