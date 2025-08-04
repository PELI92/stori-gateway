FROM golang:1.24

WORKDIR /app

COPY . .

RUN go build -o mock .

EXPOSE 5678

CMD ["./mock"]
