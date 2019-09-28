package youtrack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func LoadVersions() (*VersionBundle, error) {
	resp, err := http.Get("https://jinya.myjetbrains.com/youtrack/api/admin/customFieldSettings/bundles/version/71-0?fields=values(id,name)")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionBundle VersionBundle
	err = json.Unmarshal(body, &versionBundle)
	if err != nil {
		return nil, err
	}

	return &versionBundle, nil

}

func LoadIssues(version string) ([]Issue, error) {
	requestUrl := "https://jinya.myjetbrains.com/youtrack/api/issues?fields=summary%2CidReadable&query=project%3AJGCMS+Fix+versions:" + version
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
