FROM golang:1.8-alpine
ADD . /go/src/kube-app
RUN go install kube-app

FROM alpine:latest
COPY --from=0 /go/bin/kube-app .
ENV PORT 8080
CMD ["./kube-app"]
