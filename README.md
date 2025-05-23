# URL Shortener

Goal of this project is shortening an URL using Go language and its Fiber backend framework:

## How to run:

1. Clone the repository with `git clone github.com/kelma01/url_shortener`
2. Make sure that your PC has Go. Type `go version` in terminal to check, if it is not installed, (in Linux) you can type `sudo snap install go` on terminal to install.
3. Get the Fiber backend development framework by typing `go get -u github.com/gofiber/v2`
4. Install PostgreSQL to manage DB relations with `sudo apt install postgresql` and create your user in there with the name of database you will use. Start service with `sudo systemctl start postgresql` With `sudo -u postgres psql`, enter the PostgreSQL terminal and `\c url_shortener` to use the db. Do not forget to change DB configuration files in code.
5. Install Redis with `sudo apt install redis-tools` and `sudo apt install redis-server`. Then start with `sudo systemctl start redis-server`. After that with `redis-cli`, the redis terminal can be used.
6. With `go run cmd/server/main.go` you can run the project in your localhost. Host is served on `localhost:8080`.
7. Dockerized. Run `./start_docker`. To access redis client terminal in dockerized project, run `redis-cli -h <ip of your docker machine> -p 6380`. To see docker machine ip, `docker network inspect bridge`.(If your PostgreSQL service in background, turn it off with `sudo systemctl stop postgresql`)


Kısayollar:

curl get req:
curl -X GET http://localhost:8080

curl post payload: 
curl --header "Content-Type: application/json" --request POST --data '{"original_url": "https://www.google.com"}' http://localhost:8080

curl delete paylaod:
curl -X DELETE http://localhost:8080/shorturl




miniube start
eval $(minikube docker-env)
docker build -t <image_name> ./path-to-your-app
kubectl apply -f k8s/
docker image ile imageler gönrüntlkenebilir
kubectl get svc ile servis isimleri görülebilir
minikube service <servi-adi> ile deploy edilir.
kubectl get pods ve kubectl logs <servis_name>


minikube deploy etme olayı

systemctl stop postgresql
systemctl stop redis
minikube start --driver=docker
eval $(minikube -p minikube docker-env)
docker build -t urltest .
kubectl apply -f /k8s
docker build -t urltest:latest . //????
kubectl apply -f /k8s
minikube service url-shortener-app