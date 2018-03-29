package crawler

import "fmt"

type TrendingRepoResult struct {
	Language	string `json:"language"`
	Since		string `json:"since"`
	Repos		[]*TrendingRepo `json:"repos"`
}

type TrendingRepo struct{
	Ranking		int `json:"ranking"`
	FullName	string `json:"full_name"`
	Language	string `json:"language"`
	Description	string `json:"description"`
	Stars 		string `json:"stars"`
	Forkers		string `json:"forkers"`
	Gains		string `json:"gains"`
}

func (repo *TrendingRepo) String() string {
	return fmt.Sprintf("\nRanking: %d\nFullName: %v\nLanguage: %v\nDescription: %v\nStars: %v\nForkers: %v\nGains: %v\n",
		repo.Ranking,repo.FullName, repo.Language, repo.Description, repo.Stars, repo.Forkers, repo.Gains)
}