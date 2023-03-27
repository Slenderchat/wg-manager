package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func status() error {
	e := saveClients()
	if e != nil {
		return fmt.Errorf("error trying to save status: %v", e)
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\t\tPrivate key\t\tPublic Key\t\tIP\t\tEndpoint\t\tRoute table\t\tRoute metric")
	fmt.Fprintln(w, "----\t\t-----------\t\t----------\t\t--\t\t--------\t\t-----------\t\t-----------")
	for _, v := range clients {
		fmt.Fprintf(w, "%v\t\t%v\t\t%v\t\t%v\t\t%v\t\t%v\t\t%v\n", v.Name, v.PrivateKey, v.PublicKey, v.IP, v.Endpoint, v.RouteTable, v.RouteMetric)
	}
	w.Flush()
	fmt.Println()
	return nil
}
