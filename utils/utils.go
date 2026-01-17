package utils

import (
	"os/user"

	"github.com/deeper-x/chatgo/settings"
)

func GetUsername() string {
	username, err := user.Current()
	if err != nil {
		return settings.DefaultName
	}

	return username.Username
}

func IsAllowed(line string) bool {
	return line[:settings.LenMOTD] == settings.MOTD
}
