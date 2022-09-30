#!/bin/sh
echo "Executing database migrations ..."
/app/migrate.linux-amd64 -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?tls=$MYSQL_TLS&query" -path /app/db/migrations up

echo "Executing app ..."
/app/app