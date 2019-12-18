FROM alpine

EXPOSE 8091

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD cluster_name /go/bin/cluster_name

RUN addgroup -S cngroup && adduser -S cnuser -G cngroup

USER cnuser

ENTRYPOINT ["/go/bin/cluster_name"]