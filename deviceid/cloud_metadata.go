package deviceid

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func metaClient() *http.Client {
	return &http.Client{Timeout: 800 * time.Millisecond}
}

func tryAWS() (identity, bool) {
	c := metaClient()
	token, err := awsIMDSv2Token(c)
	if err == nil && token != "" {
		if id, err := awsGet(c, "/latest/meta-data/instance-id", token); err == nil && id != "" {
			return identity{source: SourceAWS, rawID: strings.TrimSpace(id)}, true
		}
	}
	if id, err := awsGet(c, "/latest/meta-data/instance-id", ""); err == nil && id != "" {
		return identity{source: SourceAWS, rawID: strings.TrimSpace(id)}, true
	}
	return identity{}, false
}

func awsIMDSv2Token(c *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodPut, "http://169.254.169.254/latest/api/token", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("aws token status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

func awsGet(c *http.Client, path, token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://169.254.169.254"+path, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("X-aws-ec2-metadata-token", token)
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("aws metadata status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func tryGCP() (identity, bool) {
	c := metaClient()
	req, err := http.NewRequest(http.MethodGet, "http://metadata.google.internal/computeMetadata/v1/instance/id", nil)
	if err != nil {
		return identity{}, false
	}
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := c.Do(req)
	if err != nil {
		return identity{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return identity{}, false
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 256))
	if err != nil || len(b) == 0 {
		return identity{}, false
	}
	return identity{source: SourceGCP, rawID: strings.TrimSpace(string(b))}, true
}

func tryAzure() (identity, bool) {
	c := metaClient()
	url := "http://169.254.169.254/metadata/instance/compute/vmId?api-version=2021-02-01&format=text"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return identity{}, false
	}
	req.Header.Set("Metadata", "true")
	resp, err := c.Do(req)
	if err != nil {
		return identity{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return identity{}, false
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil || len(b) == 0 {
		return identity{}, false
	}
	return identity{source: SourceAzure, rawID: strings.TrimSpace(string(b))}, true
}
