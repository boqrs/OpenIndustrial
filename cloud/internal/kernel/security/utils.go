package security

import(
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"crypto/subtle"
	"strings"
	"errors"
	"context"

	"github.com/google/uuid"

)

func verifySecret(secret string,expectedHash string) bool {
	actualHash := hashSecret(secret)

	if len(actualHash) != len(expectedHash) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(actualHash),[]byte(expectedHash)) == 1
}

func parseBootstrapToken(
	token string,
) (uuid.UUID, string, error) {

	parts := strings.SplitN(
		token,
		".",
		2,
	)

	if len(parts) != 2 {
		return uuid.Nil, "", ErrCredentialInvalid
	}

	id, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, "", ErrCredentialInvalid
	}

	if parts[1] == "" {
		return uuid.Nil, "", ErrCredentialInvalid
	}

	return id, parts[1], nil
}

func hashSecret(
	secret string,
) string {

	sum := sha256.Sum256(
		[]byte(secret),
	)

	return base64.RawURLEncoding.EncodeToString(
		sum[:],
	)
}

func generateSecret(
	size int,
) (string, error) {

	buf := make([]byte, size)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		buf,
	), nil
}

func validateCSRForResource(
	csr *ParsedCSR,
	resourceID uuid.UUID,
) error {

	if csr == nil {
		return errors.New(
			"csr is nil",
		)
	}

	expectedURI :=
		"urn:openindustrial:resource:" +
			resourceID.String()

	for _, uri := range csr.URIs {

		if uri == expectedURI {
			return nil
		}
	}

	return errors.New(
		"csr identity does not match resource",
	)
}

type resourceAuthContextKey struct{}

type ResourceAuthContext struct {
	ResourceID uuid.UUID

	CertificateID string
}

func withResourceAuth(
	ctx context.Context,
	auth ResourceAuthContext,
) context.Context {

	return context.WithValue(
		ctx,
		resourceAuthContextKey{},
		auth,
	)
}

func resourceIDFromContext(
	ctx context.Context,
) (uuid.UUID, bool) {

	value := ctx.Value(
		resourceAuthContextKey{},
	)

	auth, ok :=
		value.(ResourceAuthContext)

	if !ok {
		return uuid.Nil, false
	}

	return auth.ResourceID, true
}