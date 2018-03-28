package main

import "fmt"

type TrendingRepo struct{
	FullName	string `json:"full_name"`
	Language	string `json:"language"`
	Description	string `json:"description"`
	Stars 		string `json:"stars"`
	Forkers		string `json:"forkers"`
	Gains		string `json:"gains"`
}

func (repo *TrendingRepo) String() string {
	return fmt.Sprintf("\nFullName: %v\nLanguage: %v\nDescription: %v\nStars: %v\nForkers: %v\nGains: %v\n",
		repo.FullName, repo.Language, repo.Description, repo.Stars, repo.Forkers, repo.Gains)
}