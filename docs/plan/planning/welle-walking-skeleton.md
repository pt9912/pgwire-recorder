# Welle welle-walking-skeleton: Walking Skeleton Record und Replay

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M1.

**Verantwortlich:** pt9912. **Datum:** 2026-10-03.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Ein `SELECT 1;` läuft end-to-end durch das gesamte Hexagon: Im Record-Modus gegen eine reale PostgreSQL-Instanz entsteht ein Recording, und mit gestoppter PostgreSQL liefert der Replay-Modus dem Client dasselbe beobachtbare Ergebnis. Das ist [Abnahmeszenario 1](../../../spec/lastenheft.md#7-abnahmekriterien-für-v1) und [Abnahmeszenario 2](../../../spec/lastenheft.md#7-abnahmekriterien-für-v1) des Lastenhefts, begrenzt auf eine Query. Die Welle validiert die Architekturgrenzen, bevor weitere PGWire-Nachrichten implementiert werden.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Die Architekturentscheidungen [ADR-0001](../adr/0001-hexagonale-architektur.md), [ADR-0003](../adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0004](../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0005](../adr/0005-recording-store-ist-driven-adapter.md), [ADR-0006](../adr/0006-kanonisches-domain-model.md), [ADR-0009](../adr/0009-implementierungssprache-go.md) und [ADR-0010](../adr/0010-verwendung-von-pgproto3.md) sind `Accepted`; die übrigen ADRs aus dem Index sind mindestens `Proposed`.
- `make gates` ist auf dem Stand der Spec-Erstfassung grün.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- **Walking-Skeleton-Smoke** (das *Mehr* gegenüber den Slice-DoDs): `SELECT 1;` wird über den Recorder gegen eine reale PostgreSQL-Instanz aufgezeichnet, die Instanz wird gestoppt, und derselbe Client erhält im Replay dasselbe Ergebnis; Ablauf ohne interaktive Eingaben.
- `make gates` grün, einschließlich Architektur-Gate.
- Closure-Notiz in `welle-walking-skeleton-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-walking-skeleton-build-gates | Docker-only Build, Test und Architektur-Gate | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) |
| slice-walking-skeleton-record | Record-Pfad für `SELECT 1;` | [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings) |
| slice-walking-skeleton-replay | Replay-Pfad für `SELECT 1;` | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle welle-extended-query (baut auf Recorder und Replay-Pfad auf).
- Wird blockiert von: — (nach erfüllten Start-Triggern keine weitere Welle).
- Innerhalb der Welle: `slice-walking-skeleton-record` und `slice-walking-skeleton-replay` setzen `slice-walking-skeleton-build-gates` voraus; `slice-walking-skeleton-replay` nutzt das Recording aus `slice-walking-skeleton-record`.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Mismatch-Diagnose, Exit-Code-Zuordnung und Fehlerreplay — Gegenstand von welle-replay-semantik.
- Mehrere Client-Sessions und Parallelität — Gegenstand von welle-v1-abschluss.
- Signalbehandlung, Container-Image und Betriebsdokumentation — Gegenstand von welle-v1-abschluss.
- Extended Query Protocol — Gegenstand von welle-extended-query.
- TLS, COPY, CancelRequest — nicht Teil von v1.
- Entscheidungen zu den offenen Spec-Punkten (Exit-Code bei Signal, vorhandenes `--output`, PGWire-Versionen, Strictness-Schalter) — fallen vor der jeweiligen Folge-Welle an, nicht hier.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-walking-skeleton-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
