go.install:
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2

gorm.gen:
	go run ./cmd/gorm-gen

buf.gen:
	buf generate ./...
