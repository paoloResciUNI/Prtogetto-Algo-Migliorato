package main

import (
	"bufio"
	. "fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sort"
)

type punto struct {
	coordinataX int
	coordinataY int
}

type ostacolo struct { 
	Ss	punto	// punto in basso a sinistra
	Nd	punto	// punto in alto a destra
	diagonale	int
}


type piano *Piano

type Piano struct {
	automi  map[string]*punto // scelta presa per rendere più facile prendere un automa (tempo costante)
	ostacoli *[]ostacolo // contiene tutti gli ostacoli che sono riordinati in ordine di grandezza 
	// proviamo a mette gli ostacoli in una coda con priorità, dove la priorità è la diagonale del rettangolo 
}

type elementoPila struct {
	chiamato *punto
	distanza int
	prossimo *elementoPila
}

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
	for i := 0; i < len((*Campo.automi)); i++ {
		Printf("%s: %d,%d\n", (*Campo.automi)[i].id, (*Campo.automi)[i].coordinataX, (*Campo.automi)[i].coordinataY)
	}
	Println(")")
	Println("[")
	for i := 0; i < len((*Campo.ostacoli)); i++ {
		x0, y0, x1, y1 := estraiCoordinate((*Campo.ostacoli)[i].id)
		Printf("(%d,%d)(%d,%d)\n", x0, y0, x1, y1)
	}
	Println("]")
}

func (Campo Piano) stato(x, y int) string {
	if Campo.cercaOstacolo(x, y) != nil {
		return "O"
	}
	if Campo.cercaAutoma(x, y, "") != nil {
		return "A"
	} else {
		return "E"
	}
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

func (Campo *Piano) esistePercorso(x, y int, eta string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	if Campo.cercaOstacolo(x, y) != nil {
		Println("NO")
		return
	}
	percorrente := Campo.cercaAutoma(x, y, eta)
	if percorrente == nil {
		Println("NO")
		return
	}
	percorsoEffettuato := (avanza(Campo, percorrente, Sorgente))
	if percorsoEffettuato.coordinataX == x && percorsoEffettuato.coordinataY == y {
		Println("SI")
		return
	} else {
		Println("NO")
	}
}

func newPiano() piano {
	NewPunto := new(Piano)
	NewPunto.automi = *new(map[string]*punto)
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


func (Campo *Piano) richiamo(x, y int, alpha string) {
	Sorgente := new(punto)
	Sorgente.coordinataX = x
	Sorgente.coordinataY = y
	minDistance := math.MaxInt
	pilaChiamata := new(elementoPila)
	for i := 0; i < len((*Campo.automi)); i++ {
		if strings.HasPrefix((*Campo.automi)[i].id, alpha) {
			distanza := calcolaDistanza((*Campo.automi)[i].coordinataX, (*Campo.automi)[i].coordinataY, x, y)
			possibileAvanzamento := avanza(Campo, &(*Campo.automi)[i], Sorgente)
			if possibileAvanzamento.coordinataX == x && possibileAvanzamento.coordinataY == y {
				if distanza <= minDistance {
					minDistance = distanza
				}
				automaChiamato := new(elementoPila)
				automaChiamato.chiamato = &(*Campo.automi)[i]
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

// Rimuovere i confronti 
func (Campo *Piano) ostacoliPercorso(partenza, arrivo *punto) (distanza_O_Ascisse, distanza_O_Ordinate int) {
	ostacoloVicino := partenza.posizioneOstacoloVerticale(Campo, arrivo.coordinataY)
	if ostacoloVicino != nil {
		if arrivo.coordinataY < partenza.coordinataY {
			distanza_O_Ordinate = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, partenza.coordinataX, y1)
		} else if arrivo.coordinataY > partenza.coordinataY {
			distanza_O_Ordinate = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, partenza.coordinataX, y0)
		}
	} else {
		distanza_O_Ordinate = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, partenza.coordinataX, arrivo.coordinataY)
	}
	ostacoloVicino = partenza.posizioneOstacoloOrizzontale(Campo, arrivo.coordinataX)
	if ostacoloVicino != nil {
		if arrivo.coordinataX < partenza.coordinataX {
			distanza_O_Ascisse = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, x1, partenza.coordinataY)
		} else {
			distanza_O_Ascisse = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, x0, partenza.coordinataY)
		}
	} else {
		distanza_O_Ascisse = calcolaDistanza(partenza.coordinataX, partenza.coordinataY, arrivo.coordinataX, partenza.coordinataY)
	}
	return
}

// funzione che restituisce tutti gli ostacoli che occupano un determinato asse, orizzontale o verticale  
func ostacoliAsse(Campo piano, asse int, ascisse bool) []ostacoli {
	filtrati := new([]ostacolo)
	for i := 0; i < len(*Campo.ostacoli); i++ {
		if ascisse {
			if *Campo.ostacoli[i].Ss.coordinataX >= asse && *Campo.ostacoli[i].Nd.coordinataX <= asse {
				filtrati = append(filtrati, *Campo.ostacoli[i])
			}  
		} else {
			if *Campo.ostacoli[i].Ss.coordinataY >= asse && *Campo.ostacoli[i].Nd.coordinataY <= asse {
				filtrati = append(filtrati, *Campo.ostacoli[i])
			}  
		}
	}
	return filtrati
}

func (p *punto) posizioneOstacoloVerticale(Campo piano, y int) *punto {
	ostacoliVerticali := ostacoliAsse(Campo, p, true)
	if p.coordinataY > y {
		for i := p.coordinataY - 1; i >= y; i-- {
			ostacolo := (*Campo).cercaOstacolo(p.coordinataX, i)
			if ostacolo != nil {
				return ostacolo
			}
		}
	} else if p.coordinataY < y {
		for i := p.coordinataY + 1; i < y; i++ {
			ostacolo := (*Campo).cercaOstacolo(p.coordinataX, i)
			if ostacolo != nil {
				return ostacolo
			}
		}
	}
	return nil
}

func (p *punto) posizioneOstacoloOrizzontale(Campo piano, x int) *punto {
	if p.coordinataX > x {
		for i := p.coordinataX - 1; i >= x; i-- {
			ostacolo := (*Campo).cercaOstacolo(i, p.coordinataY)
			if ostacolo != nil {
				return ostacolo
			}
		}
	} else if p.coordinataX < x {
		for i := p.coordinataX + 1; i < x; i++ {
			ostacolo := (*Campo).cercaOstacolo(i, p.coordinataY)
			if ostacolo != nil {
				return ostacolo
			}
		}
	}
	return nil
}

func calcolaDistanza(x0, y0, x1, y1 int) int {
	Distanza := math.Abs(float64(x1-x0)) + math.Abs(float64(y1-y0))
	return int(Distanza)
}

func avanza(Campo piano, p *punto, Sorgente *punto) *punto {
	var distanzaVerticale, distanzaOrizzontale int
	passi := calcolaDistanza(p.coordinataX, p.coordinataY, Sorgente.coordinataX, Sorgente.coordinataY)
	if passi <= 0 || p.coordinataX == Sorgente.coordinataX && p.coordinataY == Sorgente.coordinataY {
		return p
	}
	possibilePasso := new(punto)
	possibilePasso.coordinataX = p.coordinataX
	possibilePasso.coordinataY = p.coordinataY
	possibilePasso.id = p.id
	distanzaOrizzontale, distanzaVerticale = (*Campo).ostacoliPercorso(possibilePasso, Sorgente)
	if distanzaVerticale < distanzaOrizzontale {
		possibilePasso = (*Campo).forwardX(possibilePasso, Sorgente)
	} else if distanzaOrizzontale == distanzaVerticale && p.coordinataX < Sorgente.coordinataX && (*Campo).cercaOstacolo(p.coordinataX+1, p.coordinataY) == nil {
		possibilePasso.coordinataX++
	} else if distanzaOrizzontale == distanzaVerticale && p.coordinataX > Sorgente.coordinataX && (*Campo).cercaOstacolo(p.coordinataX-1, p.coordinataY) == nil {
		possibilePasso.coordinataX--
	} else if distanzaVerticale > distanzaOrizzontale {
		possibilePasso = (*Campo).forwardY(possibilePasso, Sorgente)
	}
	if p.coordinataX == possibilePasso.coordinataX && p.coordinataY == possibilePasso.coordinataY {
		return possibilePasso
	}
	return avanza(Campo, possibilePasso, Sorgente)
}

func (Campo *Piano) forward(start *punto, destination *punto, forwardX bool) *punto {
	var forward punto
	forward.coordinataY = start.coordinataY // la coordinata y resta sempre la stessa perchè ci muoviamo solo sulle x

	// Vine controllato il primo ostacolo che il punto incontrerà nel suo percorso, se ne incontrerà
	if forwardX {
		ostacoloVicino := start.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
	} else {
		ostacoloVicino := start.posizioneOstacoloVerticale(Campo, destination.coordinataY)
	}
	if ostacoloVicino != nil {
		x0, y0, x1, y1 := estraiCoordinate(ostacoloVicino.id)
		if start.coordinataX < destination.coordinataX {
			forward.coordinataX = x0 - 1
		} else if start.coordinataX > destination.coordinataX {
			forward.coordinataX = x1 + 1
		}
	} else {
		forward.coordinataX = destination.coordinataX
	}

	// Viene controllato l'ostacolo sopra il punto appena spostato
	// Fallback
	ostacoloVicino = forward.posizioneOstacoloVerticale(Campo, destination.coordinataY)
	for ostacoloVicino != nil {
		var puntoE, puntoO int
		x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
		puntoE = calcolaDistanza(start.coordinataX, start.coordinataY, x1, start.coordinataY)
		puntoO = calcolaDistanza(start.coordinataX, start.coordinataY, x0, start.coordinataY)
		if puntoE < puntoO {
			forward.coordinataX = x1 + 1
			if forward.coordinataX >= start.coordinataX {
				return start
			}
		} else {
			forward.coordinataX = x0 - 1
			if forward.coordinataX <= start.coordinataX {
				return start
			}
		}
		ostacoloVicino = forward.posizioneOstacoloVerticale(Campo, destination.coordinataY)
	}
	return &forward
}

func (Campo *Piano) forwardY(start *punto, destination *punto) *punto {
	var forward punto
	forward.coordinataX = start.coordinataX

	ostacoloVicino := start.posizioneOstacoloVerticale(Campo, destination.coordinataY)
	if ostacoloVicino != nil {
		_, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
		if start.coordinataY < destination.coordinataY {
			forward.coordinataY = y0 - 1
		} else if start.coordinataY > destination.coordinataY {
			forward.coordinataY = y1 + 1
		}
	} else {
		forward.coordinataY = destination.coordinataY
	}

	ostacoloVicino = forward.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
	for ostacoloVicino != nil {
		var puntoN, puntoS int
		_, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
		puntoN = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y1)
		puntoS = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y0)
		if puntoN < puntoS {
			forward.coordinataY = y1 + 1
			if forward.coordinataY >= start.coordinataY {
				return start
			}
		} else {
			forward.coordinataY = y0 - 1
			if forward.coordinataX <= start.coordinataX {
				return start
			}
		}
		ostacoloVicino = forward.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
	}
	return &forward
}

func (Campo *Piano) cercaOstacolo(x int, y int) *ostacolo {
	testa, coda := 0, len(*Campo.ostacoli)-1
	for testa <= coda {
		x0, y0, x1, y1 := estraiCoordinate((*Campo.ostacoli)[testa].id)
		if (x <= x1 && x >= x0) && (y <= y1 && y >= y0) {
			return (*Campo.ostacoli)[testa]
		}
		x0, y0, x1, y1 = estraiCoordinate((*Campo.ostacoli)[coda].id)
		if (x <= x1 && x >= x0) && (y <= y1 && y >= y0) {
			return (*Campo.ostacoli)[coda]
		}
		testa++
		coda--
	}
	return nil
}
