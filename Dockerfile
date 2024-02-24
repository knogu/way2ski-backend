FROM golang:latest

WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . .

RUN go build -o main .
EXPOSE 8080
ENV PORT 8080
ENV HOSTNAME "0.0.0.0"
CMD ["go", "run", "main.go"]