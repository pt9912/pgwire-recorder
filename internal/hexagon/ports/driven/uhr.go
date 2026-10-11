package driven

import "context"

// Uhr ist der Uhr-Port (ARC-004, LH-FA-21.a *Uhr*). Zeit ist eine ganze Zahl
// Millisekunden (int64) seit einem Ursprung, den die Uhr wählt; ein Wert hat nur
// als Unterschied zweier Lesungen derselben Uhr eine Bedeutung.
type Uhr interface {
	// Jetzt liest die Zeit, abgeschnitten auf ganze Millisekunden und nie kleiner
	// als die vorige Lesung derselben Uhr.
	Jetzt() int64
	// WarteBis wartet, bis Jetzt den Zeitpunkt erreicht hat. Ist ctx beendet,
	// kehrt der Aufruf mit dem Fehler von ctx zurück, vor und während des
	// Wartens und auch bei einem erreichten oder vergangenen Zeitpunkt; sonst
	// kehrt er mit nil zurück, bei einem erreichten oder vergangenen Zeitpunkt
	// sofort.
	WarteBis(ctx context.Context, zeitpunkt int64) error
}
