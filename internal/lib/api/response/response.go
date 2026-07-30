package response

type Response struct {
	Status string `json:"code"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK  = "OK"
	StatusErr = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusErr,
		Error:  msg,
	}
}
