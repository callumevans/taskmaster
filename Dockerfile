FROM golang:latest as build
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/taskmaster
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o taskmaster .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/taskmaster/taskmaster .
CMD ["./taskmaster"]