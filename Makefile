build:
	go build -ldflags "-s -w" wg-manager.go status.go

clean:
	rm wg-manager servers.json clients.json