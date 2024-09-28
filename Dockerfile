# FROM golang:1.23-alpine AS builder
# WORKDIR /build
# ADD go.mod .
# ADD go.sum .
# RUN go mod download
# COPY . .
# RUN go build -o /bin/quiz

# FROM alpine
# WORKDIR /app
# COPY --from=builder /bin/quiz /app/quiz
# EXPOSE 8080

# CMD ["./quiz"]


FROM golang:1.23

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]
