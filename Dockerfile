FROM golang:1.24.6 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o built-binary .

FROM alpine:latest

WORKDIR /app

COPY --from=builder built-binary .
# Copy any additional files needed for the application
# COPY --from=builder /app/config.yaml .

# Expose port if needed
# EXPOSE 8080

# Set environment variables if needed
# ENV VAR_NAME=value

CMD ["./built-binary"]
# Entry point if needed
# ENTRYPOINT ["./built-binary"]