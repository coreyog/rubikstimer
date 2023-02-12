package config

import (
	"encoding/json"
	"math"
	"os"
	"strings"
)

// Trigger is used for indicating how to start and stop the timer
type Trigger string

const (
	TriggerModifiers Trigger = "modifiers"
	TriggerSpacebar  Trigger = "spacebar"
	TriggerAny       Trigger = "any"
)

type Config struct {
	ScrambleLength    int    `json:"scrambleLength"`
	TimerStartTrigger string `json:"timerStartTrigger"`
	TimerEndTrigger   string `json:"timerEndTrigger"`
}

var globalConfig = defaultConfig()

// LoadConfig reads and parses a config file
func LoadConfig() {
	config := defaultConfig()

	raw, err := os.ReadFile("config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return
	}

	config.ScrambleLength = int(math.Max(10, math.Min(float64(config.ScrambleLength), 50)))

	resave1 := checkTrigger(&config.TimerStartTrigger)

	resave2 := checkTrigger(&config.TimerEndTrigger)

	if resave1 || resave2 {
		SaveConfig(config)
	}

	globalConfig = config
}

// SaveConfig writes config to file and sets it as the global config for immediate use
func SaveConfig(config *Config) {
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
func GlobalConfig() *Config {
	return globalConfig
}

func defaultConfig() (config *Config) {
	return &Config{
		ScrambleLength:    20,
		TimerStartTrigger: string(TriggerModifiers),
		TimerEndTrigger:   string(TriggerModifiers),
	}
}

func checkTrigger(t *string) (resave bool) {
	*t = strings.ToLower(*t)
	if *t == "controls" {
		// temporary fix for loading old configs correctly
		*t = string(TriggerModifiers)
		resave = true
	}

	// if invalid, revert to default values
	if *t == string(TriggerAny) || *t == string(TriggerModifiers) || *t == string(TriggerSpacebar) {
		*t = string(TriggerModifiers)
	}

	return resave
}
