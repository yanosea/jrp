FROM golang:1.24
WORKDIR /jrp
COPY . .
RUN go mod download
CMD ["make", "test.local"]
