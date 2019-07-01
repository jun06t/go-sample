# go-plugin-rpc

# Quick start
## Build host
```
go build -o host main.go
```

## No authentication Plugin
### Build plugin
```
cd plugin
go build -o auth ./no-auth/plugin.go
```
### Run host
```
cd ..
./host
```
Then it returns,
```
success!!
```

## Password authentication Plugin
### Build plugin
```
cd plugin
go build -o auth ./password/plugin.go
```
### Run host
```
cd ..
./host
```
Then it returns,
```
password:
```
If you input the valid password ``hoge``.
```
password:hoge
success!!
```
If you input the invalid password ``fuga``.
```
password:fuga
fail!!
```