package tanzuclient

type MetaData struct {
	UID         string `json:"uid"`
	Description string `json:"description"`
}

type FullName struct {
	OrgID string `json:"orgId"`
	Name  string `json:"name"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
