FROM golang:alpine as BUILDER
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN apk update && apk add --no-cache git
WORKDIR /usr/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o isit-a-cat-bff .

FROM alpine:latest
COPY --from=BUILDER /usr/src/app/isit-a-cat-bff .
ENTRYPOINT ["./isit-a-cat-bff"]
