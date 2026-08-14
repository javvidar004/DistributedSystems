FROM postgres:18-alpine3.24

COPY ./sampledb.sql /docker-entrypoint-initdb.d/

ENV POSTGRES_ROOT_PASSWORD=root
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=sampledb

