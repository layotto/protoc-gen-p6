gen.example:
	go install
	  protoc -I ./example/api \
      --go_out ./example/api --go_opt=paths=source_relative \
      --p6_out ./example/api --p6_opt=paths=source_relative \
      example/api/product/app/v1/blog.proto