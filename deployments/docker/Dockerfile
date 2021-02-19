FROM golang:1.15
WORKDIR /root
COPY . ./
ENV CGO_ENABLED=0
RUN go build ./cmd/go-import-redirect

FROM scratch
COPY --from=0 /root/go-import-redirect /go-import-redirect

EXPOSE 8080
ENTRYPOINT ["/go-import-redirect"]
