package main

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/urfave/cli"
	"rsc.io/qr"
)

var version string

func main() {
	app := cli.NewApp()
	app.Name = "qrcode"
	app.Usage = "generate a QR code from a text file"
	app.Version = version
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "text file to convert to QR code",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "qr.png",
			Usage: "file to output",
		},
	}

	app.Action = func(c *cli.Context) error {
		err := qrcode(c.GlobalString("file"), c.GlobalString("output"))
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("QR code from '%s' written to '%s'\n", c.GlobalString("file"), c.GlobalString("output"))
		}
		return err
	}

	app.Run(os.Args)
}

func qrcode(readFile, outputFile string) (err error) {
	b, err := ioutil.ReadFile(readFile)
	if err != nil {
		return
	}
	q, err := qr.Encode(transformTo(b), qr.L)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(outputFile, q.PNG(), 0644)
	return
}

func transformTo(b []byte) string {
	compressed := compressByte(b)
	encrypted, _ := encrypt(compressed, "pass")
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println(len(b), len(encoded))
	return encoded
}

func transformFrom(s string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(s)
	decrypted, _ := decrypt(decoded, "pass")
	return decompressByte(decrypted)
}

// compressByte returns a compressed byte slice.
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

// decompressByte returns a decompressed byte slice.
func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

// compress uses flate to compress a byte slice to a corresponding level
func compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

// compress uses flate to decompress an io.Reader
func decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}

func encrypt(plaintext []byte, passphrase string) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))
	passphrase += hex.EncodeToString(hasher.Sum(nil))
	key := []byte(passphrase)[:32]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, passphrase string) ([]byte, error) {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))
	passphrase += hex.EncodeToString(hasher.Sum(nil))
	key := []byte(passphrase)[:32]

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
