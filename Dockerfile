FROM        golang:latest
RUN         mkdir -p /app
WORKDIR     /app
COPY        . .
RUN         go mod download
RUN         go mod tidy
RUN         go build -o app ./main.go
EXPOSE      8000
ENTRYPOINT  ["./app","--port=8000"]