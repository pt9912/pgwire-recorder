# Slice slice-v1-abschluss-antwortvergleich-doku: Handbuch und README zu: Antwortvergleich beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-antwortvergleich`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-antwortvergleich`, abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-antwortvergleich` über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt jenen Slice voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt dahinter. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit beiden.

**Übernommen aus `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, offene Frage zu V-3; Planner, 2026-10-12, `AGENTS.md` §3.13): Das Handbuch nennt `--compare-responses` erst mit `slice-v1-abschluss-antwortvergleich` und nennt dann zu jeder Option ihre Umgebungsvariable mit den Kommandos, bei denen sie gilt; die Variable beim anderen Kommando ist `PGR-E2001` (Satz zur Meldung `PGR-E2001` im Handbuch, §5 *Einstellungen*). Das Handbuch beschreibt die Variable, nachdem das gebaute Binary sie liest.

**Gebunden durch `slice-v1-abschluss-antwortvergleich`** (dort §6 *Randformen* und §8; Sendung des Architect, 2026-10-12, `AGENTS.md` §3.13, §3.11): Das Handbuch trägt erst mit diesem Slice, was das Binary nach jenem Slice zeigt, und nur das, was ein Test jenes Slice prüft: `--compare-responses` bei `play` (Variable `PGWIRE_RECORDER_COMPARE_RESPONSES`, Schlüssel `compare_responses` unter `play:`); verglichen werden Art und Reihenfolge der Antworten, Spaltenbeschreibung (Anzahl, Name, Typ), Beschreibung der Parameter, der Befehl ohne abschließende Zahlen, der SQLSTATE eines Fehlers und der Transaktionsstatus, nicht Zeilenwerte, Zeilenzahlen, Hinweise und Einstellungen des Servers; eine Abweichung ist `PGR-E5004` mit Exit-Code 5 und nennt Session, Sequenznummer, Position und Art (bei einem Fehler die SQLSTATE beider Seiten, nie den Text), je Interaktion die erste; ein Serverfehler ohne Vorbild und eine Fehlerantwort `FATAL` oder `PANIC` sind `PGR-E5004`, ein Netzwerkabbruch ohne Fehlerantwort ist `PGR-E4003`; mit Vergleich ersetzt dieser die Fehlerregel (`PGR-E4004` entfällt), `--allow-recorded-errors` bleibt ohne Wirkung, `--continue-on-error` setzt nach einer Abweichung fort; die Zeile zum Ende des Laufs trägt `eingespielt`, `verglichen` und `abweichend` auf der Stufe `info`; ohne Option findet kein Vergleich statt. Welche dieser Sätze ins Handbuch kommen, bestimmt die Probe am gebauten Binary; was sie nicht zeigt, bleibt draußen. Zum Handbuch-Teil in DoD-Punkt 1; kein neuer Liefer-Punkt, Schichten weiter eine.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-antwortvergleich` liefert.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.
- Das Verhalten bei einem aufgezeichneten Ende des Servers (`FATAL`/`PANIC` an derselben Stelle mit demselben SQLSTATE ist erwartet, nicht `PGR-E5004`) — `slice-v1-abschluss-verbindungsende-doku` (dort §1, *Übernommen aus `slice-v1-abschluss-antwortvergleich-doku`*; Entscheidung des Nutzers vom 2026-10-12, Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Sendung des Planners, `AGENTS.md` §3.13): Der Satz in „Gebunden durch“, `FATAL` und `PANIC` seien `PGR-E5004`, beschreibt den Ist-Zustand, den dieser Slice nach `slice-v1-abschluss-antwortvergleich` vorfindet; `slice-v1-abschluss-verbindungsende-vergleich` ändert das Binary dahinter, und der Doku-Slice dazu ersetzt den Satz. Dieser Slice schreibt ihn so, wie das Binary heute zeigt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [ ] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries (zu `--compare-responses` mit Umgebungsvariable und dem Verhalten beim anderen Kommando, aus `slice-v1-abschluss-umgebung-anderes-kommando`; Probe am gebauten Binary): Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [ ] `README.md` nennt, was `slice-v1-abschluss-antwortvergleich` liefert, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-antwortvergleich` liefert: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt |
| `README.md` | update | Sätze, die `slice-v1-abschluss-antwortvergleich` überholt, im Ist-Zustand |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-antwortvergleich` liegt in `done/` (Voraussetzung; Schritt 24 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter jenem Slice; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-antwortvergleich` eine weitere Randform liefert, die dort nicht steht.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** offen bis Closure.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — genau der Gegenstand dieses Slice: Das Handbuch sagt nur zu, was das Binary zeigt.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-antwortvergleich`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch), kein neuer Liefer-Punkt; Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-antwortvergleich` (Architect, 2026-10-12): Sie gehört zu DoD-Punkt 1 (Handbuch); Liefer-Punkte weiter 2, Schichten weiter eine.

Nachgezählt beim Eintragen der Abgrenzung zum aufgezeichneten Ende des Servers (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie nimmt dem Slice nichts und fügt nichts hinzu; Liefer-Punkte weiter 2, Schichten weiter eine.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
