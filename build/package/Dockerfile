FROM golang:1.25 AS builder

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download

COPY ../../ .

RUN go build -o main ./cmd/main.go


FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/main ./
#COPY ../db_data /app/db_data

ENTRYPOINT ["./main"]