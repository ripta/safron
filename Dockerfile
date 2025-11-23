FROM golang:1.24-bookworm AS builder

ENV PROJECT=github.com/ripta/safron
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /safron .


FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /safron /safron
VOLUME ["/data"]
EXPOSE 8080
CMD ["/safron", "-path=/data", "-log-format=json"]

