package model

type SshRole string

const (
	RoleUser  = SshRole("user")
	RoleAdmin = SshRole("admin")
)

type User struct {
	BaseModel
	Role          string `json:"role"`
	RealName      string `json:"real_name"`
	Account       string `json:"account"`
	Avatar        string `json:"avatar"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Secret        string `json:"-"`
	GithubAccount string `json:"github_account"`
	PublicKey     string `json:"public_key"`
}
