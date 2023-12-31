# SERVICETOOLS

This is a highly opinionated, experimental, proof of concept
library to help build services/microservices in Go.

This is very unstable as it's constantly being experimented
during interactions of other projects.

Documentation and tests will come as soon as I can decide on
what this should look like, and if it matures to that point.

## Concept

The main goal is to offer modular components,
so that services can be built from one or more.

e.g.

```go
type MyService struct{
    *server.WithDB
    *server.WithGRPC
}
```

Currently, the components are:

* WithDB: takes a database configuration and exposes `DB()` (gorm)
* WithRDB: sames as WithDB, but meant for "readonly" access. Exposes `RDB()`
* WithGRPC: starts a gRPC server internally and mounts all gRPC services that are given to it.
* WithHTTP: starts a HTTP server internally and mounts all handlers that are given to it.
* WithHealthcheck: mounts a basic HTTP server for healthchecks and exposes metrics using opentelemetry.
* WithWorker: starts tasks in the background.
