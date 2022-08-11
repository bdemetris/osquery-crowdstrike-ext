package main

import (
	"flag"
	"log"
	"time"

	"github.com/bdemetris/osquery-crowdstrike-ext/tables/crowdstrike"
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
	server.RegisterPlugin(table.NewPlugin("crowdstrike_falcon", crowdstrike.FalconColumns(), crowdstrike.FalconGenerate))
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
