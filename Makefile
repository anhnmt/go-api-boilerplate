go.install:
	go install github.com/bufbuild/buf/cmd/buf@v1.32.2
	go install github.com/google/wire/cmd/wire@v0.6.0
	go install golang.org/x/tools/cmd/deadcode@v0.22.0

gorm.gen:
	go run ./cmd/gorm-gen

buf.gen:
	buf generate ./...

wire.gen:
	wire ./...

deadcode:
	deadcode ./...