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
}

type piano *Piano

type Piano struct {
	automi   *[]punto
	ostacoli *[]punto
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
	Println("(")
	for i := 0; i < len(*Campo.automi); i++ {
		if strings.HasPrefix((*Campo.automi)[i].id, alpha) {
			Printf("%s: %d,%d\n", (*Campo.automi)[i].id, (*Campo.automi)[i].coordinataX, (*Campo.automi)[i].coordinataY)
		}
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
	NewPunto := new(Piano)
	NewPunto.automi = new([]punto)
	NewPunto.ostacoli = new([]punto)
	return NewPunto
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
		automi := append(*Campo.automi, *puntoCercato)
		Campo.automi = &automi
		return
	}
}

func (Campo *Piano) ostacolo(x0, y0, x1, y1 int) {
	testa, coda := 0, len(*Campo.automi)-1
	for testa <= coda {
		if ((*Campo.automi)[testa].coordinataX <= x1 && (*Campo.automi)[testa].coordinataX >= x0) && ((*Campo.automi)[testa].coordinataY <= y1 && (*Campo.automi)[testa].coordinataY >= y0) {
			return
		}
		if ((*Campo.automi)[coda].coordinataX <= x1 && (*Campo.automi)[coda].coordinataX >= x0) && ((*Campo.automi)[coda].coordinataY <= y1 && (*Campo.automi)[coda].coordinataY >= y0) {
			return
		}
		testa++
		coda--
	}
	newOstacolo := new(punto)
	newOstacolo.coordinataX = x0
	newOstacolo.coordinataY = y1
	newOstacolo.id = Sprintf("%d,%d,%d,%d,ostacolo", x0, y0, x1, y1)
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

// Funzione ricorsiva che simula l'avanzamento dell'automa sui suoi assi per raggiungere il segnale 
func avanza(Campo piano, p *punto, Sorgente *punto) *punto {
    passi := calcolaDistanza(p.coordinataX, p.coordinataY, Sorgente.coordinataX, Sorgente.coordinataY)
    if passi <= 0 || p.coordinataX == Sorgente.coordinataX && p.coordinataY == Sorgente.coordinataY {
        return p
    }
    
    possibilePasso := new(punto)
    possibilePasso.coordinataX = p.coordinataX
    possibilePasso.coordinataY = p.coordinataY
    possibilePasso.id = p.id
    
    // Usa la funzione forward per determinare il prossimo passo
    possibilePasso = (*Campo).forward(possibilePasso, Sorgente)
    
    // Se non è possibile avanzare, ritorna il punto corrente
    if p.coordinataX == possibilePasso.coordinataX && p.coordinataY == possibilePasso.coordinataY {
        return possibilePasso
    }
    
    // Altrimenti, continua a cercare il percorso
    return avanza(Campo, possibilePasso, Sorgente)
}

func (Campo *Piano) forward(start *punto, destination *punto) *punto {
    var forward punto
    forward.coordinataX = start.coordinataX
    forward.coordinataY = start.coordinataY
    forward.id = start.id
    
    // Determina l'asse principale su cui muoversi (quello con distanza maggiore)
    distanzaX := int(math.Abs(float64(destination.coordinataX - start.coordinataX)))
    distanzaY := int(math.Abs(float64(destination.coordinataY - start.coordinataY)))
    
    if distanzaX > distanzaY {
        // Movimento principale lungo asse X
        ostacoloVicino := start.posizioneOstacoloVerticale(Campo, destination.coordinataX)
        if ostacoloVicino != nil {
            x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
            if start.coordinataX < destination.coordinataX {
                forward.coordinataX = x1 + 1
            } else if start.coordinataX > destination.coordinataX {
                forward.coordinataX = x0 - 1
            }
        } else {
            ostacoloVicino := start.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
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
        }
        
        // Verifica se il nuovo punto ha ostacoli verticali verso la destinazione
        ostacoloVicino = forward.posizioneOstacoloVerticale(Campo, destination.coordinataY)
        for ostacoloVicino != nil {
            // Se c'è un ostacolo, prova a spostarsi lateralmente
            var puntoE, puntoO int
            x0, _, x1, _ := estraiCoordinate(ostacoloVicino.id)
            puntoE = calcolaDistanza(start.coordinataX, start.coordinataY, x1, start.coordinataY)
            puntoO = calcolaDistanza(start.coordinataX, start.coordinataY, x0, start.coordinataY)
            if puntoE < puntoO {
                forward.coordinataX = x1 + 1
                if forward.coordinataX >= start.coordinataX && start.coordinataX > destination.coordinataX {
                    return start // Non possiamo fare progressi
                }
            } else {
                forward.coordinataX = x0 - 1
                if forward.coordinataX <= start.coordinataX && start.coordinataX < destination.coordinataX {
                    return start // Non possiamo fare progressi
                }
            }
			ostacoloVicino = forward.posizioneOstacoloVerticale(Campo, destination.coordinataY)
        }
    } else {
        // Movimento principale lungo asse Y
        ostacoloVicino := start.posizioneOstacoloOrizzontale(Campo, destination.coordinataY)
        if ostacoloVicino != nil {
            _, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
            if start.coordinataY < destination.coordinataY {
                forward.coordinataY = y1 + 1
            } else if start.coordinataY > destination.coordinataY {
                forward.coordinataY = y0 - 1
            }
        } else {
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
        }
        
        // Verifica se il nuovo punto ha ostacoli orizzontali verso la destinazione
        ostacoloVicino = forward.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
        for ostacoloVicino != nil {
            // Se c'è un ostacolo, prova a spostarsi verticalmente
            var puntoN, puntoS int
            _, y0, _, y1 := estraiCoordinate(ostacoloVicino.id)
            puntoN = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y1)
            puntoS = calcolaDistanza(start.coordinataX, start.coordinataY, start.coordinataX, y0)
            if puntoN < puntoS {
                forward.coordinataY = y1 + 1
                if forward.coordinataY >= start.coordinataY && start.coordinataY > destination.coordinataY {
                    return start // Non possiamo fare progressi
                }
            } else {
                forward.coordinataY = y0 - 1
                if forward.coordinataY <= start.coordinataY && start.coordinataY < destination.coordinataY {
                    return start // Non possiamo fare progressi
                }
            }
			ostacoloVicino = forward.posizioneOstacoloOrizzontale(Campo, destination.coordinataX)
        }
    }
    
    // Verifica che ci sia stato effettivamente un movimento
    if forward.coordinataX == start.coordinataX && forward.coordinataY == start.coordinataY {
        return start // Nessun progresso possibile
    }
    
    return &forward
}

func (Campo *Piano) cercaOstacolo(x int, y int) *punto {
	testa, coda := 0, len(*Campo.ostacoli)-1
	for testa <= coda {
		x0, y0, x1, y1 := estraiCoordinate((*Campo.ostacoli)[testa].id)
		if (x <= x1 && x >= x0) && (y <= y1 && y >= y0) {
			return &(*Campo.ostacoli)[testa]
		}
		x0, y0, x1, y1 = estraiCoordinate((*Campo.ostacoli)[coda].id)
		if (x <= x1 && x >= x0) && (y <= y1 && y >= y0) {
			return &(*Campo.ostacoli)[coda]
		}
		testa++
		coda--
	}
	return nil
}

func (Campo *Piano) cercaAutoma(x, y int, id string) *punto {
	testa, coda := 0, len(*Campo.automi)-1
	for testa <= coda {
		if ((*Campo.automi)[testa].coordinataX == x && (*Campo.automi)[testa].coordinataY == y) || ((*Campo.automi)[testa].id == id) {
			return &(*Campo.automi)[testa]
		}

		if ((*Campo.automi)[coda].coordinataX == x && (*Campo.automi)[coda].coordinataY == y) || ((*Campo.automi)[coda].id == id) {
			return &(*Campo.automi)[coda]
		}
		testa++
		coda--
	}
	return nil
}
