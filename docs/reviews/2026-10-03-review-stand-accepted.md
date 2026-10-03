# Review-Report: stand-accepted — 2026-10-03

**Review-Art:** Plan und Design (Plan-Review gegen Spec und ADR, Design-Review der Spec-Straten, Konsistenzprüfung der Anwenderdokumente, Prüfung des Harness-Konventionsspeichers gegen die Vorlage).

**Gegenstand:** `git diff 75e14ea HEAD` (HEAD a19ce26) — der Stand seit dem letzten Review. Aufgetragen war `git diff 127d0eb HEAD`; dieser Bereich enthält die genannten Änderungen nicht vollständig (Gruppenregel, Plattformen, Registries, Homebrew, SPEC-040-Entfernung liegen in den Commits davor). Deshalb ist der Bereich ab dem Stand des letzten Reports geprüft: `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, die zwölf ADRs samt Index (alle `Accepted`), `docs/plan/planning/` (Roadmap, vier Wellen, vierzehn Slices, neu `slice-v1-abschluss-homebrew`), `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md`, `harness/conventions.md`.

**Skill:** `.harness/skills/reviewer.md` @ a19ce26
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- `AGENTS.md` (Hard Rules 3.3 bis 3.7, §5, §6), `harness/conventions.md`, `harness/README.md`
- `docs/reviews/2026-10-03-review-extended-query-planung.md` (F-46 bis F-68), `docs/reviews/2026-10-03-spec-folgereview.md` und `docs/reviews/2026-10-03-spec-erstfassung.md` (F-11, F-13, F-17, F-30, F-36, F-39, F-41, F-43)
- die zwölf ADRs, insbesondere [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) und der ADR-Index
- berührte `LH-*`: LH-FA-02, -03, -07, -10, -12, -13, -16, -18, -19, LH-QA-01, -03, -06; Abnahmeszenarien 1 bis 11
- Baseline `v6.13.0` · `regelwerk/modul-03-spec.md`, `modul-04-adrs.md`, `grundlagen-source-precedence.md`; Vorlagen `templates/harness/conventions.template.md`, `templates/spec/lastenheft.template.md`
- Anwender-Vorgabe: `docs/user/benutzerhandbuch-standard.md`

Nicht Gegenstand: DoD-/Abnahme-Konformität (Verifier), fachliche Eignung (Validator). Das Repo enthält keinen Quelltext; Code-Findings entfallen. Pfade sind als Datei plus Kennung beziehungsweise Abschnittsüberschrift angegeben.

---

## Status der offenen Findings

Stand ist HEAD. Erster Teil: F-46 bis F-68 des letzten Reports; zweiter Teil: die früheren Findings.

| Finding | Status | Beleg |
|---|---|---|
| F-46 | teilweise | Port-Tabelle, Domain-Tabelle (Gruppe), Invarianten, Sequenzdiagramm LH-FA-18 und beide Zustandsautomaten in spec/architecture.md tragen Extended; Reste siehe F-79 |
| F-47 | behoben | spec/spezifikation.md · LH-FA-18.a „Gruppen": Zuordnung über `Sync`/`ReadyForQuery`, Freigabe erst nach der letzten Client-Nachricht der Gruppe; `Flush` ohne wartenden Client als dokumentierte Grenze (siehe F-77) |
| F-48 | teilweise | spec/spezifikation.md · SPEC-041 nennt alle Nachrichtentypen, `target`/`name` und `null: true`; Feldschema der Server-Nachrichten fehlt weiter (F-80) |
| F-49 | behoben | spec/spezifikation.md · LH-FA-18.a „Interaktion" und „Abbruch" (nach `Sync` eintreffende Nachrichten, Abbruch auch nach `Terminate` oder `Flush` ohne `Sync`) |
| F-50 | behoben | spec/spezifikation.md · SPEC-001: Version 1 umfasst beide Arten, unbekanntes `type` ist `PGR-E3003`; Rest F-81 |
| F-51 | teilweise | SPEC-040 entfernt, Szenarien 8 bis 11 decken die SOLL-Anforderungen; die Prioritätsdefinition „sofern ohne unverhältnismäßigen Aufwand" steht weiter neben „fertig erst mit SOLL" (F-87) |
| F-52 | behoben | docs/plan/planning/open/slice-replay-semantik-mismatch.md · §2: „merkt sich die Klasse 5", Exit-Code-Abbildung bei `slice-v1-abschluss-betrieb` |
| F-53 | behoben | docs/plan/planning/open/slice-extended-query-replay.md · §2: Erkennung dort, Diagnose und Klasse bei `slice-replay-semantik-mismatch` |
| F-54 | behoben | docs/plan/planning/open/slice-extended-query-replay.md · §2 trägt den Extended-Fehlerreplay (Szenario 6, Extended); welle-extended-query · §3 nennt ihn |
| F-55 | teilweise | Registries benannt (LH-FA-16.a, SPEC-031), „jede Anmeldung" entfällt; Fundort des Binarys und Homebrew-Hinweise ungeklärt (F-74, F-78) |
| F-56 | behoben | LH-FA-05.a, -05.d, -06.a, -10.a, -12.a nennen Extended |
| F-57 | behoben | docs/plan/planning/open/slice-v1-abschluss-protokollrand.md · §2 schneidet den Go-Client-Handshake an `slice-extended-query-replay` ab |
| F-58, F-59, F-60 | behoben | slice-extended-query-modell (Ziel, DoD 1, §3), slice-extended-query-record (Ziel), slice-v1-abschluss-container (Risiko, §3) |
| F-61 | behoben | docs/plan/planning/in-progress/roadmap.md · Drift-Log trägt den Eintrag vom 2026-10-03 |
| F-62, F-63 | behoben | Handbuch: Autor und Gültigkeitsbereich im Kopf; „Mit einem Datenbanktreiber arbeiten" mit Voraussetzung, Vorgehen, Ergebnis; Hinweise nennen alle Vergleichsfelder |
| F-64 | behoben | docs/maintainer/releasing.md · §4 „Festgelegt ist … noch nicht ausgeführt", M3 benannt, §5 „gesetzte Tags"; neue Reste F-71, F-74 |
| F-65 | behoben | spec/spezifikation.md · SPEC-033 (Datenschutz): Klartext-SQL in Diagnosen ausdrücklich, nur Parameterwerte verborgen |
| F-66, F-67 | INFO bleibt | unverändert |
| F-68 | behoben | Roadmap „In Arbeit: nichts"; die Aussage zu „offenen Spec-Punkten" in welle-walking-skeleton ist entfernt |
| F-11 | offen (verlagert) | SPEC-040 entfällt, die Forderung „reproduzierbar gebaut" steht jetzt in slice-v1-abschluss-container · §2, welle-v1-abschluss · §3 und releasing.md · §2, weiter in keiner Lastenheft-Aussage (F-75) |
| F-13 | teilweise | unverändert; kein Diff berührt die Stelle |
| F-17 | teilweise | Extended-Zustände ergänzt; Verbindungsende ohne `Terminate` und Abschluss mit unverbrauchten Interaktionen fehlen weiter (F-79) |
| F-30 | teilweise | Record-Seite gelöst (Verbindung ohne Anfrage wird nicht aufgezeichnet, Mehr-Session-Ablehnung entfällt); Replay-Zuordnung der Probe-Verbindung neu offen (F-73) |
| F-36 | behoben | spec/spezifikation.md · SPEC-035 nennt Linux, macOS, Windows; Setzung siehe F-77 |
| F-39 | teilweise | `PGR-E2002`, `PGR-E5002`, `PGR-E5003`, `PGR-E6002` tragen weiter keine `SPEC-*`-Spalte; `PGR-E6003` entfallen (F-76) |
| F-41 | teilweise | Verweis auf `harness/conventions.md` im Kopf und Abschnitt „Nicht in diesem Dokument entschiedene Punkte" entfernt; die Regel-Prosa der Historie nennt weiter die beiden Spec-Dateien, ADR, Slice, Welle (Vorlagen-Wortlaut als Negation); Architect-Entscheidung steht aus |
| F-43 | offen (verlagert) | Zählregel aus `harness/conventions.md` in den Lastenheft-Kopf verschoben; der Beleg fehlt weiter, jetzt zusätzlich gegen die Baseline-Stelle (F-70, F-71) |
| F-24, F-25, F-44, F-45 | INFO bleibt | unverändert |

## Findings

Neue Findings aus diesem Lauf (Output-Schema des Reviewer-Skills; Nummern setzen die des letzten Reports fort).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-69 | MEDIUM | Die Tabellen „Aktive Adaptionen", „Aufgelöste Adaptionen" und „Zusatzklassen-Deklaration" tragen unausgefüllte Platzhalterzeilen (`<NNN>`, `<Titel>`, `<z. B. LH-Bindung>`, Anker `mr-<NNN>`) als Inhalt; vorher stand dort „Keine …". Die Datei deklariert damit eine Adaption und eine Sensors-Klasse, die es nicht gibt. Gegenüber der Vorlage fehlen außerdem die Link-Form der MR-Zeilen und die Kürzel-Spalte (letztere zulässig, die Vorlage erlaubt das Streichen). | `v6.13.0` · `templates/harness/conventions.template.md` (Zusatzklassen: Tabelle entfernen oder „— keine —"); Hard-Rule-Muster „stille Setzung" | harness/conventions.md · §Adaptions-Block (Aktive/Aufgelöste Adaptionen), §Zusatzklassen-Deklaration | ja — Suche nach `<NNN>`/`<Titel>` | Platzhalter als Inhalt |
| F-70 | MEDIUM | Die Zählregel der Lastenheft-Version (Major/Minor/Patch) steht jetzt im Lastenheft-Kopf; die Baseline verlangt, dass das Repo „welche Stelle steigt" im Adaptions-Block von `harness/conventions.md` deklariert, die Lastenheft-Vorlage verweist dorthin. Die Abweichung hat keinen MR-Eintrag; der Abschnitt „Versionierung des Lastenhefts" in der Konventionsdatei ist gestrichen. | `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §Spec-Stratifizierung („gehört in den Adaptions-Block"); `AGENTS.md` §1 (Strukturregeln leben in conventions) | spec/lastenheft.md · Kopf (Version); harness/conventions.md · §Adaptions-Block | nein | Deklarationsort weicht von Baseline ab |
| F-71 | LOW | `releasing.md` §3 sagt, die Zählregel der Lastenheft-Version stehe in `harness/conventions.md`; dort steht sie nicht mehr. Der Verweis in Inline-Code löst das Doku-Gate nicht auf. | AGENTS.md §3.7 (Kommentar beschreibt, was da ist); Maintainability | docs/maintainer/releasing.md · §3 Versionierung | ja — Suche nach „conventions" im Dokument | Verweis auf entfernte Stelle |
| F-72 | MEDIUM | Der ADR-Index nennt für [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) den Titel „Extended Query als ereignisbasierte Interaktion"; die ADR heißt „Extended Query als Gruppen aus Client- und Server-Nachrichten" und verwirft die ereignisbasierte Lesart als Alternative C. Der Index behauptet damit die verworfene Entscheidung. Die ADR ist `Accepted`, der Index ist fortschreibbar. | AGENTS.md §5 Nr. 3 (Index); `v6.13.0` · `regelwerk/modul-04-adrs.md` | docs/plan/adr/README.md · Zeile der ADR zur Extended-Query-Gruppenregel | ja — Titelvergleich Index gegen Dateikopf | Index-Titel weicht von ADR ab |
| F-73 | MEDIUM | Im Record wird eine Verbindung ohne Anfrage nicht aufgezeichnet (LH-FA-02 Boundary, LH-FA-12.a); im Replay erhält die n-te Verbindung die n-te Session, gezählt über alle Verbindungen. Öffnet die Anwendung im Replay dieselbe Probe- oder Pool-Verbindung ohne Anfrage wie beim Record, verschiebt sie die Zuordnung: die Probe-Verbindung bekommt die erste Session, die echte Verbindung die zweite oder `PGR-E5003`. Das Handbuch nennt die Probe-Verbindung als Beispiel für „nicht aufgezeichnet", ohne diese Folge. | LH-FA-12; LH-QA-01; Nachfolger von F-30 | spec/spezifikation.md · LH-FA-12.a, LH-FA-06/07.a „Schreibzeitpunkt"; docs/user/benutzerhandbuch.md · „Eine Aufzeichnung wiedergeben" (Hinweise) | ja — späterer Integrationstest mit Pool | Aufzeichnungs- und Wiedergabezählung asymmetrisch |
| F-74 | MEDIUM | Die Veröffentlichung von Binaries (mit SHA-256-Summe, von der die Homebrew-Formel liest) und der Images in `ghcr.io` und `docker.io` hat keinen Träger: `slice-v1-abschluss-container` schließt den Registry-Mechanismus aus und verweist auf `releasing.md`, das ausdrücklich keinen Mechanismus kennt; `slice-v1-abschluss-homebrew` startet bei „Release-Artefakte liegen vor". Zusätzlich zirkulär: v1.0.0 wird bei M3 getaggt, M3 verlangt Abnahmeszenario 11 (Installation über Homebrew), die Formel entsteht nur bei einem stabilen Release, eine Vorabversion ändert den Tap nicht. | LH-FA-16, LH-FA-19; Plan-Reihenfolge | docs/plan/planning/open/slice-v1-abschluss-container.md · §1; docs/plan/planning/open/slice-v1-abschluss-homebrew.md · §2, §4; docs/maintainer/releasing.md · §1, §4; docs/plan/planning/in-progress/roadmap.md · M3 | nein | Lieferpfad ohne tragenden Slice, Meilenstein zirkulär |
| F-75 | MEDIUM | `releasing.md` §2 behauptet, die Anforderungen an Container und Reproduzierbarkeit stünden im Lastenheft (`LH-FA-16`, `LH-QA-03`); das Lastenheft fordert Reproduzierbarkeit des Builds nirgends. Dieselbe Forderung steckt als Closure-Trigger in `welle-v1-abschluss` und als DoD in `slice-v1-abschluss-container`. Nach der Streichung von SPEC-040 ist die Spezifikation nicht mehr die Quelle, die Aussage ist eine Setzung ohne Vertragsgrundlage. | Source Precedence (Lastenheft vor Spec); `v6.13.0` · `regelwerk/modul-03-spec.md` („Präzisieren, nie erweitern"); F-11 | docs/maintainer/releasing.md · §2; docs/plan/planning/welle-v1-abschluss.md · §3; docs/plan/planning/open/slice-v1-abschluss-container.md · §2 | nein | Setzung als Vertragsaussage ausgewiesen |
| F-76 | MEDIUM | `PGR-E5003` gilt als „Replay-Mismatch" (LH-FA-03.a, LH-FA-12.a), SPEC-027 bildet die Klasse aber nur auf `PGR-E5001` ab; die Diagnose nach LH-FA-10.a (erwartete Interaktion, Nachrichtentyp) ist ohne erwartete Interaktion nicht erfüllbar. Eine aufgezeichnete, nie verbundene Session unter `--fail-on-unconsumed` hat keine Verbindung, die „als fehlerhaft beendet" zählt (LH-FA-03.b); Zeitpunkt der Feststellung ungeregelt. `PGR-W2001` beschreibt weiter nur eine endende Sitzung. Codes `PGR-E5002`, `PGR-E5003` ohne `SPEC-*`-Bezug. | LH-FA-03, LH-FA-10, LH-FA-13; F-39 | spec/spezifikation.md · LH-FA-03.a Nr. 4, LH-FA-03.b, LH-FA-10.a, SPEC-027, SPEC-034 (Katalog) | ja — Spec-Lesung | Fehlerbedingung ohne Definition |
| F-77 | MEDIUM | Spezifikation, Handbuch und Releasing legen einen Lieferumfang fest, den das Lastenheft nicht nennt (Setzungen, Bestätigung des Auftraggebers nötig, siehe Tabelle unten): Plattformmatrix Linux/macOS/Windows × `amd64`/`arm64` (SPEC-035; das Lastenheft nennt Plattformen nur in LH-QA-03 als „typische Umgebungen" und in LH-FA-19 Negative als „Windows"), Registries und deren Namen, Tap-Name und Formel-Aufbau, Tag `v1.0.0` bei M3, Windows-`SIGINT`-Gleichsetzung, Flush-Grenze der Aufzeichnung. | `v6.13.0` · `regelwerk/modul-03-spec.md` §Ziel-Form: Spezifikation („Präzisieren, nie erweitern"); Source Precedence | spec/spezifikation.md · SPEC-035, LH-FA-16.a, SPEC-031, LH-FA-19.a, SPEC-042, LH-FA-13.a; docs/maintainer/releasing.md · §2, §4 | nein | Spec erweitert Vertrag |
| F-78 | MEDIUM | Das Handbuch behauptet Drittverhalten ohne Beleg: Aktuelle Homebrew-Versionen luden Drittanbieter-Taps nur nach `brew trust`; Lastenheft und Spezifikation nennen den Schritt nicht, der Slice führt Homebrews Verhalten selbst als offenes Risiko. Abnahmeszenario 11 bindet den Vertrag an „dokumentierte Befehle". Zusätzlich nennt der Abschnitt „Binary" keinen Fundort (Schritt 1 „Legen Sie das Binary … in ein Verzeichnis"). | `docs/user/benutzerhandbuch-standard.md` §13 (Funktionen vorhanden/belegt); LH-FA-19 | docs/user/benutzerhandbuch.md · §2 Installation (Binary, Homebrew); docs/plan/planning/open/slice-v1-abschluss-homebrew.md · §6 | nein | Handbuch behauptet über die Spec hinaus |
| F-79 | LOW | Architektur-Reste: Das Sequenzdiagramm LH-FA-18 lässt Server-Nachrichten „bis ReadyForQuery" enden, eine Gruppe endet aber auch mit `Flush`; die Automaten kennen keinen Übergang bei Verbindungsende oder `unsupported` aus `RECORDING_GROUP`/`MATCH_GROUP` und keinen Abschluss mit unverbrauchten Interaktionen; „Verbindung *n* → Application Session *n*" gilt im Record nicht, weil Verbindungen ohne Anfrage keine Session erzeugen. | LH-FA-12; LH-FA-18; F-17 | spec/architecture.md · Use-Case LH-FA-18, §Zustandsautomaten Record/Replay, §Nebenläufigkeit | nein | Sicht deckt Spec-Pfad nicht ab |
| F-80 | LOW | SPEC-041 gibt für `error_response`, `notice_response`, `parameter_description`, `row_description` keine Feldform an („protokollrelevante Felder" als Kommentar im Beispiel, wie in SPEC-002); `close` und `flush` fehlen im Beispiel. | LH-QA-06 | spec/spezifikation.md · SPEC-041 | nein | Format-Beispiel ohne Ereignisvollständigkeit |
| F-81 | LOW | SPEC-001 nennt „`type: query`, `type: extended`" als Arten der Interaktion; `query` steht aber unter `request.type`, `extended` auf der Interaktion selbst (SPEC-002, SPEC-041). | LH-QA-06 | spec/spezifikation.md · SPEC-001, SPEC-002, SPEC-041 | nein | Schlüsselname an unterschiedlicher Ebene |
| F-82 | LOW | Satzfehler nach Teilersetzung: „den Index des Gruppe und Nachricht in der Interaktion". | Maintainability | spec/spezifikation.md · LH-FA-18.a „Mismatch" | ja — Lektorat | Satz nach Teilersetzung gebrochen |
| F-83 | LOW | Das Handbuch verlangt für die Wiedergabe als Voraussetzung „Eine Aufzeichnung mit einer Verbindung liegt vor" — Rest der Ein-Session-Lesart. | LH-FA-12 | docs/user/benutzerhandbuch.md · „Eine Aufzeichnung wiedergeben" (Voraussetzung) | ja — Suche | Single-Session-Rest |
| F-84 | LOW | Tap-Name uneinheitlich: SPEC-042 und `releasing.md` nennen „Homebrew-Tap `pt9912/homebrew-pgwire-recorder`" (Repository-Name), Handbuch und LH-FA-19.a verwenden `pt9912/pgwire-recorder` als Tap (`brew tap`). | Maintainability | spec/spezifikation.md · SPEC-042, LH-FA-19.a; docs/maintainer/releasing.md · §2; docs/user/benutzerhandbuch.md · Homebrew | ja — Suche | Repository- und Tap-Name vermischt |
| F-85 | LOW | Das Lastenheft schreibt „Verbindung ohne Anfrage" nicht aufzeichnen (LH-FA-02 Boundary); die Spezifikation verlangt „mindestens eine abgeschlossene Interaktion". Eine Verbindung mit abgebrochener erster Anfrage ist damit im Lastenheft aufzeichnungsfähig, in der Spezifikation nicht. | LH-FA-02; LH-FA-12 | spec/spezifikation.md · LH-FA-12.a, LH-FA-07.a „Schreibzeitpunkt"; spec/lastenheft.md · LH-FA-02 | nein | Spec enger als Vertrag |
| F-86 | LOW | Commit-Nachrichten tragen die falsche oder keine Anforderungs-Kennung: `97363c0` nennt LH-FA-18 für eine Versionsregel; `f74057d` ändert Inhalt zu LH-FA-12 und LH-FA-16 ohne diese Kennungen. | AGENTS.md §5 Nr. 1 | Commits 97363c0, f74057d | ja — `tools/harness/commit-msg-traceability.sh` prüft nur die Form | Kennung in Commit unpassend |
| F-87 | LOW | Priorität SOLL ist „erwünscht, sofern ohne unverhältnismäßigen Aufwand realisierbar", Zielversion und Kopf sagen „fertig erst mit SOLL"; wer „unverhältnismäßig" feststellt, ist nicht benannt. | Source Precedence (Lastenheft) | spec/lastenheft.md · Priorität (Kopf von §3), Zielversion | nein | Vertragsschwelle ohne Entscheider |
| F-88 | INFO | Die Alternativen der [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) bis [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) zählen drei Optionen einschließlich der gewählten (zwei verworfene); [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) vier. Die Vorlage verlangt „mindestens drei Optionen", das Regelwerk „drei verglichene Alternativen"; „nichts tun" fehlt überall. Seit `Accepted` nicht mehr änderbar. | `v6.13.0` · `regelwerk/modul-04-adrs.md` §Ziel-Form: ADR | docs/plan/adr/ · §Verglichene Alternativen aller ADRs | nein | Alternativen-Zählung mehrdeutig |
| F-89 | INFO | Die Umbenennung von [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) erzwang einen Nachtrag der Links im eingefrorenen Report `docs/reviews/2026-10-03-review-extended-query-planung.md` (nur Link-Ziele, kein Befund geändert). Außerdem verweisen die `Accepted`-ADRs auf Anforderungen eines Lastenhefts im Status `Draft`, das bis `Accepted` frei änderbar bleibt. | Kennung-statt-Adresse-Regel der Report-Vorlage | docs/reviews/2026-10-03-review-extended-query-planung.md · Kopf; docs/plan/adr/ · Feld `Bezug` | nein | Adresse im Lauf-Beleg |
| F-90 | INFO | `SPEC-040` und `PGR-E6003` sind ersatzlos gestrichen; die Lücke ist nirgends als vergeben vermerkt. Die Spezifikation sagt für Meldungscodes „entfallen: zurückgezogen und bleibt vergeben"; beide waren nie veröffentlicht. | `v6.13.0` · `regelwerk/grundlagen-source-precedence.md` §Vergabe (Lücken werden nicht nachbelegt) | spec/spezifikation.md · §Meldungscodes (Stabilität), Abschnitt Historie | nein | Lücke ohne Vermerk |

## Setzungen, die der Auftraggeber bestätigen muss

Setzungen des Assistenten oder Reviewers in diesem Stand, soweit sie über den Wortlaut des Lastenhefts hinausgehen (Vertragsrang; Details zu F-77).

| Nr. | Setzung | Fundort | Überschreitet das Lastenheft? |
|---|---|---|---|
| S-1 | Zielplattformen Linux, macOS, Windows jeweils `amd64`/`arm64`, Image Linux `amd64`/`arm64` als Manifestliste | SPEC-035, LH-FA-16.a, Handbuch, Releasing | ja — Lastenheft nennt keine Plattformen (LH-QA-03 „typische Umgebungen"); verpflichtet zu Windows-Build und -Test |
| S-2 | Veröffentlichung in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder` | LH-FA-16.a, SPEC-031, Releasing | ja — LH-FA-16 verlangt „bereitgestellt", keine Registry |
| S-3 | Homebrew über eigenen Tap, Formel installiert veröffentlichtes Binary mit SHA-256-Prüfung, kein Quellbau, Formel nur bei stabilem Release | LH-FA-19.a, SPEC-042 | teilweise — LH-FA-19 fordert Installierbarkeit und Vorabversions-Grenze, nicht den Aufbau |
| S-4 | n-te Verbindung = n-te Session; Verbindung ohne Anfrage wird nicht aufgezeichnet; Mehr-Verbindungen als Sessions = Mismatch `PGR-E5003` | LH-FA-12 (Lastenheft selbst), LH-FA-12.a | nein — im Lastenheft entschieden; Folge für Pool-Probe-Verbindungen offen (F-73) |
| S-5 | `Flush` ohne wartenden Client: Aufzeichnung zeitabhängig, bekannte Grenze von v1 | LH-FA-18.a, [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | nicht ausgeschlossen — LH-QA-01 prüft Replay, nicht Record; Lastenheft nennt die Grenze nicht |
| S-6 | Zählregel Major/Minor/Patch und „SemVer 2.0" für die Lastenheft-Version | Lastenheft-Kopf | Lastenheft-intern; Ort weicht von Baseline ab (F-70) |
| S-7 | Software-Version `v<SemVer>`, `v1.0.0` bei M3, keine Veröffentlichung bis dahin | Releasing | Prozess, nicht Vertrag; Zirkel mit Szenario 11 (F-74) |
| S-8 | Windows-Konsolenabbruch (`Strg+C`, `Strg+Break`) entspricht `SIGINT` | LH-FA-13.a | ja — Lastenheft nennt Signale und Windows nicht |
| S-9 | Reproduzierbarer Build für Binary und Image als Fertigstellungs- und Closure-Bedingung | Releasing, Welle, Slice | ja — Lastenheft fordert ihn nicht (F-75) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Zählung LH-FA-01 bis LH-FA-19, LH-QA-01 bis -06, LH-RB-01 und -02 | geprüft, ohne Befund; Überschriften lückenlos, Querverweise lösen auf |
| Zählung SPEC-001 bis SPEC-042 | geprüft, ohne Befund außer der dokumentierten Lücke SPEC-040 (F-90); SPEC-041, SPEC-042 je einmal definiert |
| Zählung ARC-001 bis ARC-013 | geprüft, ohne Befund |
| Abnahmeszenarien 1 bis 11 (Lastenheft §7) und Zählung „elf" in welle-v1-abschluss, Roadmap M3 | geprüft, ohne Befund; Szenarien 8 bis 11 decken LH-FA-14, -16, -17, -19 |
| Reste `PGR-E6003`, `SPEC-040`, „zehn" (Abnahme), „Linux only", „Single-Session", `Mehr-Session` | geprüft; nur F-83 (Voraussetzung im Handbuch) und die Fundstellen in älteren, eingefrorenen Reports |
| spec/lastenheft.md — Decken-Regel (Hard Rule 3.4) | geprüft; kein Verweis auf ADR, Slice, Welle in Anforderungen; Rest der Regel-Prosa siehe F-41 im Status |
| spec/architecture.md — Hard Rule 3.4 | geprüft, ohne Befund; keine Welle, Slice, ADR, kein Hash; Pfade zu Code-Modulen erlaubt |
| spec/spezifikation.md — SPEC-041 gegen [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | geprüft, ohne Befund bei Inhalt: Gruppe endet mit `Flush` oder `Sync`, Freigabe nach letzter Client-Nachricht, strict sequential; Schärft-Anker lösen auf |
| ADR-Konformität: Status, Immutabilität, Geschichte | geprüft, ohne Befund; [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) bis [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) ändern nur Statuszeile und Geschichte; Umbenennung von [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) als reiner Move (`4efc06e`) |
| Aussagen in Spec oder Plan gegen `Accepted`-ADRs | geprüft, ohne Widerspruch; welle-walking-skeleton nennt alle ADRs als `Accepted`; [ADR-0007](../plan/adr/0007-strict-replay.md) und [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) mit n-ter Session verträglich |
| ADR-Index, Statusspalte, Schärft-Links aller ADRs | geprüft; Status und Links ohne Befund, Titel von [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) siehe F-72 |
| harness/conventions.md gegen `conventions.template.md` | geprüft; ausgefüllte Felder (Stand `v6.13.0`, Datum, Modus-Deklaration, Glossar) ohne Befund; strukturelle Abweichungen F-69 und F-70 |
| docs/plan/planning — Wellen, Slices, Roadmap, `slice-v1-abschluss-homebrew` §1 bis §8 | geprüft, ohne Befund bei Form (≤ 3 Liefer-Punkte, Risiko-Ausgang, Rückführungen); inhaltliche Befunde F-74, F-75 |
| docs/user/benutzerhandbuch.md — interne Kennungen | geprüft, ohne Befund; nur `PGR-…` |
| `make docs-check` | 0 Befunde (fängt keines der Findings oben) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 9 |
| LOW | 10 |
| INFO | 3 |

(Zusätzlich offen aus früheren Läufen: F-11 (verlagert, F-75), F-13, F-17, F-30, F-39, F-41, F-43; INFO F-24, F-25, F-44, F-45, F-66, F-67.)

**Finding-Klassen dieses Laufs:** Platzhalter als Inhalt · Deklarationsort weicht von Baseline ab · Verweis auf entfernte Stelle · Index-Titel weicht von ADR ab · Aufzeichnungs- und Wiedergabezählung asymmetrisch · Lieferpfad ohne tragenden Slice · Setzung als Vertragsaussage ausgewiesen · Fehlerbedingung ohne Definition · Spec erweitert Vertrag · Handbuch behauptet über die Spec hinaus · Sicht deckt Spec-Pfad nicht ab · Format-Beispiel ohne Ereignisvollständigkeit · Schlüsselname an unterschiedlicher Ebene · Satz nach Teilersetzung gebrochen · Single-Session-Rest · Repository- und Tap-Name vermischt · Spec enger als Vertrag · Kennung in Commit unpassend · Vertragsschwelle ohne Entscheider · Alternativen-Zählung mehrdeutig · Adresse im Lauf-Beleg · Lücke ohne Vermerk

Hinweis zum Zähler: „Fehlerbedingung ohne Definition" (F-31, F-49, F-76) tritt zum dritten Mal auf — Steering-Loop-Schwelle erreicht (Klassifikation prüfen, Folge-ADR oder Fitness Function erwägen). „Spec erweitert Vertrag" (F-11 des Erstreviews, F-36, F-77) ebenfalls dreimal. „Sicht deckt Spec-Pfad nicht ab" (F-17, F-46, F-79) dreimal. „Zuständigkeit … doppelt" trat nicht erneut auf.

## Verdikt

**Merge-blockierend:** ja — für F-69, F-70, F-72 bis F-78 (MEDIUM): Konventionsspeicher mit Platzhalterinhalt und Baseline-Abweichung, ADR-Index mit der verworfenen Lesart, Session-Zuordnung bei Probe-Verbindungen, fehlender Lieferpfad samt Release-Zirkel (Szenario 11), Setzungen jenseits des Lastenhefts. Kein HIGH: Es besteht kein ADR-Verstoß, keine Verletzung von Hard Rule 3.3 bis 3.6, kein Artefakt-Verweis in den Spec-Straten nach Skill-Wortlaut außer dem bekannten Vorlagen-Rest (F-41). [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) passt zur Gruppenregel in LH-FA-18.a und SPEC-041.

**Übergabe:** Findings gehen an den Implementer; F-70, F-73, F-74, F-76, F-77, F-87 und die Setzungstabelle (S-1 bis S-9) an den Architect beziehungsweise den Auftraggeber (Vertragsrang); F-74 an den Planner (Rückkante Review → Plan). Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
