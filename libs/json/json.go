package json

import (
	"encoding/json"
	"fmt"
)

func Prettify(data interface{}) ([]byte, error) {
	bytes, err := json.MarshalIndent(data, "  ", "  ")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func PrettifyStr(data interface{}) (string, error) {
	bytes, err := Prettify(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func PrettyPrint(data interface{}) error {
	prettifiedData, err := PrettifyStr(data)
	if err != nil {
		return err
	}
	fmt.Println(prettifiedData)
	return nil
}

func Print(data interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal message: %w", err)
	}

	fmt.Printf("%v\n", string(buf))
	
	return nil
}
