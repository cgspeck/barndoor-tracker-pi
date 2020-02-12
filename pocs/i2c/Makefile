build:
	GOOS=linux GOARCH=arm go build ./cmd/accel
	@chmod +x ./accel
	GOOS=linux GOARCH=arm go build ./cmd/basic
	@chmod +x ./basic
	GOOS=linux GOARCH=arm go build ./cmd/mag
	@chmod +x ./mag
	GOOS=linux GOARCH=arm go build ./cmd/whichway
	@chmod +x ./whichway

deps:
	go get "github.com/davecgh/go-spew/spew"
	go get "github.com/kidoman/embd"
