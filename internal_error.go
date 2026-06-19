package mo

func InternalServerError()*HttpError{
	return &HttpError{
		&Response{
			ContentType: JSON,
			Body: [...]string{"Internal server error"},
		},
	}
}