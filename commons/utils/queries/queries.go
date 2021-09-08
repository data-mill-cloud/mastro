package queries

type ByTags struct {
	Tags []string `json:"tags,omitempty"`
}

type ByLabels struct {
	Labels map[string]string `json:"labels,omitempty"`
}
