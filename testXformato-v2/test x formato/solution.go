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
/*
Possibile modifica per la struttura punto:
	inannzitutto la struttua può avere un puntatora all'elemento precendente
	inoltre all'interno dell'id si può inserire l'indice di posizione nella lista
	In questo modo si può operativamente implementare una ricerca dicotomica all'interno della lista
	Viso che lo scorrimento della lista degli ostacoli e di quella degli autommi è l'operazione che si svolge più di
	frequente, mi sembra opportuno renderela di costo logaritmico invece che lieare.

	Questo imlica la creazione di due nuove strutture dati per automi e ostacoli, esse terranno un puntato alla coda e alla testa
	delle rispettive liste.
*/
type punto struct {
	coordinataX int
	coordinataY int
	id          string
	indice      int
	successivo  *punto
	precendente *punto
}

type gruppo struct {
	inizio *punto
	fine   *punto
}

// Tipo piano è un puntatore a Piano
type piano *Piano

// Struttura che rappresenta il piano con automi e ostacoli
type Piano struct {
	automi   *gruppo
	ostacoli *gruppo
}

// Struttura per una pila
type elPIla struct {
	p        *punto
	distanza int
	next     *elPIla
}

// Funzione esegui interpreta il comando e chiama la funzione corrispondente
func esegui(p piano, s string) {
	comandi := strings.Fields(s)
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
	case "t":
		fmt.Println("0")
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
	percorrente := Campo.automi.inizio
	for percorrente != nil {
		fmt.Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		percorrente = percorrente.successivo
	}
	fmt.Println(")")

	fmt.Println("[")
	percorrente = Campo.ostacoli.inizio
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
	percorrente := Campo.automi.inizio
	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			fmt.Printf("%s: %d,%d\n", percorrente.id, percorrente.coordinataX, percorrente.coordinataY)
		}
		percorrente = percorrente.successivo
	}
	fmt.Println(")")
}

func (Campo *Piano) esistePercorso(x, y int, eta string) {
	if Campo.cercaOstacolo(x, y) != nil {
		fmt.Println("NO")
		return
	}

	automa := Campo.cercaAutoma(x, y, eta)
	if automa == nil {
		fmt.Println("NO")
		return
	}

	distanzaMinima := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

	if Campo.trovaPercorsoMinimo(automa.coordinataX, automa.coordinataY, x, y, distanzaMinima) {
		fmt.Println("SI")
	} else {
		fmt.Println("NO")
	}
}

func newPiano() piano {
	var nuovPiano Piano
	nuovPiano.automi = &gruppo{
		inizio: nil,
		fine:   nil,
	}
	nuovPiano.ostacoli = &gruppo{
		inizio: nil,
		fine:   nil,
	}
	return &nuovPiano
}

func (Campo *Piano) automa(x, y int, eta string) {
	percorrente := Campo.automi.inizio
	ripercorrente := Campo.automi.fine

	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	if Campo.cercaAutoma(x, y, "") != nil {
		return
	}

	// Algoritmo per una ricerca dicotomica
	for percorrente != nil && ripercorrente != nil {
		if percorrente.id == eta {
			percorrente.coordinataX = x
			percorrente.coordinataY = y
			return
		}

		if ripercorrente.id == eta {
			ripercorrente.coordinataX = x
			ripercorrente.coordinataY = y
			return
		}

		if ripercorrente.indice < percorrente.indice {
			break
		}

		percorrente = percorrente.successivo
		ripercorrente = ripercorrente.precendente
	}

	nuovoAutoma := &punto{
		coordinataX: x,
		coordinataY: y,
		id:          eta,
	}

	if Campo.automi.inizio != nil && Campo.automi.fine == nil {
		nuovoAutoma.indice = Campo.automi.inizio.indice + 1
		nuovoAutoma.precendente = Campo.automi.inizio
		Campo.automi.fine = nuovoAutoma
		Campo.automi.inizio.successivo = nuovoAutoma
		return
	} else if Campo.automi.inizio == nil && Campo.ostacoli.fine == nil {
		nuovoAutoma.indice = 1
		nuovoAutoma.precendente = nil
		Campo.automi.inizio = nuovoAutoma
		return
	} else {
		nuovoAutoma.indice = Campo.automi.fine.indice + 1
		nuovoAutoma.precendente = Campo.automi.fine
		Campo.automi.fine = nuovoAutoma
		return
	}
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	percorrente := Campo.ostacoli.inizio
	ripercorrente := Campo.ostacoli.fine

	// Ricerca dicotomica
	for percorrente != nil && ripercorrente != nil {
		if percorrente.coordinataX >= x0 && percorrente.coordinataX <= x1 &&
			percorrente.coordinataY >= y0 && percorrente.coordinataY <= y1 {
			return
		}

		if ripercorrente.coordinataX >= x0 && ripercorrente.coordinataX <= x1 &&
			ripercorrente.coordinataY >= y0 && ripercorrente.coordinataY <= y1 {
			return
		}
		if ripercorrente.indice < percorrente.indice {
			break
		}
		ripercorrente = ripercorrente.precendente
		percorrente = percorrente.successivo
	}

	newOstacolo := &punto{
		coordinataX: x0,
		coordinataY: y1,
		id:          fmt.Sprintf("%d,%d,%d,%d,ostacolo", x0, y0, x1, y1),
	}

	if Campo.ostacoli.inizio != nil && Campo.ostacoli.fine == nil {
		newOstacolo.indice = Campo.ostacoli.inizio.indice + 1
		newOstacolo.precendente = Campo.ostacoli.inizio
		Campo.ostacoli.fine = newOstacolo
		Campo.ostacoli.inizio.successivo = newOstacolo
		return
	} else if Campo.ostacoli.inizio == nil {
		newOstacolo.indice = 1
		Campo.ostacoli.inizio = newOstacolo
		return
	} else {
		newOstacolo.indice = Campo.ostacoli.fine.indice + 1
		newOstacolo.precendente = Campo.ostacoli.fine
		Campo.ostacoli.fine = newOstacolo
		return
	}
}

func (Campo *Piano) richiamo(x, y int, alpha string) {
	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	var automiChiamati []*punto
	percorrente := Campo.automi.inizio

	for percorrente != nil {
		if strings.HasPrefix(percorrente.id, alpha) {
			automiChiamati = append(automiChiamati, percorrente)
		}
		percorrente = percorrente.successivo
	}

	distanzaMinima := math.MaxInt32
	var automiDaSpostare []*punto

	for _, automa := range automiChiamati {
		distanza := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

		if Campo.trovaPercorsoMinimo(automa.coordinataX, automa.coordinataY, x, y, distanza) {
			if distanza < distanzaMinima {
				distanzaMinima = distanza
				automiDaSpostare = []*punto{automa}
			} else if distanza == distanzaMinima {
				automiDaSpostare = append(automiDaSpostare, automa)
			}
		}
	}

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
	var coda elPIla
	coda.p = &punto{coordinataX: x0, coordinataY: y0}
	coda.distanza = 0

	// BFS
	elCorrente := &coda
	for elCorrente != nil {
		x := elCorrente.p.coordinataX
		y := elCorrente.p.coordinataY
		dist := elCorrente.distanza

		// Se siamo arrivati alla destinazione e la distanza è quella minima
		if x == x1 && y == y1 && dist == distanzaMinima {
			return true
		}

		// Se la distanza è già uguale alla distanza minima, non possiamo andare oltre
		if dist >= distanzaMinima {
			elCorrente = elCorrente.next
			continue
		}

		// Marca come visitato
		chiave := fmt.Sprintf("%d,%d", x, y)
		visite[chiave] = true

		direzionePossibile := [][2]int{}
		if x1 > x {
			direzionePossibile = append(direzionePossibile, [2]int{1, 0})
		} else if x1 < x {
			direzionePossibile = append(direzionePossibile, [2]int{-1, 0})
		}
		if y1 > y {
			direzionePossibile = append(direzionePossibile, [2]int{0, 1})
		} else if y1 < y {
			direzionePossibile = append(direzionePossibile, [2]int{0, -1})
		}
		for _, d := range direzionePossibile {
			nx := x + d[0]
			ny := y + d[1]
			chiaveNuova := fmt.Sprintf("%d,%d", nx, ny)

			// Se non è già visitato e non è un ostacolo
			if !visite[chiaveNuova] && Campo.cercaOstacolo(nx, ny) == nil {
				// Calcola la nuova distanza totale
				nuovaDistanzaTotale := dist + 1 + calcolaDistanza(nx, ny, x1, y1)

				// Se la nuova distanza totale è minore o uguale alla distanza minima, aggiungi alla coda
				if nuovaDistanzaTotale <= distanzaMinima {
					nuovoElemento := &elPIla{
						p:        &punto{coordinataX: nx, coordinataY: ny},
						distanza: dist + 1,
						next:     elCorrente.next,
					}
					elCorrente.next = nuovoElemento
				}
			}
		}
		elCorrente = elCorrente.next
	}
	return false
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
	percorrente := Campo.ostacoli.inizio
	ripercorrente := Campo.ostacoli.fine

	for percorrente != nil && ripercorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
			return percorrente
		}

		x0, y0, x1, y1 = estraiCoordinate(ripercorrente.id)
		if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
			return ripercorrente
		}

		if ripercorrente.indice < percorrente.indice {
			break
		}
		ripercorrente = ripercorrente.precendente
		percorrente = percorrente.successivo
	}

	if percorrente != nil {
		x0, y0, x1, y1 := estraiCoordinate(percorrente.id)
		if x >= x0 && x <= x1 && y >= y0 && y <= y1 {
			return percorrente
		}
	}
	return nil
}

func (Campo *Piano) cercaAutoma(x, y int, id string) *punto {
	percorrente := Campo.automi.inizio
	ripercorrente := Campo.automi.fine

	for percorrente != nil && ripercorrente != nil {
		if (percorrente.coordinataX == x && percorrente.coordinataY == y) ||
			(id != "" && percorrente.id == id) {
			return percorrente
		}

		if (ripercorrente.coordinataX == x && ripercorrente.coordinataY == y) ||
			(id != "" && ripercorrente.id == id) {
			return ripercorrente
		}

		if ripercorrente.indice < percorrente.indice {
			break
		}

		percorrente = percorrente.successivo
		ripercorrente = ripercorrente.precendente
	}

	if percorrente != nil && ((percorrente.coordinataX == x && percorrente.coordinataY == y) || percorrente.id == id) {
		return percorrente
	}

	return nil
}
