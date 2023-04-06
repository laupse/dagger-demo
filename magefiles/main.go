package main

import (
	"github.com/charmbracelet/log"
	"github.com/magefile/mage/mg"
)

func Test() error {
	log.Info("Test")
	return nil
}

func Build() error {
	mg.Deps(Test)
	log.Info("Build")
	return nil
}

func Push() error {
	mg.Deps(Build)
	log.Info("Push")
	return nil
}
