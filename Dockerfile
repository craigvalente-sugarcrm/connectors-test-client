FROM golang:1.13.7-alpine3.11 as builder

RUN mkdir /builder
ADD . /build/
WORKDIR /build

# RUN go build -o main .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.11
COPY --from=builder /build/main /app/
WORKDIR /app

ENV DURATION=1
ENV RATE=2

CMD ["./main"]
