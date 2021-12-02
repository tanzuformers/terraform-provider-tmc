package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NodeName struct {
	OrgID                 string `json:"orgId"`
	ClusterName           string `json:"clusterName"`
	ManagementClusterName string `json:"managementClusterName"`
	ProvisionerName       string `json:"provisionerName"`
	Name                  string `json:"name"`
}

type AwsNodeSpec struct {
	InstanceType     string `json:"instanceType"`
	AvailabilityZone string `json:"availabilityZone"`
	Version          string `json:"version"`
}

type AwsNodePool struct {
	NodeLabels      map[string]interface{} `json:"cloudLabels,omitempty"`
	CloudLabels     map[string]interface{} `json:"nodeLabels,omitempty"`
	WorkerNodeCount string                 `json:"workerNodeCount"`
	NodeTkgAws      AwsNodeSpec            `json:"tkgAws"`
}

type NodePool struct {
	FullName *NodeName    `json:"fullName"`
	Meta     *MetaData    `json:"meta"`
	Spec     *AwsNodePool `json:"spec"`
	Status   *Status      `json:"status"`
}

type NodePoolJsonObject struct {
	NodePool NodePool `json:"nodepool"`
}

func (c *Client) CreateNodePool(name string, managementClusterName string, provisionerName string, clusterName string, description string, cloudLabels map[string]interface{}, nodeLabels map[string]interface{}, nodeCount int, opts *AwsNodeSpec) (*NodePool, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools", c.baseURL, clusterName)

	newNodePool := &NodePool{
		FullName: &NodeName{
			ClusterName:           clusterName,
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
		},
		Spec: &AwsNodePool{
			NodeLabels:      nodeLabels,
			CloudLabels:     cloudLabels,
			WorkerNodeCount: fmt.Sprint(nodeCount),
			NodeTkgAws:      *opts,
		},
	}

	newNodePoolObject := &NodePoolJsonObject{
		NodePool: *newNodePool,
	}

	json_data, err := json.Marshal(newNodePoolObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil

}

func (c *Client) GetNodePool(name string, clusterName string, managementClusterName string, provisionerName string) (*NodePool, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil
}

func (c *Client) UpdateNodePool(name string, managementClusterName string, provisionerName string, clusterName string, description string, cloudLabels map[string]interface{}, nodeLabels map[string]interface{}, nodeCount int, opts *AwsNodeSpec) (*NodePool, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s", c.baseURL, clusterName, name)

	newNodePool := &NodePool{
		FullName: &NodeName{
			ClusterName:           clusterName,
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
		},
		Spec: &AwsNodePool{
			NodeLabels:      nodeLabels,
			CloudLabels:     cloudLabels,
			WorkerNodeCount: fmt.Sprint(nodeCount),
			NodeTkgAws:      *opts,
		},
	}

	newNodePoolObject := &NodePoolJsonObject{
		NodePool: *newNodePool,
	}

	json_data, err := json.Marshal(newNodePoolObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.NodePool, nil
}

func (c *Client) DeleteNodePool(name string, clusterName string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := NodePoolJsonObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) DescribeNodePool(name string, clusterName string, managementClusterName string, provisionerName string) (*Status, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/nodepools/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, clusterName, name, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Token))

	res, err := c.http.Do(req)
	if err != nil {
		return &Status{Phase: "ERROR"}, err
	}

	defer res.Body.Close()

	var target map[string]interface{}

	if err = json.NewDecoder(res.Body).Decode(&target); err != nil {
		return &Status{Phase: "ERROR"}, err
	}

	if target["error"] != nil {
		if strings.Contains(target["error"].(string), "Node Pool Not Found") {
			resp := &Status{
				Phase: "DELETED",
			}
			return resp, err
		}
		return &Status{Phase: "ERROR"}, err
	}

	return &Status{Phase: "DELETING"}, nil
}
