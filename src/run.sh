#!/bin/bash
while true
do
        inotifywait -e attrib ./src/api/*
        echo "notify"
        echo -n > /app/src/rebuild.xml
done

