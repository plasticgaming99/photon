// Dyntypes proceeds dynamic-typing on Go.

// (c)opyright 2023- plasticgaming99
// licensed under MIT license
package dyntypes

import (
	"strconv"
	"strings"
)

// DynType has 5 types, int, str, f32, f64, bool
// If type check occured error, It always return str.
// float is maybe not supported but can detect
func CheckDynType(toCheck string) string {
	var err error
	if toCheck == "True" || toCheck == "true" || toCheck == "False" || toCheck == "false" {
		return "bool"
	} else {
		_, err = strconv.Atoi(toCheck)
		if err == nil {
			return "int"
		} else if strings.Contains(toCheck, ".") {
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
		} else {
			return "str"
		}
	}
}

func IsDynTypeMatch(targ1 string, targ2 string) bool {
	return CheckDynType(targ1) == targ2
}

// Returns bool from dynbool. If input does
// not matches DynType, return false.
func DynBool(input string) bool {
	if CheckDynType(input) == "bool" {
		if input == "False" {
			return false
		} else if input == "True" {
			return true
		}
	} else if CheckDynType(input) == "int" {
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

// Returns string from dynstr. It always return string
// You can use any dyntypes as string directly so this is unneeded
// example: DynInt -> string, DynBool -> string
func DynStr(input string) string {
	return input
}

// Returns integer from dynint. If Dyntype
// does not match, It always return int(0).
func DynInt(input string) int {
	if !(CheckDynType(input) == "int") {
		return 0
	} else {
		int, err := strconv.Atoi(input)
		if err != nil {
			return 0
		}
		return int
	}
	return 0
}
