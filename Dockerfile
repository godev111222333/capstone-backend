# Build stage
FROM golang:1.21.10-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main src/cmd/api/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY config.yaml .
COPY wait-for.sh .
COPY start.sh .

RUN ["chmod", "+x", "/app/wait-for.sh"]
RUN ["chmod", "+x", "/app/start.sh"]

EXPOSE 9876
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
