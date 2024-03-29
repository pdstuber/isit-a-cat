FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.21.4-bookworm as builder

ENV USER=appuser
ENV UID=10001

ARG TARGETARCH TARGETOS

ENV LIBTENSORFLOW_FILENAME="libtensorflow-2.14.1-${TARGETARCH}.tar.xz"

WORKDIR /app

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

ADD build/${LIBTENSORFLOW_FILENAME} /usr/local/

RUN apt update && apt install binutils
RUN strip -s /usr/local/lib/libtensorflow*

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

RUN ldconfig /usr/local/lib
COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY main.go main.go

RUN --mount=type=cache,target="/root/.cache/go-build" \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GOCACHE=/root/.cache/go-build go build -o /go/bin/isit-a-cat

FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bookworm-slim

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /usr/local/lib/libtensorflow* /usr/local/lib/
COPY --from=builder /usr/local/include/tensorflow /usr/local/include/

RUN ldconfig /usr/local/lib

COPY --from=builder /go/bin/isit-a-cat /go/bin/isit-a-cat

RUN mkdir /model
ADD build/model/* /model/

USER appuser:appuser

ENTRYPOINT ["/go/bin/isit-a-cat"]