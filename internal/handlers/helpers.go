package handlers

import "fmt"

func FrequencyConvert(frequency string) (int, error) {
	switch frequency {
	case "D":
		return 60 * 24, nil
	case "3D":
		return 60 * 24 * 3, nil
	case "W":
		return 60 * 24 * 7, nil
	case "2W":
		return 60 * 24 * 7 * 2, nil
	case "3W":
		return 60 * 24 * 7 * 3, nil
	case "M":
		return 60 * 24 * 7 * 4, nil
	default:
		return 0, fmt.Errorf("only D, 3D, W, 2W, 3W, M")
	}
}
