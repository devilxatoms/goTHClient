FROM alpine:3.7

RUN adduser -D -g '' golang

COPY ./dojoClient /bin/dojoClient

RUN chmod +x /bin/dojoClient

USER golang

ENTRYPOINT ["/bin/dojoClient"]
CMD [ "-h" ]