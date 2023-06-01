### phanes-layout Document
    
    []: # Language: markdown
    []: # Path: phanes-layout/README.md

### Running prerequisite
> you must have docker installed, and start the containers which include Postgres, Redis, Etcd, RabbitMQ or Nats.
1. run `docker-compose up -d` to start the containers if you don't have base service components
2. modify the /scripts/config.json to application-compliant
3. run `make config` put /scripts/config.json to etcd, based on you have installed `ectdctl`

### Start command:
```shell
go build -o phanes 
./phanes --registry=etcd --registry_address=127.0.0.1:2379
./phanes --config=./script/config.yaml
```

#### project dir structure description
- [bll] Business logic layer
- [client] If server need request other service, the client can do. provide grpc, http, websocket etc.
- [server] Means the server provide service by grpc, http, websocket etc
- [config] Provide config file and other init
- [event] Provide event bus
- [errors] Provide unified external error handling
- [collector] Provide log, trace and metric collect
- [model] Provide model and entity
- [util] General tools
- [store] Provide storage
- [lib] Self-defined library
- [script] Some config script and sql script

#### Dependencies
- [go-micro](https://github.com/asim/go-micro)
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)

