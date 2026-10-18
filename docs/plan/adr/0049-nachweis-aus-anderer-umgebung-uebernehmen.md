# ADR-0049: Der Nachweis des vollen Laufs darf aus einer anderen Umgebung übernommen werden, wenn die Dateiliste des Baums gleich ist

**Status:** Accepted

**Datum:** 2026-10-18

**Autor:** Architect; angenommen durch Entscheidung des Nutzers am 2026-10-18 (Option B).

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0038](0038-gate-kette-voll-und-schnell.md), [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md)

**Schärft:** [`SPEC-062`](../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Stop-Hook gibt frei, wenn die Dateiliste des Arbeitsbaums mit der Liste des letzten grünen `make gates` identisch ist ([ADR-0038](0038-gate-kette-voll-und-schnell.md)). Die Liste ist ein Nachweis je Arbeitsbaum, in `.harness/state/` (ignoriert). Die Arbeit läuft in Worktrees und im Klon eines Dev-Containers; der Hauptbaum bekommt den Stand per Pull. Danach ist seine Liste ungleich der alten, und der Hook sperrt die nächste Aktion dort, bis `make gates` im Hauptbaum lief (rund 4,5 Minuten, mit Docker), obwohl der Baum Byte für Byte der ist, der im Worktree geprüft wurde. Der Nutzer will das nicht: „Keine blockierenden Aktionen“ im Hauptbaum, und zum Weg: „Nachweis aus dem Container oder Worktree übernehmen — das hört sich pragmatisch an“ (2026-10-18).

Bisher gilt die Lesart: Der Hook glaubt den Nachweis, den `make gates` in dieser Umgebung schrieb. Ein Werkzeug, das einen Nachweis aus einer anderen Umgebung ablegt, ändert diese Lesart: Der Hook glaubt dann Bytes, deren Herkunft er nicht kennt. Das ist keine Schwelle, aber eine Zusage des vollen Nachweises (`AGENTS.md` §3.6 gilt dem Sinn nach), und sie gehört in eine ADR.

## Entscheidung

Wir wählen **Option B: ein Werkzeug `make nachweis-uebernehmen` legt den Nachweis einer Quelle im Arbeitsbaum ab, genau dann, wenn die frisch berechnete Dateiliste des Arbeitsbaums byte-gleich mit der Liste der Quelle ist, die Liste unbeschädigt ist und der Hash zur Liste gehört; sonst lehnt es ab und ändert nichts.**

Der Nachweis bleibt der eines vollen Laufs von `make gates`. Neu ist nur, dass seine Bytes den Weg über eine Umgebungsgrenze nehmen dürfen, und zwar auf einem Weg, der prüft, dass sie diesen Baum beschreiben. Der Hook bleibt unverändert und liest weiter nur die Liste. Die Regeln im Einzelnen (Bedingungen, Reihenfolge der Gründe, Ablehnung ohne Änderung, Aufrufarten, Grenzen) stehen in [`SPEC-062`](../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md).

Zwei Setzungen, die der Nutzer am 2026-10-18 gemacht hat (so vom Koordinator übermittelt, der Wortlaut liegt dem Architect nicht vor): Das Werkzeug übernimmt nicht den Stand der schnellen Kette und ändert die zweite Sperre des Hooks nicht (inerte Änderungen verlangen den Stand der schnellen Kette); die Ablehnungen, die die Probe zeigt, werden gemessen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — wie heute: Handarbeit (`cp` einer Liste, `make record-gates` einzeln) | nichts zu bauen | Nichts prüft Gleichheit, Unversehrtheit oder Zusammengehörigkeit; ein Tippfehler schreibt einen Nachweis, den der Hook glaubt |
| **B — Werkzeug mit Gleichheitsprüfung** | der Weg über die Umgebungsgrenze ist ein Aufruf mit Prüfung; ein Nachweis eines anderen Baums wird nie abgelegt; der Hook bleibt wie er ist | ein zweiter Schreiber; die Herkunft wird nicht geprüft (Konsequenzen) |
| C — Nachweis mit Herkunft (Commit, Stand des Images) und ein Hook, der sie liest | die Quelle wäre benannt | Eine Herkunft ohne Signatur ist von Hand ebenso zu schreiben; das Format der Liste und der Hook ändern sich; der Aufwand steht in keinem Verhältnis zum Gewinn |
| D — der Hook gibt inerte Änderungen ohne Stand der schnellen Kette frei, oder nach einem Pull ohne Nachweis | löst die zweite Sperre | Lockerung des Hooks (`AGENTS.md` §3.6), anderer Vorgang; vom Nutzer für diesen Schritt nicht gewählt |
| E — nur noch im Container arbeiten, der Hook läuft dort | kein Transport des Nachweises | löst den Hauptbaum nicht: Der Pull dorthin bleibt gesperrt; der Standardweg „Arbeit im Container“ ist ein eigener Slice und ergänzt B |
| F — nichts ändern | kein Aufwand | Jede Aktion im Hauptbaum nach einem Pull verlangt den vollen Lauf; „Keine blockierenden Aktionen“ ist verletzt |

## Konsequenzen

- Positiv: Ein Pull, dessen Baum im Worktree oder Container schon geprüft wurde, sperrt den Hauptbaum nicht mehr. Ein Nachweis eines anderen Baums wird nicht abgelegt, ein beschädigter nicht, ein halb geschriebener nicht.
- Negativ, offen geschrieben: Der Hook glaubt Bytes, die in einer anderen Umgebung entstanden sind. Das Werkzeug prüft Gleichheit mit dem Baum, nicht, dass die Quelle `make gates` fuhr, mit denselben Gates, im selben Stand des Images und mit denselben Werkzeugversionen. Ein Gate, das nur in einer anderen Umgebung rot wird, bleibt in der Quelle unbemerkt.
- Akzeptiertes Negativ: Jeder lokale Nachweis war schon von Hand zu schreiben; das Werkzeug macht den Weg nicht fälschbarer, aber es macht die Zusage „in dieser Umgebung geschrieben“ falsch. Grund: Eine Herkunft ohne Signatur (Option C) ließe sich ebenso von Hand schreiben und kostete Format und Hook.
- Akzeptiertes Negativ: Nach einer Übernahme verlangt die nächste inerte Änderung im Hauptbaum den Stand der schnellen Kette, also Docker. Grund: Die Entscheidung des Nutzers, die zweite Sperre nicht zu ändern (Option D); die Probe misst, wie oft das eintritt.
- Folgepflicht: [`SPEC-062`](../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) und die Nachzüge in [`SPEC-052`](../../../harness/vertraege/SPEC-052-gate-kette.md) (Punkte 1, 8 und 21) liegen vor dem Code; der Slice liefert Werkzeug, Gegenprobe, Gate-Index und das Rezept für den Transport; ob und wann die Rollen das Werkzeug aufrufen, entscheidet der Nutzer nach der Lieferung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git, Temp-Repos | Bedingungen der Übernahme, Reihenfolge der Gründe, Ablehnung ohne Änderung, Aufrufarten, Kopplung an den unveränderten Hook | `make nachweis-uebernehmen-gegenprobe` |

Die Gleichheit von Quelle und Umgebung (Image, Werkzeugversionen) prüft kein Gate ([`SPEC-062`](../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 21).

## Re-Evaluierungs-Trigger

Wenn ein übernommener Nachweis einen Baum freigab, dessen `make gates` in dieser Umgebung rot gewesen wäre (dann Option C oder die Übernahme nur aus einer benannten Umgebung); wenn der Nutzer die zweite Sperre des Hooks ändern will (Option D, eigene ADR); wenn der Standardweg „Arbeit im Container“ den Pull in den Hauptbaum ersetzt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect nach dem Auftrag des Nutzers) | — |
| 2026-10-18 | Accepted: Entscheidung des Nutzers, Option B (zweiter Schreiber des Nachweises, Übernahme bei byte-gleicher Dateiliste) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
