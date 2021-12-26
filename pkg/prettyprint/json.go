package prettyprint

import "google.golang.org/protobuf/encoding/protojson"

// JSONFormat formats everythign as JSON
const JSONFormat Format = "json"

func formatJSON(pp *Content) error {
	enc := &jsonpb.Marshaler{
		EnumsAsInts: false,
		Indent:      "  ",
	}
	return enc.Marshal(pp.Writer, pp.Obj)
}
