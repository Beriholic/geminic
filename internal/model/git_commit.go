package model

import "fmt"

type GitCommit struct {
	Typ   string `json:"typ,omitempty"`
	Emoji string `json:"emoji,omitempty"`
	Scope string `json:"scope,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

func (g GitCommit) String() string {
	if g.Scope != "" {
		if g.Emoji != "" {
			return fmt.Sprintf("%s %s(%s): %s", g.Emoji, g.Typ, g.Scope, g.Msg)
		}
		return fmt.Sprintf("%s(%s): %s", g.Typ, g.Scope, g.Msg)
	}
	if g.Emoji != "" {
		return fmt.Sprintf("%s %s: %s", g.Emoji, g.Typ, g.Msg)
	}
	return fmt.Sprintf("%s: %s", g.Typ, g.Msg)
}
