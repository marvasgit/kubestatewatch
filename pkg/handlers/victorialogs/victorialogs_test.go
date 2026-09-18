package victorialogs

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/marvasgit/kubestatewatch/config"
)

func TestVictoriaLogsInit(t *testing.T) {
	v := &VictoriaLogs{}
	expectedError := fmt.Errorf("missing victoria logs url")

	var Tests = []struct {
		vl   config.VictoriaLogs
		err  error
	}{
		{config.VictoriaLogs{Url: "http://localhost:9428/insert/logline/http"}, nil},
		{config.VictoriaLogs{Url: "http://localhost:9428/insert/logline/http", Cluster: "prod"}, nil},
		{config.VictoriaLogs{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.VictoriaLogs = tt.vl
		if err := v.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): got %v, want %v", err, tt.err)
		}
	}
}