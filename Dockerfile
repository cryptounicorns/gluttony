FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/cryptounicorns/gluttony
COPY . .
RUN make

FROM fedora:latest

RUN mkdir            /etc/gluttony
COPY --from=builder /go/src/github.com/cryptounicorns/gluttony/build/gluttony /usr/bin/gluttony

RUN mkdir /var/gluttony

CMD [                                      \
    "/usr/bin/gluttony",                   \
    "--config",                            \
    "/etc/gluttony/config.json"            \
]
