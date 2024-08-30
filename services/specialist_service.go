package services

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
	"go-edge-specailist/config"
	"go-edge-specailist/models"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var dataSpecialist *[]models.SpecialistRequest

type SpecialistService struct{}

func NewSpecialistService() *SpecialistService {
	return &SpecialistService{}
}

func (s SpecialistService) ReadSpecialist(selectEnv string) error {

	// Retrieve data from sftp
	specialist, err := downloadLastFileOfSpecialist(selectEnv)
	if err != nil {
		return fmt.Errorf("failed to retrieve data : %s", err)
	}

	// Open the file
	fileHandle, err := os.Open(specialist.Path)
	if err != nil {
		return fmt.Errorf("failed to open file : %s", err)
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			return
		}
	}(fileHandle)

	scanner := bufio.NewScanner(fileHandle)

	var rowNo = 0
	var specialistReqSlice []models.SpecialistRequest
	for scanner.Scan() {
		readLine := scanner.Text()
		specialistReq, err := models.NewSpecialistRequest(readLine, rowNo)
		if err != nil {
			return err
		}

		specialistReqSlice = append(specialistReqSlice, *specialistReq)
		rowNo++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if len(specialistReqSlice) == 0 {
		return errors.New("specialist request is empty")
	}

	dataSpecialist = &specialistReqSlice
	return nil
}

func downloadLastFileOfSpecialist(envReq string) (*models.SftpInfo, error) {
	var sftpInfo *models.SftpInfo
	getConfig := config.GetConfig()

	selectedEnv, err := config.GetEnvSiteData(envReq)
	if err != nil {
		return nil, err
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

// SystemLineSeparator provides a platform-independent line separator
func SystemLineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

func (s SpecialistService) SaveSpecialist(specialistReq models.SpecialistRequest) error {

	// Update Specialist
	if specialistReq.RowNo != "" {
		for index, value := range *dataSpecialist {
			if value.RowNo == specialistReq.RowNo {
				(*dataSpecialist)[index] = specialistReq
				return nil
			}
		}

		return fmt.Errorf("rowNo %v not found in specialist request", specialistReq.RowNo)

	} else {
		// Add Specialist
		newRowNo := len(*dataSpecialist) + 1
		specialistReq.RowNo = strconv.Itoa(newRowNo)
		requests := append(*dataSpecialist, specialistReq)
		dataSpecialist = &requests
	}

	return nil
}

func (s SpecialistService) ExportSpecialist(envReq string) error {
	getConfig := config.GetConfig()

	selectedEnv, err := config.GetEnvSiteData(envReq)
	if err != nil {
		return err
	}

	const filenamePattern = "EIM_EDGE_BLACKLIST_%s.txt"

	// Get the current date in the required format
	currentDate := time.Now().Format("20060102") // yyyyMMdd format

	// Construct the full paths for body and control files
	bodyFile := fmt.Sprintf(filenamePattern, currentDate)
	remoteBodyFile := getConfig.Sftp.Export + bodyFile
	remoteCtrlFile := getConfig.Sftp.Export + "CTRL_" + bodyFile

	// SSH configuration
	sshConfig := &ssh.ClientConfig{
		User: selectedEnv.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(selectedEnv.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", selectedEnv.Host, selectedEnv.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}
	defer func(conn *ssh.Client) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	// Create SFTP client
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {
			return
		}
	}(sftpClient)

	// Helper function to create and write to a file
	createAndWriteRemoteFile := func(path string, data *[]models.SpecialistRequest) error {
		if data != nil {
			index := 0
			content := ""
			for _, line := range *data {
				if index > 0 {
					content += SystemLineSeparator() + line.String()
				} else {
					content += line.String()
				}
				index++
			}

			err := os.WriteFile(path, []byte(content), 0644)
			if err != nil {
				return err
			}
		} else {
			err := os.WriteFile(path, []byte{}, 0644)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Create and write data to the remote body file
	if err := createAndWriteRemoteFile(remoteBodyFile, dataSpecialist); err != nil {
		return err
	}

	// Create the empty remote control file
	if err := createAndWriteRemoteFile(remoteCtrlFile, nil); err != nil {
		return err
	}

	return nil
}

func (s SpecialistService) DeleteByIndex(rowNo string) error {
	if isEmpty(dataSpecialist) {
		return errors.New("dataSpecialist is empty. Please select a file to read data")
	}

	for index, value := range *dataSpecialist {
		if value.RowNo == rowNo {
			// Remove the element at the specified index
			*dataSpecialist = append((*dataSpecialist)[:index], (*dataSpecialist)[index+1:]...)
			return nil
		}
	}

	return fmt.Errorf("row number %s not found in dataSpecialist", rowNo)
}

func (s SpecialistService) FilterByCustomerNo(customerNo string) (result models.SpecialistRequest, err error) {
	if isEmpty(dataSpecialist) {
		return result, errors.New("dataSpecialist is empty. please select file for read data ")
	}

	for _, value := range *dataSpecialist {
		if customerNo == value.CustomerNo {
			return value, nil
		}
	}

	return result, fmt.Errorf("data not found for customerNo %s", customerNo)
}

func isEmpty(s *[]models.SpecialistRequest) bool {
	return s == nil || len(*s) == 0
}
