package utils

import (
	"fmt"
	"testing"
)

func TestGetCSIPluginType(t *testing.T) {

	pluginType, err := GetCSIPluginType("error_type")
	if err != nil {
		fmt.Println(err.Error())
	}

	pluginType, err = GetCSIPluginType("node")
	if err == nil {
		fmt.Println("Plugin Type is " + pluginType)
	}
}
