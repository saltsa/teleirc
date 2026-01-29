package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
)

func TestGetFullUsername(t *testing.T) {
	user := &models.User{ID: 1, FirstName: "John", Username: "jsmith"}
	correct := user.FirstName + " (@" + user.Username + ")"
	name := GetFullUsername(false, user)

	assert.Equal(t, correct, name)
}

func TestGetFullUserZwsp(t *testing.T) {
	tests := []struct {
		name     string
		user     *models.User
		expected string
	}{
		{
			name:     "ascii",
			user:     &models.User{ID: 1, FirstName: "John", Username: "jsmith"},
			expected: "John (@j\u200bsmith)",
		},
		{
			name:     "cyrillic",
			user:     &models.User{ID: 1, FirstName: "Иван", Username: "иван"},
			expected: "Иван (@и\u200bван)",
		},
		{
			name:     "japanese",
			user:     &models.User{ID: 1, FirstName: "まこと", Username: "まこと"},
			expected: "まこと (@ま\u200bこと)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := GetFullUsername(true, test.user)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestGetFullNoUsername(t *testing.T) {
	user := &models.User{ID: 1, FirstName: "John"}
	correct := user.FirstName
	name := GetFullUsername(false, user)

	assert.Equal(t, correct, name)
}

func TestGetNoUsername(t *testing.T) {
	user := &models.User{ID: 1, FirstName: "John"}
	correct := user.FirstName
	name := GetFullUsername(false, user)

	assert.Equal(t, correct, name)
}

func TestGetUsername(t *testing.T) {
	user := &models.User{ID: 1, FirstName: "John", Username: "jsmith"}
	correct := user.Username
	name := GetUsername(false, user)

	assert.Equal(t, correct, name)
}

func TestZwspUsername(t *testing.T) {
	tests := []struct {
		name     string
		user     *models.User
		expected string
	}{
		{
			name:     "ascii",
			user:     &models.User{ID: 1, FirstName: "John", Username: "jsmith"},
			expected: "j\u200bsmith",
		},
		{
			name:     "cyrillic",
			user:     &models.User{ID: 1, FirstName: "Иван", Username: "иван"},
			expected: "и\u200bван",
		},
		{
			name:     "japanese",
			user:     &models.User{ID: 1, FirstName: "まこと", Username: "まこと"},
			expected: "ま\u200bこと",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := GetUsername(true, test.user)
			assert.Equal(t, test.expected, actual)
		})
	}
}
