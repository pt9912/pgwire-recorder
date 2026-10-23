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

- [ ] **Liefer-Punkt 1 — Handbuch** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): der Abschnitt zum Vergleich der Antworten in `docs/user/benutzerhandbuch.md` (der Satz „Je Interaktion meldet `play` die erste Abweichung“ und jede Stelle, die der neuen Meldung widerspräche: Optionstabelle, Fehlercodes `PGR-E5004` und `PGR-E2001`, Abschnitt zum Ende der Verbindung, Exit-Code 5) nennt die Option, ihre drei Quellen, die Form der Meldung und die Obergrenze; jeder Satz belegt durch einen Lauf gegen das gebaute Binary (Datei, nicht gelesen gegen den Code), keine Chronik, kein Verweis auf Spezifikation, ADR, Slice oder Welle.
- [ ] **Liefer-Punkt 2 — README** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): die fünf Stellen von `--compare-responses` in `README.md` sind durchgesehen; wo sie die erste Abweichung als einzige beschreiben, ist die Option genannt, sonst bleibt der Text, mit Grund in §7.
- [ ] `make gates` am Endstand: grün.
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

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-play-feldvergleich` liegt in `done/` (erfüllt am Stand `db5d12ca`; der Wortlaut des Binaries steht fest). WIP-Limit 1.

**Rückführungen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik: `AGENTS.md` §3.14; Pfade des Plans: `docs/user README.md test`); Basis: der Commit des Übergangs nach `in-progress/`, Hash beim Übergang eintragen: `<Hash>` (Zweig `impl-diff-all-doku`).
- `in-progress` → `open`: das Binary zeigt ein Verhalten, das die Spezifikation anders sagt — Rückgabe an den Architect, Slice wartet nicht in `in-progress/`.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben; die Welle bzw. ein Release, das die Option enthält, entsteht erst nach diesem Slice.

## 6. Risiken und offene Punkte

- Ein Handbuchsatz sagt mehr zu, als das Binary zeigt (Wortlaut der Meldung, Obergrenze) — **Ausgang:** offen bis Closure (jeder Satz mit Lauf belegt, DoD 1).
- Die Grenze 20 zählt Spalten (ein Eintrag mit Name und Typ-OID), nicht Felder; ein Handbuchsatz, der „Felder“ oder „Einträge der Liste“ zählt, sagt es falsch — **Ausgang:** offen bis Closure (ein Lauf mit 20 und 21 Spalten mit beiden Feldern belegt den Satz, Geber `slice-play-feldvergleich` §7).
- Das Handbuch hinkt dem Binary hinterher, solange `slice-play-feldvergleich` ohne diesen Slice in `done/` liegt — **Ausgang:** offen bis Closure (Release erst nach beiden, `AGENTS.md` §3.11).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `REPO` (gesamtes Repo, `harness/conventions.md`); die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (2026-10-23). Treffer: `zusage-im-kommentar-weiter-als-pruefung` (verkörpert in `AGENTS.md` §3.11; DoD 1 bindet jeden Satz an einen Lauf), `folge-slice-liefert-handbuch-teil-nicht` (dieser Slice ist die Adresse; sein §1 nimmt die Doku-Hälfte an, Eintrag im selben Commit wie der Geber). Keine Schwelle neu erreicht.

**Zählung** (Planner, 2026-10-23, `AGENTS.md` §3.13): **Liefer-Punkte: 2** (Handbuch, README). **Schichten: eine** (Nutzer-Doku, `docs/user/` und `README.md`; der E2E-Test ist Beleg der Doku, nicht Produkt-Code). Sendung aus `slice-play-feldvergleich` (Doku-Hälfte) mitgezählt: Grenzen 3 Punkte und 2 Schichten gehalten. **Nachgezählt vom Architect (2026-10-23):** Die Festlegung der Randformen in `LH-FA-24.a` *Alle Abweichungen* ändert die Zählung nicht (Liefer-Punkte 2, Schichten 1); der Handbuchtext folgt dem Inhalt der Meldung (Einzelheiten je Art, Gesamtzahl, „und N weitere“), nicht einem festgelegten Wortlaut, und belegt jeden Satz mit einem Lauf.

**Nachgezählt beim Anspruch (Planner, 2026-10-23):** Der Geber ist geschlossen; Stand und Stellen des Handbuchs (Optionstabelle, `PGR-E5004`, `PGR-E2001`, Ende der Verbindung, Exit-Code 5) in §1 bis §3 nachgezogen. Liefer-Punkte 2, Schichten 1 (Nutzer-Doku), unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
