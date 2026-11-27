package domain

// Actor is the signed-in user an operation acts as: who they are and whether
// they are an administrator. A module that judges what someone may do to
// their own footing (the self-protection rules of ADR-0020: their own account,
// Staff Member record, role grants, the administrator flag) takes it as a
// plain value, so every adapter passes the same facts and the rule is written
// once, next to the write it guards.
type Actor struct {
	UserID        int
	Administrator bool
}
