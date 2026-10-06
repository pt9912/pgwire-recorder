# ADR-0033: Gate-Nachweise und geteilte Messung in der Abdeckung

**Status:** Proposed

**Datum:** 2026-10-06

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md); die Qualitätsanforderung an die Prüfbarkeit des Quellcodes kommt mit ihrer Kennung hinzu, sobald sie im Lastenheft steht

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft bekommt eine Qualitätsanforderung an die Prüfbarkeit des Quellcodes mit drei Messmethoden, die je ein Gate misst: statische Analyse, Anweisungsabdeckung der Unit-Tests und Lage der Unit-Tests. Eine vierte Messmethode verlangt, dass die drei Prüfungen ohne Eingaben in der Prüfumgebung laufen und jede nachweislich rot wird, wenn ihre Bedingung verletzt ist. Die Nachweise sind die Gegenproben der Gates; sie sind Bash-Skripte und tragen keine Go-Tests.

[ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md) kennt Deklarationen nur an Tests und für eine Qualitätsanforderung genau einen Pfad *Messung*. Damit hätte die Anforderung keinen Nachweis, oder der erste Gate-Nachweis zählte sie als vollständig, obwohl zwei Messmethoden fehlen. Der Re-Evaluierungs-Trigger von [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md) (eine Nachweisart, die keine Tests trägt) ist damit erreicht.

## Entscheidung

Wir wählen **die Nachweisart Gate mit geteilter Messung**:

- Eine Gegenprobe, die an der Gate-Kette hängt, kann im Kopf eine Abdeckungs-Deklaration in derselben Form wie ein Test tragen. Ihre Nachweisart ist *Gate*; sie belegt nur Qualitätsanforderungen und Randbedingungen, kein Produktverhalten.
- Hat eine Messmethode mehrere Teile, die verschiedene Nachweise belegen, deklariert jeder Nachweis seinen Teil mit Nummer und Gesamtzahl. Vollständig ist die Anforderung, wenn alle Teile belegt sind. Der ungeteilte Pfad *Messung* bleibt für Anforderungen mit einem Nachweis unverändert.
- Ein Teil der Messmethode, der nur sagt, wie die übrigen Nachweise laufen und rot werden, ist kein eigener Pfad: Die Nachweisart Gate erfüllt ihn, weil eine Gegenprobe ohne Eingaben in der Gate-Kette läuft und den Verstoß rot zeigt.

[ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md) gilt sonst unverändert: Deklaration je Anforderung und Pfad, Tabellen aus den Deklarationen, die RTM liest nur vollständig belegte Anforderungen. Form der Deklaration, Pfad-Schreibweise und Fehlformen stehen im Kopf des Abdeckungs-Skripts.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts ändern | kein Aufwand | die Anforderung bleibt in jeder Tabelle unbelegt, obwohl ihre Gates laufen |
| B — Deklaration an einem Go-Test | kein neues Werkzeug | kein Go-Test prüft eine Lint-Regel oder eine Abdeckungsschwelle; die Deklaration sagte zu, was ihr Test nicht prüft (`AGENTS.md` §3.11) |
| C — Nachweisart Gate mit ungeteiltem Pfad *Messung* | kleinste Erweiterung | der erste Gate-Nachweis zählt die Anforderung als vollständig (Teilmessung) |
| D — eine Deklaration erst mit dem letzten Gate | ungeteilter Pfad genügt | ein Skript deklariert, was zwei andere prüfen; bis dahin ist nichts sichtbar |
| E — Unterkennungen je Messmethode im Lastenheft | feinste Bindung | Lastenheft und Kennungs-Schema ändern sich für eine Teilung, die die Pfade schon tragen (in [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md) als Option B verworfen) |
| **F — Nachweisart Gate mit geteilter Messung** | jeder Nachweis deklariert nur, was er prüft; Teilabdeckung bleibt sichtbar; bestehende Deklarationen bleiben gültig | die Gesamtzahl der Teile steht in der Deklaration, nicht im Lastenheft; eine falsche Zahl fängt nur das Review |

## Konsequenzen

- Positiv: Die Anforderung erscheint mit jedem Gate-Nachweis teilweise und mit dem letzten vollständig; die RTM zählt sie erst dann.
- Negativ: Die Gesamtzahl der Teile und die Zuordnung eines Teils zur Messmethode sind Urteil dessen, der deklariert, wie bei [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md); ob eine Gegenprobe an der Gate-Kette hängt, prüft das Abdeckungs-Skript nicht, das bleibt Review.
- Folgepflicht: Ein eigener Harness-Slice erweitert das Abdeckungs-Skript und seine Gegenprobe und trägt die ersten Gate-Deklarationen; `harness/README.md` §Sensors nennt die Nachweisart beim Vertrag von `make abdeckung-check`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Skript | Tabellen entsprechen den Deklarationen an Tests und Gegenproben; eine geteilte Messung ist vollständig, wenn alle Teile belegt sind | `make abdeckung-check` |
| Gegenprobe | Fehlformen der Gate-Deklaration und der geteilten Messung werden abgelehnt | `make abdeckung-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn eine Messmethode einen Nachweis braucht, der weder Test noch Gegenprobe ist, oder wenn das Review eine falsche Gesamtzahl der Teile findet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Proposed | slice-lastenheft-pruefbarkeit |
