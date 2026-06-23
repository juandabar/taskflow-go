package valueobject

type UserRole string

const (
	RoleAdmin  UserRole = "ADMIN"
	RoleMember UserRole = "MEMBER"
)

func (r UserRole) IsValid() bool {
	return r == RoleAdmin || r == RoleMember
}

func (r UserRole) IsAdmin() bool {
	return r == RoleAdmin
}
