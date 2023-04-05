package main

import "github.com/magefile/mage/mg"

func Test() {

}

func Build() {
	mg.Deps(Test)
	print("build")
}

func Push() {
	mg.Deps(Build)
	print("push")
}
