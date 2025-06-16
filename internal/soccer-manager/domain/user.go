package domain

type UserRole string

const (
	UserRoleUSER  UserRole = "USER"
	UserRoleADMIN UserRole = "ADMIN"
)

func (e UserRole) Valid() bool {
	switch e {
	case UserRoleADMIN,
		UserRoleUSER:
		return true
	}
	return false
}

type User struct {
	ID int64
	UserInfo
}

type UserInfo struct {
	Username string
	Role     UserRole
}

func NewUser(id int64, username string, role string) User {
	return User{
		ID: id,
		UserInfo: UserInfo{
			Username: username,
			Role:     UserRole(role),
		},
	}
}
