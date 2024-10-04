package utils

import "os"

func Write2File(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, d := range data {
		file.WriteString(d + "\n")
	}
	return nil
}
