// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type backlog struct {
	lowerBound int
	upperBound int
	candidates int
}
type CRS struct {
	backlog            [15]backlog
	excel              *excelize.File
	previousStepSize   int
	previousDrawCutOff int
	previousInvites    int
	previousDrawDate   string
}

func parseScore(score string) backlog {
	var newbacklog backlog
	var err error

	bounds := strings.Split(score, "-")
	if len(bounds) != 2 {
		panic("Invalid score format in excel")
	}
	newbacklog.upperBound, err = strconv.Atoi(bounds[1])
	if err != nil {
		panic("Invalid score format in excel")
	}
	newbacklog.lowerBound, err = strconv.Atoi(bounds[0])
	if err != nil {
		panic("Invalid score format in excel")
	}
	return newbacklog
}

func parseCandidate(candidate string) int {
	var candidateCount int
	candidateCount, err := strconv.Atoi(candidate)
	if err != nil {
		panic("Invalid candidate formta in excel")
	}
	return candidateCount
}

func (c *CRS) GetPreviousDrawCutoff() int {
	return c.previousDrawCutOff
}

func (c *CRS) GetPreviousDrawStepsize() int {
	return c.previousStepSize
}
func (c *CRS) GetPreviousTotalInvitesSent() int {
	return c.previousInvites
}

func (c *CRS) GetPreviousDrawDate() time.Time {
	date, err := time.Parse("MM-DD-YYYY", c.previousDrawDate)
	if err != nil {
		panic("invalid date")
	}
	return date
}

func (c *CRS) readValues() {
	for i := 2; i <= 16; i++ {
		// read and parse score
		scoreAxis := fmt.Sprintf("A%d", i)
		score, err := c.excel.GetCellValue("Sheet1", scoreAxis)
		if err != nil {
			fmt.Println(err)
			return
		}

		newbacklog := parseScore(score)
		candidateAxis := fmt.Sprintf("B%d", i)
		candidates, err := c.excel.GetCellValue("Sheet1", candidateAxis)
		if err != nil {
			fmt.Println(err)
			return
		}

		candidateCount := parseCandidate(candidates)
		newbacklog.candidates = candidateCount
		c.backlog[i-2] = newbacklog
	}
	coAxis := "C1"
	tCO, err := c.excel.GetCellValue("Sheet1", coAxis)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.previousDrawCutOff, _ = strconv.Atoi(tCO)
	coStepSizeAxis := "C2"
	tCO, err = c.excel.GetCellValue("Sheet1", coStepSizeAxis)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.previousStepSize, _ = strconv.Atoi(tCO)
	coInvitesAxis := "C3"
	tCO, err = c.excel.GetCellValue("Sheet1", coInvitesAxis)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.previousInvites, _ = strconv.Atoi(tCO)

	coDate := "C4"
	c.previousDrawDate, err = c.excel.GetCellValue("Sheet1", coDate)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getFilePath() string {
	return os.TempDir() + "crs_file.xlxs"
}

func (c *CRS) Get_crs_distribution(url string) {
	//Download CRS file
	DownloadFile(url, getFilePath())
	// for now the file is static named
	f, err := excelize.OpenFile(getFilePath())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	c.excel = f
	c.readValues()

	fmt.Println(c.backlog)
}
