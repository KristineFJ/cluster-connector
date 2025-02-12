/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package transport

import (
	"crypto/tls"
	"testing"
	"time"

	"k8s.io/client-go/transport"
)

const (
	rootCACert = `-----BEGIN CERTIFICATE-----
MIIC4DCCAcqgAwIBAgIBATALBgkqhkiG9w0BAQswIzEhMB8GA1UEAwwYMTAuMTMu
MTI5LjEwNkAxNDIxMzU5MDU4MB4XDTE1MDExNTIxNTczN1oXDTE2MDExNTIxNTcz
OFowIzEhMB8GA1UEAwwYMTAuMTMuMTI5LjEwNkAxNDIxMzU5MDU4MIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAunDRXGwsiYWGFDlWH6kjGun+PshDGeZX
xtx9lUnL8pIRWH3wX6f13PO9sktaOWW0T0mlo6k2bMlSLlSZgG9H6og0W6gLS3vq
s4VavZ6DbXIwemZG2vbRwsvR+t4G6Nbwelm6F8RFnA1Fwt428pavmNQ/wgYzo+T1
1eS+HiN4ACnSoDSx3QRWcgBkB1g6VReofVjx63i0J+w8Q/41L9GUuLqquFxu6ZnH
60vTB55lHgFiDLjA1FkEz2dGvGh/wtnFlRvjaPC54JH2K1mPYAUXTreoeJtLJKX0
ycoiyB24+zGCniUmgIsmQWRPaOPircexCp1BOeze82BT1LCZNTVaxQIDAQABoyMw
ITAOBgNVHQ8BAf8EBAMCAKQwDwYDVR0TAQH/BAUwAwEB/zALBgkqhkiG9w0BAQsD
ggEBADMxsUuAFlsYDpF4fRCzXXwrhbtj4oQwcHpbu+rnOPHCZupiafzZpDu+rw4x
YGPnCb594bRTQn4pAu3Ac18NbLD5pV3uioAkv8oPkgr8aUhXqiv7KdDiaWm6sbAL
EHiXVBBAFvQws10HMqMoKtO8f1XDNAUkWduakR/U6yMgvOPwS7xl0eUTqyRB6zGb
K55q2dejiFWaFqB/y78txzvz6UlOZKE44g2JAVoJVM6kGaxh33q8/FmrL4kuN3ut
W+MmJCVDvd4eEqPwbp7146ZWTqpIJ8lvA6wuChtqV8lhAPka2hD/LMqY8iXNmfXD
uml0obOEy+ON91k+SWTJ3ggmF/U=
-----END CERTIFICATE-----`

	certData = `-----BEGIN CERTIFICATE-----
MIIC6jCCAdSgAwIBAgIBCzALBgkqhkiG9w0BAQswIzEhMB8GA1UEAwwYMTAuMTMu
MTI5LjEwNkAxNDIxMzU5MDU4MB4XDTE1MDExNTIyMDEzMVoXDTE2MDExNTIyMDEz
MlowGzEZMBcGA1UEAxMQb3BlbnNoaWZ0LWNsaWVudDCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAKtdhz0+uCLXw5cSYns9rU/XifFSpb/x24WDdrm72S/v
b9BPYsAStiP148buylr1SOuNi8sTAZmlVDDIpIVwMLff+o2rKYDicn9fjbrTxTOj
lI4pHJBH+JU3AJ0tbajupioh70jwFS0oYpwtneg2zcnE2Z4l6mhrj2okrc5Q1/X2
I2HChtIU4JYTisObtin10QKJX01CLfYXJLa8upWzKZ4/GOcHG+eAV3jXWoXidtjb
1Usw70amoTZ6mIVCkiu1QwCoa8+ycojGfZhvqMsAp1536ZcCul+Na+AbCv4zKS7F
kQQaImVrXdUiFansIoofGlw/JNuoKK6ssVpS5Ic3pgcCAwEAAaM1MDMwDgYDVR0P
AQH/BAQDAgCgMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCwYJ
KoZIhvcNAQELA4IBAQCKLREH7bXtXtZ+8vI6cjD7W3QikiArGqbl36bAhhWsJLp/
p/ndKz39iFNaiZ3GlwIURWOOKx3y3GA0x9m8FR+Llthf0EQ8sUjnwaknWs0Y6DQ3
jjPFZOpV3KPCFrdMJ3++E3MgwFC/Ih/N2ebFX9EcV9Vcc6oVWMdwT0fsrhu683rq
6GSR/3iVX1G/pmOiuaR0fNUaCyCfYrnI4zHBDgSfnlm3vIvN2lrsR/DQBakNL8DJ
HBgKxMGeUPoneBv+c8DMXIL0EhaFXRlBv9QW45/GiAIOuyFJ0i6hCtGZpJjq4OpQ
BRjCI+izPzFTjsxD4aORE+WOkyWFCGPWKfNejfw0
-----END CERTIFICATE-----`

	keyData = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAq12HPT64ItfDlxJiez2tT9eJ8VKlv/HbhYN2ubvZL+9v0E9i
wBK2I/Xjxu7KWvVI642LyxMBmaVUMMikhXAwt9/6jaspgOJyf1+NutPFM6OUjikc
kEf4lTcAnS1tqO6mKiHvSPAVLShinC2d6DbNycTZniXqaGuPaiStzlDX9fYjYcKG
0hTglhOKw5u2KfXRAolfTUIt9hcktry6lbMpnj8Y5wcb54BXeNdaheJ22NvVSzDv
RqahNnqYhUKSK7VDAKhrz7JyiMZ9mG+oywCnXnfplwK6X41r4BsK/jMpLsWRBBoi
ZWtd1SIVqewiih8aXD8k26gorqyxWlLkhzemBwIDAQABAoIBAD2XYRs3JrGHQUpU
FkdbVKZkvrSY0vAZOqBTLuH0zUv4UATb8487anGkWBjRDLQCgxH+jucPTrztekQK
aW94clo0S3aNtV4YhbSYIHWs1a0It0UdK6ID7CmdWkAj6s0T8W8lQT7C46mWYVLm
5mFnCTHi6aB42jZrqmEpC7sivWwuU0xqj3Ml8kkxQCGmyc9JjmCB4OrFFC8NNt6M
ObvQkUI6Z3nO4phTbpxkE1/9dT0MmPIF7GhHVzJMS+EyyRYUDllZ0wvVSOM3qZT0
JMUaBerkNwm9foKJ1+dv2nMKZZbJajv7suUDCfU44mVeaEO+4kmTKSGCGjjTBGkr
7L1ySDECgYEA5ElIMhpdBzIivCuBIH8LlUeuzd93pqssO1G2Xg0jHtfM4tz7fyeI
cr90dc8gpli24dkSxzLeg3Tn3wIj/Bu64m2TpZPZEIlukYvgdgArmRIPQVxerYey
OkrfTNkxU1HXsYjLCdGcGXs5lmb+K/kuTcFxaMOs7jZi7La+jEONwf8CgYEAwCs/
rUOOA0klDsWWisbivOiNPII79c9McZCNBqncCBfMUoiGe8uWDEO4TFHN60vFuVk9
8PkwpCfvaBUX+ajvbafIfHxsnfk1M04WLGCeqQ/ym5Q4sQoQOcC1b1y9qc/xEWfg
nIUuia0ukYRpl7qQa3tNg+BNFyjypW8zukUAC/kCgYB1/Kojuxx5q5/oQVPrx73k
2bevD+B3c+DYh9MJqSCNwFtUpYIWpggPxoQan4LwdsmO0PKzocb/ilyNFj4i/vII
NToqSc/WjDFpaDIKyuu9oWfhECye45NqLWhb/6VOuu4QA/Nsj7luMhIBehnEAHW+
GkzTKM8oD1PxpEG3nPKXYQKBgQC6AuMPRt3XBl1NkCrpSBy/uObFlFaP2Enpf39S
3OZ0Gv0XQrnSaL1kP8TMcz68rMrGX8DaWYsgytstR4W+jyy7WvZwsUu+GjTJ5aMG
77uEcEBpIi9CBzivfn7hPccE8ZgqPf+n4i6q66yxBJflW5xhvafJqDtW2LcPNbW/
bvzdmQKBgExALRUXpq+5dbmkdXBHtvXdRDZ6rVmrnjy4nI5bPw+1GqQqk6uAR6B/
F6NmLCQOO4PDG/cuatNHIr2FrwTmGdEL6ObLUGWn9Oer9gJhHVqqsY5I4sEPo4XX
stR0Yiw0buV6DL/moUO0HIM9Bjh96HJp+LxiIS6UCdIhMPp5HoQa
-----END RSA PRIVATE KEY-----`
)

func TestNew(t *testing.T) {
	testCases := map[string]struct {
		Config       *transport.Config
		Err          bool
		TLS          bool
		TLSCert      bool
		TLSErr       bool
		Default      bool
		Insecure     bool
		DefaultRoots bool
	}{
		"default transport": {
			Default: true,
			Config:  &transport.Config{},
		},

		"insecure": {
			TLS:          true,
			Insecure:     true,
			DefaultRoots: true,
			Config: &transport.Config{TLS: transport.TLSConfig{
				Insecure: true,
			}},
		},

		"server name": {
			TLS:          true,
			DefaultRoots: true,
			Config: &transport.Config{TLS: transport.TLSConfig{
				ServerName: "foo",
			}},
		},

		"ca transport": {
			TLS: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData: []byte(rootCACert),
				},
			},
		},
		"bad ca file transport": {
			Err: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAFile: "invalid file",
				},
			},
		},
		//"bad ca data transport": {
		//	Err: true,
		//	Config: &transport.Config{
		//		TLS: transport.TLSConfig{
		//			CAData: []byte(rootCACert + "this is not valid"),
		//		},
		//	},
		//},
		"ca data overriding bad ca file transport": {
			TLS: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData: []byte(rootCACert),
					CAFile: "invalid file",
				},
			},
		},

		"cert transport": {
			TLS:     true,
			TLSCert: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData:   []byte(rootCACert),
					CertData: []byte(certData),
					KeyData:  []byte(keyData),
				},
			},
		},
		"bad cert data transport": {
			Err: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData:   []byte(rootCACert),
					CertData: []byte(certData),
					KeyData:  []byte("bad key data"),
				},
			},
		},
		"bad file cert transport": {
			Err: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData:   []byte(rootCACert),
					CertData: []byte(certData),
					KeyFile:  "invalid file",
				},
			},
		},
		"key data overriding bad file cert transport": {
			TLS:     true,
			TLSCert: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData:   []byte(rootCACert),
					CertData: []byte(certData),
					KeyData:  []byte(keyData),
					KeyFile:  "invalid file",
				},
			},
		},
		"callback cert and key": {
			TLS:     true,
			TLSCert: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData: []byte(rootCACert),
					GetCert: func() (*tls.Certificate, error) {
						crt, err := tls.X509KeyPair([]byte(certData), []byte(keyData))
						return &crt, err
					},
				},
			},
		},
		//"cert callback error": {
		//	TLS:     true,
		//	TLSCert: true,
		//	TLSErr:  true,
		//	Config: &transport.Config{
		//		TLS: transport.TLSConfig{
		//			CAData: []byte(rootCACert),
		//			GetCert: func() (*tls.Certificate, error) {
		//				return nil, errors.New("GetCert failure")
		//			},
		//		},
		//	},
		//},
		"cert data overrides empty callback result": {
			TLS:     true,
			TLSCert: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData: []byte(rootCACert),
					GetCert: func() (*tls.Certificate, error) {
						return nil, nil
					},
					CertData: []byte(certData),
					KeyData:  []byte(keyData),
				},
			},
		},
		"callback returns nothing": {
			TLS:     true,
			TLSCert: true,
			Config: &transport.Config{
				TLS: transport.TLSConfig{
					CAData: []byte(rootCACert),
					GetCert: func() (*tls.Certificate, error) {
						return nil, nil
					},
				},
			},
		},
	}
	for k, testCase := range testCases {
		t.Run(k, func(t *testing.T) {
			rt, err := New(testCase.Config, nil, "", 5*time.Second)
			switch {
			case testCase.Err && err == nil:
				t.Fatal("unexpected non-error")
			case !testCase.Err && err != nil:
				t.Fatalf("unexpected error: %v", err)
			}
			if testCase.Err {
				return
			}

			//switch {
			//case testCase.Default && rt != http.DefaultTransport:
			//	t.Fatalf("got %#v, expected the default transport", rt)
			//case !testCase.Default && rt == http.DefaultTransport:
			//	t.Fatalf("got %#v, expected non-default transport", rt)
			//}

			// We only know how to check transport.TLSConfig on http.Transports
			transport := rt.(*NatsTransport)
			switch {
			case testCase.TLS && transport.TLS == nil:
				t.Fatalf("got %#v, expected TLSClientConfig", transport)
			case !testCase.TLS && transport.TLS != nil:
				t.Fatalf("got %#v, expected no TLSClientConfig", transport)
			}
			if !testCase.TLS {
				return
			}

			switch {
			case testCase.DefaultRoots && transport.TLS.CAData != nil:
				t.Fatalf("got %#v, expected nil root CAs", transport.TLS.CAData)
			case !testCase.DefaultRoots && transport.TLS.CAData == nil:
				t.Fatalf("got %#v, expected non-nil root CAs", transport.TLS.CAData)
			}

			switch {
			case testCase.Insecure != transport.TLS.Insecure:
				t.Fatalf("got %#v, expected %#v", transport.TLS.Insecure, testCase.Insecure)
			}

			//switch {
			//case testCase.TLSCert && transport.TLS.GetClientCertificate == nil:
			//	t.Fatalf("got %#v, expected TLSClientConfig.GetClientCertificate", transport.TLSClientConfig)
			//case !testCase.TLSCert && transport.TLS.GetClientCertificate != nil:
			//	t.Fatalf("got %#v, expected no TLSClientConfig.GetClientCertificate", transport.TLSClientConfig)
			//}
			//if !testCase.TLSCert {
			//	return
			//}
			//
			//_, err = transport.TLS.GetClientCertificate(nil)
			//switch {
			//case testCase.TLSErr && err == nil:
			//	t.Error("got nil error from GetClientCertificate, expected non-nil")
			//case !testCase.TLSErr && err != nil:
			//	t.Errorf("got error from GetClientCertificate: %q, expected nil", err)
			//}
		})
	}
}
