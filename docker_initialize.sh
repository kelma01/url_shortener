#!/bin/bash

sudo docket-compose down
sudo docker rm $(sudo docker ps -aq)
sudo docker volume prune
sudo docker network prune
sudo docker-compose up --build