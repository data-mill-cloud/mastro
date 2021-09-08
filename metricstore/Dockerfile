ARG ARTIFACT=metricstore
FROM golang:1.15-alpine AS builder

ARG ARTIFACT

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ${ARTIFACT} .

#FROM scratch
FROM alpine:3.12.4

ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# set default vars
ENV MASTRO_CONFIG=/conf
ENV GIN_MODE=release

COPY conf $MASTRO_CONFIG
COPY --from=builder /build/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]