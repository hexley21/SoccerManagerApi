package json

import (
	"io"
)

type Processor interface {
	Serializer
	SerializeWriter
	Deserializer
	DeserializeWriter
}

type Serializer interface {
	Serialize(v interface{}) ([]byte, error)
}

type SerializeWriter interface {
	SerializeWriter(writer io.Writer, i any) error
}

type Deserializer interface {
	Deserialize(data []byte, v interface{}) error
}

type DeserializeWriter interface {
	DeserializeReader(reader io.Reader, obj interface{}) error
}
