# Slice slice-v1-abschluss-antwortvergleich: Antwortvergleich beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-24.a` · `LH-FA-17.a` · `LH-FA-20.a` · `LH-FA-02.b` · `SPEC-018` · `SPEC-027` · `SPEC-034` · `SPEC-041` · `ARC-002`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play --compare-responses` vergleicht die Struktur der Serverantworten und Fehler mit der Aufzeichnung und meldet Abweichungen als `PGR-E5004` mit Exit-Code 5.

**Gebunden durch `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, Paarungen, offene Frage zu V-3; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): `--compare-responses` liest der Stand noch nicht; `nichtGelesen()` in `internal/adapters/driving/cli/fremd_test.go` führt sie. Dieser Slice liefert den Leser, kürzt `nichtGelesen()` um diese Option(en) und bindet ihn in `liesUmgebung` (`internal/adapters/driving/cli/cli.go`) an; die Zuordnung `zuordnung()` nennt sie für `play` (Spalte *Kommando* der Optionstabelle, `LH-FA-17.a`). Sonst wird `TestLeserUmgebungNichtGelesenerOptionen` rot, weil die Variable beim eigenen Kommando dann geprüft wird (`TestLeserZuordnungGleichTabelle` bleibt grün). Beim anderen Kommando (`record` und `replay`) ist die Variable schon jetzt `PGR-E2001`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-antwortvergleich-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Vergleich von Zeilenwerten und Toleranzregeln für Felder — Out-of-Scope von LH-FA-24.
- Einspielen selbst und Fehlersemantik der Serverfehler — `slice-v1-abschluss-einspielen`; seit dessen zweitem Schnitt vom 2026-10-09 die drei Optionen der Laufsteuerung ohne Vergleich `slice-v1-abschluss-einspielen-laufsteuerung` und das Einspielen von Extended-Interaktionen `slice-v1-abschluss-einspielen-extended` (je §1, *Übernimmt*). Dieser Slice setzt alle drei voraus: DoD-Punkt 3 baut auf den Optionen auf, DoD-Punkt 1 und 2 vergleichen auch Extended-Interaktionen.
- Anmeldung und TLS beim Einspielen — `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` (je §1, *Übernimmt*; seit dem Schnitt von `slice-v1-abschluss-einspielen` vom 2026-10-09); der Vergleich hängt an keinem der beiden.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen): Gegen eine Instanz mit gleicher Antwortstruktur endet das Einspielen mit Vergleich mit Erfolg, gegen eine abweichende mit `PGR-E5004` und Exit-Code 5; ein Serverfehler ohne Vorbild in der Aufzeichnung ist `PGR-E5004`, nicht `PGR-E4004`; eine Fehlerantwort `FATAL` oder `PANIC` ist `PGR-E5004`, auch bei gleichem SQLSTATE wie eine aufgezeichnete (die Aufzeichnung kennt kein Verbindungsende, `LH-FA-24.a` *Ende der Aufzeichnung*), ein Netzwerkabbruch ohne Fehlerantwort ist `PGR-E4003`; bei einer Fehlerantwort `FATAL`, `PANIC` oder einem Verbindungsverlust sendet `play` keine weiteren Nachrichten oder Gruppen der Interaktion, bei jeder anderen Abweichung wird die Interaktion bis zu ihrem `ReadyForQuery` verglichen und der Abbruch erfolgt vor der nächsten Interaktion (`LH-FA-20.a` Schritt 6); bei Extended wird nach einem Fehler bis zum `Sync` weitergesendet (Abnahmeszenario 16; Integrationstest).
- [ ] Jede Art der Abweichung (Nachrichtenart, Spaltenbeschreibung, Befehl ohne Zahlen, Fehler, Transaktionsstatus) wird erkannt und mit Session, Sequenznummer, Position und Art gemeldet (Position ab 1, bei fehlender oder zusätzlicher Nachricht die erste ohne Gegenstück, bei der Art Fehler der SQLSTATE beider Seiten und nie der Text, `LH-FA-24.a` *Randformen des Vergleichs*), Unterschiede nur in Zeilenwerten, Zeilenzahlen, Hinweisen und den abschließenden Zahlen des Befehls nicht (Test, einfach und Extended, auch `Flush`-Gruppen).
- [ ] `--continue-on-error` und `--allow-recorded-errors` wirken wie spezifiziert, `--continue-on-error` läuft nach einer Abweichung weiter, der Exit-Code am Ende ist 5, ein Verbindungsfehler beendet immer mit 4; ein Fehler mit gleichem SQLSTATE wie aufgezeichnet gilt als erwartet und bricht nicht ab, ein anderer SQLSTATE ist `PGR-E5004`; `--allow-recorded-errors` ist mit Vergleich ohne Wirkung; ein Abbruchsignal endet nach einer Abweichung mit 5, ein Verbindungsfehler nach einer Abweichung mit 4; bei jedem Ende des Laufs (auch Abbruch und Signal) trägt die Zeile zum Ende (`play beendet`, keine zweite Zeile) die Attribute eingespielt, verglichen, abweichend nach den Zählregeln von `LH-FA-24.a` und `LH-FA-20.a` Schritt 4; ohne `--compare-responses` findet kein Vergleich statt; der Leser der Option (Kommandozeile, Umgebungsvariable und Schlüssel `compare_responses` unter `play:`) ist angebunden und `nichtGelesen()` in `fremd_test.go` führt sie nicht mehr (aus `slice-v1-abschluss-umgebung-anderes-kommando`) (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Play-Service) | update | Normalisierung und Vergleich der Antworten, ohne PGWire-Typen |
| `internal/hexagon/model` | update | Katalogeintrag `PGR-E5004` (eine Konstante; zählt mit dem Play-Service als Kern, wie die Codes der früheren Play-Slices) |
| `internal/adapters/driving/cli` | update | Option `--compare-responses` (Kommandozeile, Umgebungsvariable, Schlüssel), Hilfetext, Abbildung auf Exit-Code 5 |
| `internal/bootstrap` | update | reicht die Zähler des Services an die Zeile zum Ende durch; reine Verdrahtung, keine Schicht (`AGENTS.md` §3.13), kein Lesen und kein Einstufen |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-24 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-einspielen-laufsteuerung` und `slice-v1-abschluss-einspielen-extended` sind `done`, und [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Normalisierung der Antworten erweist sich als größer als geplant — zurück zur Zerlegung. **Größenprobe** (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`): Steht nach dem Commit, der DoD-Punkte 1 und 2 erfüllt, der Diff (ohne Plan, Spezifikation und erzeugte Tabellen) über 1400 Zeilen mit Tests, hält der Implementer an und meldet es dem Planner; Schnittlinie: DoD-Punkte 1 und 2 mit Option, Leser und Hilfetext | DoD-Punkt 3 (Laufsteuerung mit Vergleich, Zähler, Signal) als eigener Slice.
- `in-progress` → `open`: Die aufgezeichneten Antworten tragen zu wenig Felder für den Vergleich der Spaltenbeschreibung — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Zusage, was nicht verglichen wird, ist schwer zu ändern; Treiber mit abweichenden Hinweisen oder Parameterstatus sind im Test zu belegen — **Ausgang:** offen bis Closure.
- Das Halten der Antworten bis zum Ende der Interaktion kostet Speicher bei großen Resultsets — **Ausgang:** offen bis Closure.
- Die Teile von `LH-FA-24.a` §Unvollständige Aufzeichnung, die aufgezeichnete Interaktionen ohne `ReadyForQuery` behandeln, sind seit der Entscheidung in `LH-FA-02.b` §Fehlerantwort vor dem Abbruch (`slice-replay-semantik-fehlerreplay`) nicht mehr erreichbar: Der Recorder schreibt keine solche Interaktion, und der Leser lehnt sie als beschädigt ab (`PGR-E3003`). Dieser Slice bereinigt sie in der Spezifikation, dazu in `LH-FA-20.a` den Teil der Tabellenzeile *Verbindung endet nach dem ersten `ReadyForQuery`* mit Vergleich, der ein aufgezeichnetes Ende behandelt, und den Satz „Wo die Aufzeichnung ein Verbindungsende zeigt, wartet der Recorder …“ (übergeben von `slice-v1-abschluss-einspielen`, das Schritt 3 bereinigt hat; die Klausel in `LH-FA-24` des Lastenhefts bleibt Sache des Nutzers); seine DoD sagt für sie nichts zu — **Ausgang:** entfallen: Der Architect hat die Spezifikation vor dem Code bereinigt (2026-10-12; `LH-FA-24.a` *Ende der Aufzeichnung*, Fehlerregeln und Katalogzeile in `LH-FA-20.a` und `SPEC-034`); der Implementer schreibt dazu nichts. Die Klauseln im Lastenheft bleiben Sache des Nutzers (Empfehlung des Architect: anpassen, siehe Bericht), sie sind ohne Wirkung auf diesen Slice, weil sie nur Fälle nennen, die keine Aufzeichnung enthält.

**Randformen** (`AGENTS.md` §3.12) — geprüft und entschieden vom Architect am 2026-10-12, vor dem Code, in der Spezifikation; der Implementer entscheidet keine:

| Randform | entschieden in |
|---|---|
| Variable beim eigenen und beim anderen Kommando (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`, §1): `--compare-responses` nur `play`; beim eigenen Kommando liest der Leser die Variable, beim anderen (`record`, `replay`) bleibt sie `PGR-E2001`; Test `TestLeserUmgebungNichtGelesenerOptionen` mit gekürztem `nichtGelesen()`, `TestLeserUmgebungAndererKommandos`. Wert der Option, Variable und Schlüssel: leer, `1`, `True` wie bei jeder booleschen Option | `LH-FA-17.a` (Tabelle, Absatz zu Umgebungsvariablen und Wertemenge, Konfigurationsdatei) |
| Normalisierung: `data_row`, `notice_response`, `parameter_status` entfallen; was je Nachricht verglichen wird | `LH-FA-24.a` (Tabelle) |
| Befehl mit und ohne Zahlen, Tag aus Ziffern allein | `LH-FA-24.a` *Randformen des Vergleichs*, *Befehl* |
| Fehlerantwort ohne Feld `C`, Schweregrad, `FATAL`/`PANIC` | `LH-FA-24.a` *Randformen des Vergleichs*, *Fehler*, und *Ende der Aufzeichnung* |
| Fehlende oder zusätzliche Nachricht, Position (ab 1), Art | `LH-FA-24.a` *Randformen des Vergleichs*, *Position*, *Art* |
| Inhalt der Meldung: SQLSTATE beider Seiten ja, Text und Felder nie, Werte nie | `LH-FA-24.a` *Meldung* |
| Extended: Folge aller Server-Nachrichten ohne Gruppengrenzen; Fehler mitten in der Interaktion; Abbruch erst am Ende der Interaktion, außer `FATAL`/`PANIC` | `LH-FA-24.a` (Einleitung), `LH-FA-20.a` Schritt 6 und *Gruppen* |
| Aufgezeichnete Interaktion ohne `ReadyForQuery`, aufgezeichnetes Verbindungsende: gibt es nicht (Leser lehnt ab) | `LH-FA-24.a` *Ende der Aufzeichnung*, `LH-FA-02.b` *Replay* |
| Verbindungsverlust ohne Fehlerantwort: `PGR-E4003`, nach einer Abweichung ebenfalls 4 | `LH-FA-20.a` (Tabelle, Absatz darunter) |
| Abbruch und Herunterfahren: erstes Signal (Interaktion wird verglichen), zweites Signal (ohne Vergleich), `--finish-session-on-interrupt` | `LH-FA-20.a` Schritt 7 und *Abbruchsignal*, `LH-FA-24.a` *Randformen des Vergleichs*, *Signal* |
| Zähler: wann eine Interaktion `eingespielt`, `verglichen`, `abweichend` zählt; die Zeile zum Ende trägt sie, keine zweite Zeile; Stufe `info`; Lauf ohne Interaktion; ohne Option alle `verglichen`/`abweichend` 0 | `LH-FA-24.a` *Meldung*, `LH-FA-20.a` Schritt 4 und *Meldungen* |
| Reihenfolge der Zeilen `error` bei Abweichungen: je Interaktion eine Zeile, in der Reihenfolge des Auftretens, abbrechende zuerst | `LH-FA-20.a` *Meldungen*, *Exit-Code* |
| `--allow-recorded-errors` mit Vergleich ohne Wirkung, keine Fehlermeldung für die Kombination | `LH-FA-20.a` (Optionstabelle, Fehlertabelle) |
| Verhalten je Serverversion: kein Ausgleich | `LH-FA-24.a` *Randformen des Vergleichs*, *Serverversion* |
| Aufruf-Voraussetzung: keine Ordnungszusage über mehrere Aufrufe; Zähler gelten für einen Lauf mit nacheinander eingespielten Interaktionen | `LH-FA-24.a` *Randformen des Vergleichs*, *Aufruf-Voraussetzung* |

**Hilfetext** (`AGENTS.md` §3.11): Der Hilfetext von `--compare-responses` sagt nur zu, was ein Test dieses Slice prüft (Vergleich der Antwortstruktur, Exit-Code 5); die Wirkung auf `--continue-on-error` und `--allow-recorded-errors` steht im Hilfetext nur, wenn ein Test sie prüft.

**Akzeptierte Negative** (Grund im Satz, `AGENTS.md` §3.12): (a) Speicher bei großen Resultsets (Risiko oben): `data_row` entfällt bei der Normalisierung, der Service muss Zeilen beim Lesen nicht halten; ob er es tut, prüft kein Test, der Plan sagt es nicht zu. (b) Der Vergleich über Serverversionen hinweg (Tag, OID) ist eine Abweichung, kein Fehler des Vergleichs; kein Test mit zwei Versionen in diesem Slice — `slice-v1-abschluss-postgres-versionen` fährt mehrere Versionen, ist aber keine Adresse dafür (kein Eintrag dort).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-antwortvergleich-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: zwei Schichten (Kern mit Play-Service, CLI-Adapter).

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 3 (Leser der Option und Kürzung von `nichtGelesen()` im Test des CLI-Adapters) und fügt keinen Liefer-Punkt hinzu: Liefer-Punkte 3. Schichten: zwei (Kern mit Play-Service, CLI-Adapter); der CLI-Adapter war schon gezählt, `fremd_test.go` liegt in ihm.

Nachgezählt bei der Prüfung des Architect (2026-10-12): Liefer-Punkte 3 (DoD 1 bis 3). Schichten zwei: Kern (Play-Service mit Katalog) und CLI-Adapter; der Bootstrap reicht nur die Zähler durch (Verdrahtung, `AGENTS.md` §3.13, keine Schicht; liest er oder stuft er ein, zählt er, und der Plan wird neu geprüft). Die Sendung aus `slice-v1-abschluss-umgebung-anderes-kommando` liegt in DoD-Punkt 3 und ist nachgezählt. **Größe:** Schätzung des Architect 1800 bis 2300 Zeilen mit Tests und Abdeckungstabelle (Produkt-Code etwa 450 bis 550: Normalisierung und Vergleich 200, Einbindung in `interaktion`/`extended` 120, Zähler und Verdrahtung 80, Option, Hilfetext und Katalog 60; Tests der Zusagen je Mutation, Unit und Integration, 1300 bis 1700); der Slice liegt an der Obergrenze einer Sitzung, nicht darüber: Er liegt unter den 2200 bis 3000 Zeilen der beiden Belege von `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und über `slice-v1-abschluss-einspielen-laufsteuerung` (etwa 900). Kein dritter Beleg; die Größenprobe in §4 fängt eine Fehlschätzung vor dem Rest der Arbeit. Tritt sie ein, ist das der dritte Beleg, den die Closure einträgt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
