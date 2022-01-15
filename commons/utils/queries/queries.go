package queries

type ByTags struct {
	Tags  []string `json:"tags,omitempty"`
	Limit int      `json:"limit,omitempty"`
	Page  int      `json:"page,omitempty"`
}

type ByLabels struct {
	Labels map[string]string `json:"labels,omitempty"`
	Limit  int               `json:"limit,omitempty"`
	Page   int               `json:"page,omitempty"`
}
