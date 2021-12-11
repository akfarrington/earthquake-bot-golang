package taiwanmap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// declare some constants for the image i'm working with

const degreeToLong float64 = 279.065
const degreeToLat float64 = 256.364
const baseLong float64 = 21.8898
const baseLat float64 = 119.2021
const picWidth int = 950
const picHeight int = 1000
const eqStationBoxSize int = 12
const eqStationBoxBorder int = 2
const epicenterIconSize int = 40

// Changeable allows me to change an image easily
type Changeable interface {
	Set(x, y int, c color.Color)
}

// GetBasePic returns an image variable of the taiwan map
// then later I can add stuff
func GetBasePic() image.Image {
	f, err := os.Open("eqmap.png")
	if err != nil {
		fmt.Println("error opening image")
		os.Exit(1)
	}

	image, err := png.Decode(f)
	if err != nil {
		fmt.Println("error decoding image")
		os.Exit(1)
	}

	return image
}

// MarkEpicenter adds a wavy kind of icon to where the epicenter is
func MarkEpicenter(tmap image.Image, epicenter Coords) image.Image {
	if cimg, ok := tmap.(Changeable); ok {
		// just return if the epicenter is outside of the image
		if epicenter.X > 950 || epicenter.X < 0 || epicenter.Y > 1000 || epicenter.Y < 0 {
			fmt.Println("epicenter out of bounds")
			fmt.Println(epicenter)
			return tmap
		}

		f, err := os.Open("eq-epi.png")
		if err != nil {
			fmt.Println("failed to open the epicenter png")
			os.Exit(1)
		}
		epicenterIcon, err := png.Decode(f)
		if err != nil {
			fmt.Println("failed to decode the epicenter icon")
		}

		iconStartx := epicenter.X - (epicenterIconSize / 2)
		iconStarty := epicenter.Y - (epicenterIconSize / 2)

		for xi := 0; xi < epicenterIconSize; xi++ {
			for yi := 0; yi < epicenterIconSize; yi++ {
				// get the color of the pixel at xi, yi
				pixColor := epicenterIcon.At(xi, yi)

				iconr, icong, iconb, icona := pixColor.RGBA()

				// combine pix colors and handle transparency
				if icona > 0 {
					// this pixel is somewhat transparent..
					// first get the pixels rgb from the target pixel on image
					tpicr, tpicg, tpicb, _ := tmap.At(xi+iconStartx, yi+iconStarty).RGBA()
					// get a weight in percent from icona / 65535
					// image.Image uses a max of 65535, not 255 like normal people use
					iconweight := float32(icona) / 65535.0
					tpicweight := 1.0 - iconweight
					// find the kind of weighted average for the new pixel
					// multiply by 255, then divide by 65535 to convert to the regular 0-255 for rgba values
					targetr := uint8((((float32(iconr) * iconweight) + (float32(tpicr) * tpicweight)) * 255) / 65535)
					targetg := uint8((((float32(icong) * iconweight) + (float32(tpicg) * tpicweight)) * 255) / 65535)
					targetb := uint8((((float32(iconb) * iconweight) + (float32(tpicb) * tpicweight)) * 255) / 65535)

					// now set the pixel with the new averaged out pixel
					cimg.Set(xi+iconStartx, yi+iconStarty, color.RGBA{targetr, targetg, targetb, 255})
				}
			}
		}
	}
	return tmap
}

// MarkPicList takes a list of coordinates and adds to the image
func MarkPicList(tmap image.Image, coordList []Coords) image.Image {
	for _, loc := range coordList {
		markPicLoc(tmap, loc)
	}

	return tmap
}

func markPicLoc(tmap image.Image, loc Coords) {
	// shouldn't have an eq out of bounds of the picture. If so, something's
	// messed up and the program should die
	if cimg, ok := tmap.(Changeable); ok {
		if loc.X > picWidth || loc.Y > picHeight {
			fmt.Println("eq out of bounds in picture")
			fmt.Println(loc.X, loc.Y)
			os.Exit(1)
		}

		// make a box (under, but ends up being a border)
		backBoxSize := eqStationBoxSize + (eqStationBoxBorder * 2)
		boxstartx := loc.X - backBoxSize/2
		boxstarty := loc.Y - backBoxSize/2
		// for loop for black background
		for xi := 0; xi < backBoxSize; xi++ {
			for yi := 0; yi < backBoxSize; yi++ {
				cimg.Set(boxstartx+xi, boxstarty+yi, color.Black)
			}
		}

		// this is the inner, colored box
		boxstartx = loc.X - eqStationBoxSize/2
		boxstarty = loc.Y - eqStationBoxSize/2

		// https://stackoverflow.com/questions/36573413/change-color-of-a-single-pixel-golang-image
		for xi := 0; xi < eqStationBoxSize; xi++ {
			for yi := 0; yi < eqStationBoxSize; yi++ {
				if loc.Intensity >= 7 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{130, 0, 140, 255}) // dark purple
				} else if loc.Intensity >= 6 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{240, 0, 255, 255}) // light purple
				} else if loc.Intensity >= 5 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{255, 0, 0, 255}) // red
				} else if loc.Intensity >= 4 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{255, 130, 0, 255}) // orange
				} else if loc.Intensity >= 3 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{255, 220, 0, 255}) // yellow
				} else if loc.Intensity >= 2 {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{0, 140, 0, 255}) // dark green
				} else {
					cimg.Set(boxstartx+xi, boxstarty+yi, color.RGBA{0, 190, 0, 255}) // light green
				}
			}
		}
	}
}
