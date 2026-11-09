// Copyright 2015-2022 Jeremy Rand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func writeJSONTLSA(priv any) {
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey(priv))
	if err != nil {
		log.Fatalf("Failed to marshal CA public key: %v", err)
	}

	pubHash := sha256.Sum256(pubBytes)

	// See the IANA DANE Parameters registry.
	tlsa := make([]any, 4)
	tlsa[0] = 2 // DANE-TA
	tlsa[1] = 1 // SPKI
	tlsa[2] = 1 // SHA-256
	tlsa[3] = pubHash[:]

	tlsaBytes, err := json.Marshal(tlsa)
	if err != nil {
		log.Fatalf("Failed to marshal Namecoin record: %v", err)
	}

	err = ioutil.WriteFile("namecoin.json", tlsaBytes, 0600)
	if err != nil {
		log.Fatalf("Failed to write data to namecoin.json: %v", err)
	}

	log.Print("wrote namecoin.json")

	// Same record in DNS zone file format, for DNS domains.
	record := fmt.Sprintf("*.%s. IN TLSA 2 1 1 %x\n", *host, pubHash)

	err = ioutil.WriteFile("dns.zone", []byte(record), 0600)
	if err != nil {
		log.Fatalf("Failed to write data to dns.zone: %v", err)
	}

	log.Print("wrote dns.zone")
}
