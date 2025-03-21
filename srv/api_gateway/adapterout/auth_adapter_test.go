package adapterout

import (
	"context"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func createAuthAdapter(t *testing.T, nc *nats.Conn) portout.AuthenticationPortOut {
	brk, err := broker.NewNatsMessageBroker(nc, zaptest.NewLogger(t))
	require.NoError(t, err)
	return NewAuthenticationAdapter(brk, zaptest.NewLogger(t))
}

func startAuthMock(t *testing.T, nc *nats.Conn, issuer string) func() {
	app := fx.New(
		fx.Supply(zaptest.NewLogger(t)),
		fx.Supply(issuer),
		fx.Supply(nc),
		fx.Supply(DefaultUsers),
		fx.Provide(broker.NewNatsMessageBroker),
		fx.Provide(NewAuthMock),
		fx.Invoke(StartAuthMock),
	)

	// Start microservice, with timeout of 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	t.Cleanup(cancel)
	require.NoError(t, app.Start(ctx))

	// Stop microservice, with timeout of 1 second
	stopFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		t.Cleanup(cancel)
		require.NoError(t, app.Stop(ctx))
	}

	t.Cleanup(stopFunc)
	return stopFunc
}

func TestLogin(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	startAuthMock(t, nc, "test")

	token, err := aa.GetToken("admin")
	require.NoError(t, err)
	res, err := aa.VerifyToken(token)
	require.NoError(t, err)

	username, err := aa.GetUsername(res)
	require.NoError(t, err)
	require.Equal(t, "admin", username)

	role, err := aa.GetRole(res)
	require.NoError(t, err)
	require.Equal(t, types.RoleGlobalAdmin, role)
}

func TestGetTokenMarshalError(t *testing.T) {
	nc, cancel := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	startAuthMock(t, nc, "test")

	cancel()

	token, err := aa.GetToken(string([]byte{255}))
	require.Error(t, err)
	require.Zero(t, token)
}

func TestGetUsernameInvalidToken(t *testing.T) {
	t.Run("WrongType", func(t *testing.T) {
		nc, _ := broker.NewInProcessNATSServer(t)
		aa := createAuthAdapter(t, nc)
		startAuthMock(t, nc, "test")

		username, err := aa.GetUsername(42)
		require.Error(t, err)
		require.Zero(t, username)
	})
	t.Run("WrongJwtStructure", func(t *testing.T) {
		nc, _ := broker.NewInProcessNATSServer(t)
		aa := createAuthAdapter(t, nc)
		startAuthMock(t, nc, "test")

		token := jwt.NewWithClaims(
			jwt.SigningMethodES256,
			jwt.MapClaims{"sub": 42},
		)

		username, err := aa.GetUsername(token)
		require.Error(t, err)
		require.Zero(t, username)
	})
}

func TestGetRoleInvalidToken(t *testing.T) {
	t.Run("WrongType", func(t *testing.T) {
		nc, _ := broker.NewInProcessNATSServer(t)
		aa := createAuthAdapter(t, nc)
		startAuthMock(t, nc, "test")
		token, err := aa.GetRole(42)
		require.Error(t, err)
		require.Zero(t, token)
	})

	var cases = []struct {
		name   string
		claims jwt.MapClaims
	}{
		{"JwtMissingRole", jwt.MapClaims{"sub": 42}},
		{"JwtRoleWrongType", jwt.MapClaims{"role": 42}},
		{"JwtRoleNonExistent", jwt.MapClaims{"role": "nonexistent_role"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			nc, _ := broker.NewInProcessNATSServer(t)

			aa := createAuthAdapter(t, nc)
			startAuthMock(t, nc, "test")

			token := jwt.NewWithClaims(
				jwt.SigningMethodES256,
				c.claims,
			)

			role, err := aa.GetRole(token)
			require.Error(t, err)
			require.Zero(t, role)
		})
	}
}

func TestRevocation(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	stopEarly := startAuthMock(t, nc, "test")

	token, err := aa.GetToken("admin")
	require.NoError(t, err)
	_, err = aa.VerifyToken(token)
	require.NoError(t, err)

	stopEarly()
	startAuthMock(t, nc, "test")

	res, err := aa.VerifyToken(token)
	require.ErrorIs(t, err, portout.ErrTokenInvalid)
	require.Nil(t, res)
}

func TestGetKeyStreamCreateError(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(nc)
	require.NoError(t, err)
	_, err = js.CreateStream(
		context.Background(),
		jetstream.StreamConfig{Name: "auth_keys", Subjects: []string{"hello"}},
	)
	require.NoError(t, err)

	aa := createAuthAdapter(t, nc)
	casted, ok := aa.(*AuthenticationAdapter)
	require.True(t, ok)

	key, err := getValidationKey(context.Background(), casted, "test")
	require.Error(t, err)
	require.Zero(t, key)
}

func TestGetKeyWrongMessageFormat(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	js, err := jetstream.New(nc)
	require.NoError(t, err)
	_, err = js.CreateStream(
		context.Background(),
		jetstream.StreamConfig{Name: "auth_keys", Subjects: []string{"keys.>"}},
	)
	require.NoError(t, err)
	_, err = js.Publish(context.Background(), "keys.test", []byte("ciao"))
	require.NoError(t, err)

	aa := createAuthAdapter(t, nc)
	casted, ok := aa.(*AuthenticationAdapter)
	require.True(t, ok)

	key, err := getValidationKey(context.Background(), casted, "test")
	require.Error(t, err)
	require.Zero(t, key)
}

func TestGetKeyNoMessages(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)
	casted, ok := aa.(*AuthenticationAdapter)
	require.True(t, ok)

	key, err := getValidationKey(context.Background(), casted, "test")
	require.Error(t, err)
	require.Zero(t, key)
}

func TestVerifyTokenMissingIssuer(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)

	key, err := aa.VerifyToken("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA")
	require.ErrorIs(t, err, portout.ErrTokenInvalid)
	require.Zero(t, key)
}

func TestVerifyTokenGetKeyError(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)

	key, err := aa.VerifyToken("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiaXNzIjoiSm9obiBEb2UiLCJhZG1pbiI6dHJ1ZSwiaWF0IjoxNTE2MjM5MDIyfQ.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA")
	require.ErrorIs(t, err, portout.ErrTokenInvalid)
	require.Zero(t, key)
}
func TestVerifyTokenExpired(t *testing.T) {
	nc, _ := broker.NewInProcessNATSServer(t)

	aa := createAuthAdapter(t, nc)

	js, err := jetstream.New(nc)
	require.NoError(t, err)
	_, err = js.CreateStream(
		context.Background(),
		jetstream.StreamConfig{Name: "auth_keys", Subjects: []string{"keys.>"}},
	)
	require.NoError(t, err)
	_, err = js.Publish(context.Background(), "keys.test", []byte("{\"crv\":\"P-256\",\"kty\":\"EC\",\"x\":\"EpsYjz0Av4mjECTPdzSWrHxOKV8-fuMIy5cB-_LdhpU\",\"y\":\"lYpNIOYMsGP1Di0cKCqBpqPkhFvtPtrEABxlQq01i5k\"}"))
	require.NoError(t, err)

	key, err := aa.VerifyToken("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEzNTUzMTA3MzIsImlzcyI6InRlc3QiLCJyb2xlIjoibG9jYWxfYWRtaW4iLCJzdWIiOiJ0ZXN0In0.Q_6QLbT6-Fvlc-5Ss5vkB0UCYDUCS7oOsPxsqSrlIhJ1j-_Cov9YAWa0Zm8vuSvdAPlX7yGqd0gMDA3mJv3wJQ")
	require.ErrorIs(t, err, portout.ErrTokenExpired)
	require.Zero(t, key)
}
