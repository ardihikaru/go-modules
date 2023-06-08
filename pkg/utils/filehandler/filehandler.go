package filehandler

import "os"

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
