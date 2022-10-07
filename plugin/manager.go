package plugin

import (
	"fmt"
	plugin "github.com/hashicorp/go-plugin"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Type string

const (
	TypeSource   Type = "sources"
	TypeExecutor Type = "executors"
)

type Metadata struct {
	ID     string
	Path   string
	Client *plugin.Client
	Type   Type
}

type Manager struct {
	Type           Type
	Path           string
	Plugins        map[string]*Metadata
	implementation plugin.Plugin
}

func (m *Manager) Glob() string {
	return fmt.Sprintf("%s-*", m.Type)
}

func NewManager(dir string, pType Type, implementation plugin.Plugin) *Manager {
	return &Manager{
		Type:           pType,
		Path:           dir,
		Plugins:        map[string]*Metadata{},
		implementation: implementation,
	}
}
func (m *Manager) Initialize() error {
	plugins, err := plugin.Discover(m.Glob(), m.Path)
	if err != nil {
		return err
	}

	for _, p := range plugins {
		pluginID := m.getPluginID(p)
		m.Plugins[pluginID] = &Metadata{
			ID:   pluginID,
			Path: p,
			Type: m.Type,
		}

	}

	return nil
}

func (m *Manager) Start() error {
	for id, metadata := range m.Plugins {
		log.Printf("Registering plugin %s[id=%s]", id, metadata.Path)
		client := plugin.NewClient(&plugin.ClientConfig{
			Plugins:          m.pluginMap(id),
			VersionedPlugins: nil,
			Cmd:              exec.Command(metadata.Path),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "BOTKUBE_MAGIC_COOKIE",
				MagicCookieValue: "BOTKUBE_BASIC_PLUGIN",
			},
		})
		pinfo := m.Plugins[id]
		pinfo.Client = client
	}
	return nil
}

func (m *Manager) Dispose() {
	var wg sync.WaitGroup
	for _, p := range m.Plugins {
		wg.Add(1)

		go func(client *plugin.Client) {
			client.Kill()
			wg.Done()
		}(p.Client)
	}

	wg.Wait()

}

func (m *Manager) GetAdapter(id string) (interface{}, error) {

	if _, ok := m.Plugins[id]; !ok {
		return nil, fmt.Errorf("plugin: %s not found in registered plugins", id)
	}

	client := m.Plugins[id].Client

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(id)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (m *Manager) getPluginID(path string) string {
	_, file := filepath.Split(path)
	return strings.TrimPrefix(file, strings.Split(m.Glob(), "*")[0])
}

func (m *Manager) pluginMap(id string) map[string]plugin.Plugin {
	plugins := map[string]plugin.Plugin{}
	plugins[id] = m.implementation

	return plugins
}
