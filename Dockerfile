FROM alpine

EXPOSE 8090

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD cluster_name /go/bin/cluster_name

RUN addgroup -S cngroup && adduser -S cnuser -G cngroup

USER cnuser

ENV CLUSTER_NAME_USER blah
ENV CLUSTER_NAME_PASS blah

ENTRYPOINT ["/go/bin/cluster_name"]