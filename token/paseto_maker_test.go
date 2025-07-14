package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/volskyi-dmytro/st-bank/util"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoTokenKeySize(t *testing.T) {
	testCases := []struct {
		name    string
		keySize int
	}{
		{
			name:    "TooShort",
			keySize: chacha20poly1305.KeySize - 1,
		},
		{
			name:    "TooLong",
			keySize: chacha20poly1305.KeySize + 1,
		},
		{
			name:    "Empty",
			keySize: 0,
		},
		{
			name:    "VeryShort",
			keySize: 16,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			maker, err := NewPasetoMaker(util.RandomString(tc.keySize))
			require.Error(t, err)
			require.EqualError(t, err, "invalid key size: must be exactly 32 characters")
			require.Nil(t, maker)
		})
	}
}

func TestInvalidPasetoTokenMalformed(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	testCases := []string{
		"invalid.token.format",
		"",
		"random-string-not-paseto",
		"v2.local.invalid-base64",
		"v1.local.some-token", // wrong version
		"v2.public.some-token", // wrong purpose
	}

	for _, token := range testCases {
		payload, err := maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrInvalidToken.Error())
		require.Nil(t, payload)
	}
}

func TestPasetoTokenWithDifferentKey(t *testing.T) {
	maker1, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	maker2, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker1.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker2.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestPasetoTokenExpiredPayload(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Millisecond * 100

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	time.Sleep(time.Millisecond * 200)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestPasetoTokenEdgeCases(t *testing.T) {
	t.Run("EmptyUsername", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.RandomString(32))
		require.NoError(t, err)

		token, err := maker.CreateToken("", time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, "", payload.Username)
	})

	t.Run("ZeroDuration", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.RandomString(32))
		require.NoError(t, err)

		token, err := maker.CreateToken(util.RandomOwner(), 0)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.VerifyToken(token)
		require.Error(t, err)
		require.EqualError(t, err, ErrExpiredToken.Error())
		require.Nil(t, payload)
	})

	t.Run("VeryLongDuration", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.RandomString(32))
		require.NoError(t, err)

		username := util.RandomOwner()
		duration := time.Hour * 24 * 365 // 1 year

		token, err := maker.CreateToken(username, duration)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, username, payload.Username)
	})

	t.Run("Exactly32CharacterKey", func(t *testing.T) {
		key := util.RandomString(32)
		require.Len(t, key, 32)

		maker, err := NewPasetoMaker(key)
		require.NoError(t, err)
		require.NotNil(t, maker)

		username := util.RandomOwner()
		token, err := maker.CreateToken(username, time.Minute)
		require.NoError(t, err)

		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, username, payload.Username)
	})
}

func TestPasetoTokenValidationOrder(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)

	// Verify token multiple times
	for i := 0; i < 5; i++ {
		payload, err := maker.VerifyToken(token)
		require.NoError(t, err)
		require.Equal(t, username, payload.Username)
	}
}

func TestPasetoTokenFormat(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// PASETO v2.local tokens should start with "v2.local."
	require.Contains(t, token, "v2.local.")
	require.True(t, len(token) > 9) // Minimum length check
}

func TestPasetoTokenSecurity(t *testing.T) {
	t.Run("TokensAreUnique", func(t *testing.T) {
		maker, err := NewPasetoMaker(util.RandomString(32))
		require.NoError(t, err)

		username := util.RandomOwner()
		duration := time.Minute

		// Generate multiple tokens for the same user
		tokens := make([]string, 10)
		for i := 0; i < 10; i++ {
			token, err := maker.CreateToken(username, duration)
			require.NoError(t, err)
			tokens[i] = token
		}

		// All tokens should be unique (due to random nonce in PASETO)
		uniqueTokens := make(map[string]bool)
		for _, token := range tokens {
			require.False(t, uniqueTokens[token], "Duplicate token found")
			uniqueTokens[token] = true
		}
	})

	t.Run("KeyIsolation", func(t *testing.T) {
		// Test that tokens created with one key cannot be verified with another
		makers := make([]*PasetoMaker, 5)
		for i := 0; i < 5; i++ {
			maker, err := NewPasetoMaker(util.RandomString(32))
			require.NoError(t, err)
			makers[i] = maker.(*PasetoMaker)
		}

		username := util.RandomOwner()
		duration := time.Minute

		// Create token with first maker
		token, err := makers[0].CreateToken(username, duration)
		require.NoError(t, err)

		// Try to verify with other makers - should all fail
		for i := 1; i < 5; i++ {
			payload, err := makers[i].VerifyToken(token)
			require.Error(t, err)
			require.EqualError(t, err, ErrInvalidToken.Error())
			require.Nil(t, payload)
		}
	})
}

func TestPasetoTokenStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	const numTokens = 1000
	tokens := make([]string, numTokens)
	usernames := make([]string, numTokens)

	// Create many tokens
	for i := 0; i < numTokens; i++ {
		username := util.RandomOwner()
		usernames[i] = username

		token, err := maker.CreateToken(username, time.Hour)
		require.NoError(t, err)
		tokens[i] = token
	}

	// Verify all tokens
	for i := 0; i < numTokens; i++ {
		payload, err := maker.VerifyToken(tokens[i])
		require.NoError(t, err)
		require.Equal(t, usernames[i], payload.Username)
	}
}