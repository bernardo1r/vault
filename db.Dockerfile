FROM postgres
COPY backend/queries/tables.sql /docker-entrypoint-initdb.d/
