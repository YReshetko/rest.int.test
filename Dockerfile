FROM  byrnedo/alpine-curl:latest

COPY ./run /usr/local/bin/run

ENTRYPOINT ["/usr/local/bin/run"]