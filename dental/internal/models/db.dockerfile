FROM mysql:8.0.31


COPY ./sql-scripts /docker-entrypoint-initdb.d/
RUN chown -R mysql:mysql /docker-entrypoint-initdb.d/
