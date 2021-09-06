package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Provisioner struct {
	// The name of the provisioner.
	FullName *FullName `json:"fullName"`
	// The metadata of the provisioner.
	Meta *SimpleMetaData `json:"meta"`
}

type ProvisionerResponse struct {
	Provisioner Provisioner `json:"provisioner"`
}

type AllProvisioners struct {
	Provisioners []Provisioner `json:"provisioners"`
}

func (c *Client) GetProvisioner(mgmtClusterName, name string) (*Provisioner, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s/provisioners/%s", c.baseURL, mgmtClusterName, name)

	req, err := http.NewRequest("GET", tmcURL, nil)
	if err != nil {
		return nil, err
	}

	res := ProvisionerResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Provisioner, nil
}

func (c *Client) GetAllProvisioners(mgmtClusterName string, labels map[string]interface{}) ([]Provisioner, error) {
	queryString := buildLabelQuery(labels)

	tmcURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s/provisioners?query=%s", c.baseURL, mgmtClusterName, queryString)

	req, err := http.NewRequest("GET", tmcURL, nil)
	if err != nil {
		return nil, err
	}

	res := AllProvisioners{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res.Provisioners, nil
}

func (c *Client) CreateProvisioner(mgmtClusterName string, name string, description string, labels map[string]interface{}) (*Provisioner, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s/provisioners", c.baseURL, mgmtClusterName)

	provisioner := &Provisioner{
		FullName: &FullName{
			SimpleFullName: &SimpleFullName{
				Name: name,
			},
			ManagementClusterName: mgmtClusterName,
		},
		Meta: &SimpleMetaData{
			Labels: labels,
		},
	}

	provisionerResponse := &ProvisionerResponse{
		Provisioner: *provisioner,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(provisionerResponse) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tmcURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ProvisionerResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Provisioner, nil
}

func (c *Client) DeleteProvisioner(mgmtClusterName, name string) error {
	tmcURL := fmt.Sprintf("%s/v1alpha1/managementclusters/%s/provisioners/%s", c.baseURL, mgmtClusterName, name)

	req, err := http.NewRequest("DELETE", tmcURL, nil)
	if err != nil {
		return err
	}

	res := ProvisionerResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
