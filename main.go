package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	osquery "github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

func main() {
	var (
		flSocketPath = flag.String("socket", "", "")
		flTimeout    = flag.Int("timeout", 0, "")
		_            = flag.Int("interval", 0, "")
		_            = flag.Bool("verbose", false, "")
	)
	flag.Parse()

	// allow for osqueryd to create the socket path otherwise it will error
	time.Sleep(3 * time.Second)

	server, err := osquery.NewExtensionManagerServer(
		"crowdstrike_falcon",
		*flSocketPath,
		osquery.ServerTimeout(time.Duration(*flTimeout)*time.Second))
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	// Create and register a new table plugin with the server.
	// table.NewPlugin requires the table plugin name,
	// a slice of Columns and a Generate function.
	server.RegisterPlugin(table.NewPlugin("crowdstrike_falcon", falconColums(), FalconGenerate))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

func FalconGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	stats, err := falconStats()
	if err != nil {
		log.Println(err)
	}

	agent := stats.agentInfo()
	cloud := stats.cloudInfo()

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
