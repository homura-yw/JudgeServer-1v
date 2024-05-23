package loadutil

import "fmt"

type Config struct {
	OSS struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyID"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		BucketName      string `yaml:"bucketName"`
	} `yaml:"OSS"`
	Redis struct {
		Url      string `yaml:"url"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
	RpcUrl    string `yaml:"rpcUrl"`
	GroupSize int    `yaml:"groupSize"`
	LocalPath string `yaml:"datapath"`
	Source    string `yaml:"source"`
}

func LoadutilFactory(config Config) (Loadutil, error) {
	var util Loadutil
	var err error
	if config.Source == "OSS" {
		util, err = OSSUtilFactory(
			config.OSS.Endpoint,
			config.OSS.AccessKeyID,
			config.OSS.AccessKeySecret,
			config.OSS.BucketName,
		)
	} else if config.Source == "local" {
		util, err = LocalDataUtilFactory(config.LocalPath)
	} else {
		err = fmt.Errorf("failed to load source:%s", config.Source)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return util, nil
}

type Loadutil interface {
	LoadToFile(string, string) error
}
