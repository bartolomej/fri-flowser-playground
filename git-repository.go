package main

import "fmt"

type GitRepository struct{}

func (gh GitRepository) Clone(url string) {
	fmt.Println("Cloning Git repository: " + url)
}

func (gh GitRepository) Commit() {
	fmt.Println("Commiting")
}

func (gh GitRepository) Push() {
	fmt.Println("Pushing")
}
