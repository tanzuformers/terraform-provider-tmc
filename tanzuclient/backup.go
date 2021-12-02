package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TmcClusterBackup struct {
	FullName *FullName          `json:"fullName"`
	Meta     *MetaData          `json:"meta"`
	Status   *Status            `json:"status"`
	Spec     *ClusterBackupSpec `json:"spec"`
}

type ClusterBackupSpec struct {
	IncludedNamespaces      []string       `json:"includedNamespaces,omitempty"`
	ExcludedNamespaces      []string       `json:"excludedNamespaces,omitempty"`
	IncludedResources       []string       `json:"includedResources,omitempty"`
	ExcludedResources       []string       `json:"excludedResources,omitempty"`
	SnapshotVolumes         bool           `json:"snapshotVolumes"`
	TTL                     string         `json:"ttl"`
	LabelSelector           *LabelSelector `json:"labelSelector,omitempty"`
	IncludeClusterResources bool           `json:"includeClusterResources"`
	StorageLocation         string         `json:"storageLocation"`
	VolumeSnapshotLocations []string       `json:"volumeSnapshotLocations"`
}

type TmcClusterBackupResponse struct {
	Backup TmcClusterBackup `json:"backup"`
}

func (c *Client) GetClusterBackup(name string, mgmt_cluster_name string, cluster_name string, provisioner_name string) (*TmcClusterBackup, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/dataprotection/backups/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, cluster_name, name, mgmt_cluster_name, provisioner_name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := TmcClusterBackupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Backup, nil
}

func (c *Client) DeleteClusterBackup(name string, mgmt_cluster_name string, cluster_name string, provisioner_name string) error {
	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/dataprotection/backups/%s?fullName.managementClusterName=%s&fullName.provisionerName=%s", c.baseURL, cluster_name, name, mgmt_cluster_name, provisioner_name)

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return nil
	}

	res := TmcClusterBackupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil
	}

	return nil
}

func (c *Client) CreateClusterBackup(clusterName string, backup *TmcClusterBackup) (*TmcClusterBackup, error) {

	requestURL := fmt.Sprintf("%s/v1alpha1/clusters/%s/dataprotection/backups", c.baseURL, clusterName)

	newbackupObject := &TmcClusterBackupResponse{
		Backup: *backup,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newbackupObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := TmcClusterBackupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Backup, nil
}
