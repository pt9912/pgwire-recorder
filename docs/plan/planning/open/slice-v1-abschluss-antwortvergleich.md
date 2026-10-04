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

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md), [ADR-0022](../../adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md)

**Berührte Spec-Stellen:** `LH-FA-24.a` · `LH-FA-20.a` · `SPEC-018` · `SPEC-027` · `SPEC-034` · `SPEC-041` · `ARC-002`

**Verantwortlich:** —

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

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Vergleich von Zeilenwerten und Toleranzregeln für Felder — Out-of-Scope von LH-FA-24.
- Einspielen selbst, Anmeldung und Fehlersemantik der Serverfehler — `slice-v1-abschluss-einspielen`; dieser Slice setzt es voraus.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen): Gegen eine Instanz mit gleicher Antwortstruktur endet das Einspielen mit Vergleich mit Erfolg, gegen eine abweichende mit `PGR-E5004` und Exit-Code 5; ein Serverfehler ohne Vorbild in der Aufzeichnung ist `PGR-E5004`, nicht `PGR-E4004`; eine unvollständige Aufzeichnung (auch ohne Server-Nachrichten, mit Fehlerantwort nach dem Präfix) wird verglichen, soweit sie reicht (Abnahmeszenario 16; Integrationstest).
- [ ] Jede Art der Abweichung (Nachrichtenart, Spaltenbeschreibung, Befehl ohne Zahlen, Fehler, Transaktionsstatus) wird erkannt und mit Session, Sequenznummer und Art gemeldet, Unterschiede nur in Zeilenwerten, Zeilenzahlen und Hinweisen nicht (Test, einfach und Extended, auch `Flush`-Gruppen und eine Aufzeichnung ohne `ReadyForQuery`).
- [ ] `--continue-on-error` und `--allow-recorded-errors` wirken wie spezifiziert, `--continue-on-error` läuft nach einer Abweichung weiter, der Exit-Code am Ende ist 5, ein Verbindungsfehler beendet immer mit 4; ein Fehler mit gleichem SQLSTATE wie aufgezeichnet gilt als erwartet und bricht nicht ab, ein anderer SQLSTATE ist `PGR-E5004`; `--allow-recorded-errors` ist mit Vergleich ohne Wirkung; ein Abbruchsignal endet nach einer Abweichung mit 5, ein Verbindungsfehler nach einer Abweichung mit 4; am Ende des Laufs steht die Zusammenfassung (eingespielt, verglichen, abweichend); ohne `--compare-responses` findet kein Vergleich statt (Test).
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
| `internal/adapters/driving/cli` | update | Option `--compare-responses`, Abbildung auf Exit-Code 5 |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-24 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` ist `done`, und [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Normalisierung der Antworten erweist sich als größer als geplant — zurück zur Zerlegung.
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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
