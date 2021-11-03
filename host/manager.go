package host

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/cockroachdb/errors"
	"github.com/youthlin/plugin/common"
)

const metaFile = "plugin.json"

// Manager 宿主端插件管理器：扫描插件列表、启用插件、停用插件、调用插件
type Manager struct {
	// key is plugin id
	Plugins map[string]*Stub
	// key is action name
	Actions map[string]sortableActions
}

// Scan the plugins dir and load unloaded plugins' meta
//  @param root: the plugins dir
//  @return: error message
//
//  root
//  |-- plugin1/
//  |      |- plugin1-exe
//  |      \- plugin.json
//  \-- plugin2/plugin2.json
func (m *Manager) Scan(root string) error {
	pluginsDir, err := os.Open(root)
	if err != nil {
		return errors.Wrapf(err,
			// TRANSLATORS: %s is the plugins root dir
			ts.T("failed to open plugins dir: %s", root))
	}
	entries, err := pluginsDir.ReadDir(0)
	if err != nil {
		return errors.Wrapf(err,
			// TRANSLATORS: %s is the plugins root dir
			ts.T("failed to read plugins dir: %s", root))
	}
	var errors []error
	for _, entry := range entries {
		if entry.IsDir() {
			dir := filepath.Join(root, entry.Name())
			m.Load(dir)
		}
	}
	return common.CombineErrors(errors...)
}

func (m *Manager) Load(dir string) error {
	name := filepath.Join(dir, metaFile)
	file, err := os.Open(name)
	if err != nil {
		return errors.Wrapf(err, ts.T("failed to open plugin meta file"))
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err,
			// TRANSLATORS: %s is the plugin's meta file: "plugin.json"
			ts.T("failed to read plugin meta file: %s", name))
	}
	var meta Meta
	err = json.Unmarshal(content, &meta)
	if err != nil {
		return err
	}
	meta.Dir = dir
	if err = meta.Validate(); err != nil {
		return errors.Wrapf(err,
			// TRANSLATORS: %s is plugin meta file
			ts.T("plugin is not valid: %s", name))
	}
	stub := &Stub{Meta: &meta}
	m.Plugins[stub.ID()] = stub
	return nil
}

func (m *Manager) StartPlugin(ctx context.Context, id string) {}
func (m *Manager) StopPlugin(ctx context.Context, id string)  {}

func (m *Manager) DoAction(ctx context.Context, actionName string, args ...interface{}) {
	m.ApplyFilter(ctx, actionName, nil, args...)
}

func (m *Manager) ApplyFilter(ctx context.Context, filterName string, value interface{}, args ...interface{}) interface{} {
	actions := m.Actions[filterName]
	for _, action := range actions {
		stub := m.Plugins[action.id]
		_ = stub
	}
	return nil
}

type action struct {
	*option
	id string
}

type option struct {
	priority int
	before   string
	after    string
}

type Option func(*option)

func WithPriority(priority int) Option {
	return func(o *option) { o.priority = priority }
}
func BeforePluginID(id string) Option {
	return func(o *option) { o.before = id }
}
func AfterPluginID(id string) Option {
	return func(o *option) { o.after = id }
}

func (m *Manager) AddAction(ctx context.Context, actionName string, pluginID string, options ...Option) {
	m.AddFilter(ctx, actionName, pluginID, options...)
}
func (m *Manager) AddFilter(ctx context.Context, filterName string, pluginID string, options ...Option) {
	option := &option{}
	for _, fun := range options {
		fun(option)
	}
	actions := m.Actions[filterName]
	actions = append(actions, action{
		option: option,
		id:     pluginID,
	})
	sort.Sort(actions)
	m.Actions[filterName] = actions
}

type sortableActions []action

func (a sortableActions) Len() int      { return len(a) }
func (a sortableActions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortableActions) Less(i, j int) bool {
	left := a[i]
	right := a[j]
	if left.priority < right.priority {
		return true
	}
	if left.before == right.id {
		return true
	}
	if right.after == left.id {
		return true
	}
	return false
}
