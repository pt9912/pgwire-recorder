# Review-Report: tls-client-antwortvergleich, vierte Nachprüfung — 2026-10-04

**Review-Art:** Nachprüfung der Korrekturen zu `docs/reviews/2026-10-04-review-tls-client-antwortvergleich-nachpruefung-3.md` (F-187 bis F-195), Konsistenzprüfung Lastenheft, Spezifikation, ADRs, Handbuch, Slices und README, Prüfung der Hard Rules einschließlich der neuen Regel 3.8.

**Gegenstand:** `git diff 090f83c HEAD` ohne `docs/reviews/`: `spec/lastenheft.md` (LH-FA-24 Boundary), `spec/spezifikation.md` (LH-FA-20.a Schritte 3 bis 7 und Tabelle „Fehlerregeln beim Einspielen"; LH-FA-23.a Zeilenumbruch; LH-FA-24.a „Unvollständige Aufzeichnung" und „Meldung"; Codetabelle `PGR-E5004`), `docs/user/benutzerhandbuch.md`, `README.md`, beide Slices (Antwortvergleich, Einspielen), ADR-Index, [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) (Status), [ADR-0022](../plan/adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md) (jetzt Superseded), neue [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) (Accepted), `AGENTS.md` (neue Hard Rule 3.8). `spec/architecture.md` und [ADR-0014](../plan/adr/0014-konfigurationsdatei.md) bis [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) sind im Diff unverändert.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:** `AGENTS.md` (Hard Rules), `harness/README.md`, `harness/conventions.md`; der Vorgänger-Report; [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md), [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md); berührte `LH-*`: LH-FA-20, -23, -24, LH-QA-05. Dieser Lauf setzt mit F-196 und S-45 fort. Nicht Gegenstand: DoD-Konformität (Verifier), fachliche Eignung (Validator). `make gates` wurde auftragsgemäß nicht ausgeführt; das d-check-Verhalten ist gelesen, nicht gelaufen. Lastenheft-Status `Draft`.

---

## Stand der Findings des Vorgänger-Reports

| ID | Stand | Anmerkung |
|---|---|---|
| F-187 | behoben, mit Rest | Lastenheft (Boundary), LH-FA-24.a, Fehlerregeln-Tabelle, `PGR-E5004`, Handbuch, Slice und [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) sagen übereinstimmend: eine Fehlerantwort ohne Vorbild ist stets eine Abweichung, auch nach dem aufgezeichneten Ende und ohne aufgezeichnete Server-Nachrichten. Rest: ein `FATAL` ohne Vorbild, F-196. |
| F-188 | teilweise | Der Fall „aufgezeichnetes Verbindungsende, Server verursacht es mit gleichem SQLSTATE" ist in Lastenheft, Tabellenzeile 3, LH-FA-24.a, Handbuch und Slice geregelt (erwartet, Session endet, nächste läuft). Offene Ränder: F-196, F-197. |
| F-189 | behoben | Schritte 5 und 7 verweisen auf die Tabelle; der Signalfall hat eine Tabellenzeile; die Kopfzeile nennt Tabelle und Absatz als Quelle und sagt, die Schritte wiederholen nicht. |
| F-190 | teilweise | Schritt 6 sagt: bei Abbruch keine weiteren Nachrichten, auch keine weiteren Gruppen. Das Zusammenspiel mit Schritt 4 (Antworten bis Interaktionsende gehalten) bleibt offen, F-201. |
| F-191 | behoben, mit Rest | Slice Antwortvergleich trägt Extended-Fortsetzung bis `Sync`; Slice Einspielen trägt das abgelaufene oder ungültige Serverzertifikat (`PGR-E4005`). Neue Zusagen ohne benannten Test: F-200. |
| F-192 | behoben, mit Rest | Die Zusammenfassung ist eine Log-Zeile der Stufe `info` ohne Meldungscode, bei jedem Ende des Laufs, mit festen Attributen und Zählregel für `verglichen`; konsistent in Spezifikation, Handbuch und Slice. Rest: F-198. |
| F-193 | behoben | README: „wertet Fehlerantworten des Servers und vergleicht auf Wunsch die Struktur der Antworten". |
| F-194 | behoben | Index-Zeilen 0022 und 0023 nennen die Teilersetzung von [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und die Ergänzung von [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md); die Beziehung ist aus dem Index lesbar. Neuer Rest in der Statuskette: F-199. |
| F-195 | teilweise | LH-FA-23.a ist umgebrochen. Im Handbuch sind die Zeilen um die Log-Zeile und um „mit `--compare-responses` gilt stattdessen der Vergleich" weiter unregelmäßig lang; „Anfragen" (Handbuch) gegen „Interaktionen" (Spezifikation) besteht fort. Siehe F-203. |
| F-156, F-157, F-158, F-161 | offen, unverändert | bekannte Punkte, nicht neu bewertet |
| F-162, F-163 | INFO, unverändert | |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-196 | MEDIUM | Ein `FATAL`, das die Aufzeichnung nicht oder mit anderem SQLSTATE enthält, ist an zwei Stellen verschieden geregelt. LH-FA-24.a und das Lastenheft: eine `ErrorResponse` ohne Vorbild ist stets eine Abweichung (`PGR-E5004`, Exit-Code `5`). Tabellenzeile 3 („Verbindung bricht nach dem ersten `ReadyForQuery` ab, auch durch einen `FATAL`"): mit Vergleich „gleich", also `PGR-E4003`, Exit-Code `4`, immer sofort; nur das aufgezeichnete Verbindungsende mit gleichem SQLSTATE ist die Ausnahme. Welche Zeile für einen fremden `FATAL` gilt (Zeile 3 oder 4/5) und welcher Exit-Code folgt, steht nirgends. Das ist das dritte Auftreten von Allgemeinregel gegen Ausnahmeregel in dieser Reihe. | LH-FA-24; LH-FA-20; LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Fehlerregeln-Tabelle (Zeilen 3 bis 5); LH-FA-24.a „Unvollständige Aufzeichnung"; spec/lastenheft.md · LH-FA-24 Boundary und Negative | nein | Ausnahmeregel und Allgemeinregel widersprüchlich (drittes Auftreten, siehe F-177, F-187) |
| F-197 | MEDIUM | Das erwartete Verbindungsende ist an den Rändern unbestimmt. (a) „Gleicher SQLSTATE" setzt einen aufgezeichneten Fehler voraus; wie ein aufgezeichnetes Verbindungsende ohne `ErrorResponse` (Server beendet die Verbindung ohne Fehler) verglichen wird, ist nicht gesagt. (b) Schritt 3 wartet bei einer Interaktion ohne `ReadyForQuery` auf das aufgezeichnete Verbindungsende; sendet der Server stattdessen ein `ReadyForQuery` (LH-FA-24.a: „keine Abweichung") oder beendet er die Verbindung nicht, steht nicht da, ob und wie lange der Recorder wartet und wie die Interaktion endet. (c) Ohne `--compare-responses` ist dasselbe Verbindungsende stets `PGR-E4003`, Exit-Code `4`; das Lastenheft gibt der Erwartung nur im Vergleichsfall Raum, die Spezifikation sagt es aber nur implizit mit „gleich". | LH-FA-24; LH-FA-20; LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Schritt 3, Tabellenzeile 3; LH-FA-24.a „Unvollständige Aufzeichnung" | nein | Unklare Fehlerbehandlung am Rand: erwartetes Verbindungsende |
| F-198 | LOW | Die Log-Zeile „bei jedem Ende des Laufs" ist an drei Stellen nicht bestimmt: Sie ist unter `--log-level warn` oder `error` nicht sichtbar, was der DoD („steht die Log-Zeile") nicht erwähnt; ob sie bei einem Ende vor dem ersten Einspielen erscheint (Aufzeichnung ungültig, Startfehler `PGR-E2007`, Exit-Code `3`/`2`), ist offen; ohne `--compare-responses` ist `verglichen` nicht definiert (null oder Attribut fehlt). Das Format der Attribute (Schreibweise, Ausgabeform) hat keine Entsprechung in `SPEC-005`/`SPEC-006`. | LH-QA-05; Maintainability | spec/spezifikation.md · LH-FA-24.a „Meldung"; docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md · Liefer-Punkt 3 | nein | Neue Ausgabe ohne Randfälle |
| F-199 | LOW | Statuskette der ADRs: Die Datei [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) trägt `Superseded by` [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md), der Index (Spalte Status) für 0021 weiter `Superseded by` 0022, die Geschichte von 0021 ebenfalls nur „Superseded by [ADR-0022](../plan/adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md)", und die Geschichte von [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) nennt nur 0022 als Ersetzte. Drei Fundstellen sagen für 0021 zwei verschiedene Nachfolger. | Maintainability; AGENTS.md §3.5 | docs/plan/adr/0021 · Status gegen Geschichte; docs/plan/adr/README.md · Tabelle, Zeile 0021 | ja (Lesen) | Statusangabe der abgelösten ADR nicht mit Index und Geschichte abgeglichen |
| F-200 | LOW | Neue Zusagen ohne benannten Test. Der Liefer-Punkt 1 des Slice für den Antwortvergleich bündelt sieben Zusagen in einem Satz (Erfolg, `PGR-E5004`, Fehler ohne Vorbild, unvollständige Aufzeichnung, erwartetes Verbindungsende, Weitersenden bis `Sync`). Nicht genannt sind: Abbruch ohne weitere Nachrichten und Gruppen (Schritt 6), Verbindungsende mit anderem SQLSTATE, Interaktion ohne `ReadyForQuery` ohne Vergleich (Schritt 3). | LH-QA-05; Maintainability | docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md · Liefer-Punkt 1; docs/plan/planning/open/slice-v1-abschluss-einspielen.md · Liefer-Punkt 2 | nein | Spec-Zusage ohne Test an benanntem Ort |
| F-201 | LOW | Zeitpunkt der Entscheidung bei Vergleich unbestimmt. Schritt 4 hält die Antworten bis zum Ende der Interaktion; die Tabelle sagt für Fehler ohne Vorbild und Abweichung „ohne `--continue-on-error` sofort", Schritt 6 „keine weiteren Nachrichten der Interaktion". Wird eine Abweichung erst am Interaktionsende erkannt, ist die Aussage „keine weiteren Gruppen" für sie leer; ob eine `ErrorResponse` mitten in einer Extended-Interaktion sofort (vor den weiteren Gruppen) entscheidet, ist nicht gesagt. Zudem ist in der Signalzeile „nach einem Verbindungsfehler `4`" nicht erreichbar, weil ein Verbindungsfehler immer abbricht und kein Signal mehr abgewartet wird. | LH-FA-20; LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Schritt 4, 6; Tabelle (Spalte Abbruch, Zeile Abbruchsignal) | nein | Abbruchzeitpunkt bei mehrstufiger Interaktion nicht bestimmt (zweites Auftreten, siehe F-190) |
| F-202 | INFO | Hard Rule 3.8 und Bestand: Die angenommenen [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0019](../plan/adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) und [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) führen in „Entscheidung" Regeldetails (Nachrichtenfolge, Exit-Codes, SQLSTATE-Regeln), die die neue Regel dem Spezifikationsstratum zuweist. Sie sind unveränderlich (3.5) und vor der Regel entstanden; kein Fehler der Korrektur. Die Regel selbst ist klar, steht im Abschnitt der Hard Rules hinter 3.7 und begründet sich aus 3.5. Grenzfall in [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md): „Unvollständige Aufzeichnungen werden verglichen, soweit sie reichen" ist eine Regelaussage, die Spezifikation und Lastenheft stützen; sie bleibt als Entscheidungsumfang vertretbar. | AGENTS.md §3.8, §3.5 | AGENTS.md · 3.8; docs/plan/adr/0016, 0017, 0019, 0020 · Entscheidung; docs/plan/adr/0023 · Entscheidung | nein | Regeldetails in Bestands-ADR (Hinweis) |
| F-203 | INFO | Zeilenumbruch und Wortwahl: Im Handbuch (Einspielen, Hinweise) steht eine überlange Zeile um die Log-Zeile („eingespielten, verglichenen und abweichenden Anfragen") und eine unregelmäßig umbrochene Klammer bei `PGR-E4004`; das Handbuch zählt „Anfragen", die Spezifikation „Interaktionen"; im Lastenheft bleibt Abnahmeszenario 16 unverändert und deckt Fehler ohne Vorbild nach dem Präfix und das erwartete Verbindungsende nicht ab (Szenario ist Beispiel, kein Vollständigkeitsanspruch). Keine semantische Wirkung. | Maintainability | docs/user/benutzerhandbuch.md · Einspielen, Hinweise; spec/lastenheft.md · Abnahmeszenario 16 | nein | Zeilenumbruch nach Einfügung unregelmäßig |

## Setzungen

Die früheren Setzungen S-1 bis S-44 bleiben unbestätigt, soweit nicht unten anders vermerkt.

| Nr. | Stand / Setzung | Fundort | Überschreitet das Lastenheft? |
|---|---|---|---|
| S-43 | erledigt: eine Fehlerantwort ohne Vorbild ist auch nach dem aufgezeichneten Ende eine Abweichung; Lastenheft und Spezifikation stimmen überein | LH-FA-24.a, Lastenheft LH-FA-24 | nein |
| S-44 | geändert: die Zusammenfassung ist keine Ausgabe auf `stderr` mehr, sondern eine Log-Zeile der Stufe `info` ohne Meldungscode, bei jedem Ende des Laufs; Randfälle siehe F-198 | LH-FA-24.a „Meldung" | ja — das Lastenheft verlangt keine Zusammenfassung |
| S-45 | Ein erwartetes Verbindungsende verlangt den gleichen SQLSTATE; die Session endet, die nächste läuft weiter, und das Einspielen zählt nicht als Fehler | LH-FA-20.a Tabellenzeile 3, LH-FA-24.a | teilweise — das Lastenheft sagt „genauso verursacht … gilt als erwartet"; den SQLSTATE-Vergleich und die Fortsetzung mit der nächsten Session ergänzt die Spezifikation |
| S-46 | Ohne `--compare-responses` ist ein aufgezeichnetes Verbindungsende, das der Server wiederholt, stets `PGR-E4003` mit Exit-Code `4` | LH-FA-20.a Tabellenzeile 3 | nein — das Lastenheft nennt die Erwartung nur für den Vergleich |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Hard Rule 3.3 — Move und Inhaltsänderung | geprüft, ohne Befund; der Diff enthält keine Umbenennung |
| Hard Rule 3.4 — `spec/lastenheft.md` | geprüft, ohne Befund; kein Verweis auf ADR, Slice, Welle oder Spec-Datei; Aussagen ohne Wie |
| Hard Rule 3.4 — `spec/spezifikation.md` (geänderte Abschnitte) | geprüft, ohne Befund; kein Verweis nach unten (ADR, Slice, Welle, Commit) |
| Hard Rule 3.4 — `spec/architecture.md` | geprüft, ohne Befund; im Diff unverändert |
| Hard Rule 3.5 — [ADR-0014](../plan/adr/0014-konfigurationsdatei.md) bis [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) | geprüft, ohne Befund; im Diff unverändert |
| Hard Rule 3.5 — [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) und [ADR-0022](../plan/adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md) | geprüft, ohne Befund zur Immutabilität; geändert sind bei 0021 nur die Statuszeile, bei 0022 Statuszeile und eine Geschichtszeile; Statuskette siehe F-199 |
| Hard Rule 3.5 — [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) | geprüft, ohne Befund; neue Datei, Kopf, Alternativen (A bis D), Konsequenzen, Fitness Function, Trigger, Geschichte vollständig |
| Hard Rule 3.7 — Kommentare und Zustandsfelder | geprüft, ohne Befund; Indikativ, keine Chronik in Zustandsfeldern; die Slice-DoDs nennen Zustand, keine Historie |
| Hard Rule 3.8 — [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) gegen die Regel | geprüft, ohne Befund; Entscheidung, Alternativen, Gründe, Zeiger auf die Spezifikation; Grenzfall siehe F-202 |
| [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) — Aussagen gegen die Spezifikation | geprüft, ohne Befund; der Vergleich entscheidet über Fehlerantworten (Tabelle), eigener Code und Exit-Code `5` (`PGR-E5004`), Vergleich der Struktur, nicht der Werte (LH-FA-24), unvollständige Aufzeichnungen (LH-FA-24.a); keine Aussage ohne Stütze |
| [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) — Schärft-Anker | geprüft, ohne Befund; `#lh-fa-24a--vergleich-der-antworten` und `#lh-fa-20a--einspielen` lösen auf die Überschriften „LH-FA-24.a — Vergleich der Antworten" und „LH-FA-20.a — Einspielen" auf; die Bezug-Anker des Lastenhefts sind unverändert |
| ADR-Index — Tabelle und Beziehungen | geprüft; Nummerierung lückenlos bis 0023, 0022 `Superseded by 0023`, 0023 `Accepted`, Beziehungszeilen 0022 und 0023 konsistent; Abweichung nur bei 0021, siehe F-199 |
| d-check: Referenzen auf inaktive [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) und [ADR-0022](../plan/adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md) | geprüft, ohne Befund (gelesen, Gate nicht gelaufen); außerhalb der Reviews stehen sie nur im Index, in den ADR-Dateien selbst und in der Geschichte; beide Slices verweisen jetzt auf [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md); Spezifikation, Handbuch und README verweisen nicht auf ADRs |
| Lastenheft LH-FA-20/LH-FA-24 gegen Spezifikation | geprüft; Fehler ohne Vorbild stets Abweichung, erwartetes Verbindungsende (Session endet, nächste läuft), Abbruch ohne weitere Nachrichten, Exit-Code-Rangfolge einschließlich Signal stimmen überein; Ausnahme: F-196 |
| Fehlerregeln-Tabelle gegen Schritte 3 bis 7, `SPEC-017`, `SPEC-018`, `SPEC-022`, `SPEC-027`, `SPEC-034` | geprüft; Codes und Exit-Codes stimmen überein (`PGR-E4002` bis `PGR-E4005`, `PGR-E5004`); `SPEC-017` und `SPEC-022` nennen `PGR-E4004` nur ohne Vergleich; Abweichungen siehe F-196, F-197, F-201 |
| Handbuch gegen Spezifikation und Tabellen | geprüft; Fehler ohne Vorbild, erwartetes Verbindungsende, Log-Zeile, `--allow-recorded-errors` ohne Wirkung mit Vergleich, Fehlertabelle (`PGR-E4003`, `PGR-E4004`, `PGR-E5004`) stimmen überein; keine internen Kennungen (`LH-`, `SPEC-`, `ARC-`, `ADR-`, `slice-`) |
| README | geprüft, ohne Befund; Aussage zu `play` stimmt mit der Spezifikation überein |
| Slice-DoDs | geprüft; höchstens drei Liefer-Punkte je Slice, Bezug nennt [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md); Inhalt stimmt mit der Spezifikation überein, Lücken siehe F-200 |
| Zählungen — Abnahmeszenarien (16), Optionen, Codes | geprüft, ohne Befund; keine Option und kein Code hinzugekommen oder entfallen; `PGR-E5004`-Zeile erweitert, Klasse und Exit-Code unverändert |
| Zählungen — ADR-Index | geprüft, ohne Befund bis auf F-199 |
| Core-Reinheit, Gate-Lockerung (Hard Rule 3.2, 3.6), `.a-check.yml` | geprüft, ohne Befund; nicht im Diff |
| F-156, F-157, F-158, F-161 | nicht neu bewertet (Auftrag) |
| `make gates` | nicht ausgeführt (Auftrag) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

Von den neun nachgeprüften Findings (F-187 bis F-195) sind F-189, F-193 und F-194 behoben; F-187, F-191 und F-192 behoben mit Rest; F-188, F-190 und F-195 teilweise. Neu: F-196 bis F-203.

**Finding-Klassen dieses Laufs:** Ausnahmeregel und Allgemeinregel widersprüchlich (drittes Auftreten) · Unklare Fehlerbehandlung am Rand (erwartetes Verbindungsende) · Neue Ausgabe ohne Randfälle · Statusangabe der abgelösten ADR nicht mit Index und Geschichte abgeglichen · Spec-Zusage ohne Test an benanntem Ort · Abbruchzeitpunkt bei mehrstufiger Interaktion nicht bestimmt (zweites Auftreten) · Regeldetails in Bestands-ADR.

Hinweis zum Zähler (Steering-Loop): „Ausnahmeregel und Allgemeinregel widersprüchlich" ist mit F-177, F-187 und F-196 dreimal belegt; damit ist der Schwellenwert der Pflege erreicht. Die Klassifikation (MEDIUM) bleibt richtig. Die Spezifikation führt die Fehlerregeln jetzt in einer Tabelle, deren Zeilen sich aber nicht gegenseitig ausschließen (Verbindungsende, Fehlerantwort, Abweichung eines `FATAL`); ein Ausschluss-Prinzip („erste zutreffende Zeile" oder eine Vorrangregel) fehlt. Für ein ADR- oder `AGENTS.md`-Update gibt es keinen Anlass: Hard Rule 3.8 verhindert das ADR-Seitige der Drift bereits. Als Fitness Function gibt es weiterhin keinen Sensor.

## Verdikt

**Merge-blockierend:** ja, für die MEDIUM-Findings F-196 und F-197 (Verhalten bei fremdem `FATAL` und an den Rändern des erwarteten Verbindungsendes: Widerspruch zwischen Tabellenzeile 3 und LH-FA-24.a beziehungsweise Lastenheft, Wartezeit unbestimmt). Kein HIGH: die Hard Rules 3.3 bis 3.8 sind gewahrt, [ADR-0014](../plan/adr/0014-konfigurationsdatei.md) bis [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) sind unverändert, [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) und [ADR-0022](../plan/adr/0022-vergleichsregeln-beim-einspielen-praezisiert.md) nur in Status und Geschichte geändert, [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md) trägt nur Entscheidung und Gründe. Die Hauptlinie (Fehler ohne Vorbild stets Abweichung, Weiterlauf nach erwartetem Verbindungsende, Exit-Code-Rangfolge einschließlich Signal, Index und ADR-Kette bis auf 0021) ist über Lastenheft, Spezifikation, Handbuch, Slices und README widerspruchsfrei.

**Übergabe:** F-196, F-197, F-198, F-201 an den Implementer (Autor der Spec; bei F-196 vorab Entscheidung, welche Zeile für einen fremden `FATAL` gilt, über den Architect); F-199 an den Architect (Statuskette; eine inhaltliche Änderung an einer Accepted-ADR wäre eine neue ADR, die Statuszeile und Geschichte von 0021 sind davon nicht berührt); F-200 an den Planner; F-202 zur Kenntnis (Bestand), F-203 an den Implementer; Setzungen S-44 bis S-46 an den Auftraggeber beziehungsweise Architect.
