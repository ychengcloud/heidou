package models

//MediaFile ...
type MediaFile struct {
}

type AliPolicyToken struct {
	AccessKeyId string `json:"accessId"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

func (m *MediaFile) Upload() {

}
