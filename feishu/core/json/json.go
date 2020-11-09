package json

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//ToJSON 转换成 json 格式
func toJSON(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

//ToJSONIgnoreError 转换为 JSON 格式，并忽略报错
func ToJSONIgnoreError(obj interface{}) string {
	jsonStr, _ := toJSON(obj)
	return jsonStr
}

//FromJSON a
func fromJSON(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

//FromJSONIgnoreError a
func FromJSONIgnoreError(jsonStr string, obj interface{}) {
	_ = fromJSON(jsonStr, obj)
}
