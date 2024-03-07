# Stage 1: Build the application
FROM golang:1.22 AS builder

LABEL maintainer="Elias Bourgess <elias@ebourgess.dev>"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o todolist .

# Stage 2: Create a smaller image to run the application
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todolist ./
RUN chmod +x todolist
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["/app/todolist"]

