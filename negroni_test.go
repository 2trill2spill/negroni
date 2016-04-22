package negroni

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func writeTestCertAndKey() {

	const key string = "-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIIEowIBAAKCAQEAucMkMlmXRu30mWGbEpvK1/AOYIqeYzPeTQ/13LU0iH/1vp5w\n" +
		"LyRnUDm/MCUZtELndlia2/QrZlPCmh2RDJcQuW60SMk97wXNZlWJPHBaZlMo//PC\n" +
		"nguJRr/uYeMjqLOD1WaRq1KI5IKj3OVDU4mgGUG0cmPW44rM8Gari97ImEWtL+bi\n" +
		"Gux0wGt4n4+fbcuYho2QB1lIvdgf7qLOOk1OzLYg6+le7zQo/ejgGjww/V9pTHa9\n" +
		"Kupfydpid954N4rlxHjkv1JUpbtVyNNFBis7jYeyCpeKYHoZs7vzyBjX6TupaUZh\n" +
		"jAyUNX4vKOUISwLOJMcBMFWHLPiKn+z/rMRWxwIDAQABAoIBABjWofFlj1e5m5oi\n" +
		"tfSABlhnTdTU1CbZmaDeQHa5eAPEcFyDnOaAzJuqjQjn/Q7gX0iwwkwk0uTE0GGB\n" +
		"EJMJZAAXcF3HItPiB7vmNgpxz3SdK/9FXcF6M2nVYg+Oqob1MeyH4XRx/cHoQfbY\n" +
		"ZO83mMOnyHP/K42TUqHWaJu0N0g81iPJLbn9heQkm/7PMDWRycLT4QPlZVmckc4U\n" +
		"CWYieAktRoNaqwrw6cwQq/IDtw74GSnn4e6KuuO3dGEn/0pP7C9VImMSkaGI8vcc\n" +
		"RUNH6dM3lwcN17DWIbOucAqHENucnOP2LdQ8uFVe0gP5Q5Y81BQgL+Gm45JQIG04\n" +
		"/oTjktkCgYEA89l1PaNrb2GkQMeOtIrIFRskPi2GMxzLzP4wKs1pfFOw7RLRe7rj\n" +
		"HNgEFDjQdtFMEoUujRxlDoolIqcLNSMrcs858sW/GvsqvVFLpDwbzSy12jkkuQ9X\n" +
		"xTJUWEbJChlVTStcA7OSxlHWLDkEb/T5aZ11PaKfKXVSGLUhwjuiuLUCgYEAwwS5\n" +
		"Se2yC4tJCjqYNwZcx+hl9vk3oYiRwpq9Q/qjjAUA41EKJT0j8vdmprNcVFeg7yDD\n" +
		"JzFLhR6Rz2mhnWZAs6E3xE7K+UayAGDh28F3XtylBaK9UfRTxfKpjQubQHTMMwig\n" +
		"qptuqDrpyh9P8QOryFpwd6G7ASmQOI91n0DAKwsCgYAqPvBa32cLguUL/AazLKAB\n" +
		"WBRtWUG4tTJxr0/0+mVeDrnGOM5mGzihlKMQRc+H5jbBtqUb+WFgpXpNiJcee5tZ\n" +
		"ZqFpd+zl5cG/zsfGCkvevfI7fk7oaMoR9eg66viFcWIf3nUwhvnUtfTe8HneU5iq\n" +
		"PYdESFo+un6gnTDeD6rfSQKBgAVxS06B4LcrwvUTH45hkhNOLBJRcDkE27SHpwKP\n" +
		"qyLMPPMbHuJdK3SrbOT4GnqMG4Sw8GFWodnsOXuFYipHYUTiGfFTlZyRWYRjL8p3\n" +
		"QsKV+9EFcq0n2XKrzAmQqluJJu3BruI0BkmRo1atuwhp5tBnnb3o+JQWudyqPEke\n" +
		"gH5zAoGBAOMD3KhTNwjwCvXhdub+mqJapQSSXWRwLQ7foolWPAkjyjSh/eqMaCwk\n" +
		"dTatX0FEYhC683w9PID6m/RlcBQZQbu9MARTmr+P5FzWci235/PBZK2x5rR6aLFT\n" +
		"Xr3Ya9Cr6wW+1IVGIuqZweulJ6kyCaOs+L3/MstmwDxlSBHivTpA\n" +
		"-----END RSA PRIVATE KEY-----\n"

	const cert string = "-----BEGIN CERTIFICATE-----\n" +
		"MIID9DCCAtygAwIBAgIJAO+0au5QKhtvMA0GCSqGSIb3DQEBDQUAMFkxCzAJBgNV\n" +
		"BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX\n" +
		"aWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMTCWxvY2FsaG9zdDAeFw0xNjA0MjIxOTM0\n" +
		"MzRaFw0xNjA3MjExOTM0MzRaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21l\n" +
		"LVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNV\n" +
		"BAMTCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALnD\n" +
		"JDJZl0bt9JlhmxKbytfwDmCKnmMz3k0P9dy1NIh/9b6ecC8kZ1A5vzAlGbRC53ZY\n" +
		"mtv0K2ZTwpodkQyXELlutEjJPe8FzWZViTxwWmZTKP/zwp4LiUa/7mHjI6izg9Vm\n" +
		"katSiOSCo9zlQ1OJoBlBtHJj1uOKzPBmq4veyJhFrS/m4hrsdMBreJ+Pn23LmIaN\n" +
		"kAdZSL3YH+6izjpNTsy2IOvpXu80KP3o4Bo8MP1faUx2vSrqX8naYnfeeDeK5cR4\n" +
		"5L9SVKW7VcjTRQYrO42HsgqXimB6GbO788gY1+k7qWlGYYwMlDV+LyjlCEsCziTH\n" +
		"ATBVhyz4ip/s/6zEVscCAwEAAaOBvjCBuzAdBgNVHQ4EFgQU0WpFAZ5jAx8Lmb/r\n" +
		"WIkfjGmFFn4wgYsGA1UdIwSBgzCBgIAU0WpFAZ5jAx8Lmb/rWIkfjGmFFn6hXaRb\n" +
		"MFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJ\n" +
		"bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMTCWxvY2FsaG9zdIIJAO+0\n" +
		"au5QKhtvMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQENBQADggEBAHqstZ9uH1JY\n" +
		"m51xkqVinfPkPY00eqe0Mv9T0omXetM9VQeWBi+XB7cHZxvWfk6j+YfdijYP9nW+\n" +
		"gSP9Lht8WvE82yR3HJSZWC6yUstUE1sMKYCL+Kp6iTxy1l1md9jvO0YJoIVli7W+\n" +
		"YKEv3ScsvQ4Z7FKtL1bZCjsBaR7Im0kDQMCfoH2rS1FZe6Ig1+ZNaqnMHmmd2ptx\n" +
		"wLPcB9HyN/nWt7+Dw3S6JQyESYpTnZAcYD2VM7RU7ei83U8nxqboZSYnRtqn9PSC\n" +
		"+j1WcoHBl7zqgOxbxSxzeEq14XZSSUFvDcGuSPoLw3WW910audlJFXqQZiGWcizM\n" +
		"KNmj2VATsFQ=\n" +
		"-----END CERTIFICATE-----\n"

	ioutil.WriteFile("/tmp/cert.pem", []byte(cert), 0644)
	ioutil.WriteFile("/tmp/key.pem", []byte(key), 0644)
}

func TestNegroniRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":3000")
}

func TestNegroniRunTLS(t *testing.T) {

	/* Write a test cert and test key to disk so
	   we can test RunTLS. */
	writeTestCertAndKey()

	// just test that RunTLS doesn't bomb
	go New().RunTLS(":3001", "/tmp/cert.pem", "/tmp/key.pem")
}

func TestNegroniServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n := New()
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "foo"
		next(rw, r)
		result += "ban"
	}))
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bar"
		next(rw, r)
		result += "baz"
	}))
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bat"
		rw.WriteHeader(http.StatusBadRequest)
	}))

	n.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}

// Ensures that a Negroni middleware chain
// can correctly return all of its handlers.
func TestHandlers(t *testing.T) {
	response := httptest.NewRecorder()
	n := New()
	handlers := n.Handlers()
	expect(t, 0, len(handlers))

	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.WriteHeader(http.StatusOK)
	}))

	// Expects the length of handlers to be exactly 1
	// after adding exactly one handler to the middleware chain
	handlers = n.Handlers()
	expect(t, 1, len(handlers))

	// Ensures that the first handler that is in sequence behaves
	// exactly the same as the one that was registered earlier
	handlers[0].ServeHTTP(response, (*http.Request)(nil), nil)
	expect(t, response.Code, http.StatusOK)
}
