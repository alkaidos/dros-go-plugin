package api_report

import (
	"dros-go-plugin/plugins"
	pluginConfig "dros-go-plugin/plugins"
	"dros-go-plugin/proxy/base_proxy_vo"
	"dros-go-plugin/proxy/user_proxy_vo"
	"encoding/json"
	"errors"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/logger"
	"net/http"
)

var register = false

type Post struct { //带结构标签，反引号来包围字符串
	Swagger             string         `json:"swagger"`
	Info                map[string]any `json:"info"`
	Tags                []any          `json:"tags"`
	Path                map[string]any `json:"paths"`
	SecurityDefinitions map[string]any `json:"securityDefinitions"`
	Definitions         map[string]any `json:"definitions"`
}

func init() {
	// api注册
	if !register && pluginConfig.PluginConfig.ApiConf.Enable {
		registerApi()
		register = true
	}
}

func registerApi() {
	//处理请求查询字段
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	token := "607504c5-75d9-4f43-a729-d7991d40e6ca"
	header.Add("token", token)
	parameterMap := map[string]string{}
	//parameterMap["type"] = "UDMP"
	url := pluginConfig.PluginConfig.ApiConf.RegisterHost
	swaggerData, err := plugins.GetFileInfo(pluginConfig.PluginConfig.ApiConf.SwaggerFilePath)
	if err != nil {
		logger.Error("注册服务api信息:读取swaggerInfo失败", err.Error())
		return
	}
	var post Post
	json.Unmarshal(swaggerData, &post)
	_, _, data, err := baseHttp.Post(url+"/api/orchestration/capc/import/dynamic", header, parameterMap, post)
	if err != nil {
		logger.Error("注册服务api信息:请求能力中心接口失败", err.Error())
		return
	}
	var rsp base_proxy_vo.HttpResult[any]
	err = json.Unmarshal(data.([]byte), &rsp)
	if err != nil {
		logger.Error("注册服务api信息:请求能力中心解析结果失败", err.Error())
		return
	}
	if rsp.Code != 0 {
		logger.Error("注册服务api信息:请求能力中心返回结果失败，原因:" + rsp.Message)
		return
	}
	logger.Warn("注册服务api信息:执行成功")
}

func getToken() (string, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	bodyMap := map[string]any{}
	bodyMap["appCode"] = "tddm"
	bodyMap["appSecret"] = "e6874c0a4b397e3d6b59"
	url := "http://10.101.12.156:38080"
	_, _, data, err := baseHttp.Post(url+"/api/permission/app/token/grant", header, nil, bodyMap)
	if err != nil {
		logger.Error("上报api：请求获取token失败", err.Error())
		return "", errors.New("上报api：请求获取token失败:" + err.Error())
	} else {
		var rsp base_proxy_vo.HttpResult[user_proxy_vo.AppAuthenticationDTO]
		err = json.Unmarshal(data.([]byte), &rsp)
		if err != nil {
			return "", errors.New("上报api-请求获取token：解析返回结果失败:" + err.Error())
		}

		if rsp.Code != 0 {
			return "", errors.New("上报api-请求获取token：请求返回失败:" + rsp.Message)
		}
		return rsp.Data.Token, nil
	}
}
