FROM debian

RUN mkdir /app

COPY ./static /app/static
COPY ./views /app/views
COPY ./dd-vote /app/dd-vote

EXPOSE 8081
WORKDIR /app
ENTRYPOINT ["/app/dd-vote"]
