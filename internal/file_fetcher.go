package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/xuri/excelize/v2"
)

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetCRSDatesFilePath() string {
	tempDir := os.TempDir()
	tempCRSfile := "/crs_cutoff_dates.xlxs"
	return tempDir + tempCRSfile
}

func DownloadCRSDates() {
	const crsCutOffFile = "https://github.com/gandharvas/crs/blob/main/files/crs_cutoff_dates.xlsx?raw=true"
	filePath := GetCRSDatesFilePath()
	if _, err := os.Stat("/path/to/whatever"); err != nil {
		os.Remove(filePath)

	}
	DownloadFile(crsCutOffFile, filePath)

}
func readValues(excel *excelize.File) map[string]string {
	dates := make(map[string]string)
	for i := 1; i <= 10; i++ {
		// read and parse score
		dateAxis := fmt.Sprintf("A%d", i)
		date, err := excel.GetCellValue("Sheet1", dateAxis)
		if err != nil {
			fmt.Println(err)
			return dates
		}

		linksAxis := fmt.Sprintf("B%d", i)
		link, err := excel.GetCellValue("Sheet1", linksAxis)
		if err != nil {
			fmt.Println(err)
			return dates
		}
		dates[date] = link
	}
	return dates
}

func GetCRSDates() map[string]string {
	// for now the file is static named
	f, err := excelize.OpenFile(GetCRSDatesFilePath())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	excel := f
	return readValues(excel)
}
