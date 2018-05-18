FROM alpine:3.6

RUN adduser -D deployment-test-operator
USER deployment-test-operator

ADD tmp/_output/bin/deployment-test-operator /usr/local/bin/deployment-test-operator
