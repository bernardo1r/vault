FROM golang:alpine AS builder
COPY ./backend/ /home/backend
WORKDIR /home/backend
RUN go build -ldflags="-s -w" -trimpath -o backend cmd/main.go

FROM scratch
COPY --from=builder /home/backend/backend /home/backend/backend
WORKDIR /home/backend
COPY ./cert/backend/cert.pem .
COPY ./cert/backend/key.pem .
COPY ./cert/backend/client_certs.json .
COPY ./cert/frontend/cert.pem client.pem

