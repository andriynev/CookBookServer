#!/bin/bash
cat /docker-entrypoint-initdb.d/init-schema.sql | mysql -uroot -proot