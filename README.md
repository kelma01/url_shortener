# URL Shortener

Goal of this project is shortening an URL using Go language and its Fiber backend framework:

## How to run on your "localhost:8080" :

1. Clone the repository with `git clone github.com/kelma01/url_shortener`
2. Make sure that your PC has Go. Type `go version` in terminal to check, if it is not installed, (in Linux) you can type `sudo snap install go` on terminal to install.
3. Get the libraries that project needs
4. Install PostgreSQL to manage DB relations with `sudo apt install postgresql` and create your user in there with the name of database you will use. Start service with `sudo systemctl start postgresql` With `psql -u <username> -d url_shortener`, enter the PostgreSQL terminal and `\c url_shortener` to use the db. Do not forget to change DB configuration files in code.
5. Install Redis with `sudo apt install redis-tools` and `sudo apt install redis-server`. Then start with `sudo systemctl start redis-server`. After that with `redis-cli`, the redis terminal can be used.
6. With `go run cmd/server/main.go` you can run the project in your localhost. Host is served on `localhost:8080`.


## API Requests for Copying :


```bash
#GET REQUEST
curl -X GET http://localhost:8080

#POST REQUEST: 
curl --header "Content-Type: application/json" --request POST --data '{"original_url": "https://www.google.com"}' http://localhost:8080

#DELETE REQUEST:
curl -X DELETE http://localhost:8080/shorturl
```


## Deploying with Docker :

```bash
sudo docker-compose down
sudo docker rm -f $(sudo docker ps -aq) 
sudo docker volume prune -f             
sudo docker network prune -f            
sudo docker-compose up --build
```

## Deploying with Minicube :

```bash
systemctl stop postgresql
systemctl stop redis
minikube start --driver=docker
eval $(minikube -p minikube docker-env)
docker build -t urltest:latest .
kubectl apply -f /k8s
minikube service url-shortener-app
```