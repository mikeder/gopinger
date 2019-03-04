FROM golang:onbuild AS fetch 

ADD . / 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gopinger .
RUN ["./gopinger", "-version"]

FROM scratch

COPY --from=fetch /go/src/app/gopinger /bin/gopinger
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 3001
CMD ["gopinger"]
