package debug

import (
	"encoding/json"
	"fmt"
	"os"
)

func prettyfy(i any) []byte {
	pp, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		panic(err)
	}
	return pp
}

// PrettyPrint prints structs in a readable way in the terminal.
func PrettyPrint(i interface{}) {
	fmt.Println(string(prettyfy(i)))
}

// PrettyPrintToFile writes the struct to a file in a readable way.
func PrettyPrintToFile(i interface{}, filepath string) {
	content := prettyfy(i)

	// 2. Create (or truncate) the file and handle any file I/O errors.
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 3. Write the JSON bytes to the file.
	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}
}

// PrintToFile writes the struct to a file.
func PrintToFile(i interface{}, filepath string) {
	content, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		panic(err)
	}
}

// WriteToFile creates (or truncates) a file and writes the provided content to it.
func WriteToFile(filename, content string) error {
	// Open the file for writing, creating it if it doesn't exist, and truncating it if it does.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed after writing

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
