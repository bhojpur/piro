package prettyprint

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONFormat formats everything as JSON
const JSONFormat Format = "json"

func formatJSON(pp *Content) error {
	ProtobufToJSON(pp.Obj)
	return nil
}

func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		Indent:          "  ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	b, err := marshaler.Marshal(message)
	return string(b), err
}
