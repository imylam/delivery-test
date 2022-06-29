# Delivery Test

## Prerequisite
The project requires Docker. Make sure you have Docker and Docker-compose installed.

## Run The Project

#### Configs:
Input your Google Maps API key for the field `GOOGLE_MAP_API_KEY` in `docker-compose.yml`:

```
docker-compose.yml

services:
  app:
    ...
    environment:
      ...
      - GOOGLE_MAP_API_KEY=<YOUR GOOLGE MAP API KEY HERE>
```

#### Start the server:
```sh
$ ./start.sh
```

The command should take around 30 seconds to set up. The server will be listening to PORT 8080 by default.
You may change the port by changing the `APP_PORT` field in `docker-compose.yml` file.


## Stop The Project

```sh
$ docker-compose down -v
```


## Run Unit Tests
```sh
$ go test ./...
```

## Run Integration Tests
> Note: For integration tests to pass, a fresh DB is needed.

#### Run Integration Tests locally
```sh
$ ./start.sh
cd ../integration_tests
go test ./... -tags=integration
```

#### Run Integration Tests in docker
```sh
$ ./run-integration-test.sh
```

### Brief Explaination
The project tries to follow clean architecture design pattern (https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

