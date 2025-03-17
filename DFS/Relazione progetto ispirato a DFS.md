# Relazione progetto d'esame di algoritmi e strutture dati (revisione)

## Introduzione

Questa relazione presenta le specifiche delle funzioni implementate nel programma Go fornito, che risolve un problema simile a quello descritto nella traccia d'esame "*Automi e segnali*". Il programma gestisce il movimento di automi in un piano cartesiano, tenendo conto di ostacoli e punti di richiamo. Gli automi possono spostarsi solo in determinate direzioni, evitando gli ostacoli presenti nel piano.

Nella relazione si farà riferimento alla distanza di Manhattan fra due punti del piano con $D$, al numero di automi nel piano con $a$ e al numero di ostacoli con $m$.

## Strutture dati e scelte progettuali

Per rappresentare il piano, il programma utilizza una struttura con due liste concatenate: una per gli automi e una per gli ostacoli. Questa scelta permette di gestire dinamicamente l'aggiunta e la modifica di automi e ostacoli, mantenendo una struttura leggera in termini di consumo di memoria e facile da manipolare.

### Strutture dati principali

- **`punto`**: rappresenta un punto nel piano. Contiene le coordinate `x` e `y`, un identificativo `id` e un puntatore a un altro `punto`. Questa struttura è utilizzata per rappresentare sia automi che ostacoli.
- **`Piano`**: struttura principale che mantiene riferimenti a due liste concatenate: una per gli automi e una per gli ostacoli.
- **`piano`**: alias di un tipo puntatore a una variabile `Piano`.
- **`elementoPila`**: struttura utilizzata per gestire l'operazione di richiamo, memorizzando gli automi candidati allo spostamento. Contiene un puntatore a un `punto` (l'automa candidato), un intero che rappresenta la distanza dal richiamo e un puntatore al prossimo elemento della pila.
- **`cella`**: struttura utilizzata per rappresentare una cella nel piano durante l'esecuzione dell'algoritmo DFS. Contiene le coordinate `x` e `y`, il numero di passi effettuati e un puntatore alla cella genitore.
- **`visitata`**: struttura utilizzata per tenere traccia delle celle già visitate durante l'esecuzione dell'algoritmo DFS.

## Implementazione e tempi delle operazioni richieste

L'operazione `crea` viene implementata dalla funzione `newPiano`, che restituisce un nuovo piano vuoto. Questa operazione impiega tempo costante $O(1)$.

L'operazione `stato` viene implementata dal metodo `stato`, che scorre le liste degli automi e degli ostacoli per determinare se una posizione specifica è occupata da un automa, un ostacolo o è libera. Questa operazione richiede tempo $O(a + m)$ nel caso peggiore.

L'operazione `stampa` è implementata dal metodo `stampa`, che scorre entrambe le liste degli automi e degli ostacoli e stampa le loro coordinate. Questa operazione impiega tempo $\Theta(a + m)$.

L'operazione `automa` è implementata dal metodo `automa`, che aggiunge o rimuove un automa dalla lista degli automi. Questa operazione impiega tempo $\Theta(a + m)$ nel caso peggiore, poiché deve verificare la presenza di ostacoli nella posizione specificata.

L'operazione `ostacolo` è implementata dal metodo `ostacolo`, che aggiunge un ostacolo alla lista degli ostacoli. Questa operazione impiega tempo $\Theta(a)$, poiché deve verificare che nessun automa si trovi all'interno dell'area dell'ostacolo.

L'operazione `posizioni` è implementata dal metodo `posizioni`, che scorre la lista degli automi e stampa le posizioni degli automi il cui identificativo inizia con un determinato prefisso. Questa operazione impiega tempo $O(a)$.

### Movimenti e percorsi degli automi

L'operazione `richiamo` è implementata dal metodo `richiamo`, che controlla gli automi più vicini a un punto di richiamo e li sposta verso di esso. Il metodo utilizza una pila per memorizzare gli automi candidati e calcola la distanza di Manhattan tra ciascun automa e il punto di richiamo. Gli automi con la distanza minima vengono spostati. Il tempo di esecuzione di questa operazione è $O(a \cdot D^2 \cdot m)$ nel caso peggiore, poiché deve eseguire la funzione `dfs` per ciascun automa.

L'operazione `esistePercorso` è implementata dal metodo `esistePercorso`, che verifica se esiste un percorso libero tra un automa e un punto di destinazione. Il metodo utilizza la funzione `dfs` per determinare se il percorso esiste. Il tempo di esecuzione di questa operazione è $O(D^2 \cdot m)$, poiché deve eseguire la funzione `dfs`.

#### Funzione `dfs`

La funzione `dfs` implementa un algoritmo di ricerca in profondità (DFS) per determinare se esiste un percorso libero tra due punti. La funzione tiene conto degli ostacoli presenti nel piano e verifica se è possibile raggiungere la destinazione. Il tempo di esecuzione di questa funzione è $O(D^2 \cdot m)$, poiché nel caso peggiore deve esplorare tutte le possibili direzioni per ogni passo.

## Conclusioni

Il programma Go implementa le operazioni richieste per gestire il movimento degli automi in un piano cartesiano con ostacoli e punti di richiamo. Le strutture dati scelte, come le liste concatenate, permettono una gestione dinamica degli automi e degli ostacoli, mentre le funzioni di ricerca e spostamento garantiscono che gli automi si muovano in modo efficiente evitando gli ostacoli. I tempi di esecuzione delle operazioni sono stati ottimizzati per garantire prestazioni accettabili anche in scenari complessi. L'uso dell'algoritmo DFS per la ricerca del percorso offre una soluzione alternativa all'approccio BFS, con tempi di esecuzione simili ma con una diversa strategia di esplorazione del piano.
