package password

import "golang.org/x/crypto/bcrypt"


type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *PasswordService) ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
