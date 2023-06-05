FROM golang:1.20 as build

RUN apt-get update && \
    apt-get install -y ca-certificates libssl-dev cpio
COPY --from=prom/prometheus /bin/promtool /go/bin/promtool

WORKDIR /app
ADD . /app

RUN make check build

FROM busybox

COPY --from=build /app/bin/po-test /go/bin/po-test
COPY --from=prom/prometheus /bin/promtool /go/bin/promtool
ENV PATH $PATH:/go/bin
#Required to communicate with pubsub... and likely anything external
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/po-test"]