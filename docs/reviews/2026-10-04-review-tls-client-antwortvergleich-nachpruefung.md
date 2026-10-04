# Review-Report: tls-client-antwortvergleich, Nachprüfung — 2026-10-04

**Review-Art:** Nachprüfung der Korrekturen zu `docs/reviews/2026-10-04-review-tls-client-antwortvergleich.md` (F-141 bis F-163), Konsistenzprüfung Lastenheft, Spezifikation, Architektur, Handbuch und Slices, Prüfung der Hard Rules.

**Gegenstand:** `git diff 6efb99b HEAD` (HEAD ca35c70; Commits 4076bc4, 14ee691, ca35c70): `spec/lastenheft.md` (LH-FA-24, Abnahmeszenario 16), `spec/spezifikation.md` (LH-FA-01.a, LH-FA-03.a, LH-FA-20.a Schritte 2, 6, 7, LH-FA-23.a, LH-FA-24.a, SPEC-018, Codetabelle SPEC-034), `docs/user/benutzerhandbuch.md`, `slice-v1-abschluss-tls-client`, `slice-v1-abschluss-antwortvergleich`, `slice-v1-abschluss-einspielen`. `spec/architecture.md` und alle ADRs sind im Diff unverändert.

**Skill:** `.harness/skills/reviewer.md` @ ca35c70
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:** `AGENTS.md` (Hard Rules), `harness/README.md`, `harness/conventions.md`; der Vorgänger-Report; [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0018](../plan/adr/0018-tls-zum-client.md), [ADR-0019](../plan/adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md), [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md); berührte `LH-*`: LH-FA-05, -13, -17, -18, -20, -23, -24, LH-QA-05. Dieser Lauf setzt mit F-164 und S-35 fort. Nicht Gegenstand: DoD-Konformität (Verifier), fachliche Eignung (Validator). `make gates` wurde auftragsgemäß nicht ausgeführt. Lastenheft-Status `Draft`.

---

## Stand der Findings des Vorgänger-Reports

| ID | Stand | Anmerkung |
|---|---|---|
| F-141 | behoben | Schritt 7 nennt „der zuerst aufgetretenen Ursache (4, mit `--compare-responses` auch 5)"; deckt sich mit Schritt 6. Rest siehe F-165. |
| F-142 | teilweise | Vorrang `PGR-E5004` vor `PGR-E4004` ist festgelegt, Lastenheft und Handbuch tragen ihn. Neue Widersprüche: F-164, F-165, F-167, F-170. |
| F-143 | teilweise | `Flush`-Gruppen sind durch den Vergleich der Gesamtfolge erledigt (passt zur Architektur-Sequenz „bis ReadyForQuery"); der Fehler in Extended ist geregelt. Neuer Widerspruch: F-166. |
| F-144 | teilweise | Regel vorhanden (nicht verglichen, Log `info`), aber mit Widersprüchen: F-167. |
| F-145 | behoben, mit Rest | Klartext-Startup nach `GSSENCRequest` ist `PGR-E6003`, Klartext-`CancelRequest` bleibt `PGR-W3001`. Rest: F-171. Exit-Code-Wirkung bleibt Setzung S-30. |
| F-146 | behoben | `--upstream-ca` wird beim Start gelesen (`PGR-E2007`, Exit-Code 2), Pfade relativ zum aktuellen Verzeichnis, verschlüsselter Schlüssel ist `PGR-E2007`; Handbuch gleichlautend. Rest: F-174 (INFO). |
| F-147 | behoben | Abgelaufen oder noch nicht gültig ist ungültig (`PGR-E2007`); Spezifikation, Codetabelle, Handbuch und Slice stimmen überein. Namen prüft der Recorder nicht, das ist im Lastenheft nicht verlangt. Die frühere Setzung S-29 ist in ihrem Kern erledigt. |
| F-148 | teilweise | Lastenheft trägt jetzt Datentypen, Parameter, Transaktionsstatus, Zeilenanzahl, Hinweise und Einstellungen. Rest: F-168. |
| F-149 | behoben, mit Rest | Negativfälle sind in den DoD der drei Slices ergänzt (`--tls-key` allein, `--allow-plaintext` ohne Zertifikat, `GSSENCRequest`, Klartext-`CancelRequest`, mehrere Verbindungen, `--allow-recorded-errors`, `Flush`, Aufzeichnung ohne `ReadyForQuery`, `--upstream-ca` ohne TLS). Rest: F-169. |
| F-150 | behoben | Start-Trigger nennt `slice-v1-abschluss-sessions`. |
| F-151 | teilweise | SPEC-018 trägt jetzt „Einspielen: Abbruch beziehungsweise Exit-Code 5 am Ende". Die Klasse „Replay" für `PGR-E5004` in SPEC-034 und im Handbuch bleibt; nicht neu bewertet. |
| F-152 | behoben, mit Rest | Lastenheft verlangt jetzt einen Fehlerstatus, der sich von dem eines Serverfehlers ohne Vergleich unterscheidet (Exit-Code 5 gegen 4). Rest: F-165. |
| F-153 | behoben | Handbuch zeigt die Erzeugung (OpenSSL), Passwort-Hinweis und die Aussage, dass Namen und Vertrauen der Client prüft. |
| F-154 | behoben | Handbuch nennt Vorrang und Exit-Code der zuerst aufgetretenen Ursache; die Exit-Code-Tabelle führt Zeile 5 mit `--compare-responses`. |
| F-155 | behoben | Warnungstabelle ist aufsteigend. |
| F-159 | behoben | Optionslisten LH-FA-01.a und LH-FA-03.a nennen `--tls-cert`, `--tls-key`, `--allow-plaintext`; Tabelle LH-FA-17.a stimmt. |
| F-160 | behoben | Kopf von `slice-v1-abschluss-tls-client` führt SPEC-018, SPEC-019, SPEC-020; `slice-v1-abschluss-antwortvergleich` ergänzt SPEC-041. |
| F-156, F-157, F-158, F-161 | offen, unverändert | bleiben bekannte Punkte; die Korrekturen ändern ihre Lage nicht. Zu F-156: der Wortlaut in LH-FA-24 enthält weiter Protokollnähe („Spaltenbeschreibung", „Fehlercode"), nicht neu bewertet. |
| F-162, F-163 | INFO, unverändert | F-163: Stand-Zeile des Handbuchs ist auf 04.10.2026 nachgezogen. |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-164 | MEDIUM | LH-FA-24 (Happy Path) verlangt Erfolg, wenn Antworten „in Struktur und Fehlern der Aufzeichnung entsprechen"; LH-FA-20.a Schritt 6 lässt bei gleichem SQLSTATE die Regel von `PGR-E4004` (Abbruch, Exit-Code 4) unverändert gelten, außer mit `--allow-recorded-errors`. Ein Lauf mit Vergleich gegen eine Aufzeichnung mit Fehler, die der Server genau so wiedergibt, endet damit nicht mit Erfolg; LH-FA-20 (Boundary) verlangt für „erwartet" die ausdrückliche Wahl, LH-FA-24 nicht. Die Spezifikation entscheidet zwischen zwei Lastenheft-Aussagen, ohne es als Setzung zu kennzeichnen; das Handbuch nennt die Folge nicht. | LH-FA-24; LH-FA-20 | spec/spezifikation.md · LH-FA-20.a Schritt 6; spec/lastenheft.md · LH-FA-24 Happy Path, LH-FA-20 Boundary | nein | Spec entscheidet zwischen zwei Lastenheft-Aussagen ohne Setzung |
| F-165 | MEDIUM | Der Exit-Code am Ende ist der der zuerst aufgetretenen Ursache (Schritt 6 und 7). Tritt mit `--continue-on-error` zuerst ein Serverfehler mit gleichem SQLSTATE (`PGR-E4004`) und danach eine Abweichung (`PGR-E5004`) auf, endet der Lauf mit Exit-Code 4, obwohl eine Abweichung gemeldet wurde. LH-FA-24 (Negative) verlangt für eine Abweichung einen Fehlerstatus, der sich von dem eines Serverfehlers ohne Vergleich unterscheidet; Handbuch („Bei einer Abweichung endet das Einspielen mit … Exit-Code 5") sagt dazu nichts Einschränkendes. | LH-FA-24; LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Schritte 6 und 7; spec/lastenheft.md · LH-FA-24 Negative; docs/user/benutzerhandbuch.md · Einspielen, Exit-Codes | nein | Rangfolgeregel hebt zugesagte Unterscheidbarkeit auf |
| F-166 | MEDIUM | LH-FA-24.a vergleicht bei Extended „die Folge aller Server-Nachrichten ihrer Gruppen, ohne die Gruppengrenzen"; derselbe Abschnitt nennt in der Meldung von `PGR-E5004` „bei Extended den Gruppenindex". Ohne Gruppengrenzen im Vergleich ist kein Gruppenindex bestimmbar (die Zuordnung ist laut Text zeitabhängig). Auch die Art der Abweichung („Nachrichtenart, …") ist ohne Bezug zur Gruppe nicht eindeutig. | LH-FA-24; LH-QA-05 | spec/spezifikation.md · LH-FA-24.a | nein | Meldungsinhalt setzt entfallene Struktur voraus |
| F-167 | MEDIUM | Eine aufgezeichnete Interaktion ohne `ready_for_query` am Ende oder ohne Server-Nachrichten wird „nicht verglichen" (Log `info`). (a) LH-FA-20.a Schritt 6 sagt für eine aufgezeichnete Interaktion ohne `ErrorResponse`, eine Serverfehlerantwort sei `PGR-E5004`; für eine nicht verglichene Interaktion bleibt offen, ob `PGR-E5004` oder `PGR-E4004` gilt. (b) LH-FA-24.a zählt „eine zusätzliche oder fehlende Nachricht" zur Abweichung; ein Lauf mit Exit-Code 0 kann Interaktionen enthalten, die nicht verglichen wurden, und das Lastenheft kennt diese Ausnahme nicht. (c) Das Handbuch nennt die Ausnahme, aber nicht, dass sie nur im Log erscheint. | LH-FA-24; LH-QA-05 | spec/spezifikation.md · LH-FA-24.a, LH-FA-20.a Schritt 6; docs/user/benutzerhandbuch.md · Einspielen | nein | Ausnahmeregel kollidiert mit allgemeiner Fehlerregel |
| F-168 | LOW | LH-FA-24 nennt jetzt „der ausgeführte Befehl" als Vergleichsgegenstand; LH-FA-24.a vergleicht nur das erste Wort des Tags. `CREATE TABLE` gegen `CREATE INDEX` (oder `ALTER TABLE`/`ALTER INDEX`) bleibt unerkannt, der Befehl ist laut Lastenheft verglichen. Rest von F-148, Setzung S-31. | LH-FA-24 | spec/spezifikation.md · LH-FA-24.a Tabelle `command_complete`; spec/lastenheft.md · LH-FA-24 | nein | Spec vergleicht weniger als das Lastenheft nennt |
| F-169 | LOW | `slice-v1-abschluss-einspielen` DoD (Abbruchsignal) nennt „Exit-Code 0 ohne vorherigen Fehler, sonst 4"; LH-FA-20.a Schritt 7 sagt jetzt „der zuerst aufgetretenen Ursache (4, mit `--compare-responses` auch 5)". Der Slice wurde an dieser Stelle nicht nachgezogen; der Slice `antwortvergleich` testet die Rangfolge nur allgemein. | Maintainability | docs/plan/planning/open/slice-v1-abschluss-einspielen.md · §2 dritter Liefer-Punkt | nein | DoD folgt Spec-Änderung nicht |
| F-170 | LOW | „Je Interaktion entsteht höchstens eine Meldung" lässt offen, welche Meldung gilt, wenn eine Interaktion gleichzeitig strukturell abweicht und eine Serverfehlerantwort mit gleichem SQLSTATE trägt (`PGR-E5004` oder `PGR-E4004`); davon hängt nach F-165 der Exit-Code ab. | LH-FA-24 | spec/spezifikation.md · LH-FA-20.a Schritt 6 | nein | Tie-Break innerhalb einer Interaktion nicht festgelegt |
| F-171 | LOW | LH-FA-23.a sagt, ein `CancelRequest` „trägt keine Daten"; er trägt Prozess-ID und Schlüssel des Abbruchs. Gemeint ist offenbar „keine Nutzdaten und kein Vertraulichkeitsbedarf"; der Satz begründet die Ausnahme von der Ablehnungsregel ungenau. | Maintainability | spec/spezifikation.md · LH-FA-23.a „Unverschlüsselte Clients" | nein | Begründung der Ausnahme sachlich ungenau |
| F-172 | LOW | Redaktion: In LH-FA-20.a Schritt 2 steht „Er authentifiziert sich" jetzt hinter dem Satz zum relativen Pfad, das Subjekt ist nicht mehr der Recorder; mehrere eingefügte Sätze (LH-FA-23.a „Unverschlüsselte Clients", LH-FA-24.a, Abnahmeszenario 16 im Lastenheft) brechen den Zeilenumbruch der Umgebung. | Maintainability | spec/spezifikation.md · LH-FA-20.a Schritt 2, LH-FA-23.a, LH-FA-24.a; spec/lastenheft.md · §7 Abnahmeszenario 16 | nein | Bezug nach Einfügung unklar |
| F-173 | INFO | Der Vorrang von `PGR-E5004` vor `PGR-E4004`, der Vergleich ohne Gruppengrenzen und das Nichtvergleichen unvollständiger Aufzeichnungen stehen nur in der Spezifikation; [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) (Accepted, unverändert, Hard Rule 3.5 gewahrt) sagt dazu nichts und nennt als Vergleichsumfang nicht die `parameter_description`. Wer diese Entscheidungen begründet sucht, findet sie nirgends; ein Folge-ADR mit `Supersedes` oder eine ergänzende ADR wäre der Ort. Verwandt mit F-161 und F-157, nicht neu bewertet. | AGENTS.md §3.5 | spec/spezifikation.md · LH-FA-20.a, LH-FA-24.a; docs/plan/adr/0020 | nein | Entscheidung ohne ADR-Träger |
| F-174 | INFO | `PGR-E2007` gilt für `--upstream-ca` bei einer Datei, die „kein gültiges Zertifikat enthält"; ob ein abgelaufenes CA-Zertifikat darunter fällt (für das eigene Zertifikat ist Ablauf ausdrücklich ungültig), und was geschieht, wenn ein Zertifikat während eines langen Laufs abläuft, ist nicht gesagt. Hinweis, keine Aktion erwartet. | LH-FA-23; LH-FA-20 | spec/spezifikation.md · LH-FA-20.a Schritt 2, LH-FA-23.a | nein | Gültigkeitsbegriff nicht einheitlich geführt |

## Setzungen

Stand der früheren Setzungen und neue Setzungen. Zählung der unbestätigten Setzungen: S-1 bis S-27 der früheren Reports und S-28 bis S-34 des Vorgänger-Reports bleiben unbestätigt, soweit nicht unten anders vermerkt.

| Nr. | Stand / Setzung | Fundort | Überschreitet das Lastenheft? |
|---|---|---|---|
| S-29 | im Kern erledigt: Ablauf und „noch nicht gültig" gelten jetzt als ungültig (Lastenheft „ungültiges Zertifikat"); die Namen prüft der Recorder nicht | LH-FA-23.a | nein mehr; Restsetzung siehe S-38 |
| S-30 | weiter teilweise: Exit-Code 6 für eine abgewiesene Klartext-Verbindung; `GSSENCRequest` mit `N` und Klartext-Startup danach: `PGR-E6003` | LH-FA-23.a, SPEC-019 | teilweise |
| S-31 | weiter: nur das erste Wort des Befehls wird verglichen; alle übrigen Grenzen trägt jetzt das Lastenheft | LH-FA-24.a | ja, nach unten (Lastenheft: „der ausgeführte Befehl"), siehe F-168 |
| S-32 | im Kern erledigt: Exit-Code 5 und „Abweichung" für Serverfehler ohne Vorbild trägt das Lastenheft; Restsetzung S-36 | LH-FA-20.a Schritt 6 | nein mehr; Restsetzung siehe S-36 und F-165 |
| S-35 | Eine aufgezeichnete Interaktion ohne `ready_for_query` oder ohne Server-Nachrichten wird nicht verglichen, nur Log `info`, kein Einfluss auf den Exit-Code | LH-FA-24.a | ja — Lastenheft zählt „fehlende Nachricht" zur Abweichung und kennt keine Ausnahme |
| S-36 | Mit `--compare-responses` bleibt ein Serverfehler mit gleichem SQLSTATE wie in der Aufzeichnung ein Fehler (`PGR-E4004`, Abbruch), außer mit `--allow-recorded-errors`; bei mehreren Ursachen gilt die zuerst aufgetretene | LH-FA-20.a Schritte 6 und 7 | ja — Lastenheft LH-FA-24 Happy Path (Erfolg bei entsprechenden Fehlern) und eigener Fehlerstatus |
| S-37 | Relative Pfade von `--tls-cert`, `--tls-key`, `--upstream-ca` gelten ab dem aktuellen Verzeichnis, auch aus der Konfigurationsdatei | LH-FA-20.a Schritt 2, LH-FA-23.a | nein — Lastenheft schweigt |
| S-38 | Ein verschlüsselter privater Schlüssel (mit Passwort) ist `PGR-E2007`; die Dateien werden beim Start geprüft, vor dem Öffnen eines Listen-Ports beziehungsweise vor der ersten Verbindung | LH-FA-23.a, LH-FA-20.a Schritt 2 | ja, leicht — Lastenheft: „Zertifikat und Schlüssel bereitgestellt", „ungültig"; Passwortschutz nicht erwähnt |
| S-39 | Bei Extended wird die Gesamtfolge der Server-Nachrichten verglichen, nicht gruppenweise | LH-FA-24.a | nein — Lastenheft: „Art und Reihenfolge der Antworten"; Folge siehe F-166 |
| S-40 | Ein Klartext-`CancelRequest` wird auch bei konfiguriertem TLS wie sonst behandelt (`PGR-W3001`), nicht als Ablehnung nach `PGR-E6003` | LH-FA-23.a, LH-FA-05.e | teilweise — Lastenheft: „lehnt die Verbindung ab"; die Verbindung wird geschlossen, aber ohne Exit-Code-Wirkung |

Die Setzungen S-28 (TLS 1.2), S-33 und S-34 des Vorgänger-Reports bleiben unverändert; S-28 überschreitet weiter das Lastenheft.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Hard Rule 3.4 — `spec/lastenheft.md` (LH-FA-24, Abnahmeszenario 16) | geprüft, ohne Befund; kein ADR-, Slice-, Wellen- oder Spec-Datei-Verweis |
| Hard Rule 3.4 — `spec/spezifikation.md` (geänderte Abschnitte) | geprüft, ohne Befund; kein Verweis nach unten (ADR, Slice, Welle) |
| Hard Rule 3.4 — `spec/architecture.md` | geprüft, ohne Befund; im Diff unverändert; Play-Sequenz „bis ReadyForQuery" je Interaktion deckt den Vergleich der Gesamtfolge |
| Hard Rule 3.5 — ADRs 0014 bis 0020 | geprüft, ohne Befund; keine ADR im Diff verändert, Index unverändert |
| Hard Rule 3.3 — Move und Inhaltsänderung | geprüft, ohne Befund; Diff enthält keine Umbenennung |
| Hard Rule 3.7 — Kommentare und Zustandsfelder in Slices, Handbuch-Kopf | geprüft, ohne Befund; Stand-Zeile nennt Datum, keine Chronik |
| Lastenheft ohne Wie — LH-FA-24 | geprüft, ohne neuen Befund; keine Nachrichtentypen oder Feldnamen; Protokollnähe einzelner Wörter ist F-156 |
| Handbuch ohne interne Kennungen | geprüft, ohne Befund; nur `PGR-…`, keine `LH-`, `SPEC-`, `ARC-`, `ADR-`, `slice-`, `welle-` |
| Handbuch als Ist-Zustand | geprüft, ohne Befund; Kopf („gilt ab der ersten veröffentlichten Version") unverändert; Handbuch und Spezifikation stimmen in `PGR-E2007`, `PGR-E5004`, Exit-Codes, Pfaden und Vergleichsumfang überein, Abweichungen siehe F-165, F-167 |
| Zählungen — Abnahmeszenarien (Lastenheft §7, Roadmap M3, Welle: „1 bis 10 und 12 bis 16") | geprüft, ohne Befund; unverändert |
| Zählungen — Optionen (Optionslisten LH-FA-01.a, 03.a, Tabelle LH-FA-17.a, Tabelle LH-FA-20.a, Handbuch) | geprüft, ohne Befund; `--tls-cert`, `--tls-key`, `--allow-plaintext`, `--upstream-ca`, `--compare-responses` überall gleich benannt, Geltung je Kommando stimmt |
| Zählungen — Codes (Codetabelle SPEC-034, SPEC-017 bis SPEC-019, Handbuch) | geprüft, ohne Befund; kein Code hinzugekommen oder entfallen; Klasse und Exit-Code von `PGR-E2007`, `PGR-E5004`, `PGR-E6003`, `PGR-W3001`, `PGR-W3002` stimmen |
| Exit-Code-Rangfolge Schritt 6, Schritt 7, Handbuch-Tabelle (F-141, F-154) | geprüft; Übereinstimmung der drei Stellen; Mängel siehe F-165, F-170 |
| Klartext-`CancelRequest`/`GSSENCRequest` bei TLS gegen LH-FA-05.e (F-145) | geprüft, ohne Widerspruch zu LH-FA-05.e; Wortlaut siehe F-171 |
| Zeitpunkt der `PGR-E2007`-Prüfung und Auflösung relativer Pfade (F-146) | geprüft; Konfigurationsdatei kennt Abschnitte je Kommando, die Aussage zu Pfaden aus der Datei hat damit einen Gegenstand |
| docs/plan/planning — drei Slices (Kopf, Liefer-Punkte, Trigger) | geprüft, ohne Befund bei Form; Inhalt siehe F-169 |
| Core-Reinheit, Gate-Lockerung (Hard Rule 3.2, 3.6), `.a-check.yml` | geprüft, ohne Befund; nicht im Diff |
| `make gates` | nicht ausgeführt (Auftrag) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 5 |
| INFO | 2 |

Von den 17 nachgeprüften Findings (F-141 bis F-155, F-159, F-160) sind 12 behoben (F-141, F-145, F-146, F-147, F-149, F-150, F-152, F-153, F-154, F-155, F-159, F-160; F-145, F-149, F-152 mit einem Rest), vier teilweise (F-142, F-143, F-144, F-148) und eines (F-151) mit bekanntem Rest. Reste und Folgewidersprüche sind F-164 bis F-172.

**Finding-Klassen dieses Laufs:** Spec entscheidet zwischen zwei Lastenheft-Aussagen ohne Setzung · Rangfolgeregel hebt zugesagte Unterscheidbarkeit auf · Meldungsinhalt setzt entfallene Struktur voraus · Ausnahmeregel kollidiert mit allgemeiner Fehlerregel · Spec vergleicht weniger als das Lastenheft nennt · DoD folgt Spec-Änderung nicht · Tie-Break nicht festgelegt.

Hinweis zum Zähler (Steering-Loop): „Fehlerbedingung ohne Definition" (zuletzt F-167) und „Spec überschreitet Lastenheft" (F-164, S-35, S-36) liegen weiter über der Schwelle von drei; der Folgeschritt aus dem Vorgänger-Report bleibt fällig. Neu: Korrekturen an einer Fehlerregel erzeugen Widersprüche an benachbarten Regeln (F-165, F-166, F-167, F-169, F-170), also ein viertes Mal „Korrektur zieht Nachbarstellen nicht nach".

## Verdikt

**Merge-blockierend:** ja — für F-164 bis F-167 (MEDIUM): Konflikt zwischen LH-FA-24 und LH-FA-20.a bei übereinstimmenden Fehlern (F-164), Exit-Code der zuerst aufgetretenen Ursache gegen den zugesagten eigenen Fehlerstatus (F-165), Gruppenindex ohne Gruppengrenzen (F-166), unvollständige Aufzeichnung gegen Fehlerregel und Lastenheft (F-167). Kein HIGH: keine ADR-Verletzung, kein Verstoß gegen Hard Rules 3.3 bis 3.7.

**Übergabe:** F-164, F-165, F-167 und die Setzungen S-35, S-36 an den Architect beziehungsweise den Auftraggeber (Vertragsrang, Entscheidung zwischen zwei Lastenheft-Aussagen); F-166, F-168, F-170, F-171, F-172 an den Implementer (Autor der Spec); F-169 an den Planner; F-173 an den Architect.
