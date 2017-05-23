FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

EXPOSE 8014

COPY besto /

ENTRYPOINT ["/besto", "--data", "/data/besto.db"]

