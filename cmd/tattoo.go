package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/edmundofuentes/tattoo/internal"
	"log"
	"math/rand"
	"os"
	"time"
)

const CONFIG_FILE string = "config.toml"

type Connection struct {
	a int
	b int
}

func main() {
	// Read Config TOML file
	cfg := internal.Config{}

	_, err := toml.DecodeFile(CONFIG_FILE, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var seed int64
	if cfg.Seed != 0 {
		seed = cfg.Seed
	} else {
		seed = time.Now().UTC().UnixNano()
	}

	rand.Seed(seed)
	fmt.Printf("Attempting design from seed %d ..\n", seed)

	// STEP 1: Generate the design
	constellation := internal.Generate(cfg)

	// debug the constellation
	//fmt.Printf("%V\n", constellation)

	// STEP 2: Draw it
	image := internal.Draw(cfg, constellation)

	os.MkdirAll("./output", os.ModePerm)
	filename := fmt.Sprintf("./output/%d.png", seed)

	err = image.SavePNG(filename)

	if err != nil {
		fmt.Println("Error writing output PNG file. Check permissions and try again.")
	}

	fmt.Printf("Success! Image generated on output/%d.png\n", seed)

	os.Exit(0)
}
