FROM golang:1.17
WORKDIR /root
COPY . ./
ENV CGO_ENABLED=0
RUN go build -tags nosystemd

FROM scratch
COPY --from=0 /root/go-import-redirect /go-import-redirect

EXPOSE 8080
ENTRYPOINT ["/go-import-redirect", "-addr", ":8080"]
