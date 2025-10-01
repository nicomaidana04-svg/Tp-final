package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	constCantFilasTablero    = 25
	constCantColumnasTablero = 29

	constCantColumnas = 2
	constY            = 0
	constX            = 1

	constCantColumnasOvni = 4
	constTipoOvni         = 0
	constOvniY            = 1
	constOvniX            = 2
	constEnDescenso       = 3

	constTiempoDeDisparoOvni   = 3
	constTiempoLiberarcionOvni = 10

	constSimboloVacío       = ""
	constSimboloNave        = "N"
	constSimboloDisparoNave = "*"
	constSimboloDisparoOvni = "."
	constSimboloOvniLider   = "L"
	constSimboloOvniComun   = "C"
	constSimboloBorde       = "X"

	constCantColumnasDisparos = 2
)

// Vector global con las direcciones posibles
var (
	quieto    = [constCantColumnas]int{0, 0}
	izquierda = [constCantColumnas]int{0, -1}
	derecha   = [constCantColumnas]int{0, 1}
	arriba    = [constCantColumnas]int{-1, 0}
	abajo     = [constCantColumnas]int{1, 0}
)

// Vector global con las direccion de la nave
var direccionNave [constCantColumnas]int

// Variable global que indica si se presiono la barra espaciadora lo que ejecuta un disparo de la nave
var disparoNave bool

// Función para enviar actualizaciones a los clientes
func generarEventos() {
	var (
		tablero [constCantFilasTablero][constCantColumnasTablero]string

		nave         [constCantColumnas]int
		disparosNave [][constCantColumnas]int

		ovnis         [][constCantColumnasOvni]int
		disparosOvnis [][constCantColumnas]int

		ultimaEjecucionDisparoOvni    time.Time
		ultimaEjecucionLiberacionOvni time.Time

		puntos int
		vidas  int
	)

	rand.Seed(time.Now().Unix())

	//Se inicializa variables
	ultimaEjecucionDisparoOvni = time.Now()
	ultimaEjecucionLiberacionOvni = time.Now()

	disparoNave = false

	vidas = 3

	// Se genera tablero por primera vez con los bordes
	tablero = generarTablero()

	// Se genera la nave (posición inicial) por primera vez
	nave, direccionNave = inicializarNave(constCantFilasTablero, constCantColumnasTablero)

	// Se generan los ovnis (posiciones iniciales) por primera vez
	ovnis = inicializarOvnis(constCantFilasTablero, constCantColumnasTablero)

	// Se actualiza nave y ovnis en el tablero por primera vez
	actualizarTablero(&tablero, nave, disparosNave, ovnis, disparosOvnis)

	for {
		// Se actualizan las posiciones de la nave según la dirección
		calcularNuevaPosicionNave(tablero, &nave, &direccionNave)

		// Se crea un nuevo disparo si corresponde
		crearDisparoNave(nave, &disparoNave, &disparosNave)

		//Cada "constTiempoDeDisparoOvni" segundos, se crea un disparo de un ovni
		if time.Since(ultimaEjecucionDisparoOvni) >= constTiempoDeDisparoOvni*time.Second {
			crearDisparoOvni(ovnis, &disparosOvnis)
			ultimaEjecucionDisparoOvni = time.Now()
		}

		//Cada "constTiempoLiberarcionOvni" segundos, se libera un obvni de la formación
		if time.Since(ultimaEjecucionLiberacionOvni) >= constTiempoLiberarcionOvni*time.Second {
			liberarOvni(ovnis)
			ultimaEjecucionLiberacionOvni = time.Now()
		}

		// Se calcula la nueva posición de los ovnis liberados
		calcularNuevaPosicionOvnisLiberados(ovnis)

		// Se calcula las nuevas posiciones de los disparos de la nave y de los ovnis
		calcularNuevasPosicionesDisparos(tablero, disparosNave, disparosOvnis)

		// Se verifica el estado del juego y eliminan elementos si corresponde
		if !verificarEstadoDeJuego(tablero, nave, &ovnis, &disparosNave, &disparosOvnis, &puntos) {
			// Si no tiene más vidas, se devuelve pantalla gameOver
			vidas--

			if vidas <= 0 {
				enviarGameOver(puntos)
				return
			}
		} else {
			if len(ovnis) == 0 {
				enviarWin(puntos)

				return
			}

			enviarActualizacionTexto(fmt.Sprint("Puntaje: ", puntos, ". Vidas: ", vidas))
		}

		//Se actualiza el tablero con los valores de la nave, ovnis y disparos en sus nuevas posiciones
		actualizarTablero(&tablero, nave, disparosNave, ovnis, disparosOvnis)

		// Se envía actualización de tablero al cliente para mostrar en pantalla
		enviarActualizacionTablero(tablero)

		// Espera un tiempo antes de generar un nuevo movimiento
		time.Sleep(150 * time.Millisecond)
	}
}

func generarTablero() [constCantFilasTablero][constCantColumnasTablero]string {
	var tablero [constCantFilasTablero][constCantColumnasTablero]string

	//PROGRAMAR
	for f := 0; f < constCantFilasTablero; f++ {
		for c := 0; c < constCantColumnasTablero; c++ {
			// Bordes en X
			if f == 0 || f == constCantFilasTablero-1 || c == 0 || c == constCantColumnasTablero-1 {
				tablero[f][c] = constSimboloBorde
			} else {
				tablero[f][c] = constSimboloVacío
			}
		}
	}

	return tablero
}

func inicializarNave(cantFilasTablero int, cantColumnasTablero int) ([constCantColumnas]int, [constCantColumnas]int) {
	//PROGRAMAR
	posicion := [constCantColumnas]int{
		constY: cantFilasTablero - 2,
		constX: cantColumnasTablero / 2,
	}

	return posicion, quieto
	//return [constCantColumnas]int{}, quieto
}

func inicializarOvnis(cantFilasTablero int, cantColumnasTablero int) [][constCantColumnasOvni]int {
	var (
		ovnis [][constCantColumnasOvni]int
	)
	for f := 1; f <= 3; f++ {
		for c := 2; c < cantColumnasTablero-2; c += 2 {
			ovni := [constCantColumnasOvni]int{
				constTipoOvni:   0,
				constOvniY:      f,
				constOvniX:      c,
				constEnDescenso: 0,
			}
			ovnis = append(ovnis, ovni)
		}
	}
	//PROGRAMAR

	return ovnis
}

func actualizarTablero(tablero *[constCantFilasTablero][constCantColumnasTablero]string,
	nave [constCantColumnas]int,
	disparosNave [][constCantColumnas]int,
	ovnis [][constCantColumnasOvni]int,
	disparosOvnis [][constCantColumnas]int) {

	for f := 0; f < constCantFilasTablero; f++ {
		for c := 0; c < constCantColumnasTablero; c++ {
			if f == 0 || f == constCantFilasTablero-1 || c == 0 || c == constCantColumnasTablero-1 {
				// esto asegura que siempre hay borde en X
				(*tablero)[f][c] = constSimboloBorde
			} else {
				// todo lo demás arranca vacío
				(*tablero)[f][c] = constSimboloVacío
			}
		}
	}

	// Pongo la nave en el tablero
	naveFila := nave[0]
	naveColumna := nave[1]
	(*tablero)[naveFila][naveColumna] = constSimboloNave

	// Ahora dibujo los ovnis
	for i := 0; i < len(ovnis); i++ {
		ovniFila := ovnis[i][constOvniY]
		ovniColumna := ovnis[i][constOvniX]
		ovniTipo := ovnis[i][constTipoOvni]

		// redundancia: elijo símbolo con ifs separados
		if ovniTipo == 0 {
			(*tablero)[ovniFila][ovniColumna] = constSimboloOvniComun
		}
		if ovniTipo == 1 {
			(*tablero)[ovniFila][ovniColumna] = constSimboloOvniLider
		}
	}

	// Redundancia extra: los disparos todavía no se usan, pero igual los recorro
	if len(disparosNave) > 0 {
		for i := 0; i < len(disparosNave); i++ {
			dy := disparosNave[i][0]
			dx := disparosNave[i][1]
			(*tablero)[dy][dx] = constSimboloDisparoNave
		}
	}

	if len(disparosOvnis) > 0 {
		for i := 0; i < len(disparosOvnis); i++ {
			dy := disparosOvnis[i][0]
			dx := disparosOvnis[i][1]
			(*tablero)[dy][dx] = constSimboloDisparoOvni
		}
	}

	//PROGRAMAR
}

func calcularNuevaPosicionNave(tablero [constCantFilasTablero][constCantColumnasTablero]string,
	nave *[constCantColumnas]int, direccionNave *[constCantColumnas]int) {
	var nuevax, nuevay int

	nuevax = nave[constX] + direccionNave[constX]
	nuevay = nave[constY] + direccionNave[constY]

	if nuevax >= 0 && nuevax < constCantFilasTablero && nuevay >= 0 && nuevay < constCantColumnasTablero {

		if tablero[nuevax][nuevay] != constSimboloBorde {
			nave[constX] = nuevax
			nave[constY] = nuevay
		} else {

			*direccionNave = quieto

		}
	}

	//PROGRAMAR

}

func crearDisparoNave(nave [constCantColumnas]int,
	disparoNave *bool,
	disparosNave *[][constCantColumnasDisparos]int) {

	//PROGRAMAR
}

func crearDisparoOvni(ovnis [][constCantColumnasOvni]int,
	disparosOvnis *[][constCantColumnasDisparos]int) {

	//PROGRAMAR
}

func calcularNuevasPosicionesDisparos(tablero [constCantFilasTablero][constCantColumnasTablero]string,
	disparosNave [][constCantColumnasDisparos]int,
	disparosOvnis [][constCantColumnasDisparos]int) {

	//PROGRAMAR
}

func verificarEstadoDeJuego(tablero [constCantFilasTablero][constCantColumnasTablero]string,
	nave [constCantColumnas]int,
	ovnis *[][constCantColumnasOvni]int,
	disparosNave *[][constCantColumnasDisparos]int,
	disparosOvnis *[][constCantColumnasDisparos]int,
	puntos *int) bool {

	//PROGRAMAR

	return true
}

func eliminarDisparo(slice [][constCantColumnasDisparos]int, coordenadaY int, coordenadaX int) [][2]int {
	var nuevoSlice [][constCantColumnasDisparos]int
	for f := 0; f < len(slice); f++ {
		if slice[f][constY] != coordenadaY &&
			slice[f][constX] != coordenadaX {
			nuevoSlice = append(nuevoSlice, slice[f])
		}
	}
	return nuevoSlice
}

func eliminarOvni(slice [][constCantColumnasOvni]int, coordenadaY int, coordenadaX int) [][4]int {
	var nuevoSlice [][constCantColumnasOvni]int
	for f := 0; f < len(slice); f++ {
		if slice[f][constOvniY] != coordenadaY ||
			slice[f][constOvniX] != coordenadaX {
			nuevoSlice = append(nuevoSlice, slice[f])
		}
	}
	return nuevoSlice
}

func liberarOvni(ovnis [][constCantColumnasOvni]int) {
	//PROGRAMAR
}

func calcularNuevaPosicionOvnisLiberados(ovnis [][constCantColumnasOvni]int) {
	//PROGRAMAR
}

func Personalizacion() {

}
