# Review-Report: slice-v1-abschluss-zeitangaben-doku — 2026-10-12

**Review-Art:** Code-Review (Dokumentations-Slice). Geprüft wird gegen Plan §1, §3, §6 und §8, gegen [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen), gegen [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) und gegen `AGENTS.md` §3.3, §3.7, §3.9, §3.11, §3.13. Die DoD prüft der Verifier.

**Gegenstand:** `56d7797..c2c7272` (zwei Commits: `15583ea` Roadmap-Zeile, `c2c7272` Handbuch 0.4, README, Plan §3/§7/§8). Geändert: `docs/user/benutzerhandbuch.md`, `README.md`, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-doku.md`.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-zeitangaben-doku` ganz gelesen, mit §7 *Belege des Implementers* (Probe-Matrix, Befehlsfolge, `grep`); den Bericht des Implementers habe ich nicht gelesen.
- [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen); die Spezifikation nur für den *Messpunkt* (Abschnitt in `spec/spezifikation.md`).
- Vorherige Reports am selben Modul: `slice-v1-abschluss-zeitangaben` und seine Teile (Nummern dieses Laufs ab F-660).
- `AGENTS.md` (Hard Rules), Handbuch und README am Stand `c2c7272` ganz im Diff gelesen, die Umgebung (Zeilen 520 bis 545 des Handbuchs) gegen die neuen Sätze gehalten.

**Ausgeführte Läufe** (Docker-only; `make build` am Stand `c2c7272`, Image `pgwire-recorder:dev`; eigenes Netz, ein PostgreSQL aus dem gepinnten Image von `harness/mk/integration.mk` mit `log_statement=all`; Aufzeichnungen und Skripte im Sitzungsverzeichnis außerhalb des Repos; Container mit `docker rm -fv`, Netz per Name entfernt, keine hostweite Bereinigung):

- `record --record-timing` mit `psql` in einer Sitzung (Pausen 2 s und 1 s): `offset_ms` 4, 2006, 3011; erste Anfrage nicht 0. Die zweite Aufzeichnung ohne Pause: 4.
- `play --keep-timing` (`relative`) und `--timing-mode absolute`, Zeitstempel am Server: erste Anfrage 2 ms nach dem Verbindungsaufbau (ohne Wartezeit, aufgezeichnet 4 ms), danach +2,0 s und +1,0 s in beiden Modi.
- Aufzeichnung mit Werten 3000, (fehlt), 500: `play --keep-timing` endet mit `PGR-E3003`, Exit 3, `Session 1, Interaktion 3`; der Vergleich läuft also gegen den vorigen vorhandenen Wert, wie das Handbuch sagt.
- `PGWIRE_RECORDER_RECORD_TIMING=true` bei `play`: läuft ohne Meldung und ohne Wirkung; `replay --keep-timing`: `PGR-E2001`; `PGWIRE_RECORDER_KEEP_TIMING=TRUE`: `PGR-E2001` ohne Wert in der Meldung.
- `SIGTERM` bei der ersten Wartezeit, mit `--finish-session-on-interrupt`: übrige Anfragen sofort (Abstand 3 ms), Exit 0; `SIGINT` ohne die Option: Ende nach der ersten Anfrage, Exit 0.
- `grep` über das Handbuch auf `ADR-`, `SPEC-`, `LH-`, `[Ss]lice`, `Welle`, `Spezifikation`, `Lastenheft`, `[Rr]eview`: kein Treffer; Treffer auf „kommt“ und „später“ liegen im Bestand (Zeilen 14, 508, 783, 894) und betreffen nicht die Zeitangaben; in den neuen Zeilen kein „noch nicht“, „bisher“, „geplant“, „compare-responses“, Windows, macOS.
- **Commit-Messages** `56d7797..HEAD`: keine `SPEC-`/`ARC-`-Kennung, beide nennen `slice-v1-abschluss-zeitangaben-doku` und `LH-FA-21`. Kein Move im Diff (§3.3).
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-660 | LOW | Handbuch und README sagen, `play --keep-timing` warte „vor jeder Anfrage“ auf die aufgezeichnete Zeit. Beim Standard `relative` sendet `play` die erste Anfrage jeder Verbindung ohne Wartezeit (am Server 2 ms nach dem Verbindungsaufbau bei aufgezeichneten 4 ms); die Ausnahme steht erst im Abschnitt *Hinweise zum Einspielen*. Die Probe-Matrix in §7 hat zu dem Satz „vor jeder Anfrage“ keine Gegenprobe in Richtung der ersten Anfrage. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · `Beim Einspielen mit `--keep-timing` wartet `play` vor jeder Anfrage auf ihre aufgezeichnete`; ebenso `Mit `--keep-timing` wartet `play` vor jeder`, `es vor jeder Anfrage auf ihre aufgezeichnete Zeit` und `README.md` · `wartet vor jeder Anfrage auf diese` | ja — `play --keep-timing` von einer Aufzeichnung mit großem ersten `offset_ms`, Zeitstempel am Server | Kettenaussage „jede“ ohne die eigene Ausnahme |
| F-661 | LOW | Der neue Satz zum Abbruchsignal endet mit „Der Exit-Code ist 0.“; der Bestand zwei Abschnitte darüber sagt für dasselbe Signal „Der Exit-Code ist 0, nach einem früheren Fehler 4.“. Der neue Satz gilt ohne den Nachsatz, und keine Zeile der Probe-Matrix läuft mit einem früheren Fehler vor dem Signal. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · `der Option spielt es die übrigen Anfragen der Sitzung ohne Wartezeit ein und endet dann.`; Bestand: `Der Exit-Code ist 0, nach einem früheren Fehler 4.` | ja — Aufzeichnung mit einer fehlschlagenden Anfrage vor der Wartezeit, Signal in der Wartezeit | Zustandswort im Satz weggelassen, Gegenprobe fehlt |
| F-662 | LOW | Einleitung, README, Abschnittsanfang, Hinweis und Glossar nennen als Bezug von `offset_ms` den „Beginn der Verbindung“; erst ein späterer Punkt sagt, dass die Zeit bis zum Senden der Startup-Nachricht nicht zählt (Probe in §7: 2 s vor dem Startup ergeben 4 ms). Der Glossar-Eintrag und die Sätze der Einleitung tragen die Einschränkung nicht, obwohl §1 des Plans den Messpunkt (Beginn der Session nach TLS und Startup des Clients) als Lesart verlangt. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · `die Millisekunden vom Beginn der Verbindung bis zum Eintreffen der Anfrage, aufgezeichnet mit`; `README.md` · `die Millisekunden seit dem Beginn` | ja — Roh-Client mit Pause vor dem Startup, `offset_ms` bleibt klein | Bezugspunkt ungenau, Einschränkung nur an einer Stelle |
| F-663 | INFO | Länge und Gliederung: Der neue Abschnitt umfasst rund 120 Zeilen. *Hinweise zum Einspielen* trägt zehn Punkte, darunter den Punkt der Kombinationen mit acht Sätzen und Verschachtelung („Als gesetzt gilt …“, „ebenso …“, „auch mit dem Standard“); *Ergebnis* für das Einspielen ist ein Satz ohne beobachtbare Ausgabe. Die Sätze sind nach meinen Proben richtig; sie sind schwer zu lesen. Hinweis, keine Aktion erwartet. | Maintainability | `docs/user/benutzerhandbuch.md` · `Ungültig sind außerdem (`PGR-E2001`, Exit-Code 2)` | nein | Dichte Regelaufzählung in einem Punkt |
| F-664 | INFO | Zu den zwei Hinweisen des Implementers: Das Handbuch sagt zu beiden nichts, und nichts davon ist falsch oder unbelegt. Die Tabelle in §5 führt `PGWIRE_RECORDER_RECORD_TIMING` nur bei `record` (bei `play` ohne Meldung und ohne Wirkung, von mir bestätigt); die Sitzung ohne Interaktion berührt kein Satz der Zeitangaben. Beide gehören an den Planner, wie §7 sie führt. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-doku.md` · `Hinweise an Architect und Planner` | nein | Verhalten neben der Spezifikation, dem Planner gemeldet |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md`, neuer Abschnitt: Ordnung der Startfehler (`PGR-E3003` vor `PGR-E2003`, Exit 3 gegen 2), Kombinationen („gesetzt“ = irgendeine Quelle, Standardwert zählt, Reihenfolge ohne Wirkung), Wertemenge (Kleinbuchstaben; Option und Variable `PGR-E2001`, Datei `PGR-E2004`; Meldung ohne Wert), Messpunkt (nach Startup), Priorität Option vor Umgebung vor Datei, `replay` wartet nie, Vergleich gegen den vorigen vorhandenen Wert | geprüft, ohne Befund (F-660 bis F-662 betreffen nur Formulierungen); Stichproben am Binary bestätigt |
| `docs/user/benutzerhandbuch.md`, Tabelle §5, Fehlercodes `PGR-E2001`/`E2003`/`E2004`/`E3003` (Reihenfolge der Codes, Ursache und Lösung), SQLite-Tabelle, FAQ, Glossar, Einleitung | geprüft, ohne Befund (F-662 für den Glossar-Eintrag) |
| Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews; ohne Chronik oder Zielstand; keine Zusage zu Windows, macOS, Netzwerk-Dateisystemen, großen Dateien, `--compare-responses`; Grenze „nur Linux im Container“ unverändert | geprüft, ohne Befund (`grep`, siehe oben) |
| `README.md`: ein Absatz, Ist-Zustand, Verweise auf `spec/` unverändert | geprüft, ohne Befund (F-660, F-662 betreffen zwei Formulierungen) |
| Meldungen nennen nie einen Wert (Option darf): `PGR-E2003`, `PGR-E3003`, Umgebung und Datei | geprüft, ohne Befund |
| Plan §3, §7, §8 gegen den Diff (§3.9): §3 nennt die Roadmap-Zeile und die zusätzlichen Stellen im Handbuch; §8 zählt zwei Liefer-Punkte und eine Schicht, der Diff ändert kein Produkt, keinen Test, kein Gate; §7 führt die Befehlsfolge mit Mutanten und die Namen der genannten Dateien (Namensabgleich) | geprüft, ohne Befund |
| Probe-Matrix in §7 (Lehre V-180/F-636): jeder Verhaltenssatz des neuen Abschnitts hat eine Zeile, bis auf die in F-660 und F-661 genannten Richtungen | geprüft, mit Befund (F-660, F-661) |
| Nehmer-Sendungen aus `slice-v1-abschluss-zeitangaben` und `-aufzeichnen` (Ordnung der Fehlercodes, Kombinationen, Wertemenge, Messpunkt-Lesart) im Handbuch umgesetzt; nachgezählt in §8 | geprüft, ohne Befund (Messpunkt: F-662) |
| Commits: reiner Inhalts-Commit je Schritt, kein Move (§3.3), keine `SPEC-`/`ARC-`-Kennung ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)); Roadmap-Zeile *In Arbeit* nennt Zustand ohne Chronik (§3.7) | geprüft, ohne Befund |
| Produkt-Code, Tests, Gates, `docs/user/abdeckung-*.md`, `spec/`, ADRs | im Diff unberührt, geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kettenaussage „jede“ ohne die eigene Ausnahme · Zustandswort im Satz weggelassen, Gegenprobe fehlt · Bezugspunkt ungenau, Einschränkung nur an einer Stelle · Dichte Regelaufzählung in einem Punkt · Verhalten neben der Spezifikation, dem Planner gemeldet

## Verdikt

**Merge-blockierend:** nein — keine HIGH- oder MEDIUM-Findings; die drei LOW sind Formulierungen, die der Implementer vor der Closure schärfen kann.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Gate-Ergebnis siehe Bericht des Reviewers.
