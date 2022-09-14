FROM alpine:3.7

COPY ./dojoClient /bin/dojoClient

RUN adduser -D -g '' golang

USER golang

ENTRYPOINT ["/bin/dojoClient"]
CMD [ "-h" ]