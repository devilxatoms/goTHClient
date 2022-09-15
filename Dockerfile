FROM golang:1.9.2-alpine


COPY ./dojoClient /bin/dojoClient
RUN chmod +x /bin/dojoClient

ENTRYPOINT ["/bin/dojoClient"]
CMD [ "-h" ]