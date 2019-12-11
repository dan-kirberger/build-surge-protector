FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD build-surge-protector /bin/
ENTRYPOINT ["/bin/build-surge-protector"]