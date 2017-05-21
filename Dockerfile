FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY besto /

EXPOSE 8014

ENTRYPOINT ["/besto", "--data", "/data/besto.db"]

