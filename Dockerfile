FROM golang AS builder
LABEL maintainer "sanshirookazaki"

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o app ./

FROM alpine
COPY --from=builder /go/src/app/app /usr/local/bin/app
ENTRYPOINT ["/usr/local/bin/app"]