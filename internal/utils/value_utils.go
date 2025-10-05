package value_utils

func GetStrngOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
