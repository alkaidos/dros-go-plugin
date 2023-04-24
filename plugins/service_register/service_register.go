package service_register

import (
	pluginConfig "dros-go-plugin/plugins"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/logger"
	"net/http"
	"strings"
)

var register = false

func init() {

	// api注册
	if !register && pluginConfig.PluginConfig.AppConf.Enable {
		registerServiceRoute()
		register = true
	}

}

func registerServiceRoute() {
	registerService()
	registerRoute()
}

func registerService() {
	logger.Warn("开始请求注册服务信息")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	parameter := map[string]any{}
	parameter["serviceId"] = pluginConfig.PluginConfig.AppConf.ServiceId
	parameter["serviceName"] = pluginConfig.PluginConfig.AppConf.ServiceName
	parameter["servicePath"] = pluginConfig.PluginConfig.AppConf.ServicePath
	parameter["serviceUrl"] = pluginConfig.PluginConfig.AppConf.ServiceUrl
	parameter["serviceEnable"] = true
	_, _, _, err := baseHttp.Post(pluginConfig.PluginConfig.AppConf.RegisterAppUrl+"/api/rc-application/open/service/register", header, nil, parameter)
	if err != nil {
		logger.Error("请求注册服务信息失败", err.Error())
	} else {
		logger.Warn("请求注册服务信息成功")
	}
}

func registerRoute() {

	logger.Warn("开始请求注册路由信息")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	parameter := map[string]any{}
	parameter["serviceId"] = pluginConfig.PluginConfig.AppConf.ServiceId
	parameter["url"] = pluginConfig.PluginConfig.AppConf.ServiceUrl
	parameter["path"] = pluginConfig.PluginConfig.AppConf.ServicePath
	parameter["excludeUrl"] = strings.Split(pluginConfig.PluginConfig.AppConf.ExcludeUrl, ";")
	_, _, _, err := baseHttp.Put(pluginConfig.PluginConfig.AppConf.RegisterAppUrl+"/api/route/update/exclude", header, nil, parameter)
	if err != nil {
		logger.Error("请求注册路由信息失败", err.Error())
	} else {
		logger.Warn("请求注册路由信息成功")
	}
}
