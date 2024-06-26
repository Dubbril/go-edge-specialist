package services

import (
	"bufio"
	"errors"
	"fmt"
	"go-edge-specailist/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var dataSpecialist *[]models.SpecialistRequest

type SpecialistService struct{}

func NewSpecialistService() *SpecialistService {
	return &SpecialistService{}
}

func (s SpecialistService) ReadSpecialist(file *multipart.FileHeader) error {

	// Open the file
	fileHandle, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file : %s", err)
	}
	defer func(fileHandle multipart.File) {
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

		specialistReq := models.NewSpecialistRequest(readLine, rowNo)
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

		rowNo, err := strconv.Atoi(specialistReq.RowNo)
		if err != nil {
			return errors.New("Cannot to convert rowNo to number with error : " + err.Error())
		}

		(*dataSpecialist)[rowNo-1] = specialistReq

	} else {
		// Add Specialist
		newRowNo := len(*dataSpecialist) + 1
		specialistReq.RowNo = strconv.Itoa(newRowNo)
		requests := append(*dataSpecialist, specialistReq)
		dataSpecialist = &requests
	}

	return nil
}

func (s SpecialistService) ExportSpecialist() error {
	//const directory = "/home/edw/outbox/edge"
	const directory = "C:\\Users\\dubbril\\Desktop\\data"
	const filenamePattern = "EIM_EDGE_BLACKLIST_%s.txt"

	// Get the current date in the required format
	currentDate := time.Now().Format("20060102") // yyyyMMdd format

	// Construct the full paths for body and control files
	bodyFile := fmt.Sprintf(filenamePattern, currentDate)
	fullPathBodyFile := filepath.Join(directory, bodyFile)
	fullPathCtrlFile := filepath.Join(directory, "CTRL_"+bodyFile)

	// Helper function to create and write to a file
	createAndWriteFile := func(path string, data *[]models.SpecialistRequest) error {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		if data != nil {
			index := 0
			for _, line := range *data {
				if index > 0 {
					if _, err := file.WriteString(SystemLineSeparator() + line.String()); err != nil {
						return fmt.Errorf("failed to write to file %s: %w", path, err)
					}
				} else {
					if _, err := file.WriteString(line.String()); err != nil {
						return fmt.Errorf("failed to write to file %s: %w", path, err)
					}
				}
				index++
			}
		}
		return nil
	}

	// Create and write data to the body file
	if err := createAndWriteFile(fullPathBodyFile, dataSpecialist); err != nil {
		return err
	}

	// Create the empty control file
	if err := createAndWriteFile(fullPathCtrlFile, nil); err != nil {
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
