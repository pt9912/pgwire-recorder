# Die Diagnose des Replays nennt nach einem gescheiterten `parse` auf das unbenannte Statement ein Statement, das der Server nicht mehr kennt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein scheiterndes `parse` auf das unbenannte Statement gibt das alte frei (PostgreSQL 14 bis 18, `26000` beim folgenden `bind`); das Replay führt das alte SQL weiter und nennt es in der Diagnose einer späteren Abweichung. Betroffen ist der Text der Diagnose, nicht die Antworten des Replays. Die Spezifikation führt es als Grenze von v1 (`LH-FA-18.a`, *Grenze (unbenanntes Statement)*), als Empfehlung des Architect; eine Korrektur müsste den Erfolg von `execute` und `describe` je Nachricht nachbilden.
