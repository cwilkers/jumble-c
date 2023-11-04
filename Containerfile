FROM golang
WORKDIR /app

COPY words *.go *.html ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /jumble

EXPOSE 8080

CMD ["/jumble"]