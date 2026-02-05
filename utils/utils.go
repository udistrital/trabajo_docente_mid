package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/udistrital/trabajo_docente_mid/models"
)

// Checkeo de Id si cumple
//   - param: string de Id tabla relacional
//
// Retorna:
//   - int64 si id valido o cero
//   - error si no es valido
func CheckIdInt(param string) (int64, error) {
	paramInt, err := strconv.ParseInt(param, 10, 64)
	if paramInt <= 0 && err == nil {
		err = fmt.Errorf("no valid Id: %d > 0 = false", paramInt)
	}
	return paramInt, err
}

// Checkeo de _id si cumple
//   - param: string de _id tabla no relacional
//
// Retorna:
//   - string si _id valido o " "
//   - error si no es valido
func CheckIdString(param string) (string, error) {
	pattern := `^[0-9a-fA-F]{24}$`
	regex := regexp.MustCompile(pattern)
	if regex.MatchString(param) {
		return param, nil
	} else {
		return "", fmt.Errorf("no valid Id: %s", param)
	}
}

// Formatea data en base a modelo de datos; ver en ~/models
//   - data: data en interface{}
//   - &tipo: variable con tipo de dato especificado (no olvide "&")
//
// Retorna por referencia:
//   - la data en &tipo
func ParseData(data interface{}, tipo interface{}) {
	inbytes, err := json.Marshal(data)
	if err == nil {
		json.Unmarshal(inbytes, &tipo)
	} else {
		fmt.Println(err)
	}
}

// Formatea nombre completo de tercero como nombre1 nombre2 apellido1 apellido2
//   - tercero: registro de tercero
//
// Retorna:
//   - string con nombre completo
func FormatNameTercero(tercero models.Tercero) string {
	nombreFormateado := ""
	if tercero.PrimerNombre != "" {
		str := strings.ToLower(tercero.PrimerNombre)
		nombreFormateado += strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.SegundoNombre != "" {
		str := strings.ToLower(tercero.SegundoNombre)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.PrimerApellido != "" {
		str := strings.ToLower(tercero.PrimerApellido)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if tercero.SegundoApellido != "" {
		str := strings.ToLower(tercero.SegundoApellido)
		nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
	}
	if nombreFormateado == "" {
		splittedStr := strings.Split(strings.ToLower(tercero.NombreCompleto), " ")
		for _, str := range splittedStr {
			if str != "" {
				nombreFormateado += " " + strings.ToUpper(str[0:1]) + str[1:]
			}
		}
	}
	return strings.Trim(nombreFormateado, " ")
}

// Formatea cualquier texto a primera letra mayúscula
//   - badString: texto cualquiera
//
// Retorna:
//   - string capitalizado
func Capitalize(badString string) string {
	if badString == "" {
		return ""
	}
	str := strings.ToLower(badString)
	return strings.ToUpper(str[0:1]) + str[1:]
}

// Une lo que hay igual en dos arreglos, tomado como referencia el primero
//   - A: Array que se evalua en su totalidad, es el valor que toma
//   - B: Array para comparar, en coincidencia hace break
//   - from: "A" toma valor de A, "B" toma valor de B
//   - compare: Función que recibe item de A y B, debe retornar el parametro a comparar en forma de interface{}
//
// Retorna:
//   - Array con iguales de A en B
func JoinEqual(A, B interface{}, from string, compare func(item interface{}) interface{}) (igual []interface{}) {
	igual = []interface{}{}
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _A.Len(); i++ {
		for j := 0; j < _B.Len(); j++ {
			if compare(_A.Index(i).Interface()) == compare(_B.Index(j).Interface()) {
				if from == "B" {
					igual = append(igual, _B.Index(j).Interface())
				} else {
					igual = append(igual, _A.Index(i).Interface())
				}
				break
			}
		}
	}
	return igual
}

// Extrae lo que hay diferente en dos arreglos, tomado como referencia el primero
//   - A: Array para comparar con JoinEqual(A, B)
//   - B: Array para comparar con JoinEqual(A, B)
//   - from: "A" toma valor de A, "B" toma valor de B, from "AB" toma valor de A y B
//   - compare: Función que recibe item de A y B, debe retornar el parametro a comparar en forma de interface{}
//
// Retorna:
//   - Array con la diferencia de A o B o AB
func SubstractDiff(A, B interface{}, from string, compare func(item interface{}) interface{}) (diferente []interface{}) {
	diferente = []interface{}{}
	_I := reflect.ValueOf(JoinEqual(A, B, "", compare))
	_A := reflect.ValueOf(A)
	_B := reflect.ValueOf(B)
	for i := 0; i < _I.Len(); i++ {
		if from == "A" || from == "AB" {
			for j := 0; j < _A.Len(); j++ {
				if compare(_I.Index(i).Interface()) != compare(_A.Index(j).Interface()) {
					diferente = append(diferente, _A.Index(j).Interface())
				}
			}
		}
		if from == "B" || from == "AB" {
			for k := 0; k < _B.Len(); k++ {
				if compare(_I.Index(i).Interface()) != compare(_B.Index(k).Interface()) {
					diferente = append(diferente, _B.Index(k).Interface())
				}
			}
		}
	}
	return diferente
}

// Encuentra una coincidencia en un arreglo, tomando como validador una funcion de filtro dada
//   - A: Array para buscar según criterio de función
//   - compare: Función que recibe item de A, debe retornar booleano resultado de comprobación de item
//
// Retorna:
//   - item interface{} si coincide o nil en caso contrario
func Find(A interface{}, compare func(item interface{}) bool) interface{} {
	_A := reflect.ValueOf(A)
	for i := 0; i < _A.Len(); i++ {
		if compare(_A.Index(i).Interface()) {
			return _A.Index(i).Interface()
		}
	}
	return nil
}

// Remueve duplicados en un arreglo, tomando como validador una funcion de filtro dada
//   - A: Array para buscar según criterio de función
//   - compare: Función que recibe item de A, debe retornar el parametro de comparación
//
// Retorna:
//   - Array con los items únicos de A
func RemoveDuplicated(A interface{}, compare func(item interface{}) interface{}) (unicos []interface{}) {
	_A := reflect.ValueOf(A)
	mapeoUnicos := make(map[interface{}]bool)
	unicos = []interface{}{}
	for i := 0; i < _A.Len(); i++ {
		if _, encontrado := mapeoUnicos[compare(_A.Index(i).Interface())]; !encontrado {
			mapeoUnicos[compare(_A.Index(i).Interface())] = true
			unicos = append(unicos, _A.Index(i).Interface())
		}
	}
	return unicos
}

// SplitTrimSpace use Split function to slices s into all substrings separated
// by sep and use TrimSpace to remove space and return a slice of the substrings.
func SplitTrimSpace(s, sep string) []string {
	substrings := strings.Split(s, sep)

	for i, elementString := range substrings {
		substrings[i] = strings.TrimSpace(elementString)
	}
	return substrings
}

func GetOrDefault(value interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return defaultValue
		}
	case int:
		if v == 0 {
			return defaultValue
		}
		// Añadir más tipos según sea necesario
	}
	return value
}
