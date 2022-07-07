package keymanager

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"math"
	"time"

	"github.com/cloudslit/cfssl/helpers"
	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/cloudslit/cloudslit/ca/pkg/memorycacher"
	"gorm.io/gorm"
)

// Keeper ...
type Keeper struct {
	DB     *gorm.DB
	cache  *memorycacher.Cache
	logger *logger.Logger
}

var (
	Std *Keeper
)

const RootCert = `-----BEGIN CERTIFICATE-----
MIIFDjCCAvagAwIBAgIUJYidzBLfVSlkMiqvJJs+H0UFY08wDQYJKoZIhvcNAQEN
BQAwHzEdMBsGA1UEChMUQ0kxMjMgUk9PVCBBVVRIT1JJVFkwHhcNMjEwMTEzMTE0
MjAwWhcNNDEwMTA4MTE0MjAwWjAfMR0wGwYDVQQKExRDSTEyMyBST09UIEFVVEhP
UklUWTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAKMTs/+vIjnKDMKM
9E9sutG8q6xnPAAH1bsvlNrXdxH/5FNYCyInWqN3ou9LmDgO0bPwvghI76EhQ8Dw
KtYJ6JsQ51bJksOz4b2XQO9q8fzZQ7IRpX2s/eX4ck1XfCdxhY1bGDruFX0krG6M
tcSQKmw7pDfx4cHaeswgeKAu+/OuPeT9ZBzS2iuvOr6ph1B2V+Nw779aiT6nxmpF
XZcMnjQgKLmrJORq5lnnhwd0/wcrwTp7IYSeuZq3c2UIxtF8cKOWvE8G5xHwvXtx
P8UCI0j0qyITn4JyyTX0TxuU5t/jyylO0O6N1PBmqdf7rY70KE5Ot9O+UCq2shKA
wT9TNZFCv/omXWQClKBkPaR2LKtwnZRwKhA/+5lt5b2zXB8YER5sPgtn1e/e4qNZ
BFMAzCJPiEd2rc/yjqqcF/bA8j4ZRX6tS6/JHSagfkWf0KKmMg0+DwA/S7gHH2fL
5uHlPpCD5VSYbSQhfZzoTN0aDr+4OxGPWCA3cLGSwsfl0diDpDhHLl/GJ6TKcAve
32QzxC6NegR13VF8XXwfYyYeHXyQI4ugSZ4q0jzZ6gdK3cQ6d4kYv7xwMtr0xJ/4
ibNi6ptNUJlm6zfLFzL0nkPAQvcyKa7s3G62UxbJrCR+DtszzQ9DgSjQepkrGyzp
dbDXsUiXh7U+9ZjGad0ShikY1G5pAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSENRvuAYt6s0gnKPkElnGDQ+vpTDAN
BgkqhkiG9w0BAQ0FAAOCAgEAgWgdIv7AuF3NF0GY4v3CI2gmlVfThQxPfh78E03D
+ur66AQ43cBFL4g5U9hkCTNlCtgBHYJNHYT1a9eNWEGEqqcvFiubGJ9wCGd6bOA+
ucjub2KkdE1lfk2anWhba0iseO6MrYcLxpIDo9Vle7Tv00a7QpQrIAw7BeFJUduQ
LKgrqK1DBc+sjZVB3g1Nbf3c1TzXCbmJ1r2wJbTdOBOGhFpn02odzgR9717Q69hf
TQ+SOrnhhxP25wycNndvRtW+G0tPFNmDVHgpYwJWJG8WNBrNi4XhRDJoCQ8OdIc5
LyvfgqmteUGTvMRPkp4t/ar0xGb0Zole41mgRIpukIkUSIj/FPexvpLCYV7vBuPC
DmhjNVfcHtJwhNpd/SCNoQSLIrUPFhC011KaEQObYDN+5yzVDuBzi4VBfGC15t6Z
UTd7KWvboF7Ido45UNuWb7ts2OuTHpIMjRiKQVyDIIIbOg0vKh2qCwtfi81GQYfA
6+vaUaE7CFwnoKrlvZoMGCF38EGmFwRDxi/eUisMql1QRIzisL7o/lpjW0JJI+SL
l8hS4ybA8XXb3RaFz4XNgoLHp3FoPtrxmuoGJxl2ZnBAj8/B9rjj8l7brpUGCQbb
tqo/wKzI9SQPZXCMs+8az26lOA6wbcyjY4rhzEAtPWnQTRlCTjicakmER9c8GVeo
upg=
-----END CERTIFICATE-----
`
const RootKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

const SelfCert = `-----BEGIN CERTIFICATE-----
MIIFSzCCAzOgAwIBAgIUME6urC5H9CuOcaUYD7rch8xyBe8wDQYJKoZIhvcNAQEN
BQAwHzEdMBsGA1UEChMUQ0kxMjMgUk9PVCBBVVRIT1JJVFkwHhcNMjEwMTEzMTE0
MjAwWhcNMjMwMTEzMTE0MjAwWjA7MRkwFwYDVQQKExBTSVRFIENBIElERU5USUZZ
MR4wHAYDVQQLExVzcGlmZmU6Ly9zaXRlL2NsdXN0ZXIwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQDTe40aiw7eMjLbSQ4ozfIcOcnysHYdmiaj5waq4nVE
dA9mBy13Cn6a5ygS+5hiHCUtFms7xv1sLenyV6WHOMk/ykpZO8WSk2cCg9Mz4HGV
zs7GxkrqydOwx/GhM/gHXif2QtrUc7E+6EBIrN5Z162xxn2zqdhQPQdrSY7g9H6b
c9GMMZ8VXZvh2GLpdmq9749MRtpBPaWFARQduGkh1v18yczfeCXfPqXXob8Bt3fe
EoU7v9PnLb3D+6I5Lc9Fe7zYjNrDgODPo8psQWwNs/5UMmzPIw5c66Slc0ybKrtI
RkPYDhkE36Oq8+1jGJqDc7ranen79J6J9OrxNNt+8l4o18nxuTu02VfJ+IdlDsFp
8RjrnDjSd3zfuhnVu5lT+2WRfcXYLsORchfG5s77ZYW0JnGg9eWRYcmVMckOvTkO
CkHBxNnN0X3hM0SFcjy39CItZ+P2d/0/Eh1mxiF/vBnGtaainu4lARMgpf7lPPYy
qg4t3JYFvO2JMBSfBGiJ7g5wAa9WSvMeErETOd3hn4Im6Uay8pcvid41i3NFHj3T
kRNA8eCBXMBPXld9jPTSHfAIyvu/qr1i3bkdIBX4CV3kG0GfDRTL1jxcy6Hy5Ngs
ztRORVJ8odaGYIXokBtE5EMVYimiMU0vYl49xEyEdG8ne4TV0I/DirGk19r5+0pF
8QIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV
HQ4EFgQU059L1fTEbv8U13NhAe+KCswUvbUwHwYDVR0jBBgwFoAUhDUb7gGLerNI
Jyj5BJZxg0Pr6UwwDQYJKoZIhvcNAQENBQADggIBAAjovPqTpAd+rsDi/7snSwsq
bpzOSi1pQ+en+aE4vivmyfmoMPOLNgaCOP1fu1XfGjvZO09QMrqPmTN2mtysFBvr
sxgRmlUIlq08LYOHW/geu7oYrBF+dpxOz476oAV7hKyJqEJLk7hPouoq8KAo0Kyu
baZL/Vv9r9cHP6qrJMgbsEW+98vGdsuFhJXSq/Z+vziVI4FElUlF3UkQa5c4NUFx
nS8wEEi7DnpvoQLqnqGYTSfZgi6ud8lC6M3ksjOO86+mUvmr5I1hWgHh7p9SbWYn
D+7HMCoInSPY48DdSHd1ITwLXHB1vOV3LbhsgdU+9PH1/INO5Fg8KLD7CXhO8i9H
xz3kJl+HHKlMzIJbbpv/O47kdU340rLt4R8nR4PBk1uUGN22tLty9gOhwmYcs1Tp
+hGAVhNE1LNJKu0qclJ3fFghrwhJPvgJVdlRB090vtPVZgBNezrt5Ab6gxyowAcY
dQwm1oU+YTsQlCwgh2nMzWsizpLDagXy6rTDgLbxq6RK3pxtwOmJlL6WWv9s6sPb
op0XJ7KmFeDRzoi2f9w+wQ/ftZF8IXOwD4n3C1n9sFSxtC0QcwNHZ3iCED1N1Epv
8Jd5HdpnPdcmdGOiW1dhTSkR6w6X2jJpkNb3rifgwVkO2/btIYfyAqv/UggamY31
irIAFQZ5jp9jdUvFvfKB
-----END CERTIFICATE-----
`
const SelfKey = `-----BEGIN RSA PRIVATE KEY-----
MIIJKwIBAAKCAgEA03uNGosO3jIy20kOKM3yHDnJ8rB2HZomo+cGquJ1RHQPZgct
dwp+mucoEvuYYhwlLRZrO8b9bC3p8lelhzjJP8pKWTvFkpNnAoPTM+Bxlc7OxsZK
6snTsMfxoTP4B14n9kLa1HOxPuhASKzeWdetscZ9s6nYUD0Ha0mO4PR+m3PRjDGf
FV2b4dhi6XZqve+PTEbaQT2lhQEUHbhpIdb9fMnM33gl3z6l16G/Abd33hKFO7/T
5y29w/uiOS3PRXu82Izaw4Dgz6PKbEFsDbP+VDJszyMOXOukpXNMmyq7SEZD2A4Z
BN+jqvPtYxiag3O62p3p+/SeifTq8TTbfvJeKNfJ8bk7tNlXyfiHZQ7BafEY65w4
0nd837oZ1buZU/tlkX3F2C7DkXIXxubO+2WFtCZxoPXlkWHJlTHJDr05DgpBwcTZ
zdF94TNEhXI8t/QiLWfj9nf9PxIdZsYhf7wZxrWmop7uJQETIKX+5Tz2MqoOLdyW
BbztiTAUnwRoie4OcAGvVkrzHhKxEznd4Z+CJulGsvKXL4neNYtzRR4905ETQPHg
gVzAT15XfYz00h3wCMr7v6q9Yt25HSAV+Ald5BtBnw0Uy9Y8XMuh8uTYLM7UTkVS
fKHWhmCF6JAbRORDFWIpojFNL2JePcRMhHRvJ3uE1dCPw4qxpNfa+ftKRfECAwEA
AQKCAgEAv87fO6PD6GM/HQQ1g2zGmjMPpr3FYzPct+thcLvUADIDhVYdNkqeaYAe
KJlorBM65JngyGbCkstR1CsKRoqKfQDMTgKYP0jRtMY7WGHIo0be5AoVCL3k8gFm
df7chlIvjHs2XlpP9+5N35xqTrH/J64PdFQnjR7NC9G7dOxLqFJzS+P1lTtBlqTb
kUOFrJ8KKPRaH7H53ZgM1qfbMoX7gqLE3inqr3/yipB43OW6IgBKNtbVKmeiJY6k
o5eY1jxUG1QlVgwd2jWD8oujpQoLMfJKHdR9wmk2+5iHcnRfSD3yapLevjzYBMUY
GTjHa1Ibtwxim9JRuG4aaIq4SXspSdo4YNHWC4ud5nYw1IVQc7fxKVwqMXolr2bM
IjFw2QvuHBgcjlaa+XfYPyGpuKGSDlWjkn6OsKKS2P3t7ctY2i/wNJClpGyLvU6a
LLRs62ZKRpK9Xx9EcgMniX2b86WZQpSrZ/eaj1gtfilkrfWHieoOhLzy+B3yzaiK
AtyBPTejgvXKXwzPmZtsmxes+iIzGSd22Fkg19QZNS2vAdNz74/brBLCNtS6yU+/
pNv3d4CGDZ4aIqbf+oO90Hw3c8GKR9iXGRCTW40DopcbNJPGdVi62/2PEaEtoNG6
gNcN+HNwNgCLK98rpkJFE0Rv+Na1zjWaezKwj5YHhZfQnO6s5cECggEBAP4W8SCW
2I1BpOukcZdkPjzYTjzMxKQjlzmrlWHPu6BF1qQ6E1uC3rDIbQ0oMNDzqhuhyrZS
jG8XDbLFT7oFj7jMDTaWOoE0LkSvnpVRpyuSpSd1kZeRRuxYQ8EigxEyJYQ+Y5gy
XfBJHSlOk7QWawB2WcAReCU2hGtaDbrxReK/QkMGBYPf3zbPa/y7kqhTayejLwP3
kjhG0hkn0wBxG1Vbj0ryvWt3l3D52IYttgPK8qAjYx2t7OrxGOOn4714D8I/FqLr
vAN0aKbve0QfC8cSlJU7Z3D9QqgHhtZtirl6Z21ws4ogw1wwI9ahCDqCQYd/3pGG
5S1EggqacGxT0KsCggEBANUSmgOzRiKzLmeKWxSGnZNuKJUExLGNVoFqJZ/JaYLa
pc41HxVdtX/Wl0ZsCSu2bIzWJD6+BcT+wzzU1kOaa75ciKgPqh0ErYMwG2DWzowa
MTs++d/JRIkqct0kiZTd8nJUNTpnrXz9vNMoiNNF6By0gA+ZNEgTuKSnFfqvLb+B
XklAI2dvJMN1602jOx+BuIqCyLRJRDJwjWIDky66CDijG4Eayocn7QLlYIyMLUuQ
UXsXesxO67MA/Jzem/Dx/HDbdDDpZ+N18+Pub/u0nFA/vbefY9DcSbf2aV4aQyde
uXXsWrpd8dshrHVohEKLgep48AXblcZGxaN3TNoZ29MCggEBAP3IiFqWkAC2mjTK
cLJXa1p2ad6MX7PZ4Ie3e2LQi4SPfM1XPFJgqnBEH7fOdsOdPECRHtlwJdgnXIU8
Ul9ogp5/IItvDUxThAsSpgBaJ/B7bf21jg+nCQGzPyk+gU7BmXs2nV88n1sKi9fg
JeLvqTwy+X2/dRMmGqjmr2QS7EyH33T2JLgM+PPTxPYPm3IIr7RNZv49XoxbICoD
/tooHrbo1nxzawJV1qr+wWdzbKLpJ+EOt0bDmykmWke4Pt3Vd1f6j292qLmCoaxq
8eGeaLMTOdi8FptiOht+OQ0fKDoNqhRDRvAlvTrs0j7jZEacJzthWjpcU7cdAA37
J2LrbDUCggEBALnZu9U2bhqeR/+wQrmooZGTKHqy6g4kxiujtqWlPQ8SQEWZOD+e
uU2Ek/atDWK/f/doYb0Iamfl/83zp+DXtNsaQ2i0ISGmjuI69+aD9y1lO0P+Ll3w
ZINwLziNQiRDY9IteTA0drLrb+SPGqmN9GP1XS4958hmy0tzIkzCuBiucttZwofZ
/isvk8rocg2NTLYkVYRL09xbKDcx/xNm2Pzt6HO4NqelP6qjAJAXRPsAKtI/LLFA
tX4xgiYiVcrYh+S4xqRTMnbIz68krzDR3PZyYrzjnmDzhKmdmVKnfaO9j+839ftR
LkCBIrhWLecNYIhwbIvveLi1ynZG/RXQMFkCggEBAJXB0s2r49zeYs6/5uWIW+BW
yG3Msyly/lRU7KIQQqX3CVD/ba7Pc8SoMFOaPPWFJts9PmhPWBQQ6k8XUo41/NOo
pK2vcPeHN/9YD3B0rKDeCXopPae/Zdy3wNC4AOrKaC2xAsCjmvBWoW4MOc6w36Qb
iKyT9BrrxMb/TqXGM8OD+w1iNdjRwzLa/Eh/5ML9ntM399mkDhxGOO0ADmtZ/ScN
yFIyk39Ca4GbY/24rFtO9UE1FId6B4jJAUQ/EwY1H8M4zedp3ywAOdlxH0tLKtBk
KHpa0CTtQFwoEaDMOFVUvfBzg0hY4Lm7xWc5uXlEKu7TmQ0qECtw/N2iuUlA6Cg=
-----END RSA PRIVATE KEY-----
`

// ...
const (
	cacheKey  = "key"
	cacheCert = "cert"
)

// InitKeeper ...
func InitKeeper() error {
	db := core.Is.Db
	Std = &Keeper{
		DB:     db,
		logger: logger.Named("keeper"),
		cache:  memorycacher.New(time.Hour, memorycacher.NoExpiration, math.MaxInt64),
	}
	return nil
}

// GetKeeper ...
func GetKeeper() *Keeper {
	defer func() {
		if err := recover(); err != nil {
			logger.Named("keeper").Fatal("Uninitialized")
		}
	}()
	return Std
}

// GetCachedTLSKeyPair ...
func (k *Keeper) GetCachedTLSKeyPair() (*tls.Certificate, error) {
	keyPEM, certPEM, err := k.GetCachedSelfKeyPairPEM()
	if err != nil {
		k.logger.Errorf("tls.Cert Get error： %v", err)
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		k.logger.Errorf("tls.X509 error: %v", err)
		return nil, err
	}
	return &cert, nil
}

// GetCachedSelfKeyPair ...
func (k *Keeper) GetCachedSelfKeyPair() (key crypto.Signer, cert *x509.Certificate, err error) {
	if cachedKey, ok := k.cache.Get(cacheKey); ok {
		if v, ok := cachedKey.(crypto.Signer); ok {
			key = v
		}
	}
	if cachedCert, ok := k.cache.Get(cacheCert); ok {
		if v, ok := cachedCert.(*x509.Certificate); ok {
			cert = v
		}
	}

	if key != nil && cert != nil {
		return
	}

	keyPEM, certPEM, err := k.GetCachedSelfKeyPairPEM()
	if err != nil {
		k.logger.Errorf("Error getting cache keypair PEM: %v", err)
		return
	}
	priv, err := helpers.ParsePrivateKeyPEM(keyPEM)
	if err != nil {
		k.logger.With("key", string(keyPEM)).Errorf("Certificate key parsing error: %v", err)
		return
	}
	key = priv

	cert, err = helpers.ParseCertificatePEM(certPEM)
	if err != nil {
		k.logger.With("cert", string(certPEM)).Errorf("Certificate PEM parsing error: %v", err)
		return
	}

	k.cache.SetDefault(cacheKey, key)
	k.cache.SetDefault(cacheCert, cert)
	return
}

// GetCachedSelfKeyPairPEM ...
func (k *Keeper) GetCachedSelfKeyPairPEM() (key, cert []byte, err error) {
	key = []byte(SelfKey)
	cert = []byte(SelfCert)
	return key, cert, nil
}

// GetL3CachedTrustCerts Memory > multi level cache > remote process > certificate
func (k *Keeper) GetL3CachedTrustCerts() (certs []*x509.Certificate, err error) {
	cert, err := helpers.ParseCertificatePEM([]byte(RootCert))
	return []*x509.Certificate{cert}, nil
}
