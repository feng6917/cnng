package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var Conf = &YamlFile{}

type YamlFile struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Port    int    `yaml:"port"`
	} `yaml:"app"`
}

func UnmarshalYamlFile(filePath string) error {
	ymlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败；err: %v", err)
		return err
	}
	err = yaml.Unmarshal(ymlFile, &Conf)
	if err != nil {
		log.Fatalf("解析文件失败；err: %v", err)
		return  err
	}
	return  err
}
