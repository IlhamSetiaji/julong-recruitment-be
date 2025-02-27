package response

type DashboardResponse struct {
	TotalRecruitmentTargetResponse      TotalRecruitmentTargetResponse      `json:"total_recruitment_target"`
	TotalRecruitmentRealizationResponse TotalRecruitmentRealizationResponse `json:"total_recruitment_realization"`
	TotalBilingualResponse              TotalBilingualResponse              `json:"total_bilingual"`
	ChartDurationRecruitmentResponse    ChartDurationRecruitmentResponse    `json:"chart_duration_recruitment"`
	ChartJobLevelResponse               ChartJobLevelResponse               `json:"chart_job_level"`
	ChartDepartmentResponse             ChartDepartmentResponse             `json:"chart_department"`
	AvgTimeToHireResponse               AvgTimeToHireResponse               `json:"avg_time_to_hire"`
}

type TotalRecruitmentTargetResponse struct {
	TotalRecruitmentTarget int `json:"total_recruitment_target"`
	Percentage             int `json:"percentage"`
}

type TotalRecruitmentRealizationResponse struct {
	TotalRecruitmentRealization int `json:"total_recruitment_realization"`
	Percentage                  int `json:"percentage"`
}

type TotalBilingualResponse struct {
	TotalBilingual int `json:"total_bilingual"`
}

type ChartDurationRecruitmentResponse struct {
	Label    string `json:"label"`
	Datasets []int  `json:"datasets"`
}

type ChartJobLevelResponse struct {
	Label    string `json:"label"`
	Datasets []int  `json:"datasets"`
}

type ChartDepartmentResponse struct {
	Label    string `json:"label"`
	Datasets []int  `json:"datasets"`
}

type AvgTimeToHireResponse struct {
	AvgTimeToHire int `json:"avg_time_to_hire"`
}
