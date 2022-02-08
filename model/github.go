package model

type Result struct {
	RepositoryOwner struct {
		Login string `json:"login"`
		User  struct {
			Bio                     string `json:"bio"`
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"total_contributions"`
					Weeks              []struct {
						ContributionDays []struct {
							ContributionCount int    `json:"contribution_count"`
							Date              string `json:"date"`
						} `json:"contribution_days"`
					} `json:"weeks"`
				} `json:"contribution_calendar"`
			} `json:"contributions_collection"`
		} `json:"user"`
	} `json:"repository_owner"`
}

type ResultContributionDays struct {
	ContributionCount int
	Date              string
}

type ResultWeeks struct {
	ContributionDays []ResultContributionDays
}
