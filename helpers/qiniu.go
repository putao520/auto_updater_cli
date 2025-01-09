package helpers

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"path/filepath"
	"time"
)

func PutWithRewrite(distributeOptions map[string]interface{}, spec map[string]interface{}, releaseName string, jobName string, localFile string) error {
	vars, _ := distributeOptions["variables"].(map[string]interface{})
	ak, _ := vars["QINIU_ACCESS_KEY"].(string)
	sk, _ := vars["QINIU_SECRET_KEY"].(string)
	mac := credentials.NewCredentials(ak, sk)

	pub := GetReleaseItemInfo(distributeOptions, releaseName, jobName, "publish")
	target, _ := pub["target"].(string)
	if target != "qiniu" {
		panic("publish target must be 'qiniu'")
	}
	pubArgs, _ := pub["args"].(map[string]interface{})
	bucket, _ := pubArgs["bucket"].(string)
	// bucketDomain := pubArgs["bucket-domain"].(string)
	keyPrefix, _ := pubArgs["savekey-prefix"].(string)

	localFile, _ = filepath.Abs(localFile)
	fileName := filepath.Base(localFile)
	keyToOverwrite := fmt.Sprintf("%s%s", keyPrefix, fileName)

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	putPolicy, err := uptoken.NewPutPolicyWithKey(bucket, keyToOverwrite, time.Now().Add(1*time.Hour))
	if err != nil {
		return err
	}
	return uploadManager.UploadFile(context.Background(), localFile, &uploader.ObjectOptions{
		UpToken:    uptoken.NewSigner(putPolicy, mac),
		ObjectName: &keyToOverwrite,
		FileName:   fileName,
		CustomVars: map[string]string{},
	}, nil)
}
