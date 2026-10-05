# Slice slice-replay-semantik-mismatch: Nicht verbrauchte Interaktionen als Fehler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-replay-semantik.

**Bezug:** [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)

**Berührte Spec-Stellen:** `LH-FA-03.b` · `LH-FA-13.b` · `LH-FA-10.a` · `SPEC-012` · `SPEC-018` · `ARC-002` · `ARC-006`

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

**Ziel:** Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions im Replay ein Fehler (`PGR-E5002`, Exit-Code 5) statt der Warnung `PGR-W2001`. Erkennung und Diagnose der Abweichung sind für beide Protokollvarianten geliefert (siehe unten).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Meldungscodes und Fehlertext-Kopf — `slice-replay-semantik-meldungscodes`.
- Diagnose einer Abweichung — geliefert (siehe „Bereits geliefert“); dieser Slice ändert sie nicht.
- Warnung bei nicht verbrauchten Interaktionen — geliefert von `slice-walking-skeleton-replay` (`PGR-W2001`).

**Bereits geliefert** von `slice-walking-skeleton-replay`: Mismatch einfacher Anfragen mit Diagnose (Session, erwartete Nummer, erwartete und empfangene Anfrage), `ErrorResponse` mit `PGR-E5001` und Exit-Code 5 beim Herunterfahren. Von `slice-extended-query-replay`: Mismatch der Extended-Nachrichten mit Diagnose nach `LH-FA-10.a` (Session, Interaktion, Gruppe, Nachricht, erwarteter und empfangener Nachrichtentyp, abweichendes Feld, SQL der erwarteten und der empfangenen Anweisung, ohne Parameterwerte; `TestReplayExtendedDiagnoseAnweisung`), auch für eine einfache Anfrage, wo eine Extended-Nachricht erwartet ist, und umgekehrt (`TestReplayExtendedAbweichung`, `TestReplayExtendedFalscheArt`, `TestE2EReplayExtendedAbweichung`). Dieser Slice ergänzt `--fail-on-unconsumed` (`PGR-E5002`). Titel und Bezug folgen diesem Rest: Die Folgepflicht aus [ADR-0007](../../adr/0007-strict-replay.md), die Diagnose der Abweichung, ist für beide Protokollvarianten geliefert; `--fail-on-unconsumed` schärft [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), nicht das Matching.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions `PGR-E5002` mit Exit-Code 5 statt der Warnung `PGR-W2001` (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Replay-Service) | update | nicht verbrauchte Interaktionen und Sessions als `PGR-E5002` bei gesetzter Option |
| `internal/adapters/driving/pgwire` | update | `PGR-E5002` beim Ende der Verbindung als Verbindungsfehler merken |
| `internal/adapters/driving/cli` | update | Option `--fail-on-unconsumed` und Umgebungsvariable |
| `test/integration` | update | Exit-Code 5 mit, Warnung ohne die Option |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-extended-query` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: `--fail-on-unconsumed` verlangt eine Änderung am Recording-Format — zurück zur Zerlegung.
- `in-progress` → `open`: Der Replay-Pfad des Walking Skeleton ist unvollständig — Carveout.


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

- Exit-Code erst beim Prozessende (Verbindungsfehler beenden nur die Verbindung) — **Ausgang:** offen bis Closure.

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
