FROM golang:1.15-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o mvc .

FROM scratch

COPY --from=builder /build/mvc ./

ENTRYPOINT ["./mvc"]