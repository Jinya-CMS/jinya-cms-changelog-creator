package youtrack

type VersionBundleElement struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"$type"`
}

type VersionBundle struct {
	Values []VersionBundleElement `json:"values"`
	Type   string                 `json:"$type"`
}

type Issue struct {
	Summary string `json:"summary"`
	Id      string `json:"idReadable"`
	Type    string `json:"$type"`
}

type issueType struct {
	Name string `json:"name"`
	Type string `json:"$type"`
}
