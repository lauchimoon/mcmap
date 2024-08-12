package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

const (
	BlockData = `GRASS|127, 178, 56, 255|Grass Block
SAND|247, 233, 163, 255|Sand
WOOL|199, 199, 199, 255|Cobweb
FIRE|255, 0, 0, 255|Redstone Block
ICE|160, 160, 255, 255|Ice
METAL|167, 167, 167, 255|Block of Iron
PLANT|0, 124, 0, 255|Wheat
SNOW|255, 255, 255, 255|White Wool
CLAY|164, 168, 184, 255|Clay
DIRT|151, 109, 77, 255|Dirt
STONE|112, 112, 112, 255|Stone
WATER|64, 64, 255, 255|Water
WOOD|143, 119, 72, 255|Oak Planks
QUARTZ|255, 252, 245, 255|Quartz Block
COLOR_ORANGE|216, 127, 51, 255|Orange Wool
COLOR_MAGENTA|178, 76, 216, 255|Magenta Wool
COLOR_LIGHT_BLUE|102, 153, 216, 255|Light Blue Wool
COLOR_YELLOW|229, 229, 51, 255|Yellow Wool
COLOR_LIGHT_GREEN|127, 204, 25, 255|Lime Wool
COLOR_PINK|242, 127, 165, 255|Pink Wool
COLOR_GRAY|76, 76, 76, 255|Gray Wool
COLOR_LIGHT_GRAY|153, 153, 153, 255|Light Gray Wool
COLOR_CYAN|76, 127, 153, 255|Cyan Wool
COLOR_PURPLE|127, 63, 178, 255|Purple Wool
COLOR_BLUE|51, 76, 178, 255|Blue Wool
COLOR_BROWN|102, 76, 51, 255|Brown Wool
COLOR_GREEN|102, 127, 51, 255|Green Wool
COLOR_RED|153, 51, 51, 255|Red Wool
COLOR_BLACK|25, 25, 25, 255|Black Wool
GOLD|250, 238, 77, 255|Block of Gold
DIAMOND|92, 219, 213, 255|Block of Diamond
LAPIS|74, 128, 255, 255|Block of Lapis Lazuli
EMERALD|0, 217, 58, 255|Block of Emerald
PODZOL|129, 86, 49, 255|Podzol
NETHER|112, 2, 0, 255|Netherrack
TERRACOTTA_WHITE|209, 177, 161, 255|White Terracotta
TERRACOTTA_ORANGE|159, 82, 36, 255|Orange Terracotta
TERRACOTTA_MAGENTA|149, 87, 108, 255|Magenta Terracotta
TERRACOTTA_LIGHT_BLUE|112, 108, 138, 255|Light Blue Terracotta
TERRACOTTA_YELLOW|186, 133, 36, 255|Yellow Terracotta
TERRACOTTA_LIGHT_GREEN|103, 117, 53, 255|Lime Terracotta
TERRACOTTA_PINK|160, 77, 78, 255|Pink Terracotta
TERRACOTTA_GRAY|57, 41, 35, 255|Gray Terracotta
TERRACOTTA_LIGHT_GRAY|135, 107, 98, 255|Light Gray Terracotta
TERRACOTTA_PURPLE|122, 73, 88, 255|Purple Terracotta
TERRACOTTA_BLUE|76, 62, 92, 255|Blue Terracotta
TERRACOTTA_BROWN|76, 50, 35, 255|Brown Terracotta
TERRACOTTA_GREEN|76, 82, 42, 255|Green Terracotta
TERRACOTTA_RED|142, 60, 46, 255|Red Terracotta
TERRACOTTA_BLACK|37, 22, 16, 255|Black Terracotta
CRIMSON_NYLIUM|189, 48, 49, 255|Crimson Nylium
CRIMSON_STEM|148, 63, 97, 255|Crimson Stem
CRIMSON_HYPHAE|92, 25, 29, 255|Crimson Hyphae
WARPED_NYLIUM|22, 126, 134, 255|Warped Nylium
WARPED_STEM|58, 142, 140, 255|Warped Stem
WARPED_HYPHAE|86, 44, 62, 255|Warped Hyphae
WARPED_WART_BLOCK|20, 180, 133, 255|Warped Wart Block
DEEPSLATE|100, 100, 100, 255|Deepslate
RAW_IRON|216, 175, 147, 255|Block of Raw Iron
GLOW_LICHEN|127, 167, 150, 255|Glow Lichen`
)

type MapBlock struct {
	Name  string
	Color color.RGBA
	Block string
}

type Triple struct {
	First  color.Color
	Second MapBlock
	Third  bool
}

type Point struct {
	First  int
	Second int
}

func NewMapBlock(name string, color color.RGBA, block string) MapBlock {
	return MapBlock{
		Name:  name,
		Color: color,
		Block: block,
	}
}

func NewPoint(x, y int) Point {
	return Point{
		First:  x,
		Second: y,
	}
}

func StringToUint8(s string) uint8 {
	n, _ := strconv.Atoi(s)
	return uint8(n)
}

func loadMapBlocks() []MapBlock {
	mapBlocks := []MapBlock{}

	lines := strings.Split(BlockData, "\n")
	for _, line := range lines {
		fields := strings.Split(line, "|")
		name := fields[0]
		colorComponents := strings.Split(fields[1], ", ")
		color := color.RGBA{
			StringToUint8(colorComponents[0]),
			StringToUint8(colorComponents[1]),
			StringToUint8(colorComponents[2]),
			StringToUint8(colorComponents[3]),
		}
		block := fields[2]

		mapBlocks = append(mapBlocks, NewMapBlock(name, color, block))
	}

	return mapBlocks
}

func resizeImage(i image.Image, w, h int) image.Image {
	rect := image.Rect(0, 0, w, h)
	dst := image.NewRGBA(rect)
	draw.NearestNeighbor.Scale(dst, rect, i, i.Bounds(), draw.Over, nil)

	return dst
}

func areColorsSimilar(c1, c2 color.Color) (bool, int64) {
	var epsilon int64 = 550
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	u8r1 := uint8(r1)
	u8r2 := uint8(r2)
	u8g1 := uint8(g1)
	u8g2 := uint8(g2)
	u8b1 := uint8(b1)
	u8b2 := uint8(b2)

	var rDiff uint8 = 0
	if u8r1 > u8r2 {
		rDiff = u8r1 - u8r2
	} else {
		rDiff = u8r2 - u8r1
	}

	var gDiff uint8 = 0
	if u8g1 > u8g2 {
		gDiff = u8g1 - u8g2
	} else {
		gDiff = u8g2 - u8g1
	}

	var bDiff uint8 = 0
	if u8b1 > u8b2 {
		bDiff = u8b1 - u8b2
	} else {
		bDiff = u8b2 - u8b1
	}

	sum := int64(rDiff) + int64(gDiff) + int64(bDiff)
	return sum < epsilon, sum
}

func mostSimilar(t []Triple) MapBlock {
	if len(t) == 1 {
		return t[0].Second
	}

	_, minimum := areColorsSimilar(t[0].Second.Color, t[0].First)
	mb := t[0].Second
	for _, tp := range t {
		_, sum := areColorsSimilar(tp.Second.Color, tp.First)
		if minimum > sum {
			minimum = sum
			mb = tp.Second
		}
	}

	return mb
}

func generateMap(img image.Image) map[Point]MapBlock {
	//func generateMap(img image.Image) map[Point][]Triple {
	// Find similar pixels on mapblocks
	// Method:
	// - check each pixel against all 61 mapblocks' color and check if colors are similar (done by areColorsSimilar)
	// - store the similarities as triples (color, mapblock color, similar?) in a map of slices
	// - the mapblock associated to that pixel will be picked if the two colors are similar for each slice
	finalMap := map[Point]MapBlock{}
	mapBlocks := loadMapBlocks()
	differences := map[Point][]Triple{}
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			for _, mb := range mapBlocks {
				color := img.At(x, y)
				similar, _ := areColorsSimilar(mb.Color, color)

				t := Triple{color, mb, similar}
				p := NewPoint(x, y)

				if similar {
					differences[p] = append(differences[p], t)
				}
			}
		}
	}

	for p, t := range differences {
		finalMap[p] = mostSimilar(t)
	}

	return finalMap
}

func encodeMap(m map[Point]MapBlock) {
	newImage := image.NewRGBA(image.Rect(0, 0, 128, 128))
	for p, mb := range m {
		newImage.Set(p.First, p.Second, mb.Color)
	}

	exported, err := os.Create("image-map.png")
	if err != nil {
		log.Fatalf("there was an error: %v", err)
	}
	defer exported.Close()

	png.Encode(exported, newImage)
}

func main() {
	args := os.Args
	filename := args[1]
	rawImage, err := os.Open(filename)
	if err != nil {
		log.Fatalf("there was an error: %v", err)
	}
	defer rawImage.Close()

	img, err := png.Decode(rawImage)
	if err != nil {
		log.Fatalf("there was an error: %v", err)
	}

	img = resizeImage(img, 128, 128)
	mapp := generateMap(img)
	encodeMap(mapp)
}
