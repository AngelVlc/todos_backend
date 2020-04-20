#!/bin/sh
echo "Executing database migrations ..."
rm -f conn.txt
echo $CLEARDB_DATABASE_URL > conn.txt 

sed -i "s|@|@tcp(|g" conn.txt
sed -i "s|\/|:3306)\/|3" conn.txt
sed -i "s|?reconnect=true|?query|g" conn.txt
conn=`cat conn.txt`
/app/migrate.linux-amd64 -database "$conn"  -path /app/db/migrations up
echo "Executing app ..."
/app/app