package main

import (
	"github.com/charmbracelet/log"
	"github.com/magefile/mage/mg"
)

func Test() error {
	log.Info("Test")
	// Starting dagger engine && api session

	// Reading dir including file

	// Testing in a golang container

	return nil
}

func Build() error {
	mg.Deps(Test)
	log.Info("Build")
	// Starting dagger engine && api session

	// Reading dir exluding file

	// Building in a golang container

	return nil
}

func Push() error {
	mg.Deps(Build)
	log.Info("Push")
	// Starting dagger engine && api session

	// Reading dir exluding file

	// Building in golang container

	// Service container ?

	// Publish

	return nil
}
