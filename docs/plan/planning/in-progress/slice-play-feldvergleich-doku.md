# Slice slice-play-feldvergleich-doku: Handbuch und README zur Option für alle Abweichungen beim Vergleich

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** — (wellenlos), wie `slice-play-feldvergleich`.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)

**Berührte Spec-Stellen:** `LH-FA-24.a` (*Meldung* und *Alle Abweichungen*: Optionsname `--diff-all`, Einzelheiten je Art, Obergrenze 20, `PGR-E2001` ohne `--compare-responses`; gelesen, der Slice ändert die Spezifikation nicht) · `LH-FA-24`

**Verantwortlich:** Implementer (Zweig `impl-diff-all-doku`)
**Autor:** Planner. **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**.

**Ziel:** Benutzerhandbuch und `README.md` nennen die Option `--diff-all` zur Meldung aller Abweichungen (Name, Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all`, Meldung, Obergrenze 20, Kombination mit `--compare-responses`) so, wie das gebaute Binary sie zeigt, jeder Satz mit einem Lauf gegen das Binary belegt (`AGENTS.md` §3.11).

**Sendung von `slice-play-feldvergleich`:** dessen Handbuch- und README-Teil (`AGENTS.md` §3.11 und §3.13); dieser Slice ist sein Doku-Folge-Slice und nimmt sie an. Der Geber liegt in `done/`; sein Stand (§7, `LH-FA-24.a` *Alle Abweichungen*) ist der Inhalt: ein Eintrag je abweichender Spalte mit Name und Typ-OID (Name zuerst), die Obergrenze 20 zählt Spalten, nicht Felder; höchstens 20 Positionen je Interaktion mit „und N weitere“; ohne `--compare-responses` `PGR-E2001`, aus jeder Quelle; `record` und `replay` kennen die Option nicht. Das Handbuch sagt davon nur, was ein Lauf gegen das gebaute Binary zeigt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Spezifikation, Lastenheft und Tests des Produkts — Schicht-Abgrenzung (Klasse 4). Grund: Sie gehören `slice-play-feldvergleich`; dieser Slice ist reine Doku und beim Review sofort prüfbar.
- Beschreibung von Verhalten, das das Binary nicht zeigt (Zielstand, „kommt mit“) — Bestand bleibt stehen (Klasse 2). Grund: Handbuch und README beschreiben den Ist-Zustand (`AGENTS.md` §3.11).
- Ein Vergleich der Zeilenwerte — anderer Vorgang (Klasse 3). Grund: Das Lastenheft schließt ihn aus; das Handbuch sagt weiter, dass Werte nicht verglichen werden.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Handbuch** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): der Abschnitt zum Vergleich der Antworten in `docs/user/benutzerhandbuch.md` (der Satz „Je Interaktion meldet `play` die erste Abweichung“ und jede Stelle, die der neuen Meldung widerspräche: Optionstabelle, Fehlercodes `PGR-E5004` und `PGR-E2001`, Abschnitt zum Ende der Verbindung, Exit-Code 5) nennt die Option, ihre drei Quellen, die Form der Meldung und die Obergrenze; jeder Satz belegt durch einen Lauf gegen das gebaute Binary (Datei, nicht gelesen gegen den Code), keine Chronik, kein Verweis auf Spezifikation, ADR, Slice oder Welle.
- [x] **Liefer-Punkt 2 — README** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): die fünf Stellen von `--compare-responses` in `README.md` sind durchgesehen; wo sie die erste Abweichung als einzige beschreiben, ist die Option genannt, sonst bleibt der Text, mit Grund in §7.
- [x] `make gates` am Endstand: grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitt zum Vergleich, Optionstabelle, Fehlercodes `PGR-E5004`/`PGR-E2001`, Stellen zu „erste Abweichung“: Option, Meldung, Obergrenze, Spalte als ein Eintrag mit Name und Typ-OID; Kopf (Version des Dokuments, Stand) |
| `README.md` | update | Stellen zu `--compare-responses`, soweit sie die erste Abweichung beschreiben |
| `test/integration/play_vergleich_handbuch_e2e_test.go` | neu | E2E-Fälle gegen das gebaute Binary für jeden Satz des Handbuchs (nur Test der Doku, kein Produkt-Code) |
| `docs/user/abdeckung-e2e.md` | update (erzeugt) | `make abdeckung` schreibt sie aus den Deklarationen des neuen Tests |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-play-feldvergleich` liegt in `done/` (erfüllt am Stand `db5d12ca`; der Wortlaut des Binaries steht fest). WIP-Limit 1.

**Rückführungen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik: `AGENTS.md` §3.14; Pfade des Plans: `docs/user README.md test`); Basis: der Commit des Übergangs nach `in-progress/`, Hash beim Übergang eintragen: `3d3832d9` (Zweig `impl-diff-all-doku`).
- `in-progress` → `open`: das Binary zeigt ein Verhalten, das die Spezifikation anders sagt — Rückgabe an den Architect, Slice wartet nicht in `in-progress/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben; die Welle bzw. ein Release, das die Option enthält, entsteht erst nach diesem Slice.

## 6. Risiken und offene Punkte

- Ein Handbuchsatz sagt mehr zu, als das Binary zeigt (Wortlaut der Meldung, Obergrenze) — **Ausgang:** offen bis Closure (jeder Satz mit Lauf belegt, DoD 1).
- Die Grenze 20 zählt Spalten (ein Eintrag mit Name und Typ-OID), nicht Felder; ein Handbuchsatz, der „Felder“ oder „Einträge der Liste“ zählt, sagt es falsch — **Ausgang:** offen bis Closure (ein Lauf mit 20 und 21 Spalten mit beiden Feldern belegt den Satz, Geber `slice-play-feldvergleich` §7).
- Das Handbuch hinkt dem Binary hinterher, solange `slice-play-feldvergleich` ohne diesen Slice in `done/` liegt — **Ausgang:** offen bis Closure (Release erst nach beiden, `AGENTS.md` §3.11).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*
### Belege des Implementers (Stand `5f35239f`)

**Läufe:** `make gates` am Endstand grün (Exit 0, einschließlich `test-integration`, `test-postgres-versionen`, `docs-check`, `abdeckung-check` nach `make abdeckung`, `lint`, `test`). Engster Fall davor: `make test-integration` mit `PGR_TESTS_1` auf die vier neuen Tests, alle PASS (die Phase ohne PostgreSQL dieser Teilläufe ist ohne die Aufzeichnungen der vollen Phase 1 rot und zählt nicht). **Marke 1500** (§3.14, Basis `3d3832d9`, Pfade `docs/user README.md test`): 340 Zeilen.

**Handbuchaussage · Test (Datei `test/integration/play_vergleich_handbuch_e2e_test.go`; Beispiele als Zeile aus dem Lauf gegen das gebaute Binary):**

| Aussage im Handbuch | Test | Mutation, rot gesehen |
|---|---|---|
| Option, Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all` wirken (Beispiel und Konfiguration) | `TestE2EPlayHandbuchDiffAllQuellen` (Option, Datei, Umgebung) | — (Geltung der Quellen gehört dem Geber; hier Beleg des Satzes) |
| Beispielzeile Typ-OID, Beispielzeile COMMIT gegen BEGIN (Zeile samt `play beendet`-Zähler) | `…Einzelheiten` (Fälle Typ-OID, Transaktionsstatus: Zeile `Replay [PGR-E5004]: …"` vollständig, Zähler) | — |
| Einzelheiten je Art (`Nachrichtenart` keine; Spalten: Zahl, Nummer, Name zuerst, Typ-OID; Befehl; Fehler; Transaktionsstatus), keine Zeilenwerte, kein Fehlertext | `…Einzelheiten` (sechs Fälle) | — |
| Form auch bei einer Abweichung; eine Zeile je Interaktion; Interaktion ohne Abweichung keine Zeile | `…Einzelheiten` (`Abweichungen: 1`), `…Fortsetzung` (zwei Zeilen bei zwei Interaktionen von fünf) | — |
| höchstens 20 Positionen, dahinter `und N weitere`; Gesamtzahl auch bei Kürzung | `…Obergrenzen` (20, 21, 23 Abweichungen) | `maxPositionen` 20→21: 21_ und 23_Abweichungen rot; Kürzung entfernt (`&& false`): beide rot |
| je Spaltenbeschreibung höchstens 20 Spalten, dahinter `und N weitere`; Grenze zählt Spalten, nicht Angaben (20 Spalten mit Name und Typ-OID = 40 Angaben, alle genannt) | `…Obergrenzen` (20, 21, 23 Spalten) | `maxEinzelne` 20→21: 21_ und 23_Spalten rot (20_Spalten bleibt grün, es bricht keine Zusage) |
| `--continue-on-error`: Liste je Interaktion, weiter, `abweichend` zählt je Interaktion einmal; ohne: Ende nach der ersten Interaktion mit vollständiger Liste; Exit 5 | `…Fortsetzung` (5/5/2 und 2/2/1) | — |
| ohne `--compare-responses` `PGR-E2001`, Exit 2, auch `=false`, auch aus Umgebung und Datei (nicht `PGR-E2004`), nennt beide Optionen | `…Quellen` (vier Fälle) | Kombinationsprüfung entfernt (`cli.go`, Bedingung `false`): die vier Fälle ohne `--compare-responses` rot |
| `record`/`replay`: Option und Umgebungsvariable `PGR-E2001`, Schlüssel `PGR-E2004` | `…Quellen` (fünf Fälle) | — |

Mutationen nur an den Zusagen, die der Handbuchsatz über §3.10 hinaus an eine Grenze bindet (Obergrenzen, Kombination); die übrigen Zeilen sind Belege der Doku (§3.11, Test je Zusage), die Mutationen des Produkts liegen beim Geber. Der Mutant stand je in einer Kopie derselben Datei mit `touch`, danach zurückgesetzt (`git status` zeigte nur Doku und Test).

**Durchsicht der Quantoren (Schritt 17):** Treffer zu „erste Abweichung“: Zeile 1096 (eng gefasst: „Ohne `--diff-all`“), Satz zur Meldung davor (eng gefasst: „Ohne `--diff-all` nennt die Meldung …“), Fehlercode-Zeile `PGR-E5004` (ergänzt um `--diff-all`), Zeile zur Art `Fehler` bei FATAL („nennt die Meldung ohne `--diff-all` diese Abweichung“, eng gefasst). „jede Interaktion mit einer Abweichung“ (Fortsetzung): unverändert, gilt mit und ohne Option, durch `…Fortsetzung` belegt. „ist das keine Abweichung“ (erwartetes Ende): unverändert, die Option ändert nicht, was zählt. Zeile 1655 (`PGR-E5001`): anderer Gegenstand, unverändert. Optionstabelle und `PGR-E2001`-Zeile: ergänzt. Abschnitte zum Ende der Verbindung (Art `Nachrichtenart`/`Fehler`): unverändert, die Art ist weiter die der Position; die Form der Zeile mit der Option steht im neuen Abschnitt.

**README (Liefer-Punkt 2):** von den fünf Stellen beschreibt nur die erste (Meldung als `PGR-E5004`) die Abweichung als einzelne; dort steht jetzt `--diff-all`. Zeilen 26 und 37 sowie 121 nennen `--compare-responses` ohne Aussage über die Zahl der Abweichungen: unverändert.

**Nicht belegt, nicht im Handbuch:** Parameterbeschreibung beim erweiterten Protokoll, Wertemenge und Rangfolge der Quellen von `--diff-all`, Verhalten mit FATAL in der Liste. Das Handbuch sagt dazu nichts.

**Nacharbeit nach dem Review (F-1 bis F-5):** F-1 enger gefasst („Ist `--diff-all` wirksam `true` und `--compare-responses` nicht …“, belegt durch die vier Fälle). F-2: Satz zur Form ohne Option aus dem Abschnitt gestrichen, Testkommentar enger. F-3: Fall `record Umgebungsvariable` ergänzt. F-4: Satz zum Ende der Verbindung auf den Wortlaut vor diesem Slice zurückgesetzt. F-5: `diff_all` in der Zeile `PGR-E2004` des Handbuchs, Fall `play Schlüssel ungültig` (`diff_all: 1`, `PGR-E2004`) ergänzt; `docs/user/abdeckung-e2e.md` in §3 genannt. Klasse (Merkmal: Handbuchsatz nennt eine Quelle oder Betriebsart, der Test eine Auswahl): übrige Sätze des Abschnitts mit Quelle oder Betriebsart gegen die Fälle gelesen, keine weitere Lücke. F-5 (b) und (c) bleiben unverändert (Handbuch sagt dazu nichts Zusätzliches).

**Keine Beobachtung angefallen** (Implementer; Closure-Notiz und Register schreibt der Planner).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `REPO` (gesamtes Repo, `harness/conventions.md`); die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (2026-10-23). Treffer: `zusage-im-kommentar-weiter-als-pruefung` (verkörpert in `AGENTS.md` §3.11; DoD 1 bindet jeden Satz an einen Lauf), `folge-slice-liefert-handbuch-teil-nicht` (dieser Slice ist die Adresse; sein §1 nimmt die Doku-Hälfte an, Eintrag im selben Commit wie der Geber). Keine Schwelle neu erreicht.

**Zählung** (Planner, 2026-10-23, `AGENTS.md` §3.13): **Liefer-Punkte: 2** (Handbuch, README). **Schichten: eine** (Nutzer-Doku, `docs/user/` und `README.md`; der E2E-Test ist Beleg der Doku, nicht Produkt-Code). Sendung aus `slice-play-feldvergleich` (Doku-Hälfte) mitgezählt: Grenzen 3 Punkte und 2 Schichten gehalten. **Nachgezählt vom Architect (2026-10-23):** Die Festlegung der Randformen in `LH-FA-24.a` *Alle Abweichungen* ändert die Zählung nicht (Liefer-Punkte 2, Schichten 1); der Handbuchtext folgt dem Inhalt der Meldung (Einzelheiten je Art, Gesamtzahl, „und N weitere“), nicht einem festgelegten Wortlaut, und belegt jeden Satz mit einem Lauf.

**Nachgezählt beim Anspruch (Planner, 2026-10-23):** Der Geber ist geschlossen; Stand und Stellen des Handbuchs (Optionstabelle, `PGR-E5004`, `PGR-E2001`, Ende der Verbindung, Exit-Code 5) in §1 bis §3 nachgezogen. Liefer-Punkte 2, Schichten 1 (Nutzer-Doku), unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
