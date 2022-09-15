FROM golang:1.9.2-alpine

RUN adduser -D -g '' golang

COPY ./dojoClient /bin/dojoClient

RUN chmod +x /bin/dojoClient

USER golang

ENTRYPOINT ["/bin/dojoClient"]
CMD [ "-h" ]