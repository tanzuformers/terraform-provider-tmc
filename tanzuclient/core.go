package tanzuclient

import (
	"fmt"
	"strings"
)

type FullName struct {
	OrgID                 string `json:"orgId"`
	Name                  string `json:"name"`
	ManagementClusterName string `json:"managementClusterName"`
	ProvisionerName       string `json:"provisionerName,omitempty"`
	ClusterName           string `json:"clusterName,omitempty"`
}

type MetaData struct {
	UID             string                 `json:"uid"`
	Description     string                 `json:"description"`
	Labels          map[string]interface{} `json:"labels,omitempty"`
	ResourceVersion string                 `json:"resourceVersion,omitempty"`
	Annotations     map[string]string      `json:"annotations,omitempty"`
}

type Status struct {
	Phase string `json:"phase,omitempty"`
}

type LabelSelector struct {
	MatchLabels      map[string]interface{} `json:"matchLabels,omitempty"`
	MatchExpressions []MatchExpressions     `json:"matchExpressions,omitempty"`
}

type MatchExpressions struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func buildLabelQuery(labels map[string]interface{}) string {

	var query strings.Builder
	var labelArray []string

	for k, v := range labels {
		newFilter := fmt.Sprintf("meta.labels.%s:%s", k, v)
		labelArray = append(labelArray, newFilter)
	}

	for i, label := range labelArray {
		query.WriteString(label)
		if i == len(labelArray)-1 {
			break
		}
		query.WriteString(" and ")
	}

	return query.String()
}
