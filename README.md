# Microservices

This implements a simple TODO list microservices using jwt authentication.


## Installation

The easiest way to install it is to have the [latest bleeding edge Docker engine](https://github.com/docker/docker/releases/tag/v1.10.0-rc3) and the [latest bleeding edge Docker compose](https://github.com/docker/compose/releases/tag/1.6.0-rc2). (We need them because of Docker compose latest features requiring an engine running a v1.22 API. The features enable us automatic host discovery through a docker network, making the whole process easier).

Once you have both, the way to go is :

```
go get -v github.com/dolanor/microservices
cd $GOPATH/src/github.com/dolanor/microservices
make # It will build each microservice binary and put them in a container
docker-compose up
```

Then, fire your browser and go to [http://localhost:8080](http://localhost:8080)
Or alternatively, you could `go test -v ./...`

I used a Makefile to simplify and fasten the development process. If I used a golang:onbuild base image, every code modification would imply a new `go get -v` for all the dependencies. With a local build, it's way faster but could add some more dependency bugs. We should add a PROD Dockerfile for that case.


## Docker compose

Docker compose helps you putting all the components together.
A very nice addition in the context of microservices: scaling.

`docker-compose scale auth=3 todo=2` will create 3 containers for authentication services and 2 for todo services.

It makes scaling your services very simple without the need of a full blown orchestration solution for a simple system. One backward is that it is slow (I tried to scale 20 auth, 10 todo and 10 data, and it took a few minutes to scale and more time to launch.
