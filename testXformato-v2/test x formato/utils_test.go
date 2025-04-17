package main

import (
	"fmt"
	"sort"
	"strings"
)

type CasoTest struct {
	nome   string
	input  string
	atteso string
}

// si assume che solo il penultimo comando nella stringa di input del test stampi su stdin
func eseguiTest(input string) string {
	lines := strings.Split(input, "\n")
	var p piano

	// si eseguono comandi che precedono la stampa su stdin
	for _, line := range lines[:len(lines)-1] {
		if line[0] == 'c' {
			p = newPiano()
		} else {
			esegui(p, line)
		}
	}

	// cattura stringa stampata su stdin
	comando := lines[len(lines)-1]
	cOper := strings.Fields(comando)[0]
	output := CaptureOutput(esegui, p, lines[len(lines)-1])

	// ordina output nel caso di comandi "p" e "S"
	if cOper == "p" {
		output = ordinaOutputPosizioni(output[:len(output)-1]) // rimuove il carattere di newline finale
	} else if cOper == "S" {
		output = ordinaOutputStampa(output[:len(output)-1]) // rimuove il carattere di newline finale
	}
	return output
}

func ordinaLineeTraDelimitatori(input string, delSx string, delDx string) string {
	lines := strings.Split(input, "\n")
	lines = lines[1 : len(lines)-1]
	if len(lines) == 0 {
		return input
	}

	sort.Strings(lines)
	sortedInput := strings.Join(lines, "\n")
	return fmt.Sprintf("%s\n%s\n%s", delSx, sortedInput, delDx)

}
func ordinaOutputPosizioni(input string) string {
	return fmt.Sprintf("%s\n", ordinaLineeTraDelimitatori(input, "(", ")"))
}

func ordinaOutputStampa(input string) string {
	separator := ")\n["
	parts := strings.SplitN(input, separator, 2)
	if len(parts) != 2 {
		return input // se la separazione tra lista di automi e lista di ostacoli non è andata a buon fine, l'input iniziale verrà comunque rilevato come errato
	}

	return fmt.Sprintf("%s\n%s\n",
		ordinaLineeTraDelimitatori(parts[0]+")", "(", ")"),
		ordinaLineeTraDelimitatori("["+parts[1], "[", "]"))
}
