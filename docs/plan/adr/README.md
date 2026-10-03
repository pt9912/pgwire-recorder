# ADR-Index — pgwire-recorder

> **Derivativ:** Quelle der Wahrheit sind die ADR-Dateien; dieser Index ist eine Bequemlichkeits-Sicht — bei jedem neuen/akzeptierten ADR mitziehen.

| ID | Titel | Status | Bezug |
|---|---|---|---|
| [0001](0001-hexagonale-architektur.md) | Hexagonale Architektur | Proposed | [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus) |
| [0002](0002-driving-driven-terminologie.md) | Driving/Driven für Adapter, Inbound/Outbound für Ports | Proposed | — |
| [0003](0003-pgwire-server-ist-driving-adapter.md) | PGWire Server ist Driving Adapter | Proposed | [`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients) |
| [0004](0004-postgresql-upstream-ist-driven-adapter.md) | PostgreSQL Upstream ist Driven Adapter | Proposed | [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus) |
| [0005](0005-recording-store-ist-driven-adapter.md) | Recording Store ist Driven Adapter | Proposed | [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings) |
| [0006](0006-kanonisches-domain-model.md) | Kanonisches Domain Model | Proposed | [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten) |
| [0007](0007-strict-replay.md) | Strict Replay | Proposed | [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| [0008](0008-kein-sql-parser-in-v1.md) | Kein SQL-Parser in v1 | Proposed | [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| [0009](0009-implementierungssprache-go.md) | Implementierungssprache Go | Proposed | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung) |
| [0010](0010-verwendung-von-pgproto3.md) | Verwendung von pgproto3 | Proposed | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| [0011](0011-meldungscodes-praefix-pgr.md) | Meldungscodes mit Präfix PGR | Proposed | [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler) |

## Konventionen

- ADRs sind nach `Accepted` **immutable** (siehe Baseline-Regelwerk `modul-04-adrs.md`).
- Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN`.
- Bei `Accepted`: diesen Index aktualisieren (Status, Datum).
- Jede ADR deklariert im `**Schärft:**`-Feld *aufwärts*, welche Spec-Stelle
  sie verbindlich macht (Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)) —
  als Kennung (`SPEC-*`, `ARC-*`, `<PREFIX>-FA-*.<Buchstabe>`), ersatzweise als
  Abschnitt, wo die Sektion keine Kennungen vergibt.
  Prozess-ADRs ohne Spec-Stratum tragen `—`.
