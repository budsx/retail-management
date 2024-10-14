# Golang Binary Build Stage
FROM golang:1.20.7-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -installsuffix cgo -o RetailManagement
RUN ls -all && pwd

FROM alpine:latest

COPY --from=builder /app/RetailManagement /RetailManagement

ENTRYPOINT [ "/RetailManagement" ]