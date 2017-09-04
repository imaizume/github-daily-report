package main

import("github.com/google/go-github/github")

type GHEvent []github.Event

func (e GHEvent) Len() int {
	return len(e)
}

func (e GHEvent) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e GHEvent) Less(i, j int) bool {
	return e[i].GetCreatedAt().After(e[i].GetCreatedAt())
}
