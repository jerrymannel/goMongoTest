FROM alpine:3

COPY goMongoTest-x86 /

ENTRYPOINT ["./goMongoTest-x86"]