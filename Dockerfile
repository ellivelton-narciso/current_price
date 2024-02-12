FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD=mysql
EXPOSE 3306
COPY ./scripts/ /docker-entrypoint-initdb.d/

CMD ["mysqld"]
