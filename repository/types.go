// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

// Registration
type GetRegistrationInput struct {
	FullName    string
	PhoneNumber string
	Password    string
}

type GetRegistrationOutput struct {
	Id int
}

// Login
type GetLoginInput struct {
	PhoneNumber string
}

type GetLoginOutput struct {
	Id          int
	Password    string
	PhoneNumber string
	FullName    string
}

type PostUpdateUserSuccesLoginInput struct {
	Id int
}

// UpdateUser/Profile
type UpdateUserInput struct {
	Name           *string
	PhoneNumber    *string
	OldPhoneNumber string
}

type UpdateUserOutput struct {
	Id int
}
