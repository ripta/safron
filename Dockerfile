FROM golang:1.10-alpine3.7 AS builder

ENV PROJECT=github.com/ripta/safron
ENV CGO_ENABLED=0

# RUN mkdir -p $GOPATH/src/$PROJECT
COPY . $GOPATH/src/$PROJECT

RUN apk add --update --no-cache git \
    && go get $PROJECT \
    && apk del git
RUN go build -o /safron $PROJECT


FROM scratch
COPY --from=builder /safron /safron
VOLUME ["/data"]
EXPOSE 8080
CMD ["/safron", "-path=/data", "-log-format=json"]

