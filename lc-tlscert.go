/*
* Copyright 2025 DAIPING.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
* Derived from Golang src/pkg/crypto/tls/generate_cert.go
* Copyright 2009 The Go Authors. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
*/

package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// input is a global buffered reader for user input.
var input *bufio.Reader

func init() {
	input = bufio.NewReader(os.Stdin)
}

// readString prompts the user for a string and returns the trimmed input.
func readString(prompt string) string {
	fmt.Printf("%s: ", prompt)
	line, _ := input.ReadString('\n')
	return strings.TrimSpace(line)
}

// readNumber prompts the user for a numerical value and returns it as int64.
// It handles invalid input by prompting the user again.
func readNumber(prompt string) (num int64) {
	var err error
	for {
		if num, err = strconv.ParseInt(readString(prompt), 0, 64); err != nil {
			fmt.Println("Please enter a valid numerical value.")
			continue
		}
		break
	}
	return
}

// anyKey waits for the user to press any key before proceeding.
func anyKey() {
	fmt.Print("Press any key to continue...")
	input.ReadRune()
}

func main() {
	var err error

	// Define a basic certificate template with key usage and CA properties.
	template := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"DP-SELF"},
			Country:      []string{"CN"},
			Province:     []string{"Shanghai"},
			Locality:     []string{"Shanghai"},
		},
		NotBefore: time.Now(),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		IsCA: true,
	}

	fmt.Println("Specify the Common Name for the certificate. The common name")
	fmt.Println("can be anything, but is usually set to the server's primary")
	fmt.Println("DNS name. Even if you plan to connect via IP address you")
	fmt.Println("should specify the DNS name here.")
	fmt.Println()

	template.Subject.CommonName = readString("Common Name")
	fmt.Println()

	// Prompt the user for DNS names and IP addresses (Subject Alternative Names).
	fmt.Println("The next step is to add any additional DNS names and IP")
	fmt.Println("addresses that clients may use to connect to the server.")
	fmt.Println("You MUST specify the server's IP address if you plan to connect")
	fmt.Println("via IP. When you are finished, just press enter.")
	fmt.Println()

	var cnt = 0
	var val string
	for {
		cnt++
		if val = readString(fmt.Sprintf("DNS or IP address %d (leave blank to finish)", cnt)); val == "" {
			break
		}
		if ip := net.ParseIP(val); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, val)
		}
	}

	fmt.Println()

	// Prompt the user for certificate validity period in days.
	fmt.Println("How long should the certificate be valid for? A year (365")
	fmt.Println("days) is a common choice.")
	fmt.Println()

	template.NotAfter = template.NotBefore.Add(time.Duration(readNumber("Number of days")) * time.Hour * 24)

	// Display the final certificate details to the user for confirmation.
	fmt.Println("\n--- Certificate Details ---")
	fmt.Println("Common Name:", template.Subject.CommonName)
	fmt.Println("Organization:", template.Subject.Organization[0])
	fmt.Println("Validity:", template.NotBefore.Format("2006-01-02"), "to", template.NotAfter.Format("2006-01-02"))
	fmt.Println("DNS SANs:")
	if len(template.DNSNames) == 0 {
		fmt.Println("    None")
	} else {
		for _, e := range template.DNSNames {
			fmt.Println("    ", e)
		}
	}
	fmt.Println("IP SANs:")
	if len(template.IPAddresses) == 0 {
		fmt.Println("    None")
	} else {
		for _, e := range template.IPAddresses {
			fmt.Println("    ", e)
		}
	}
	fmt.Println("---------------------------")
	fmt.Println()

	anyKey()

	// Generate a 4096-bit RSA private key.
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("Failed to generate private key:", err)
		os.Exit(1)
	}

	// Generate a unique serial number for the certificate.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Println("Failed to generate serial number:", err)
		os.Exit(1)
	}

	// Create the self-signed certificate.
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Println("Failed to create certificate:", err)
		os.Exit(1)
	}

	// Save the certificate to a file.
	certOut, err := os.Create("selfsigned.crt")
	if err != nil {
		fmt.Println("Failed to open selfsigned.crt for writing:", err)
		os.Exit(1)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		fmt.Println("Failed to encode certificate:", err)
		os.Exit(1)
	}
	certOut.Close()

	// Save the private key to a file.
	keyOut, err := os.OpenFile("selfsigned.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("Failed to open selfsigned.key for writing:", err)
		os.Exit(1)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		fmt.Println("Failed to encode private key:", err)
		os.Exit(1)
	}
	keyOut.Close()

	fmt.Println("\nSuccessfully generated certificate and private key.")
	fmt.Println("    Certificate: selfsigned.crt")
	fmt.Println("    Private Key: selfsigned.key")
	fmt.Println("-------------------------------------------")
	fmt.Println("Next steps:")
	fmt.Println("1. Copy the 'selfsigned.crt' and 'selfsigned.key' files to your GitLab server.")
	fmt.Println("2. Configure GitLab to use the new certificate and key.")
	fmt.Println("3. Import 'selfsigned.crt' into SonarQube's Java trust store.")
	fmt.Println("4. Restart both GitLab and SonarQube services.")
}
