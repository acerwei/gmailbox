package encoder

import (
	"crypto/cipher"

	"golang.org/x/crypto/blowfish"
)

func blowfishChecksizeAndPad(pt []byte) []byte {
	modulus := (len(pt) + 1) % blowfish.BlockSize
	if modulus != 0 {
		padnglen := blowfish.BlockSize - modulus
		for i := 0; i < padnglen; i++ {
			pt = append(pt, 0)
		}
		pt = append(pt, byte(padnglen))
	} else {
		pt = append(pt, 0)
	}
	return pt
}

func blowfishDecrypt(et []byte, cp *blowfish.Cipher) []byte {
	div := et[:blowfish.BlockSize]
	decrypted := et[blowfish.BlockSize:]
	if len(decrypted)%blowfish.BlockSize != 0 {
		panic("decrypted is not a multiple of blowfish.BlockSize")
	}
	dcbc := cipher.NewCBCDecrypter(cp, div)
	dcbc.CryptBlocks(decrypted, decrypted)
	n := len(decrypted)
	padlen := int(decrypted[n-1])
	return decrypted[:n-padlen-1]
}

func blowfishEncrypt(ppt []byte, cp *blowfish.Cipher) []byte {
	et := make([]byte, blowfish.BlockSize+len(ppt))
	eiv := et[:blowfish.BlockSize]
	ecbc := cipher.NewCBCEncrypter(cp, eiv)
	ecbc.CryptBlocks(et[blowfish.BlockSize:], ppt)
	return et
}

const (
	maxKeyLengthC  = 56
	defaultEncKeyC = "1234567890qwertyuiopasdfghjklzxcvbnmASGHHJK"
)

//BlowFishEncoder BlowFishEncoder
type BlowFishEncoder struct {
	cipher *blowfish.Cipher
}

//NewBlowFishEncoder NewBlowFishEncoder
func NewBlowFishEncoder(enckey string) (*BlowFishEncoder, error) {
	if len(enckey) > maxKeyLengthC {
		enckey = enckey[:maxKeyLengthC]
	}
	if enckey == "" {
		enckey = defaultEncKeyC
	}
	key := []byte(enckey)
	cp, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &BlowFishEncoder{cipher: cp}, nil
}

//Encode Encode
func (e *BlowFishEncoder) Encode(data []byte) []byte {
	data = blowfishChecksizeAndPad(data)
	return blowfishEncrypt(data, e.cipher)
}

//Decode Decode
func (e *BlowFishEncoder) Decode(data []byte) []byte {
	return blowfishDecrypt(data, e.cipher)
}
