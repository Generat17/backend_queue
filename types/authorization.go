package types

type ParseTokenWorkstationResponse struct {
	UserId        int
	WorkstationId int
}

type SessionInfo struct {
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at" db:"expires_at"`
	Workstation  int    `json:"workstation_id" db:"workstation_id"`
}

type SignUpResponse struct {
	Id int `json:"id" binding:"required"`
}

type AuthorizationResponse struct {
	AccessToken  string      `json:"accessToken" binding:"required"`
	RefreshToken string      `json:"refreshToken" binding:"required"`
	Employee     Employee    `json:"employee" binding:"required"`
	Workstation  Workstation `json:"workstation" binding:"required"`
}

type LogoutResponse struct {
	StatusResponse bool `json:"status_response" binding:"required"`
}
