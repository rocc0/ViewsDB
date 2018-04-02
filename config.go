package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var config *configuration

const (
	grpcAddress = "localhost:50051"
)

type configuration struct {
	Listen      string `yaml:"listen"`
	MySQL       string `yaml:"mysql"`
	Assets      string `yaml:"assets"`
	ImagePath   string `yaml:"imgpath"`
	CPUProf     string `yaml:"cpuprofile"`
	MemProf     string `yaml:"memprofile"`
	ElasticURL  string `yaml:"elastic-url"`
	ElasticLog  string `yaml:"elastic-log"`
	ElasticPass string `yaml:"elastic-pass"`
	Mongo       string `yaml:"mongo"`
	MinioKay    string `yaml:"minio-k"`
	MinioSecret string `yaml:"minio-s"`
}

func (c *configuration) getConf() error {
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
