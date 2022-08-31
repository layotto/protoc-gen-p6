# protoc-gen-p6
<img src="https://user-images.githubusercontent.com/26001097/187589228-23f04889-6c0e-41d9-abbb-89057a6d777d.png" width="20%" height="20%">

p6 is a protoc plugin. 

She is good at coding and writes code for [Layotto](https://github.com/mosn/layotto) project. 

p6 works hard.

## Install

Please make sure you have tools below:

- [go 1.16](https://golang.org/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf)
- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go)

Go version >= 1.16

```shell
go install github.com/seeflood/protoc-gen-p6@latest
```

## How to give p6 new work
Suppose you have a new proto file [example/api/product/app/v1/blog.proto](example/api/product/app/v1/blog.proto) , and you want to implement this API in [Layotto](https://github.com/mosn/layotto) . 

It's a tedious job because you have to write lots of boring code. You don't want to do it yourself.

Then you can ask p6 to do it. For example:

```shell
protoc -I ./example/api \
      --go_out ./example/api --go_opt=paths=source_relative \
      --p6_out ./example/api --p6_opt=paths=source_relative \
      example/api/product/app/v1/blog.proto
```

or, you can give her an example work by:

```shell
make work.example
```

And p6 will write the code very quickly:

![image](https://user-images.githubusercontent.com/26001097/187620570-d4666c28-8390-437d-8e9c-9ab17064e508.png)

You can then move these code to Layotto project.