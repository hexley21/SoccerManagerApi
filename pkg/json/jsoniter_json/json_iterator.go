package jsoniter_json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type jsonIterSerializer struct{}

func New() *jsonIterSerializer {
	return &jsonIterSerializer{}
}

func (j *jsonIterSerializer) Serialize(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (j *jsonIterSerializer) Deserialize(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func (j *jsonIterSerializer) SerializeWriter(writer io.Writer, i interface{}) error {
	return jsoniter.NewEncoder(writer).Encode(i)
}

func (j *jsonIterSerializer) DeserializeReader(reader io.Reader, obj interface{}) error {
	return jsoniter.NewDecoder(reader).Decode(obj)
}

type jsonIterEchoAdapter struct {
	*jsonIterSerializer
}

func NewEcho(j *jsonIterSerializer) *jsonIterEchoAdapter {
	return &jsonIterEchoAdapter{
		jsonIterSerializer: j,
	}
}

func (j *jsonIterEchoAdapter) Serialize(c echo.Context, i interface{}, indent string) error {
	return j.jsonIterSerializer.SerializeWriter(c.Response(), i)
}

func (j *jsonIterEchoAdapter) Deserialize(c echo.Context, i interface{}) error {
	return j.jsonIterSerializer.DeserializeReader(c.Request().Body, i)
}
