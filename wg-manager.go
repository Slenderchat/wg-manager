package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"wg-manager/client"
	"wg-manager/genConfig"
)

var clients []*client.Client
var action = "Status"
var name = ""

func saveClients() error {
	f, e := os.OpenFile("clients.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if e != nil {
		return fmt.Errorf("failed to open clients database: %v", e)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	e = enc.Encode(clients)
	if e != nil {
		return fmt.Errorf("failed to update clients database: %v", e)
	}
	return nil
}

func readClients() error {
	f, e := os.OpenFile("clients.json", os.O_CREATE|os.O_RDONLY, 0600)
	if e != nil {
		return fmt.Errorf("failed to open clients database: %v", e)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	e = dec.Decode(&clients)
	if e != nil {
		return fmt.Errorf("failed to read clients database: %v", e)
	}
	for _, v := range clients {
		if v.Name == "" {
			return errors.New("unnamed clients in clients database")
		} else if v.PrivateKey == "" && v.PublicKey == "" {
			removeClient(v.Name)
			client.NewClient(v.Name)
			v = nil
		} else {
			if v.PrivateKey == "" {
				return fmt.Errorf("client %v does not have private key", v.Name)
			}
			if v.PublicKey == "" {
				return fmt.Errorf("client %v does not have public key", v.Name)
			}
			if v.IP == "" {
				fmt.Fprintf(os.Stderr, "Client %v does not have assigned IP.\n"+
					"You can either recreate client (Keys will change) or delete it from database.\n"+
					"To recreate type `r`, to delete type `d`, then press `Enter`: ", v.Name)
				for {
					var a string
					_, e := fmt.Scanln(&a)
					if e != nil {
						return fmt.Errorf("failed to get answer: %v", e)
					}
					if a == `r` {
						removeClient(v.Name)
						newClient(v.Name)
						v = nil
						break
					} else if a == `d` {
						removeClient(v.Name)
						v = nil
						break
					} else {
						fmt.Fprintf(os.Stderr, "Invalid answer `%v`.\n"+
							"Please type exactly (without quotes) `r` to recreate or `d` to delete, then press `Enter`: ", a)
						continue
					}
				}
			}
		}
	}
	return nil
}

func readArgs() error {
	var isNextValue bool = false
	var arg [2]string
	for _, v := range os.Args[1:] {
		if isNextValue {
			arg[1] = v
		} else {
			if v[0:2] != "--" {
				return errors.New("all arguments should start with --")
			} else if strings.Contains(v, "=") {
				arg = ([2]string)(strings.Split(v, "="))
			} else {
				arg[0] = v
				arg[1] = ""
			}
		}
		switch arg[0] {
		case "--newclient":
			action = "New client"
		case "--delclient":
			action = "Delete client"
		case "--genserverconfig":
			action = "Generate server config"
		case "--genclientconfig":
			action = "Generate client config"
		case "--name":
			if isNextValue {
				name = arg[1]
				isNextValue = false
			} else {
				isNextValue = true
				continue
			}
		default:
			return fmt.Errorf("unknown argument %v", arg[0])
		}
	}
	if isNextValue {
		return fmt.Errorf("argument %v is missing value", arg[0])
	}
	return nil
}

func newClient(name string) error {
	if name == "" {
		return errors.New("please use --name to specify new client's name")
	}
	for _, v := range clients {
		if v.Name == name {
			return errors.New("client with desired name already exists")
		}
	}
	newClientObject := client.NewClient(name)
	clients = append(clients, newClientObject)
	e := saveClients()
	if e != nil {
		return e
	}
	fmt.Printf("Added new client %v\nPrivate key: %v\nPublic key: %v\n", newClientObject.Name, newClientObject.PrivateKey, newClientObject.PublicKey)
	return nil
}

func removeClient(name string) error {
	if name == "" {
		return errors.New("please use --name to specify which client to delete")
	}
	var clientToDelete int = -1
	var clientToDeleteName string
	for i, v := range clients {
		if v.Name == name {
			clientToDelete = i
			clientToDeleteName = v.Name
		}
	}
	if clientToDelete < 0 || clientToDelete > len(clients) || clientToDeleteName == "" {
		return errors.New("client with specified name not found")
	}
	clients[clientToDelete] = clients[len(clients)-1]
	clients = clients[:len(clients)-1]
	e := saveClients()
	if e != nil {
		return e
	}
	fmt.Printf("Removed client %v\n", clientToDeleteName)
	return nil
}

func main() {
	e := readArgs()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", e)
		os.Exit(1)
	}
	e = readClients()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error reading clients database: %v\n", e)
		os.Exit(2)
	}
	if action == "Status" {
		e := status()
		if e != nil {
			fmt.Fprintf(os.Stderr, "Failed to determine status: %v\n", e)
			os.Exit(3)
		}
	}
	if action == "New client" {
		e := newClient(name)
		if e != nil {
			fmt.Fprintf(os.Stderr, "Failed to create new client: %v\n", e)
			os.Exit(4)
		}
	}
	if action == "Delete client" {
		e := removeClient(name)
		if e != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove client: %v\n", e)
			os.Exit(5)
		}
	}
	if action == "Generate server config" {
		config, e := genConfig.GenServerConfig(&clients)
		if e != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate server config: %v\n", e)
			os.Exit(6)
		}
		fmt.Println(config)
	}
	if action == "Generate client config" {
		if name == "" {
			fmt.Fprintln(os.Stderr, "Please use --name to specify new client's name for which config will be generated")
		}
		for _, v := range clients {
			if v.Name == name {
				config, e := genConfig.GenClientConfig(v)
				if e != nil {
					fmt.Fprintf(os.Stderr, "Failed to generate config for client %v: %v\n", v.Name, e)
					os.Exit(7)
				}
				fmt.Println(config)
			}
		}
	}
}
