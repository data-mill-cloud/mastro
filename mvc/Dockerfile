FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY commons commons
COPY mvc mvc

WORKDIR /build/mvc
RUN go mod download
RUN go build -o mvc .

FROM scratch
COPY --from=builder /build/mvc/mvc ./
ENTRYPOINT ["./mvc"]