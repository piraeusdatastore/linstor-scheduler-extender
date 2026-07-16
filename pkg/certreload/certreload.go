// Package certreload provides a TLS certificate source that reloads the
// certificate and key from disk when their contents change, so a long-running
// server keeps serving a valid certificate after rotation without being
// restarted.
package certreload

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Reloader loads a TLS certificate and key from disk and reloads them when the
// files change on disk. It is safe for concurrent use.
//
// Kubernetes Secret volumes (and cert-manager) rotate a mounted certificate by
// atomically swapping a symlink in the mount directory rather than rewriting the
// existing file in place. Reloader reads the files by path on every handshake,
// which follows those symlinks to the current target, and detects a rotation by
// comparing a hash of the file contents rather than an open file descriptor,
// inode, or file metadata such as size and mtime, none of which reliably changes
// on a symlink swap. Reading a small PEM file per TLS handshake is negligible
// for an admission webhook, whose handshakes are infrequent.
type Reloader struct {
	certFile string
	keyFile  string
	logger   logrus.FieldLogger

	mu       sync.Mutex
	cert     *tls.Certificate
	certHash [sha256.Size]byte
	keyHash  [sha256.Size]byte
}

// New returns a Reloader serving the certificate and key at certFile and
// keyFile. The pair is loaded once up front so that a misconfiguration fails
// fast at startup instead of on the first TLS handshake. A nil logger falls
// back to the logrus standard logger.
func New(certFile, keyFile string, logger logrus.FieldLogger) (*Reloader, error) {
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	r := &Reloader{
		certFile: certFile,
		keyFile:  keyFile,
		logger:   logger,
	}

	if _, err := r.GetCertificate(nil); err != nil {
		return nil, fmt.Errorf("loading initial TLS certificate: %w", err)
	}

	return r, nil
}

// GetCertificate returns the current certificate, reloading it from disk first
// if the files changed since the last successful load. It satisfies the
// tls.Config.GetCertificate callback signature; the ClientHelloInfo is ignored
// because a single certificate is served for every client.
//
// A reload failure (for example a partially written file observed mid-rotation)
// is logged and the last known-good certificate keeps being served, so a
// transient error never takes the webhook down. The cached hashes are left
// unchanged on failure, so the next handshake retries the reload until it
// succeeds.
func (r *Reloader) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	certPEM, keyPEM, err := r.readPair()
	if err != nil {
		if r.cert != nil {
			r.logger.Warnf("Keeping last known-good TLS certificate, cannot read certificate files: %v", err)
			return r.cert, nil
		}
		return nil, err
	}

	certHash := sha256.Sum256(certPEM)
	keyHash := sha256.Sum256(keyPEM)

	unchanged := r.cert != nil && certHash == r.certHash && keyHash == r.keyHash
	if unchanged {
		return r.cert, nil
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		if r.cert != nil {
			r.logger.Warnf("Keeping last known-good TLS certificate, reload failed: %v", err)
			return r.cert, nil
		}
		return nil, fmt.Errorf("loading TLS key pair: %w", err)
	}

	firstLoad := r.cert == nil
	r.cert = &cert
	r.certHash = certHash
	r.keyHash = keyHash

	if firstLoad {
		r.logger.Infof("Loaded TLS certificate from %s", r.certFile)
	} else {
		r.logger.Infof("Reloaded TLS certificate from %s after change on disk", r.certFile)
	}

	return r.cert, nil
}

func (r *Reloader) readPair() (certPEM, keyPEM []byte, err error) {
	certPEM, err = os.ReadFile(r.certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", r.certFile, err)
	}
	keyPEM, err = os.ReadFile(r.keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", r.keyFile, err)
	}
	return certPEM, keyPEM, nil
}
