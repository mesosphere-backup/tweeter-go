# Oinker-Go

Example Go (golang) web app with dependency injection and graceful shutdown. Acts like a mini Twitter clone.


## Dependencies

Building requires [Go](https://golang.org/doc/install), [Godep](https://github.com/tools/godep), [Docker](https://docs.docker.com/installation/), &amp; Make.

Code dependencies are vendored with Godep:

- [Facebook's Grace library](http://github.com/facebookgo/grace) - graceful shutdown
- [Inject](http://github.com/karlkfi/inject) - dependency injection
- [Humanize](http://github.com/dustin/go-humanize) - readable units
- [Logrus](http://github.com/Sirupsen/logrus) - structured, leveled logging


## Compilation

Build a local binary:

```
make
```

Build a docker image:

```
make image
```


## Operation

There are several ways to launch the Oinker server.

### Source

Run from local source code:

```
go run main.go
```

(ctrl-c to quit)

### Docker

Run in Docker:

```
docker run -d --name oinker -p 0.0.0.0:8080:8080 karlkfi/oinker-go:latest
```

With Cassandra:

```
docker run -d --name cassandra cassandra:2.2.3
docker run -d --name oinker --link cassandra:cassandra -p 0.0.0.0:8080:8080 karlkfi/oinker-go:latest --cassandra-addr=cassandra
```

Find Oinker IP:

```
docker inspect --format "{{.NetworkSettings.IPAddress}}" oinker
```

### Marathon

Run in [Marathon](https://mesosphere.github.io/marathon/):

```
curl -H 'Content-Type: application/json' -X POST -d @"marathon.json" ${MARATHON_URL}/v2/apps
```

### Kubernetes

Run in [Kubernetes](http://kubernetes.io/):

```
kubectl create -f kubernetes.yaml
```


## Usage

Visit the home page at http://localhost:8080/

Enter handle &amp; message &amp; hit Oink.

See past oinks on the right-hand side of the home page.


## Future

- DNS SRV record Cassandra discovery
- Analytics page


## License

   Copyright 2015 Karl Isenberg

   Licensed under the [Apache License Version 2.0](LICENSE) (the "License");
   you may not use this project except in compliance with the License.

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
