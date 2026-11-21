#!/bin/sh
echo "Executing database migrations ..."
if [ -n "$MYSQL_UNIX_SOCKET" ]; then
  /app/migrate.linux-amd64 -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@unix($MYSQL_UNIX_SOCKET)/$MYSQL_DATABASE?query" -path /app/db/migrations up
else
  /app/migrate.linux-amd64 -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE?ssl-mode=PREFERRED&query" -path /app/db/migrations up
fi

echo "Executing app ..."
/app/app