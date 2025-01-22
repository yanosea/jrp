FROM golang:1.23
RUN apt-get update && apt-get install -y curl unzip
RUN curl -L -O https://github.com/sagiegurari/cargo-make/releases/download/0.37.14/cargo-make-v0.37.14-x86_64-unknown-linux-musl.zip \
  && unzip cargo-make-v0.37.14-x86_64-unknown-linux-musl.zip \
  && cp cargo-make-v0.37.14-x86_64-unknown-linux-musl/cargo-make /usr/local/bin/cargo-make \
  && rm -rf cargo-make-v0.37.14-x86_64-unknown-linux-musl.zip cargo-make-v0.37.14-x86_64-unknown-linux-musl
WORKDIR /jrp
COPY . .
RUN go mod download
CMD ["cargo-make", "make", "test.local"]
