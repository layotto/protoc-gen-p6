# protoc-gen-p6
<img src="https://user-images.githubusercontent.com/26001097/188037681-005b8104-823e-45ea-82a9-3f77cd371636.png" width="30%" height="30%">

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
      --go-grpc_out=./example/api \
      --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative \
      --p6_out ./example/api --p6_opt=paths=source_relative \
      example/api/product/app/v1/blog.proto
```

or, you can give her an example work by:

```shell
make work.example
```

And p6 will write the code very quickly:

![image](https://user-images.githubusercontent.com/26001097/188781256-690b6d47-3d5a-4f09-9dcf-e9dda3ae151f.png)

You can then move these code to Layotto project.