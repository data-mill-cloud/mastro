ARG SOURCE_FOLDER=metricstore
ARG ARTIFACT=metricstore
FROM golang:1.18-alpine AS builder

ARG SOURCE_FOLDER
ARG ARTIFACT

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY commons commons
COPY ${SOURCE_FOLDER} ${ARTIFACT}

WORKDIR /build/${ARTIFACT}
RUN go mod download
RUN go build -o ${ARTIFACT} .

#FROM scratch
FROM alpine:3.12.4

ARG SOURCE_FOLDER
ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# set default vars
ENV MASTRO_CONFIG=/conf
ENV GIN_MODE=release

COPY ${SOURCE_FOLDER}/conf $MASTRO_CONFIG
COPY --from=builder /build/${ARTIFACT}/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]