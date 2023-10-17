package util

import (
	"strings"

	"crypto/ed25519"
	"encoding/hex"
	com "github.com/aureontu/MRWebServer/mr_services/common"
	"github.com/aureontu/MRWebServer/mr_services/mpberr"
)

func ReadNonceFromAptosFullMsg(msg string) string {
	for _, v := range strings.Split(msg, "\n") {
		if strings.HasPrefix(v, "nonce:") {
			nonce := v[7:]
			for len(nonce) < com.NonceLen {
				nonce = "0" + nonce
			}
			return nonce
		}
	}
	return ""
}

func EncodeAptosPubKey(pubKey []byte) string {
	return "0x" + hex.EncodeToString(pubKey)
}

func DecodeAptosPubKey(pubKeyHex string) ([]byte, error) {
	if pubKeyHex[:2] == "0x" {
		pubKeyHex = pubKeyHex[2:]
	}
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func VerifySignature(publicKey []byte, msg string, signatureHex string) bool {
	if signatureHex[:2] == "0x" {
		signatureHex = signatureHex[:2]
	}
	var signature = make([]byte, 64)
	n, err := hex.Decode(signature, []byte(signatureHex))
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, []byte(msg), signature[:n])
}

func FixHashId(hashId string) (string, error) {
	if len(hashId) == 0 {
		return "", mpberr.ErrNFTHashId
	}
	hashId = strings.ToLower(hashId)

	if strings.HasPrefix(hashId, "0x") {
		hashId = hashId[2:]
	}

	if len(hashId) == 0 || len(hashId) > com.AptosNFTTokenIdLen {
		return "", mpberr.ErrNFTHashId
	}

	if len(hashId) == com.AptosNFTTokenIdLen {
		return "0x" + hashId, nil
	}

	var fixZero = make([]rune, com.AptosNFTTokenIdLen-len(hashId))
	for i := range fixZero {
		fixZero[i] = '0'
	}

	return "0x" + string(fixZero) + hashId, nil
}
