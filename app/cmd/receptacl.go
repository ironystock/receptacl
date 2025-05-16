package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const config_file = "cfg/receptacl.yaml"

type receptaclCfg struct {
	ServiceName  string `yaml:"servicename"`
	RemoteServer struct {
		IP   string `yaml:"ip"`
		Port string `yaml:"port"`
		PW   string `yaml:"pw,omitempty"`
	} `yaml:"remoteserver"`
	LocalPath string `yaml:"localpath"`
	LocalPort string `yaml:"localport"`
	Timeout   int
	MsgBuffer int
}

type hxMesg struct {
}

func configure() (receptaclCfg, error) {
	var Cfg receptaclCfg
	f, err := os.Open(config_file)
	if err != nil {
		return receptaclCfg{}, err
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	err = d.Decode(&Cfg)
	if err != nil {
		return receptaclCfg{}, err
	}

	return Cfg, nil

}
