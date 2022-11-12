package three

import (
	"fmt"
	"math"
)

type Branch struct {
	input float64
	// width     float64
	nextThree []Branch
}

var Three = []Branch{}

func Add(prevInput float64, input float64) {
	branch := Branch{input: input, nextThree: []Branch{}}
	if prevInput != 0 {
		closestInput := Branch{input: 0, nextThree: []Branch{}}

		for _, value := range Three {
			if value.input == input {
				return
			}

			if closestInput.input == 0 || math.Abs(closestInput.input-prevInput) > math.Abs(value.input-prevInput) {
				closestInput = value
			}

		}

		closestInput.nextThree = append(closestInput.nextThree, branch)

		// fmt.Println(input, closestInput.input, closestInput.nextThree)
	}
	Three = append(Three, branch)
	fmt.Println(Three)
	// return r.width * r.height
}
