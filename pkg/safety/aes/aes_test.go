package aes

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	aes := CreateCipher("PrZz5jDCGNydGfvi")
	decrypt, err := aes.AesBase64Decrypt("EQ9FGzB1LlvcwN8R2BhLwPeMjLXWubpMDepn+ORrqAQ=")
	if err != nil {
		panic(err)
	}
	fmt.Println(decrypt)

}
