
gonic1:  main_server.go
	go build .

upd:
	go get -u -v
	go mod tidy

clean:
	rm gonic1
