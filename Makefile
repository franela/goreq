export GOPATH=$(shell pwd)
test:
	go get github.com/onsi/gomega
	go get github.com/franela/goblin
	go test -v
