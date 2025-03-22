package main

import (
	"bufio"
	. "fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

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
	percorsoEffettuato := (avanza(Campo, percorrente, Sorgente))
	if percorsoEffettuato.coordinataX == x && percorsoEffettuato.coordinataY == y {
		Println("SI")
		return
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
			possibileAvanzamento := avanza(Campo, percorrente, Sorgente)
			if possibileAvanzamento.coordinataX == x && possibileAvanzamento.coordinataY == y {
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

func (Campo *Piano) ostacoliPercorso(partenza, arrivo *punto) (distanza_O_Ascisse, distanza_O_Ordinate int) {
	ostacoloVicino := partenza.posizioneOstacoloVerticale(Campo, arrivo.coordinataY)
	if ostacoloVicino != nil {
		_, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
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
		x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
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

func (p *punto) posizioneOstacoloVerticale(Campo piano, y int) *punto {
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

func (Campo *Piano) forwardX(start *punto, destination *punto) *punto {
	var forward punto
	ostacoloVicino := start.posizioneOstacoloVerticale(Campo, destination.coordinataY)
	if ostacoloVicino != nil {
		var puntoE, puntoO int
		x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
		puntoE = calcolaDistanza(destination.coordinataX, destination.coordinataY, x1, destination.coordinataY)
		puntoO = calcolaDistanza(destination.coordinataX, destination.coordinataY, x0, destination.coordinataY)
		if puntoE < puntoO {
			forward.coordinataX = x1 + 1
			forward.coordinataY = start.coordinataY
		} else {
			forward.coordinataX = x0 - 1
			forward.coordinataY = start.coordinataY
		}
		return &forward
	}
	ostacoloVicino = start.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
	if ostacoloVicino != nil {
		x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
		if start.coordinataX < destination.coordinataX {
			forward.coordinataX = x0 - 1
		} else if start.coordinataX > destination.coordinataX {
			forward.coordinataX = x1 + 1
		}
	} else {
		forward.coordinataX = destination.coordinataX
	}
	forward.coordinataY = start.coordinataY
	osostacoloVicino := forward.posizioneOstacoloVerticale(Campo, destination.coordinataX)
	if osostacoloVicino != nil {
		var puntoE, puntoO int
		x0, _, x1, _ := estraiCoordinate(osostacoloVicino.id)
		puntoE = calcolaDistanza(start.coordinataX, start.coordinataY, x1, start.coordinataY)
		puntoO = calcolaDistanza(start.coordinataX, start.coordinataY, x0, start.coordinataY)
		if puntoE < puntoO {
			forward.coordinataX = x1 + 1
		} else {
			forward.coordinataX = x0 - 1
		}
	}
	return &forward
}

func (Campo *Piano) forwardY(start *punto, destination *punto) *punto {
	var forward punto
	ostacoloVicino := start.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
	if ostacoloVicino != nil {
		var puntoN, puntoS int
		_, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
		puntoN = calcolaDistanza(destination.coordinataX, destination.coordinataY, destination.coordinataX, y1)
		puntoS = calcolaDistanza(destination.coordinataX, destination.coordinataY, destination.coordinataX, y0)
		if puntoN < puntoS {
			forward.coordinataY = y1 + 1
			forward.coordinataX = start.coordinataX
		} else {
			forward.coordinataY = y0 - 1
			forward.coordinataX = start.coordinataX
		}
		return &forward
	}
	forward.coordinataX = start.coordinataX
	ostacoloVicino = start.posizioneOstacoloVerticale(Campo, destination.coordinataY)
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
	if ostacoloVicino != nil {
		var puntoN, puntoS int
		_, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
		puntoN = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y1)
		puntoS = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y0)
		if puntoN < puntoS {
			forward.coordinataY = y1 + 1
		} else {
			forward.coordinataY = y0 - 1
		}
		return &forward
	}
	return &forward
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
