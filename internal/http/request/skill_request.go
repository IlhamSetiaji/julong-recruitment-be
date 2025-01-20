package request

import "mime/multipart"

type Skill struct {
	Name            string                `form:"name" validate:"required"`
	Description     string                `form:"description" validate:"required"`
	Certificate     *multipart.FileHeader `form:"certificate" validate:"omitempty"`
	CertificatePath string                `form:"certificate_path" validate:"omitempty"`
}
