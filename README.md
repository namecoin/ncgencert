ncgencert
=========

ncgencert generates TLS server certificates of the form expected by Namecoin, suitable for serving HTTPS on `.bit` domains.

It produces a name-constrained 3-level certificate chain (AIA Parent CA → Domain CA → End-Entity), a private key for the server, and a `namecoin.json` file containing the TLSA record to place in the Namecoin name value.

Usage
-----

~~~
ncgencert --host <domain> [flags]
~~~

### Required flags

| Flag | Description |
|------|-------------|
| `--host` | Comma-separated hostnames to generate a certificate for. Only use one unless `--parent-chain` or `--grandparent-chain` is set. |

### Optional flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ecdsa-curve` | `P256` | ECDSA curve for key generation. Valid: `P224`, `P256` (recommended), `P384`, `P521`. |
| `--ed25519` | `false` | Generate an Ed25519 key instead of ECDSA. |
| `--duration` | `8760h` (1 year) | Duration the certificate is valid for. |
| `--start-date` | now | Creation date, formatted as `Jan 1 15:04:05 2011`. |
| `--email` | `false` | Generate a certificate for an email address instead of a DNS name. |
| `--uri` | `false` | Generate a certificate for a URI instead of a DNS name. |
| `--client` | `false` | Generate a TLS client certificate (EKU: ClientAuth). |
| `--code` | `false` | Generate a code signing certificate (EKU: CodeSigning). |
| `--smime` | `false` | Generate an S/MIME certificate (EKU: EmailProtection). |
| `--parent-key` | | Path to existing CA private key to sign the end-entity cert with. |
| `--parent-chain` | | Path to existing CA cert chain to sign the end-entity cert with. |
| `--grandparent-key` | | Path to existing CA private key to sign the domain CA cert with. |
| `--grandparent-chain` | | Path to existing CA cert chain to sign the domain CA cert with. |
| `--sigs` | | Path to Namecoin message signatures to staple (saves blockchain space). |

### Output files

| File | Description |
|------|-------------|
| `cert.pem` | End-entity TLS certificate. |
| `key.pem` | End-entity private key. **Deploy this on your server.** |
| `caCert.pem` | Domain CA certificate (name-constrained to `--host`). |
| `caKey.pem` | Domain CA private key. |
| `caAIAKey.pem` | AIA parent CA private key. |
| `chain.pem` | Full certificate chain (`cert.pem` + `caCert.pem`). **Deploy this on your server.** |
| `caChain.pem` | CA chain only (`caCert.pem`). |
| `namecoin.json` | TLSA record for the Namecoin name value. |
| `caAIAMessage.txt` | Message to sign with Namecoin wallet for Option 2 (see below). |

The TLSA record in `namecoin.json` uses the format `[2, 1, 1, "<sha256-hex>"]`:

- **2** = DANE-TA (Trust Anchor)
- **1** = SPKI (SubjectPublicKeyInfo)
- **1** = SHA-256 hash

Examples
--------

### Basic: generate a certificate for a `.bit` domain

~~~
mkdir my-domain-tls && cd my-domain-tls
ncgencert --host example.bit
~~~

This creates all output files in the current directory. Check `namecoin.json` for the TLSA record.

### With a longer validity period

~~~
ncgencert --host example.bit --duration 26280h
~~~

This generates a certificate valid for 3 years.

### Using Ed25519 instead of ECDSA

~~~
ncgencert --host example.bit --ed25519
~~~

### Deploying the certificate

After running ncgencert, you have two deployment options:

**Option 1** (simple, uses more blockchain space):

1. Place `chain.pem` and `key.pem` in your HTTPS server configuration.
2. Place the contents of `namecoin.json` in the `tls` field for `*.example.bit` in your Namecoin name value.

**Option 2** (conserves blockchain space, uses AIA):

1. Edit `caAIAMessage.txt` — replace `FILL IN NAMECOIN ADDRESS HERE BEFORE SIGNING` with your Namecoin address.
2. Sign the message with your Namecoin wallet:
   ~~~
   namecoin-cli signmessage "<your-namecoin-address>" "$(cat caAIAMessage.txt)"
   ~~~
3. Save the signature to a file (e.g. `sigs.txt`).
4. Re-run ncgencert with the grandparent key and signatures:
   ~~~
   ncgencert --host example.bit --grandparent-key caAIAKey.pem --sigs sigs.txt
   ~~~
5. Place `chain.pem` and `key.pem` in your HTTPS server. No additional blockchain transaction is needed.

### Example: updating a Namecoin name

After generating the certificate:

~~~
cat namecoin.json
# Example output: [2,1,1,"a1b2c3d4..."]
~~~

Use the TLSA value in your `name_update`:

~~~
namecoin-cli name_update "d/example" '{
  "ip": "1.2.3.4",
  "map": {
    "*": {
      "tls": [[2,1,1,"a1b2c3d4..."]]
    }
  }
}'
~~~

Building
--------

Prerequisites:

1. Ensure you have the Go tools installed.

Option A: Using Go build commands without Go modules (works on any platform with Bash; only Go 1.15-1.16.x; will not work on Go 1.17+):

1. Ensure you have the `GOPATH` environment variable set. (For those not
   familiar with Go, setting it to the path to an empty directory will suffice.
   The directory will be filled with build files.)

2. Run `export GO111MODULE=off` to disable Go modules.

3. Run `go get -d -t -u github.com/namecoin/ncgencert/...`. The ncgencert source code will be
   retrieved automatically.

4. If running an old Go version, you may need to enter the `ncgencert` directory and `git checkout` an older tag, e.g. for Go 1.15.x you should checkout the `v1.14` tag.

5. Run `go get -t -u github.com/namecoin/ncgencert/...`.  ncgencert will be built. The binaries will be at `$GOPATH/bin/ncgencert`.

Option B: Using Go build commands with Go modules (works on any platform with Bash; Go 1.15+):

1. `git clone https://github.com/namecoin/ncgencert`.

2. If running an old Go version, you may need to enter the `ncgencert` directory and `git checkout` an older tag, e.g. for Go 1.15.x you should checkout the `v1.14` tag.

3. Run the following in the ncgencert directory to set up Go modules:
   
   ~~~
   go mod init github.com/namecoin/ncgencert
   go mod tidy
   ~~~

4. Run `go install ./...`.  ncgencert will be built. The binaries will be at `$GOPATH/bin/ncgencert`.

Option C: Using Makefile (non-Windows platforms):

1. Run `make`. The source repository will be retrieved via `go get`
   automatically.

Licence
-------

License matches the upstream Go standard library's BSD-style license; see `LICENSE` for details.
