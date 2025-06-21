package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	cError = color.New(color.FgRed, color.Bold)
	cInfo  = color.New(color.FgGreen)
	cWarn  = color.New(color.FgYellow)
	p      = &printer{}
)

type printer struct{}

func (p *printer) Error(msg string) {
	cError.Print("ERRO: ")
	fmt.Println(msg)
}

func (p *printer) Errorf(format string, args ...interface{}) {
	cError.Print("ERRO: ")
	fmt.Printf(format, args...)
}

func (p *printer) Info(msg string) {
	cInfo.Print("INFO: ")
	fmt.Println(msg)
}

func (p *printer) Infof(format string, args ...interface{}) {
	cInfo.Print("INFO: ")
	fmt.Printf(format, args...)
}

func (p *printer) Warn(msg string) {
	cWarn.Print("WARN: ")
	fmt.Println(msg)
}

func (p *printer) Warnf(format string, args ...interface{}) {
	cWarn.Print("WARN: ")
	fmt.Printf(format, args...)
}
