FROM golang:1.19.1-alpine3.16 as build

ENV GO111MODULE=on

WORKDIR /go/src/hello
COPY . .

RUN go build -o /go/bin/hello

FROM alpine:3.16

COPY --from=build /go/bin/hello /go/bin/

CMD [ "/go/bin/hello" ]