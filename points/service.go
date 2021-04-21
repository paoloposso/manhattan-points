package points

import (
	"os"
	"sort"

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

		point.CalculateDistanceFromOrigin(origin)

		if point.distFromOrigin <= dist {
			points = append(points, point)
		}
	}

	sort.Slice(points[:], func(i, j int) bool {
		return points[i].distFromOrigin < points[j].distFromOrigin
	})

	return points, nil
}
