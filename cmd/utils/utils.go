package utils

import (
	"fmt"
	"github.com/withmandala/go-log"
	"jaytaylor.com/html2text"
	"os"
)

func NewLogger(pwd string) (logger *log.Logger, err error) {
	var file *os.File
	if pwd == "" {
		file = os.Stderr
	} else {
		file, err = os.Create(pwd)
		if err != nil {
			fmt.Println("logger " + pwd + " not created")
			return
		}
	}
	logger = log.New(file)
	logger.WithDebug()
	logger.WithColor()
	return
}

func HtmlToString(htmlString string) string {
	text, err := html2text.FromString(htmlString, html2text.Options{PrettyTables: true})
	if err != nil {
		fmt.Println(err)
	}
	return text
}
