package services

import (
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"hash"
	"io"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	
	"{{ . }}/internal/models"
)

type MediaFileService interface {
	// Upload(c context.Context, file *models.MediaFile) error
	GetSign(c context.Context) (*models.AliPolicyToken, error)
	UploadCallback(c context.Context, bytePublicKey []byte, byteMd5 []byte, authorization []byte) bool
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

type DefaultMediaFileService struct {
	logger                *zap.Logger
	aliyunAccessKeyId     string
	aliyunAccessKeySecret string
	aliyunEndpoint        string
	aliyunCallbackUrl     string
	aliyunUploadDir       string
	aliyunExpireTime      int64
}

func NewMediaFileService(v *viper.Viper, logger *zap.Logger) MediaFileService {
	return &DefaultMediaFileService{
		logger:                logger.With(zap.String("type", "DefaultMediaFilesService")),
		aliyunAccessKeyId:     v.GetString("aliyun.accessKeyId"),
		aliyunAccessKeySecret: v.GetString("aliyun.accessKeySecret"),
		aliyunEndpoint:        v.GetString("aliyun.endpoint"),
		aliyunCallbackUrl:     v.GetString("aliyun.callbackUrl"),
		aliyunUploadDir:       v.GetString("aliyun.uploadDir"),
		aliyunExpireTime:      v.GetInt64("aliyun.expireTime"),
	}
}

func getGmtISO8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func (s *DefaultMediaFileService) GetSign(c context.Context) (*models.AliPolicyToken, error) {
	now := time.Now().Unix()
	expire_end := now + s.aliyunExpireTime
	var tokenExpire = getGmtISO8601(expire_end)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, s.aliyunUploadDir)
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, err := json.Marshal(config)
	if err != nil {
		return nil, models.ErrSignFail
	}
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(s.aliyunAccessKeySecret))
	_, err = io.WriteString(h, debyte)
	if err != nil {
		return nil, models.ErrSignFail
	}
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = s.aliyunCallbackUrl
	callbackParam.CallbackBody = "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		return nil, models.ErrSignFail
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	policyToken := &models.AliPolicyToken{}
	policyToken.AccessKeyId = s.aliyunAccessKeyId
	policyToken.Host = s.aliyunEndpoint
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = s.aliyunUploadDir
	policyToken.Policy = string(debyte)
	policyToken.Callback = string(callbackBase64)

	return policyToken, nil
}

/*  VerifySignature
*   VerifySignature需要三个重要的数据信息来进行签名验证： 1>获取公钥PublicKey;  2>生成新的MD5鉴权串;  3>解码Request携带的鉴权串;
*   1>获取公钥PublicKey : 从RequestHeader的"x-oss-pub-key-url"字段中获取 URL, 读取URL链接的包含的公钥内容， 进行解码解析， 将其作为rsa.VerifyPKCS1v15的入参。
*   2>生成新的MD5鉴权串 : 把Request中的url中的path部分进行urldecode， 加上url的query部分， 再加上body， 组合之后进行MD5编码， 得到MD5鉴权字节串。
*   3>解码Request携带的鉴权串 ： 获取RequestHeader的"authorization"字段， 对其进行Base64解码，作为签名验证的鉴权对比串。
*   rsa.VerifyPKCS1v15进行签名验证，返回验证结果。
* */
func verifySignature(bytePublicKey []byte, byteMd5 []byte, authorization []byte) bool {
	pubBlock, _ := pem.Decode(bytePublicKey)
	if pubBlock == nil {
		fmt.Printf("Failed to parse PEM block containing the public key")
		return false
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if (pubInterface == nil) || (err != nil) {
		fmt.Printf("x509.ParsePKIXPublicKey(publicKey) failed : %s \n", err.Error())
		return false
	}
	pub := pubInterface.(*rsa.PublicKey)

	errorVerifyPKCS1v15 := rsa.VerifyPKCS1v15(pub, crypto.MD5, byteMd5, authorization)
	if errorVerifyPKCS1v15 != nil {
		fmt.Printf("\nSignature Verification is Failed : %s \n", errorVerifyPKCS1v15.Error())
		//printByteArray(byteMd5, "AuthMd5(fromNewAuthString)")
		//printByteArray(bytePublicKey, "PublicKeyBase64")
		//printByteArray(authorization, "AuthorizationFromRequest")
		return false
	}

	fmt.Printf("\nSignature Verification is Successful. \n")
	return true
}

func (s *DefaultMediaFileService) UploadCallback(c context.Context, bytePublicKey []byte, byteMd5 []byte, authorization []byte) bool {

	// verifySignature and response to client
	return verifySignature(bytePublicKey, byteMd5, authorization)
}
