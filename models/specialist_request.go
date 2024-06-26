package models

import (
	"fmt"
	"strconv"
	"strings"
)

// SpecialistRequest represents your SpecialistRequest data.
type SpecialistRequest struct {
	RowNo          string
	CustomerNo     string
	CustomerNameTh string
	ChkThaiName    string
	ChkEngName     string
	CustomerNameEn string
	Dob            string
	CountryCode    string
	CustomerType   string
	Zipcode        string
	Action         string
	OldAct         string
	Ovract         string
	Pob            string
	ReasonCode     string
	RtnCustomer    string
	SrcSeq         string
}

// NewSpecialistRequest creates a new SpecialistRequest from a data line.
func NewSpecialistRequest(dataLine string, index int) *SpecialistRequest {
	dataFields := strings.Split(dataLine, "|")

	rowNo := dataFields[0]
	if index > 0 {
		rowNo = strconv.Itoa(index)
	}

	return &SpecialistRequest{
		RowNo:          rowNo,
		CustomerNo:     dataFields[1],
		CustomerNameTh: dataFields[2],
		ChkThaiName:    dataFields[3],
		ChkEngName:     dataFields[4],
		CustomerNameEn: dataFields[5],
		Dob:            dataFields[6],
		CountryCode:    dataFields[7],
		CustomerType:   dataFields[8],
		Zipcode:        dataFields[9],
		Action:         dataFields[10],
		OldAct:         dataFields[11],
		Ovract:         dataFields[12],
		Pob:            dataFields[13],
		ReasonCode:     dataFields[14],
		RtnCustomer:    dataFields[15],
		SrcSeq:         dataFields[16],
	}
}

// String returns a string representation of the SpecialistRequest.
func (s *SpecialistRequest) String() string {
	toString := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
		s.RowNo, s.CustomerNo, s.CustomerNameTh, s.ChkThaiName, s.ChkEngName, s.CustomerNameEn,
		s.Dob, s.CountryCode, s.CustomerType, s.Zipcode, s.Action, s.OldAct, s.Ovract, s.Pob,
		s.ReasonCode, s.RtnCustomer, s.SrcSeq)
	return toString
}
