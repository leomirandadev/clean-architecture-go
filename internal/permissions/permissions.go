package permissions

type Permission uint8

const (
	Free Permission = iota
	Premium
	Admin
)

var permissionsMap = map[Permission]string{
	Free:    "free",
	Premium: "premium",
	Admin:   "admin",
}

func (value Permission) String() string {
	return permissionsMap[value]
}

func Strings() []string {
	result := make([]string, 0, len(permissionsMap))
	for _, v := range permissionsMap {
		result = append(result, v)
	}
	return result
}
