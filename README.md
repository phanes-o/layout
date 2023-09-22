### phanes-layout Document
    
    []: # Language: markdown
    []: # Path: phanes-layout/README.md

### Running prerequisite
To get started, ensure that you have Docker installed and running, with containers for Postgres, Redis, Etcd, RabbitMQ, or Nats.

Follow these steps:

Start Base Service Components: If you haven't already set up the base service components, run the following command to start the containers:
```shell
docker-compose up -d
```
Configure the Application: Modify the /scripts/config.json file to make it compliant with your application's requirements.

Update Etcd Configuration: Run the following command to put the updated /scripts/config.json file into Etcd. Please make sure you have `etcdctl` installed.

```shell
make config
```

By following these steps, you'll have your environment set up and configured to work with your application.


### Start command:

```shell
go build -o phanes 
```

Run with Etcd registry and Etcd configuration center
```shell
./phanes --registry=etcd --registry_address=127.0.0.1:2379 
```

Run with local config file and Etcd registry
```shell
./phanes --registry=etcd --registry_address=127.0.0.1:2379  --config=./script/config.yaml
```

Run with local config file and mdns registry
> If you just want running on local environment, you can use the following command
```shell
./phanes --registry=mdns  --config=./script/config.yaml
```

#### project dir structure description
[bll] Business Logic Layer: This layer handles the core business logic.

[client] Client: This component allows the server to request other services and provides support for communication protocols like gRPC, HTTP, WebSocket, etc.

[server] Server: This refers to the server component responsible for providing services using protocols such as gRPC, HTTP, WebSocket, etc.

[config] Configuration: This module provides configuration files and other initialization settings.

[event] Event Handling: This component offers an event bus for managing events within the system.

[errors] Error Handling: This module provides a standardized approach to handling external errors.

[collector] Data Collection: This component collects logs, traces, and metrics.

[model] Models and Entities: This section provides the necessary data models and entities.

[util] Utilities: This includes general-purpose utility tools.

[store] Storage: This module handles storage-related functionalities.

[lib] Custom Libraries: This category consists of self-defined libraries.

[script] Configuration and SQL Scripts: This section contains various configuration and SQL scripts.

#### Dependencies
- [go-micro](https://github.com/asim/go-micro)
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://github.com/go-gorm/gorm)

