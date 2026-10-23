# Verifikationsbericht: slice-play-feldvergleich — 2026-10-23

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gegen Plan und DoD?

**Gegenstand:** Worktree `/Development/pgwire-recorder-impl-df`, Zweig `impl-diff-all`, Stand **`a78065dc`** (Basis der Metrik `bf68b876`). Lieferung: `d195a58f` (Code, Tests), `e2ecc02f` (Klarstellung der Spezifikation), `9643db95` (Grenzfall F-1), `950f95a0` (Plan nachgezogen), Review `9b91bfa9`.

**Maßstab:** Plan (Kopf, §1 bis §8), `spec/spezifikation.md` `LH-FA-24.a` *Alle Abweichungen*, Optionstabellen `LH-FA-17.a` und `LH-FA-20.a`, `AGENTS.md` §3.9 bis §3.14. Eingang: [Review-Report](2026-10-23-review-slice-play-feldvergleich.md) (0 HIGH, 0 MEDIUM, 3 LOW, 1 INFO).

**Grenze:** kein `make gates` (am Endstand `950f95a0` grün laut Implementer, kein Anlass), keine Gegenprobe berührt (Diff unter `harness/`, `tools/`, `Makefile`, `.claude/`, `.harness/`: 0 Dateien); kein Push.

## Ergebnis in einem Satz

Beide Liefer-Punkte sind am Endstand **erfüllt**: Tests der drei Pakete grün, fünf Mutanten (darunter „Grenze zählt Felder“) aus dem richtigen Grund rot, Größe 1029 gegen Marke 1500; **offen** sind `make gates` (Beleg des Implementers, nicht selbst gefahren), die Closure-Pflichten und die Entscheidung des Planners zu F-3; kein HIGH, kein MEDIUM, ein LOW (V-1).

## Selbst nachgefahren

| # | Lauf | Ergebnis |
|---|---|---|
| 1 | `go test ./internal/hexagon/services/ ./internal/adapters/driving/cli/ ./internal/bootstrap/` im Image `golang:1.27.3`, Bind-Mount, Modul-Cache im Scratchpad (Download der Module mit Netz) | alle drei `ok` |
| 2 | Mutant a (F-1): `spaltenEinzelheiten` macht je Feld einen Eintrag (Grenze zählt Felder), Kopie per `git archive` | rot: `TestVergleichAlleObergrenzeSpaltenMitBeidenFeldern`, `TestVergleichAlle` |
| 3 | Mutant b: `abweichung` ignoriert die Option (`if false && …DiffAll`) | rot: `TestVergleichAlle`, `…ObergrenzePositionen`, `…ObergrenzeEinzelne`, `…Signal`, `…Zusammenspiel`, `…EndeDesServers`, `…FlushEnde`, `…Stelle` u. a. |
| 4 | Mutant c: `maxPositionen = 21` | rot: `TestVergleichAlleObergrenzePositionen` |
| 5 | Mutant d: Ende des Servers `u.position > ende` statt `>=` | rot: `TestVergleichAlleEndeDesServers`, `TestVergleichEndeAbweichung` |
| 6 | Mutant e (CLI): Kombination ohne `--compare-responses` entfällt | rot: `TestPlayDiffAllOhneVergleich` |
| 7 | Marke nach `AGENTS.md` §3.14, Basis `bf68b876`, Pfade des Plans | **1029** (Plan: 996 am Stand `d195a58f`); Marke 1500 nicht erreicht |
| 8 | `make abdeckung-check`, `make kopf-check`, `make docs-check` | je Exit 0 (`1411 Datei(en) geprüft, 0 Befund(e)`) |
| 9 | Aufräumen | kein Container mit Präfix `pgr-vdf-`; Worktree vor diesem Bericht sauber |

## DoD-Stand (§2)

| Punkt | Verdikt |
|---|---|
| **Liefer-Punkt 1 — Vergleich sammelt alle Abweichungen** | **erfüllt.** Code `alle`/`gewoehnlich`/`endeFehler`/`listeText`/`begrenze` folgt `LH-FA-24.a` *Alle Abweichungen* (Zuordnung Position für Position, Obergrenze 20 Positionen und 20 Einzelne, Ende des Servers als letzte Position, ohne Option `erste` = `alle[0]`, Meldung unverändert). Mutanten a bis d rot aus dem richtigen Grund; die Tabelle in §7 nennt je Zusage Mutation und roten Test, Stichprobe (vier Zeilen) bestätigt. Keine grünen Mutanten, nichts als äquivalent geführt |
| **Liefer-Punkt 2 — die Option** | **erfüllt.** Option in der Tabelle von `play` (Wahrheitswert, Standard `false`), Kombinationsprüfung auf zusammengeführten Werten (Mutant e rot), `record`/`replay` kennen sie nicht (`{"diff-all", p}`), globale Hilfe ohne, Hilfe von `play` mit der Option; E2E `TestE2EPlayVergleichAlle` (Option, Umgebung, Datei) und Abdeckungs-Deklarationen vorhanden; `abdeckung-check` grün; kein neuer Meldungscode, Bootstrap nur Durchreichung (`TestRunPlayDiffAll`) |
| `make gates` am Endstand | **nicht selbst gefahren.** Beleg des Implementers in §7 (`950f95a0`); seit dort ändert `a78065dc` nur §7 des Plans |
| Review-Report liegt vor | **erfüllt** (anderer Kontext, 0 HIGH, 0 MEDIUM); F-1, F-2, F-4 sind am Endstand behoben (Grenzfall und Mutation vorhanden, Rückgabe-Absatz in §7 trägt „entschieden“, Tabellenzeile in §7 hat vier Zellen); **F-3 offen** (siehe unten) |
| Closure-Notiz · Register · Risiken-Ausgänge · Paarungen | **offen (Planner)**; §7 trägt „Closure-Notiz: wird bei der Closure gefüllt“, §6 Risiken alle „offen bis Closure“ |

## Spezifikation und Entscheidungen

Diff `spec/spezifikation.md` gegen `main`: eine Klammer in *Einzelheiten je Art* plus Zeile im Änderungsprotokoll (`e2ecc02f`); der Code folgt ihr (`Spalte %d:` ein Eintrag mit Name zuerst, `begrenze` zählt Einträge je Spalte). Keine ADR nötig oder berührt; Entscheidungen des Nutzers vom 2026-10-23 (Name, Variable, Schlüssel, Form) sind im Code wiederzufinden. Handbuch und `README.md` unberührt, Folge-Slice `slice-play-feldvergleich-doku` liegt in `open/`.

## Findings

| ID | Kategorie | Befund |
|---|---|---|
| V-1 | LOW | §7 des Plans hält die Zahl 996 für die Marke am Stand `d195a58f` fest; am Endstand misst der Befehl 1029. Beleg-Zahl in §7 ist dadurch überholt (kein Einfluss auf die Entscheidung, Marke 1500 weit). Beim Closure die Endzahl eintragen oder die Zahl streichen (§3.14: §7 trägt keine Messzeile) |
| V-2 | INFO | Der Nachweis `make gates` am Endstand ruht auf dem Beleg des Implementers; wer Wert auf einen eigenen Lauf legt, fährt ihn im Hauptbaum nach dem Merge |
| F-3 (Review, offen) | LOW | siehe Einschätzung |

## Einschätzung zu F-3 (an den Planner, keine Entscheidung)

Lage: §3 und §4 sagen wörtlich, eine Änderung der Spezifikation **während** des Slice zähle als dritte Schicht, dann greife der vorgemerkte Schnitt (Option in eigenen Slice). `e2ecc02f` ist eine solche Änderung; §8 nennt das und überlässt die Entscheidung dem Planner.

Für „Schnitt greift nicht“: Die Änderung ist eine Klammer (3 Zeilen), sie hält fest, was der Code schon tat, Rückgabe-Weg und Entscheidung des Architect sind regulär gelaufen. Die Größe liegt bei 1029 gegen 1500, zwei Liefer-Punkte, Review in einer Sitzung ohne HIGH/MEDIUM. Ein Schnitt jetzt kostete Plan, Review, Verifikation und Gate-Lauf für einen Diff, der nach Code und Tests schon geliefert ist, und trennte Option und Vergleich, die einzeln wenig tragen.

Für „Schnitt greift“: Die Regel ist wörtlich formuliert und `AGENTS.md` §3.13 zählt Schichten ohne Ausnahme für Klarstellungen; eine stille Ausnahme schwächt den vorgemerkten Schnitt als Instrument, und der Zähler `BEO-REPO/slice-waechst-durch-uebernahmen` kennt die Klasse.

Meine Lesart: Der Zweck der Regel (Größe, Prüfbarkeit in einer Sitzung) ist belegt erfüllt; der Planner kann die Regel in §3/§4 für diesen Fall ausdrücklich als „nicht ausgelöst“ entscheiden, mit Begründung (Klarstellung ohne neuen Inhalt, Code unverändert) und Datum in §7 und §8, und §3/§4 so schärfen, dass Klarstellungen, die den Code nicht ändern, ausgenommen sind, statt die Ausnahme still zu nehmen. Das ist Urteil des Planners, nicht des Verifiers.

## Verdikt

DoD-Liefer-Punkte erfüllt; Review-Findings F-1, F-2, F-4 behoben, F-3 beim Planner; keine Merge-blockierende Abweichung. Offen bis Closure: Closure-Notiz mit Lerneintrag, Register, Ausgänge der vier Risiken aus §6, Paarungen.
