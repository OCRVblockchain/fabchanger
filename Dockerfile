FROM golang:1.12 as build-stage
LABEL maintainer="Vadim Inshakov <vadiminshakov@gmail.com>"
WORKDIR /app
COPY . .
RUN GOSUMDB=off GOPROXY=direct go build -o fabchanger ./cmd/main.go

# production stage
FROM alpine:3.9 as production-stage
WORKDIR /app
# for this stage copy to current dir Fabric /bin directory with cryptogen utility
COPY generate.sh .
COPY sign.sh .
COPY config/ ./config
COPY bin /app/bin
COPY --from=build-stage /app/fabchanger .
COPY --from=build-stage /app/config .
ENV PATH $PATH/app/bin
RUN apk add --no-cache libc6-compat