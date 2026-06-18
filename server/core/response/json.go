package response

import (
	"encoding/json"
	"io"
)

func encodeJson(w io.Writer, body any) error{
	err:=json.NewEncoder(w).Encode(body)
	return err
}