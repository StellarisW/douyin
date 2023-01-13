package auth

type Token struct {
	TokenValue string `json:"token_value,omitempty"` // 令牌值
	ExpiresAt  int64  `json:"expires_at,omitempty"`  // 过期时间 (unix)
}
