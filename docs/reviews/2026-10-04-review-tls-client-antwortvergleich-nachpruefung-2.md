# Review-Report: tls-client-antwortvergleich, zweite Nachprüfung — 2026-10-04

**Review-Art:** Nachprüfung der Korrekturen zu `docs/reviews/2026-10-04-review-tls-client-antwortvergleich-nachpruefung.md` (F-164 bis F-174), Konsistenzprüfung Lastenheft, Spezifikation, Architektur, ADRs, Handbuch und Slices, Prüfung der Hard Rules.

**Gegenstand:** `git diff ca35c70 HEAD` ohne `docs/reviews/` (HEAD 972e036; Commits ec08fbc, 2bc1751, 5ca9ae0, 972e036): `spec/lastenheft.md` (LH-FA-24), `spec/spezifikation.md` (LH-FA-20.a Schritte 2, 6, 7; LH-FA-23.a „Unverschlüsselte Clients"; LH-FA-24.a), `docs/user/benutzerhandbuch.md`, `slice-v1-abschluss-antwortvergleich`, `slice-v1-abschluss-einspielen`, neue [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) und der ADR-Index. `spec/architecture.md` und die ADR-Dateien mit den Nummern 0001 bis 0020 sind im Diff unverändert.

**Skill:** `.harness/skills/reviewer.md` @ 972e036
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:** `AGENTS.md` (Hard Rules), `harness/README.md`, `harness/conventions.md`; der Vorgänger-Report; [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md), [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md); berührte `LH-*`: LH-FA-05, -20, -23, -24, LH-QA-05. Dieser Lauf setzt mit F-175 und S-41 fort. Nicht Gegenstand: DoD-Konformität (Verifier), fachliche Eignung (Validator). `make gates` wurde auftragsgemäß nicht ausgeführt. Lastenheft-Status `Draft`.

---

## Stand der Findings des Vorgänger-Reports

| ID | Stand | Anmerkung |
|---|---|---|
| F-164 | behoben | Lastenheft (LH-FA-24 Beschreibung, Boundary) sagt: mit Vergleich gilt ein Fehler wie aufgezeichnet als erwartet, der Wunsch ist die ausdrückliche Wahl nach LH-FA-20; Schritt 6 und Handbuch folgen. Die Entscheidung steht jetzt im Vertragsrang. Setzung S-36 entfällt. |
| F-165 | behoben | Mit Vergleich ist der Exit-Code am Ende `5`, ohne Vergleich `4`; ein Verbindungsfehler beendet immer sofort mit `4`. `PGR-E4004` kommt mit Vergleich nicht mehr vor, damit gibt es keine „zuerst aufgetretene Ursache" mehr zwischen `4` und `5`. Schritt 6, Schritt 7, Handbuch und Slice stimmen überein. |
| F-166 | behoben | Der Gruppenindex ist aus der Meldung entfernt, sie nennt die Position in der normalisierten Folge. Handbuch („Sitzung, Anfrage, Art") und Slice („Session, Sequenznummer und Art") stimmen. |
| F-167 | behoben, mit Rest | Die Ausnahme „nicht verglichen" ist gestrichen; eine unvollständige Aufzeichnung wird verglichen, soweit sie reicht, und das Lastenheft trägt die Regel (Boundary). Setzung S-35 entfällt. Rest: F-177. |
| F-168 | behoben | `command_complete` vergleicht den Tag ohne abschließende Zahlen; `CREATE TABLE` und `CREATE INDEX` sind unterscheidbar. Slice und Meldungsart („Befehl") stimmen. Setzung S-31 im Kern erledigt. |
| F-169 | behoben, mit Rest | `slice-v1-abschluss-einspielen` nennt „sonst 4" und verweist für die Rangfolge mit Vergleich auf den Slice für den Antwortvergleich. Rest: F-182. |
| F-170 | behoben | „Je Interaktion wird die erste Abweichung in der Reihenfolge der Nachrichten gemeldet"; `PGR-E4004` entfällt mit Vergleich, ein Tie-Break zwischen den Codes ist nicht mehr nötig. |
| F-171 | behoben | „keine Anfrage- und keine Antwortdaten" ist sachlich zutreffend. |
| F-172 | teilweise | Subjekt in Schritt 2 ist wieder der Recorder (Satz zur Authentifizierung steht vor `--upstream-ca`). Der ungleichmäßige Zeilenumbruch in LH-FA-23.a und Abnahmeszenario 16 bleibt, siehe F-184 (INFO). |
| F-173 | behoben, mit Rest | [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) trägt Vorrang, Gesamtfolge, Präfix-Vergleich und Befehl ohne Zahlen; [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) bleibt unverändert. Reste: F-178, F-179. Die `parameter_description` nennt [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) weiter nicht; in LH-FA-24.a und im Lastenheft steht sie, nicht neu bewertet. |
| F-174 | teilweise | Der Zeitpunkt ist gesagt (Start: Datei lesbar und PEM; Verbinden: Ablauf, `PGR-E4005`). Neuer Widerspruch in der Codetabelle: F-181. Ablauf während eines langen Laufs bleibt INFO: F-185. |
| F-156, F-157, F-158, F-161 | offen, unverändert | bleiben bekannte Punkte; die Korrekturen ändern ihre Lage nicht. F-157 und F-161 (ADR-Träger beziehungsweise Spec-Folgen) sind durch [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) für den Antwortvergleich berührt, aber nicht neu bewertet. |
| F-162, F-163 | INFO, unverändert | |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-175 | MEDIUM | Die Wechselwirkung `--allow-recorded-errors` mit `--compare-responses` ist nicht geregelt. Schritt 6 sagt, der Vergleich „ersetzt diese Regel für Fehlerantworten", ohne `--allow-recorded-errors` zu nennen; die Optionstabelle sagt weiter „ein Serverfehler in einer Interaktion, deren Aufzeichnung ebenfalls eine `ErrorResponse` enthält, gilt als erwartet" (ohne Bedingung SQLSTATE), das Handbuch ebenso. Daraus folgt für eine aufgezeichnete Interaktion mit Fehler A und einen Server, der Fehler B liefert: Abweichung (Vergleich) oder erwartet (Option). [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) nennt die Option „mit Vergleich überflüssig"; das steht weder in der Spezifikation noch im Handbuch noch in der Folgepflicht. Der Slice für den Antwortvergleich verlangt „wirken wie spezifiziert" für beide Optionen. | LH-FA-20; LH-FA-24; [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) | spec/spezifikation.md · LH-FA-20.a Optionstabelle, Schritt 6; docs/user/benutzerhandbuch.md · Einspielen; docs/plan/adr/0021 · Verglichene Alternativen (Option D) | nein | Optionswechselwirkung nicht geregelt |
| F-176 | MEDIUM | Schritt 6 ersetzt mit Vergleich „diese Regel" für Fehlerantworten. Zur Regel gehörte der Satz, dass bei Extended der Server nach dem Fehler bis zum `Sync` verwirft und der Recorder „weitersendet wie aufgezeichnet". Ob das für einen erwarteten Fehler mit Vergleich weiter gilt, ist nicht gesagt, ebenso, was nach einer Abweichung mit `--continue-on-error` bei Extended gesendet wird (Rest der Interaktion oder nächste Interaktion). Ein erwarteter Fehler „bricht nicht ab", sein Weiterlauf ist nicht bestimmt. | LH-FA-20; LH-FA-24 | spec/spezifikation.md · LH-FA-20.a Schritt 6 | nein | Ersetzte Regel nimmt Teilregel ungewollt mit |
| F-177 | MEDIUM | Vergleich einer unvollständigen Aufzeichnung (LH-FA-24.a „Unvollständige Aufzeichnung"): „weitere Nachrichten des Servers danach sind keine Abweichung" und „eine `ErrorResponse`, die die Aufzeichnung nicht enthält, bleibt eine Abweichung" widersprechen sich, wenn die Fehlerantwort nach dem aufgezeichneten Präfix kommt. Im Fall ohne Server-Nachrichten ist nur der zweite Satz anwendbar (jeder Serverfehler ist Abweichung, auch bei Aufzeichnungsende durch Verbindungsabbruch), der erste (jede Antwort ist zulässig) sonst; welcher gilt, ist nicht bestimmt. Auch unklar, ob für den Präfix die Normalisierung (ohne `data_row`) zuerst gilt und ob „soweit sie reicht" bei einer Aufzeichnung mit `FATAL` als letzter Nachricht die Folge danach (`ready_for_query` des Servers) ausschließt. | LH-FA-24 | spec/spezifikation.md · LH-FA-24.a „Unvollständige Aufzeichnung"; spec/lastenheft.md · LH-FA-24 Boundary | nein | Ausnahmeregel und Allgemeinregel im selben Absatz widersprüchlich |
| F-178 | MEDIUM | [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) „ersetzt die Fehlerregel von [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) für Fehlerantworten" mit Vergleich, trägt aber weder ein `Supersedes`/`Ergänzt`-Feld noch einen Vermerk im Index; die ADR-Konventionen sagen „Schärfungen entstehen als neue ADR mit `Supersedes`". [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) (beide Accepted, unverändert, Hard Rule 3.5 gewahrt) tragen keinen Rückzeiger; [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md) sagt weiter „folgt der Regel für Fehlerantworten ([ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md))", was für einen Leser nur dieser ADR der Entscheidung widerspricht. Dass 0021 eine Teilersetzung für einen Modus ist und die beiden anderen sonst gelten, steht im Text, aber nicht in einem benennbaren Feld. | AGENTS.md §3.5; ADR-Index Konventionen | docs/plan/adr/0021 · Kopf und Entscheidung; docs/plan/adr/README.md · Konventionen | nein | Teilersetzung ohne maschinenlesbare Beziehung |
| F-179 | LOW | [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) behauptet Dinge, die die Spezifikation nicht trägt oder anders sagt: (a) „Negativ: Mit Vergleich gilt jeder aufgezeichnete Fehler als erwartet" — die Spezifikation verlangt gleichen SQLSTATE an der Stelle; (b) „Ein Lauf mit Vergleich endet mit 0, 4 oder 5" — Start- und Aufzeichnungsfehler (`2`, `3`) und `6` bleiben möglich; (c) Schärft-Feld nennt `SPEC-034`, aber die Zeile `PGR-E4004` dort trägt keine Einschränkung (siehe F-180); (d) Folgepflicht „Tests decken … abgebrochene Aufzeichnungen" und Fitness-Function-Zeile „die Tests des Slice für den Antwortvergleich tragen die Regeln" werden vom DoD des Slice nur teilweise getragen (kein Fall „anderer SQLSTATE bei aufgezeichnetem Fehler", kein Fall ohne Server-Nachrichten, siehe F-182). | [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) | docs/plan/adr/0021 · Konsequenzen, Schärft, Fitness Function | nein | ADR sagt mehr oder anderes als die Spezifikation |
| F-180 | LOW | `PGR-E4004` entfällt laut Schritt 6 mit `--compare-responses`; die Codetabelle (`SPEC-017`, `SPEC-022`, `SPEC-034`) und die Handbuch-Zeile zu `PGR-E4004` nennen den Code uneingeschränkt („Server beantwortet eine eingespielte Anfrage mit einem Fehler"). `SPEC-018` und `PGR-E5004` tragen die Wechselwirkung, die Gegenstelle nicht. | LH-QA-05; [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) | spec/spezifikation.md · §4 Codetabelle `PGR-E4004`, `SPEC-017`, `SPEC-022`; docs/user/benutzerhandbuch.md · Tabelle `PGR-E4004` | nein | Codetabelle trägt Ausnahme des Codes nicht |
| F-181 | MEDIUM | Zeitpunkt der Zertifikatsprüfung: Schritt 2 sagt, ein Ablauf eines Zertifikats der Datei (`--upstream-ca`) oder des Servers wird erst beim Verbinden geprüft und ist `PGR-E4005` (Exit-Code `4`). Die Codetabelle sagt für `PGR-E2007` (Exit-Code `2`): „Zertifikat, Schlüssel oder Zertifizierungsstelle … ungültig (auch abgelaufen)". Für eine abgelaufene Zertifizierungsstelle widersprechen sich Exit-Code und Zeitpunkt. Für das eigene Zertifikat (`--tls-cert`) gilt weiter Start und `PGR-E2007` (LH-FA-23.a), das ist stimmig. Außerdem nennt die Handbuch-Zeile `PGR-E4005` das abgelaufene oder ungültige Zertifikat nicht (nur die Zeile `PGR-E2007` verweist darauf), während Schritt 5 und `SPEC-034` „ungültiges Zertifikat" für `PGR-E4005` führen. | LH-FA-23; LH-FA-20; LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Schritt 2 und Schritt 5, §4 `PGR-E2007` und `PGR-E4005`; docs/user/benutzerhandbuch.md · Tabelle `PGR-E4005`, `PGR-E2007` | nein | Codebeschreibung folgt Zeitpunktregel nicht |
| F-182 | MEDIUM | `slice-v1-abschluss-einspielen` verweist für „die Rangfolge mit Vergleich (Exit-Code 5)" bei einem Abbruchsignal auf `slice-v1-abschluss-antwortvergleich`; dessen DoD enthält keinen Test dazu (kein Abbruchsignal, kein Fall `--allow-recorded-errors` mit Vergleich, kein Fall ohne Server-Nachrichten oder Fehlerantwort nach dem Präfix, kein Fall Verbindungsfehler nach einer Abweichung mit `--continue-on-error`, der mit `4` endet). Neuer öffentlicher Vertrag (Schritt 7, Schritt 6) ohne Negativtest an einem benannten Ort. | Maintainability; LH-QA-05 | docs/plan/planning/open/slice-v1-abschluss-einspielen.md · dritter Liefer-Punkt; docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md · Liefer-Punkte | nein | Verweis auf Test, der nirgends steht |
| F-183 | LOW | Kopf und Bezug von `slice-v1-abschluss-antwortvergleich` führen [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) nicht, obwohl der Slice dessen Regeln trägt und die ADR ihn als Träger der Tests nennt; `slice-v1-abschluss-einspielen` nennt die Entscheidung ebenfalls nicht. | Maintainability | docs/plan/planning/open/slice-v1-abschluss-antwortvergleich.md · Kopf (Bezug); docs/plan/planning/open/slice-v1-abschluss-einspielen.md · Kopf | nein | Slice-Bezug folgt neuer ADR nicht |
| F-184 | INFO | Rest von F-172: In LH-FA-23.a (Satz zum `CancelRequest`) und im Lastenheft Abnahmeszenario 16 („… unterscheidet. Bezug: LH-FA-24.") sind Zeilenumbrüche ungleichmäßig; keine semantische Wirkung. | Maintainability | spec/spezifikation.md · LH-FA-23.a; spec/lastenheft.md · Abnahmeszenario 16 | nein | Zeilenumbruch nach Einfügung unregelmäßig |
| F-185 | INFO | Läuft ein Zertifikat während eines langen Laufs ab: für `play` ist es mit der Prüfung je Verbindungsaufbau gesagt (jede Session eigene Verbindung); für das eigene Zertifikat von `record` und `replay` (`--tls-cert`) ist nicht gesagt, ob eine laufende Verbindung beeinflusst wird. Hinweis für die Umsetzung, keine Aktion erwartet. | LH-FA-23 | spec/spezifikation.md · LH-FA-23.a | nein | Gültigkeit zur Laufzeit nicht gesagt |
| F-186 | INFO | Zeitpunkt der Fehlerantwort-Regeln und Zählung: Die Spezifikation legt `PGR-E5004`-Meldung auf die erste Abweichung je Interaktion fest; mit `--continue-on-error` kann damit eine Interaktion mehrere abweichende Antworten tragen, die nur einmal gemeldet werden. Das ist stimmig mit Exit-Code `5` am Ende, die Zahl der Abweichungen steht nicht im Ergebnis; Hinweis, nicht neu bewertet gegen LH-QA-05. | LH-QA-05 | spec/spezifikation.md · LH-FA-20.a Schritt 6 | nein | Hinweis ohne Aktion |

## Setzungen

Die früheren Setzungen S-1 bis S-27 und S-28 bis S-40 bleiben unbestätigt, soweit nicht unten anders vermerkt.

| Nr. | Stand / Setzung | Fundort | Überschreitet das Lastenheft? |
|---|---|---|---|
| S-31 | im Kern erledigt: Befehl ohne abschließende Zahlen; Lastenheft „der ausgeführte Befehl" wird getragen | LH-FA-24.a | nein mehr |
| S-35 | erledigt: unvollständige Aufzeichnung wird verglichen, das Lastenheft trägt es (Boundary) | LH-FA-24.a, LH-FA-24 | nein mehr; Rest siehe F-177 |
| S-36 | erledigt: Fehler wie aufgezeichnet gilt mit Vergleich als erwartet, das Lastenheft trägt es | LH-FA-20.a Schritt 6, LH-FA-24 | nein mehr |
| S-39 | bleibt: Gesamtfolge bei Extended; die Meldung nennt die Position in der normalisierten Folge | LH-FA-24.a | nein |
| S-41 | Mit Vergleich ist `--allow-recorded-errors` wirkungslos (ADR sagt „überflüssig"); die Spezifikation sagt es nicht | [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md), LH-FA-20.a Schritt 6 | nein — Lastenheft: „auf ausdrückliche Wahl"; Folge siehe F-175 |
| S-42 | Eine abgelaufene Zertifizierungsdatei (`--upstream-ca`) ist kein Startfehler, sondern `PGR-E4005` beim Verbinden; ein abgelaufenes eigenes Zertifikat ist `PGR-E2007` beim Start | LH-FA-20.a Schritt 2, LH-FA-23.a | teilweise — Lastenheft: „ungültiges Zertifikat" ohne Zeitpunkt; Folge siehe F-181 |
| S-43 | Nach dem Präfix einer unvollständigen Aufzeichnung sind weitere Nachrichten keine Abweichung, eine fremde `ErrorResponse` schon | LH-FA-24.a | teilweise — Lastenheft: „vergleicht, soweit sie reicht"; Folge siehe F-177 |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Hard Rule 3.4 — `spec/lastenheft.md` (LH-FA-24) | geprüft, ohne Befund; kein ADR-, Slice-, Wellen- oder Spec-Datei-Verweis; Aussagen sind fachlich (Wunsch, erwartet, Abweichung), keine Nachrichtentypen oder Feldnamen; Protokollnähe einzelner Wörter ist F-156 |
| Hard Rule 3.4 — `spec/spezifikation.md` (geänderte Abschnitte) | geprüft, ohne Befund; kein Verweis nach unten (ADR, Slice, Welle) |
| Hard Rule 3.4 — `spec/architecture.md` | geprüft, ohne Befund; im Diff unverändert; Play-Sequenz („bis ReadyForQuery", „Antworten vergleichen") trägt die Gesamtfolge; für eine Aufzeichnung ohne `ready_for_query` sagt die Sequenz nichts, kein Widerspruch |
| Hard Rule 3.5 — ADR-Dateien 0014 bis 0020 | geprüft, ohne Befund; im Diff unverändert, [ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md) ist neu und trägt Status Accepted mit Geschichtszeilen für Proposed und Accepted; Teilersetzung siehe F-178 |
| Hard Rule 3.3 — Move und Inhaltsänderung | geprüft, ohne Befund; Diff enthält keine Umbenennung |
| Hard Rule 3.7 — Kommentare und Zustandsfelder (Slices, Geschichte der neuen ADR, Handbuch-Kopf) | geprüft, ohne Befund; Indikativ, keine Chronik in Zustandsfeldern |
| Lastenheft ohne Wie — LH-FA-24 | geprüft, ohne neuen Befund |
| Handbuch ohne interne Kennungen | geprüft, ohne Befund; nur `PGR-…`, keine `LH-`, `SPEC-`, `ARC-`, `ADR-`, `slice-`, `welle-` |
| Handbuch ↔ Spezifikation (Vorrang, Exit-Codes `5` am Ende und `4` für Verbindungsfehler, Präfix-Vergleich, Zertifikatsablauf) | geprüft; stimmt, Abweichungen siehe F-175, F-180, F-181 |
| Exit-Code-Rangfolge Schritt 6, Schritt 7, Handbuch-Tabelle, Slices | geprüft; Übereinstimmung; „zuerst aufgetretene Ursache" ist überall entfallen; Signal-Exit-Code `0`/`4`/`5` ist konsistent (mit Vergleich ist `4` nur ein sofortiger Verbindungsfehler) |
| Verbindungsfehler (`PGR-E4002`, `PGR-E4003`, `PGR-E4005`) immer `4`, auch mit `--continue-on-error` und Vergleich | geprüft, ohne Widerspruch zu Schritt 5, Slice und Handbuch; eine vorherige Abweichung wird dabei vom Exit-Code `4` überdeckt, im Handbuch genannt |
| Zählungen — Abnahmeszenarien (Lastenheft §7: 16; Roadmap M3, Welle: „12 bis 16") | geprüft, ohne Befund; unverändert, Szenario 16 nennt den erwarteten Fehler nicht, ohne Widerspruch |
| Zählungen — Optionen (LH-FA-20.a Tabelle, Handbuch-Optionstabelle) | geprüft, ohne Befund; keine Option hinzugekommen oder entfallen |
| Zählungen — Codes (Codetabelle, Handbuch) | geprüft, ohne Befund; kein Code hinzugekommen oder entfallen; Klasse und Exit-Code von `PGR-E5004`, `PGR-E4004`, `PGR-E4005` stimmen, Beschreibung siehe F-180, F-181 |
| Zählungen — ADR-Index | geprüft, ohne Befund; Zeile 0021 vorhanden, Status Accepted, Bezug LH-FA-24, Reihenfolge lückenlos; Schärft-Anker (`LH-FA-24.a`, `LH-FA-20.a`, `SPEC-018`, `SPEC-027`, `SPEC-034`) zeigen auf vorhandene Abschnitte |
| docs/plan/planning — beide Slices (Liefer-Punkte, Form) | geprüft, Form ohne Befund; Inhalt siehe F-182, F-183 |
| Core-Reinheit, Gate-Lockerung (Hard Rule 3.2, 3.6), `.a-check.yml` | geprüft, ohne Befund; nicht im Diff |
| `make gates` | nicht ausgeführt (Auftrag) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 6 |
| LOW | 3 |
| INFO | 3 |

Von den elf nachgeprüften Findings (F-164 bis F-174) sind sechs behoben (F-164, F-165, F-166, F-168, F-170, F-171), fünf mit Rest behoben oder teilweise (F-167, F-169, F-172, F-173, F-174). Neu: F-175 bis F-186.

**Finding-Klassen dieses Laufs:** Optionswechselwirkung nicht geregelt · Ersetzte Regel nimmt Teilregel ungewollt mit · Ausnahmeregel und Allgemeinregel im selben Absatz widersprüchlich · Teilersetzung ohne maschinenlesbare Beziehung · ADR sagt mehr oder anderes als die Spezifikation · Codetabelle trägt Ausnahme des Codes nicht · Verweis auf Test, der nirgends steht.

Hinweis zum Zähler (Steering-Loop): „Korrektur zieht Nachbarstellen nicht nach" tritt weiter auf (F-180, F-181, F-182, F-183); das ist die fünfte Folge dieser Klasse über die Läufe. Der Folgeschritt (Prüfliste für Nachbarstellen bei Änderung einer Fehlerregel: Codetabelle, Handbuch-Tabelle, Slice-DoD, ADR-Bezüge) steht aus dem Vorgänger-Report offen; als Fitness Function gibt es weiterhin keinen Sensor.

## Verdikt

**Merge-blockierend:** ja, für die MEDIUM-Findings F-175, F-176, F-177, F-178, F-181 und F-182 (Wechselwirkung `--allow-recorded-errors` mit Vergleich, Fortsetzung bei Extended nach erwartetem Fehler, widersprüchliche Präfix-Regel, Teilersetzung von [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) ohne formale Beziehung, Zertifikatsablauf in der Codetabelle gegen Schritt 2, Test ohne Träger). Kein HIGH: keine Verletzung von Hard Rule 3.3 bis 3.7, die ADR-Dateien 0014 bis 0020 sind unverändert. Die Hauptlinie der Korrekturen (Vergleich ersetzt Fehlerregel, Exit-Code `5` am Ende, `4` nur für Verbindungsfehler, Gesamtfolge, Befehl ohne Zahlen, Präfix-Vergleich) ist über Lastenheft, Spezifikation, Handbuch und Slices widerspruchsfrei.

**Übergabe:** F-175, F-176, F-177 an den Implementer (Autor der Spec); F-178, F-179 an den Architect ([ADR-0021](../plan/adr/0021-vergleichsregeln-beim-einspielen.md): Form der Beziehung zu [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und [ADR-0020](../plan/adr/0020-antwortvergleich-beim-einspielen.md); eine Änderung an der Accepted-ADR wäre eine neue ADR); F-180, F-181 an den Implementer; F-182, F-183 an den Planner; Setzungen S-41 bis S-43 an den Auftraggeber beziehungsweise Architect.
