package service_register

import (
	"github.com/alkaidos/dros-go-plugin/plugins"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/logger"
	"net/http"
	"strings"
)

var register = false

func init() {

	// api注册
	if !register && plugins.PluginConfig.AppConf.Enable {
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
	parameter["serviceId"] = plugins.PluginConfig.AppConf.ServiceId
	parameter["serviceName"] = plugins.PluginConfig.AppConf.ServiceName
	parameter["servicePath"] = plugins.PluginConfig.AppConf.ServicePath
	parameter["serviceUrl"] = plugins.PluginConfig.AppConf.ServiceUrl
	parameter["serviceEnable"] = true
	_, _, _, err := baseHttp.Post(plugins.PluginConfig.AppConf.RegisterAppUrl+"/api/rc-application/open/service/register", header, nil, parameter)
	if err != nil {
		//logger.Error("请求注册服务信息失败", err.Error())
		panic("请求注册服务信息失败,原因:" + err.Error())
	} else {
		logger.Warn("请求注册服务信息成功")
	}
}

func registerRoute() {

	logger.Warn("开始请求注册路由信息")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=UTF-8")
	parameter := map[string]any{}
	parameter["serviceId"] = plugins.PluginConfig.AppConf.ServiceId
	parameter["url"] = plugins.PluginConfig.AppConf.ServiceUrl
	parameter["path"] = plugins.PluginConfig.AppConf.ServicePath
	parameter["excludeUrl"] = strings.Split(plugins.PluginConfig.AppConf.ExcludeUrl, ";")
	_, _, _, err := baseHttp.Put(plugins.PluginConfig.AppConf.RegisterAppUrl+"/api/route/update/exclude", header, nil, parameter)
	if err != nil {
		//logger.Error("请求注册路由信息失败", err.Error())
		panic("请求注册路由信息失败,原因:" + err.Error())
	} else {
		logger.Warn("请求注册路由信息成功")
	}
}
