FROM postgres:18-alpine3.24

ENV POSTGRES_ROOT_PASSWORD=root
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=sampledb


COPY ./sampledb.sql /docker-entrypoint-initdb.d/

EXPOSE 5432