package helper

func JsonMessage(status string, message string) map[string]interface{} {
	var x = map[string]interface{}{
		"status":  status,
		"message": message,
	}
	return x
}