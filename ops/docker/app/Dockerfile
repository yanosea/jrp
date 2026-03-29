FROM golang:1.26
WORKDIR /jrp
COPY . .
RUN go mod download
CMD ["make", "test.local"]
