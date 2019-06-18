#!/bin/bash

/app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" version
/app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" up

CompileDaemon -include="*.xml" -build="go build -gcflags='all=-N' -o /app/api /app/src/api/" -command="/app/api" -directory="." -color -command-stop=true
