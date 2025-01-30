package response

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
)

type AdministrativeResultResponse struct {
	ID                        uuid.UUID                         `json:"id"`
	AdministrativeSelectionID uuid.UUID                         `json:"administrative_selection_id"`
	UserProfileID             uuid.UUID                         `json:"user_profile_id"`
	Status                    entity.AdministrativeResultStatus `json:"status"`
	CreatedAt                 string                            `json:"created_at"`
	UpdatedAt                 string                            `json:"updated_at"`

	UserProfile *UserProfileResponse `json:"user_profile"`
}
