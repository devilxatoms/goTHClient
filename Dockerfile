FROM golang:latest


COPY ./dojoClient /
RUN chmod +x /dojoClient

ENTRYPOINT ["/dojoClient"]
CMD [ "-h" ]