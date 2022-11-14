<p align="center">
<img src="https://user-images.githubusercontent.com/52234994/165200623-c60e956b-5805-4088-bf58-f97ebd8ae8b4.png" 
    width="40%" border="0" alt="CA">
</p>

# CA

CA is a PKI developed based on cloudflare cfssl,Public key infrastructure (PKI) governs the issuance of digital certificates to protect sensitive data, provide unique digital identities for users, devices and applications and secure end-to-end communications.

CA includes the following components：

1. TLS service, as the CA center, is used for certificate issuance, revocation, signature and other operations.
2. API services, as some API services for certificate management.
2. OCSP service is a service that queries the online status of certificates and has been signed by OCSP.
2. SDK component, which is used for other services to access the CA SDK as a toolkit for certificate issuance and automatic rotation.

## How the Certificate Creation Process Works

The certificate creation process relies heavily on asymmetric encryption and works as follows: 

- A private key is created and the corresponding public key gets computed 
- The CA requests any identifying attributes of the private key owner and vets that information 
- The public key and identifying attributes get encoded into a Certificate Signing Request (CSR) 
- The CSR is signed by the key owner to prove possession of that private key 
- The issuing CA validates the request and signs the certificate with the CA’s own private key 

<img src="https://user-images.githubusercontent.com/52234994/165200483-0b7e9698-552c-4a9f-b9b0-afce84e8c313.png" alt="image-20220425164057905" style="width:50%;" />

Anyone can use the public portion of a certificate to verify that it was actually issued by the CA by confirming who owns the private key used to sign the certificate. And, assuming they deem that CA trustworthy, they can verify that anything they send to the certificate holder will actually go to the intended recipient and that anything signed using that certificate holder’s private key was indeed signed by that person/device. 

One important part of this process to note is that the CA itself has its own private key and corresponding public key, which creates the need for CA hierarchies. 

## How CA Hierarchies and Root CAs Create Layers of Trust

Since each CA has a certificate of its own, layers of trust get created through CA hierarchies — in which CAs issue certificates for other CAs. However, this process is not circular, as there is ultimately a root certificate. Normally, certificates have an issuer and a subject as two separate parties, but these are the same parties for root CAs, meaning that root certificates are self-signed. As a result, people must inherently trust the root certificate authority to trust any certificates that trace back to it. 

<img src="https://user-images.githubusercontent.com/52234994/165200520-842ecf88-bfea-441b-a1af-53260ce4085f.png" alt="image-20220425164028072" style="width:50%;" />

## CA overall architecture and working mode

![image-20220425165623191](https://user-images.githubusercontent.com/52234994/165200574-ac647d20-1044-4580-8378-862d4fd4af9e.png)

## Building

Building cfssl requires a [working Go 1.12+ installation](http://golang.org/doc/install).

```
$ git clone git@github.com:CloudSlit/cloudslit.git
$ cd ca
$ make
```

You can set GOOS and GOARCH environment variables to allow Go to cross-compile alternative platforms.

The resulting binaries will be in the bin folder:

```
$ tree bin
bin
├── ca
```

## Configuration reference

When CA starts each service, it needs to rely on some configurations, and the dependent configuration information has two configuration methods:

**configuration file:**

The configuration file is in the project root directory：`conf.yml` ,The file format is standard yaml format, which can be used as a reference。

**environment variable:**

In the project root directory：`.env.example`, The file describes how to configure some settings through environment variables.

**Priority:**

The configuration priority of environment variables is higher than the configuration in the configuration file.


## Service Installation

### TLS service

TLS service is used to issue certificates through control`IS_KEYMANAGER_SELF_SIGN` Environment variable to control whether to start as Root CA.

- Started as root CA, TLS service will self sign certificate.
- When starting as an intermediate CA, the TLS service needs to request the root CA signing certificate as its own CA certificate.

Start command：`CA tls`，Default listening port 8081

### OCSP service

OCSP online certificate status is used to query the certificate status information. OCSP returns the certificate online status information to quickly check whether the certificate has expired, whether it has been revoked and so on.

Start command：`CA ocsp`，Default listening port 8082

### API service

Provide CA center API service, which can be accessed after the service is started`http://localhost:8080/swagger/index.html`，View API documentation.

Start command：`CA api`，Default listening port 8080



### SDK Installation

```
$ go get github.com/cloudslit/cloudslit/casdk
```

The classic usage of the CA SDK is that the client and the server use the certificate issued by the CA center for encrypted communication. The following is the usage of the sdk between the client and the server.

See：[Demo](https://github.com/cloudslit/cloudslit/casdk/tree/master/caclient/examples)

