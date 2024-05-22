package loadutil

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
	} `yaml:"redis"`
	RpcUrl    string `yaml:"rpcUrl"`
	GroupSize int    `yaml:"groupSize"`
	LocalPath string `yaml:"datapath"`
}

func LoadutilFactory(config Config) (Loadutil, error) {
	util, err := OSSUtilFactory(
		config.OSS.Endpoint,
		config.OSS.AccessKeyID,
		config.OSS.AccessKeySecret,
		config.OSS.BucketName,
	)
	if err != nil {
		return nil, err
	}
	return util, nil
}

type Loadutil interface {
	LoadToFile(string, string) error
}
