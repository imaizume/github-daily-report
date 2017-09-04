package main

type GHPushEventCommit []github.PushEventCommit

func (e GHPushEventCommit) Len() int {
	return len(e)
}

func (e GHPushEventCommit) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e GHPushEventCommit) Less(i, j int) bool {
	return e[i].GetTimestamp().Time.After(e[i].GetTimestamp().Time)
}
