###############################################################################################
# Go building
###############################################################################################

FROM golang AS build

WORKDIR /

COPY main.go /

RUN go get go.mongodb.org/mongo-driver/mongo \
    && go get go.mongodb.org/mongo-driver/bson \
    && env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o goMongoTest-x86 .

###############################################################################################
# final image
###############################################################################################

FROM alpine:3

COPY --from=build /goMongoTest-x86 /

ENTRYPOINT ["./goMongoTest-x86"]