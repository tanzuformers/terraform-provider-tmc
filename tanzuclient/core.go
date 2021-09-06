package tanzuclient

import (
	"fmt"
	"strings"
)

type SimpleMetaData struct {
	UID    string                 `json:"uid"`
	Labels map[string]interface{} `json:"labels,omitempty"`
}

type MetaData struct {
	SimpleMetaData SimpleMetaData
	Description    string `json:"description"`
}

type SimpleFullName struct {
	OrgID string `json:"orgId"`
	Name  string `json:"name"`
}

type FullName struct {
	SimpleFullName        SimpleFullName
	ManagementClusterName string `json:"managementClusterName"`
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
