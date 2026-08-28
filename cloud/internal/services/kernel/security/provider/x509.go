package provider

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
)

// ParseCSR parses and validates a PKCS#10 CSR.
//
// This function verifies the cryptographic signature of the CSR,
// therefore it proves that the CSR was generated using the
// corresponding private key.
func ParseCSR(
	csrPEM string,
) (*ParsedCSR, error) {

	if csrPEM == "" {
		return nil, errors.New(
			"csr is empty",
		)
	}

	block, _ := pem.Decode(
		[]byte(csrPEM),
	)

	if block == nil {
		return nil, errors.New(
			"invalid csr pem",
		)
	}

	if block.Type != "CERTIFICATE REQUEST" &&
		block.Type != "NEW CERTIFICATE REQUEST" {

		return nil, fmt.Errorf(
			"invalid csr pem type: %s",
			block.Type,
		)
	}

	csr, err :=
		x509.ParseCertificateRequest(
			block.Bytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"parse csr: %w",
			err,
		)
	}

	if err :=
		csr.CheckSignature(); err != nil {

		return nil, fmt.Errorf(
			"invalid csr signature: %w",
			err,
		)
	}

	result := &ParsedCSR{
		Subject: csr.Subject.String(),

		CommonName: csr.Subject.CommonName,

		DNSNames:
			append(
				[]string(nil),
				csr.DNSNames...,
			),

		EmailAddresses:
			append(
				[]string(nil),
				csr.EmailAddresses...,
			),
	}

	for _, ip := range csr.IPAddresses {
		result.IPAddresses =
			append(
				result.IPAddresses,
				ip.String(),
			)
	}

	for _, uri := range csr.URIs {
		if uri == nil {
			continue
		}

		result.URIs =
			append(
				result.URIs,
				uri.String(),
			)
	}

	switch key := csr.PublicKey.(type) {

	case *rsa.PublicKey:

		result.PublicKeyAlgorithm = "rsa"

		result.PublicKeySize =
			key.Size() * 8

	case *ecdsa.PublicKey:

		result.PublicKeyAlgorithm = "ecdsa"

		result.PublicKeySize =
			key.Curve.Params().BitSize

	default:

		result.PublicKeyAlgorithm =
			csr.PublicKeyAlgorithm.String()
	}

	return result, nil
}

// ParseCSRDER returns the DER encoded CSR.
func ParseCSRDER(
	csrPEM string,
) ([]byte, error) {

	block, _ :=
		pem.Decode(
			[]byte(csrPEM),
		)

	if block == nil {
		return nil, errors.New(
			"invalid csr pem",
		)
	}

	if block.Type != "CERTIFICATE REQUEST" &&
		block.Type != "NEW CERTIFICATE REQUEST" {

		return nil, fmt.Errorf(
			"invalid csr pem type: %s",
			block.Type,
		)
	}

	if _, err :=
		x509.ParseCertificateRequest(
			block.Bytes,
		); err != nil {

		return nil, fmt.Errorf(
			"parse csr: %w",
			err,
		)
	}

	return block.Bytes, nil
}

// ParseIssuedCertificate parses an issued certificate
// and converts it to the provider-independent model.
func ParseIssuedCertificate(certificateID uint,certificatePEM string) (*IssuedCertificate, error) {

	if certificatePEM == "" {
		return nil, errors.New(
			"certificate pem is empty",
		)
	}

	block, _ :=
		pem.Decode(
			[]byte(certificatePEM),
		)

	if block == nil {
		return nil, errors.New(
			"invalid certificate pem",
		)
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf(
			"invalid certificate pem type: %s",
			block.Type,
		)
	}

	cert, err :=
		x509.ParseCertificate(
			block.Bytes,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"parse certificate: %w",
			err,
		)
	}

	return &IssuedCertificate{
		CertificateID:
			certificateID,

		CertificatePEM:
			certificatePEM,

		Fingerprint:
			CertificateFingerprint(
				block.Bytes,
			),

		SerialNumber:
			cert.SerialNumber.String(),

		Subject:
			cert.Subject.String(),

		Issuer:
			cert.Issuer.String(),

		NotBefore:
			cert.NotBefore,

		NotAfter:
			cert.NotAfter,
	}, nil
}

// CertificateFingerprint returns SHA-256 fingerprint
// of DER encoded X.509 certificate.
func CertificateFingerprint(
	der []byte,
) string {

	sum :=
		sha256.Sum256(der)

	return hex.EncodeToString(
		sum[:],
	)
}