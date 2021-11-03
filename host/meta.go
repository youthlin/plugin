package host

import (
	"strings"

	"github.com/cockroachdb/errors"
)

// Meta used to decribe the plugin
type Meta struct {
	Dir            string   `json:"-"` // the plugin dir, fill on runtime
	Name           string   // the plugin name
	Cmd            []string // how to start the plugin
	MinHostVersion int      // required minimal host version
}

func (m *Meta) Validate() error {
	m.Name = strings.TrimSpace(m.Name)
	if m.Name == "" {
		return errors.Errorf(ts.T("plugin name must not be empty"))
	}
	if len(m.Cmd) == 0 {
		return errors.Errorf(ts.T("start command must not be empty"))
	}
	m.ResolveCmdPath()
	return nil
}

func (m *Meta) ResolveCmdPath() {
	var cmd = make([]string, 0, len(m.Cmd))
	for _, arg := range m.Cmd {
		cmd = append(cmd, strings.ReplaceAll(arg, "${dir}", m.Dir))
	}
	m.Cmd = cmd
}
