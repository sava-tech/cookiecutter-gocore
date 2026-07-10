package domain

type Repository interface {
	UserRepository
	SessionRepository
	VerificationRepository
}
