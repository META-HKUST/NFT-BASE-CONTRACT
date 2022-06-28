package main

import(
	 "fmt"
	 "crypto/sha256"
	 "encoding/hex"
)


func main() {
	hash := sha256.New()
	hash.Write([]byte("x509:cn , sdasdah= s"))
	s := hex.EncodeToString(hash.Sum(nil))
	tokenURI := s +"-" + "8"
	fmt.Println(tokenURI)
}

 


