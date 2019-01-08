#!/bin/bash 

ROOTDIR=/home/[USERNAME]/Documents/PROGRAMMING/go/src/github.com/[GITHUB-NAME]/jblastor
COUNTER=0
while [ $COUNTER -lt 25000 ]; do
  $ROOTDIR/jblastor --files /var/tmp/JSON/ --endpoint 'http://localhost:8081/server/save'
  let COUNTER=COUNTER+1
done

