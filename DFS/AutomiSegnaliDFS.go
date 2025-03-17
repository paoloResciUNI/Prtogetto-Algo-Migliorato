package main

import (
	"bufio"
	. "fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Strutture dati esistenti
type punto struct {
	coordinataX int
	coordinataY int
	id          string
	successivo  *punto
}

type piano *Piano

type Piano struct {
	automi   *punto
	ostacoli *punto
}

type elementoPila struct {
	chiamato *punto
	distanza int
	prossimo *elementoPila
}

// Struttura per lo sostamento 
type cella struct {
	x, y   int
	passi  int
	parent *cella
}

type visitata struct {
	x, y int
	next *visitata
}

// Funzioni principali del programma
func esegui(p piano, s string) {
	comandi := strings.Split(s, " ")
	switch comandi[0] {
	case "c":
		P := newPiano()
		p.automi = P.automi
		p.ostacoli = P.ostacoli
	case "S":
		(*p).stampa()
	case "s":
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).stato(a, b)
	case "a":
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).automa(a, b, comandi[3])
	case "o":
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		c, _ := strconv.Atoi(comandi[3])
		d, _ := strconv.Atoi(comandi[4])
		(*p).ostacolo(a, b, c, d)
	case "p":
		(*p).posizioni(comandi[1])
	case "r":
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).richiamo(a, b, comandi[3])
	case "e":
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).esistePercorso(a, b, comandi[3])
	case "f":
		os.Exit(0)
	}
}

func main() {
	var Campo Piano
	pCampo := &Campo
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		esegui(pCampo, scanner.Text())
	}
}

func (Campo Piano) stampa() {
	percorrente := new(punto)
	Println("(")
	percorrente = Campo.automi
	for percorrente != nil {
		Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		percorrente = percorrente.successivo
	}
	Println(")")
	Println("[")
	percorrente = Campo.ostacoli
	for percorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		Printf("(%d,%d)(%d,%d)\n", x0, y0, x1, y1)
		percorrente = percorrente.successivo
	}
	Println("]")
}

func (Campo Piano) stato(x, y int) {
	if Campo.cercaOstacolo(x, y) != nil {
		Println("O")
		return
	}
	if Campo.cercaAutoma(x, y, "") != nil {
		Println("A")
		return
	} else {
		Println("E")
	}
}

func (Campo Piano) posizioni(alpha string) {
	percorrente := new(punto)
	Println("(")
	percorrente = Campo.automi
	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		}
		percorrente = percorrente.successivo
	}
	Println(")")
}

// Implementazione modificata con DFS
func (Campo *Piano) esistePercorso(x, y int, eta string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	Sorgente.id = eta
	if Campo.cercaOstacolo(x, y) != nil {
		Println("NO")
		return
	}
	percorrente := Campo.cercaAutoma(x, y, eta)
	if percorrente == nil {
		Println("NO")
		return
	}

	// Utilizziamo il nuovo algoritmo DFS
	destinazione := new(punto)
	destinazione.coordinataX = x
	destinazione.coordinataY = y

	percorsoTrovato := dfs(Campo, percorrente, destinazione)
	if percorsoTrovato {
		Println("SI")
	} else {
		Println("NO")
	}
}

func newPiano() piano {
	var nuovPiano Piano
	return &nuovPiano
}

func (Campo *Piano) automa(x, y int, eta string) {
	puntoCercato := Campo.cercaAutoma(x, y, eta)
	if puntoCercato != nil {
		puntoCercato.coordinataX = x
		puntoCercato.coordinataY = y
	}
	if Campo.cercaOstacolo(x, y) == nil {
		puntoCercato = new(punto)
		puntoCercato.coordinataX = x
		puntoCercato.coordinataY = y
		puntoCercato.id = eta
		if Campo.automi == nil {
			Campo.automi = puntoCercato
			return
		}
		puntoCercato.successivo = Campo.automi
		Campo.automi = puntoCercato
	}
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	percorrente := Campo.automi
	for percorrente != nil {
		if (percorrente.coordinataX <= x1 && percorrente.coordinataX >= x0) && (percorrente.coordinataY <= y1 && percorrente.coordinataY >= y0) {
			return
		}
		percorrente = percorrente.successivo
	}
	newOstacolo := new(punto)
	newOstacolo.coordinataX = x0
	newOstacolo.coordinataY = y1
	newOstacolo.id = Sprintf("%d,%d,%d,%d,ostacolo", x0, y0, x1, y1)
	if Campo.ostacoli == nil {
		Campo.ostacoli = newOstacolo
		return
	}
	newOstacolo.successivo = Campo.ostacoli
	Campo.ostacoli = newOstacolo
}

// Implementazione modificata del richiamo con DFS
func (Campo *Piano) richiamo(x, y int, alpha string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	minDistance := math.MaxInt
	pilaChiamata := new(elementoPila)
	percorrente := Campo.automi

	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			distanza := calcolaDistanza(percorrente.coordinataX, percorrente.coordinataY, x, y)

			// Verifichiamo se esiste un percorso usando DFS
			destinazione := new(punto)
			destinazione.coordinataX = x
			destinazione.coordinataY = y

			percorsoTrovato := dfs(Campo, percorrente, destinazione)

			if percorsoTrovato {
				if distanza <= minDistance {
					minDistance = distanza
				}
				automaChiamato := new(elementoPila)
				automaChiamato.chiamato = percorrente
				automaChiamato.prossimo = pilaChiamata
				automaChiamato.distanza = distanza
				pilaChiamata = automaChiamato
			}
		}
		percorrente = percorrente.successivo
	}

	attraversoPila := pilaChiamata
	for attraversoPila != nil {
		if attraversoPila.distanza == minDistance {
			attraversoPila.chiamato.coordinataX = x
			attraversoPila.chiamato.coordinataY = y
		}
		attraversoPila = attraversoPila.prossimo
	}
}

func estraiCoordinate(id string) (x0 int, y0 int, x1 int, y1 int) {
	coordinate, _ := strings.CutSuffix(id, "ostacolo")
	slCoordinate := strings.Split(coordinate, ",")
	x0, _ = strconv.Atoi(slCoordinate[0])
	y0, _ = strconv.Atoi(slCoordinate[1])
	x1, _ = strconv.Atoi(slCoordinate[2])
	y1, _ = strconv.Atoi(slCoordinate[3])
	return
}

func calcolaDistanza(x0, y0, x1, y1 int) int {
	Distanza := math.Abs(float64(x1-x0)) + math.Abs(float64(y1-y0))
	return int(Distanza)
}

// Le nuove funzioni per l'algoritmo DFS
func cellaVisitata(listaVisitate *visitata, x, y int) bool {
	temp := listaVisitate
	for temp != nil {
		if temp.x == x && temp.y == y {
			return true
		}
		temp = temp.next
	}
	return false
}

func aggiungiVisitata(listaVisitate **visitata, x, y int) {
	nuova := new(visitata)
	nuova.x = x
	nuova.y = y
	nuova.next = *listaVisitate
	*listaVisitate = nuova
}

// Funzione principale DFS
func dfs(Campo piano, partenza, destinazione *punto) bool {
	// Definiamo i limiti del quadrato di ricerca
	minX := math.Min(float64(partenza.coordinataX), float64(destinazione.coordinataX))
	maxX := math.Max(float64(partenza.coordinataX), float64(destinazione.coordinataX))
	minY := math.Min(float64(partenza.coordinataY), float64(destinazione.coordinataY))
	maxY := math.Max(float64(partenza.coordinataY), float64(destinazione.coordinataY))

	// Direzioni di movimento (su, giù, sinistra, destra)
	dx := []int{0, 0, -1, 1}
	dy := []int{-1, 1, 0, 0}

	// Inizializzazione della pila per DFS
	var pila []*cella
	var listaVisitate *visitata = nil

	// Cella iniziale
	start := &cella{
		x:      partenza.coordinataX,
		y:      partenza.coordinataY,
		passi:  0,
		parent: nil,
	}

	pila = append(pila, start)
	aggiungiVisitata(&listaVisitate, start.x, start.y)

	for len(pila) > 0 {
		// Estrai l'ultimo elemento della pila
		current := pila[len(pila)-1]
		pila = pila[:len(pila)-1]
		// Se abbiamo raggiunto la destinazione
		if current.x == destinazione.coordinataX && current.y == destinazione.coordinataY {
			return true
		}
		// Esplora i vicini
		for i := 0; i < 4; i++ {
			newX := current.x + dx[i]
			newY := current.y + dy[i]
			// Verifica che la nuova posizione sia all'interno del quadrato di ricerca
			if float64(newX) < minX || float64(newX) > maxX || float64(newY) < minY || float64(newY) > maxY {
				continue
			}
			// Verifica che la nuova posizione non sia già stata visitata e non sia un ostacolo
			if !cellaVisitata(listaVisitate, newX, newY) && (*Campo).cercaOstacolo(newX, newY) == nil {
				nextCell := &cella{
					x:      newX,
					y:      newY,
					passi:  current.passi + 1,
					parent: current,
				}

				pila = append(pila, nextCell)
				aggiungiVisitata(&listaVisitate, newX, newY)
			}
		}
	}
	// Se non abbiamo trovato un percorso
	return false
}

func (Campo *Piano) cercaOstacolo(x int, y int) *punto {
	percorrente := Campo.ostacoli
	for percorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		if (x <= x1 && x >= x0) && (y <= y1 && y >= y0) {
			return percorrente
		}
		percorrente = percorrente.successivo
	}
	return nil
}

func (Campo *Piano) cercaAutoma(x, y int, id string) *punto {
	percorrente := Campo.automi
	for percorrente != nil {
		if percorrente.coordinataX == x && percorrente.coordinataY == y || percorrente.id == id {
			return percorrente
		}
		percorrente = percorrente.successivo
	}
	return nil
}
