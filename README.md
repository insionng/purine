#purine

A Fast and Powerful Blog Framework built in Go.

### Features

* Support markdown content
* Single binary file

### Installation

If you want to use `purine` as your blog, simply install the `purine` binaries without dependencies.

[Binaries Download](#)

##### Development

1.install requirements:

    - `Go`
    - `Git`
    - `gcc` for compiling SQLite

2.download source code:

```
go get -u -v github.com/fuxiaohei/purine
```

3.run `purine`:

```
cd $GOPATH/src/github.com/fuxiaohei/purine
go run purine.go
```

### Started

If you get the binary release `purine.exe`, install and run blog framework with following steps:

1.install blog, create default config and database

```
purine.exe install
```

2.start http server

```
purine.exe server
```

the blog is running on http://localhost:9999.

3.enter into administrator

visit http://localhost:9999/admin/login with **admin:123456789**