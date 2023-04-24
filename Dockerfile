FROM golang:1.19

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the app
RUN go build -o main .

# Set the CMD instruction to run the app with the APP_NAME argument
CMD ["./main"]
