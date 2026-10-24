# Die Abdeckungs-Deklaration eines Tests nennt weniger Pfade, als der Test prüft

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test mit Teilfällen deklariert nur einen Pfad (`Boundary`), obwohl einzelne Fälle Negativpfade derselben Zusage prüfen (Name löst nicht auf, Port belegt). `make abdeckung-check` ist grün, weil es Deklaration gegen Tabelle hält, nicht Deklaration gegen Test. Die Tabelle sagt weniger zu, als der Test prüft; ein Verstoß gegen `AGENTS.md` §3.11 ist es nicht (der Satz sagt nur zu, was geprüft ist), aber die Anforderung steht in der Tabelle mit weniger Nachweis, als vorhanden ist.
