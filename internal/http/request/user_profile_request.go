package request

import "mime/multipart"

type FillUserProfileRequest struct {
	ID              string                `form:"id" validate:"omitempty,uuid"`
	MaritalStatus   string                `form:"marital_status" validate:"required,marital_status_validation"`
	Gender          string                `form:"gender" validate:"required,user_gender_validation"`
	PhoneNumber     string                `form:"phone_number" validate:"required"`
	Age             int                   `form:"age" validate:"required"`
	BirthDate       string                `form:"birth_date" validate:"required,datetime=2006-01-02"`
	BirthPlace      string                `form:"birth_place" validate:"required"`
	Ktp             *multipart.FileHeader `form:"ktp" validate:"required"`
	CurriculumVitae *multipart.FileHeader `form:"curriculum_vitae" validate:"required"`
	KtpPath         string                `form:"ktp_path" validate:"omitempty"`
	CvPath          string                `form:"cv_path" validate:"omitempty"`
	WorkExperiences []struct {
		Name            string                `form:"name" validate:"required"`
		CompanyName     string                `form:"company_name" validate:"required"`
		YearExperience  int                   `form:"year_experience" validate:"required"`
		JobDescription  string                `form:"job_description" validate:"required"`
		Certificate     *multipart.FileHeader `form:"certificate" validate:"omitempty"`
		CertificatePath string                `form:"certificate_path" validate:"omitempty"`
	} `form:"work_experiences" validate:"omitempty,dive"`
	Educations []struct {
	} `form:"educations" validate:"omitempty,dive"`
}
