package client

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/petarov/query-apple-firmware-updates/config"
)

const INFO_URL = "https://api.ipsw.me/v2.1/%s/latest/info.json"

type IPSWInfo struct {
	Identifier  string `json:"identifier"`
	Version     string `json:"version"`
	Device      string `json:"device"`
	BuildId     string `json:"buildid"`
	SHA1Sum     string `json:"sha1sum"`
	MD5Sum      string `json:"md5sum"`
	Size        int64  `json:"size"`
	ReleaseDate string `json:"releasedate"`
	UploadDate  string `json:"uploaddate"`
	Url         string `json:"url"`
	Signed      bool   `json:"signed"`
	Filename    string `json:"filename"`
}

func NewIPSWClient() (client *http.Client, err error) {
	tlsConfig := &tls.Config{}
	client = &http.Client{Timeout: 20 * time.Second}
	client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	return client, nil
}

func getResponseBody(resp *http.Response) io.ReadCloser {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	return reader
}

func IPSWGetInfo(client *http.Client, product string) (ipsw []*IPSWInfo, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(INFO_URL, product), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("User-Agent", fmt.Sprintf("qados-v%s", config.VERSION))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 201 {
		read, _ := io.ReadAll(getResponseBody(resp))
		return nil, fmt.Errorf("HTTP error when requesting update info: (%d) %s", resp.StatusCode, string(read))
	}

	ipsw = make([]*IPSWInfo, 0)
	if err := json.NewDecoder(getResponseBody(resp)).Decode(&ipsw); err != nil {
		return nil, fmt.Errorf("Error decoding ipsw json body: %w", err)
	}

	return ipsw, nil
}
