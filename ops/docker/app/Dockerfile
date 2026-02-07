FROM golang:1.25
WORKDIR /jrp
COPY . .
RUN go mod download
CMD ["make", "test.local"]
