```sh
$ docker run --rm -v "$(pwd)":/app -w /app golang:1.21 go run main.go
Go Version: go1.21.13
```

```
$ docker run --rm -v "$(pwd)":/app -w /app golang:1.22 go run main.go
Go Version: go1.22.9
```

```
$ docker run --rm -v "$(pwd)":/app -w /app -e GOTOOLCHAIN=auto golang:1.21 go run main.go
Go Version: go1.22.1
```

```
$ docker run --rm -v "$(pwd)":/app -w /app -e GOTOOLCHAIN=go1.21.3 golang:1.22 go run main.go
Go Version: go1.21.3
```


## go directive 1.22
```
$ docker run --rm -v "$(pwd)":/app -w /app -e GOTOOLCHAIN=go1.21.3 golang:1.22 go run main.go
go: downloading go1.21.3 (linux/arm64)
go: go.mod requires go >= 1.22 (running go 1.21.3; GOTOOLCHAIN=go1.21.3)
```