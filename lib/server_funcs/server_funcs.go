package serverfuncs

import (
	"kato-studio/katoengine/lib/utils"

	"github.com/tidwall/gjson"
)

func FuncsSwitch(func_string string) string {

	utils.Debug(func_string)

	return "server_func_string"
}

func Fetch(func_string string) gjson.Result {

	utils.Debug(func_string)

	return gjson.Get(func_string, "data")
}