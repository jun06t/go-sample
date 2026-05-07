# go-sample

A personal collection of small, self-contained Go samples that exercise language
features, the standard library, third-party libraries, and surrounding
middleware. Each top-level directory is an independent topic that can be
run on its own with minimal setup.

Module: `github.com/jun06t/go-sample` (Go 1.24)

## Index

### Language / Tooling

| Directory | Description |
| --- | --- |
| [cli-interface](./cli-interface) | Building a CLI with the `flag` package (uppercase converter) |
| [cfs](./cfs) | Concurrent prime search used as a runtime-tracing demo |
| [pipe](./pipe) | A `head`-like tool driven by Unix-style stdin pipes |
| [unix](./unix) | HTTP server over a Unix domain socket |
| [shutdown-hook](./shutdown-hook) | Graceful shutdown triggered by OS signals |
| [modules](./modules) | Go Modules basics |
| [toolchain](./toolchain) / [toolchain-library](./toolchain-library) | Experiments with the `toolchain` directive in `go.mod` |
| [go-cmp](./go-cmp) | Struct comparison and mock generation with `google/go-cmp` |
| [yaml](./yaml) | YAML parsing with environment-variable expansion |

### HTTP / Networking

| Directory | Description |
| --- | --- |
| [http2](./http2) | HTTP/2 server and client behavior |
| [keep-alive](./keep-alive) | Verifying HTTP keep-alive behavior |
| [sse](./sse) | Server-Sent Events server and client |
| [compress](./compress) | Response compression with memory profiling |
| [request-body](./request-body) | (placeholder) |
| [error-middleware](./error-middleware) | HTTP error-handling middleware |
| [log-response](./log-response) | Capturing HTTP responses via middleware |

### Logging / Observability

| Directory | Description |
| --- | --- |
| [zap](./zap) | `uber-go/zap` configuration and multi-output examples |
| [log-masking](./log-masking) | Masking sensitive fields in log output |

### Datastores

| Directory | Description |
| --- | --- |
| [bigtable](./bigtable) | Cloud Bigtable client with App Profile usage |
| [preparedstatement](./preparedstatement) | PostgreSQL prepared-statement performance comparison |
| [nested-mongo-document](./nested-mongo-document) | Nested MongoDB documents (including BSON flattening) |
| [causal-consistency](./causal-consistency) | MongoDB causal consistency (read-your-writes) |
| [distributed-lock](./distributed-lock) | Distributed locking with etcd |
| [firebase](./firebase) | Firebase backend with a Vue.js front end |

### Security / Authorization / Feature Flags

| Directory | Description |
| --- | --- |
| [crypto](./crypto) | AES-GCM encryption and Argon2 hashing |
| [open-policy-agent](./open-policy-agent) | Policy evaluation with OPA |
| [openfga](./openfga) | ReBAC authorization with OpenFGA |
| [openfeature](./openfeature) | Feature-flag management with OpenFeature |
| [pgv](./pgv) | Protocol Buffers validation with `protoc-gen-validate` |

### Algorithms / Data Structures

| Directory | Description |
| --- | --- |
| [bandit](./bandit) | Multi-armed bandit (Thompson sampling) |
| [bloomfilter](./bloomfilter) | Bloom filter |
| [radix-tree](./radix-tree) | Radix tree vs. regex path-matching benchmark |
| [weighted-random-choice](./weighted-random-choice) | Weighted random selection |
| [guid](./guid) | UUID vs. Sonyflake vs. xid generation benchmark |

### DI / Plugins

| Directory | Description |
| --- | --- |
| [wire](./wire) | Dependency injection with Google Wire (struct / interface / packages variants) |
| [go-plugin-rpc](./go-plugin-rpc) | RPC-based plugins via `hashicorp/go-plugin` |

## Usage

Each directory has its own entry point and can typically be run with
`go run` or `go test` from inside it.

```bash
cd radix-tree
go test -bench . -benchmem
```

Samples that depend on external middleware (PostgreSQL, MongoDB, etcd,
Bigtable, etc.) carry their own README and/or `docker-compose.yml` in the
subdirectory.

## License

[MIT](./LICENSE)
