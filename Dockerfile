FROM golang:1.11 AS build
RUN mkdir -p /go/src/github.com/cs3238-tsuzu/arrsup
WORKDIR /go/src/github.com/cs3238-tsuzu/arrsup/
COPY ./ /go/src/github.com/cs3238-tsuzu/arrsup/
ENV GO111MODULE=on
ENV CGO_ENABLED=0
RUN go build . 

FROM scratch
COPY --from=build /go/src/github.com/cs3238-tsuzu/arrsup/arrsup /bin/
ENTRYPOINT [ "/bin/arrsup" ]
CMD [ "--help" ]