package jwt

import (
	"github.com/dimfu/finch/authentication/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExpiredAccessToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "bersama baaayanganmu kasiiih seakan akan ku terjaga dari mimpi mimpi dari kehidupan yang semu dan melenakanku...")
	user := models.User{ID: "1337"}
	token, err := Generate(user.ID, time.Now().Unix())
	assert.Equal(t, nil, err, "should not throw error on Generate")
	assert.NotEqual(t, "", token.AccessToken)
	assert.NotEqual(t, "", token.RefreshToken)

	_, err = ValidateAccessToken(token.AccessToken)
	assert.NotErrorIs(t, jwt.ErrTokenExpired, err, "should throw because token is expired")
}

func TestValidRefreshToken(t *testing.T) {
	os.Setenv("SECRET_KEY", "esok kan masih adaaaa.... huuuuuu uuu, esok kannn masih... adaaaaaaaaaaaaaaaaaaaa")
	user := models.User{ID: "1337"}
	token, err := Generate(user.ID, time.Now().Add(ACCESS_LIFETIME_DURATION).Unix())
	assert.Equal(t, nil, err, "should not throw error on Generate")
	_, err = ValidateRefreshToken(token.RefreshToken)
	assert.Equal(t, nil, err, "should not throw error on ValidateRefreshToken")
}
