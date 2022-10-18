package plugin

import (
	"encoding/json"
	"fmt"
	plugin "github.com/hashicorp/go-plugin"
	botkubeexecutorplugin "github.com/huseyinbabal/botkube-plugins/api/executor"
	botkubesourceplugin "github.com/huseyinbabal/botkube-plugins/api/source"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Type string

const (
	TypeSource   Type = "source"
	TypeExecutor Type = "executor"

	PluginIndexUrl             = "https://raw.githubusercontent.com/huseyinbabal/botkube-plugins/main/index.json"
	PluginExecutableUrlPattern = "https://github.com/huseyinbabal/botkube-plugins/releases/download/%s/%s/%s"
)

type IndexInfo struct {
	Name        string
	Type        Type
	Description string
	Version     string
}

type Metadata struct {
	Name    string
	Path    string
	Client  *plugin.Client
	Type    Type
	Version string
}

type Manager struct {
	PluginsFolder      string
	Plugins            []Metadata
	http               *http.Client
	fileDownloadClient *http.Client
}

func (m *Manager) Initialize(plugins []string) error {
	err := m.RefreshPluginIndex()
	if err != nil {
		return err
	}

	for _, p := range plugins {
		pl := m.GetPlugin(p)
		if pl == nil {
			log.Warnf("Plugin: %s not found.", p)
		} else {
			err := m.Download(pl.Name, pl.Version)
			if err != nil {
				log.Warnf("Failed to download plugin: %s", p)
			}
		}
	}
	return nil
}

func (m *Manager) Download(name, version string) error {
	pluginFile := filepath.Join(m.PluginsFolder, name)
	if _, err := os.Stat(pluginFile); err == nil {
		log.Infof("plugin %s already exists, skipping download.", name)
		return nil
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", m.PluginsFolder, name))
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf(PluginExecutableUrlPattern, name, version, name)
	resp, err := m.fileDownloadClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	err = os.Chmod(pluginFile, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (m *Manager) RefreshPluginIndex() error {
	resp, err := http.Get(PluginIndexUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result []IndexInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}
	for _, index := range result {
		m.Plugins = append(m.Plugins, Metadata{
			Name:    index.Name,
			Type:    index.Type,
			Path:    fmt.Sprintf("%s/%s", m.PluginsFolder, index.Name),
			Version: index.Version,
		})
	}
	return nil
}

func (m *Manager) GetPlugin(name string) *Metadata {
	for _, p := range m.Plugins {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

func NewManager(pluginsFolder string) *Manager {
	if _, err := os.Stat(pluginsFolder); os.IsNotExist(err) {
		err := os.Mkdir(pluginsFolder, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &Manager{
		Plugins:       []Metadata{},
		PluginsFolder: pluginsFolder,
		http:          &http.Client{Timeout: time.Duration(1) * time.Second},
		fileDownloadClient: &http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		},
	}
}

func (m *Manager) Start() error {
	for i, metadata := range m.Plugins {
		log.Printf("Registering plugin %s[id=%s]", metadata.Name, metadata.Path)
		client := plugin.NewClient(&plugin.ClientConfig{
			Plugins:          m.pluginMap(metadata),
			VersionedPlugins: nil,
			Cmd:              exec.Command(metadata.Path), // Plugin specific params goes here
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "BOTKUBE_MAGIC_COOKIE",
				MagicCookieValue: "BOTKUBE_BASIC_PLUGIN",
			},
		})
		m.Plugins[i].Client = client
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

func (m *Manager) GetAdapter(name string) (interface{}, error) {

	pl := m.GetPlugin(name)
	if pl == nil {
		return nil, fmt.Errorf("plugin: %s not found in registered plugins", name)
	}

	client := pl.Client

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (m *Manager) pluginMap(metadata Metadata) map[string]plugin.Plugin {
	plugins := map[string]plugin.Plugin{}
	if metadata.Type == TypeExecutor {
		plugins[metadata.Name] = &botkubeexecutorplugin.ExecutorPlugin{}
	} else if metadata.Type == TypeSource {
		plugins[metadata.Name] = &botkubesourceplugin.SourcePlugin{}

	}

	return plugins
}
