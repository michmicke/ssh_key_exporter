package metrics

import (
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/michmicke/ssh_key_exporter/internal/config"
	"github.com/michmicke/ssh_key_exporter/internal/ssh"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	KeyCount *prometheus.GaugeVec
	Keys     *prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		KeyCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ssh_key_count",
			Help: "Number of authorized SSH keys",
		}, []string{"path"}),
		Keys: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "ssh_keys",
			Help: "Number of authorized SSH keys",
		}, []string{"path", "keytype", "fingerprint", "comment"}),
	}

	reg.MustRegister(m.KeyCount)
	reg.MustRegister(m.Keys)
	return m
}

func (m *Metrics) Extract(path string) {
	keys, err := ssh.ParseAuthorizedKeysFile(path)
	if err != nil {
		log.Printf("Error parsing %s: %v", path, err)
		return
	}
	m.KeyCount.WithLabelValues(path).Set(float64(len(keys)))
	for _, key := range keys {
		m.Keys.WithLabelValues(path, key.Keytype, key.Fingerprint, key.Comment).Set(1)
	}
}

func (m *Metrics) WatchAuthorizedKeys(c chan bool) {
	var paths = []string{path.Join(config.GetConfig().HostPath, "/root/.ssh/authorized_keys")}

	homePath := path.Join(config.GetConfig().HostPath, "home")
	homeDirs, err := os.ReadDir(homePath)
	if err != nil {
		log.Printf("Error reading /home directory: %v", err)
		return
	}
	for _, homeDir := range homeDirs {
		path := path.Join(homePath, homeDir.Name(), ".ssh/authorized_keys")
		log.Printf("Adding %v to the watched paths", path)
		paths = append(paths, path)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating watcher: %v", err)
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					m.Extract(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Error watching files: %v", err)
			}
		}
	}()

	log.Printf("Watching authorized_keys files: %s", paths)
	for _, path := range paths {
		m.Extract(path)
		if err := watcher.Add(path); err != nil {
			log.Printf("Error watching %s: %v", path, err)
		}
	}

	for range c {
	}
}
