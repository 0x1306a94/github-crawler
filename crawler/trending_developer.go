package crawler

import "fmt"

type TrendingDeveloperResult struct {
	Language		string `json:"language"`
	Since			string `json:"since"`
	Developers		[]*TrendingDeveloper `json:"developers"`
}

type TrendingDeveloper struct {
	Ranking		int `json:"ranking"`
	Avatar			string `json:"avatar"`
	Login			string `json:"login"`
	NickName		string `json:"nickName"`
	RepoName		string `json:"repoName"`
	RepoDescription string `json:"repoDescription"`
}

func (developer *TrendingDeveloper) String() string {
	return fmt.Sprintf("\nRanking: %d\nLogin: %v\nNickName: %v\nAvatar: %v\nRepoName: %v\nRepoDescription: %v\n",
		developer.Ranking,developer.Login, developer.NickName, developer.Avatar, developer.RepoName, developer.RepoDescription)
}