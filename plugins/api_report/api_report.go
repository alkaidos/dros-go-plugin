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
	"io"
	"net/http"
	"strconv"
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
		logger.Warn("注册服务api插件执行开始")
		err := registerApi()
		if err != nil {
			panic(err.Error())
		}
		logger.Warn("注册服务api插件执行成功")
		register = true
	}
}

func registerApi() error {
	//处理请求查询字段
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	token, err := getToken()
	if err != nil {
		return err
	}
	logger.Warn("请求能力中心注册服务api信息:获取应用授权token成功，" + token)
	header.Add("token", token)
	parameterMap := map[string]string{}
	if plugins.PluginConfig.ApiConf.Type != "" {
		parameterMap["type"] = plugins.PluginConfig.ApiConf.Type
	}
	url := plugins.PluginConfig.ApiConf.RegisterHost
	swaggerData, err := plugins.GetFileInfo(plugins.PluginConfig.ApiConf.SwaggerFilePath)
	if err != nil {
		return errors.New("请求能力中心注册服务api信息:读取swaggerInfo失败，原因" + err.Error())
	}
	logger.Warn("请求能力中心注册服务api信息:读取swagger文件成功")
	var post Post
	json.Unmarshal(swaggerData, &post)
	_, _, data, err := baseHttp.Post(url+"/api/orchestration/capc/import/dynamic", header, parameterMap, post)
	if err != nil {
		return errors.New("请求能力中心注册服务api信息:请求能力中心接口失败，原因" + err.Error())
	}
	var rsp base_proxy_vo.HttpResult[any]
	err = json.Unmarshal(data.([]byte), &rsp)
	if err != nil {
		return errors.New("请求能力中心注册服务api信息:能力中心返回结果解析失败，原因" + err.Error())
	}
	if rsp.Code != 0 {
		return errors.New("请求能力中心注册服务api信息:请求能力中心返回结果失败，原因" + rsp.Message)
	}
	logger.Warn("注册服务api插件执行情况:往能力中心上报api信息成功，开始获取api分组id")
	categoryId, err := getCategoryId(token)
	if err != nil {
		return err
	}
	if categoryId <= -1 {
		return errors.New("获取api分组id:获取不到应用在能力中心的注册的分组id")
	}
	logger.Warn("注册服务api插件执行情况:往能力中心上报api信息成功，读取上报后的api分组id成功，值为" + strconv.Itoa(categoryId))
	logger.Warn("注册服务api插件执行情况:往开放平台上报分组id信息，开始")
	err = updateAppInfoByCategoryId(categoryId)
	if err != nil {
		return err
	}
	logger.Warn("注册服务api插件执行情况:往开放平台上报分组id信息，成功")
	return nil
}

func getCategoryId(token string) (int, error) {
	header := http.Header{}
	header.Add("token", token)
	header.Set("Content-Type", "application/json")
	client := http.Client{}
	// 获取上传的categoryId
	httpRequest, err := http.NewRequest("GET",
		plugins.PluginConfig.ApiConf.RegisterHost+"/api/orchestration/capc/category/page", nil)
	if err != nil {
		return -1, errors.New("获取api分组id:构造请求参数失败" + err.Error())
	}
	q := httpRequest.URL.Query()
	q.Add("name", plugins.PluginConfig.ApiConf.GroupName)
	q.Add("type", plugins.PluginConfig.ApiConf.Type)
	httpRequest.URL.RawQuery = q.Encode()
	if header != nil {
		httpRequest.Header = header
	}
	resp, err := client.Do(httpRequest)
	if err != nil {
		return -1, errors.New("获取api分组id:请求能力中心接口失败，原因" + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return -1, errors.New("获取api分组id:请求查询应用分组id失败,status code error:" + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	var caIdRsp base_proxy_vo.HttpResult[[]isc_category_vo.IscCategoryDTO]
	err = json.Unmarshal(body, &caIdRsp)
	if err != nil {
		return -1, errors.New("获取api分组id:能力中心返回结果解析失败，原因" + err.Error())
	}
	if caIdRsp.Code != 0 {
		return -1, errors.New("获取api分组id:请求能力中心返回结果失败，原因" + caIdRsp.Message)
	}

	if len(caIdRsp.Data) <= 0 {
		return -1, errors.New("获取api分组id:找不到上传api信息后的分组")
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
		return errors.New("开放平台上报分组id信息:请求失败，原因" + err.Error())
	}
	var rsp base_proxy_vo.HttpResult[api_manage_vo.ApiProducerInfo]
	err = json.Unmarshal(data.([]byte), &rsp)
	if err != nil {
		return errors.New("开放平台上报分组id信息:解析返回结果失败，原因" + err.Error())
	}
	if rsp.Code != 0 {
		return errors.New("开放平台上报分组id信息:开放平台返回结果失败，原因" + rsp.Message)
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
		return "", errors.New("获取token:请求用户中心获取token失败，原因" + err.Error())
	} else {
		var rsp base_proxy_vo.HttpResult[user_proxy_vo.AppAuthenticationDTO]
		err = json.Unmarshal(data.([]byte), &rsp)
		if err != nil {
			return "", errors.New("获取token:返回结果解析失败，原因" + err.Error())
		}

		if rsp.Code != 0 {
			return "", errors.New("获取token:返回失败，原因" + rsp.Message)
		}
		return rsp.Data.Token, nil
	}
}
