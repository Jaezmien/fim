package utilities

import "strconv"

func FloatAsString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func StringAsFloat(str string) float64 {
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic("strconv.ParseFloat threw an error")
	}
	return float
}

func BoolAsString(b bool) string {
	if b == true {
		return "true"
	}
	return "false"
}
func StringAsBool(str string) bool {
	return str == "true"
}
