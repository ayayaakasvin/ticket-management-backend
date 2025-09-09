FROM golang:1.24.6 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o built-binary .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/config ./config
COPY --from=builder /app/built-binary .

RUN chmod +x ./built-binary
# Copy any additional files needed for the application
# COPY --from=builder /app/config.yaml .

# Expose port if needed
# EXPOSE 8080

# Set environment variables if needed
# ENV VAR_NAME=value
ENV JWT_SECRET_KEY=${JWT_SECRET_KEY}
ENV CONFIG_PATH=${CONFIG_PATH}

EXPOSE ${PORT}

CMD ["./built-binary"]
# Entry point if needed
# ENTRYPOINT ["./built-binary"]