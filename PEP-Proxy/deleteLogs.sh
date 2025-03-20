#!/bin/bash

for container in $(docker ps -aq)
do

echo $container
echo "" > $(docker inspect --format='{{.LogPath}}' $container)

done
