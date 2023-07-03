package config_center

import (
	"encoding/json"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
	"time"
)

var tenantId = "system"

func GetConfig(appName, key string) (*RespStringDTO[[]SysCommonConfiguration], error) {
	configUrl := config.GetValueString("app.configUrl")
	url := configUrl + "/config/item/getConfigList?appName=" + appName
	if key != "" {
		url = url + "&key=" + key
	}
	req, requestErr := http.NewRequest("GET", url, nil)
	if requestErr != nil {
		return nil, errors.Wrap(requestErr, requestErr.Error())
	}
	header := req.Header.Clone()
	header.Add("Content-Type", "application/json")
	header.Add("isc-biz-tenant-id", tenantId)
	req.Header = header

	client := http.Client{Timeout: time.Second * 5}
	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, errors.Wrap(resErr, resErr.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("%s", err.Error())
		}
	}(res.Body)
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, errors.Wrap(readErr, readErr.Error())
	}
	var resDto RespStringDTO[[]SysCommonConfiguration]
	if unmarshalErr := json.Unmarshal(body, &resDto); unmarshalErr != nil {
		return nil, errors.Wrap(unmarshalErr, unmarshalErr.Error())
	}
	if len(resDto.Data) == 0 {
		switch key {
		case "logo":
			vo := SysCommonConfigurationVo{Key: "logo", AppName: "isc-sys-config-service", Profile: "default"}
			rs, err := UpdateConfig(&vo)
			if err != nil {
				return nil, errors.Wrap(err, err.Error())
			}
			if rs.Code != "success" && isc.ToString(rs.Code) != "0" {
				return nil, errors.New(rs.Message)
			}
		}
	}
	return &resDto, nil
}

func InitConfig() {
	InitConfigAndRegister("system")
}

func InitConfigAndRegister(tenantId string) {
	values := []ConfigAppValueReq{{Key: "logo", Value: "", AllowPush: 1}, {Key: "enterprise", Value: "", AllowPush: 1},
		{Key: "phoneNum", Value: "", AllowPush: 1}, {Key: "powerby", Value: "", AllowPush: 1}}
	v3Req := ConfigUploadV3Req{Version: 1, Profile: "default", AppName: "isc-sys-config-service", Group: "default", ValueList: values, ProjectType: "isc-os"}
	configUrl := config.GetValueString("app.configUrl")
	marshal, marshalErr := json.Marshal(v3Req)
	if marshalErr != nil {
		logger.Error(marshalErr.Error())
	}
	req, requestErr := http.NewRequest("POST", configUrl+"/client/uploadConfigV3", strings.NewReader(string(marshal)))
	if requestErr != nil {
		logger.Error(requestErr.Error())
		return
	}
	header := req.Header.Clone()
	header.Add("Content-Type", "application/json")
	header.Add("isc-biz-tenant-id", tenantId)
	req.Header = header

	client := http.Client{Timeout: time.Second * 5}
	res, resErr := client.Do(req)
	if resErr != nil {
		//retry
		logger.Error(resErr.Error())
		time.Sleep(time.Minute)
		InitConfigAndRegister("system")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("%s", err.Error())
		}
	}(res.Body)
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		logger.Error(readErr.Error())
	}
	var resDto RespStringDTO[int]
	if unmarshalErr := json.Unmarshal(body, &resDto); unmarshalErr != nil {
		logger.Error(unmarshalErr.Error())
	}
	if resDto.Code != "success" && isc.ToString(resDto.Code) != "0" {
		logger.Error(resDto.Message)
	}
}

func UpdateConfig(configurationVo *SysCommonConfigurationVo) (*RespStringDTO[int], error) {
	configUrl := config.GetValueString("app.configUrl")
	marshal, marshalErr := json.Marshal(configurationVo)
	if marshalErr != nil {
		return nil, errors.Wrap(marshalErr, marshalErr.Error())
	}
	req, requestErr := http.NewRequest("PUT", configUrl+"/config/item/updateFromKey", strings.NewReader(string(marshal)))
	if requestErr != nil {
		return nil, errors.Wrap(requestErr, requestErr.Error())
	}
	header := req.Header.Clone()
	header.Add("Content-Type", "application/json")
	header.Add("isc-biz-tenant-id", tenantId)
	req.Header = header

	client := http.Client{Timeout: time.Second * 5}
	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, errors.Wrap(resErr, resErr.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("%s", err.Error())
		}
	}(res.Body)
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, errors.Wrap(readErr, readErr.Error())
	}
	var resDto RespStringDTO[int]

	if unmarshalErr := json.Unmarshal(body, &resDto); unmarshalErr != nil {
		return nil, errors.Wrap(unmarshalErr, unmarshalErr.Error())
	}
	return &resDto, nil
}

func GetConfigStringValue(appName, key string) (string, error) {
	rs, error := GetConfig(appName, key)
	if error != nil {
		return "", error
	}
	if rs.Code != "success" && isc.ToString(rs.Code) != "0" {
		return "", errors.New(rs.Message)
	}
	var data string
	if len(rs.Data) > 0 {
		data = rs.Data[0].Value
	}
	return data, nil
}

type RespStringDTO[V any] struct {
	Code    any
	Message string
	Data    V
}

type ConfigAppValueReq struct {
	Key          string `json:"key"`          //
	Value        string `json:"value"`        //
	PasswordType int    `json:"passwordType"` //
	AllowPush    int    `json:"allowPush"`    //
	Desc         string `json:"desc"`         //
	ValueType    string `json:"valueType"`    //
}

type ConfigUploadV3Req struct {
	Profile     string              `json:"profile"`     //
	AppName     string              `json:"appName"`     //
	Version     int                 `json:"version"`     //
	Group       string              `json:"group"`       //
	UpgradeKeys []string            `json:"upgradeKeys"` //
	ProjectType string              `json:"projectType"` //
	ValueList   []ConfigAppValueReq `json:"valueList"`   //
}

type SysCommonConfiguration struct {
	Id          int64  `json:"id"`
	Key         string `json:"key"`         //配置名称
	Value       string `json:"value"`       //配置内容
	Application string `json:"application"` //服务名称
	Profile     string `json:"profile"`     //配置
	Description string `json:"description"` //描述
}

type SysCommonConfigurationVo struct {
	Key     string `json:"key"`     //配置名称
	Value   string `json:"value"`   //配置内容
	AppName string `json:"appName"` //服务名称
	Profile string `json:"profile"` //配置 default
}
