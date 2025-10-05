package dto

type CommitDTO struct {
	Commit string   `json:"commit,omitempty"`
	Diff   string   `json:"diff,omitempty"`
	Files  []string `json:"files,omitempty"`
}
