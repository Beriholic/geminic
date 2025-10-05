package dto

import (
	"fmt"
	"reflect"

	"github.com/Beriholic/geminic/internal/config"
	"google.golang.org/genai"
)

type GitCommit struct {
	Typ   string `json:"typ" desc:"type of commit" required:"true"`
	Emoji string `json:"emoji" desc:"emoji of commit" required:"false"`
	Scope string `json:"scope" desc:"scope of commit" required:"false"`
	Msg   string `json:"msg" desc:"msg of commit" required:"true"`
}

func (g GitCommit) String() string {
	if g.Scope != "" {
		if g.Emoji != "" {
			return fmt.Sprintf("%s %s(%s): %s", g.Typ, g.Emoji, g.Scope, g.Msg)
		}
		return fmt.Sprintf("%s(%s): %s", g.Typ, g.Scope, g.Msg)
	}
	if g.Emoji != "" {
		return fmt.Sprintf("%s %s: %s", g.Typ, g.Emoji, g.Msg)
	}
	return fmt.Sprintf("%s: %s", g.Typ, g.Msg)
}

func (g GitCommit) ToGeminiGenerateStruct() *genai.Schema {
	schema := &genai.Schema{
		Type:       genai.TypeObject,
		Properties: map[string]*genai.Schema{},
	}

	useEmoji := config.Get().Emoji

	t := reflect.TypeOf(g)
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		desc := t.Field(i).Tag.Get("desc")

		if name == "emoji" && !useEmoji {
			continue
		}

		schema.Properties[name] = &genai.Schema{
			Type:        genai.TypeString,
			Description: desc,
		}
	}

	return schema
}
