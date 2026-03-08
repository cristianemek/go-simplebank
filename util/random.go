package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

//* Go 1.20 al usar math/rand ya usa una semilla aleatoria, no es necesario
// func init() { //funcion init se va a llamar automaticamente cuando se importe el paquete
// 	var r = rand.New(rand.NewSource(time.Now().UnixNano())) //creamos una nueva instancia de rand con una semilla basada en el tiempo actual

// }

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) //generamos un numero aleatorio entre min y max
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6) //generamos un string aleatorio de 6 caracteres para el nombre del propietario
}

func RandomMoney() int64 {
	return RandomInt(0, 1000) //generamos un numero aleatorio entre 0 y 1000 para el balance de la cuenta
}

func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)] //seleccionamos una moneda aleatoria de la lista de monedas
}
