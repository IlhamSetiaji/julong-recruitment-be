package request

type CreateOrUpdateDocumentVerificationLine struct {
	DocumentVerificationHeaderID string `json:"document_verification_header_id" validate:"required,uuid"`
	DocumentVerificationLines    []struct {
		ID                     string `json:"id" validate:"omitempty,uuid"`
		DocumentVerificationID string `json:"document_verification_id" validate:"required,uuid"`
	} `json:"document_verification_lines" validate:"required,dive"`
	DeletedDocumentVerificationLineIDs []string `json:"deleted_document_verification_line_ids" validate:"omitempty,dive,uuid"`
}

type UploadDocumentVerificationLine struct {
	ID   string `form:"id" validate:"required,uuid"`
	Path string `form:"path" validate:"omitempty"`
}

type UpdateAnswer struct {
	Answer string `json:"answer" validate:"omitempty"`
}
