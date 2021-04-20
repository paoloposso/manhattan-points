package points

import (
	"os"

	"github.com/bcicen/jstream"
	"github.com/mitchellh/mapstructure"
)

func GetPointsInsideManhattanDistance(origin Point, dist float64) ([]Point, error) {
	jsonFile, _ := os.Open("./data/points.json")
	decoder := jstream.NewDecoder(jsonFile, 1)

	var points []Point

	for mv := range decoder.Stream() {
		var point Point
		mapstructure.Decode(mv.Value, &point)

		if point.DistanceFromOrigin(origin) <= dist {
			points = append(points, point)
		}
	}

	return points, nil
}
