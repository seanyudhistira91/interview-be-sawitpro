// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	CreateNewUser(ctx context.Context, input GetRegistrationInput) (output GetRegistrationOutput, err error)
	GetUserByPhoneNumber(ctx context.Context, input GetLoginInput) (output GetLoginOutput, err error)
	UpdateUserSuccesLogin(ctx context.Context, input PostUpdateUserSuccesLoginInput) error
	UpdateUserByPhoneNumber(ctx context.Context, input UpdateUserInput) (output UpdateUserOutput, err error)
}
