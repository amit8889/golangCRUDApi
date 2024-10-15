package types

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,gt=18"`
}
