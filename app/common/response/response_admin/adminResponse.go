package response_admin

type DeptResponse struct {
	ID       uint   `json:"id"`
	DeptName string `json:"dept_name"`
}

type ManagerResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type PermissionResponse struct {
	ID         uint   `json:"id"`
	ModuleName string `json:"module_name"`
	ActionName string `json:"action_name"`
}

type PostResponse struct {
	ID       uint   `json:"id"`
	PostName string `json:"post_name"`
	PostCode string `json:"post_code"`
}

type RoleResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
