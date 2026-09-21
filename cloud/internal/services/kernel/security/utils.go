package security

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

func verifySecret(
	secret string,
	expectedHash string,
) bool {
	actualHash := hashSecret(secret)

	if len(actualHash) != len(expectedHash) {
		return false
	}

	return subtle.ConstantTimeCompare(
		[]byte(actualHash),
		[]byte(expectedHash),
	) == 1
}

func formatBootstrapToken(
	credentialID uint,
	secret string,
) string {
	return strconv.FormatUint(
		uint64(credentialID),
		10,
	) + "." + secret
}

func parseBootstrapToken(
	token string,
) (uint, string, error) {
	parts := strings.SplitN(
		token,
		".",
		2,
	)

	if len(parts) != 2 {
		return 0, "", ErrCredentialInvalid
	}

	if parts[0] == "" {
		return 0, "", ErrCredentialInvalid
	}

	if parts[1] == "" {
		return 0, "", ErrCredentialInvalid
	}

	id64, err := strconv.ParseUint(
		parts[0],
		10,
		64,
	)
	if err != nil {
		return 0, "", ErrCredentialInvalid
	}

	if id64 == 0 {
		return 0, "", ErrCredentialInvalid
	}

	return uint(id64), parts[1], nil
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
	resourceID uint,
) error {
	if csr == nil {
		return errors.New(
			"csr is nil",
		)
	}

	expectedURI :=
		"urn:openindustrial:resource:" +
			strconv.Itoa(int(resourceID))

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
	ResourceID    uint
	CertificateID uint
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
) (uint, bool) {
	value := ctx.Value(
		resourceAuthContextKey{},
	)

	auth, ok :=
		value.(ResourceAuthContext)

	if !ok {
		return 0, false
	}

	return auth.ResourceID, true
}
