#!/bin/sh
for i in `seq 1 11`
do
  ./myAwesomeBot -team=away -number=$i -wshost=$1 &
  ./myAwesomeBot -team=home -number=$i -wshost=$1 &
done



