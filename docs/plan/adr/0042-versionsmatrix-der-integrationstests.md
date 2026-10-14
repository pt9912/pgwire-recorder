# ADR-0042: Versionsmatrix der Integrationstests als Teilmenge gegen die jüngsten Hauptversionen

**Status:** Accepted

**Datum:** 2026-10-14

**Autor:** Architect

**Bezug:** [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-RB-02`](../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md)

**Schärft:** [`SPEC-053`](../../../harness/vertraege/SPEC-053-postgres-versionen.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0026](0026-build-und-test-im-multistage-dockerfile.md) lässt die Integrationstests gegen *ein* gepinntes PostgreSQL-Image laufen. Das Produkt soll gegen mehrere Hauptversionen aufzeichnen; ein Teil seines Verhaltens hängt am Server (was PostgreSQL als leere Anfrage behandelt, welche Anmeldung die Instanz annimmt). Eine Zusage über Versionen ohne Gate dafür ist eine Behauptung. Der Nutzer hat am 2026-10-05 entschieden: unterstützt werden die fünf jüngsten erschienenen Hauptversionen, die Liste wird beim Start gegen die Release-Liste von PostgreSQL bestätigt, und der Laufzeit-Zuwachs des Gates ist akzeptiert.

## Entscheidung

Wir wählen **Option C**: `make test-integration` bleibt unverändert die volle Suite gegen die Referenzversion. Ein weiteres Gate, `make test-postgres-versionen`, hängt an `make gates` und fährt gegen jede andere Version der Liste, je per Digest gepinnt, nur die Teilmenge der Tests, deren Aussage vom Server abhängt. Die Liste steht an einer Stelle im Fragment; ihre Einzelregeln (Form, Rot-Bedingungen, Teilmenge, Aufräumen) stehen im Vertrag [`SPEC-053`](../../../harness/vertraege/SPEC-053-postgres-versionen.md). Das Gate liest die Release-Liste nicht; die Bestätigung des Bereichs ist eine Handlung beim Pflegen der Liste.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — volle Suite gegen jede Version in `make gates` | belegt am meisten | Laufzeit wächst mit der Zahl der Versionen, für Tests ohne Serverabhängigkeit ohne Gewinn |
| B — Matrix als eigenes Ziel außerhalb von `make gates` | billig | ein Lauf, den niemand ruft, belegt nichts; eine Zusage stünde ohne Gate (`AGENTS.md` §3.6, §3.11) |
| **C — volle Suite gegen die Referenz, Teilmenge gegen die übrigen Versionen** | jede Zusage hat ein Gate; keine bestehende Prüfung fällt weg; Teilmenge erweiterbar | Laufzeit wächst um einen Serverstart und zwei Phasen je Version; die Teilmenge wählt Tests über ihre Namen und braucht eine Prüfung, dass sie lief |
| D — Gate liest die Release-Liste und leitet die Versionen ab | keine Pflege von Hand | Netz zur Laufzeit des Gates, Ergebnis hängt vom Tag des Laufs ab (`LH-QA-01`) |

## Konsequenzen

- Positiv: Was die Produktspezifikation über Versionen zusagt, belegt ein Gate; die Referenz bleibt unverändert.
- Negativ: Die Gate-Laufzeit wächst; die Liste der Versionen muss gepflegt werden, wenn eine neue Hauptversion erscheint.
- Akzeptiertes Negativ: Die Teilmenge belegt das Verhalten der übrigen Tests nicht für jede Version. Der Grund: Tests ohne Serverabhängigkeit zeigen gegen eine weitere Version nichts Neues; ein Lauf von Hand mit anderem `POSTGRES_IMAGE` bleibt möglich.
- Folgepflicht: `harness/README.md` §Sensors führt das Ziel mit Bindung an diese ADR; die Ausnahmeliste der schnellen Kette ([ADR-0038](0038-gate-kette-voll-und-schnell.md)) nennt es, weil es PostgreSQL-Container startet.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Docker, Integrationstests | die Teilmenge besteht gegen jede Version der Liste | `make test-postgres-versionen` |
| Gegenprobe | Eingabefehler, fehlender Test, Aufräumen und Aufrufarten des Treibers werden rot | `make test-postgres-versionen-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn die Laufzeit des Ziels die Kette dominiert, wenn eine Hauptversion erscheint (die Liste verschiebt sich) oder wenn ein serverabhängiger Abnahmetest hinzukommt, der nicht in der Teilmenge steht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-14 | Proposed | — |
| 2026-10-14 | Accepted (Entscheidung des Nutzers vom 2026-10-14) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
