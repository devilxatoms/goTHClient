FROM golang:1.9.2-alpine


COPY ./dojoClient /
RUN chmod +x /dojoClient

ENTRYPOINT ["/dojoClient"]
CMD [ "-h" ]