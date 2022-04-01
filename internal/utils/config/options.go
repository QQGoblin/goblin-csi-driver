package config

import (
	"goblin-csi-driver/internal"
	"goblin-csi-driver/internal/utils/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Options struct {
	HostPathConfig *HostPathConfig `yaml:"hostPathConfig"`
	PosConfig      *PosConfig      `yaml:"posConfig"`
}

type HostPathConfig struct {
	State string `yaml:"state"`
}

type PosConfig struct {
}

func NewOptions() *Options {
	return &Options{
		HostPathConfig: &HostPathConfig{State: internal.DefaultHostPathState},
	}
}

func LoadOptions(path string) *Options {

	cf, err := os.Open(path)
	if err != nil {
		log.FatalLogMsg("load config file from %s failed, err %v", path, err)
	}
	defer cf.Close()

	bs, err := ioutil.ReadAll(cf)
	if err != nil {
		log.FatalLogMsg("read config file from %s failed, err %v", path, err)
	}

	opt := &Options{}

	if err := yaml.Unmarshal(bs, opt); err != nil {
		log.FatalLogMsg("unmarshal config failed, content: %v，", bs)
	}

	return opt
}
