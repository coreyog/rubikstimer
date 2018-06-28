package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Trigger is used for indicating how to start and stop the timer
type Trigger string

// Trigger strings
const (
	TriggerControls Trigger = "controls"
	TriggerSpacebar Trigger = "spacebar"
	TriggerAny      Trigger = "any"
)

// Config holds options that can be set
type Config struct {
	ScrambleLength    int    `json:"scrambleLength"`
	TimerStartTrigger string `json:"timerStartTrigger"`
	TimerEndTrigger   string `json:"timerEndTrigger"`
}

var globalConfig = defaultConfig()

// LoadConfig reads and parses a config file
func LoadConfig() {
	var config Config
	globalConfig = defaultConfig()
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return
	}

	if config.ScrambleLength < 10 || config.ScrambleLength > 30 {
		config.ScrambleLength = defaultConfig().ScrambleLength
	}
	if !checkTrigger(config.TimerStartTrigger) {
		config.TimerStartTrigger = defaultConfig().TimerStartTrigger
	}
	if !checkTrigger(config.TimerEndTrigger) {
		config.TimerEndTrigger = defaultConfig().TimerEndTrigger
	}

	globalConfig = config
}

// GlobalConfig gives access to a
func GlobalConfig() Config {
	return globalConfig
}

func defaultConfig() (config Config) {
	config.ScrambleLength = 20
	config.TimerStartTrigger = string(TriggerControls)
	config.TimerEndTrigger = string(TriggerControls)
	return config
}

func checkTrigger(t string) (isAllowed bool) {
	t = strings.ToLower(t)
	return t == string(TriggerAny) || t == string(TriggerControls) || t == string(TriggerSpacebar)
}
