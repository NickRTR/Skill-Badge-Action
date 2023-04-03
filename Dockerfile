FROM golang:1.20.2

WORKDIR /app

COPY ./ ./

RUN go build -o /bin/app ./

ENTRYPOINT ["app"]