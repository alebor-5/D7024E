### Build stage ###
FROM golang:latest as builder

ADD ./src /go/src
WORKDIR /go/src
RUN CGO_ENABLED=0 go build -o main 

### New lightweight stage ###
FROM alpine:latest

COPY --from=builder /go/src/main .

#EXPOSE <portNumber>

CMD ["./main"]
