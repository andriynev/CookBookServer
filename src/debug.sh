#!/bin/bash
while true
do
        inotifywait -e attrib /app/api
        echo "run debugger"
        sleep 3
        echo $(pidof api)
        dlv --headless --listen=:40000 --api-version=2 attach $(pidof api)
done

