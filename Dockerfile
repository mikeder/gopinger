FROM golang:onbuild AS fetch 

ADD . / 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gopinger .
RUN ["./gopinger", "-version"]

FROM alpine:3.8 
RUN apk --no-cache add ca-certificates

COPY --from=fetch /go/src/app/gopinger /bin/gopinger

EXPOSE 3001
CMD ["gopinger"]
