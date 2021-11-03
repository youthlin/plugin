package host

// Stub 是宿主端存根，调用 Stub 就像调用插件一样
type Stub struct {
	*Meta
}

func (m *Stub) ID() string { return m.Name }
func (m *Stub) OnStart()   {}
func (m *Stub) OnStop()    {}
