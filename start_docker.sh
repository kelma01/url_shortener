#!/bin/bash

sudo docker-compose down
sudo docker rm -f $(sudo docker ps -aq) 
sudo docker volume prune -f             
sudo docker network prune -f            
sudo docker-compose up --build