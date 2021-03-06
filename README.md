# Delivery Test

### Prerequisite
The project requires Docker. Make sure you have Docker and Docker-compose installed.

### Run The Project

##### Configs:
```sh
$ cp .env-example .env
```

Input your Google Maps API key in the `GOOGLE_MAP_API_KEY` field in `.env` file.


##### Start the server:
```sh
$ ./start.sh
```

The command should take around 1-2 minute to set up. The server will be listening to PORT 8080 by default.
You may change the port by changing the `APP_PORT` field in `.env` file.


### Stop The Project

```sh
$ docker-compose down
```


### Run Unit Tests
```sh
$ go test ./...
```

### Brief Explaination

#### Code Architecture
The project tries to follow clean architecture design pattern (https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). I found it makes writing unit tests very easy.

