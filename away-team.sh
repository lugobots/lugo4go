#!/bin/sh
for i in `seq 1 11`
do
  go run main.go -team=away -number=$i&
  go run main.go -team=home -number=$i&
  sleep 1
done



