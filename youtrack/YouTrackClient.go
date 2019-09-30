package youtrack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LoadError struct {
}

func (err LoadError) Error() string {
	return "Load failed"
}

func LoadVersions() (*VersionBundle, error) {
	resp, err := http.Get("https://jinya.myjetbrains.com/youtrack/api/admin/customFieldSettings/bundles/version/71-2?fields=values(id,name)")
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

func LoadIssueTypes() ([]string, error) {
	requestUrl := "https://jinya.myjetbrains.com/youtrack/api/admin/customFieldSettings/bundles/enum/66-7/values?$includeArchived=false&fields=$type,name"
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	var types []issueType
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &types)
	if err != nil {
		return nil, err
	}

	resultTypes := make([]string, len(types))
	for idx, _type := range types {
		resultTypes[idx] = _type.Name
	}

	return resultTypes, nil
}

func LoadIssues(version string, issueType string) ([]Issue, error) {
	requestUrl := "https://jinya.myjetbrains.com/youtrack/api/issues?fields=summary%2CidReadable&query=project%3AJGCMS+Fix+versions:" + version + "+Type:%22" + url.QueryEscape(issueType) + "%22"
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, LoadError{}
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
