package printer

import (
	"encoding/json"
	"fmt"
	"io"
)

func FprintJSON(w io.Writer, data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal message: %w", err)
	}

	if _, err = fmt.Fprintf(w, "%v\n", string(buf)); err != nil {
		return fmt.Errorf("couldn't print data to %v: %w", w, err)
	}

	return nil
}
