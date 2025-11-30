package accesshandler

type CheckInPayload struct {
	QrJWT    string `json:"qr_jwt"`
	BranchID string `json:"branch_id"`
}
