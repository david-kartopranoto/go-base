FROM library/postgres
COPY db/*.sql /docker-entrypoint-initdb.d/