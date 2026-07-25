package utils

// Response adalah bentuk standar untuk setiap response API, sehingga klien
// tidak perlu menebak format dari endpoint tertentu.
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func BuildSuccess(message string, data any) Response {
	return Response{Success: true, Message: message, Data: data}
}

func BuildError(message string) Response {
	return Response{Success: false, Message: message}
}
