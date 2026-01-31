package util

import "file/internal/consts"

func GetFileType(fileExt string) string {

	if fileExt == "jpg" || fileExt == "jpeg" || fileExt == "png" {
		return consts.IMAGE
	} else {
		return consts.FILE
	}
}
