package utils

import (
	"errors"
	"goblin-csi-driver/internal"
	"strings"
)

func GetCSIPluginType(t string) (internal.CSIPluginType, error) {

	key := strings.ToLower(t)

	csiPlugin := internal.CSIPluginTypeMap[key]
	if csiPlugin == "" {
		return "", errors.New("unsupport csi plugin type: " + t)
	}
	return csiPlugin, nil

}
