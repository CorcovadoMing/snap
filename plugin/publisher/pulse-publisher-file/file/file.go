package file

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"time"
	// "fmt"
	"os"
	// "strings"
	// "time"

	"github.com/intelsdilabs/pulse/control/plugin"
	"github.com/intelsdilabs/pulse/control/plugin/cpolicy"
	"github.com/intelsdilabs/pulse/core/ctypes"
)

const (
	name       = "file"
	version    = 1
	pluginType = plugin.PublisherPluginType
)

type filePublisher struct {
}

func NewFilePublisher() *filePublisher {
	return &filePublisher{}
}

func (f *filePublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue, logger *log.Logger) error {
	logger.Println("Publishing started")
	var metrics []plugin.PluginMetricType

	switch contentType {
	case plugin.ContentTypes[plugin.PulseGobContentType]:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return errors.New(fmt.Sprintf("Unknown content type '%s'", contentType))
	}

	logger.Printf("publishing %v to %v", metrics, config)
	file, err := os.OpenFile(config["file"].(ctypes.ConfigValueStr).Value, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		logger.Printf("Error: %v", err)
		return err
	}
	w := bufio.NewWriter(file)
	for _, m := range metrics {
		w.WriteString(fmt.Sprintf("%v|%v|%v\n", time.Now().Local(), m.Namespace(), m.Data()))
	}
	w.Flush()

	return nil
}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType)
}

func (f *filePublisher) GetConfigPolicyNode() cpolicy.ConfigPolicyNode {
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule("file", true)
	handleErr(err)
	r1.Description = "Absolute path to the output file for publishing"

	config.Add(r1)
	return *config
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}