all: build

test: build integration


build: server.go
	go build -o doccer server.go

clean:
	rm doccer

integration: 
	nosetests

