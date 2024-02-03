package repository

import (
	"context"
	"fmt"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateNewUser(ctx context.Context, input GetRegistrationInput) (output GetRegistrationOutput, err error) {
	err = r.Db.QueryRowContext(ctx, `INSERT INTO users(phone_number, full_name, hash_password) VALUES ($1, $2, $3) RETURNING id`, input.PhoneNumber, input.FullName, input.Password).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, input GetLoginInput) (output GetLoginOutput, err error) {
	err = r.Db.QueryRowContext(ctx, `SELECT id, hash_password, phone_number, full_name FROM users WHERE phone_number = $1`, input.PhoneNumber).Scan(&output.Id, &output.Password, &output.PhoneNumber, &output.FullName)
	if err != nil {
		return output, err
	}
	return
}

func (r *Repository) UpdateUserSuccesLogin(ctx context.Context, input PostUpdateUserSuccesLoginInput) error {
	err := r.Db.QueryRowContext(ctx, `UPDATE users SET count_login = count_login + $1 WHERE id = $2`, 1, input.Id)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (r *Repository) UpdateUserByPhoneNumber(ctx context.Context, input UpdateUserInput) (output UpdateUserOutput, err error) {
	var query string
	if input.Name != nil && input.PhoneNumber != nil {
		query = fmt.Sprintf(`UPDATE users SET full_name = '%s', phone_number = '%s' WHERE phone_number = '%s'`, *input.Name, *input.PhoneNumber, input.OldPhoneNumber)
	}

	if input.Name != nil && input.PhoneNumber == nil {
		query = fmt.Sprintf(`UPDATE users SET full_name = '%s' WHERE phone_number = '%s'`, *input.Name, input.OldPhoneNumber)
	}

	if input.Name == nil && input.PhoneNumber != nil {
		query = fmt.Sprintf(`UPDATE users SET phone_number = '%s' WHERE phone_number = '%s'`, *input.PhoneNumber, input.OldPhoneNumber)
	}

	err = r.Db.QueryRowContext(ctx, query+"RETURNING id").Scan(&output.Id)
	if err != nil {
		return
	}
	return

}
