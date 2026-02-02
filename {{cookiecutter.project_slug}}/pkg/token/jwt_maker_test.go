package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"{{ cookiecutter.module_path }}/utils"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)

	email := utils.RandomEmail()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(email, "personal", duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, email, payload.Email)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomEmail(), "personal", -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

func TestInvalideJWTTokenALgNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomEmail(), "personal", time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalideToken.Error())
	require.Nil(t, payload)

}

func TestTamperedJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomEmail(), "personal", time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// Tamper token by changing a character
	tamperedToken := token[:len(token)-1] + "X"

	payload, err = maker.VerifyToken(tamperedToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalideToken.Error())
	require.Nil(t, payload)
}
