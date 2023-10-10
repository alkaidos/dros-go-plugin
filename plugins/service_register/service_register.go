package service_register

import (
	"encoding/json"
	"github.com/alkaidos/dros-go-plugin/plugins"
	"github.com/alkaidos/dros-go-plugin/proxy/base_proxy_vo"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/logger"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

var register = false
var appCodeRegexp = "[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"
var servicePathPrefix = "/api/app/"

func init() {

	// api注册
	if !register && plugins.PluginConfig.AppConf.Enable {
		if plugins.PluginConfig.AppConf.IsMainService {
			registerServiceRouteForMain()
		} else {
			registerServiceRouteForSubService()
		}
		register = true
	}

}

func registerServiceRouteForMain() {
	// 检查应用基本信息
	checkInfoForMainService()
	// 注册应用信息（应用信息+服务基本信息+服务路由信息）
	registerMainService()
}

func registerServiceRouteForSubService() {
	// 检查服务基本信息
	checkInfoForSubService()
	// 注册服务基本信息
	registerSubServiceInfo()
	// 注册服务路由信息
	registerSubServiceRoute()
}

func registerMainService() {
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	header.Set("isc-api-version", "2.1")
	header.Set("isc-tenant-id", "system")
	parameter := map[string]any{}
	parameter["appId"] = plugins.PluginConfig.AppConf.AppId
	parameter["code"] = plugins.PluginConfig.AppConf.AppCode
	parameter["name"] = plugins.PluginConfig.AppConf.AppName
	parameter["version"] = plugins.PluginConfig.AppConf.AppVersion
	parameter["inMenu"] = plugins.PluginConfig.AppConf.InMenu
	parameter["serviceList"] = plugins.PluginConfig.AppConf.ServiceList
	if plugins.PluginConfig.AppConf.Type > 0 {
		parameter["type"] = plugins.PluginConfig.AppConf.Type
		if plugins.PluginConfig.AppConf.Type == 10 {
			parameter["inMenu"] = 2
		}
	}
	if plugins.PluginConfig.AppConf.SecondType > 0 {
		parameter["secondType"] = plugins.PluginConfig.AppConf.SecondType
	}

	if plugins.PluginConfig.AppConf.RedirectUrl != "" {
		parameter["redirectUrl"] = plugins.PluginConfig.AppConf.RedirectUrl
	}
	_, _, body, err := baseHttp.Post(plugins.PluginConfig.AppConf.RegisterAppUrl+"/api/rc-application/application/register", header, nil, parameter)
	if err != nil {
		panic("请求注册应用信息失败,原因:" + err.Error())
	} else {
		var rsp base_proxy_vo.HttpResult[any]
		err = json.Unmarshal(body.([]byte), &rsp)
		if rsp.Code != 0 {
			panic("注册服务接口返回失败，原因：" + rsp.Message)
		}
		logger.Warn("请求注册应用信息成功")
	}
}

func checkInfoForSubService() {
	service := plugins.ServiceRegisterInfo{}
	service.ServiceId = plugins.PluginConfig.AppConf.ServiceId
	service.ServiceName = plugins.PluginConfig.AppConf.ServiceName
	service.ServicePath = plugins.PluginConfig.AppConf.ServicePath
	service.ServiceUrl = plugins.PluginConfig.AppConf.ServiceUrl
	errorMsg := checkServiceInfoIllegal(service)
	if errorMsg != "" {
		panic(errorMsg)
	}
}

func checkInfoForMainService() {
	if plugins.PluginConfig.AppConf.AppCode == "" {
		panic("appCode 不能为空")
	}

	match, err := regexp.MatchString(appCodeRegexp, plugins.PluginConfig.AppConf.AppCode)
	if err != nil {
		panic("appCode 校验异常" + err.Error())
	}
	if !match {
		panic("appCode必须配置并且满足正则表达式:[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\n")
	}

	for _, service := range plugins.PluginConfig.AppConf.ServiceList {
		if !strings.HasSuffix(service.ServiceId, "-ui") {
			errorMsg := checkServiceInfoIllegal(service)
			if errorMsg != "" {
				panic(errorMsg)
			}
		}
	}
}

func checkServiceInfoIllegal(service plugins.ServiceRegisterInfo) string {
	if service.ServiceId == "" {
		return service.ServiceId + ":serviceId不能为空"
	}

	for _, v := range service.ServiceId {
		if unicode.Is(unicode.Han, v) {
			return service.ServiceId + ":serviceI不能有汉字"
		}
		if unicode.IsUpper(v) {
			return service.ServiceId + ":serviceId不能输入大写字母"
		}
	}

	if service.ServicePath != "" && !strings.HasPrefix(service.ServicePath, servicePathPrefix) {
		return service.ServiceId + ":servicePath必须以/api/app/作为前缀, 推荐配置为/api/app/${isyscore.appCode}/**\n"
	}

	return ""
}

func registerSubServiceInfo() {
	logger.Warn("开始请求注册服务信息")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	parameter := map[string]any{}
	parameter["serviceId"] = plugins.PluginConfig.AppConf.ServiceId
	parameter["serviceName"] = plugins.PluginConfig.AppConf.ServiceName
	parameter["servicePath"] = plugins.PluginConfig.AppConf.ServicePath
	parameter["serviceUrl"] = plugins.PluginConfig.AppConf.ServiceUrl
	parameter["serviceEnable"] = true
	_, _, body, err := baseHttp.Post(plugins.PluginConfig.AppConf.RegisterAppUrl+"/api/rc-application/open/service/register", header, nil, parameter)
	if err != nil {
		//logger.Error("请求注册服务信息失败", err.Error())
		panic("请求注册服务信息失败,原因:" + err.Error())
	} else {
		var rsp base_proxy_vo.HttpResult[any]
		err = json.Unmarshal(body.([]byte), &rsp)
		if rsp.Code != 0 {
			panic("注册服务接口返回失败，原因：" + rsp.Message)
		}
		logger.Warn("请求注册服务信息成功")
	}
}

func registerSubServiceRoute() {

	logger.Warn("开始请求注册路由信息")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	parameter := map[string]any{}
	parameter["serviceId"] = plugins.PluginConfig.AppConf.ServiceId
	parameter["url"] = plugins.PluginConfig.AppConf.ServiceUrl
	parameter["path"] = plugins.PluginConfig.AppConf.ServicePath
	parameter["excludeUrl"] = strings.Split(plugins.PluginConfig.AppConf.ExcludeUrl, ";")
	_, _, body, err := baseHttp.Put(plugins.PluginConfig.AppConf.RegisterAppUrl+"/api/route/update/exclude", header, nil, parameter)
	if err != nil {
		//logger.Error("请求注册路由信息失败", err.Error())
		panic("请求注册路由信息失败,原因:" + err.Error())
	} else {
		var rsp base_proxy_vo.HttpResult[any]
		err = json.Unmarshal(body.([]byte), &rsp)
		if rsp.Code != 0 {
			panic("注册路由信息接口返回失败，原因：" + rsp.Message)
		}
		logger.Warn("请求注册路由信息成功")
	}
}
