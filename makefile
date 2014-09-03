all: build

test: build unittests integration

build: server.go
	go build -o doccer server.go

clean:
	rm doccer

rmproject:
	rm accounts; rm content

integration: 
	nosetests

integration-xunit:
	nosetests --with-xunit --xunit-file=integration.xml

unittests:
	go test

jenkins: build unittests integration-xunit rmproject
