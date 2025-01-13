package request

type CreateMPRequest struct {
	MPRCloneID string `json:"mpr_clone_id" validate:"required"`
}
