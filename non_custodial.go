package main

import (
    "fmt"
    "log"

    "github.com/tyler-smith/go-bip39"
    "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/btcsuite/btcd/btcutil"
)

func main() {
    // Generate a random 24-word mnemonic (seed phrase)
    entropy, err := bip39.NewEntropy(256)
    if err != nil {
        log.Fatalf("Error generating entropy: %v", err)
    }
    mnemonic, err := bip39.NewMnemonic(entropy)
    if err != nil {
        log.Fatalf("Error generating mnemonic: %v", err)
    }

    // Generate a private key from the mnemonic
    seed := bip39.NewSeed(mnemonic, "")
    privateKey, _ := btcec.PrivKeyFromBytes(seed[:32])

    // Create WIF
    wif, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
    if err != nil {
        log.Fatalf("Error creating WIF: %v", err)
    }

    // Generate public key and address
    pubKey := privateKey.PubKey()
    address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
    if err != nil {
        log.Fatalf("Error generating address: %v", err)
    }

    // Print results
    fmt.Println("Seed Phrase (Mnemonic):")
    fmt.Println(mnemonic)
    fmt.Println("\nPrivate Key (WIF):")
    fmt.Println(wif.String())
    fmt.Println("\nBitcoin Address:")
    fmt.Println(address.EncodeAddress())
}