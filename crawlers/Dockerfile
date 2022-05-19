ARG SOURCE_FOLDER=crawlers
ARG ARTIFACT=crawlers
FROM golang:1.18-alpine AS builder

ARG SOURCE_FOLDER
ARG ARTIFACT

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add --no-cache build-base krb5 krb5-dev
#krb5-pkinit

WORKDIR /build
COPY commons commons
COPY ${SOURCE_FOLDER} ${ARTIFACT}

WORKDIR /build/${ARTIFACT}
RUN go mod download
RUN go build -o ${ARTIFACT} .

# Build the application
RUN go build -tags=kerberos -o ${ARTIFACT} .

# multistage build - we only copy the result (binary) into a fresh image which is super light
FROM alpine:3.12.4

ARG SOURCE_FOLDER
ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

RUN apk add --no-cache krb5 krb5-dev

# set default vars
ENV MASTRO_CONFIG=/conf
ENV GIN_MODE=release

# set config.yaml using wget or local copy
COPY ${SOURCE_FOLDER}/conf $MASTRO_CONFIG
COPY --from=builder /build/${ARTIFACT}/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]