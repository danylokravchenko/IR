package main

import(
	"fmt"
	"github.com/dotcypress/phonetics"
)

func main() {
	fmt.Println(phonetics.EncodeSoundex("Miller")) //M460
	fmt.Println(phonetics.EncodeSoundex("Muller")) //M460
}