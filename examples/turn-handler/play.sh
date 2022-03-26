#!/usr/bin/env bash

if [ -z "$1" ]
  then
    echo "Please, pass the first argument (home or away) to set the team side"
    exit 1
fi

go build -o myAwesomeBot main.go || { echo "building has failed"; exit 1; }
for i in `seq 1 11`
do
  ./myAwesomeBot -team=$1 -number=$i &
  sleep 0.1
done
