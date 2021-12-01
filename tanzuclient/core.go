package tanzuclient

import (
	"fmt"
	"strings"
)

type MetaData struct {
	UID         string                 `json:"uid"`
	Labels      map[string]interface{} `json:"labels,omitempty"`
	Description string                 `json:"description,omitempty"`
}

type FullName struct {
	OrgID                 string `json:"orgId"`
	Name                  string `json:"name"`
	ManagementClusterName string `json:"managementClusterName,omitempty"`
}

type FullNameProvisioned struct {
	FullName        `json:",inline"`
	ProvisionerName string `json:"provisionerName,omitempty"`
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
