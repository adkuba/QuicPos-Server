FROM golang:1.15.2-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin .
ENTRYPOINT ["/bin/QuicPos"]