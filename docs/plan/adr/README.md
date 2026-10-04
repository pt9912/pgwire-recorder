# ADR-Index — pgwire-recorder

> **Derivativ:** Quelle der Wahrheit sind die ADR-Dateien; dieser Index ist eine Bequemlichkeits-Sicht — bei jedem neuen/akzeptierten ADR mitziehen.

| ID | Titel | Status | Bezug |
|---|---|---|---|
| [0001](0001-hexagonale-architektur.md) | Hexagonale Architektur | Accepted | [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus) |
| [0002](0002-driving-driven-terminologie.md) | Driving/Driven für Adapter, Inbound/Outbound für Ports | Accepted | — |
| [0003](0003-pgwire-server-ist-driving-adapter.md) | PGWire Server ist Driving Adapter | Accepted | [`LH-FA-04`](../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients) |
| [0004](0004-postgresql-upstream-ist-driven-adapter.md) | PostgreSQL Upstream ist Driven Adapter | Accepted | [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus) |
| [0005](0005-recording-store-ist-driven-adapter.md) | Recording Store ist Driven Adapter | Accepted | [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings) |
| [0006](0006-kanonisches-domain-model.md) | Kanonisches Domain Model | Accepted | [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten) |
| [0007](0007-strict-replay.md) | Strict Replay | Accepted | [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| [0008](0008-kein-sql-parser-in-v1.md) | Kein SQL-Parser in v1 | Accepted | [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| [0009](0009-implementierungssprache-go.md) | Implementierungssprache Go | Accepted | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung) |
| [0010](0010-verwendung-von-pgproto3.md) | Verwendung von pgproto3 | Accepted | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| [0011](0011-meldungscodes-praefix-pgr.md) | Meldungscodes mit Präfix PGR | Accepted | [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler) |
| [0012](0012-extended-query-gruppen.md) | Extended Query als Gruppen aus Client- und Server-Nachrichten | Accepted | [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) |
| [0013](0013-sqlite-recording-backend.md) | SQLite als zweites Aufzeichnungsformat | Accepted | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| [0014](0014-konfigurationsdatei.md) | Konfigurationsdatei mit benannten Verbindungen | Accepted | [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| [0015](0015-uhr-port.md) | Zeit über einen Uhr-Port | Accepted | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| [0016](0016-einspielen-anmeldung-und-tls.md) | Einspielen: Anmeldung und TLS als Client | Accepted | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| [0017](0017-einspielen-sequenziell-und-fehlersemantik.md) | Einspielen: sequenziell, Antworten verworfen, Fehlersemantik | Accepted | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| [0018](0018-tls-zum-client.md) | TLS-Terminierung zum Client | Accepted | [`LH-FA-23`](../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client) |
| [0019](0019-eigene-zertifizierungsstelle-beim-einspielen.md) | Eigene Zertifizierungsstelle beim Einspielen | Accepted | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| [0020](0020-antwortvergleich-beim-einspielen.md) | Antwortvergleich beim Einspielen | Accepted | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| [0021](0021-vergleichsregeln-beim-einspielen.md) | Vergleichsregeln beim Einspielen | Proposed | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |

## Konventionen

- ADRs sind nach `Accepted` **immutable** (siehe Baseline-Regelwerk `modul-04-adrs.md`).
- Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN`.
- Bei `Accepted`: diesen Index aktualisieren (Status, Datum).
- Jede ADR deklariert im `**Schärft:**`-Feld *aufwärts*, welche Spec-Stelle
  sie verbindlich macht (Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)) —
  als Kennung (`SPEC-*`, `ARC-*`, `<PREFIX>-FA-*.<Buchstabe>`), ersatzweise als
  Abschnitt, wo die Sektion keine Kennungen vergibt.
  Prozess-ADRs ohne Spec-Stratum tragen `—`.
