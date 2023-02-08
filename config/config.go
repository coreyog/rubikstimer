package config

import (
	"encoding/json"
	"os"
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
	raw, err := os.ReadFile("config.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return
	}

	if config.ScrambleLength < 10 {
		config.ScrambleLength = 10
	} else if config.ScrambleLength > 50 {
		config.ScrambleLength = 50
	}
	if !checkTrigger(config.TimerStartTrigger) {
		config.TimerStartTrigger = defaultConfig().TimerStartTrigger
	}
	if !checkTrigger(config.TimerEndTrigger) {
		config.TimerEndTrigger = defaultConfig().TimerEndTrigger
	}

	globalConfig = config
}

// SaveConfig writes config to file and sets it as the global config for immediate use
func SaveConfig(config Config) {
	raw, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("config.json", raw, 0666)
	if err != nil {
		panic(err)
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
