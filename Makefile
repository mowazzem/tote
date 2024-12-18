GO_MODULE= "github.com/mowazzem/tote"

.PHONY: protoc
protoc-go:
	protoc --go_opt=module=$(GO_MODULE) --go_out=. \
	--go-grpc_opt=module=$(GO_MODULE) --go-grpc_out=. \
	./proto/common/*.proto \
	./proto/auth/*.proto

.PHONY: clean
clean:
	rm -rf protogen