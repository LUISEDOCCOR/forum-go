package utils

func CreateResponse(mode string, msg string) map[string]string {
	var response = map[string]string{
		"mode": mode,
		"msg":  msg,
	}
	return response

}

func CreateResponseAuth(token string, username string, id uint) map[string]map[string]interface{} {
	var response = map[string]map[string]interface{}{
		"data": {
			"jwt":      token,
			"username": username,
			"id":       id,
		},
	}
	return response
}
