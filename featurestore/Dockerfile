ARG ARTIFACT=mastro
FROM golang:1.15-alpine AS builder

ARG ARTIFACT

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

RUN apk update && apk upgrade && \
    apk add --no-cache git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ${ARTIFACT} .

FROM scratch

COPY --from=builder /build/${ARTIFACT} ./

ARG ARTIFACT
ENV ARTIFACT=${ARTIFACT}

# set default vars
ENV MASTRO_CONFIG=/conf
ENV GIN_MODE=release

# set config.yaml using wget or local copy
COPY conf $MASTRO_CONFIG
# copy binary
COPY --from=builder /build/${ARTIFACT} ./

# Command to run when starting the container
ENTRYPOINT ["sh", "-c", "./${ARTIFACT}"]