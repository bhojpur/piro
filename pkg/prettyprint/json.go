package prettyprint

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// JSONFormat formats everything as JSON
const JSONFormat Format = "json"

func formatJSON(pp *Content) error {
	enc := &protojson.MarshalOptions{
		UseEnumNumbers: false,
		Indent:      "  ",
	}
	return enc.Marshal(pp.Writer, pp.Obj)
}
