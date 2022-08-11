package crowdstrike

import (
	"context"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/osquery/osquery-go/plugin/table"
	"gopkg.in/ini.v1"
)

const falconPath = "/Applications/Falcon.app/Contents/Resources/falconctl"

// Stats is a reusable struct of the CFS stats output
type Stats struct {
	falconStats *ini.File
}

func FalconStats() (Stats, error) {

	out, err := exec.Command(falconPath, "stats").Output()
	if err != nil {
		log.Println(err)
	}

	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
		KeyValueDelimiters:      ":",
	}, out)
	if err != nil {
		log.Println(err)
	}

	stats := Stats{
		falconStats: cfg,
	}

	return stats, nil
}

func (s Stats) AgentInfo() agentInfo {
	return agentInfo{
		Version:           s.falconStats.Section("").Key("version").String(),
		AgentID:           s.falconStats.Section("").Key("agentID").String(),
		CustomerID:        s.falconStats.Section("").Key("customerID").String(),
		SensorOperational: s.falconStats.Section("").Key("Sensor operational").String(),
	}
}

func (s Stats) CloudInfo() cloudInfo {
	return cloudInfo{
		Host:  s.falconStats.Section("").Key("Host").String(),
		Port:  s.falconStats.Section("").Key("Port").String(),
		State: s.falconStats.Section("").Key("State").String(),
	}
}

func FalconGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	stats, err := FalconStats()
	if err != nil {
		log.Println(err)
	}

	agent := stats.AgentInfo()
	cloud := stats.CloudInfo()

	ft := FalconTable{
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

func FalconColumns() []table.ColumnDefinition {
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

type agentInfo struct {
	Version           string
	AgentID           string
	CustomerID        string
	SensorOperational string
}

type cloudInfo struct {
	Host  string
	Port  string
	State string
}

type FalconTable struct {
	Version           string
	AgentID           string
	CustomerID        string
	SensorOperational string
	Host              string
	Port              string
	State             string
}
