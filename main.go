package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/table"
)

func main() {
	// socket := flag.String("socket", "", "Path to osquery socket file")
	socket := "/Users/brettdemetris/.osquery/shell.em"
	// flag.Parse()
	// if *socket == "" {
	// 	log.Fatalf(`Usage: %s --socket SOCKET_PATH`, os.Args[0])
	// }

	server, err := osquery.NewExtensionManagerServer("crowdstrike_falcon", socket)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// table.NewPlugin requires the table plugin name,
	// a slice of Columns and a Generate function.
	server.RegisterPlugin(table.NewPlugin("crowdstrike_falcon", falconColums(), falconGenerate))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

func falconGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	stats, err := falconStats()
	if err != nil {
		log.Println(err)
	}
	agent, err := Stats.agentInfo(stats)
	if err != nil {
		log.Println(err)
	}
	cloud, err := Stats.cloudInfo(stats)
	if err != nil {
		log.Println(err)
	}

	ft := falconTable{
		Version:           agent.Version,
		AgentID:           agent.AgentID,
		CustomerID:        agent.CustomerID,
		SensorOperational: agent.SensorOperational,
		Host:              cloud.Host,
		Port:              cloud.Port,
		State:             cloud.State,
	}

	var values []map[string]string

	j, _ := json.Marshal(ft)
	m := make(map[string]string)
	err = json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}

	values = append(values, m)

	return values, nil
}

func falconColums() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("Version"),
		table.TextColumn("AgentID"),
		table.TextColumn("CustomerID"),
		table.TextColumn("SensorOperational"),
		table.TextColumn("Host"),
		table.TextColumn("Port"),
		table.TextColumn("State"),
	}
}

type falconTable struct {
	Version           string
	AgentID           string
	CustomerID        string
	SensorOperational string
	Host              string
	Port              string
	State             string
}
