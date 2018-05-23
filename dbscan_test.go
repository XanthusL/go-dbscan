package dbscan

import (
	"fmt"
	"log"
	"math"
	"testing"
	"image"
	"image/color"
	"image/png"
	"os"
	"math/rand"
	"image/draw"
)

type SimpleClusterable struct {
	position float64
}

func (s SimpleClusterable) Distance(c interface{}) float64 {
	distance := math.Abs(c.(SimpleClusterable).position - s.position)
	return distance
}

func (s SimpleClusterable) GetID() string {
	return fmt.Sprint(s.position)
}

func TestPutAll(t *testing.T) {
	testMap := make(map[string]Clusterable)
	clusterList := []Clusterable{
		SimpleClusterable{10},
		SimpleClusterable{12},
	}
	putAll(testMap, clusterList)
	mapSize := len(testMap)
	if mapSize != 2 {
		t.Errorf("Map does not contain expected size 2 but was %d", mapSize)
	}
}

//Test find neighbour function
func TestFindNeighbours(t *testing.T) {
	log.Println("Executing TestFindNeighbours")
	clusterList := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{-1},
		SimpleClusterable{1.5},
		SimpleClusterable{-0.5},
	}

	eps := 1.0
	neighbours := findNeighbours(clusterList[0], clusterList, eps)

	assertEquals(t, 3, len(neighbours))
}

func TestMerge(t *testing.T) {
	log.Println("Executing TestMerge")
	expected := 6
	one := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{2.1},
		SimpleClusterable{2.2},
		SimpleClusterable{2.3},
	}

	two := []Clusterable{
		one[0],
		one[1],
		SimpleClusterable{2.4},
	}

	output := merge(one, two)
	assertEquals(t, expected, len(output))
}

func TestExpandCluster(t *testing.T) {
	log.Println("Executing TestExpandCluster")
	expected := 4
	clusterList := []Clusterable{
		SimpleClusterable{0},
		SimpleClusterable{1},
		SimpleClusterable{2},
		SimpleClusterable{2.1},
		SimpleClusterable{5},
	}

	eps := 1.0
	minPts := 3
	visitMap := make(map[string]bool)
	cluster := make(Cluster, 0)
	cluster = expandCluster(cluster, clusterList, visitMap, minPts, eps)
	assertEquals(t, expected, len(cluster))
}

func TestClusterize(t *testing.T) {
	log.Println("Executing TestClusterize")
	clusterList := []Clusterable{
		SimpleClusterable{1},
		SimpleClusterable{0.5},
		SimpleClusterable{0},
		SimpleClusterable{5},
		SimpleClusterable{4.5},
		SimpleClusterable{4},
	}
	eps := 1.0
	minPts := 2
	clusters := Clusterize(clusterList, minPts, eps)
	assertEquals(t, 2, len(clusters))
	if 2 == len(clusters) {
		assertEquals(t, 3, len(clusters[0]))
		assertEquals(t, 3, len(clusters[1]))
	}
}

func TestClusterizeNoData(t *testing.T) {
	log.Println("Executing TestClusterizeNoData")
	clusterList := []Clusterable{}
	eps := 1.0
	minPts := 3
	clusters := Clusterize(clusterList, minPts, eps)
	assertEquals(t, 0, len(clusters))
}

//Assert function. If  the expected value not equals result, function
//returns error.
func assertEquals(t *testing.T, expected, result int) {
	if expected != result {
		t.Errorf("Expected %d but got %d", expected, result)
	}
}

type P struct {
	Lng float64
	Lat float64
	ID  int64
}

func (p *P) GetID() int64 {
	return p.ID
}

func (p *P) Distance(t interface{}) float64 {
	target := t.(*P)
	f := p.Lat - target.Lat
	f2 := p.Lng - target.Lng
	return math.Sqrt(f*f + f2*f2)
}

func TestOPTISC(t *testing.T) {
	const count = 60
	const size = 500
	src := make([]Point, count)
	for i := int64(0); i < count; i++ {
		p := &P{
			ID:  i,
			Lat: float64(rand.Int31n(size)),
			Lng: float64(rand.Int31n(size)),
		}
		if i%10 < 3 {
			p.Lng = float64(i*size/count + (int64(p.Lng))%10)
			p.Lat = float64(i*size/count + (int64(p.Lng))%6)
		}

		src[i] = p
	}

	result, distances := OPTISC(src, 3, 30)
	uniform := image.NewUniform(color.White)
	gray := image.NewGray(image.Rect(0, 0, size, size))
	draw.Draw(gray, gray.Bounds(), uniform, image.ZP, draw.Src)
	points := image.NewGray(image.Rect(0, 0, size, size))
	draw.Draw(points, points.Bounds(), uniform, image.ZP, draw.Src)
	for i, r := range result {
		t.Logf("%#v %f", r, distances[i])
		x := i * size / count
		y := int(math.Ceil(distances[i]))
		// draw result
		drawPoint(x, y, gray)

		// draw the origin point
		pX := int(r.(*P).Lng)
		pY := int(r.(*P).Lat)
		drawPoint(pX, pY, points)
	}
	file, _ := os.Create("result.png")
	pointsF, _ := os.Create("points.png")

	e := png.Encode(file, gray)
	t.Log(e)
	e = png.Encode(pointsF, points)
	t.Log(e)
	t.Log(len(result))
	file.Close()
	pointsF.Close()
}

func drawPoint(x, y int, img draw.Image) {
	// 4 points as a bigger point
	y = img.Bounds().Dy() - y
	img.Set(x, y, color.Black)
	img.Set(x, y+1, color.Black)
	img.Set(x+1, y, color.Black)
	img.Set(x+1, y+1, color.Black)

	// image:
	//  ------->x
	//  |
	//  |
	// \|/
	//  Y  y
}
