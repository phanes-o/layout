### phanes-layout Document

[]: # Language: markdown
[]: # Path: phanes-layout/README.md


### Running prerequisite
> you must have docker installed, and start the containers which include Postgres, Redis, Etcd, RabbitMQ or Nats.
1. run `docker-compose up -d` to start the containers if you don't have base service components
2. modify the /scripts/config.json to application-compliant
3. run `make config` put /scripts/config.json to etcd, based on you have installed `ectdctl`

### Srart command:
```shell
go build -o phanes 
./phanes --registry=etcd --registry_address=127.0.0.1:2379
```

[]: # Language: bash
[]: # Path: phanes-layout/README.md

#### project dir structure description
- bll business logic layer
- client provide grpc, http, websocket etc.
- server serve grpc, http, websocket
- config provide config file and other init
- event provide event bus
- collector provide log, trace and metric collect
- model provide model and entity
- util general tools
- store provide storage
- lib self-defined library
- script some config script and sql script


#### Feature
- [x]  Load Balancing - Client side load balancing built on service discovery. Once we have the addresses of any number of instances of a service we now need a way to decide which node to route to. We use random hashed load balancing to provide even distribution across the services and retry a different node if there's a problem.

- [x] Service Discovery - Automatic service registration and name resolution. Service discovery is at the core of micro service development. When service A needs to speak to service B it needs the location of that service. The default discovery mechanism is multicast DNS (mdns), a zeroconf system.

- [x] Data Storage - A simple data store interface to read, write and delete records. It includes support for memory, file and CockroachDB by default. State and persistence becomes a core requirement beyond prototyping and Micro looks to build that into the framework.

- [x] Dynamic Config - Load and hot reload dynamic config from anywhere. The config interface provides a way to load application level config from any source such as env vars, file, etcd. You can merge the sources and even define fallbacks.

- [x] Event Bus - A simple event bus interface to publish and subscribe to events. 

- [x] Metric Collector - A simple metric collector interface to collect metrics.

- [x] Log Collector - A simple log collector interface to collect logs.

- [x] Trace Collector - A simple trace collector interface to collect traces.

#### Dependencies
- [go-micro](https://github.com/asim/go-micro)
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)

