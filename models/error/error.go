package error

type Error struct {
    ErrorCode int `json:"error_code"`
    ErrorMessage string `json:"error_message"`
}
