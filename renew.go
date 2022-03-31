package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

// read the ca cert and key from bolt and write them to disk as pem
func exportCA(boltDBPath string) error {
	db, err := bolt.Open(boltDBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		certbytes []byte
		keybytes  []byte
	)

	// read the ca cert and key from bolt db bucket
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scep_certificates"))
		certbytes = b.Get([]byte("ca_certificate"))
		keybytes = b.Get([]byte("ca_key"))

		return nil
	})

	// write the ca cert to file
	b := pem.Block{Type: "CERTIFICATE", Bytes: certbytes}
	f, _ := os.OpenFile("./out.pem", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	pem.Encode(f, &b)

	// write the ca key to file
	keyBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: keybytes}
	keyFile, _ := os.OpenFile("./key.pem", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	pem.Encode(keyFile, &keyBlock)

	return nil
}

// read ca certificate from pem encoded file and shove into scep_certificates bucket, overwriting the current ca certificate
func importCA(boltDBPath string, caCertPath string) error {
	certData, _ := ioutil.ReadFile(caCertPath)
	certPEM, _ := pem.Decode(certData)
	certX509, err := x509.ParseCertificate(certPEM.Bytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("is ca: %v\norg: %v\nbytes: %v\n", certX509.IsCA, certX509.Issuer.Organization, certPEM.Bytes)

	// shove the bytes back into bolt
	db, err := bolt.Open(boltDBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// testBytes := []byte("asdf")
	// certBucket := "scep_certficates"
	crtBytes := certPEM.Bytes

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scep_certificates"))
		err := b.Put([]byte("ca_certificate"), crtBytes)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func showCA(boltDBPath string) error {
	db, err := bolt.Open(boltDBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		certbytes []byte
		keybytes  []byte
		key       *rsa.PrivateKey
		cert      *x509.Certificate
	)

	// read the ca cert and key from bolt db bucket
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("scep_certificates"))
		certbytes = b.Get([]byte("ca_certificate"))
		keybytes = b.Get([]byte("ca_key"))

		return nil
	})
	// parse bytes from bolt
	key, _ = x509.ParsePKCS1PrivateKey(keybytes)
	cert, _ = x509.ParseCertificate(certbytes)

	fmt.Println(key)
	fmt.Println(cert)

	// fmt.Println(certbytes)

	// fmt.Println(key.PublicKey)
	fmt.Printf("is ca: %v\nnot valid before: %v\nnot valid after: %v\n", cert.IsCA, cert.NotBefore, cert.NotAfter)

	return nil
}

func main() {
	var (
		flBoltDBPath string
		flExportCA   bool
		flImportCA   bool
		flShowCA     bool
		flCACert     string
	)

	flag.StringVar(&flBoltDBPath, "boltdb", "", "path to bolt db")
	flag.BoolVar(&flExportCA, "export-ca", false, "export ca key and certificate")
	flag.StringVar(&flCACert, "ca-cert", "", "path to certificate to import into scep_certificates bucket")
	flag.BoolVar(&flImportCA, "import-ca", false, "import ca into scep_certificates bucket")
	flag.BoolVar(&flShowCA, "show-ca", false, "print information about current certificate authority certificate")

	flag.Parse()

	if flBoltDBPath == "" {
		log.Fatal("must provide db path")
	}

	if flExportCA {
		exportCA(flBoltDBPath)
	} else if flImportCA {
		importCA(flBoltDBPath, flCACert)
	} else if flShowCA {
		showCA(flBoltDBPath)
	} else {
		log.Fatal("must provide instructions")
	}
}
