FROM alpine:3.7

COPY ./dojoClient /bin/dojoClient

ENTRYPOINT ["/bin/dojoClient"]
CMD [ "-h" ]