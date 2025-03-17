package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Struttura per rappresentare un punto nel piano
type punto struct {
	coordinataX int
	coordinataY int
	id          string
	successivo  *punto
}

// Tipo piano è un puntatore a Piano
type piano *Piano

// Struttura che rappresenta il piano con automi e ostacoli
type Piano struct {
	automi   *punto
	ostacoli *punto
}

// Struttura per la coda BFS
type nodo struct {
	p        *punto
	distanza int
	next     *nodo
}

// Funzione esegui interpreta il comando e chiama la funzione corrispondente
func esegui(p piano, s string) {
	comandi := strings.Fields(s) // Usa Fields invece di Split per gestire spazi multipli
	if len(comandi) == 0 {
		return
	}

	switch comandi[0] {
	case "c":
		P := newPiano()
		p.automi = P.automi
		p.ostacoli = P.ostacoli
	case "S":
		(*p).stampa()
	case "s":
		if len(comandi) < 3 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).stato(a, b)
	case "a":
		if len(comandi) < 4 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).automa(a, b, comandi[3])
	case "o":
		if len(comandi) < 5 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		c, _ := strconv.Atoi(comandi[3])
		d, _ := strconv.Atoi(comandi[4])
		(*p).ostacolo(a, b, c, d)
	case "p":
		if len(comandi) < 2 {
			return
		}
		(*p).posizioni(comandi[1])
	case "r":
		if len(comandi) < 4 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).richiamo(a, b, comandi[3])
	case "e":
		if len(comandi) < 4 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).esistePercorso(a, b, comandi[3])
	case "t": // Operazione tortuosità (per l'appello di febbraio)
		if len(comandi) < 4 {
			return
		}
		a, _ := strconv.Atoi(comandi[1])
		b, _ := strconv.Atoi(comandi[2])
		(*p).tortuosita(a, b, comandi[3])
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
	fmt.Println("(")
	percorrente := Campo.automi
	for percorrente != nil {
		fmt.Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		percorrente = percorrente.successivo
	}
	fmt.Println(")")

	fmt.Println("[")
	percorrente = Campo.ostacoli
	for percorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		fmt.Printf("(%d,%d)(%d,%d)\n", x0, y0, x1, y1)
		percorrente = percorrente.successivo
	}
	fmt.Println("]")
}

func (Campo Piano) stato(x, y int) {
	if Campo.cercaOstacolo(x, y) != nil {
		fmt.Println("O")
		return
	}
	if Campo.cercaAutoma(x, y, "") != nil {
		fmt.Println("A")
	} else {
		fmt.Println("E")
	}
}

func (Campo Piano) posizioni(alpha string) {
	fmt.Println("(")
	percorrente := Campo.automi
	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			fmt.Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		}
		percorrente = percorrente.successivo
	}
	fmt.Println(")")
}

// Implementazione migliorata di esistePercorso utilizzando BFS
func (Campo *Piano) esistePercorso(x, y int, eta string) {
	// Verifica se la destinazione è un ostacolo
	if Campo.cercaOstacolo(x, y) != nil {
		fmt.Println("NO")
		return
	}

	// Cerca l'automa
	automa := Campo.cercaAutoma(x, y, eta)
	if automa == nil {
		fmt.Println("NO")
		return
	}

	// Calcola la distanza manhattan
	distanzaMinima := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

	// Verifica se esiste un percorso libero di distanza minima
	if Campo.trovaPercorsoMinimo(automa.coordinataX, automa.coordinataY, x, y, distanzaMinima) {
		fmt.Println("SI")
	} else {
		fmt.Println("NO")
	}
}

// Implementazione per l'operazione di tortuosità
func (Campo *Piano) tortuosita(x, y int, eta string) {
	// Verifica se la destinazione è un ostacolo
	if Campo.cercaOstacolo(x, y) != nil {
		fmt.Println("-1")
		return
	}

	// Cerca l'automa
	automa := Campo.cercaAutoma(-1, -1, eta)
	if automa == nil {
		fmt.Println("-1")
		return
	}

	// Calcola la distanza manhattan
	distanzaMinima := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

	// Trova la tortuosità minima
	tortuositaMin := Campo.trovaTortuositaMinima(automa.coordinataX, automa.coordinataY, x, y, distanzaMinima)
	fmt.Println(tortuositaMin)
}

func newPiano() piano {
	var nuovPiano Piano
	return &nuovPiano
}

func (Campo *Piano) automa(x, y int, eta string) {
	// Se c'è già un automa con questo nome, lo rimuovo dalla lista
	var prev *punto
	curr := Campo.automi

	for curr != nil {
		if curr.id == eta {
			if prev == nil {
				Campo.automi = curr.successivo
			} else {
				prev.successivo = curr.successivo
			}
			break
		}
		prev = curr
		curr = curr.successivo
	}

	// Se la posizione è all'interno di un ostacolo, non faccio niente
	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	// Creo un nuovo automa
	nuovoAutoma := &punto{
		coordinataX: x,
		coordinataY: y,
		id:          eta,
		successivo:  Campo.automi,
	}

	Campo.automi = nuovoAutoma
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	// Verifica se ci sono automi nel rettangolo
	percorrente := Campo.automi
	for percorrente != nil {
		if percorrente.coordinataX >= x0 && percorrente.coordinataX <= x1 &&
			percorrente.coordinataY >= y0 && percorrente.coordinataY <= y1 {
			return
		}
		percorrente = percorrente.successivo
	}

	// Crea un nuovo ostacolo
	newOstacolo := &punto{
		coordinataX: x0,
		coordinataY: y1,
		id:          fmt.Sprintf("%d,%d,%d,%d,ostacolo", x0, y0, x1, y1),
		successivo:  Campo.ostacoli,
	}

	Campo.ostacoli = newOstacolo
}

func (Campo *Piano) richiamo(x, y int, alpha string) {
	// Se la posizione è all'interno di un ostacolo, nessun automa può raggiungerla
	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	// Trova tutti gli automi che devono rispondere al richiamo
	var automiChiamati []*punto
	percorrente := Campo.automi

	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			automiChiamati = append(automiChiamati, percorrente)
		}
		percorrente = percorrente.successivo
	}

	// Per ogni automa, calcola la distanza minima e verifica se può raggiungerla
	distanzaMinima := math.MaxInt32
	var automiDaSpostare []*punto

	for _, automa := range automiChiamati {
		distanza := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

		// Verifica se esiste un percorso libero di lunghezza minima
		if Campo.trovaPercorsoMinimo(automa.coordinataX, automa.coordinataY, x, y, distanza) {
			if distanza < distanzaMinima {
				distanzaMinima = distanza
				automiDaSpostare = []*punto{automa}
			} else if distanza == distanzaMinima {
				automiDaSpostare = append(automiDaSpostare, automa)
			}
		}
	}

	// Sposta gli automi che hanno la distanza minima
	for _, automa := range automiDaSpostare {
		automa.coordinataX = x
		automa.coordinataY = y
	}
}

// Funzione che trova se esiste un percorso libero di lunghezza minima utilizzando BFS
func (Campo *Piano) trovaPercorsoMinimo(x0, y0, x1, y1, distanzaMinima int) bool {
	// Se la distanza è 0, siamo già arrivati
	if x0 == x1 && y0 == y1 {
		return true
	}

	// Crea una mappa di visite
	visite := make(map[string]bool)

	// Crea la coda per BFS
	var coda nodo
	coda.p = &punto{coordinataX: x0, coordinataY: y0}
	coda.distanza = 0

	// Direzioni possibili: destra, sinistra, su, giù
	dx := []int{1, -1, 0, 0}
	dy := []int{0, 0, 1, -1}

	// BFS
	nodoCorrente := &coda
	for nodoCorrente != nil {
		x := nodoCorrente.p.coordinataX
		y := nodoCorrente.p.coordinataY
		dist := nodoCorrente.distanza

		// Se siamo arrivati alla destinazione e la distanza è quella minima
		if x == x1 && y == y1 && dist == distanzaMinima {
			return true
		}

		// Se la distanza è già uguale alla distanza minima, non possiamo andare oltre
		if dist >= distanzaMinima {
			nodoCorrente = nodoCorrente.next
			continue
		}

		// Marca come visitato
		chiave := fmt.Sprintf("%d,%d", x, y)
		visite[chiave] = true

		// Prova tutte le direzioni
		for i := 0; i < 4; i++ {
			nx := x + dx[i]
			ny := y + dy[i]
			chiaveNuova := fmt.Sprintf("%d,%d", nx, ny)

			// Se non è già visitato e non è un ostacolo
			if !visite[chiaveNuova] && Campo.cercaOstacolo(nx, ny) == nil {
				// Calcola la nuova distanza totale
				nuovaDistanzaTotale := dist + 1 + calcolaDistanza(nx, ny, x1, y1)

				// Se la nuova distanza totale è minore o uguale alla distanza minima, aggiungi alla coda
				if nuovaDistanzaTotale <= distanzaMinima {
					nuovoNodo := &nodo{
						p:        &punto{coordinataX: nx, coordinataY: ny},
						distanza: dist + 1,
						next:     nodoCorrente.next,
					}
					nodoCorrente.next = nuovoNodo
				}
			}
		}

		nodoCorrente = nodoCorrente.next
	}

	return false
}

// Funzione che trova la tortuosità minima di un percorso libero
func (Campo *Piano) trovaTortuositaMinima(x0, y0, x1, y1, distanzaMinima int) int {
	// Se non esiste un percorso di distanza minima, restituisci -1
	if !Campo.trovaPercorsoMinimo(x0, y0, x1, y1, distanzaMinima) {
		return -1
	}

	// Crea una matrice di visite con la tortuosità minima
	type statoVisita struct {
		tortuosita int
		direzione  int // 0: nessuna, 1: orizzontale, 2: verticale
	}

	visite := make(map[string]statoVisita)

	// Crea la coda per BFS
	type nodoTortuosita struct {
		x, y       int
		distanza   int
		tortuosita int
		direzione  int // 0: nessuna, 1: orizzontale, 2: verticale
		next       *nodoTortuosita
	}

	var coda nodoTortuosita
	coda.x = x0
	coda.y = y0
	coda.distanza = 0
	coda.tortuosita = 0
	coda.direzione = 0

	// Direzioni possibili: destra, sinistra, su, giù
	dx := []int{1, -1, 0, 0}
	dy := []int{0, 0, 1, -1}
	dir := []int{1, 1, 2, 2} // 1: orizzontale, 2: verticale

	minTortuosita := math.MaxInt32

	// BFS
	nodoCorrente := &coda
	for nodoCorrente != nil {
		x := nodoCorrente.x
		y := nodoCorrente.y
		dist := nodoCorrente.distanza
		tort := nodoCorrente.tortuosita
		dirCorrente := nodoCorrente.direzione

		// Se siamo arrivati alla destinazione e la distanza è quella minima
		if x == x1 && y == y1 && dist == distanzaMinima {
			if tort < minTortuosita {
				minTortuosita = tort
			}
		}

		// Se la distanza è già uguale alla distanza minima, non possiamo andare oltre
		if dist >= distanzaMinima {
			nodoCorrente = nodoCorrente.next
			continue
		}

		// Marca come visitato
		chiave := fmt.Sprintf("%d,%d,%d", x, y, dirCorrente)
		stato, trovato := visite[chiave]
		if trovato && stato.tortuosita <= tort {
			nodoCorrente = nodoCorrente.next
			continue
		}
		visite[chiave] = statoVisita{tortuosita: tort, direzione: dirCorrente}

		// Prova tutte le direzioni
		for i := 0; i < 4; i++ {
			nx := x + dx[i]
			ny := y + dy[i]

			// Se non è un ostacolo
			if Campo.cercaOstacolo(nx, ny) == nil {
				// Calcola la nuova distanza totale
				nuovaDistanza := dist + 1

				// Calcola la nuova tortuosità
				nuovaTortuosita := tort
				if dirCorrente != 0 && dirCorrente != dir[i] {
					nuovaTortuosita++
				}

				// Aggiungi alla coda
				nuovoNodo := &nodoTortuosita{
					x:          nx,
					y:          ny,
					distanza:   nuovaDistanza,
					tortuosita: nuovaTortuosita,
					direzione:  dir[i],
					next:       nodoCorrente.next,
				}
				nodoCorrente.next = nuovoNodo
			}
		}

		nodoCorrente = nodoCorrente.next
	}

	if minTortuosita == math.MaxInt32 {
		return -1
	}
	return minTortuosita
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
	return int(math.Abs(float64(x1-x0)) + math.Abs(float64(y1-y0)))
}

func (Campo *Piano) cercaOstacolo(x int, y int) *punto {
	percorrente := Campo.ostacoli
	for percorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
			return percorrente
		}
		percorrente = percorrente.successivo
	}
	return nil
}

func (Campo *Piano) cercaAutoma(x, y int, id string) *punto {
	percorrente := Campo.automi
	for percorrente != nil {
		if (percorrente.coordinataX == x && percorrente.coordinataY == y) ||
			(id != "" && percorrente.id == id) {
			return percorrente
		}
		percorrente = percorrente.successivo
	}
	return nil
}
