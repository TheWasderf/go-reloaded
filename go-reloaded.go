package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) == 3 {
		input, err := os.ReadFile(os.Args[1])
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}
		text := string(input)

		array := strings.Fields(text)
		for i := range array {
			array[i] = LastControl(array[i])
		}

		var changearr []string
		for i, word := range array {
			if word == "(cap)" {
				if len(changearr) > 0 {
					changearr[len(changearr)-1] = Captalized(changearr[len(changearr)-1])
				}
			} else if word == "(bin)" || word == "(bin)," {
				if len(changearr) > 0 {
					changearr[len(changearr)-1] = Binary(changearr[len(changearr)-1])
				}
			} else if word == "(up)" || word == "(up)," {
				if len(changearr) > 0 {
					changearr[len(changearr)-1] = ToUpper(changearr[len(changearr)-1])
				}
			} else if word == "(low)" || word == "(low)," {
				if len(changearr) > 0 {
					changearr[len(changearr)-1] = ToLower(changearr[len(changearr)-1])
				}
			} else if word == "(hex)" || word == "(hex)," {
				if len(changearr) > 0 {
					changearr[len(changearr)-1] = HexDecimal(changearr[len(changearr)-1])
				}
			} else if word == "a" || word == "an" || word == "A" || word == "An" {
				if i+1 < len(array) {
					changearr = append(changearr, CorretMulti(array[i+1], array[i]))
				} else {
					changearr = append(changearr, word)
				}
			} else if word == "(up," || word == "(low," || word == "(cap," {
				numstr := strings.Trim(array[i+1], ")")
				num, _ := strconv.Atoi(numstr)
				for j := 0; j < num; j++ {
					in := len(changearr) - num + j
					if in >= 0 && in < len(changearr) {
						switch word {
						case "(up,":
							changearr[in] = ToUpper(changearr[in])
						case "(low,":
							changearr[in] = ToLower(changearr[in])
						case "(cap,":
							changearr[in] = Captalized(changearr[in])
						}
					}
				}
				i++
				continue
			} else if string(word[len(word)-1]) == ")" {
				continue
			} else {
				changearr = append(changearr, word)
			}
		}
		result := strings.Join(changearr, " ")
		result = CorrectPunctuation2(result)
		result = CorrectPunctuation(result)
		result = strings.TrimSpace(result)
		err = os.WriteFile(os.Args[2], []byte(result), 0o644)
		if err != nil {
			log.Fatalf("Error writing to output file: %v", err)
		}
	} else {
		log.Println("Usage: go run main.go <input_file> <output_file>")
	}
}

func CorrectPunctuation2(s string) string {
	// punct := regexp.MustCompile(`\s*([.,:;!?]+)\s*`)
	// line := punct.ReplaceAllString(s, "$1 ")
	// return line
	// s = strings.ReplaceAll(s, ".", ". ")
	s = strings.ReplaceAll(s, " .", ".")
	s = strings.ReplaceAll(s, " ,", ", ")
	s = strings.ReplaceAll(s, ",  ", ", ")
	s = strings.ReplaceAll(s, " ;", ";")
	s = strings.ReplaceAll(s, " :", ":")
	s = strings.ReplaceAll(s, ":", ": ")
	s = strings.ReplaceAll(s, " !", "!")
	s = strings.ReplaceAll(s, " ?", "?")
	s = strings.ReplaceAll(s, " ...", "...")
	s = strings.ReplaceAll(s, ". . . ", "...")
	s = strings.ReplaceAll(s, " '", "'")
	s = strings.ReplaceAll(s, "' ", "'")
	s = strings.ReplaceAll(s, ". '", ".'")
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, s)
}

func LastControl(str string) string {
	str = strings.ReplaceAll(str, "\"", " ")
	str = strings.ReplaceAll(str, " \"", "\"")
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, str)
}

func CorrectPunctuation(s string) string {
	// Düzenli ifadeyi derle
	re := regexp.MustCompile(`\s*([,:;!?]+)\s*`)

	// Düzenli ifadeyi kullanarak noktalama işaretlerini düzelt
	s = re.ReplaceAllString(s, "$1 ")

	// Sonuçta kalan fazla boşlukları kaldır
	s = strings.TrimSpace(s)

	return s
}

func CorretMulti(s, bir string) string {
	vovel := []rune{'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U'}
	isVovel := false
	for _, char := range vovel {
		if rune(s[0]) == char {
			isVovel = true
			break
		}
	}
	if isVovel {
		if bir == "a" {
			bir = "an"
		} else if bir == "A" {
			bir = "An"
		}
	} else {
		if bir == "an" || bir == "An" {
			bir = "a"
		}
	}
	return bir
}

func Captalized(s string) string {
	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

func Binary(s string) string {
	üs := len(s) - 1
	sum := 0
	for i := 0; i < len(s); i++ {
		eklenecek := Power(2, üs)
		if s[i] == '1' {
			sum += eklenecek
		}
		üs--
	}
	return strconv.Itoa(sum)
}

func Power(base, power int) int {
	result := 1
	for i := 0; i < power; i++ {
		result *= base
	}
	return result
}

func HexDecimal(hexString string) string {
	var HexString string
	if len(hexString) == 3 {
		for i := 1; i <= len(hexString)-1; i++ {
			HexString += string(hexString[i])
		}
		decimal, _ := strconv.ParseInt(HexString, 16, 64)
		return strconv.Itoa(int(decimal))
	}
	decimal, _ := strconv.ParseInt(hexString, 16, 64)
	return strconv.Itoa(int(decimal))
}

func ToUpper(s string) string {
	sonuc := ""
	for _, karakter := range s {
		if karakter >= 'a' && karakter <= 'z' {
			sonuc += string(karakter - 32)
		} else {
			sonuc += string(karakter)
		}
	}
	return sonuc
}

func ToLower(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}
