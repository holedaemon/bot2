FROM golang:1.21.1 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

ARG version
RUN go build -ldflags="-X github.com/holedaemon/bot2/internal/version.version=${version}" -o bot2

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /app/bot2 /bot2
ENTRYPOINT [ "/bot2" ]