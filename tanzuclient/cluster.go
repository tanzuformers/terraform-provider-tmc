package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Subnet struct {
	Id       string `json:"id"`
	IsPublic bool   `json:"isPublic"`
}

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
			Id        string `json:"id,omitempty"`
			CidrBlock string `json:"cidrBlock"`
		} `json:"vpc"`
		Subnets []Subnet `json:"subnets"`
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
			HighAvailability  bool     `json:"highAvailability,omitempty"`
			InstanceType      string   `json:"instanceType"`
		} `json:"controlPlane"`
	} `json:"topology"`
}

type ClusterSpec struct {
	ClusterGroupName string     `json:"clusterGroupName"`
	TkgAws           AWSCluster `json:"tkgAws,omitempty"`
}

type ClusterStatus struct {
	Status        `json:",inline"`
	InstallerLink string `json:"installerLink"`
}

type Cluster struct {
	FullName *FullName      `json:"fullName"`
	Meta     *MetaData      `json:"meta"`
	Spec     *ClusterSpec   `json:"spec"`
	Status   *ClusterStatus `json:"status"`
}

type ClusterJSONObject struct {
	Cluster Cluster `json:"cluster"`
}

// Options interface for passing arguments to the
// functions neccessary to perform on the Cluster
type ClusterOpts struct {
	Region           string
	Version          string
	CredentialName   string
	ControlPlaneSpec map[string]interface{}
	PodCidrBlock     string
	ServiceCidrBlock string
	SshKey           string
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

func (c *Client) CreateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, labels map[string]interface{}, opts *ClusterOpts) (*Cluster, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters", c.baseURL)

	awsSpec := buildAwsJsonObject(opts)

	newCluster := &Cluster{
		FullName: &FullName{
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: cluster_group,
			TkgAws:           awsSpec,
		},
	}

	newClusterObject := &ClusterJSONObject{
		Cluster: *newCluster,
	}

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

func (c *Client) DeleteCluster(name string, managementClusterName string, provisionerName string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

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

func (c *Client) UpdateCluster(name string, managementClusterName string, provisionerName string, cluster_group string, description string, resourceVersion string, labels map[string]interface{}, opts *ClusterOpts) (*Cluster, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, name, managementClusterName, provisionerName)

	awsSpec := buildAwsJsonObject(opts)

	newCluster := &Cluster{
		FullName: &FullName{
			Name:                  name,
			ManagementClusterName: managementClusterName,
			ProvisionerName:       provisionerName,
		},
		Meta: &MetaData{
			ResourceVersion: resourceVersion,
			Description:     description,
			Labels:          labels,
		},
		Spec: &ClusterSpec{
			ClusterGroupName: cluster_group,
			TkgAws:           awsSpec,
		},
	}

	newClusterObject := &ClusterJSONObject{
		Cluster: *newCluster,
	}

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

func (c *Client) DescribeCluster(fullName string, managementClusterName string, provisionerName string) (*Status, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, fullName, managementClusterName, provisionerName)

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
		if strings.Contains(target["error"].(string), fmt.Sprintf("Resource not found - cluster(%s)", fullName)) {
			resp := &Status{
				Phase: "DELETED",
			}
			return resp, err
		}
		return &Status{Phase: "ERROR"}, err
	}

	return &Status{Phase: "DELETING"}, nil
}

func buildAwsJsonObject(opts *ClusterOpts) AWSCluster {

	var newAwsSpec AWSCluster

	newAwsSpec.Distribution.ProvisionerCredentialName = opts.CredentialName
	newAwsSpec.Distribution.Region = opts.Region
	newAwsSpec.Distribution.Version = opts.Version

	newAwsSpec.Settings.Network.ClusterNetwork.Pods = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Pods[0].CidrBlocks = opts.PodCidrBlock
	newAwsSpec.Settings.Network.ClusterNetwork.Services = make([]struct {
		CidrBlocks string "json:\"cidrBlocks\""
	}, 1)
	newAwsSpec.Settings.Network.ClusterNetwork.Services[0].CidrBlocks = opts.ServiceCidrBlock

	newAwsSpec.Settings.Security.SshKey = opts.SshKey

	cp_spec := opts.ControlPlaneSpec

	newAwsSpec.Topology.ControlPlane.InstanceType = cp_spec["instance_type"].(string)
	var azList []string
	for i := 0; i < len((cp_spec["availability_zones"]).([]interface{})); i++ {
		azList = append(azList, (cp_spec["availability_zones"]).([]interface{})[i].(string))
	}

	newAwsSpec.Topology.ControlPlane.AvailabilityZones = azList
	if len(azList) > 1 {
		newAwsSpec.Topology.ControlPlane.HighAvailability = true
	}

	if cp_spec["vpc_id"] != nil {
		newAwsSpec.Settings.Network.Provider.Vpc.Id = cp_spec["vpc_id"].(string)

		pvt_subnets := cp_spec["private_subnets"].([]interface{})
		pub_subnets := cp_spec["public_subnets"].([]interface{})

		newAwsSpec.Settings.Network.Provider.Subnets = make([]Subnet, 0)

		for i := 0; i < len(pvt_subnets); i++ {
			public_subnet_map := &Subnet{
				Id:       pub_subnets[i].(string),
				IsPublic: true,
			}
			private_subnet_map := &Subnet{
				Id:       pvt_subnets[i].(string),
				IsPublic: false,
			}
			newAwsSpec.Settings.Network.Provider.Subnets = append(newAwsSpec.Settings.Network.Provider.Subnets, *private_subnet_map, *public_subnet_map)
		}

	} else {
		newAwsSpec.Settings.Network.Provider.Vpc.CidrBlock = cp_spec["vpc_cidrblock"].(string)
	}

	return newAwsSpec
}
