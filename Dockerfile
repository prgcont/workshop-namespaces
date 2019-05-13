FROM golang:1.11 as builder

WORKDIR /go/src/github.com/prgcont/workshop-namespaces/
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/workshop-namespaces ./cmd/api
RUN useradd goappp -u 10001 && \
  chown 10001:10001 /go/bin/workshop-namespaces

FROM scratch
COPY --from=builder /go/bin/workshop-namespaces /workshop-namespaces
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./static /static
USER 10001
EXPOSE 9090
ENTRYPOINT ["/workshop-namespaces"]
