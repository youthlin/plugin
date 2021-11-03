package common

// Plugin 插件接口
type Plugin interface {
	ID() string
	OnStart()
	OnStop()
}
