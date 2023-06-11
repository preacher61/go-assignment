# 1.Build
FROM golang:1.16-alpine AS builder

WORKDIR /app
COPY . /app

RUN apk add --no-cache make
RUN make build

# Step 2: Execute
FROM alpine:latest
WORKDIR /data

COPY . /app
COPY --from=builder /app/build/* /data/