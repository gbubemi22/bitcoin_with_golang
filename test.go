package test

import (
   // "crypto/rand"
    "fmt"
    "log"

   "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/btcec/v2"
   "github.com/btcsuite/btcd/btcutil"
)

func main() {
    // Choose the network: Testnet3 or Mainnet
    network := &chaincfg.TestNet3Params // Change to chaincfg.MainNetParams for mainnet

    // Generate a new private key using the secp256k1 elliptic curve
    privateKey, err := btcec.NewPrivateKey()
    if err != nil {
        log.Fatalf("Error generating private key: %v", err)
    }

    // Create the WIF (Wallet Import Format) structure from the private key
    wif, err := btcutil.NewWIF(privateKey, network, true)
    if err != nil {
        log.Fatalf("Error creating WIF: %v", err)
    }

    // Derive the public key
    pubKey := privateKey.PubKey()

    // Generate a new P2PKH address
    address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), network)
    if err != nil {
        log.Fatalf("Error generating address: %v", err)
    }

    fmt.Println("Private Key (WIF):", wif.String())
    fmt.Println("Address:", address.EncodeAddress())
}
