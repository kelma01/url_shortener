# URL Shortener

Goal of this project is shortening an URL using Go language and its Fiber backend framework:

## How to run:

1. Clone the repository with `git clone github.com/kelma01/url_shortener`
2. Make sure that your PC has Go. Type `go version` in terminal to check, if it is not installed, (in Linux) you can type `sudo snap install go` on terminal to install it
3. Install the Fiber framework by typing `go get -u github.com/gofiber/v2`
4. Install PostgreSQL to manage DB relations with `sudo apt install postgresql` and create your user in there with the name of database you will use. With `sudo -u postgres psql`, enter the PostgreSQL terminal and `\c url_shortener` to use the db.
5. With `go run cmd/server/main.go` you can run the project in your localhost. Host is served on `localhost:8080`.
