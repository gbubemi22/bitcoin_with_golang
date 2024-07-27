package main

import (
    "fmt"
    "log"
    "encoding/hex"
    "strconv"

    "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/txscript"
    "github.com/btcsuite/btcd/wire"
    "github.com/btcsuite/btcd/btcutil"
    "github.com/btcsuite/btcd/chaincfg/chainhash"
    "github.com/btcsuite/btcd/btcec/v2"
    "github.com/imroc/req/v3"
)

type UTXO struct {
    TxID         string  `json:"txid"`
    OutputNo     int     `json:"output_no"`
    ScriptHex    string  `json:"script_hex"`
    Value        string  `json:"value"`
    Confirmations int    `json:"confirmations"`
    Time         int     `json:"time"`
}

type UTXOResponse struct {
    Status string `json:"status"`
    Data   struct {
        Network string  `json:"network"`
        Address string  `json:"address"`
        Txs     []UTXO  `json:"txs"`
    } `json:"data"`
}

type FeeResponse struct {
    HourFee float64 `json:"hourFee"`
}

func sendBitcoin(receiverAddress string, amountToSend float64) error {
    sochainNetwork := "BTCTEST"
    privateKeyHex := "4b48e0d11b191b5e685359b09747c9c24f79834fca273edf9e5e9767011a36b8"
    sourceAddress := "1GwB4FCSbSxwvsaNtykfxp76DGqaGfWGUP"
    satoshiToSend := int64(amountToSend * 100000000)

    // Fetch UTXOs
    utxoResp := &UTXOResponse{}
    _, err := req.Get(fmt.Sprintf("https://sochain.com/api/v2/get_tx_unspent/%s/%s", sochainNetwork, sourceAddress)).
        SetResult(utxoResp).
        Send()
    if err != nil {
        return fmt.Errorf("error fetching UTXOs: %v", err)
    }

    // Fetch recommended fee
    feeResp := &FeeResponse{}
    _, err = req.Get("https://bitcoinfees.earn.com/api/v1/fees/recommended").
        SetResult(feeResp).
        Send()
    if err != nil {
        return fmt.Errorf("error fetching fee: %v", err)
    }

    // Create transaction
    tx := wire.NewMsgTx(wire.TxVersion)
    var totalAmountAvailable int64
    var inputCount int
    for _, utxo := range utxoResp.Data.Txs {
        txHash, err := chainhash.NewHashFromStr(utxo.TxID)
        if err != nil {
            return fmt.Errorf("error parsing txid: %v", err)
        }
        outPoint := wire.NewOutPoint(txHash, uint32(utxo.OutputNo))
        txIn := wire.NewTxIn(outPoint, nil, nil)
        tx.AddTxIn(txIn)

        amount, err := strconv.ParseFloat(utxo.Value, 64)
        if err != nil {
            return fmt.Errorf("error parsing amount: %v", err)
        }
        totalAmountAvailable += int64(amount * 100000000)
        inputCount++
    }

    outputCount := 2
    transactionSize := inputCount*180 + outputCount*34 + 10 - inputCount
    fee := int64(float64(transactionSize) * feeResp.HourFee / 3)

    if totalAmountAvailable-satoshiToSend-fee < 0 {
        return fmt.Errorf("balance is too low for this transaction")
    }

    // Add output
    receiverAddr, err := btcutil.DecodeAddress(receiverAddress, &chaincfg.TestNet3Params)
    if err != nil {
        return fmt.Errorf("error decoding receiver address: %v", err)
    }
    pkScript, err := txscript.PayToAddrScript(receiverAddr)
    if err != nil {
        return fmt.Errorf("error creating pkScript: %v", err)
    }
    tx.AddTxOut(wire.NewTxOut(satoshiToSend, pkScript))

    // Add change output
    changeAddr, err := btcutil.DecodeAddress(sourceAddress, &chaincfg.TestNet3Params)
    if err != nil {
        return fmt.Errorf("error decoding change address: %v", err)
    }
    changePkScript, err := txscript.PayToAddrScript(changeAddr)
    if err != nil {
        return fmt.Errorf("error creating change pkScript: %v", err)
    }
    changeAmount := totalAmountAvailable - satoshiToSend - fee
    tx.AddTxOut(wire.NewTxOut(changeAmount, changePkScript))

    // Sign transaction
    privKeyBytes, err := hex.DecodeString(privateKeyHex)
    if err != nil {
        return fmt.Errorf("error decoding private key: %v", err)
    }
    privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

    for i, txIn := range tx.TxIn {
        sigScript, err := txscript.SignatureScript(tx, i, changePkScript, txscript.SigHashAll, privKey, true)
        if err != nil {
            return fmt.Errorf("error signing transaction: %v", err)
        }
        txIn.SignatureScript = sigScript
    }

    // Serialize and send transaction
    serializedTx, err := tx.SerializeToHex()
    if err != nil {
        return fmt.Errorf("error serializing transaction: %v", err)
    }

    resp, err := req.Post(fmt.Sprintf("https://sochain.com/api/v2/send_tx/%s", sochainNetwork)).
        SetBody(map[string]string{"tx_hex": serializedTx}).
        Send()
    if err != nil {
        return fmt.Errorf("error sending transaction: %v", err)
    }

    fmt.Printf("Transaction sent. Response: %s\n", resp.String())
    return nil
}

func main() {
    err := sendBitcoin("mzJR1zKcZCZvMJj87rVqmRPUYbgJyFxtSL", 0.001)
    if err != nil {
        log.Fatalf("Error sending Bitcoin: %v", err)
    }
}