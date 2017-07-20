FROM golang:latest

VOLUME /data
EXPOSE 80
EXPOSE 443

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/cmd/redirector
RUN go get -v -d -u . && go install -v .

WORKDIR /data
CMD ["/go/bin/redirector"]
