package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

var config *Config

type Config struct {
	Listen string `yaml:"listen"`
	MySql string `yaml:"mysql"`
	Assets string `yaml:"assets"`
	ImagePath string `yaml:"imgpath"`
	CpuProf string `yaml:"cpuprofile"`
	MemProf string `yaml:"memprofile"`
	ElasticUrl string `yaml:"elastic-url"`
	ElasticLog string `yaml:"elastic-log"`
	ElasticPass string `yaml:"elastic-pass"`
	Mongo string `yaml:"mongo"`
}

func (c *Config) getConf() error {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	return nil
}