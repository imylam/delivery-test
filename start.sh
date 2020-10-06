source .env

docker-compose up -d 

echo "waiting MariaDB to be initialized..."
sleep 30

echo "setting up the db..."
docker exec -i app-mariadb-container mysql -uroot -p$MYSQL_ROOT_PASSWORD -e "CREATE USER IF NOT EXISTS 'delivery'@'%' IDENTIFIED BY \"$MYSQL_PASSWORD\";"
docker exec -i app-mariadb-container mysql -uroot -p$MYSQL_ROOT_PASSWORD < dbSetup.sql

docker-compose up -d 