package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checkSum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checkSum...)
	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(publicSHA256[:])
	if err != nil {
		panic(err)
	}
	publicRIPEMD160 := ripemd160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	 firstHash := sha256.Sum256(payload)
	 secondHash := sha256.Sum256(firstHash[:])
	 return secondHash[:32]
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}
