package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const (
	ErrTmplReadCode     = 1
	ErrTmplCreationCode = 1
	ErrFileCreationCode = 1
	ErrGenerateCode     = 1
)

func readFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contentString := string(content)

	return contentString, nil
}

func generateFile(templateFile, generateFile string, data interface{}, regen bool) {
	tmplString, err := readFileContent(templateFile)
	if err != nil {
		fmt.Println("Ошибка при прочтении файла: "+templateFile, err)
		os.Exit(ErrTmplReadCode)
	}

	tmpl, err := template.New(templateFile).Parse(tmplString)
	if err != nil {
		fmt.Println("Ошибка при создании шаблона:", err)
		os.Exit(ErrTmplCreationCode)
	}

	if _, err := os.Stat(generateFile); !errors.Is(err, os.ErrNotExist) && !regen {
		return
	}
	file, err := os.Create(generateFile)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		os.Exit(ErrFileCreationCode)
	}
	defer file.Close()

	if err = tmpl.Execute(file, data); err != nil {
		fmt.Println("Ошибка при генерации кода:", err)
		os.Exit(ErrGenerateCode)
	}

	if err = exec.Command("gofmt", "-s", "-w", generateFile).Run(); err != nil {
		fmt.Println("Ошибка при форматировании кода:", err)
		os.Exit(ErrGenerateCode)
	}

}
