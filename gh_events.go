package main

import("github.com/google/go-github/github")

type GHEvents []github.Event

func (e GHEvents) Len() int {
	return len(e)
}

func (e GHEvents) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e GHEvents) Less(i, j int) bool {
	return e[i].GetCreatedAt().After(e[i].GetCreatedAt())
}
