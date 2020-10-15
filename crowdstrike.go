package main

import (
	"log"
	"os/exec"

	"gopkg.in/ini.v1"
)

const falconPath = "/Applications/Falcon.app/Contents/Resources/falconctl"

// Stats is a reusable struct of the CFS stats output
type Stats struct {
	falconStats *ini.File
}

func falconStats() (Stats, error) {

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

func (s Stats) agentInfo() (agentInfo, error) {
	// op, err := strconv.ParseBool(s.falconStats.Section("").Key("Sensor operational").String())
	// if err != nil {
	// 	log.Println(err)
	// }
	info := agentInfo{
		Version:           s.falconStats.Section("").Key("version").String(),
		AgentID:           s.falconStats.Section("").Key("agentID").String(),
		CustomerID:        s.falconStats.Section("").Key("customerID").String(),
		SensorOperational: s.falconStats.Section("").Key("Sensor operational").String(),
	}
	return info, nil
}

func (s Stats) cloudInfo() (cloudInfo, error) {
	// p, err := strconv.Atoi(s.falconStats.Section("").Key("Port").String())
	// if err != nil {
	// 	log.Println(err)
	// }
	info := cloudInfo{
		Host:  s.falconStats.Section("").Key("Host").String(),
		Port:  s.falconStats.Section("").Key("Port").String(),
		State: s.falconStats.Section("").Key("State").String(),
	}

	return info, nil
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
