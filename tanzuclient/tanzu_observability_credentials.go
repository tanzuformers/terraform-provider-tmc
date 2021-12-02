package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TmcObservabilityCredential struct {
	FullName *FullName                    `json:"fullName"`
	Meta     *MetaData                    `json:"meta"`
	Spec     *ObservabilityCredentialSpec `json:"spec"`
	Status   struct {
		Phase string `json:"phase,omitempty"`
	} `json:"status,omitempty"`
}

type ObservabilityCredentialSpec struct {
	Capability string                      `json:"capability"`
	Data       ObservabilityCredentialData `json:"data"`
}

type ObservabilityCredentialData struct {
	KeyValue ObservabilityKey `json:"keyValue"`
}

type ObservabilityKey struct {
	Data WaveFrontData `json:"data,omitempty"`
}

type WaveFrontData struct {
	Token string `json:"wavefront.token"`
}

type TmcObservabilityCredentialResponse struct {
	TmcObservabilityCredential TmcObservabilityCredential `json:"credential"`
}

func (c *Client) GetObservabilityCredential(name string) (*TmcObservabilityCredential, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := TmcObservabilityCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcObservabilityCredential, nil
}

// Deletes an already existing Tanzu Observability credential with a given name.
func (c *Client) DeleteObservabilityCredential(name string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, name)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := TmcObservabilityCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateObservabilityCredential(cred *TmcObservabilityCredential) (*TmcObservabilityCredential, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials", c.baseURL)

	newCredObject := &TmcObservabilityCredentialResponse{
		TmcObservabilityCredential: *cred,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newCredObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := TmcObservabilityCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcObservabilityCredential, nil
}

func (c *Client) UpdateObservabilityCredential(cred *TmcObservabilityCredential) (*TmcObservabilityCredential, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, cred.FullName.Name)

	credObject := &TmcObservabilityCredentialResponse{
		TmcObservabilityCredential: *cred,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(credObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := TmcObservabilityCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcObservabilityCredential, nil
}
