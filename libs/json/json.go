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
