# OPA

## Get Started

### Run OPA server

```
$ opa run -s policy.rego data.json
```

### Run API server

```
$ go run main.go
```

### Call API

#### GET /articles

```
$ curl -i -XGET localhost:3000/articles
HTTP/1.1 200 OK
Date: Sun, 12 Sep 2021 22:26:47 GMT
Content-Length: 7
Content-Type: text/plain; charset=utf-8

listed
```

#### GET /articles/foo

```
$ curl -i -XGET localhost:3000/articles/foo
HTTP/1.1 200 OK
Date: Mon, 13 Sep 2021 02:31:12 GMT
Content-Length: 8
Content-Type: text/plain; charset=utf-8

got foo
```

#### DELETE /articles/foo

```
$ curl -i -XDELETE localhost:3000/articles/hoge
HTTP/1.1 403 Forbidden
Date: Mon, 13 Sep 2021 02:31:32 GMT
Content-Length: 0
```
