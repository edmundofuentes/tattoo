package main

import (
	"encoding/json"
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

	// Initialize folders
	os.MkdirAll("./output", os.ModePerm)
	os.MkdirAll("./output/img", os.ModePerm)
	os.MkdirAll("./output/json", os.ModePerm)

	// STEP 1: Generate the design
	rand.Seed(seed)
	fmt.Printf("Attempting design from seed %d ..\n", seed)

	constellation := internal.Generate(cfg)

	// write the constellation json
	b, err := json.Marshal(constellation)
	if err != nil {
		fmt.Printf("Error marshalling constellation to JSON: %v\n", err)
		return
	}

	jsonFilename := fmt.Sprintf("./output/json/%d.json", seed)
	err = os.WriteFile(jsonFilename, b, os.ModePerm)

	if err != nil {
		fmt.Println("Error writing output JSON file. Check permissions and try again.")
	}

	// debug the constellation
	//fmt.Printf("%V\n", constellation)

	// STEP 2: Draw it
	image := internal.Draw(cfg, constellation)

	// write the image
	imageFilename := fmt.Sprintf("./output/img/%d.png", seed)
	err = image.SavePNG(imageFilename)

	if err != nil {
		fmt.Println("Error writing output PNG file. Check permissions and try again.")
	}

	fmt.Printf("Success! Image generated on output/%d.png\n", seed)

	os.Exit(0)
}
