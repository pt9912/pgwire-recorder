# Review-Report: spec-erstfassung — 2026-10-03

**Review-Art:** Design (Spec-Erstfassung gegen Source Precedence, Decken-Regel,
Hard Rules und Spec-Template; Konsistenz der drei Straten untereinander)

**Gegenstand:** Working-Tree-Diff gegenüber `HEAD` (d027ecc) für
`spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`
(übernommen aus den externen Quellen LASTENHEFT, PFLICHTENHEFT, ARCHITECTURE).

**Skill:** `.harness/skills/reviewer.md` @ d027ecc — **Skill ist noch die
unausgefüllte Vorlage** (Platzhalter in der HIGH-Liste, siehe F-27); die
Klassifikation folgt daher dem Skill-Gerüst plus `v6.13.0` ·
`regelwerk/modul-03-spec.md` und `regelwerk/grundlagen-referenz-richtung.md`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- `AGENTS.md` (Hard Rules 3.4, 3.7; Source Precedence)
- `harness/conventions.md`
- `v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation, §Ziel-Form: Architektur-Sicht
- `v6.13.0` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)
- `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §ID-Schema als Klammer, §Spec-Stratifizierung
- Spec-Vorlagen unter `v6.13.0` · `templates/spec/`
- externe Quellen LASTENHEFT, PFLICHTENHEFT (nur für die Abgrenzung „übernommen vs. hinzugefügt")
- `make docs-check`: 20 Dateien, 0 Befunde (Gate grün; fängt keine der Findings unten)

Nicht Gegenstand: DoD-/Abnahme-Konformität (Verifier), fachliche Eignung (Validator).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die Architektur-Sicht ist an Go gebunden: Go-Code-Blöcke (Ports, Domain-Typen, Matcher, Composition Root), `.go`-Dateinamen im Package-Baum, `os.Exit`, `context.Context`, Goroutine-Aussage, „Go-Binary". | AGENTS.md Hard Rule 3.4 (sprach- und meilensteinfrei); `v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Architektur-Sicht | spec/architecture.md:49, 112, 124-125, 130-178, 186-207, 219-255, 266-293, 414-429, 515-527, 544 | nein — Doku-Gate prüft Sprachfreiheit nicht | Sicht trägt Sprach-Artefakte |
| F-2 | MEDIUM | Das Lastenheft (Decke) nennt `spec/spezifikation.md` und `spec/architecture.md` namentlich und verlagert ganze Abnahmekriterien nach unten („in der Spezifikation festgelegt"); §8 „Offene Punkte für die Spezifikation" ist eine Liste von Abwärts-Verweisen. Die Matrix kennt für Vertrag → Technik/Sicht nur ❌; das Gate meldet 0 Befunde, die Vorlage trägt denselben Satz in einer Messmethode. | `v6.13.0` · `regelwerk/grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP), Matrix-Zeile Vertrag | spec/lastenheft.md:27, 208, 482, 506, 522, 650 (Prosa); 218, 310, 336, 380, 402, 425, 490 (Kriterien); 647-668 (§8) | nein — `spec-straten`/`no-downward` fängt Klartext-Pfade nicht | Abwärts-Zeiger im Vertrag |
| F-3 | MEDIUM | Die Geltungsgrenze des Vertrags wird von der frei fortschreibbaren Spezifikation bestimmt: „definierter Funktionsumfang … wird in der Spezifikation konkretisiert", Out-of-Scope „Garantien für Parallelität über die Spezifikation hinaus", Glossar „Unterstützt = innerhalb des für v1 definierten Funktionsumfangs". Die Spezifikation darf den Umfang ohne Change Request verschieben. | `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §Spec-Stratifizierung | spec/lastenheft.md:206-209, 386, 600 | nein | Vertragsumfang hängt am Technik-Stratum |
| F-4 | MEDIUM | Exit-Code bei `SIGINT`/`SIGTERM` ist nirgends festgelegt, obwohl LH-FA-13 Boundary das ausdrücklich verlangt. Da ein Record-/Replay-Server nur per Signal endet, ist der Regelfall des Prozessendes ungeregelt (0? 130/143? 1?). | LH-FA-13 | spec/spezifikation.md:336-343, 500 (SPEC-013); spec/lastenheft.md:401-402 | ja — Spec-Lesung; später Integrationstest | Boundary-Kriterium ohne Spec-Festlegung |
| F-5 | MEDIUM | Verhalten bei bereits existierendem `--output` ist nicht festgelegt (Überschreiben oder Ablehnen); LH-FA-08 Boundary verlangt die Festlegung. Auch die Wechselwirkung mit „temporäre Datei, dann atomar verschieben" bleibt offen. | LH-FA-08 | spec/spezifikation.md:54-60, 237-261; spec/lastenheft.md:286-288 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-6 | MEDIUM | Unterstützte PostgreSQL-/PGWire-Versionen stehen auf „offen" (SPEC-029); LH-FA-05 nennt PGWire-Versionen ausdrücklich als Gegenstand der Konkretisierung, und LH-FA-05 Boundary verlangt einen eindeutigen Status am Rand des Umfangs. Protokollversion 3.0 wird nirgends genannt, das Verhalten bei einem anderen `ProtocolVersion` im Startup fehlt. | LH-FA-05 | spec/spezifikation.md:537; spec/lastenheft.md:208-209, 216-218, 665 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-7 | MEDIUM | Der Default des Strictness-Schalters ist „bei Implementierung festzulegen" (SPEC-012) — ein Wert-Eintrag ohne Wert in der Defaults-Tabelle. Es fehlen Optionsname, Env-Variable und der Exit-Code, falls der Schalter nicht verbrauchte Interaktionen als Fehler wertet. „strict" trägt zwei Bedeutungen (Matching, SPEC-011; Unverbrauchtes, SPEC-012); „standardmäßig strict" suggeriert einen nicht definierten Nicht-strict-Modus. | LH-FA-09, LH-FA-17 | spec/spezifikation.md:130-132, 267-268, 485, 576 | ja — Spec-Lesung | Default offen gelassen |
| F-8 | MEDIUM | Prozess- und Verbindungsebene sind vermischt: SPEC-017/018/019 „Prozess endet" bzw. „Verbindung wird beendet beziehungsweise Prozess endet" für Fehler, die je Client-Verbindung auftreten (Upstream-Verbindung entsteht je Session, Mismatch je Query). Offen: ob ein Mismatch den ganzen Replay-Server beendet, was der Client sieht, und wann das Recording als „abgeschlossen" geschrieben wird (Auslöser: Verbindungsende, Signal, erste Session?). Abnahmeszenario 1 setzt dafür einen Zeitpunkt voraus. | LH-FA-07, LH-FA-13, LH-QA-05 | spec/spezifikation.md:83-84, 172, 237-251, 305-307, 505-506; spec/lastenheft.md:146-147, 609-613 | ja — Spec-Lesung | Prozess- vs. Verbindungsfehler ungeklärt |
| F-9 | MEDIUM | Multi-Session ist für Replay nur negativ umrissen („erst dann … unterstützt, wenn Strategie implementiert"). Offen: zweite Verbindung gegen ein Ein-Session-Recording, Verbindung gegen ein Mehr-Session-Recording, Cursor-Zuordnung. LH-FA-12 Boundary verlangt die Festlegung in der Spezifikation. | LH-FA-12, LH-FA-03 | spec/spezifikation.md:119-121, 327-332; spec/lastenheft.md:379-381 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-10 | MEDIUM | LH-QA-03 verweist für den Plattformumfang auf die Spezifikation; die Spezifikation führt keinen (kein Betriebssystem, keine Architektur, keine Mindestumgebung). Messmethode ist damit nicht ausführbar. | LH-QA-03 | spec/lastenheft.md:520-522; spec/spezifikation.md (keine Stelle) | ja — Spec-Lesung | Verweisziel in Spec fehlt |
| F-11 | MEDIUM | Die Spezifikation erweitert: Fertigstellungskriterium 8 verlangt Container-Image und reproduzierbaren Build als Bedingung für „v1 technisch fertig"; LH-FA-16 ist SOLL, LH-FA-16.a selbst sagt „kann". „Reproduzierbar gebaut" steht in keiner Lastenheft-Aussage. | „Präzisieren, nie erweitern" (`v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation) | spec/spezifikation.md:660, 369; spec/lastenheft.md:454-459 | nein | Spec erweitert Vertrag |
| F-12 | MEDIUM | Log-Level werden als unterstützt benannt, aber weder Option noch Variable benennen sie; LH-FA-17 verlangt alle Einstellungen per CLI, LH-FA-14 Boundary „konfigurierte Detailstufe". Die CLI-Optionstabellen enthalten nur Pflichtoptionen von `record`/`replay`; „genaue CLI-Kommandos und Optionen" (Lastenheft §8) bleibt teilweise offen. | LH-FA-14, LH-FA-17 | spec/spezifikation.md:35-45, 54-60, 100-105, 347-352, 390-393; spec/lastenheft.md:424-425, 486-490, 652 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-13 | MEDIUM | Drei neu hinzugefügte Akzeptanzkriterien haben kein Spec-Gegenstück: LH-FA-02 Boundary (leere Aufzeichnung oder „klar benannter Zustand"), LH-FA-07 Boundary (Recording auf anderem Rechner „unverändert" nutzbar), LH-FA-11 Negative („nicht unterstützte Fehlerart") — letzteres steht im Widerspruch zu LH-FA-11.a/05.a, wonach jede `ErrorResponse` aufgezeichnet wird und die Nachrichtenmenge „nicht beschränkt" ist. | LH-FA-02, LH-FA-07, LH-FA-11 | spec/lastenheft.md:143-145, 265-267, 360-361; spec/spezifikation.md:153-167, 221-256, 311-315 | ja — Spec-Lesung | Neu ergänztes Kriterium ohne Spec-Deckung |
| F-14 | MEDIUM | Randverhalten am Protokollrand bleibt unbestimmt: TLS „darf" abgelehnt werden; CancelRequest, COPY, LISTEN/NOTIFY, GSSENCRequest sind „nicht zugesichert", aber nicht als abgelehnt/Exit 6/ignoriert definiert. LH-FA-05 Boundary verlangt eindeutigen Status, LH-FA-10/RB-02 klar erkennbare Ablehnung. | LH-FA-05 | spec/spezifikation.md:198-206, 625-642; spec/lastenheft.md:216-221, 571-572 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-15 | MEDIUM | `version: 1` ist eine einzelne Ganzzahl, SPEC-001 lehnt aber „unbekannte Major-Versionen" ab; ein Major/Minor-Modell ist nirgends definiert. Verhalten bei `version: 2` bzw. fehlender Formatkennung (LH-QA-06: „unbekannte Version wird erkannt") ist damit nicht eindeutig. | LH-QA-06 | spec/spezifikation.md:414-421, 483; spec/lastenheft.md:547-548 | ja — Spec-Lesung | Versionsmodell unscharf |
| F-16 | MEDIUM | Das Architektur-Gate wird in Gegenwartsform beschrieben („werden geprüft", „schlägt fehl", „Die Konfiguration liegt in der Repository-Wurzel"); im Repo existieren weder `.a-check.yml` noch ein Make-Target oder Gate-Fragment dazu. | AGENTS.md §4 (kein Gate nennen, das nicht existiert); Hard Rule 3.7 (Kommentar beschreibt, was da ist) | spec/architecture.md:572-596; spec/spezifikation.md:593-598 | ja — `make help`, Dateisuche | Beschriebenes Gate existiert nicht |
| F-17 | LOW | Zustandsautomat Replay kennt weder `unsupported → FAILED` noch den Abschluss mit unverbrauchten Interaktionen, obwohl LH-FA-05.a/03.b beides für Replay festlegen; der Record-Automat führt `unsupported`. | LH-FA-03.b, LH-FA-05.a | spec/architecture.md:451-462; spec/spezifikation.md:125-132, 169-174 | nein | Sicht deckt Spec-Pfad nicht ab |
| F-18 | LOW | Composition-Root-Beispiel ruft `service.NewStrictMatcher`/`service.NewReplayService`; Package heißt `services` (ARC-002, Package-Baum). | Maintainability | spec/architecture.md:517-521 vs 89, 144 | nein | Bezeichner-Drift innerhalb der Sicht |
| F-19 | LOW | Die Sicht behauptet, der Record-Service halte kein komplettes Resultset im Speicher, zeigt aber `Query(...) ([]model.Response, error)`, `Interaction.Responses []Response` und `Save(…, *model.Recording)` über das ganze Recording; ein Abbruch mit „als unvollständig erkennbarer Datei" (LH-FA-07.a) ist mit einem einzigen `Save` am Ende nicht darstellbar. | SPEC-032, LH-FA-07.a | spec/architecture.md:189-193, 237-262, 283, 347; spec/spezifikation.md:243-251, 548-550 | nein | Streaming-Zusage ohne tragende Signatur |
| F-20 | LOW | Spezifikation §8-§11 (Leitentscheidungen, Tests, „Nicht zugesichert", Fertigstellung) tragen keine `SPEC-*`; nur SPEC-030/-011 werden darin zitiert. Eine ADR hat dort kein benennbares `Schärft:`-Ziel. | `v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation (Zwei Kennungs-Arten) | spec/spezifikation.md:566-660 | nein | Spec-Inhalt ohne Adresse |
| F-21 | LOW | Das Lastenheft verweist für die Stelle des Version-Bumps auf `harness/conventions.md` („deklariert dein Repo"); dort ist sie nicht deklariert. Vorlagen-Anrede („dein Repo") steht im bindenden Text. | `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §Spec-Stratifizierung | spec/lastenheft.md:3-6; harness/conventions.md:150-173 | nein | Vorlagen-Prosa im Vertragstext |
| F-22 | LOW | Messmethoden teils nicht messbar: „jedes Mal gleich" ohne Wiederholungszahl (QA-01), „typische" Umgebungen (QA-03). LH-QA-05 bindet die Unterscheidbarkeit an Prozessstatus, LH-FA-13 Negative an „Fehlerklassen aus LH-QA-05" — passt (5 Klassen auf Codes 2/3/4/5/6), aber QA-01 und QA-05 sind keine der in der Vorlage genannten Kategorien. | Vorlage `lastenheft.template.md` §4 (messbare Anforderung) | spec/lastenheft.md:504-506, 520-522, 538-539 | nein | Messmethode unscharf |
| F-23 | LOW | LH-RB-02 „dürfen mit einem klar erkennbaren Fehler abgelehnt werden" (Erlaubnis) steht neben LH-FA-05 Negative „wird … abgelehnt" (Pflicht). | Maintainability | spec/lastenheft.md:571-572 vs 219-221 | nein | Modalverb-Widerspruch |
| F-24 | INFO | **Von der Übernahme hinzugefügte Inhalte** (die Quellen kennen weder Given/When/Then, Boundary, Negative, Out-of-Scope je Anforderung noch Messmethoden). Einordnung siehe Tabelle unter „Übernahme-Zusätze". Überschreiten das Quell-Lastenheft: LH-FA-02 Boundary, LH-FA-06 Boundary, LH-FA-07 Boundary, LH-FA-08 Boundary, LH-FA-11 Negative, LH-FA-12 Boundary, LH-FA-13 Boundary/Negative (Exit-Bindung), LH-QA-05/-03 Messmethode. Das Lastenheft ist `Draft`; Freigabe durch den Auftraggeber steht aus. | Maintainability; `v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Akzeptanzkriterium | spec/lastenheft.md:107-496, 500-548 | nein | Übernahme-Zusatz ohne Quellenbeleg |
| F-25 | INFO | Die `LH-*`-Kennungen sind im Doku-Gate nicht aktiviert (Muster auskommentiert, `ids` kennt nur `ADR-*`); Querverweise `LH-*`/`SPEC-*`/`ARC-*` werden von keinem Sensor aufgelöst. Manuelle Prüfung dieses Laufs: alle Verweise lösen auf; Zählung fortlaufend (LH-FA-01–17, LH-QA-01–06, LH-RB-01–02, Verfeinerungen `.a`–`.d` nur zu existierenden FA, SPEC-001–033 in Dateireihenfolge, ARC-001–013 mit §3 als Fortsetzung). | `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §Vergabe | .d-check.yml:13-15 | ja — `make docs-check` nach Aktivierung des Musters | Referenz-Auflösung ungeprüft |
| F-26 | INFO | Spezifikation legt Go und `pgproto3` fest (Technik-Stratum zulässig, aber Entscheidungs-Charakter wie eine ADR); §10 „Nicht zugesichert in v1" benennt COPY/CancelRequest u. a., die das Lastenheft nicht ausschließt (LH-FA-05 delegiert) — Zuständigkeit Architect, siehe F-3. | Maintainability | spec/spezifikation.md:451, 570-572, 625-642 | nein | Technik-Entscheidung im Spec-Stratum |
| F-27 | INFO | `.harness/skills/reviewer.md` und `harness/conventions.md` (Baseline-Block, `MR-000`-Datum, Vertragspräfix `LH`) sind unausgefüllte Vorlagen; nicht Teil des Diffs. Zuständig: Architect/Implementer; Verifier-Frage, nicht Reviewer-Frage. | Maintainability | .harness/skills/reviewer.md:1-47; harness/conventions.md:29-39, 155 | nein | Harness-Datei unausgefüllt |

## Übernahme-Zusätze (Prüfpunkt 5)

Alle Boundary-/Negative-/Out-of-Scope-Zeilen je LH-FA und alle Messmethoden sind
Zusätze (Vorlagenpflicht), keine Quellen-Aussagen. „Überschreitet" heißt: das
Kriterium begründet eine neue vertragliche Pflicht, die in LASTENHEFT/PFLICHTENHEFT
nicht steht.

| Zusatz | Stelle | Überschreitet Quell-Lastenheft? |
|---|---|---|
| FA-01 Boundary/Negative (Start ohne Betriebsart, ungültige Parameter) | lastenheft.md:119-123 | nein — folgt aus Exit-Code-Semantik |
| FA-02 Boundary (leere Aufzeichnung gültig oder benannter Zustand; „bleibt nicht hängen") | lastenheft.md:143-145 | ja — neue Pflicht, offene Alternative im Vertrag (F-13) |
| FA-03 Boundary/Negative | lastenheft.md:167-171 | nein |
| FA-04 Boundary (Nichtunterstützung melden) | lastenheft.md:192-194 | nein |
| FA-06 Boundary (mehrere Ergebnismengen, ohne Datensätze) | lastenheft.md:241-243 | ja — leicht, präzisiert „genügend Informationen" |
| FA-07 Boundary (kopierte Aufzeichnung, anderer Rechner) | lastenheft.md:265-267 | ja — Portabilitätspflicht, kein Spec-Gegenstück (F-13) |
| FA-08 Boundary (existierendes Ziel) | lastenheft.md:286-288 | ja — erzeugt Spec-Pflicht (F-5) |
| FA-09 Boundary (identische Queries → eindeutige Zuordnung) | lastenheft.md:308-310 | nein — Quellen-Offene-Punkte nennen das Thema |
| FA-10 Boundary (ähnliche Anfrage) | lastenheft.md:333-336 | nein — deckt sich mit „kein Fuzzy" der Spec |
| FA-11 Negative (nicht unterstützte Fehlerart) | lastenheft.md:360-361 | ja — widerspricht Spec (F-13) |
| FA-12 Boundary (mehrere Verbindungen) | lastenheft.md:379-381 | nein — Quellen-Offene-Punkte nennen Sessions |
| FA-13 Boundary (Signal) / Negative (Klassen unterscheiden) | lastenheft.md:401-405 | ja — bindet Exit-Code an Signal und an die QA-05-Klassen |
| FA-14/15/16/17 Boundary/Negative | lastenheft.md:424-427, 445-448, 466-469, 489-492 | nein |
| QA-01..06 Messmethoden | lastenheft.md:504-547 | ja für QA-03 (Lauf lokal und CI, Plattform in Spec), QA-05 (Szenario je Klasse über Meldung und Prozessstatus), QA-06 (unbekannte Version erkannt) |
| SPEC-020..028 Fehlerklasse → Exit-Code-Zuordnung | spezifikation.md:511-521 | nein — in PFLICHTENHEFT §12/§13 getrennt vorhanden; die Zuordnung je Klasse ist Zusatz, aber konsistent mit der Code-Tabelle |
| Doppelte-Query-Zuordnung (Position am Cursor, höchstens einmal verbraucht) | spezifikation.md:289-290 | nein — Spec-Präzisierung von FA-09 Boundary; konsistent mit architecture.md:402-408 |
| Abnahmeszenario 1-6 (Klartext, „Bezug:"-Zeilen) | lastenheft.md:607-645 | nein — „Bezug:" ist neu, präzisiert nur |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Referenz-Richtung architecture.md (ADR/Slice/Welle/Commit-Bezüge im Körper) | geprüft, ohne Befund (Treffer nur in der Hard-Rule-Kopfzeile und in Vorlagen-Hinweisen) |
| Chronik/Zustandsfelder architecture.md (`Letzte Änderung`, `Status: Aktiv`) | geprüft, ohne Befund |
| Historie lastenheft.md §9 | geprüft, ohne Befund (Verweis `—`, kein ADR/Slice/Spec-Bezug) |
| Historie spezifikation.md §12 | geprüft, ohne Befund |
| ID-Zählung je Datei (LH, SPEC, ARC) | geprüft, ohne Befund (F-25) |
| Querverweise LH-FA-xx.y / SPEC-0xx / ARC-0xx | geprüft, ohne Befund (alle auflösbar) |
| Exit-Code-Tabelle ↔ Fehlerklassen (SPEC-013..028) ↔ LH-FA-13/LH-QA-05 | geprüft, ohne Widerspruch (Lücken siehe F-4, F-8) |
| Matching-Aussagen (LH-FA-09.a, LH-FA-10.a, architecture §4.1) | geprüft, ohne Widerspruch |
| Package-Pfade und Import-Regeln (spezifikation §8 ↔ architecture §1/§2/§6, Mermaid ↔ Tabelle) | geprüft, ohne Befund |
| Vorlagen-Konformität Happy/Boundary/Negative/Out-of-Scope je LH-FA (17 von 17) | geprüft, ohne Befund |
| Messmethode je LH-QA (6 von 6), Art/Vorgabe/Nachweis je LH-RB (2 von 2) | geprüft, vorhanden (Qualität: F-22) |
| `make docs-check` | 0 Befunde |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 15 |
| LOW | 7 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Sicht trägt Sprach-Artefakte · Abwärts-Zeiger im Vertrag · Vertragsumfang hängt am Technik-Stratum · Boundary-Kriterium ohne Spec-Festlegung (F-4, F-5, F-6, F-9, F-12, F-14: sechs Vorkommen) · Default offen gelassen · Prozess- vs. Verbindungsfehler ungeklärt · Neu ergänztes Kriterium ohne Spec-Deckung · Spec erweitert Vertrag · Beschriebenes Gate existiert nicht

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH, Hard Rule 3.4). Die MEDIUM-Findings F-4 bis
F-9 und F-12/F-14 sind die bekannten Lücken bzw. nicht erfüllte
„in der Spezifikation festgelegt"-Kriterien des Lastenhefts; F-2/F-3/F-11/F-13
berühren den Vertragsrang und gehören vor `Accepted` an den Architect, weil das
Lastenheft im Status `Draft` frei änderbar, danach nur per Change Request
änderbar ist. Die Klasse „Boundary-Kriterium ohne Spec-Festlegung" tritt sechsmal
auf und erreicht damit die Schwelle des Steering-Loop-Zählers.

**Übergabe:** Findings gehen an den Implementer; F-2, F-3, F-11, F-13, F-26 an
den Architect (Vertrags-/Stratum-Frage). Die Finding-Klassen gehen in die
Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
