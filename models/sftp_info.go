package models

type SftpInfo struct {
	FileName string `json:"fileName"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}
