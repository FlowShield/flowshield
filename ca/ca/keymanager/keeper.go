package keymanager

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"math"
	"time"

	"github.com/flowshield/cfssl/helpers"
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/flowshield/flowshield/ca/pkg/memorycacher"
)

// Keeper ...
type Keeper struct {
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
MIIFSzCCAzOgAwIBAgIUQNRKnnBYb52joiWyCa1fXNrtJlUwDQYJKoZIhvcNAQEN
BQAwHzEdMBsGA1UEChMUQ0kxMjMgUk9PVCBBVVRIT1JJVFkwHhcNMjIwNzA3MDU0
OTAwWhcNMjQwNzA2MDU0OTAwWjA7MRkwFwYDVQQKExBTSVRFIENBIElERU5USUZZ
MR4wHAYDVQQLExVzcGlmZmU6Ly9zaXRlL2NsdXN0ZXIwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQDl3f1WOWWUOx3lQg8ovTzATm/UfFEjUZr6e7mV2BsD
YGZH7Z+rrx7yWJg1DgVb9tY1NjYaNkQI8aX/Iom7qZyb3rxPHS1pywD5aICwNSng
QGDDU+2cK4tx7GWSvrafXn4aHMnjLP6KaIG+GUwyF5UFHkRsWrjOR662jyCrlreX
lbpLwoiFIQiBQQzxspoFzqG3RJhIFYo8q3pC27z+U9uMcfIJ0qcK6/9SgzB0Gr6N
AgQnEpRGXTahZ0Knk342SNE2BcLLnQDfmJwqoCW/keUehgoxyPmYfOw8YTG33Ydh
jH8cV+oJMPjAr46V1qnAKile8QvyYijuaVa0MzeJjGjXtBFyaWtVnkVZq426lhZ5
YIq2zLXZgk28b27QG4YbIOQYTSV4t0fV8XBTZTvT3ynyMy481qx3lUIqHgom5TPB
phMAPwkQX/dDF3ray/RyY3Ax0c4JLd5WQBGH16gLsuJb290savB+duZFK4eTQKA4
pYtQ3poIBDtrcGPUyCArMwCwnh0YYm0V3unWYd1mAmAp/6V7kaWQziPmuJncqZex
xsqgGJnkDa5iAdR+tNCoeBpR2Akmm+OuG7t3nVu6MytMO3G9tS0rdO6uiTwqp3FO
EUiOnU71kWw523vXFXR6hYZOTWjbVUCHj8d260tmG8SwfpCdiJAdrVMPIRzkd8y6
mwIDAQABo2MwYTAOBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zAdBgNV
HQ4EFgQU8+KCWZBvl9Ttm48jOIPp2LU81m8wHwYDVR0jBBgwFoAUhDUb7gGLerNI
Jyj5BJZxg0Pr6UwwDQYJKoZIhvcNAQENBQADggIBAEYZ11npQmY9wNjjajhQT9fj
wlfNYf4BOhCo0IAUZlozZxLxujN0O9MNsWJfv2kIy2dAjNt3IsALldI30kvUYXAt
GebM95T8bU/+cFIKiyjkXxW72wrE2I4Y4kIKWhxW/EwtA6Y75OPm+8t3ao0rLcNd
jUSlLd/VIdlUiZ5OybQgW98CQwihy1KvyRn+XelfQWM3Fi9dlYJbR3lHobxE0HlM
hd6PE5IuUPK1HuOv9zVJD2XhuSEhoOMC7c4W5zSnkvKoOzrNjj0cIoAs8AEcrpiZ
+mtW2I0y243nN+rZohHxnQgBzkUlJ6PgsZNvSoddd5slAkY76MdTPfWBysFCHaji
P/qziQQD+uhGxe/ct3yb4eTHeIO6+kK4PP3ssl7NyW/InD/8ppNPB4PEeMWm1ofh
FLRNV6ksyhszzbaFoFoWe0MEcuOJrn5xKIjwsWQUxN4/Iw4R2yARFgwPOlolo9a7
jSSQNm5wCYvoIIiNQMOkxvLthqNTjcVxxDgniJiLl3lCKQmhz+7Xl9zI9HiPqqGd
ZTEks8/nnB/BLxBav+tsBz9GF730QSFUEAvM6DXsazUFG4YF4sgHUywvZzs5VfNt
2mpswhOfAR+D8Eg1RXJS2BycpSjFleKUvmYwyibYaufM3YBbgKCHclZfOLgEurgv
MBJGugOVtaF5ALu/asNr
-----END CERTIFICATE-----
`
const SelfKey = `-----BEGIN RSA PRIVATE KEY-----
MIIJKgIBAAKCAgEA5d39VjlllDsd5UIPKL08wE5v1HxRI1Ga+nu5ldgbA2BmR+2f
q68e8liYNQ4FW/bWNTY2GjZECPGl/yKJu6mcm968Tx0tacsA+WiAsDUp4EBgw1Pt
nCuLcexlkr62n15+GhzJ4yz+imiBvhlMMheVBR5EbFq4zkeuto8gq5a3l5W6S8KI
hSEIgUEM8bKaBc6ht0SYSBWKPKt6Qtu8/lPbjHHyCdKnCuv/UoMwdBq+jQIEJxKU
Rl02oWdCp5N+NkjRNgXCy50A35icKqAlv5HlHoYKMcj5mHzsPGExt92HYYx/HFfq
CTD4wK+OldapwCopXvEL8mIo7mlWtDM3iYxo17QRcmlrVZ5FWauNupYWeWCKtsy1
2YJNvG9u0BuGGyDkGE0leLdH1fFwU2U7098p8jMuPNasd5VCKh4KJuUzwaYTAD8J
EF/3Qxd62sv0cmNwMdHOCS3eVkARh9eoC7LiW9vdLGrwfnbmRSuHk0CgOKWLUN6a
CAQ7a3Bj1MggKzMAsJ4dGGJtFd7p1mHdZgJgKf+le5GlkM4j5riZ3KmXscbKoBiZ
5A2uYgHUfrTQqHgaUdgJJpvjrhu7d51bujMrTDtxvbUtK3Turok8KqdxThFIjp1O
9ZFsOdt71xV0eoWGTk1o21VAh4/HdutLZhvEsH6QnYiQHa1TDyEc5HfMupsCAwEA
AQKCAgEAp9+FgpEMZEMAREwIMiJx2afu9+mKgBa67i+pi4t1vvOJ/wHFWGbeXPLi
xexlcZJFQjtDK8Vxzm5cYoXgvNLT8ump8MVIQzjkj0EIqqdI2+NoR35ly2Xpwtt4
GsF5Mc6anYtkzaowgvhruF9VTEf4pvJB2jgvO0cSi3Tf0TCTB/trQKTjewZy5aKS
R3m+DnunkuZfqFVCzagV8/nyfnLTpjAZKZ9UKx/LKYFxw3k8rFJKohZpfzHYEewN
u8B1rkNjtuZiKr2Qw7r6Yg7vQobUI4SUsWMrFXg2Nqps2vDITC+FLTfvTaMcJ+yO
X1OUSSkBabr0lKGnbuYzUgsNhkCq12rvmQ32Jre2aEkDE4FOA/5cAnLGBGb0JG4e
ex8i3RYs61GqH1Seu1dNlo+AegrtHeoQ+3P6dQemp9DgjH4SKQY4cQo6uFtlmjeR
MPe5DtCIrNq/ZzzxoZS/fUKrCU1jL0wBS2hDXr3M1x0TS9Ryg1kHmjlSf5gozDoX
tpprgROiBF+1jAfIWbpsKQfZkBD/4BW6L2iVTCtZimNWxbMc9+UMpDGCIXKIjPdc
lUzoY2Pj3pbVTqmT6Oi1Gski1AIiiPpwPk+1SwsMq3Zhuh2uzcfazfMWbkQu0nqr
BUL7aiEfXFaCm7BqN2jm0mnulwUBqDTKL83YR5C7HOPUM4NQrSECggEBAPD5NKFd
5qGiyAVjNVrd4b6CjUSxxcHvCeqv1RBOykImSlgHj6lTWzGYBlvZwUuTE5CML3w5
J96ubYPj77b2hc9cKEBHF1VZhiuPBXzwa7Bk4SYt2v3AkKPau8naX2qoYXycI5JM
ShzRQt+rm20SJEV+pl5USjkFYqlHRBW2eDsWJksqgY2rnZ55UI4vxr7vCRwk4MA2
oD8ODV8t/wFspuZh7IlFJzA/lFXOfKgSA+OrahA6atsFK4nwPX/xuIDVTpqHGX6j
Jb4FUSGjDV31RI2RFfncGr502lIN3SClr7eSW1IUwoW8F6+ZHnnfr5KgQARs8c1/
pE2BFehRsuo0SFMCggEBAPQzfNsWYjcBPbsddMQWIQSWU58ZNX8iVcf8MgJb0zOw
itwx4IvMlAxxtqAExaeQ3+bmL8XBupG0a5JZUgyaXs565jrV8Re4gUb1nlclSBWP
OnG62WM8EG/ILMzxf7FpX+fYpD1HUvjrbFcfFuzXhLuPkwljW45MQDJEmhloNwAo
dPvvaddf4/W9c2qjFeqOrf1jgG8kRIT/zZPStI7avgP5nNvxt6FgpIe4BS/VevRk
n2MlNs0VI6BoDc4NTq5KFBv+nApcZNAQEL3H5vq5pYA0bcq3g6iAnhkb/eHZLGen
/VkmiRiSDkbzLoQM4BNyGJ3mEmeHGdoDYw6eX/jIW5kCggEAWsHRI7GFCn2Pjg/m
aXnF6wWqhSPbUoZgGsbsnT6/iJh5SZxXbDOb1hrm0jM6TOdw5/EUdVnlfUX+szzm
7Ob+ULHp6wObcybLlJ5CN+Mo+/+SmNmOcCHVmBDqx0R6yWXviYoZD6GyDBZ4dFti
p3q8tHvV1xMx/TXdCdpwdykJMV/PPmIc9ymarLQONe8ikIjgynvNNjectQLq4F0n
fPbaCUz6VFz6PH7FtGeXpYlbc6T8xm3qDuFsm4Ai+YwlrEgqWaLmZD64GVBRVTTe
9PNKRXNObpOKcw75pwvXq9MbUi1KPajZ9pp35UfrJYWsz7GRStlpXLdlP1eN07jp
hLH6RwKCAQEA1eZnngwsOVv383douP1dLIR0eK5/In4zvzmToGdIR0WDTD7QHgQz
Rfcw7VdgvlbzGHBWNhVsU4ZCl21vpiRtmNUj5zNQu+NcMYihilnYmzHbEpWFJxwM
la00OMvserz/SbiEnDxmXCzyuBk+XnSlChlHxPhn0OvPa3iVtl2Hl7bYSXk7L3EJ
301z1FtHri4ODx5h+Hg/IHRkYPA2Qc8uk5LIKAvBOjTJySuSN4T57ypYRmLpbpfu
nemm9e7IFXhDxwWtLpIhp/H6iBGaq9GDOxoxdVhrlWQbl4jiSDqPX7hQ/Q75FTGS
GemWvAn/GGlfUKefRVmcdk9zK/HjxKnTOQKCAQEA4R1aS6yGu+lZ/vQ18ZsTpvEz
KEDVus4szf0WEpxmrRvvF1174sqtYZHPxm2IcZPE8aeD3kRWYRztcwqVEYnLEbJX
hTRvzzm8frBuSnkydODAeo73hLd690cyYqRC7O1sneIVz0cfR+HCflC8ymer47dt
pzgwLkhjKlwU/aRpkyxoxwKEJky2GCQravG7jm1GO4xAo0c8AJ9rFS05VLINSeAF
stVtk9/FD6doY+36Fhv7xGw/5MFONtTJJ4vgoSqyzQ8NZSq2x15U4jPdggc2ymiF
+VHeiT1HO0ltFXwE9kRho8HCz1QzFdYK8jaYBrmnSSEQsSwLn17cQQyvFfCu1Q==
-----END RSA PRIVATE KEY-----
`

// ...
const (
	cacheKey  = "key"
	cacheCert = "cert"
)

// InitKeeper ...
func InitKeeper() error {
	Std = &Keeper{
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
	//_, cert, err := GetKeeper().GetCachedSelfKeyPair()
	//if err != nil {
	//	logger.Errorf("Error getting priv key: %v", err)
	//}
	cert, err := helpers.ParseCertificatePEM([]byte(RootCert))
	return []*x509.Certificate{cert}, nil
}
