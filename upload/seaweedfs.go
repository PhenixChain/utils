package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/go-ini/ini"
)

// SeaweedfsConf 上传配置
type SeaweedfsConf struct {
	SrvAddr       string
	MaxUploadSize int64
	AllowType     string
}

// InitSeaweedfs 初始化
func InitSeaweedfs() *SeaweedfsConf {
	cfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	cfgSec := cfg.Section("upload")
	seaweedfsConf := new(SeaweedfsConf)
	seaweedfsConf.SrvAddr = cfgSec.Key("srv_url").String()
	maxSize, _ := cfgSec.Key("max_size").Int64()
	seaweedfsConf.MaxUploadSize = maxSize * 1024 // KB
	seaweedfsConf.AllowType = cfgSec.Key("allow_type").String()
	return seaweedfsConf
}

// Upload 上传
func Upload(fileBytes []byte, srvAddr string) (string, error) {
	resp, err := http.Get(srvAddr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var sg SeaweedfsDir
	err = json.Unmarshal(body, &sg)
	if err != nil {
		return "", err
	}
	fileURL := sg.Publicurl + "/" + sg.Fid

	b, err := postBytes(fileURL, "file", fileBytes)
	if err != nil {
		return "", err
	}

	var sw SeaweedfsUpload
	err = json.Unmarshal(b, &sw)
	if err != nil {
		return "", err
	}
	if sw.Size <= 0 {
		return "", fmt.Errorf("upload file fail")
	}

	return fileURL, nil
}

func postBytes(url, fieldname string, data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormField(fieldname)
	if err != nil {
		return nil, err
	}

	if _, err = fw.Write(data); err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SeaweedfsDir ...
type SeaweedfsDir struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	Publicurl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

// SeaweedfsUpload ...
type SeaweedfsUpload struct {
	Etag string `json:"eTag"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}
