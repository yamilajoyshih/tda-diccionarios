package diccionario

import (
	TDAPila "tdas/pila"
)

type abb[K comparable, V any] struct {
	raiz *nodoAbb[K, V]
	cant int
	cmp  func(K, K) int
}
type nodoAbb[K comparable, V any] struct {
	izq, der *nodoAbb[K, V]
	clave    K
	dato     V
}

type iterABB[K comparable, V any] struct {
	pila         *TDAPila.Pila[*nodoAbb[K, V]]
	actual       *nodoAbb[K, V]
	abb          *abb[K, V]
	desde, hasta *K
	cmp          func(K, K) int
}

// Crear ABB
func CrearABB[K comparable, V any](funcion_cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: funcion_cmp}
}

// Crear Nodo
func crearNodo[K comparable, V any](clave K, dato V) *nodoAbb[K, V] {
	return &nodoAbb[K, V]{izq: nil, der: nil, clave: clave, dato: dato}
}

// Funcion Primitiva: Cantidad
func (abb *abb[K, V]) Cantidad() int {
	return abb.cant
}

// Funcion Primitiva: Guardar
func (abb *abb[K, V]) Guardar(clave K, dato V) {
	if abb.raiz == nil {
		abb.raiz = crearNodo[K, V](clave, dato)
		abb.cant++
	} else {
		abb.raiz.guardarWrapper(clave, dato, abb.cmp, abb)
	}
}

// Wrapper del guardar
func (nodo *nodoAbb[K, V]) guardarWrapper(clave K, dato V, cmp func(K, K) int, abb *abb[K, V]) {
	cmpResultado := cmp(clave, nodo.clave)
	//Un entero menor que 0 si la primera clave es menor que la segunda.
	if cmpResultado < 0 {
		if nodo.izq == nil {
			nodo.izq = crearNodo[K, V](clave, dato)
			abb.cant++
		} else {
			nodo.izq.guardarWrapper(clave, dato, cmp, abb)
		}
	} else if cmpResultado > 0 { //Un entero mayor que 0 si la primera clave es mayor que la segunda.
		if nodo.der == nil {
			nodo.der = crearNodo[K, V](clave, dato)
			abb.cant++
		} else {
			nodo.der.guardarWrapper(clave, dato, cmp, abb)
		}
	} else { //son iguales
		nodo.dato = dato
	}
}

// Funcion Primitiva: Pertenece
func (abb *abb[K, V]) Pertenece(clave K) bool {
	_, encontrado := abb.raiz.buscarClave(clave, abb.cmp)
	return encontrado
}

// Funcion Primitiva: Obtener
func (abb *abb[K, V]) Obtener(clave K) V {
	dato, encontrado := abb.raiz.buscarClave(clave, abb.cmp)
	if !encontrado {
		panic("La clave no pertenece al diccionario")
	}
	return *dato
}

// Wrapper para buscar
func (nodo *nodoAbb[K, V]) buscarClave(clave K, cmp func(K, K) int) (*V, bool) {
	if nodo == nil {
		return nil, false
	}
	cmpResultado := cmp(clave, nodo.clave)
	if cmpResultado < 0 {
		return nodo.izq.buscarClave(clave, cmp)
	} else if cmpResultado > 0 {
		return nodo.der.buscarClave(clave, cmp)
	} else {
		return &nodo.dato, true
	}
}

func (abb *abb[K, V]) Borrar(clave K) V {
	_, encontrado := abb.raiz.buscarClave(clave, abb.cmp)
	if !encontrado {
		panic("La clave no pertenece al diccionario")
	}
	valorBorrado := abb.Obtener(clave)
	abb.cant--
	abb.raiz = abb.raiz.borrarWrapper(clave, abb.cmp)
	return valorBorrado
}

func (nodo *nodoAbb[K, V]) borrarWrapper(clave K, cmp func(K, K) int) *nodoAbb[K, V] {
	if nodo == nil {
		return nil
	}
	cmpResultado := cmp(clave, nodo.clave)
	if cmpResultado < 0 {
		nodo.izq = nodo.izq.borrarWrapper(clave, cmp)
	} else if cmpResultado > 0 {
		nodo.der = nodo.der.borrarWrapper(clave, cmp)
	} else {
		if nodo.izq == nil {
			return nodo.der
		}
		if nodo.der == nil {
			return nodo.izq
		}
		//si ambos son nulos, tengo que buscar quien va a sustiyuir el nodo actual
		reemplazante := reemplaza(nodo.der)
		nodo.clave, nodo.dato = reemplazante.clave, reemplazante.dato
		nodo.der = nodo.der.borrarWrapper(reemplazante.clave, cmp)
	}
	return nodo
}

func reemplaza[K comparable, V any](n *nodoAbb[K, V]) *nodoAbb[K, V] {
	for n.izq != nil {
		n = n.izq
	}
	return n
}

func (abb *abb[K, V]) Iterador() IterDiccionario[K, V] {
	return abb.IteradorRango(nil, nil)
}

func (abb *abb[K, V]) Iterar(visitar func(clave K, dato V) bool) {
	iterarwrapper(abb.raiz, visitar)
}

func iterarwrapper[K comparable, V any](nodo *nodoAbb[K, V], visitar func(clave K, dato V) bool) bool {
	if nodo == nil {
		return true
	}
	if !iterarwrapper(nodo.izq, visitar) || !visitar(nodo.clave, nodo.dato) || !iterarwrapper(nodo.der, visitar) {
		return false
	}
	return true
}

func (abb *abb[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	iterarRangoWrapper(abb.raiz, desde, hasta, visitar, abb.cmp)
}

func iterarRangoWrapper[K comparable, V any](nodo *nodoAbb[K, V], desde, hasta *K, visitar func(clave K, dato V) bool, cmp func(K, K) int) bool {
	if nodo == nil {
		return true
	}
	siguiente := true
	if desde == nil || cmp(*desde, nodo.clave) < 0 {
		siguiente = iterarRangoWrapper(nodo.izq, desde, hasta, visitar, cmp)
	}
	if siguiente && (desde == nil || cmp(*desde, nodo.clave) <= 0) && (hasta == nil || cmp(*hasta, nodo.clave) >= 0) {
		siguiente = visitar(nodo.clave, nodo.dato)
	}
	if !siguiente {
		return false
	}
	if hasta == nil || cmp(*hasta, nodo.clave) > 0 {
		siguiente = iterarRangoWrapper(nodo.der, desde, hasta, visitar, cmp)
	}
	return siguiente
}

func (abb *abb[K, V]) IteradorRango(desde *K, hasta *K) IterDiccionario[K, V] {
	pila := TDAPila.CrearPilaDinamica[*nodoAbb[K, V]]()
	abb.raiz.rango(&pila, desde, hasta, abb.cmp)
	iterador := &iterABB[K, V]{abb: abb, pila: &pila, desde: desde, hasta: hasta, cmp: abb.cmp}
	return iterador
}

func (nodo *nodoAbb[K, V]) rango(pila *TDAPila.Pila[*nodoAbb[K, V]], desde, hasta *K, cmp func(K, K) int) {
	if nodo == nil {
		return
	}
	if desde == nil || cmp(*desde, nodo.clave) <= 0 {
		(*pila).Apilar(nodo)
		nodo.izq.rango(pila, desde, hasta, cmp)
	} else if hasta == nil || cmp(*hasta, nodo.clave) > 0 {
		nodo.der.rango(pila, desde, hasta, cmp)
	}
}

// Primitivas del iterador
func (iter *iterABB[K, V]) HaySiguiente() bool {
	return !(*iter.pila).EstaVacia() && (iter.hasta == nil || iter.cmp(*iter.hasta, (*iter.pila).VerTope().clave) >= 0)
}

func (iter *iterABB[K, V]) VerActual() (K, V) {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	nodo := (*iter.pila).VerTope()
	return nodo.clave, nodo.dato
}

func (iter *iterABB[K, V]) Siguiente() {
	if !iter.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	nodo := (*iter.pila).Desapilar()
	if nodo.der == nil {
		return
	}
	nodo.der.rango(iter.pila, iter.desde, iter.hasta, iter.cmp)
}
