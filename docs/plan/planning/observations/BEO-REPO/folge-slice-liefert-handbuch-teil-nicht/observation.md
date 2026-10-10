# Ein Folge-Slice liefert seinen Teil von Handbuch und README nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`AGENTS.md` §3.11 verlangt für Handbuch und README den Ist-Zustand; der Teil, den ein Slice
liefert, steht in seiner DoD oder in einem Doku-Folge-Slice direkt hinter ihm. Das Risiko: Ein
Nehmer oder ein Doku-Folge-Slice schließt, ohne dass sein Teil im Handbuch steht, weil nur eine
Zeile der DoD ihn trägt und das Review sie übersieht. Handbuch und Binary laufen dann
auseinander, ohne dass ein Gate es meldet (ein Sensor, der das Handbuch gegen das Binary prüft,
ist ausgeschlossen, `slice-doku-ist-stand` §1).

Abgrenzung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Dort sagt ein Text mehr zu, als
geprüft ist. Hier fehlt der Text, den der Slice schuldet.
