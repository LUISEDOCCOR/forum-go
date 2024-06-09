package utils

func CreateResponse(mode string, msg string) map[string]string {
	var response = map[string]string{
		"mode": mode,
		"msg":  msg,
	}
	return response

}
