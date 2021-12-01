package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Network struct {
	ClusterNetwork struct {
		Pods []struct {
			CidrBlocks string `json:"cidrBlocks"`
		} `json:"pods"`
		Services []struct {
			CidrBlocks string `json:"cidrBlocks"`
		} `json:"services"`
	} `json:"cluster"`
	Provider struct {
		Vpc struct {
			CidrBlock string `json:"cidrBlock"`
		} `json:"vpc"`
	} `json:"provider"`
}

type AWSCluster struct {
	Distribution struct {
		ProvisionerCredentialName string `json:"provisionerCredentialName"`
		Region                    string `json:"region"`
		Version                   string `json:"version"`
	} `json:"distribution"`
	Settings struct {
		Network  Network `json:"network"`
		Security struct {
			SshKey string `json:"sshKey"`
		} `json:"security"`
	} `json:"settings"`
	Topology struct {
		ControlPlane struct {
			AvailabilityZones []string `json:"availabilityZones"`
			InstanceType      string   `json:"instanceType"`
		} `json:"controlPlane"`
	} `json:"topology"`
}

type ClusterSpec struct {
	ClusterGroupName string `json:"clusterGroupName"`
	//TkgAws           AWSCluster `json:"tkgAws,omitempty"`
}

type ClusterStatus struct {
	InstallerLink string `json:"installerLink"`
}

type Cluster struct {
	FullName *FullNameProvisioned `json:"fullName"`
	Meta     *MetaData            `json:"meta"`
	Spec     *ClusterSpec         `json:"spec"`
	Status   *ClusterStatus       `json:"status"`
}

type ClusterJSONObject struct {
	Cluster Cluster `json:"cluster"`
}

func (c *Client) GetCluster(fullName string, managementClusterName string, provisionerName string) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, fullName, managementClusterName, provisionerName)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) CreateCluster(name string, description string, managementCluster string, provisionerName string, clusterGroupName string, labels map[string]interface{}) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters", c.baseURL)

	newCluster := &Cluster{
		FullName: &FullNameProvisioned{
			FullName: FullName{
				Name:                  name,
				ManagementClusterName: managementCluster,
			},
			ProvisionerName: provisionerName,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: clusterGroupName,
		},
	}

	newClusterObject := ClusterJSONObject{
		Cluster: *newCluster,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) UpdateCluster(name string, description string, managementCluster string, provisionerName string, clusterGroupName string, labels map[string]interface{}) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s", c.baseURL, name)

	newCluster := &Cluster{
		FullName: &FullNameProvisioned{
			FullName: FullName{
				Name:                  name,
				ManagementClusterName: managementCluster,
			},
			ProvisionerName: provisionerName,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: clusterGroupName,
		},
	}

	newClusterObject := ClusterJSONObject{
		Cluster: *newCluster,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newClusterObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Cluster, nil
}

func (c *Client) DeleteCluster(name string, managementCluster string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementCluster, provisionerName)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := ClusterJSONObject{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}
