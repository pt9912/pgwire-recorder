# ADR-0031: Lebendprüfungen im Replay außerhalb der Reihe

**Status:** Proposed

**Datum:** 2026-10-05

**Autor:** pt9912

**Bezug:** [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-QA-02`](../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [ADR-0007](0007-strict-replay.md), [ADR-0008](0008-kein-sql-parser-in-v1.md), [ADR-0012](0012-extended-query-gruppen.md)

**Schärft:** [`LH-FA-09.a`](../../../spec/spezifikation.md#lh-fa-09a--strict-sequential-matching), [`LH-FA-18.a`](../../../spec/spezifikation.md#lh-fa-18a--extended-query-ablauf-aufzeichnung-matching), [`LH-FA-12.a`](../../../spec/spezifikation.md#lh-fa-12a--reihenfolge-sessions-und-parallelität), [`LH-FA-03.b`](../../../spec/spezifikation.md#lh-fa-03b--nicht-verbrauchte-interaktionen), [`SPEC-011`](../../../spec/spezifikation.md#3-defaults-und-konstanten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`pgxpool` und `database/sql` mit pgx prüfen eine Verbindung, die länger als 1 s geruht hat, mit der einfachen Anfrage `-- ping`. Ob eine solche Lebendprüfung in der Aufzeichnung steht und ob sie in der Wiedergabe kommt, hängt damit am Zeitverhalten; das strenge Replay ([ADR-0007](0007-strict-replay.md)) meldet jede Verschiebung als Abweichung (`PGR-E5001`), in beiden Richtungen ([Validierung des Extended-Replay](../../reviews/2026-10-05-validierung-slice-extended-query-replay.md), Frage 2, Sonden S5 bis S7). Die Abhilfe am Treiber (`ShouldPing`) wäre ein Eingriff über Host und Port hinaus, gegen [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) und [`LH-QA-02`](../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung). Der Nutzer hat entschieden, solche Anfragen zu tolerieren, statt die Zusage zu verengen. Der Re-Evaluierungs-Trigger von [ADR-0007](0007-strict-replay.md) („Anwender scheitern regelmäßig an der Strenge“) ist damit für diese eine Klasse von Anfragen eingetreten.

PostgreSQL beantwortet eine einfache Anfrage, deren Text keine Anweisung enthält, mit `EmptyQueryResponse` und `ReadyForQuery`, ohne den Zustand der Session zu ändern, auch in einer Transaktion und nach einem Fehler in ihr. Die Antwort hängt also nur am Transaktionsstatus, nicht an Daten. Bestehende Aufzeichnungen enthalten Lebendprüfungen als gewöhnliche Interaktionen (S6).

Annahmen: Verbreitete Treiber senden Lebendprüfungen als einfache Anfrage aus Kommentar oder Leerraum, nicht über das Extended Query Protocol. Das Recording-Format bleibt unverändert.

## Entscheidung

Wir wählen: **Das Replay behandelt eine Lebendprüfung als Nicht-Interaktion, eingehend wie aufgezeichnet; der Record bleibt unverändert.** Diese ADR ergänzt [ADR-0007](0007-strict-replay.md) und ersetzt sie nicht: Für jede andere Anfrage gilt strict sequential unverändert.

- **Erkennung.** Eine Lebendprüfung ist eine einfache Anfrage (`Query`), deren Text nach den lexikalischen Regeln von PostgreSQL nur aus Leerraum, Zeilenkommentaren und geschlossenen, auch verschachtelten Blockkommentaren besteht; der leere Text zählt dazu. Was PostgreSQL mit einem Fehler beantwortete (etwa ein nicht geschlossener Blockkommentar) oder was ein anderes Zeichen enthält, auch ein einzelnes `;`, ist keine. Die Klasse ist rein lexikalisch und gilt für jede solche Anfrage, gleich ob sie als Lebendprüfung gemeint ist: Auf Protokollebene ist der Unterschied nicht sichtbar, und PostgreSQL antwortet auf beide gleich.
- **Eingehend.** Kommt eine Lebendprüfung zwischen zwei Interaktionen, auch vor der ersten und nach der letzten, beantwortet das Replay sie außerhalb der Reihe mit `EmptyQueryResponse` und `ReadyForQuery`; der Cursor bleibt stehen, und sie löst keine Session-Zuordnung aus. Das `ReadyForQuery` trägt den Transaktionsstatus des letzten `ReadyForQuery`, das das Replay auf dieser Verbindung gesendet hat, also den Stand nach der letzten beantworteten Interaktion, nach dem Handshake `I`.
- **Aufgezeichnet.** Das Replay liest eine Aufzeichnung so, als stünden die Lebendprüfungen nicht darin: für Cursor, Session-Zuordnung und nicht verbrauchte Interaktionen. Ihre aufgezeichneten Antworten verwendet es nicht. Diagnosen nennen weiter die aufgezeichnete Interaktionsnummer.
- **Grenze.** Mitten in einer Extended-Interaktion, nach deren erster Nachricht und vor deren `Sync`, ist auch eine Lebendprüfung eine Abweichung (`PGR-E5001`), ebenso außerhalb davon jede andere Nachricht der falschen Protokollart am Cursor, etwa eine einfache Anfrage mit Anweisung, wo eine Extended-Interaktion erwartet wird. Lebendprüfungen über das Extended Query Protocol (`Parse` mit leerem Text) toleriert diese ADR nicht.

Die Einzelregeln (Zeichen des Leerraums, Kommentarformen, Zuordnung und Zählung) stehen in der Spezifikation; bei Abweichung gilt sie.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, Zusage verengen (`ShouldPing` abschalten, im Lastenheft und Handbuch genannt) | [ADR-0007](0007-strict-replay.md) bleibt ohne Ausnahme; kein Code | bricht „nur Host und Port“ für den häufigsten Einsatz von pgx; vom Nutzer verworfen |
| B — Der Record zeichnet Lebendprüfungen nicht auf, das Replay beantwortet eingehende außer der Reihe | Aufzeichnungen werden vom Zeitverhalten unabhängig | vorhandene Aufzeichnungen mit Lebendprüfungen (S6) bleiben rot, solange das Replay sie nicht überspringt; zweite Stelle im Record-Service; `play` sendet die Lebendprüfungen der Anwendung nicht mehr |
| C — Wie B, dazu eine Markierung aufgezeichneter Lebendprüfungen im Format | aufgezeichnete Lebendprüfungen ohne Erkennung beim Lesen kenntlich | Formatänderung in YAML und SQLite samt Version; dritte Schicht, Rückführung des Slice zur Zerlegung; die Erkennung beim Lesen ist ohnehin nötig |
| D — Toleranz für jede Anfrage, die PostgreSQL mit `EmptyQueryResponse` beantwortet (auch `;`) | deckt sich mit dem Verhalten des Servers | verlangt, Anweisungsgrenzen zu deuten, näher an einem SQL-Parser ([ADR-0008](0008-kein-sql-parser-in-v1.md)); kein beobachteter Treiber sendet `;` als Lebendprüfung |
| **E — Erkennung nur aus Leerraum und Kommentaren; Replay beantwortet eingehende und überspringt aufgezeichnete; Record unverändert (gewählt)** | keine Formatänderung; vorhandene Aufzeichnungen werden verträglich; eine Stelle im Code (Replay-Service); die Antwort ist genau die des Servers | Aufzeichnungen derselben Anwendung bleiben je nach Pausen verschieden (nur in Lebendprüfungen); eine fachlich gemeinte Kommentar-Anfrage wird nicht mehr streng geprüft |

## Konsequenzen

- Positiv: `pgxpool` und `database/sql` laufen im Replay unabhängig von Pausen beim Aufzeichnen und Wiedergeben (S4 bis S7); vorhandene Aufzeichnungen bleiben ohne Umwandlung verwendbar. Der Record und `play` bleiben unverändert und geben weiter wieder, was die Anwendung gesendet hat.
- Negativ: Eine Kommentar- oder Leer-Anfrage, die außer der Reihe kommt oder fehlt, erkennt das Replay nicht mehr als Abweichung. Das ist ein akzeptiertes Negativ: Sie ändert den Zustand der Session nicht, und der Client erhält dieselbe Antwort wie von PostgreSQL. Die Antwort außerhalb der Reihe stammt nicht aus der Aufzeichnung; sie ist wie der Handshake im Replay (`LH-FA-05.b`) eine feste Protokollantwort, keine frei definierte Antwort im Sinne des Lastenhefts. Die Erkennung von Kommentaren ist lexikalisch und keine semantische Analyse; [ADR-0008](0008-kein-sql-parser-in-v1.md) bleibt unberührt.
- Folgepflicht: `LH-FA-09.a`, `LH-FA-18.a` §Replay, `LH-FA-12.a` und `LH-FA-03.b` nennen die Ausnahme und ihre Grenze; das Lastenheft nennt sie in [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), soweit der Nutzer das bestätigt; das Handbuch beschreibt das Verhalten. Unit-Tests decken jede Kommentarform, die Grenzfälle der Erkennung, den Transaktionsstatus und die Grenze mitten in einer Extended-Interaktion ab.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Unit-Tests am Replay-Service | Lebendprüfung zwischen Interaktionen bewegt den Cursor nicht und erhält `EmptyQueryResponse` mit dem Transaktionsstatus der letzten Antwort; aufgezeichnete Lebendprüfung ist kein Mismatch und nicht unverbraucht; keine Lebendprüfung (`;`, offener Blockkommentar, Anweisung nach Kommentar) und jede einfache Anfrage mitten in einer Extended-Interaktion bleiben `PGR-E5001` | `make test` |
| Integrationstest gegen PostgreSQL | `pgxpool` und `database/sql` mit und ohne Pausen über 1 s beim Aufzeichnen und Wiedergeben (Lagen S4 bis S7) ohne Abweichung | `make test-integration` |

## Re-Evaluierungs-Trigger

Wenn ein verbreiteter Treiber Lebendprüfungen anders sendet (über das Extended Query Protocol oder als Anweisung wie `SELECT 1`), wenn eine vom Zeitverhalten unabhängige Aufzeichnung als eigene Anforderung aufgenommen wird (dann Option B neu prüfen) oder wenn ein toleranter Matching-Modus als eigene Anforderung aufgenommen wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-05 | Proposed | [Validierungsbeleg](../../reviews/2026-10-05-validierung-slice-extended-query-replay.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
