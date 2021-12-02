package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type CredentialMetaData struct {
	Provider string `json:"provider,omitempty"`
}

type CredentialData struct {
	AwsCredential *AwsCredential    `json:"awsCredential"`
	KeyValue      *AwsCredentialKey `json:"keyValue,omitempty"`
}

type AwsCredentialKey struct {
	Type string        `json:"type"`
	Data *AwsAccessKey `json:"data"`
}

type AwsCredential struct {
	AccountID string   `json:"accountId,omitempty"`
	IamRole   *IamRole `json:"iamRole,omitempty"`
}

type AwsAccessKey struct {
	AccessKeyId     string `json:"aws_access_key_id,omitempty"`
	SecretAccessKey string `json:"aws_secret_access_key,omitempty"`
}

type IamRole struct {
	Arn string `json:"arn"`
}

type CredentialSpec struct {
	MetaData   *CredentialMetaData `json:"meta"`
	Capability string              `json:"capability"`
	Data       *CredentialData     `json:"data"`
}

type TmcAwsCredential struct {
	FullName *FullName       `json:"fullName"`
	Meta     *MetaData       `json:"meta"`
	Spec     *CredentialSpec `json:"spec"`
	Status   struct {
		Phase string `json:"phase,omitempty"`
	} `json:"status,omitempty"`
}

type TmcAwsCredentialResponse struct {
	TmcAwsCredential TmcAwsCredential `json:"credential"`
}

func (c *Client) GetAwsCredential(name string) (*TmcAwsCredential, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := TmcAwsCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcAwsCredential, nil
}

// Deletes an already existing AWS Data protection credential account with a given name.
func (c *Client) DeleteAwsCredential(name string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials/%s", c.baseURL, name)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := TmcAwsCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateAwsCredential(cred *TmcAwsCredential) (*TmcAwsCredential, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/account/credentials", c.baseURL)

	newCredObject := &TmcAwsCredentialResponse{
		TmcAwsCredential: *cred,
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

	res := TmcAwsCredentialResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.TmcAwsCredential, nil
}
