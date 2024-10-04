FROM golang:1.22
WORKDIR /app
COPY . /app
RUN go get -d -v ./...
RUN go build -o main .
EXPOSE 8000
CMD ["/app/main"]
