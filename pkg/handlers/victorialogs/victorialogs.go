package victorialogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/marvasgit/kubestatewatch/config"
	"github.com/marvasgit/kubestatewatch/pkg/event"
	"github.com/marvasgit/kubestatewatch/pkg/message"
	"github.com/sirupsen/logrus"
)

type VictoriaLogs struct {
	Url     string
	Cluster string
	Title   string
}

type victoriaLogsMessage struct {
	Msg      string `json:"_msg"`
	Cluster  string `json:"cluster"`
	Title    string `json:"title"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	NS       string `json:"namespace"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
	APIVersion string `json:"apiversion"`
	Diff     string `json:"diff"`
}

func (v *VictoriaLogs) Init(c *config.Config) error {
	v.Url = c.Handler.VictoriaLogs.Url
	v.Cluster = c.Handler.VictoriaLogs.Cluster
	v.Title = message.GetTitle(c.Message.Title, "KW_VICTORIALOGS_TITLE")

	if v.Url == "" {
		return fmt.Errorf("missing victoria logs url")
	}

	return nil
}

func (v *VictoriaLogs) Handle(e event.StatemonitorEvent) {
	msg := victoriaLogsMessage{
		Msg:       e.Message(),
		Cluster:   v.Cluster,
		Title:     v.Title,
		Kind:      e.Kind,
		Name:      e.Name,
		NS:        e.Namespace,
		Reason:    e.Reason,
		Status:    e.Status,
		APIVersion: e.ApiVersion,
		Diff:      e.DiffMarshalled,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("failed to marshal victoria logs message: %v", err)
		return
	}

	req, err := http.NewRequest("POST", v.Url, bytes.NewBuffer(data))
	if err != nil {
		logrus.Errorf("failed to create victoria logs request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("failed to send victoria logs message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logrus.Errorf("victoria logs returned error status: %d", resp.StatusCode)
		return
	}

	logrus.Debugf("Message successfully sent to victoria logs at %s", time.Now())
}