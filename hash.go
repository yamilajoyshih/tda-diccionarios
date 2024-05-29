package diccionario

import (
	"fmt"
	TDALista "tdas/lista"
)

//https://golangprojectstructure.com/hash-functions-go-code/#how-are-hash-functions-used-in-the-map-data-structure

const (
	uint64Offset uint64 = 0xcbf29ce484222325
	uint64Prime  uint64 = 0x00000100000001b3
	tamIni       int    = 10
	factorPred   int    = 6 / 5 //tiene que ser mayor a 1
)

func fvnHash(data []byte) (hash uint64) {
	hash = uint64Offset

	for _, b := range data {
		hash ^= uint64(b)
		hash *= uint64Prime

	}
	return
}

// transforma un tipo de dato genérico a un array de bytes
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

type parClaveValor[K comparable, V any] struct {
	clave K
	dato  V
}

type hashAbierto[K comparable, V any] struct {
	tabla    []TDALista.Lista[*parClaveValor[K, V]]
	cantidad int
	tam      int
}
type iteradorDiccionario[K comparable, V any] struct {
	diccionario      *hashAbierto[K, V]
	actual           int
	actualLista      TDALista.Lista[*parClaveValor[K, V]]
	actParClaveValor TDALista.IteradorLista[*parClaveValor[K, V]]
}

func (hash *hashAbierto[K, V]) Cantidad() int {
	return hash.cantidad
}

// creo la tabla con el tamaño inicial
func CrearHash[K comparable, V any]() Diccionario[K, V] {
	return &hashAbierto[K, V]{tabla: make([]TDALista.Lista[*parClaveValor[K, V]], tamIni)}
}

func crearparClaveValor[K comparable, V any](clave K, dato V) *parClaveValor[K, V] {
	return &parClaveValor[K, V]{clave, dato}
}

func obtenerPosicion[K comparable](clave K, largo int) uint64 {
	return fvnHash(convertirABytes(clave)) % uint64(largo)
}

// toma un número k como parámetro y devuelve el siguiente número primo mayor que k
func buscarNumPrimo(k int) int {
	for i := k + 1; ; i++ {
		if esPrimo(i) {
			return i
		}
	}
}

func esPrimo(num int) bool {
	if num <= 1 || num%2 == 0 || num%3 == 0 {
		return false
	}

	i := 5
	for i*i <= num {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}

// redimensiono una tabla de hash abierta
func (hash *hashAbierto[K, V]) redimensionar(largo int) {
	nuevaTabla := make([]TDALista.Lista[*parClaveValor[K, V]], largo)
	for _, lista := range hash.tabla {
		if lista == nil {
			continue
		}
		for iter := lista.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
			actual := iter.VerActual()
			posicion := obtenerPosicion(actual.clave, largo)
			if nuevaTabla[posicion] == nil {
				nuevaTabla[posicion] = TDALista.CrearListaEnlazada[*parClaveValor[K, V]]()
			}
			nuevaTabla[posicion].InsertarUltimo(crearparClaveValor(actual.clave, actual.dato))
		}
	}
	hash.tabla = nuevaTabla
}

// Guardar guarda el par clave-dato en el Diccionario. Si la clave ya se encontraba, se actualiza el dato asociado
func (hash *hashAbierto[K, V]) Guardar(clave K, dato V) {
	posicion := obtenerPosicion(clave, len(hash.tabla))
	if hash.tabla[posicion] == nil {
		hash.tabla[posicion] = TDALista.CrearListaEnlazada[*parClaveValor[K, V]]()
	}
	for iter := hash.tabla[posicion].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		if iter.VerActual().clave != clave {
			continue
		}
		iter.VerActual().dato = dato
		return
	}
	nuevoPar := crearparClaveValor(clave, dato)
	hash.tabla[posicion].InsertarUltimo(nuevoPar)
	hash.cantidad++
	factordeCarga := float64(hash.Cantidad()) / float64(len(hash.tabla))
	if factordeCarga >= 2 {
		nuevolargo := buscarNumPrimo(hash.Cantidad() * factorPred)
		hash.redimensionar(nuevolargo)
	}
}

// Pertenece determina si una clave ya se encuentra en el diccionario, o no
func (hash *hashAbierto[K, V]) Pertenece(clave K) bool {
	posicion := obtenerPosicion(clave, len(hash.tabla))
	if hash.tabla[posicion] == nil {
		return false
	}
	for iter := hash.tabla[posicion].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		if iter.VerActual().clave == clave {
			return true
		}
	}
	return false
}

// Obtener devuelve el dato asociado a una clave. Si la clave no pertenece, debe entrar en pánico con mensaje
// 'La clave no pertenece al diccionario'
func (hash *hashAbierto[K, V]) Obtener(clave K) V {
	posicion := obtenerPosicion(clave, len(hash.tabla))
	if hash.tabla[posicion] == nil {
		panic("La clave no pertenece al diccionario")
	}
	for iter := hash.tabla[posicion].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		parClaveValor := iter.VerActual()
		if parClaveValor.clave == clave {
			return parClaveValor.dato
		}
	}
	panic("La clave no pertenece al diccionario")
}

// Borrar borra del Diccionario la clave indicada, devolviendo el dato que se encontraba asociado. Si la clave no
// pertenece al diccionario, debe entrar en pánico con un mensaje 'La clave no pertenece al diccionario'
func (hash *hashAbierto[K, V]) Borrar(clave K) V {
	posicion := obtenerPosicion(clave, len(hash.tabla))
	if hash.tabla[posicion] == nil {
		panic("La clave no pertenece al diccionario")
	}
	for iter := hash.tabla[posicion].Iterador(); iter.HaySiguiente(); iter.Siguiente() {
		parClaveValor := iter.VerActual()
		if parClaveValor.clave != clave {
			continue
		}
		dato := parClaveValor.dato
		iter.Borrar()
		hash.cantidad--
		factor := float32(hash.Cantidad()) / float32(len(hash.tabla))
		aux := hash.Cantidad() * factorPred

		if hash.tabla[posicion].Largo() == 0 {
			hash.tabla[posicion] = nil
		}
		if factor <= 0.2 && aux > hash.tam {
			hash.redimensionar(buscarNumPrimo(aux))
		}
		return dato
	}
	panic("La clave no pertenece al diccionario")
}

// Iterar itera internamente el diccionario, aplicando la función pasada por parámetro a todos los elementos del
// mismo
func (hash *hashAbierto[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	for _, lista := range hash.tabla {
		if lista == nil {
			continue
		}
		for iter := lista.Iterador(); iter.HaySiguiente(); iter.Siguiente() {
			parClaveValor := iter.VerActual()
			if !visitar(parClaveValor.clave, parClaveValor.dato) {
				return
			}
		}
	}
}

// Iterador devuelve un IterDiccionario para este Diccionario
func (hash *hashAbierto[K, V]) Iterador() IterDiccionario[K, V] {
	iterador := &iteradorDiccionario[K, V]{diccionario: hash}
	for i, lista := range iterador.diccionario.tabla {
		if lista == nil { //test iterar dicc vacio
			continue
		}
		iterador.actual = i
		iterador.actualLista = lista
		iterador.actParClaveValor = iterador.actualLista.Iterador()
		break
	}
	return iterador
}

// PRIMITIVAS DE IterDiccionario
func (iter *iteradorDiccionario[K, V]) HaySiguiente() bool {
	return iter.actualLista != nil && iter.actParClaveValor.HaySiguiente()
}

func (iter *iteradorDiccionario[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	parClaveValor := iter.actParClaveValor.VerActual()
	return parClaveValor.clave, parClaveValor.dato
}

func (iter *iteradorDiccionario[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	iter.actParClaveValor.Siguiente()
	if !iter.actParClaveValor.HaySiguiente() {
		for i := iter.actual + 1; i < len(iter.diccionario.tabla); i++ {
			if iter.diccionario.tabla[i] != nil || len(iter.diccionario.tabla)-1 == i {
				iter.actual = i
				iter.actualLista = iter.diccionario.tabla[i]
				if iter.diccionario.tabla[i] != nil {
					iter.actParClaveValor = iter.actualLista.Iterador()
				}
				return
			}
		}
	}
}
