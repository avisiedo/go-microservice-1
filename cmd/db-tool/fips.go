//go:build fips
// +build fips

package main

import _ "crypto/tls/fipsonly" // https://go.googlesource.com/go/+/dev.boringcrypto/src/crypto/tls/fipsonly/fipsonly.go
