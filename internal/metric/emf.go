package metric

import (
	"encoding/json"
	"os"
	"time"
)

const namespace = "PolicyInferenceDecider"

type (
	emfPayload struct {
		AWS     emfAWS `json:"_aws"`
		Success int    `json:"success"`
	}

	emfAWS struct {
		Timestamp         int64    `json:"Timestamp"`
		CloudWatchMetrics []emfDef `json:"CloudWatchMetrics"`
	}

	emfDef struct {
		Namespace  string      `json:"Namespace"`
		Dimensions [][]string  `json:"Dimensions"`
		Metrics    []emfMetric `json:"Metrics"`
	}

	emfMetric struct {
		Name string `json:"Name"`
		Unit string `json:"Unit"`
	}
)

func EmitSuccess() {
	payload := emfPayload{
		AWS: emfAWS{
			Timestamp: time.Now().UnixMilli(),
			CloudWatchMetrics: []emfDef{{
				Namespace:  namespace,
				Dimensions: [][]string{},
				Metrics:    []emfMetric{{Name: "success", Unit: "Count"}},
			}},
		},
		Success: 1,
	}
	data, _ := json.Marshal(payload)
	os.Stdout.Write(data)
	os.Stdout.Write([]byte("\n"))
}
