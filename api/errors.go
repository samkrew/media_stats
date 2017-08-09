package api

type Error struct {
	Type  string `json:"type"`
	Field string `json:"field,omitempty"`
	Error string `json:"error"`
}

var (
	ErrorAuthForbidden = Error{
			Type:  "auth",
			Error: "forbidden",
	}

	ErrorNotAllowed = Error{
		Type:  "access",
		Error: "not allowed",
	}

	ErrorUrlInvalid = Error{
		Type:  "query",
		Field: "url",
		Error: "invalid",
	}

	ErrorHashInvalid = Error{
		Type:  "query",
		Field: "md5",
		Error: "invalid",
	}
)
