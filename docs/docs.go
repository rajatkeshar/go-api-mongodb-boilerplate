package docs

import (
    "io/ioutil"
	"github.com/swaggo/swag"
)

type s struct{}

func (s *s) ReadDoc() string {
    doc, _ := ioutil.ReadFile("docs/swagger/swagger.json")
	return string(doc)
}
func init() {
	swag.Register(swag.Name, &s{})
}
