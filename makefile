all:
	GOPATH=`pwd`:`pwd`/vendor go build -o bin/fixme fixme
