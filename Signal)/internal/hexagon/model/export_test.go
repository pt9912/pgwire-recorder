package model

// Export-Test-Brücke (SPEC-049 Punkt 7): reicht an Unexportiertes weiter, ohne
// eigenen Zustand auf Paketebene.

// Klasse reicht an klassen weiter.
func Klasse(ziffer int) string { return klassen[ziffer] }
