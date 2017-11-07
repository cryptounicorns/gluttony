FROM fedora:latest

RUN mkdir            /etc/gluttony
ADD ./build/gluttony /usr/bin/gluttony

CMD [                                      \
    "/usr/bin/gluttony",                   \
    "--config",                            \
    "/etc/gluttony/config.json"            \
]
