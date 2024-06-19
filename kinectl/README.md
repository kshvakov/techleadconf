# kinectl

## структура проекта

```
.
├── cmd
│   ├── server
│   │   ├── action
│   │   │   ├── a.go
│   │   │   └── b.go
│   │   └── main.go
│   └── task
├── config
├── docs
├── internal
│   ├── handler
│   ├── middleware
│   ├── module
│   ├── server
│   └── service
├── pkg
│   ├── lib1
│   └── lib2
├── test
│   ├── mocks
│   └── server
├── Makefile
├── version.go
├── spec.yml
├── go.mod
└── go.sum
```

## переменные окружения

Лучше давать однотипные и говорящие о назначении названия, например для адресов:

* HTTP_ADDR
* GRPC_ADDR
* PROXY_ADDR

