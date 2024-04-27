package service

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/saurabh-sde/library-task-go/model"
	"github.com/saurabh-sde/library-task-go/utility"
	"github.com/spf13/cast"
)

const (
	Admin   = "admin"
	Regular = "regular"
	DirPath = "resources/"
)

func GetBooksByUserType(userType string) (resp []model.Book, err error) {
	funcName := "GetBooksByUserType"
	utility.Print(funcName, userType)
	resp = []model.Book{}

	fileName := "regularUser.csv"
	// open file
	file, err := os.Open(DirPath + fileName)
	if err != nil {
		utility.Error(err)
		return nil, err
	}
	defer file.Close()

	// read all rows
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		utility.Error(err)
		return nil, err
	}
	// loop on rows
	for i, record := range records {
		// skip header row
		if i == 0 {
			continue
		}
		// add book data struct
		data := model.Book{
			BookName:        record[0],
			Author:          record[1],
			PublicationYear: cast.ToInt(record[2]),
		}
		// append to resp
		resp = append(resp, data)
	}

	if userType == "admin" {
		// append admin csv data also
		adminFileName := "adminUser.csv"
		// open file
		adminFile, err := os.Open(DirPath + adminFileName)
		if err != nil {
			utility.Error(err)
			return nil, err
		}
		defer adminFile.Close()

		// read all rows
		records, err := csv.NewReader(adminFile).ReadAll()
		if err != nil {
			utility.Error(err)
			return nil, err
		}
		// loop on rows
		for i, record := range records {
			// skip header row
			if i == 0 {
				continue
			}
			// add book data struct
			data := model.Book{
				BookName:        record[0],
				Author:          record[1],
				PublicationYear: cast.ToInt(record[2]),
			}
			// append to resp
			resp = append(resp, data)
		}
	}
	return resp, nil
}

func AddBook(bookName, author string, pubYear int) (err error) {
	fileName := "regularUser.csv"
	// open file in append mode
	file, err := os.OpenFile(DirPath+fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		utility.Error(err)
		return err
	}
	defer file.Close()

	book := model.Book{
		BookName:        bookName,
		Author:          author,
		PublicationYear: pubYear,
	}

	w := csv.NewWriter(file)
	defer w.Flush()
	// create row
	row := []string{book.BookName, book.Author, strconv.Itoa(book.PublicationYear)}
	// Using Write to insert row
	if err := w.Write(row); err != nil {
		utility.Error("error writing record to file: ", err)
		return err
	}

	return nil
}

// Note: There is no delete method in csv package offered by Go
/*
	1. Read existing rows
	2. Copy all rows leaving matching bookName
	3. Write to new file
	4. Rename to original file to replace it
	-> Directly loading and editing original file can lead to data loss so work on copy of it then replace it with original

	OR

	- We can also add enable flag and mark deleted book as enable = 0 to avoid data loss and optimize code complexity
	- If data was stored into DB then this would have been fine
*/
func DeleteBook(bookName string) (err error) {
	fileName := "regularUser.csv"
	// open file
	file, err := os.Open(DirPath + fileName)
	if err != nil {
		utility.Error(err)
		return err
	}
	defer file.Close()

	// read all rows
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		utility.Error(err)
		return err
	}
	updatedData := [][]string{}
	// loop on rows to copy data
	for _, record := range records {
		// skip matching book name row
		if strings.TrimSpace(strings.ToLower(record[0])) == strings.TrimSpace(strings.ToLower(bookName)) {
			continue
		}
		updatedData = append(updatedData, record)
	}
	tempFile := "temp.csv"
	// create new file
	newFile, err := os.Create(DirPath + tempFile)
	if err != nil {
		utility.Print("Error creating file: ", err)
		return
	}

	w := csv.NewWriter(newFile)
	// multi-write data to new file
	err = w.WriteAll(updatedData)
	if err != nil {
		utility.Print("Error updating file: ", err)
		return
	}

	// rename to original file
	err = os.Rename(DirPath+tempFile, DirPath+fileName)
	if err != nil {
		utility.Print("Error renaming file: ", err)
		return
	}
	return nil
}
