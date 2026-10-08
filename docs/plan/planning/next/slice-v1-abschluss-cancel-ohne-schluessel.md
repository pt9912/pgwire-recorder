# Slice slice-v1-abschluss-cancel-ohne-schluessel: CancelRequest ohne Schlüssel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)

**Berührte Spec-Stellen:** `LH-FA-05.e` · `ARC-006`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein `CancelRequest`, den ein Client ohne Abbruchschlüssel sendet (pgx bei Kontext-Abbruch, 12 Byte, weil der Recorder `BackendKeyData` nicht weitergibt), wird wie jeder `CancelRequest` behandelt: Verbindung geschlossen, Warnung `PGR-W3001`, kein Verbindungsfehler (`LH-FA-05.e`). Heute lehnt die PGWire-Bibliothek ihn als „cancel request too short“ ab, und der Recorder meldet `PGR-E6001` mit Klasse 6 (Review-Finding F-310 aus `docs/reviews/2026-10-05-review-slice-extended-query-record.md`, Sonde P3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Weitergabe von `BackendKeyData` und echtes Abbrechen am Upstream — `LH-FA-05.e` legt fest, dass ein `CancelRequest` nicht weitergeleitet wird; eine Änderung wäre eine neue Anforderung.
- Übriger Protokollrand — Gegenstand von `slice-v1-abschluss-protokollrand`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein `CancelRequest` ohne Schlüssel (Länge 12) schließt die Verbindung mit `PGR-W3001` ohne Verbindungsfehler, in `record` und `replay` (Unit-Test mit 12 Byte).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): pgx mit Kontext-Abbruch über `record` lässt den Lauf mit Exit-Code 0 enden, wenn sonst kein Fehler auftrat (Integrationstest).
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
| `internal/adapters/driving/pgwire` (`startup`) | update | Startcode `CancelRequest` vor der Bibliothek erkennen, Länge 12 und 16 annehmen |
| `internal/adapters/driving/pgwire/server_test.go` | update | Negative nach `LH-FA-05`: 12 Byte |
| `test/integration` | update | pgx mit Kontext-Abbruch, Exit-Code |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Behandlung verlangt die Weitergabe von `BackendKeyData`.
- `in-progress` → `open` (blockiert — Carveout?): die PGWire-Bibliothek lässt den Startcode nicht vor ihrer Längenprüfung lesen.

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

- Andere Clients senden einen `CancelRequest` in weiteren Längen — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — der Fall entstand, weil der Test nur den `CancelRequest` mit Schlüssel prüfte.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
