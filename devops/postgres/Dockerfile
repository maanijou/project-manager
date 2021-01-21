FROM postgres:latest
COPY ./database/db.sql /fixtures/db.sql
COPY ./devops/postgres/setup.sh /docker-entrypoint-initdb.d/setup.sh
WORKDIR /fixtures
RUN chmod -R 777 /docker-entrypoint-initdb.d/setup.sh
