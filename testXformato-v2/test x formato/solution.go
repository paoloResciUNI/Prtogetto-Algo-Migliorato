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

// Struttura per una pila
type elPIla struct {
	p        *punto
	distanza int
	next     *elPIla
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
	case "t":
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

func (Campo *Piano) tortuosita(x, y int, eta string) {
	if Campo.cercaOstacolo(x, y) != nil {
		fmt.Println("-1")
		return
	}

	automa := Campo.cercaAutoma(-1, -1, eta)
	if automa == nil {
		fmt.Println("-1")
		return
	}

	distanzaMinima := calcolaDistanza(automa.coordinataX, automa.coordinataY, x, y)

	tortuositaMin := Campo.trovaTortuositaMinima(automa.coordinataX, automa.coordinataY, x, y, distanzaMinima)
	fmt.Println(tortuositaMin)
}

func newPiano() piano {
	var nuovPiano Piano
	return &nuovPiano
}

func (Campo *Piano) automa(x, y int, eta string) {
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

	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	nuovoAutoma := &punto{
		coordinataX: x,
		coordinataY: y,
		id:          eta,
		successivo:  Campo.automi,
	}

	Campo.automi = nuovoAutoma
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	percorrente := Campo.automi
	for percorrente != nil {
		if percorrente.coordinataX >= x0 && percorrente.coordinataX <= x1 &&
			percorrente.coordinataY >= y0 && percorrente.coordinataY <= y1 {
			return
		}
		percorrente = percorrente.successivo
	}

	newOstacolo := &punto{
		coordinataX: x0,
		coordinataY: y1,
		id:          fmt.Sprintf("%d,%d,%d,%d,ostacolo", x0, y0, x1, y1),
		successivo:  Campo.ostacoli,
	}

	Campo.ostacoli = newOstacolo
}

func (Campo *Piano) richiamo(x, y int, alpha string) {
	if Campo.cercaOstacolo(x, y) != nil {
		return
	}

	var automiChiamati []*punto
	percorrente := Campo.automi

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

func (Campo *Piano) trovaTortuositaMinima(x0, y0, x1, y1, distanzaMinima int) int {
	// Verifica prima se esiste un percorso valido
	if !Campo.trovaPercorsoMinimo(x0, y0, x1, y1, distanzaMinima) {
		return -1
	}

	// Definizione della struttura per la BFS
	type Nodo struct {
		x, y      int
		passi     int
		curve     int
		direzione int // 0=nessuna, 1=orizzontale, 2=verticale
		next      *Nodo
	}

	// Struttura per memorizzare lo stato di una visita
	type StatoVisita struct {
		curve     int
		direzione int
	}

	// Mappa delle visite (chiave: "x,y,direzione")
	visite := make(map[string]StatoVisita)

	// Direzioni possibili: destra, sinistra, su, giù
	dx := []int{1, -1, 0, 0}
	dy := []int{0, 0, 1, -1}
	// Tipo di direzione: 1=orizzontale, 2=verticale
	tipiDirezione := []int{1, 1, 2, 2}

	// Inizializza la coda con il punto di partenza
	coda := Nodo{
		x:         x0,
		y:         y0,
		passi:     0,
		curve:     0,
		direzione: 0,
	}

	tortuositaMinima := math.MaxInt32

	// BFS che tiene traccia della tortuosità
	nodoCorrente := &coda
	for nodoCorrente != nil {
		x := nodoCorrente.x
		y := nodoCorrente.y
		passi := nodoCorrente.passi
		curve := nodoCorrente.curve
		dirAttuale := nodoCorrente.direzione

		// Destinazione raggiunta con distanza minima?
		if x == x1 && y == y1 && passi == distanzaMinima {
			if curve < tortuositaMinima {
				tortuositaMinima = curve
			}
		}

		// Se abbiamo superato la distanza minima, passa al prossimo nodo
		if passi >= distanzaMinima {
			nodoCorrente = nodoCorrente.next
			continue
		}

		// Controlla se abbiamo già visitato questo stato con meno curve
		chiave := fmt.Sprintf("%d,%d,%d", x, y, dirAttuale)
		if stato, trovato := visite[chiave]; trovato && stato.curve <= curve {
			nodoCorrente = nodoCorrente.next
			continue
		}

		// Aggiorna lo stato della visita
		visite[chiave] = StatoVisita{curve: curve, direzione: dirAttuale}

		// Esplora tutte le direzioni possibili
		for i := 0; i < 4; i++ {
			nx := x + dx[i]
			ny := y + dy[i]

			// Se la posizione è valida (non è un ostacolo)
			if Campo.cercaOstacolo(nx, ny) == nil {
				nuoveCurve := curve

				// Incrementa le curve se cambiamo direzione (e non è la prima mossa)
				if dirAttuale != 0 && dirAttuale != tipiDirezione[i] {
					nuoveCurve++
				}

				// Aggiungi il nuovo nodo alla coda
				nuovoNodo := &Nodo{
					x:         nx,
					y:         ny,
					passi:     passi + 1,
					curve:     nuoveCurve,
					direzione: tipiDirezione[i],
					next:      nodoCorrente.next,
				}
				nodoCorrente.next = nuovoNodo
			}
		}
		nodoCorrente = nodoCorrente.next
	}

	if tortuositaMinima == math.MaxInt32 {
		return -1
	}
	return tortuositaMinima
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
