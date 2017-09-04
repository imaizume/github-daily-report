package main

import("github.com/google/go-github/github")

type GHPushEventCommits []github.PushEventCommit

func (e GHPushEventCommits) Len() int {
	return len(e)
}

func (e GHPushEventCommits) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e GHPushEventCommits) Less(i, j int) bool {
	return e[i].GetTimestamp().Time.After(e[i].GetTimestamp().Time)
}
