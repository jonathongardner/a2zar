package xar

import (
	"crypto/x509"
	"encoding/base64"
	"strings"
)

func (xr *XarReader) readCertificates() error {
	for _, encCert := range xr.toc.Toc.Signature.KeyInfo.X509Data.X509Certificate {

		cb64 := []byte(strings.ReplaceAll(encCert, "\n", ""))
		cder := make([]byte, base64.StdEncoding.DecodedLen(len(cb64)))
		ndec, err := base64.StdEncoding.Decode(cder, cb64)
		if err != nil {
			return err
		}

		if cert, err := x509.ParseCertificate(cder[0:ndec]); err != nil {
			return err
		} else {
			xr.certs = append(xr.certs, cert)
		}
	}

	return nil
}

func (xr *XarReader) CheckCertificatesSignatures() error {
	// cache certificates
	if len(xr.certs) == 0 {
		if err := xr.readCertificates(); err != nil {
			return err
		}
	}

	// in order to verify certificates we check signatures in the chain
	// this assumes certificates are ordered
	for ii := 0; ii < len(xr.certs); ii++ {
		if ii == len(xr.certs)-1 {
			continue
		}
		if err := xr.certs[ii].CheckSignatureFrom(xr.certs[ii+1]); err != nil {
			return err
		}
	}

	return nil
}
