# Review-Report: slice-v1-abschluss-sessions-doku — 2026-10-13

**Review-Art:** Code-Review (Dokumentations-Slice). Geprüft wird gegen Plan §1, §3, §6, §8, gegen [`LH-FA-12`](../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen) und [`LH-FA-13`](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) (über die Randformen der Spezifikation) und gegen `AGENTS.md` §3.3, §3.7, §3.9, §3.11, §3.13. Die DoD prüft der Verifier.

**Gegenstand:** `c7bce82..HEAD` (drei Commits: `1fe572a` Handbuch 0.6 und README, `a4f956f` Belege in §7, `1f34991` Stand der Gates). Geändert: `docs/user/benutzerhandbuch.md`, `README.md`, `docs/plan/planning/in-progress/slice-v1-abschluss-sessions-doku.md` (nur §7).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-sessions-doku` ganz gelesen, mit §7 *Belege des Implementers* (Probe-Tabelle Zeilen 1 bis 40, Sätze neben Bestand, Befehlsfolge); den Bericht des Implementers habe ich nicht gelesen.
- [`LH-FA-12`](../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-13`](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus); aus der Spezifikation die Abschnitte *Randformen des Recordings*, *Replay*, *Randformen der Zuordnung* und die Zeilen zu Verbindungsfehlern und Startnachrichten.
- Vorherige Reports am selben Modul: `slice-v1-abschluss-zeitangaben-doku` (Lehre: Kettenaussage ohne eigene Ausnahme, Zustandswort weggelassen).
- `AGENTS.md` (Hard Rules); Diff von Handbuch und README ganz gelesen, die Nachbarsätze (Exit-Codes, Fehlerbehandlung, Fehlercodes, Glossar, Treiber-Abschnitt, Fehlerbehebung, FAQ) gegen den neuen Abschnitt gehalten.

**Ausgeführte Läufe:**

- Befehlsfolge aus `done/slice-doku-ist-stand` §7 selbst gefahren, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`, gegen das vorhandene Image `pgwire-recorder:dev` (`deb48422487d`, wie in §7): Optionen ohne `--help`-Treffer nur `--h --help --name --rm --version`; Variablen ohne Option: leer; Codes außerhalb des Katalogs: leer. Entspricht §7.
- `grep` über die Plus-Zeilen von Handbuch und README auf `spec/`, Lastenheft, Slice, Welle, Review, `LH-`, `SPEC-`, ADR, „noch nicht“, „geplant“, „bisher“, „derzeit“: kein Treffer; ein Treffer auf „kommt“ ist Bestandswortlaut (`kommt` in der TLS-Zeile, nur umgebrochen). README-Plus-Zeilen: kein Verweis.
- Am Binary (nur Optionsebene, ohne Datenbank): `replay --session-assignment` mit `Connection`, `first_request` und leerem Wert: `PGR-E2001`, Meldung „erlaubt sind first-request und connection“; `record --session-assignment=connection`: `PGR-E2001` „flag provided but not defined“; `--help` nennt `--session-assignment first-request|connection` bei `replay` und `--record-empty-sessions[=true|false]` samt Variable bei `record`. Weitere Proben (Datei-Schlüssel, Variablen) scheiterten in meinem Aufbau an einer Rechtefrage der Arbeitsverzeichnis-Datei und zählen nicht.
- **Nicht am Binary gefahren:** alle Proben mit Datenbank oder Netz: `record` mit mehreren Verbindungen, Reihenfolge der `id`, `empty_sessions` in YAML und SQLite, Zwangsende, `COPY`, Passwort, die `printf`-Fälle (`SSLRequest`, `CancelRequest`, 1/4/7/8/20 Bytes), `replay` mit `connection` und `first-request`, `PGR-E5001`/`E5002`/`E5003`/`W2001`, `play`. Diese Aussagen habe ich gegen die Probe-Tabelle in §7 und gegen die Spezifikation gelesen, nicht selbst bestätigt.
- **Commit-Messages** `c7bce82..HEAD`: alle drei nennen `slice-v1-abschluss-sessions-doku` und `LH-FA-12`; kein Move im Diff (§3.3). `make gates`: nicht gefahren (Auftrag); `make docs-check` siehe Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-686 | LOW | Der Satz in *Verbindungsfehler* sagt, ein später folgender Fehler einer anderen Klasse ändere den Exit-Code nicht. Der Abschnitt *Exit-Codes* nennt zwei Ausnahmen: Scheitert beim Beenden das Schreiben der Aufzeichnung, ist der Exit-Code 3, auch wenn vorher ein anderer Fehler auftrat; ein Abbruch der Verbindung bei `--compare-responses` ist 4 trotz früherer Abweichung. Probe-Zeile 33 deckt nur das Paar `PGR-E6001` dann `PGR-E4006`. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · `ein später folgender Fehler einer anderen Klasse ändert ihn nicht.`; Gegenstück: `Aufzeichnung, ist der Exit-Code 3, auch wenn vorher ein anderer Fehler auftrat.` | ja — erst `COPY`, dann nicht schreibbares Ziel beim Beenden | Kettenaussage ohne die Ausnahmen des Nachbarsatzes |
| F-687 | LOW | Mit mehreren Verbindungen sagen „der Exit-Code ist der des ersten gemerkten Fehlers“ (neuer Abschnitt) und „den Exit-Code des ersten aufgetretenen Fehlers“ (Exit-Codes) zu, welcher Fehler zuerst zählt. Die Spezifikation schränkt für die Randform Verbindungsfehler von [`LH-FA-13`](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) ein: Entstehen Fehler verschiedener Verbindungen gleichzeitig, ist nicht zugesichert, welcher der erste ist. Kein Satz des Handbuchs trägt die Voraussetzung; Probe-Zeile 33 läuft mit zeitlich getrennten Fehlern. „gemerkt“ ist im neuen Abschnitt ohne Erklärung (Bestand verwendet es nur bei `PGR-E4006`). | `AGENTS.md` §3.12 (Aufruf-Voraussetzung), §3.11 | `docs/user/benutzerhandbuch.md` · `ersten gemerkten Fehlers, im ersten Fall 5, im zweiten 6`; `liefert den Exit-Code des ersten aufgetretenen Fehlers` | nein — Gleichzeitigkeit ist nicht deterministisch zu proben | Aufruf-Voraussetzung (gleichzeitig/nacheinander) fehlt im Nachbarsatz |
| F-688 | LOW | „Die `id` der Sitzungen zählt ab 1 in der Reihenfolge, in der die Sitzungen enden“ gilt ohne Bedingung. Die Spezifikation (*Gleichzeitiges Ende*) sagt: Enden Sessions gleichzeitig, bleiben die `id`s eindeutig und ohne Lücken, welche Session welche erhält, ist nicht zugesichert. Probe-Zeile 7 zeigt nur Enden im Abstand von 3 s. | `AGENTS.md` §3.11, §3.12 | `docs/user/benutzerhandbuch.md` · `Die `id` der Sitzungen zählt ab 1 in der Reihenfolge, in der die Sitzungen enden` | nein | Aufruf-Voraussetzung nicht an die Zusage gebunden |
| F-689 | LOW | „Das gilt auch für eine Verbindung, die beim zwangsweisen Beenden nach Ablauf von `--shutdown-timeout` … keine gestellt hatte“ sagt, eine ruhende Verbindung ohne Anfrage werde beim Zwangsende mit der Option zur Sitzung. Probe-Zeile 11 beendet diese Verbindung mit `docker stop` vor der Frist (Exit 0), nicht durch das Zwangsende. Die Spezifikation nimmt den Fall aus: Trifft das Zwangsende eine solche Session, bevor ihr Ende bemerkt ist (zweites Signal), wird sie auch mit der Option nicht geschrieben. Ein Satz neben dem Bestandssatz „Eine Verbindung ohne begonnene Anfrage … endet ohne Meldung“. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · `` `--shutdown-timeout` ihre erste Anfrage noch ausführte (`PGR-E4006`) oder keine gestellt `` | ja — Verbindung ohne Anfrage offen lassen, zweites Signal senden, mit Option | Satz weiter als die Probe (Zwangsende vs. geordnetes Ende) |
| F-690 | LOW | §7 *Ungenannt geblieben* führt „der Satz ‚nur zugesichert, wenn nacheinander‘ … keine Probe“ als nicht im Handbuch stehend. Der neue Abschnitt *Voraussetzung* enthält ihn (`ist nur zugesichert, wenn Ihre Anwendung ihre Verbindungen nacheinander aufbaut und nutzt`), ohne Zeile in der Probe-Tabelle; der Bestandssatz bei `first-request` sagt dasselbe enger. Plan und Handbuch widersprechen sich. | `AGENTS.md` §3.9, §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-sessions-doku.md` · `die Zuordnung bei gleichzeitigem Aufbau (der Satz „nur zugesichert, wenn nacheinander“ sagt zu`; `docs/user/benutzerhandbuch.md` · `zugesichert, wenn Ihre Anwendung ihre Verbindungen nacheinander aufbaut und nutzt.` | nein | Plan §7 widerspricht dem Diff |
| F-691 | LOW | §1 sagt, ob das Handbuch die Verbindungen nennt, die keine Session sind, entscheide der Architect vor dem ersten Commit. §7 *Rückgaben* hält fest, das Handbuch nenne sie, soweit Zeilen 17, 19 und 20 sie zeigen, und „keine Randform entschieden“; eine Entscheidung des Architect ist in Plan oder Diff nicht festgehalten. Ob sie mündlich fiel, liegt mir nicht vor. | `AGENTS.md` §3.12, §3.13 | `docs/plan/planning/in-progress/slice-v1-abschluss-sessions-doku.md` · `Ob das Handbuch diese Fälle nennt, entscheidet der Architect vor dem ersten Commit`; `Zur Frage aus §1 (nennt das Handbuch die Verbindungen, die keine Session sind): Es nennt sie` | nein | Entscheidung des Architect ohne Artefakt |
| F-692 | INFO | Probe-Tabelle (Aufbau, Wiederholbarkeit): Keine Zeile benennt eine von Hand editierte Aufzeichnung; alle Aufzeichnungen (`a` bis `g`) entstehen durch `record`, die vom Server ausgelösten Fälle (Zeilen 11, 15 bis 18) laufen gegen echte Aufbauten. Wiederholbar ist der Aufbau über Image-Stand, gepinntes PostgreSQL und Beschreibung der Clients; Skripte und Dateien liegen bewusst außerhalb des Repos, einige Zeilen nennen den Aufruf nur verkürzt (Zeile 29: `PGSSLMODE=prefer`, Zeile 20: Byte-Längen). Hinweis, keine Aktion erwartet. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-sessions-doku.md` · `**Probe-Aufbau.** `make build` am Stand `cd57943`` | nein | Probe-Aufbau nur beschrieben |
| F-693 | INFO | Die Zeile `PGR-E4006` des Fehlercode-Katalogs nennt die Sitzung, nicht die `id`, unter der sie mit der Option geschrieben wird; die Spezifikation nennt die `id`. Das Handbuch sagt dazu nichts Falsches. Hinweis, keine Aktion erwartet. | Maintainability | `docs/user/benutzerhandbuch.md` · `Die Meldung nennt die Sitzung und die Anfrage; beim Aufzeichnen fehlt diese Anfrage in der Aufzeichnung.` | nein | Lücke ohne Falschaussage |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md`, neuer Abschnitt *Mehrere Verbindungen aufzeichnen und wiedergeben*: Wertemenge `--session-assignment` (genau `first-request`/`connection`, klein, leerer Wert `PGR-E2001`, Datei `PGR-E2004`, Variable nennt den Namen), `connection` verlangt `empty_sessions` (`PGR-E2001` beim Start), `PGR-E5001`/`PGR-E5003`/`PGR-E3004`, Lebendprüfung beantwortet, Sitzung ohne Interaktionen nicht verbraucht, Zählung der Verbindungen | geprüft, ohne Befund gegen die Spezifikation; jeder Satz hat eine Zeile in der Probe-Tabelle (Optionsebene teils selbst bestätigt) |
| Fälle „keine Session“ (`SSLRequest`/`GSSENCRequest`, `CancelRequest`, Ende vor der Startnachricht, 7 Bytes still, ab 8 Bytes `PGR-E6001` Exit 6, Zwangsende im Aufbau, Aufbau scheitert `PGR-E4002`, `COPY`) | geprüft, ohne Befund: jeder genannte Fall steht in der Spezifikation und in den Zeilen 15 bis 20 der Probe-Tabelle; das Handbuch nennt nicht mehr als dort gedeckt (`GSSENCRequest` nur unter „Anfrage nach Verschlüsselung“). Die Wortwahl „Zwangsende im Aufbau“ ist gedeckt, die Ausweitung auf ruhende Verbindungen nicht (F-689) |
| Aufgezeichnetes Verbindungsende des Servers (`FATAL`, `PANIC`) | geprüft, ohne Befund: Handbuch und README sagen dazu nichts (akzeptierte Lücke) |
| Nachbarsätze: *Hinweise* der Abschnitte *Aufzeichnen*, *Wiedergeben*, *Einspielen*, SQLite-Abschnitt und Tabelle `meta`/`session`, Satz zu `--fail-on-unconsumed`, Tabelle §5 und Absatz zu Variablen beim anderen Kommando, Exit-Code 5, Zeilen `PGR-E2001`, `E2004`, `E3004`, `E5001`, `E5002`, `E5003`, `W2001`, Glossar, Einleitung, Treiber-Abschnitt, Fehlerbehebung, FAQ | geprüft, ohne Befund bis auf F-686 bis F-689; keine Bestandsaussage „jede Verbindung ohne Anfrage ist keine Sitzung“ blieb stehen; die „Verbindung“-/„Sitzung“-Wörter im Glossar sind mit der Kennzeichnung der Sitzung ohne Interaktionen konsistent |
| `PGR-E2001` für Variablen beim anderen Kommando (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`) und Schlüssel im anderen Abschnitt `PGR-E2004`; Satz zu Variablen mit Präfix | geprüft, ohne Befund: Zeile 34 deckt Option, Variable (auch `false`) und Schlüssel; der Bestandssatz zu Variablen ohne Option bleibt gedeckt; die Befehlsfolge zeigt beide Variablen mit Option |
| `README.md`: ein Absatz im Ist-Zustand, keine Chronik, kein Zielstand, kein neuer Verweis; Sätze zu Standard, `connection`, `PGR-E2001`, `PGR-E5003`, `PGR-E5001` | geprüft, ohne Befund (Zeilen 1, 7, 9, 10, 22 bis 27 der Probe-Tabelle) |
| Plan §1, §3, §6, Kopf, §8 gegen den Diff (§3.9): der Diff ändert nur §7; §3 nennt die zwei Dateien; Liefer-Punkte 2 und Schicht 1 in §8 nachgezählt; die zwei Risiken in §6 stehen mit Ausgang „offen bis Closure“ | geprüft, ohne Befund (F-690 und F-691 betreffen §7 und §1) |
| Befehlsfolge in §7: Variablen-Zeile mit der benannten Ausnahme `PGWIRE_RECORDER_PASSWORD`, keine zweite Ausnahme, Optionen-Zeile nennt nur `--h`, `--help`, `--version`, `--name`, `--rm` | geprüft, ohne Befund; selbst gefahren, Ausgabe wie in §7 |
| Adressen und Übernahmen aus `slice-v1-abschluss-umgebung-anderes-kommando` und `slice-v1-abschluss-sessions` im Handbuch umgesetzt (Variable beim anderen Kommando, Umgebungsvariablen, Schlüssel, Kennzeichnung in YAML und SQLite, Grenzen) | geprüft, ohne Befund (§3.13 Nachzählen in §8) |
| Commits: reiner Inhalts-Commit je Schritt, kein Move (§3.3), jede Message mit Slice-Kennung und `LH-FA-12`, keine `SPEC-`/`ARC-`-Kennung | geprüft, ohne Befund |
| Produkt-Code, Tests, Gates, `docs/user/abdeckung-*.md`, `spec/`, ADRs | im Diff unberührt, geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kettenaussage ohne die Ausnahmen des Nachbarsatzes · Aufruf-Voraussetzung (gleichzeitig/nacheinander) fehlt im Nachbarsatz · Aufruf-Voraussetzung nicht an die Zusage gebunden · Satz weiter als die Probe (Zwangsende vs. geordnetes Ende) · Plan §7 widerspricht dem Diff · Entscheidung des Architect ohne Artefakt · Probe-Aufbau nur beschrieben · Lücke ohne Falschaussage

## Verdikt

**Merge-blockierend:** nein — keine HIGH- oder MEDIUM-Findings; die sechs LOW sind Formulierungen und Plan-Nachzug, die der Implementer vor der Closure schärfen kann. F-687 und F-688 sind die Fortsetzung der Klasse „aus einem Nachbarsatz übernommene Zusage ohne seine Voraussetzung“ und gehören in die Closure §7.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Gate-Ergebnis: `make gates` nicht gefahren (Auftrag des Aufrufers).
