package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Tnze/go-mc/nbt"
	"golang.org/x/image/draw"
)

const (
	BlockData = `6|GRASS|127, 178, 56, 255|Grass Block
10|SAND|247, 233, 163, 255|Sand
14|WOOL|199, 199, 199, 255|Cobweb
18|FIRE|255, 0, 0, 255|Redstone Block
22|ICE|160, 160, 255, 255|Ice
26|METAL|167, 167, 167, 255|Block of Iron
30|PLANT|0, 124, 0, 255|Wheat
34|SNOW|255, 255, 255, 255|White Wool
38|CLAY|164, 168, 184, 255|Clay
42|DIRT|151, 109, 77, 255|Dirt
46|STONE|112, 112, 112, 255|Stone
50|WATER|64, 64, 255, 255|Water
54|WOOD|143, 119, 72, 255|Oak Planks
58|QUARTZ|255, 252, 245, 255|Quartz Block
62|COLOR_ORANGE|216, 127, 51, 255|Orange Wool
66|COLOR_MAGENTA|178, 76, 216, 255|Magenta Wool
70|COLOR_LIGHT_BLUE|102, 153, 216, 255|Light Blue Wool
74|COLOR_YELLOW|229, 229, 51, 255|Yellow Wool
78|COLOR_LIGHT_GREEN|127, 204, 25, 255|Lime Wool
82|COLOR_PINK|242, 127, 165, 255|Pink Wool
86|COLOR_GRAY|76, 76, 76, 255|Gray Wool
90|COLOR_LIGHT_GRAY|153, 153, 153, 255|Light Gray Wool
94|COLOR_CYAN|76, 127, 153, 255|Cyan Wool
98|COLOR_PURPLE|127, 63, 178, 255|Purple Wool
102|COLOR_BLUE|51, 76, 178, 255|Blue Wool
106|COLOR_BROWN|102, 76, 51, 255|Brown Wool
110|COLOR_GREEN|102, 127, 51, 255|Green Wool
114|COLOR_RED|153, 51, 51, 255|Red Wool
118|COLOR_BLACK|25, 25, 25, 255|Black Wool
122|GOLD|250, 238, 77, 255|Block of Gold
126|DIAMOND|92, 219, 213, 255|Block of Diamond
130|LAPIS|74, 128, 255, 255|Block of Lapis Lazuli
134|EMERALD|0, 217, 58, 255|Block of Emerald
138|PODZOL|129, 86, 49, 255|Podzol
142|NETHER|112, 2, 0, 255|Netherrack
146|TERRACOTTA_WHITE|209, 177, 161, 255|White Terracotta
150|TERRACOTTA_ORANGE|159, 82, 36, 255|Orange Terracotta
154|TERRACOTTA_MAGENTA|149, 87, 108, 255|Magenta Terracotta
158|TERRACOTTA_LIGHT_BLUE|112, 108, 138, 255|Light Blue Terracotta
162|TERRACOTTA_YELLOW|186, 133, 36, 255|Yellow Terracotta
166|TERRACOTTA_LIGHT_GREEN|103, 117, 53, 255|Lime Terracotta
170|TERRACOTTA_PINK|160, 77, 78, 255|Pink Terracotta
174|TERRACOTTA_GRAY|57, 41, 35, 255|Gray Terracotta
178|TERRACOTTA_LIGHT_GRAY|135, 107, 98, 255|Light Gray Terracotta
182|TERRACOTTA_CYAN|187, 92, 92, 255|Cyan Terracotta
186|TERRACOTTA_PURPLE|122, 73, 88, 255|Purple Terracotta
190|TERRACOTTA_BLUE|76, 62, 92, 255|Blue Terracotta
194|TERRACOTTA_BROWN|76, 50, 35, 255|Brown Terracotta
198|TERRACOTTA_GREEN|76, 82, 42, 255|Green Terracotta
202|TERRACOTTA_RED|142, 60, 46, 255|Red Terracotta
206|TERRACOTTA_BLACK|37, 22, 16, 255|Black Terracotta
210|CRIMSON_NYLIUM|189, 48, 49, 255|Crimson Nylium
214|CRIMSON_STEM|148, 63, 97, 255|Crimson Stem
218|CRIMSON_HYPHAE|92, 25, 29, 255|Crimson Hyphae
222|WARPED_NYLIUM|22, 126, 134, 255|Warped Nylium
226|WARPED_STEM|58, 142, 140, 255|Warped Stem
230|WARPED_HYPHAE|86, 44, 62, 255|Warped Hyphae
234|WARPED_WART_BLOCK|20, 180, 133, 255|Warped Wart Block
238|DEEPSLATE|100, 100, 100, 255|Deepslate
242|RAW_IRON|216, 175, 147, 255|Block of Raw Iron
246|GLOW_LICHEN|127, 167, 150, 255|Glow Lichen`
)

type MapData struct {
	Version int32 `nbt:"DataVersion"`
	M       Map   `nbt:"data"`
}

type Map struct {
	ZCenter           int32  `nbt:"zCenter"`
	UnlimitedTracking byte   `nbt:"unlimitedTracking"`
	TrackingPosition  byte   `nbt:"trackingPosition"`
	Frames            []byte `nbt:"frames,list"`
	Scale             byte   `nbt:"scale"`
	Locked            byte   `nbt:"locked"`
	Dimension         string `nbt:"dimension"`
	Banners           []byte `nbt:"banners,list"`
	XCenter           int32  `nbt:"xCenter"`
	Colors            []byte `nbt:"colors"`
}

type MapBlock struct {
	Code  byte
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

func NewMapBlock(code byte, name string, color color.RGBA, block string) MapBlock {
	return MapBlock{
		Code:  code,
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
		code, _ := strconv.Atoi(fields[0])
		name := fields[1]
		colorComponents := strings.Split(fields[2], ", ")
		color := color.RGBA{
			StringToUint8(colorComponents[0]),
			StringToUint8(colorComponents[1]),
			StringToUint8(colorComponents[2]),
			StringToUint8(colorComponents[3]),
		}
		block := fields[3]

		mapBlocks = append(mapBlocks, NewMapBlock(byte(code), name, color, block))
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

func export(m map[Point]MapBlock, id uint) error {
	// Prepare to export map to nbt file
	exportMap := Map{
		ZCenter:           0,
		UnlimitedTracking: 0,
		TrackingPosition:  0,
		Frames:            []byte{},
		Scale:             0,
		Locked:            1,
		Dimension:         "minecraft:overworld",
		Banners:           []byte{},
		XCenter:           0,
	}

	// Place colors
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			p := NewPoint(x, y)
			exportMap.Colors = append(exportMap.Colors, m[p].Code)
		}
	}

	// Export
	exportData := MapData{3955, exportMap}
	var buffNBT bytes.Buffer
	encoder := nbt.NewEncoder(&buffNBT)
	err := encoder.Encode(exportData, "")
	if err != nil {
		return err
	}

	var buffGz bytes.Buffer
	gw := gzip.NewWriter(&buffGz)
	gw.Write(buffNBT.Bytes())
	gw.Close()

	f, err := os.Create(fmt.Sprintf("map_%d.dat", id))
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(buffGz.Bytes())

	return nil
}

func main() {
	log.SetFlags(0)
	args := os.Args
	progName := args[0]
	filename := flag.String("img", "", "png file to export into a map.")
	mapId := flag.Uint("id", 0, "set map id. default is 0")
	flag.Parse()

	if len(*filename) == 0 {
		log.Fatalf("usage: %s -img=<image.png>", progName)
	}

	rawImage, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("%s: %v", progName, err)
	}
	defer rawImage.Close()

	img, err := png.Decode(rawImage)
	if err != nil {
		log.Fatalf("%s: %v", progName, err)
	}

	img = resizeImage(img, 128, 128)
	mapp := generateMap(img)
	err = export(mapp, *mapId)
	if err != nil {
		log.Fatalf("%s: %v", progName, err)
	}
}
