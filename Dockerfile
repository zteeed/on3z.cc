FROM golang:1.18 AS build
WORKDIR /go/src
COPY go ./go
COPY go.mod .
COPY go.sum .
COPY app.go .
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o main .

FROM golang:1.18 AS runtime
RUN apt-get update && \
    apt-get install -yq tzdata && \
    ln -fs /usr/share/zoneinfo/Europe/Paris /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata
ENV TZ="Europe/Paris"
COPY --from=build /go/src/main ./
EXPOSE 8888/tcp
ENTRYPOINT ["./main"]
