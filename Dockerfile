FROM golang:1.20
WORKDIR /
COPY . .

RUN /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

EXPOSE 80

CMD ["./main"]
