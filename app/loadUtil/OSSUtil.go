package loadutil

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSUtil struct {
	bucket *oss.Bucket
}

func (ou *OSSUtil) LoadToFile(url, path string) error {

	bts := []byte(url)
	if bts[0] == '/' {
		bts = bts[1:]
		url = string(bts)
	}

	err := ou.bucket.GetObjectToFile(url, path)
	return err
}

func OSSUtilFactory(Endpoint, AccessKeyID, AccessKeySecret, BucketName string) (*OSSUtil, error) {
	ossClient, err := oss.New(
		Endpoint,
		AccessKeyID,
		AccessKeySecret,
	)
	if err != nil {
		return nil, err
	}
	bucket, err := ossClient.Bucket(BucketName)
	if err != nil {
		return nil, err
	}
	return &OSSUtil{bucket: bucket}, nil
}
