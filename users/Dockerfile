FROM golang:1.22-alpine3.19 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/main

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]