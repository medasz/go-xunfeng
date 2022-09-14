package tools

import "os"

func OpenOrCreateTargetFile(content string) error {
	f, err := os.OpenFile("target.log", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
