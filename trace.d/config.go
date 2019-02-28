package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
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
	MinioUrl    string `yaml:"minio-url"`
	Consul      string `yaml:"consul"`
	PgUser      string `yaml:"pg-user"`
	PgPass      string `yaml:"pg-pass"`
	PgHost      string `yaml:"pg-host"`
}

var config *configuration

const grpcAddress = "localhost:50051"

func (c *configuration) getConf() error {
	yamlFile, err := ioutil.ReadFile("./trace.d/conf.yaml")
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return err
	}
	return nil
}
