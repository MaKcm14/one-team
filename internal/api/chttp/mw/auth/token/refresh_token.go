package token

type RefreshToken struct {
}

func NewRefreshToken() RefreshToken {
	return RefreshToken{}
}

func (r RefreshToken) IssueRefreshToken() string {
	return ""
}
