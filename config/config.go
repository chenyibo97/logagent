package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	KafkaAddr string `json:"kafkaServerAddr"`
	LogFileName string `json:"logFileName"`
	Test int `json:"test"`
	EtcdAddr string `json:"etcdAddr"`
	EtcdTimeout int `json:"etcdTimeout"`
	EtcdAgentKey string `json:"etcdAgentKey"`
}
func ReadConfigFile() (*Config,error){
	LogagentConfig:=new(Config)
    file,err:=os.Open("./config/config.json")
	if err != nil {
		return nil,err
	}

	defer file.Close()
    decoder:=json.NewDecoder(file)
    err=decoder.Decode(LogagentConfig)
	if err != nil {
		fmt.Println("decode err:",err)
		return nil, err
	}
	fmt.Println(LogagentConfig)
	return LogagentConfig,nil
}