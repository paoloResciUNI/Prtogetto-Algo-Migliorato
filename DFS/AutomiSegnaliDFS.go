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
}

type ostacolo struct {
	Ss        punto // punto in basso a sinistra
	Nd        punto // punto in alto a destra
	diagonale int
}

type piano *Piano

type Piano struct {
	automi   map[string]*punto // scelta presa per rendere più facile prendere un automa (tempo costante)
	ostacoli *[]ostacolo       // contiene tutti gli ostacoli che sono riordinati in ordine di grandezza
	// proviamo a mette gli ostacoli in una coda con priorità, dove la priorità è la diagonale del rettangolo
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
	Println("(")
	for k, v := range Campo.automi {
		Printf("%s: %d,%d\n", k, v.coordinataX, v.coordinataY)
	}
	Println(")")
	Println("[")
	for i := 0; i < len((*Campo.ostacoli)); i++ {
		Printf("(%d,%d)(%d,%d)\n", (*Campo.ostacoli)[i].Ss.coordinataX, (*Campo.ostacoli)[i].Ss.coordinataY, (*Campo.ostacoli)[i].Nd.coordinataX, (*Campo.ostacoli)[i].Nd.coordinataY)
	}
	Println("]")
}

func (Campo Piano) stato(x, y int) {
	if Campo.cercaOstacolo(x, y) != nil {
		Println("O")
		return
	}
	for _, v := range Campo.automi {
		if v.coordinataX == x && v.coordinataY == y {
			Println("A")
			return
		}
	}
	Println("E")
}

func (Campo Piano) posizioni(alpha string) {
	Println("(")
	for k, v := range Campo.automi {
		if strings.HasPrefix(k, alpha) {
			Printf("%s: %d,%d\n", k, v.coordinataX, v.coordinataY)
		}
	}
	Println(")")
}

// Implementazione modificata con DFS
func (Campo *Piano) esistePercorso(x, y int, eta string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	if Campo.cercaOstacolo(x, y) != nil {
		Println("NO")
		return
	}
	percorrente := Campo.automi[eta]
	if percorrente == nil {
		Println("NO")
		return
	}
	percorsoEffettuato := (avanza(Campo, percorrente, Sorgente))
	if percorsoEffettuato {
		Println("SI")
		return
	} else {
		Println("NO")
	}
}

func newPiano() piano {
	NewPunto := new(Piano)
	NewPunto.automi = make(map[string]*punto)
	NewPunto.ostacoli = new([]ostacolo)
	return NewPunto
}

func (Campo *Piano) automa(x, y int, eta string) {
	puntoCercato := new(punto)
	if Campo.automi[eta] == nil {
		for _, v := range Campo.automi {
			if v.coordinataX == x && v.coordinataY == y {
				return
			}
		}
	}
	if Campo.cercaOstacolo(x, y) == nil {
		puntoCercato = new(punto)
		puntoCercato.coordinataX = x
		puntoCercato.coordinataY = y
		Campo.automi[eta] = puntoCercato
		return
	}
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	for _, v := range Campo.automi {
		if (v.coordinataX <= x1 && v.coordinataX >= x0) && (v.coordinataY <= y1 && v.coordinataY >= y0) {
			return
		}
	}
	newOstacolo := new(ostacolo)
	newOstacolo.Ss = punto{
		coordinataX: x0,
		coordinataY: y0,
	}
	newOstacolo.Nd = punto{
		coordinataX: x1,
		coordinataY: y1,
	}
	newOstacolo.diagonale = calcolaDistanza(x0, y0, x1, y1)

	*Campo.ostacoli = append(*Campo.ostacoli, *newOstacolo)

	return
}

// Implementazione modificata del richiamo con DFS
func (Campo *Piano) richiamo(x, y int, alpha string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	minDistance := math.MaxInt
	pilaChiamata := new(elementoPila)
	for k, v := range Campo.automi {
		if strings.HasPrefix(k, alpha) {
			distanza := calcolaDistanza(v.coordinataX, v.coordinataY, x, y)
			possibileAvanzamento := avanza(Campo, v, Sorgente)
			if possibileAvanzamento {
				if distanza <= minDistance {
					minDistance = distanza
				}
				automaChiamato := new(elementoPila)
				automaChiamato.chiamato = v
				automaChiamato.prossimo = pilaChiamata
				automaChiamato.distanza = distanza
				pilaChiamata = automaChiamato
			}
		}
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
func avanza(Campo piano, partenza, destinazione *punto) bool {
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

func (Campo *Piano) cercaOstacolo(x int, y int) *ostacolo {
	testa, coda := 0, len(*Campo.ostacoli)-1
	for testa <= coda {
		if (x <= (*Campo.ostacoli)[testa].Nd.coordinataX && x >= (*Campo.ostacoli)[testa].Ss.coordinataX) && (y <= (*Campo.ostacoli)[testa].Nd.coordinataY && y >= (*Campo.ostacoli)[testa].Ss.coordinataY) {
			return &(*Campo.ostacoli)[testa]
		}
		if (x <= (*Campo.ostacoli)[coda].Nd.coordinataX && x >= (*Campo.ostacoli)[coda].Ss.coordinataX) && (y <= (*Campo.ostacoli)[coda].Nd.coordinataY && y >= (*Campo.ostacoli)[coda].Ss.coordinataY) {
			return &(*Campo.ostacoli)[coda]
		}
		testa++
		coda--
	}
	return nil
}
