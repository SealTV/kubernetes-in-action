#!/bin/bash
trap "exit" SIGINT
mkdir /var/htdocs

echo Configure to generate fortune every $INTERVAL seconds

while :
do
    echo $(date) Writing fortune to /var/htdocs/index.html
    /usr/games/fortune > /var/htdocs/index.html
    sleep $INTERVAL
done