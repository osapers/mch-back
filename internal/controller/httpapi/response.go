package httpapi

type response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func jsonResp(data interface{}, err error) interface{} {
	resp := response{}

	if data != nil {
		resp.Data = data
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp
}
