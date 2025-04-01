# --- Build Stage ---
    FROM golang:1.24-alpine AS builder

    # Set working directory inside the container
    WORKDIR /app
    
    # Copy go.mod and go.sum to cache dependencies
    COPY go.*  ./
    
    # Download all dependencies. Caching is leveraged here.
    RUN go mod download
    
    # Copy the source from the host to the container
    COPY . .
    
    # Build the application
    RUN CGO_ENABLED=0 GOOS=linux go build -o main .
    
    # --- Final Stage ---
    FROM alpine:latest
    
    # Set working directory
    WORKDIR /app
    
    # Copy the binary from the builder stage
    COPY --from=builder /app/main .
    
    # Expose the port the app listens on
    EXPOSE 9008
    
    # Run the executable
    CMD ["./main"]