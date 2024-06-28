package services

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
	"go-edge-specailist/config"
	"go-edge-specailist/models"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"sort"
	"strings"
)

var sftpInfo *models.SftpInfo

type SftpClientService struct{}

func NewSftpClientService() *SftpClientService {
	return &SftpClientService{}
}

func (s SftpClientService) DownloadLastFileOfSpecialist(envReq string) (*models.SftpInfo, error) {
	getConfig := config.GetConfig()
	var selectedEnv config.SftpConfig

	switch envReq {
	case "dev":
		selectedEnv = getConfig.Sftp.Dev
		break
	case "uat":
		selectedEnv = getConfig.Sftp.Uat
		break
	case "pre":
		selectedEnv = getConfig.Sftp.Pre
		break
	default:
		return nil, fmt.Errorf("Invalid env value: %s", envReq)

	}

	// Establish an SSH connection
	configSftp := &ssh.ClientConfig{
		User: selectedEnv.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(selectedEnv.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", selectedEnv.Host, selectedEnv.Port), configSftp)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to SSH:")
		return nil, err
	}
	defer func(conn *ssh.Client) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// Create an SFTP client
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		log.Error().Err(err).Msg("Error creating SFTP client:")
		return nil, err
	}
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {
			return
		}
	}(sftpClient)

	// List files in the remote directory
	files, err := sftpClient.ReadDir(getConfig.Sftp.RemoteDir)
	if err != nil {
		log.Error().Err(err).Msg("Error listing remote directory:")
		return nil, err
	}

	// Filter files starting with "EIM_EDGE_BLACKLIST"
	filteredFiles := make([]os.FileInfo, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "EIM_EDGE_BLACKLIST") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Sort filtered files by filename
	sort.Slice(filteredFiles, func(i, j int) bool {
		return strings.Compare(filteredFiles[i].Name(), filteredFiles[j].Name()) > 0
	})

	// Download the latest file
	if len(filteredFiles) > 0 {
		latestFile := filteredFiles[0]
		remoteFilePath := getConfig.Sftp.RemoteDir + latestFile.Name()
		localFilePath := getConfig.Sftp.LocalDir + latestFile.Name()

		// Open the remote file for reading
		remoteFile, err := sftpClient.Open(remoteFilePath)
		if err != nil {
			log.Error().Err(err).Msg("Error opening remote file:")
			return nil, err
		}
		defer func(remoteFile *sftp.File) {
			err := remoteFile.Close()
			if err != nil {
				return
			}
		}(remoteFile)

		// Create the local file for writing
		localFile, err := os.Create(localFilePath)
		if err != nil {
			log.Error().Err(err).Msg("Error creating local file:")
			return nil, err
		}
		defer func(localFile *os.File) {
			err := localFile.Close()
			if err != nil {
				return
			}
		}(localFile)

		// Copy the contents from the remote file to the local file
		_, err = io.Copy(localFile, remoteFile)
		if err != nil {
			log.Error().Err(err).Msg("Error copying file contents:")
			return nil, err
		}

		sftpInfoData := models.SftpInfo{FileName: latestFile.Name(),
			Size: latestFile.Size(),
			Path: localFilePath,
		}

		sftpInfo = &sftpInfoData

		log.Info().Msgf("Latest file downloaded successfully from %s to %s\n", remoteFilePath, localFilePath)
	} else {
		log.Info().Msg("No files found in the remote directory starting with 'EIM_EDGE_BLACKLIST'")
	}

	return sftpInfo, nil
}
