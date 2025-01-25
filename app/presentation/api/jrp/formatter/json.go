package formatter

import (
	"errors"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

// JsonFormatter is a struct that formats the output of jrp server.
type JsonFormatter struct{}

// ResponseOutputDto is a struct that represents the response of jrp server.
type ResponseOutputDto struct {
	Body []byte `json:"body"`
}

// JrpJsonOutputDto is a struct that represents the output json of jrp server.
type JrpJsonOutputDto struct {
	Phrase string `json:"phrase"`
}

var (
	// Ju is a variable that contains the JsonUtil struct for injecting dependencies in testing.
	Ju = utility.NewJsonUtil(proxy.NewJson())
)

// NewJsonFormatter returns a new instance of the JsonFormatter struct.
func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

// Format formats the output of jrp server.
func (f *JsonFormatter) Format(result interface{}) ([]byte, error) {
	var formatted []byte
	var err error
	switch v := result.(type) {
	case *jrpApp.GenerateJrpUseCaseOutputDto:
		jjoDto := JrpJsonOutputDto{Phrase: v.Phrase}
		gjj, err := Ju.Marshal(jjoDto)
		if err != nil {
			return nil, err
		}
		gjroDto := ResponseOutputDto{Body: gjj}
		formatted = gjroDto.Body
	default:
		formatted = nil
		err = errors.New("invalid result type")
	}
	return formatted, err
}
