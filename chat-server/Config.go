package main

import (
	"fmt"
	"github.com/larspensjo/config"
	"log"
)

const TOPICSECTION = "topicList"
const WINDOWS = "windows"
const LINUX = "linux"


var (
	//topic list
	TOPIC = make(map[string]string)
)

type ConfigHelper struct {
}

func (this *ConfigHelper) LoadConfig(configFile *string) {
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find config file", *configFile)
	}

	//Initialized topic from the configuration
	if cfg.HasSection(TOPICSECTION) {
		section, err := cfg.SectionOptions(TOPICSECTION)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(TOPICSECTION, v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	fmt.Println(TOPIC)
}
