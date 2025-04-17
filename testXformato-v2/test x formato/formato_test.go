package main

import (
	"testing"
)

var prog = "./solution"
var verbose = true

func TestBase(t *testing.T) {

	LanciaGenericaConFileInOutAtteso(
		t,
		prog,
		"base-in",
		"base-out",
		verbose,
	)

}

func TestFormatoStato(t *testing.T) {
	casiTest := []CasoTest{
		{"automa in posizione x,y",
			`c
a 1 2 1
o 2 0 4 3
s 1 2`,
			"A\n"},
		{"ostacolo in posizione x,y",
			`c
a 1 2 1
o 2 0 4 3
s 3 1`,
			"O\n"},
		{"vuoto in posizione x,y",
			`c
a 1 2 1
o 2 0 4 3
s 7 8`,
			"E\n"},
	}
	for _, ct := range casiTest {
		t.Run(ct.nome, func(t *testing.T) {
			output := eseguiTest(ct.input)
			if output != ct.atteso {
				t.Errorf("\nInput:\n%s \n\nESECUZIONE:\n<<<<<\n%s\n>>>>\n\nATTESO:\n<<<<<\n%s\n>>>>", ct.input, output, ct.atteso)
			}
		})
	}
}

func TestFormatoStampa(t *testing.T) {
	casiTest := []CasoTest{
		{"stampa con automi e ostacoli",
			`c
a 1 2 1
a 0 3 0
o 2 0 4 3
o -2 4 1 6
S`,
			`(
0: 0,3
1: 1,2
)
[
(-2,4)(1,6)
(2,0)(4,3)
]
`},
		{"stampa con soli automi",
			`c
a 1 2 1
a 0 3 0
S`,
			`(
0: 0,3
1: 1,2
)
[
]
`},
		{"stampa con soli ostacoli",
			`c
o 2 0 4 3
o -2 4 1 6
S`,
			`(
)
[
(-2,4)(1,6)
(2,0)(4,3)
]
`},
	}

	for _, ct := range casiTest {
		t.Run(ct.nome, func(t *testing.T) {
			output := eseguiTest(ct.input)
			if output != ct.atteso {
				t.Errorf("\nInput:\n%s \n\nESECUZIONE:\n<<<<<\n%s\n>>>>\n\nATTESO:\n<<<<<\n%s\n>>>>", ct.input, output, ct.atteso)
			}
		})
	}
}

func TestFormatoEsistePercorso(t *testing.T) {
	casiTest := []CasoTest{
		{"esiste percorso libero",
			`c
a 1 2 1
e 10 2 1`,
			"SI\n"},
		{"non esiste percorso libero",
			`c
a 1 2 1
o 2 0 4 3
e 10 2 1`,
			"NO\n"},
	}

	for _, ct := range casiTest {
		t.Run(ct.nome, func(t *testing.T) {
			output := eseguiTest(ct.input)
			if output != ct.atteso {
				t.Errorf("\nInput:\n%s \n\nESECUZIONE:\n<<<<<\n%s\n>>>>\n\nATTESO:\n<<<<<\n%s\n>>>>", ct.input, output, ct.atteso)
			}
		})
	}
}

func TestFormatoPosizioni(t *testing.T) {
	casiTest := []CasoTest{
		{"posizione automi prefisso 1",
			`c
a 1 2 1
a 1 -2 10
a 0 0 11
p 1`,
			`(
10: 1,-2
11: 0,0
1: 1,2
)
`},
		{"posizione automi prefisso 0",
			`c
a 1 2 0
a 1 -2 00
a 0 0 01
p 0`,
			`(
00: 1,-2
01: 0,0
0: 1,2
)
`},
		{"posizione solo alcuni automi (prefisso 1)",
			`c
a 1 2 1
a 1 -2 010
a 0 0 11
p 1`,
			`(
11: 0,0
1: 1,2
)
`},
	}

	for _, ct := range casiTest {
		t.Run(ct.nome, func(t *testing.T) {
			output := eseguiTest(ct.input)
			if output != ct.atteso {
				t.Errorf("\nInput:\n%s \n\nESECUZIONE:\n<<<<<\n%s\n>>>>\n\nATTESO:\n<<<<<\n%s\n>>>>", ct.input, output, ct.atteso)
			}
		})
	}
}

func TestFormatoTortuosita(t *testing.T) {
	casiTest := []CasoTest{
		{"tortuosità 0",
			`c
a 1 2 1 
t 10 2 1`,
			"0\n"},
	}

	for _, ct := range casiTest {
		t.Run(ct.nome, func(t *testing.T) {
			output := eseguiTest(ct.input)
			if output != ct.atteso {
				t.Errorf("\nInput:\n%s \n\nESECUZIONE:\n<<<<<\n%s\n>>>>\n\nATTESO:\n<<<<<\n%s\n>>>>", ct.input, output, ct.atteso)
			}
		})
	}
}
