package dyntypes

import (
	"strconv"
	"strings"
)

// DynType has 5 types, int, str, f32, f64, bool
// If type check occured error, It always return str.
func checkDynType(toCheck string) string {
	var err error
	if toCheck == "True" || toCheck == "False" {
		return "bool"
	} else {
		if strings.Contains(toCheck, ".") {
			_, err = strconv.ParseFloat(toCheck, 64)
			if err != nil {
				_, err = strconv.ParseFloat(toCheck, 32)
				if err != nil {
					return "str"
				} else {
					return "f32"
				}
			} else {
				return "f64"
			}
		}
		_, err = strconv.Atoi(toCheck)
		if err != nil {
			return "int"
		} else {
			return "str"
		}
	}
}

func IsDynTypeMatch(targ1 string, targ2 string) bool {
	return checkDynType(targ1) == targ2
}

// Return bool from dynbool. If input does
// not matches DynType, return false.
func DynBool(input string) bool {
	if checkDynType(input) == "bool" {
		if input == "False" {
			return false
		} else if input == "True" {
			return true
		}
	} else if checkDynType(input) == "int" {
		if input == "1" {
			return true
		} else if input == "0" {
			return false
		}
	} else {
		return false
	}
	return false
}
