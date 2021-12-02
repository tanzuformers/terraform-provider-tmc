package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ClusterResourceFullName struct {
	FullName
	Provisioner `json:"provisionerName"`
}

type TmcAwsAccountCredential struct {
	FullName *ClusterResourceFullName `json:"fullName"`
	Meta     *MetaData                `json:"meta"`
	Spec     *CredentialSpec          `json:"spec"`
	Status   struct {
		Phase string `json:"phase,omitempty"`
	} `json:"status,omitempty"`
}

type TmcAwsAccountCredentialResponse struct {
	TmcAwsAccountCredential TmcAwsAccountCredential `json:"credential"`
}

func (c *Client) GetAwsAccountCredential(name string, mgmtClusterName string) (*TmcAwsAccountCredential, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/managementcluster/credentials/%s?fullName.managementClusterName=%s", c.baseURL, name, mgmtClusterName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := TmcAwsAccountCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcAwsAccountCredential, nil
}

// Deletes an already existing AWS Data protection credential account with a given name.
func (c *Client) DeleteAwsAccountCredential(name string, mgmtClusterName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/managementcluster/credentials/%s?fullName.managementClusterName=%s", c.baseURL, name, mgmtClusterName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := TmcAwsAccountCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateAwsAccountCredential(cred *TmcAwsAccountCredential) (*TmcAwsAccountCredential, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/account/managementcluster/credentials", c.baseURL)

	newCredObject := &TmcAwsAccountCredentialResponse{
		TmcAwsAccountCredential: *cred,
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

	res := TmcAwsAccountCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcAwsAccountCredential, nil
}
