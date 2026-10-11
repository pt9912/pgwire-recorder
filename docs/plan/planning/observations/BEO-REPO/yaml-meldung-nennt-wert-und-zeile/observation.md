# Die Meldungen des Lesens von `yaml` nennen Wert und Zeile, die von `sqlite` nicht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Für das Lesen von `sqlite` gilt: Die Meldung zu einem Fehler des Inhalts nennt Datei und Art, nie einen Wert
aus der Datei und keine Zeilennummer. Für das Lesen von `yaml` gilt das nicht; die Meldung enthält Zeile und
einen Teil des Inhalts, den die YAML-Bibliothek nennt. Entscheidung des Nutzers vom 2026-10-11: die Regel
gilt vorerst nur für `sqlite`. Die Beobachtung hält fest, dass dieselbe Zusage an den zwei Formaten
verschieden ausfällt.
