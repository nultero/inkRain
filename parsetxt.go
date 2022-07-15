package main

import (
	"os"
	"strings"
)

type coord struct {
	x, y int
}

type coordmap map[coord]struct{}

func getCoords(file string, minX, minY int) coordmap {
	bytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	coords := map[coord]struct{}{}

	cont := string(bytes)
	lines := strings.Split(cont, "\n")

	for y, ln := range lines {
		for x, r := range ln {
			if r == '|' {
				c := coord{x: x + minX/5, y: y + minY/5}
				coords[c] = struct{}{}
			}
		}
	}

	return coords
}
