package utils

import "github.com/pquerna/ffjson/ffjson"

// JsonEncodeEncode return json string
func JsonEncode(v interface{}) (string, error) {
	buf, err := ffjson.Marshal(v)
	if err != nil {
		return "", err
	}

	return Byte2Str(buf), nil
}

// JsonDecode json to interface, obj must be a pointer
func JsonDecode(buf []byte, obj interface{}) error {
	return ffjson.Unmarshal(buf, obj)
}
