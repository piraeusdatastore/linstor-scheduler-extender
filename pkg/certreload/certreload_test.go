package certreload_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/piraeusdatastore/linstor-scheduler-extender/pkg/certreload"
)

// generateCertPEM returns a fresh self-signed certificate/key pair encoded as
// PEM, carrying the given serial so tests can tell two certificates apart.
func generateCertPEM(t *testing.T, serial int64) (certPEM, keyPEM []byte) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(serial),
		Subject:      pkix.Name{CommonName: "linstor-scheduler-admission-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("creating certificate: %v", err)
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshaling key: %v", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

	return certPEM, keyPEM
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func symlink(t *testing.T, target, link string) {
	t.Helper()

	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink %s -> %s: %v", link, target, err)
	}
}

func leafSerial(t *testing.T, cert *tls.Certificate) *big.Int {
	t.Helper()

	if cert == nil || len(cert.Certificate) == 0 {
		t.Fatal("certificate is empty")
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parsing leaf: %v", err)
	}
	return leaf.SerialNumber
}

func discardLogger() logrus.FieldLogger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestReloaderReloadsAfterRotation(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	got, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate before rotation: %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 1 {
		t.Fatalf("expected serial 1 before rotation, got %d", serial)
	}

	// Rotate: overwrite both files with a different pair.
	certB, keyB := generateCertPEM(t, 2)
	writeFile(t, certPath, certB)
	writeFile(t, keyPath, keyB)

	got, err = r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate after rotation: %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 2 {
		t.Fatalf("expected serial 2 after rotation, got %d", serial)
	}
}

// TestReloaderServesRotatedCertOverHandshake drives real TLS handshakes over a
// loopback listener to prove that crypto/tls actually calls GetCertificate and
// that a rotation reaches a freshly connecting client. Each handshake uses a new
// client, so session resumption cannot mask the rotated certificate.
func TestReloaderServesRotatedCertOverHandshake(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// A real loopback listener is used rather than net.Pipe: an unbuffered
	// synchronous pipe can deadlock a multi-record TLS flight, whereas a TCP
	// socket's send buffer absorbs it.
	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{GetCertificate: r.GetCertificate})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = ln.Close() }()

	handshakeSerial := func(t *testing.T) *big.Int {
		t.Helper()

		errc := make(chan error, 1)
		go func() {
			conn, err := ln.Accept()
			if err != nil {
				errc <- err
				return
			}
			defer func() { _ = conn.Close() }()
			errc <- conn.(*tls.Conn).Handshake()
		}()

		// InsecureSkipVerify is fine here: the test asserts on the presented
		// leaf, not on trust, and the certificates are ephemeral self-signed.
		client, err := tls.Dial("tcp", ln.Addr().String(), &tls.Config{InsecureSkipVerify: true}) //nolint:gosec
		if err != nil {
			t.Fatalf("dial: %v", err)
		}
		defer func() { _ = client.Close() }()

		if err := <-errc; err != nil {
			t.Fatalf("server handshake: %v", err)
		}

		state := client.ConnectionState()
		if len(state.PeerCertificates) == 0 {
			t.Fatal("client received no peer certificate")
		}
		return state.PeerCertificates[0].SerialNumber
	}

	if serial := handshakeSerial(t); serial.Int64() != 1 {
		t.Fatalf("expected serial 1 over handshake, got %d", serial)
	}

	certB, keyB := generateCertPEM(t, 2)
	writeFile(t, certPath, certB)
	writeFile(t, keyPath, keyB)

	if serial := handshakeSerial(t); serial.Int64() != 2 {
		t.Fatalf("expected serial 2 over handshake after rotation, got %d", serial)
	}
}

// TestReloaderFollowsSymlinkSwap reproduces the directory layout kubelet
// creates for a mounted Secret and rotates it exactly as kubelet does: a fresh
// data directory plus an atomic rename of the ..data symlink. It guards against
// reading the files with os.Lstat/os.Readlink (which would pin the stale
// target) instead of following the symlink to the current contents.
func TestReloaderFollowsSymlinkSwap(t *testing.T) {
	dir := t.TempDir()

	// writeVersion lays down <dir>/<name>/{tls.crt,tls.key}, mimicking one of
	// kubelet's timestamped ..data_<ts> directories, and returns its base name.
	writeVersion := func(name string, certPEM, keyPEM []byte) {
		vdir := filepath.Join(dir, name)
		if err := os.Mkdir(vdir, 0o700); err != nil {
			t.Fatalf("mkdir %s: %v", vdir, err)
		}
		writeFile(t, filepath.Join(vdir, "tls.crt"), certPEM)
		writeFile(t, filepath.Join(vdir, "tls.key"), keyPEM)
	}

	certA, keyA := generateCertPEM(t, 1)
	writeVersion("..data_1", certA, keyA)

	// Relative targets, exactly as kubelet writes them:
	//   <dir>/..data  -> ..data_1
	//   <dir>/tls.crt -> ..data/tls.crt
	//   <dir>/tls.key -> ..data/tls.key
	symlink(t, "..data_1", filepath.Join(dir, "..data"))
	symlink(t, filepath.Join("..data", "tls.crt"), filepath.Join(dir, "tls.crt"))
	symlink(t, filepath.Join("..data", "tls.key"), filepath.Join(dir, "tls.key"))

	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	got, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate before swap: %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 1 {
		t.Fatalf("expected serial 1 before swap, got %d", serial)
	}

	// Rotate: new data directory, then atomically replace the ..data symlink by
	// renaming a temporary symlink over it, exactly as kubelet does.
	certB, keyB := generateCertPEM(t, 2)
	writeVersion("..data_2", certB, keyB)
	tmp := filepath.Join(dir, "..data_tmp")
	symlink(t, "..data_2", tmp)
	if err := os.Rename(tmp, filepath.Join(dir, "..data")); err != nil {
		t.Fatalf("atomic ..data swap: %v", err)
	}

	got, err = r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate after swap: %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 2 {
		t.Fatalf("expected serial 2 after symlink swap, got %d", serial)
	}
}

func TestReloaderReturnsSameCertWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	first, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate first: %v", err)
	}
	second, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate second: %v", err)
	}

	// With the files unchanged, the reloader must return the cached certificate
	// (the same pointer) instead of re-parsing and allocating a new one.
	if first != second {
		t.Fatal("expected the same cached certificate when files are unchanged")
	}
}

func TestNewNilLoggerDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, nil)
	if err != nil {
		t.Fatalf("New with nil logger: %v", err)
	}
	if _, err := r.GetCertificate(nil); err != nil {
		t.Fatalf("GetCertificate with nil logger: %v", err)
	}
}

func TestReloaderKeepsLastGoodOnParseError(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Simulate a partially written / corrupt file appearing during rotation.
	writeFile(t, certPath, []byte("-----BEGIN CERTIFICATE-----\nnot a certificate\n"))

	got, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("expected no error when reload fails with a last-good cert, got %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 1 {
		t.Fatalf("expected last known-good serial 1 after a bad reload, got %d", serial)
	}

	// Once a valid pair reappears, the next call must pick it up.
	certB, keyB := generateCertPEM(t, 2)
	writeFile(t, certPath, certB)
	writeFile(t, keyPath, keyB)

	got, err = r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("GetCertificate after recovery: %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 2 {
		t.Fatalf("expected serial 2 after recovery, got %d", serial)
	}
}

func TestReloaderKeepsLastGoodOnReadError(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := os.Remove(keyPath); err != nil {
		t.Fatalf("removing key: %v", err)
	}

	got, err := r.GetCertificate(nil)
	if err != nil {
		t.Fatalf("expected no error when a file is temporarily missing, got %v", err)
	}
	if serial := leafSerial(t, got); serial.Int64() != 1 {
		t.Fatalf("expected last known-good serial 1 when a file is missing, got %d", serial)
	}
}

func TestNewFailsWhenCertificateMissing(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	if _, err := certreload.New(certPath, keyPath, discardLogger()); err == nil {
		t.Fatal("expected New to fail when the certificate files are absent")
	}
}

func TestNewFailsWhenCertificateInvalid(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	// Both files exist but are not a valid PEM key pair, so the initial load
	// must fail fast rather than start serving with no certificate.
	writeFile(t, certPath, []byte("not a certificate"))
	writeFile(t, keyPath, []byte("not a key"))

	if _, err := certreload.New(certPath, keyPath, discardLogger()); err == nil {
		t.Fatal("expected New to fail when the files are not a valid key pair")
	}
}

// TestReloaderConcurrent exercises the mutex under the race detector: many
// readers call GetCertificate while the files are rotated underneath them.
func TestReloaderConcurrent(t *testing.T) {
	dir := t.TempDir()
	certPath := filepath.Join(dir, "tls.crt")
	keyPath := filepath.Join(dir, "tls.key")

	certA, keyA := generateCertPEM(t, 1)
	writeFile(t, certPath, certA)
	writeFile(t, keyPath, keyA)

	r, err := certreload.New(certPath, keyPath, discardLogger())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				if _, err := r.GetCertificate(nil); err != nil {
					t.Errorf("GetCertificate: %v", err)
					return
				}
			}
		}()
	}

	certB, keyB := generateCertPEM(t, 2)
	writeFile(t, certPath, certB)
	writeFile(t, keyPath, keyB)

	wg.Wait()
}
