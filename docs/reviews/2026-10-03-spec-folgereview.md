# Review-Report: spec-folgereview — 2026-10-03

**Review-Art:** Design (Folgelauf zu `spec-erstfassung`: Abarbeitung der Findings F-1 bis F-27 gegen
Source Precedence, Decken-Regel, Hard Rules; Konsistenz der drei Straten untereinander)

**Gegenstand:** Working-Tree-Diff gegenüber `HEAD` (2044483) für `spec/lastenheft.md`,
`spec/spezifikation.md`, `spec/architecture.md`, `docs/plan/adr/0009-implementierungssprache-go.md`,
`harness/conventions.md`, `harness/README.md`, `.harness/skills/reviewer.md`.

**Skill:** `.harness/skills/reviewer.md` @ Working Tree (ausgefüllt; Repo-HIGHs: Spec-Stratum nennt
Artefakt, ADR/Move-Regel, Core-Reinheit).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- `docs/reviews/2026-10-03-spec-erstfassung.md` (Findings F-1 bis F-27)
- `AGENTS.md` (Hard Rules 3.3 bis 3.7, §4, Source Precedence)
- `.a-check.yml` (Schichten, Kanten, `tech`-Muster)
- berührte `LH-*`: LH-FA-02, -03, -05, -07, -08, -09, -10, -11, -12, -13, -14, -17, LH-QA-01, -03, -05, -06, LH-RB-02
- [ADR-0009](../plan/adr/0009-implementierungssprache-go.md) (Bezug des Diffs), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) (Verweisziele)
- `make help` (vorhandene Targets), `make docs-check`: 37 Dateien, 0 Befunde (fängt keines der Findings unten)

Nicht Gegenstand: DoD-/Abnahme-Konformität (Verifier), fachliche Eignung (Validator).

---

## Status der Findings aus dem Erstreport

Die Tabelle ist die Prüfung je Erst-Finding; Stand ist der Working Tree.

| Erst-ID | Status | Beleg |
|---|---|---|
| F-1 | behoben (Rest: N-2) | spec/architecture.md enthält keine Code-Blöcke, `.go`, `os.Exit`, `context.Context`, Goroutine-Aussagen mehr (Suche ohne Treffer); Bibliotheksnamen verbleiben, siehe N-2 |
| F-2 | behoben (Rest: N-9) | spec/lastenheft.md:26-27, 653-657: nur noch „nachgelagerte Dokumente"; Kriterien tragen den Inhalt selbst (Suche nach „Spezifikation" nur in der Regel-Prosa, :698) |
| F-3 | behoben | Funktionsumfang steht in LH-FA-05 (:206-212), Out-of-Scope LH-FA-12 (:394) und Glossar „Unterstützt" (:603) im Vertrag |
| F-4 | behoben | LH-FA-13 Boundary (:409-411); Spec LH-FA-13.b (:401-420), SPEC-013; Rest N-6 |
| F-5 | behoben (Rest: N-11) | LH-FA-08 Boundary (:292-294); Spec :281-284, `PGR-E2002`, `--force`; temporäre Datei im Zielverzeichnis (:294-296) |
| F-6 | behoben | LH-FA-05.e (spec :234-240), `PGR-E6002`, SPEC-029 |
| F-7 | behoben | Option `--fail-on-unconsumed`, Env-Variable, Default `false` (SPEC-012), `PGR-E5002`/Exit 5; „strict" nur noch für Matching |
| F-8 | teilweise | Fehlerebenen und Schreibzeitpunkt festgelegt (spec :281-292, :401-420); offen: Begriff „unerwartetes Verbindungsende" (N-4), laufende Sessions beim Herunterfahren (N-5), Architektur-Sequenz (N-1) |
| F-9 | behoben (Rest: N-3) | spec :370-385, `PGR-E6003`, Cursor je Verbindung |
| F-10 | behoben (Rest: N-8) | SPEC-035 (spec :715-718); LH-QA-03 trägt keinen Zeiger mehr |
| F-11 | teilweise | Container-Teil jetzt an SOLL gebunden (SPEC-040 Nr. 8, spec :824-826); „Binary reproduzierbar gebaut" steht weiter in keiner Lastenheft-Aussage (Suche: nur :248, :302-307, :509 anderer Bedeutung) |
| F-12 | behoben | `--log-level`/`PGWIRE_RECORDER_LOG_LEVEL` (spec :428-429), Optionstabelle (spec :476-490) |
| F-13 | teilweise | FA-02-Boundary: spec :286-292 deckt; FA-07-Boundary: SPEC-004 (spec :560-561) deckt, aber siehe N-12; FA-11-Negative neu gefasst (lastenheft :366-368), Spec-Gegenstück für nicht verlustfrei darstellbare Serverantwort fehlt (N-7) |
| F-14 | behoben (Rest: N-10) | Protokollrand-Tabelle (spec :242-253); TLS ist Pflicht-Ablehnung (spec :215-219) |
| F-15 | behoben | `version` ganzzahlig, Ablehnung unbekannter Version `PGR-E3002`, falsche Formatkennung `PGR-E3003` (spec :512-516) |
| F-16 | teilweise | spec/architecture.md:478-480 in Vorgabe-Form umgestellt; spec/spezifikation.md:751-753 beschreibt das Gate weiter in Gegenwartsform („werden … geprüft", „meldet"); kein Make-Target dafür (`make help`) — `.a-check.yml` existiert (tracked) |
| F-17 | teilweise | `unsupported → FAILED` im Replay-Automaten ergänzt (architecture :376); Abschluss mit unverbrauchten Interaktionen und Verbindungsende ohne `Terminate` in keinem Automaten |
| F-18 | behoben | `service.New…`-Beispiel nicht mehr vorhanden |
| F-19 | offen | siehe N-1: Streaming-Zusage (architecture :199-204) gegen Ganz-Recording-`Save` und Session-Verwurf |
| F-20 | behoben | SPEC-035 bis SPEC-040 vergeben, [ADR-0009](../plan/adr/0009-implementierungssprache-go.md) `Schärft:` zeigt auf SPEC-036; Zählung SPEC-001 bis SPEC-040 lückenlos, je einmal definiert |
| F-21 | behoben (Rest: N-9, N-13) | lastenheft :3-6 ohne Vorlagen-Anrede; Stelle in harness/conventions.md deklariert |
| F-22 | teilweise | QA-01: zehn Läufe (lastenheft :510-512); QA-03 weiter „typisch" (:523-528) und verweist auf ein „Abnahmeprotokoll", das nirgends definiert ist |
| F-23 | behoben | lastenheft :576-578 „werden … abgelehnt" |
| F-24 | offen (INFO bleibt) | Lastenheft `Draft`; Freigabe des Auftraggebers steht aus, Zusätze wachsen (siehe „Entscheidungen zur Bestätigung") |
| F-25 | offen (INFO bleibt) | `.d-check.yml` unverändert: `LH-*`-Muster weiter auskommentiert; manuelle Prüfung dieses Laufs: alle neuen Verweise (`LH-FA-13.b`, `LH-FA-05.e`, `SPEC-035`…`SPEC-040`, [ADR-0009](../plan/adr/0009-implementierungssprache-go.md)-Anker) lösen auf |
| F-26 | teilweise (INFO bleibt) | COPY/CancelRequest sind jetzt in LH-FA-05.e eindeutig abgelehnt, SPEC-039 führt sie weiter als „nicht zugesichert" (N-10) |
| F-27 | behoben | Platzhalter in reviewer.md, harness/README.md, harness/conventions.md ersetzt (nur Formbeispiele wie `<tag>`, `MR-<NNN>` bleiben); alle genannten Targets existieren (`docs-check`, `baseline-verify`, `gates`, `slice-mv`, `archive-welle`, `hooks-install`); Rest N-13 |

## Findings

Neue Findings aus diesem Lauf (Output-Schema des Reviewer-Skills; Nummern setzen die des Erstreports fort).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-28 | MEDIUM | Die Architektur zeigt `Save(Recording)` im Record-Sequenzdiagramm direkt nach der ersten Query und Schritt 7 nach „ReadyForQuery"; die neue Festlegung ruft `RecordingRepository` erst nach dem Ende jeder Session und beim Herunterfahren auf. Die Zusage „kein komplettes Resultset im Speicher" widerspricht dem Ganz-Recording-Speicherstand: Eine Interaktion wird mit `ReadyForQuery` Teil des Recordings, und eine Session mit nicht unterstützter Interaktion wird verworfen, was die Session bis zum Ende im Speicher voraussetzt (SPEC-032 „keine … kompletter Datenbanksitzungen im RAM"). | LH-FA-07; SPEC-032; Maintainability | spec/architecture.md:196, 199-204, 270, 284; spec/spezifikation.md:286-292, 185-186, 696-698 | nein | Streaming-Zusage ohne tragende Struktur |
| F-29 | HIGH | Die Sicht nennt Bibliotheks- und Sprachnamen: `pgproto3` (u. a. :38, 109-122, 249, 402-415, 496, 524, 526), `YAML-Bibliothek` als eigene Komponente (ARC-013) sowie die qualifizierten Typnamen `model.Query`/`model.Response`/`model.Recording`. Der Skill rechnet Bibliotheksnamen zu den Sprach-Artefakten der Sicht; AGENTS.md §3.4 erlaubt nur Pfade zu Code-Modulen. Ob ein Format-/Protokollname unter das Verbot fällt, ist eine Architect-Frage; nach Skill-Wortlaut Treffer. | AGENTS.md §3.4; Skill HIGH „Spec-Stratum nennt Artefakt" | spec/architecture.md:38, 109-122, 249-250, 411-415, 496-497, 524-526 | nein — kein Gate prüft Sprachfreiheit der Sicht | Sicht trägt Sprach-Artefakte |
| F-30 | MEDIUM | Ein Mehr-Session-Recording wird im Replay abgelehnt, im Record entsteht aber jede Verbindung als eigene Session, auch ohne Anfrage (Probe-Verbindung, Pool, abgebrochene Verbindung). Eine Aufzeichnung mit einer leeren Zweit-Session ist damit nicht replay-fähig; Verhalten bei einem Recording mit null Sessions (gültig laut :288-290, „genau einer Session" laut :378) ist ungeregelt. Berührt LH-FA-07 Happy Path und Abnahmeszenario 1/2. | LH-FA-07, LH-FA-12, Abnahmeszenario 1/2 | spec/spezifikation.md:127-128, 286-290, 378-385; spec/lastenheft.md:143-145, 394, 610-615 | ja — Spec-Lesung; später Integrationstest | Recording-Erzeugung und Replay-Annahme nicht deckungsgleich |
| F-31 | MEDIUM | „Unerwartetes Verbindungsende" (`PGR-E4003`, SPEC-028) ist ein Verbindungsfehler, aber nirgends definiert (Ende ohne `Terminate`? mitten in einer Interaktion? im Replay vor Verbrauch, dort auch `PGR-W2001`?). Folgen offen: Exit-Code beim Herunterfahren und ob die Session ins Recording kommt; berührt LH-FA-02 Boundary („Client ohne Anfrage, Verbindung endet → gültige Aufzeichnung"). | LH-FA-02, LH-FA-13, LH-QA-05 | spec/spezifikation.md:409, 618, 653, 135-138; spec/lastenheft.md:143-145 | ja — Spec-Lesung | Fehlerbedingung ohne Definition |
| F-32 | MEDIUM | Behandlung von Sessions, die beim Signal noch laufen, ist nicht festgelegt („enthält alle bis dahin beendeten Sessions" gegen „laufende Schreiboperationen soweit möglich abschließen"); „Recording als unvollständig behandeln" (LH-FA-13.a) hat keinen Gegenstand mehr, weil die Zieldatei nach LH-FA-07.a nie unvollständig ist. | LH-FA-13, LH-FA-07 | spec/spezifikation.md:286-288, 299-303, 394-399 | ja — Spec-Lesung | Zustand beim Herunterfahren ungeklärt |
| F-33 | MEDIUM | Verbindungsfehler beenden nur die Verbindung, der Prozess endet erst bei Signal; LH-FA-13 Negative lautet „when er auftritt, then signalisiert der Prozessstatus Misserfolg" und LH-QA-05 misst „über Meldung und Prozessstatus". Der Prozessstatus liegt damit erst nach einem späteren Signal vor (bei `SIGKILL` nie). Außerdem hängt der Exit-Code bei parallelen Verbindungen von „dem ersten" Fehler ab; ein harmloser `CancelRequest` (Klasse 6) bestimmt den Exit-Code des ganzen Laufs. | LH-FA-13, LH-QA-05, LH-QA-01 | spec/spezifikation.md:401-420, 248; spec/lastenheft.md:412-414, 540-543 | ja — Spec-Lesung | Prozess- vs. Verbindungsfehler ungeklärt |
| F-34 | MEDIUM | Lastenheft-Wortlaut „die n-te gleiche Anfrage erhält die n-te aufgezeichnete Antwort" lässt eine SQL-bezogene Zählung zu; die Spec ist positionsbasiert (strict sequential, SPEC-011). Bei Aufzeichnung A,B,A und Eingabe A,A gilt nach Lastenheft-Lesart ein Treffer, nach Spec ein Mismatch. Das Lastenheft gewinnt im Konfliktfall. | LH-FA-09, SPEC-011 | spec/lastenheft.md:314-316; spec/spezifikation.md:316-335 | ja — Spec-Lesung | Vertragswortlaut weiter als Spec-Verfahren |
| F-35 | MEDIUM | LH-FA-11 Negative verlangt „nicht unterstützte Protokollinteraktion" für eine nicht verlustfrei aufzeichenbare Fehlerantwort; Spec LH-FA-11.a und LH-FA-05.a (:165-179) legen für nicht darstellbare Serverantworten keine Behandlung fest (Tabelle in LH-FA-05.e kennt nur Client-Interaktionen und `NotificationResponse`). | LH-FA-11, LH-FA-05 | spec/lastenheft.md:366-368; spec/spezifikation.md:165-179, 362-368 | ja — Spec-Lesung | Neu ergänztes Kriterium ohne Spec-Deckung |
| F-36 | MEDIUM | SPEC-035 beschränkt die Zielplattformen auf Linux `amd64`/`arm64`; LH-QA-03 verlangt Bereitstellbarkeit für „typische lokale Entwicklungs- und CI-Umgebungen" (macOS-/Windows-Entwicklerrechner sind typisch). Die Spec kann den Vertragsumfang dadurch verengen. | LH-QA-03 | spec/spezifikation.md:715-718; spec/lastenheft.md:523-528 | nein | Spec verengt Vertrag |
| F-37 | LOW | Spec-Gate in Gegenwartsform trotz fehlendem Make-Target; architecture.md:478-490 endet mit „Eine Verletzung … ist ein Befund des Gates". | AGENTS.md §4; Hard Rule 3.7 | spec/spezifikation.md:751-753; spec/architecture.md:492-494 | ja — `make help` | Beschriebenes Gate existiert nicht |
| F-38 | LOW | Wortlaut-Drift: LH-FA-05.e legt CancelRequest, COPY, NotificationResponse als „nicht unterstützt, Exit 6" fest; SPEC-039 führt dieselben Punkte als „nicht zugesichert". | Maintainability | spec/spezifikation.md:248-252, 789-806 | nein | Modalverb-Widerspruch |
| F-39 | LOW | SPEC-014 („ungültige CLI-Verwendung oder Konfiguration") und LH-FA-07.a „Fehlermodi" nennen den Startfehler „`--output` existiert" (`PGR-E2002`, Exit 2) nicht; `PGR-E2002`, `PGR-E5002`, `PGR-E6002`, `PGR-E6003` tragen im Gegensatz zu den anderen Codes keine `SPEC-*`-Spalte. | Maintainability | spec/spezifikation.md:598, 305-306, 645, 656, 659-660 | nein | Zuordnungstabelle unvollständig |
| F-40 | LOW | Exit-Code-Vorrang bei gleichzeitigem Schreibfehler und gemerktem Verbindungsfehler ist nicht festgelegt (`3` oder Klasse des Fehlers). | LH-FA-13 | spec/spezifikation.md:416-419 | ja — Spec-Lesung | Boundary-Kriterium ohne Spec-Festlegung |
| F-41 | LOW | Lastenheft §9 zitiert in der Regel-Prosa `spezifikation.md`, `architecture.md`, ADR, Slice, Welle (Negation der Decken-Regel, kein Zeiger); Skill-HIGH „nennt … die darunterliegenden Spec-Dateien" trifft dem Wortlaut nach. Zusätzlich nennt :3-6 `harness/conventions.md`. §8 listet Punkte als „nicht in diesem Dokument entschieden", die das Lastenheft jetzt selbst entscheidet (Sessions, Signalbehandlung, TLS, Funktionsumfang). | AGENTS.md §3.4; Skill HIGH „Spec-Stratum nennt Artefakt" | spec/lastenheft.md:3-6, 653-675, 698 | nein | Vorlagen-Prosa im Vertragstext |
| F-42 | LOW | SPEC-004 sagt „keine rechnerspezifischen Angaben (Hostnamen, Adressen)", das Recording trägt aber Startup-Parameter und `ParameterStatus` des Servers, die Host-/Umgebungswerte enthalten können; die Zusage ist nicht prüfbar formuliert. | LH-FA-07 | spec/spezifikation.md:560-561, 524-535 | nein | Zusage ohne Prüfbarkeit |
| F-43 | LOW | Neue Regel „Major bei Änderung/Streichung, Minor bei neuer Anforderung, Patch bei Klarstellung" in harness/conventions.md hat keine Quelle im Repo oder in der Baseline-Stelle (die Baseline verlangt nur die Deklaration); es ist eine Setzung des Implementers. | Maintainability | harness/conventions.md:186-191 | nein | Setzung ohne Beleg |
| F-44 | INFO | Das Lastenheft nennt jetzt PGWire-Nachrichtennamen (`Query`, `Terminate`, `ReadyForQuery`) und TLS-Details im Vertragstext; passt zum Produkt (PGWire-Recorder), aber an der Grenze zu „technische Entscheidungen gehören nicht in dieses Dokument". Entscheidung Architect/Auftraggeber. | Maintainability | spec/lastenheft.md:206-212 | nein | Technikbegriff im Vertrag |
| F-45 | INFO | Die Architektur bleibt ohne ADR-, Slice-, Welle-, Commit-Bezug im Körper (Treffer nur in der Hard-Rule-Kopfzeile und in „Regeln dieser Sektion"-Prosa); Kanten und Pfade stimmen mit `.a-check.yml` überein (services nur nach model und ports-driven; `pgproto3` nur in `adapters/driving/pgwire` und `adapters/driven/postgres`, YAML nur in `adapters/driven/recording`; Composition Root `cmd/**`, `internal/bootstrap/**`). Die Kette „Rand merkt sich die Klasse des ersten Fehlers" (architecture :442) benennt den Meldeweg vom PGWire-Adapter zum CLI-Rand nicht. | Maintainability | spec/architecture.md:109-125, 442, 492-497 | ja — nach Aktivierung des Gates | Sicht deckt Meldeweg nicht ab |

## Entscheidungen zur Bestätigung

Entscheidungen des Implementers, die der Auftraggeber bestätigen muss. „Rang" nennt, wo die Festlegung
steht; „Überschreitung" beantwortet, ob sie über den Wortlaut des Lastenhefts hinausgeht oder ihn berührt.

| # | Entscheidung | Rang | Überschreitung des Lastenhefts |
|---|---|---|---|
| E-1 | Exit 0 bei Signal ohne vorherigen Fehler | Lastenheft (LH-FA-13 Boundary :409-411) und SPEC-013 | Nein, sie steht im Lastenheft selbst; sie ist neu gegenüber der Quelle (F-24) und braucht deshalb die Vertragsbestätigung. Siehe F-33 zur Wirkung auf LH-FA-13 Negative |
| E-2 | `--output` existiert → Ablehnung (Exit 2, `PGR-E2002`) außer `--force` | Prinzip im Lastenheft (LH-FA-08 Boundary :292-294); Exit 2 und `--force` nur Spec | Nein; das Lastenheft verlangt „nicht stillschweigend überschreiben, ausdrücklich verlangbar" und deckt die Spec. Die Prinzip-Aussage ist neu gegenüber der Quelle |
| E-3 | Mehr-Session-Recording im Replay abgelehnt (Exit 6, `PGR-E6003`) | nur Spec (:127-128, :378-385); Lastenheft nur Out-of-Scope (:394) | Wortlaut nein (Out-of-Scope verbietet keine Ablehnung); berührt aber „Aufzeichnung … für Replay verwendbar" (LH-FA-07 Happy, Abnahmeszenario 1/2), siehe F-30 |
| E-4 | `--fail-on-unconsumed` Default aus (Warnung `PGR-W2001`) | nur Spec (SPEC-012, :135-144) | Nein; das Lastenheft führt unverbrauchte Interaktionen nicht |
| E-5 | PGWire 3.0; andere Version → `PGR-E6002` (Exit 6) | nur Spec (:236-240, SPEC-029) | Nein; das Lastenheft nennt keine Protokollversion |
| E-6 | `CancelRequest`, `COPY`, `NotificationResponse` nicht unterstützt (Exit 6) | nur Spec (:242-253); Lastenheft LH-FA-05 Boundary/Negative und LH-RB-02 verlangen eindeutigen Status und Ablehnung | Nein; deckt sich mit LH-RB-02 (jetzt Pflicht), Wirkung auf den Exit-Code siehe F-33 |
| E-7 | Zielplattformen Linux `amd64`/`arm64` | nur Spec (SPEC-035) | Möglich: LH-QA-03 „typische lokale Entwicklungs- und CI-Umgebungen" wird verengt, siehe F-36 |
| E-8 | Recording-Schreibzeitpunkt nach jeder Session und beim Herunterfahren | nur Spec (:286-292); Architektur-Port :196 | Nein; das Lastenheft legt keinen Zeitpunkt fest. Widerspruch zur Architektur-Sequenz siehe F-28 |

Weitere Setzungen, die nicht in der Liste stehen und Auftraggeber-Rang berühren:

| # | Setzung | Rang | Überschreitung |
|---|---|---|---|
| E-9 | Verbindungsfehler beenden nur die Verbindung; Prozessstatus erst beim Herunterfahren (LH-FA-13.b) | nur Spec | Berührt LH-FA-13 Negative und LH-QA-05, siehe F-33 |
| E-10 | LH-RB-02 von „dürfen abgelehnt werden" auf „werden abgelehnt" (Erlaubnis → Pflicht) | Lastenheft (:576-578) | Ja, verschärft den Vertrag gegenüber der Quelle |
| E-11 | Neue Lastenheft-Festlegungen: QA-01 „mindestens zehn Läufe" (:511), FA-02 „gültige Aufzeichnung ohne Interaktionen" (:143-145), FA-05 Funktionsumfang mit TLS-Ausschluss (:206-212), FA-07 „keine rechnerspezifischen Angaben" (:270-273), FA-09 „n-te Anfrage" (:314-316, siehe F-34), FA-11 Negative (:366-368), FA-12 Out-of-Scope (:394), QA-03 „Abnahmeprotokoll" (:527-528) | Lastenheft | Ja: jede ist eine neue vertragliche Aussage gegenüber der Quelle; das Lastenheft ist `Draft` |
| E-12 | Versionierungsregel Major/Minor/Patch des Lastenhefts | harness/conventions.md:186-191 | Betrifft die Vertragsführung, siehe F-43 |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| spec/lastenheft.md — Abwärts-Zeiger (Suche nach `spezifikation`, `architecture`, ADR, Slice, Welle) | geprüft; Treffer nur in der Regel-Prosa :698 und :3-6 (F-41) |
| spec/spezifikation.md — ID-Zählung SPEC-001 bis SPEC-040 | geprüft, ohne Befund (je einmal definiert; Dateireihenfolge nicht streng aufsteigend: SPEC-034 steht vor SPEC-029, SPEC-035 nach SPEC-033, Verweise lösen unabhängig davon auf) |
| spec/spezifikation.md — Querverweise `LH-FA-05.e`, `LH-FA-13.b`, `SPEC-012`, `SPEC-013` bis `SPEC-019`, `PGR-*` | geprüft, ohne Befund (alle auflösbar; Exit-Code-Ziffer je `PGR-E*` entspricht der Klasse) |
| spec/spezifikation.md — SPEC-034 neue Codes `PGR-E2002`, `PGR-E5002`, `PGR-E6002`, `PGR-E6003` | geprüft, ohne Befund bei Schema `PGR-[EWI][0-9]{4}` und Klassenzuordnung (Lücke: F-39) |
| spec/architecture.md — ADR-/Slice-/Welle-/Commit-Bezüge, Code-Blöcke, `.go`, `os.Exit` | geprüft, ohne Befund (F-29 für Bibliotheksnamen) |
| spec/architecture.md ↔ `.a-check.yml` | geprüft, ohne Befund (F-45) |
| docs/plan/adr/0009-implementierungssprache-go.md | geprüft, ohne Befund (Status `Proposed`, also keine Immutabilitätsverletzung; neuer Anker `#spec-036--technische-leitentscheidungen` löst auf; andere ADR-Anker unverändert gültig) |
| harness/README.md | geprüft, ohne Befund (alle genannten Targets in `make help`; „Nicht behauptet" nennt Build/Test/Architektur-Gate als geplant) |
| harness/conventions.md | geprüft, ohne Befund bis auf F-43 (Platzhalter ersetzt; Baseline v6.13.0, Datum, Präfix `LH`) |
| .harness/skills/reviewer.md | geprüft, ohne Befund (Platzhalter ersetzt; drei Repo-HIGHs; kein Make-Target erfunden) |
| `make docs-check` | 0 Befunde |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 8 |
| LOW | 7 |
| INFO | 2 |

(Neue Findings F-28 bis F-45. Zusätzlich offen aus dem Erstreport: F-11, F-13, F-16, F-17, F-19, F-22 teilweise/offen, F-24 bis F-26 INFO.)

**Finding-Klassen dieses Laufs:** Sicht trägt Sprach-Artefakte · Streaming-Zusage ohne tragende Struktur · Recording-Erzeugung und Replay-Annahme nicht deckungsgleich · Fehlerbedingung ohne Definition · Zustand beim Herunterfahren ungeklärt · Prozess- vs. Verbindungsfehler ungeklärt · Vertragswortlaut weiter als Spec-Verfahren · Neu ergänztes Kriterium ohne Spec-Deckung · Spec verengt Vertrag · Beschriebenes Gate existiert nicht · Boundary-Kriterium ohne Spec-Festlegung · Vorlagen-Prosa im Vertragstext

## Verdikt

**Merge-blockierend:** ja — F-29 (HIGH nach Skill-Wortlaut, Architect bestätigt, ob Format-/Bibliotheksnamen
in der Sicht unter AGENTS.md §3.4 fallen; fällt die Antwort „nein", entfällt das HIGH). Die MEDIUM-Findings
F-28, F-30 bis F-36 betreffen entweder einen Widerspruch Architektur ↔ Spec (F-28) oder eine
Vertrags-/Spec-Spannung, die vor `Accepted` des Lastenhefts entschieden sein muss (F-30, F-33, F-34, F-36).
Die Klasse „Boundary-Kriterium ohne Spec-Festlegung" ist mit F-4, F-5, F-6, F-9, F-12, F-14 im Erstlauf behoben,
tritt aber mit F-40 erneut auf; die Klasse „Sicht trägt Sprach-Artefakte" tritt zum zweiten Mal auf.

**Übergabe:** Findings gehen an den Implementer; F-29, F-30, F-33, F-34, F-36, F-41, F-43 und die
„Entscheidungen zur Bestätigung" E-1 bis E-12 an den Architect beziehungsweise den Auftraggeber (Vertragsrang).
Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
