package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config *configuration

type configuration struct {
	ImagePath   string `yaml:"imgpath"`
	Mongo       string `yaml:"mongo"`
	MinioKay    string `yaml:"minio-k"`
	MinioSecret string `yaml:"minio-s"`
	MinioUrl    string `yaml:"minio-url"`
	Consul      string `yaml:"consul"`
}

func (c *configuration) getConf() error {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	return nil
}
