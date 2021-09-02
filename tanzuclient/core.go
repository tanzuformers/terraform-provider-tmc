package tanzuclient

import (
	"fmt"
	"strings"
)

type MetaData struct {
	UID         string                 `json:"uid"`
	Description string                 `json:"description"`
	Labels      map[string]interface{} `json:"labels,omitempty"`
}

type FullName struct {
	OrgID string `json:"orgId"`
	Name  string `json:"name"`
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
