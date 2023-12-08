FROM debian:bookworm-slim

RUN apt-get update

RUN apt-get install --no-install-recommends -y \
    golang \
    binutils \
    gcc \
    libc6-dev \
    default-jdk

RUN mkdir /app
WORKDIR /app
ENV CGO_CFLAGS="-I/usr/lib/jvm/default-java/include -I/usr/lib/jvm/default-java/include/linux"

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

CMD ./testing.sh


