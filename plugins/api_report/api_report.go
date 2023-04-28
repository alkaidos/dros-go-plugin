package api_report

import (
	"encoding/json"
	"errors"
	"github.com/alkaidos/dros-go-plugin/plugins"
	"github.com/alkaidos/dros-go-plugin/proxy/api_manage_vo"
	"github.com/alkaidos/dros-go-plugin/proxy/base_proxy_vo"
	"github.com/alkaidos/dros-go-plugin/proxy/isc_category_vo"
	"github.com/alkaidos/dros-go-plugin/proxy/user_proxy_vo"
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
	if !register && plugins.PluginConfig.ApiConf.Enable {
		logger.Warn("注册服务api信息:开始")
		err := registerApi()
		if err != nil {
			panic("请求注册服务api信息失败,原因:" + err.Error())
		}
		logger.Warn("注册服务api信息:成功")
		register = true
	}
}

func registerApi() error {
	//处理请求查询字段
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	token, err := getToken()
	if err != nil {
		return errors.New("注册服务api信息:获取应用授权token失败" + err.Error())
	}
	header.Add("token", token)
	parameterMap := map[string]string{}
	if plugins.PluginConfig.ApiConf.Type != "" {
		parameterMap["type"] = plugins.PluginConfig.ApiConf.Type
	}
	url := plugins.PluginConfig.ApiConf.RegisterHost
	swaggerData, err := plugins.GetFileInfo(plugins.PluginConfig.ApiConf.SwaggerFilePath)
	if err != nil {
		//logger.Error("注册服务api信息:读取swaggerInfo失败", err.Error())
		return errors.New("注册服务api信息:读取swaggerInfo失败" + err.Error())
	}
	var post Post
	json.Unmarshal(swaggerData, &post)
	_, _, data, err := baseHttp.Post(url+"/api/orchestration/capc/import/dynamic", header, parameterMap, post)
	if err != nil {
		//logger.Error("注册服务api信息:请求能力中心接口失败", err.Error())
		return errors.New("注册服务api信息:请求能力中心接口失败" + err.Error())
	}
	var rsp base_proxy_vo.HttpResult[any]
	err = json.Unmarshal(data.([]byte), &rsp)
	if err != nil {
		//logger.Error("注册服务api信息:请求能力中心解析结果失败", err.Error())
		return errors.New("注册服务api信息:请求能力中心解析结果失败" + err.Error())
	}
	if rsp.Code != 0 {
		//logger.Error("注册服务api信息:请求能力中心返回结果失败，原因:" + rsp.Message)
		return errors.New("注册服务api信息:请求能力中心返回结果失败，原因:" + rsp.Message)
	}
	categoryId, err := getCategoryId(header)
	if err != nil {
		return err
	}

	err = updateAppInfoByCategoryId(categoryId)
	if err != nil {
		return err
	}
	return nil
}

func getCategoryId(header http.Header) (int, error) {

	// 获取上传的categoryId
	caIDParamsMap := map[string]string{}
	caIDParamsMap["name"] = plugins.PluginConfig.ApiConf.GroupName
	caIDParamsMap["type"] = plugins.PluginConfig.ApiConf.Type
	_, _, data, err := baseHttp.Get(plugins.PluginConfig.ApiConf.RegisterHost+"/api/orchestration/capc/category/page", header, caIDParamsMap)
	if err != nil {
		return -1, errors.New("注册服务api信息:请求查询应用分组id失败" + err.Error())
	}
	var caIdRsp base_proxy_vo.HttpResult[[]isc_category_vo.IscCategoryDTO]
	err = json.Unmarshal(data.([]byte), &caIdRsp)
	if err != nil {
		//logger.Error("注册服务api信息:请求能力中心解析结果失败", err.Error())
		return -1, errors.New("注册服务api信息:查询应用分组id失败" + err.Error())
	}
	if caIdRsp.Code != 0 {
		return -1, errors.New("注册服务api信息:请求能力中心返回结果失败，原因:" + caIdRsp.Message)
	}

	if len(caIdRsp.Data) <= 0 {
		return -1, errors.New("注册服务api信息:上传api失败，原因:找不到上传后的分组")
	}

	for _, v := range caIdRsp.Data {
		if v.Name == plugins.PluginConfig.ApiConf.GroupName {
			return v.Id, nil
		}
	}
	return -1, nil
}

func updateAppInfoByCategoryId(categoryId int) error {
	paramsMap := map[string]any{}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	paramsMap["serviceCode"] = plugins.PluginConfig.ApiConf.ApiServiceCode
	paramsMap["categoryId"] = categoryId
	_, _, data, err := baseHttp.Put(plugins.PluginConfig.ApiConf.ApiManagerHost+"/api/app/api-manage/producer/update", header, nil, paramsMap)
	if err != nil {
		return errors.New("更新api服务基本信息,请求失败:" + err.Error())
	}
	var rsp base_proxy_vo.HttpResult[api_manage_vo.ApiProducerInfo]
	err = json.Unmarshal(data.([]byte), &rsp)
	if err != nil {
		//logger.Error("注册服务api信息:请求能力中心解析结果失败", err.Error())
		return errors.New("更新api服务基本信息:" + err.Error())
	}
	if rsp.Code != 0 {
		return errors.New("注册服务api信息:请求能力中心返回结果失败，原因:" + rsp.Message)
	}

	return nil
}

func getToken() (string, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	bodyMap := map[string]any{}
	bodyMap["appCode"] = "tddm"
	bodyMap["appSecret"] = "e6874c0a4b397e3d6b59"
	_, _, data, err := baseHttp.Post(plugins.PluginConfig.ApiConf.PermissionHost+"/api/permission/app/token/grant", header, nil, bodyMap)
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
