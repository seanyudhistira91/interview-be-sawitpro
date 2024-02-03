package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type JwtCustomClaims struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (*JwtCustomClaims) Valid() error {
	panic("unimplemented")
}

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
// TODO DELETE
// func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

// 	var resp generated.HelloResponse
// 	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
// 	return ctx.JSON(http.StatusOK, resp)
// }

// // (GET /hello)
// func (s *Server) Hello2(ctx echo.Context, params generated.Hello2Params) error {

// 	var resp generated.HelloResponse
// 	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
// 	return ctx.JSON(http.StatusOK, resp)
// }

// var privateKey *rsa.PrivateKey

// func init() {
// 	keyFile := "path/to/private_key.pem" // Replace with the path to your private key file
// 	keyBytes, err := os.ReadFile(keyFile)
// 	if err != nil {
// 		log.Fatalf("Failed to read private key file: %v", err)
// 	}

// 	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
// 	if err != nil {
// 		log.Fatalf("Failed to parse private key: %v", err)
// 	}
// }

// // (POST /registration)
func (s *Server) Registration(ctx echo.Context) error {
	var resp generated.RegistrationResponse
	var params generated.RegistrationParam
	var errors []string

	if err := json.NewDecoder(ctx.Request().Body).Decode(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid JSON format")
	}

	// Checking Request Body
	if params.PhoneNumber == "" {
		errors = append(errors, "PhoneNumber is required")
	}
	if params.FullName == "" {
		errors = append(errors, "FullName is required")
	}

	if params.Password == "" {
		errors = append(errors, "Password is required")
	}

	// Phone numbers must be at minimum 10 characters and maximum 13 characters.
	if len(params.PhoneNumber) < 10 || len(params.PhoneNumber) > 13 {
		errors = append(errors, "PhoneNumber cannot less than 10 or more than 13 characters")
	}

	// Phone numbers must start with the Indonesia country code “+62”.
	patternFormatPhoneNumber := `^\+62\d{3,13}$`
	re := regexp.MustCompile(patternFormatPhoneNumber)
	if !re.MatchString(params.PhoneNumber) {
		errors = append(errors, "PhoneNumber must start with the Indonesia country code “+62”")
	}

	// Full name must be at minimum 3 characters and maximum 60 characters.
	if len(params.FullName) < 3 || len(params.FullName) > 60 {
		errors = append(errors, "FullName cannot less than 3 or more than 60 characters")
	}

	// Passwords must be minimum 6 characters and maximum 64 characters,
	if len(params.Password) < 6 || len(params.Password) > 64 {
		errors = append(errors, "Passwords cannot less than 6 or more than 64 characters")
	}

	// containing at least 1 capital characters AND 1 number AND 1 special (nonalpha-numeric) characters.
	if !validatepatternFormatPassword(params.Password) {
		errors = append(errors, "Password containing at least 1 capital characters AND 1 number AND 1 special (nonalpha-numeric) characters")
	}

	if len(errors) > 0 {
		return ctx.JSON(http.StatusBadRequest, errors)
	}

	// hasing & salt password
	hashedPassword, err := hashAndSaltPassword(params.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	// map to type
	registTypeInput := repository.GetRegistrationInput{
		FullName:    params.FullName,
		PhoneNumber: params.PhoneNumber,
		Password:    hashedPassword,
	}

	res, err := s.Repository.CreateNewUser(ctx.Request().Context(), registTypeInput)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	resp.Id = res.Id
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	var resp generated.LoginResponse
	var params generated.LoginParam

	if err := json.NewDecoder(ctx.Request().Body).Decode(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON format")
	}

	// mapping login input
	loginInput := repository.GetLoginInput{
		PhoneNumber: params.PhoneNumber,
	}

	// get user by phone number
	res, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), loginInput)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// Compare password
	if !comparePasswords(res.Password, []byte(params.Password)) {
		return ctx.JSON(http.StatusBadRequest, "Invalid Password")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		res.FullName,
		res.PhoneNumber,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	// update flag user successful_login
	updateParam := repository.PostUpdateUserSuccesLoginInput{
		Id: res.Id,
	}

	err = s.Repository.UpdateUserSuccesLogin(ctx.Request().Context(), updateParam)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// map response
	resp = generated.LoginResponse{
		Id:    res.Id,
		Token: t,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) MyProfile(ctx echo.Context) error {
	claims, err := extractJWTClaims(ctx)

	if err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}

	resp := generated.MyProfileResponse{
		Name:        claims["name"].(string),
		PhoneNumber: claims["phoneNumber"].(string),
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var params generated.UpdateProfileParam
	var resp generated.UpdateProfileResponse

	claims, err := extractJWTClaims(ctx)
	oldPhoneNumber := claims["phoneNumber"].(string)

	if err != nil {
		return ctx.JSON(http.StatusForbidden, err.Error())
	}

	if err := json.NewDecoder(ctx.Request().Body).Decode(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, "invalid JSON format")
	}

	// mapping update profile input
	loginInput := repository.UpdateUserInput{
		OldPhoneNumber: oldPhoneNumber,
		Name:           params.FullName,
		PhoneNumber:    params.PhoneNumber,
	}

	// get user by phone number
	res, err := s.Repository.UpdateUserByPhoneNumber(ctx.Request().Context(), loginInput)
	if err != nil && err.Error() == "pq: duplicate key value violates unique constraint \"users_phone_number_key\"" {
		return ctx.JSON(http.StatusConflict, errors.New("duplicate phoneNumber"))
	}

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// mapping response
	resp = generated.UpdateProfileResponse{
		Id: res.Id,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// this function for validate pattern format password
// containing at least 1 capital characters AND 1 number AND 1 special (nonalpha-numeric) characters.
func validatepatternFormatPassword(password string) bool {
	hasUppercase := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUppercase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126: // ASCII printable characters excluding alphanumerics
			hasSpecial = true
		}
	}

	return hasUppercase && hasDigit && hasSpecial
}

func hashAndSaltPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func extractJWTClaims(c echo.Context) (jwt.MapClaims, error) {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("authorization token not provided")
	}

	// Extract the token from the "Bearer" prefix
	tokenString = tokenString[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method and provide the secret key
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims from token")
	}

	return claims, nil
}
