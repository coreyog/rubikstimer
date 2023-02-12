package config

import (
	"encoding/json"
	"os"

	"gopkg.in/validator.v2"
)

// Trigger is used for indicating how to start and stop the timer
type Trigger string

const (
	TriggerModifiers Trigger = "modifiers"
	TriggerSpacebar  Trigger = "spacebar"
	TriggerAny       Trigger = "any"
)

type Config struct {
	ScrambleLength    int    `json:"scrambleLength" validate:"min=10,max=50"`
	TimerStartTrigger string `json:"timerStartTrigger" validate:"oneof=modifiers spacebar any"`
	TimerEndTrigger   string `json:"timerEndTrigger" validate:"oneof=modifiers spacebar any"`
	WindowWidth       int    `json:"windowWidth" validate:"min=100,max=10000"`
	WindowHeight      int    `json:"windowHeight" validate:"min=100,max=10000"`
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

	resave := false

	err = validator.Validate(config)
	if err != nil {
		// something in the config isn't kosher, so we'll just use the defaults
		config = defaultConfig()
		resave = true
	}

	if resave {
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
		WindowWidth:       1000,
		WindowHeight:      400,
	}
}
