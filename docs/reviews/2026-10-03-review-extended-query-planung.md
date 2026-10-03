# Review-Report: extended-query-planung — 2026-10-03

**Review-Art:** Plan und Design (Plan-Review gegen Spec/ADR/Hard Rules für die drei neuen Wellen und zehn Slices; Design-Review der Spec-Straten nach der Verschiebung von Extended Query in v1; Konsistenzprüfung der beiden neuen Anwenderdokumente).

**Gegenstand:** `git diff 77a8c93 HEAD` (HEAD 75e14ea): `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) samt ADR-Index, `docs/plan/planning/` (drei Welle-Dateien, zehn Slices in `open/`, Roadmap, Änderung an `welle-walking-skeleton`), `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md`. Zusätzlich gelesen, nicht Teil des Diffs: `docs/user/benutzerhandbuch-standard.md`.

**Skill:** `.harness/skills/reviewer.md` @ 75e14ea (ausgefüllt; Repo-HIGHs: Spec-Stratum nennt Artefakt, ADR/Move-Regel, Core-Reinheit)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- `AGENTS.md` (Hard Rules 3.3 bis 3.7, §4, Source Precedence), `harness/conventions.md`
- `docs/reviews/2026-10-03-spec-folgereview.md` (Findings F-28 bis F-45, Entscheidungen E-1 bis E-12) und `docs/reviews/2026-10-03-spec-erstfassung.md`
- [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) (Bezug des Diffs), [ADR-0007](../plan/adr/0007-strict-replay.md) (Verweisziel)
- berührte `LH-*`: LH-FA-04, -05, -06, -09, -10, -11, -12, -13, -18, LH-QA-01, -03, -06; Abnahmeszenarien 1 bis 7
- Baseline `v6.13.0` · `regelwerk/modul-03-spec.md`, `modul-05-planning-harness.md`, `modul-06-roadmap.md`, `grundlagen-referenz-richtung.md`
- Anwender-Vorgabe: `docs/user/benutzerhandbuch-standard.md`
- `make help`/`Makefile` (vorhandene Targets: `gates`, `help`, `docs-check`, `baseline-verify`, `slice-mv`, `archive-welle`, `hooks-install`)

Nicht Gegenstand: DoD-/Abnahme-Konformität (Verifier), fachliche Eignung der Vorgaben (Validator). Das Repo enthält keinen Quelltext; Code-Findings entfallen.

---

## Status der offenen Findings aus dem Folgereview

Stand ist HEAD. Nur Findings, die im Folgereview offen oder teilweise waren.

| Folge-ID | Status | Beleg |
|---|---|---|
| F-11 | offen | SPEC-040 Nr. 8 fordert Binary und Image reproduzierbar gebaut (spec/spezifikation.md:925-926); das Lastenheft fordert Reproduzierbarkeit des Builds nirgends (LH-QA-03 nur „bereitstellbar") |
| F-13 | teilweise | unverändert; Spec-Gegenstück für nicht verlustfrei darstellbare Serverantwort steht jetzt in LH-FA-05.a (spec/spezifikation.md:188-190), damit F-35 behoben |
| F-16, F-37 | behoben | Gate in Vorgabe-Form (spec/spezifikation.md:853-855, spec/architecture.md:488-490); das Make-Target plant `slice-walking-skeleton-build-gates` |
| F-17 | offen | Zustandsautomaten spec/architecture.md:366-385 unverändert; siehe auch F-46 |
| F-19, F-28 | behoben | `Save` nach Session-Ende (spec/architecture.md:196, 270-271), Speicher-Aussage deckungsgleich mit SPEC-032 (spec/architecture.md:199-205) |
| F-22 | behoben | LH-QA-03 „bei der Abnahme benannt" (spec/lastenheft.md:557-558) |
| F-24 | offen (INFO bleibt) | Lastenheft `Draft`; LH-FA-18 ist ein weiterer Zusatz vor `Accepted` |
| F-25 | offen (INFO bleibt) | `.d-check.yml` ohne aktives `LH-*`-Muster; manuelle Prüfung: alle neuen Verweise lösen auf, `make docs-check` 0 Befunde |
| F-29 | behoben (Rest INFO) | `pgproto3` kommt in spec/architecture.md nicht mehr vor; „YAML-Serialisierungsbibliothek" bleibt als Technologiename (spec/architecture.md:38, 94, 109, 259) |
| F-30 | offen | LH-FA-07.a (spec/spezifikation.md:299-301) und LH-FA-12.a (:392-394): Probe-/Pool-Verbindung erzeugt leere Zweit-Session, die im Replay abgelehnt wird; das Handbuch weicht darauf aus (docs/user/benutzerhandbuch.md:130-133, 259) |
| F-31 | teilweise | „Unerwartetes Verbindungsende" jetzt für Simple definiert (spec/spezifikation.md:96-103); für Extended offen, siehe F-49 |
| F-32, F-33, F-34, F-35, F-38, F-40, F-42 | behoben | LH-FA-13.a/13.b (spec/spezifikation.md:405-440), LH-FA-09 Boundary „Position" (spec/lastenheft.md:315-316), SPEC-004 (spec/spezifikation.md:556-559), SPEC-039 (:906-909) |
| F-36 | offen | SPEC-035 nur Linux (spec/spezifikation.md:816); das Handbuch übernimmt die Verengung (docs/user/benutzerhandbuch.md:36, 52) |
| F-39 | teilweise | SPEC-014 nennt `--output`; Codes `PGR-E2002`, `PGR-E5002`, `PGR-E6002`, `PGR-E6003` tragen weiter keine `SPEC-*`-Spalte (spec/spezifikation.md:743, 755, 758-759) |
| F-41 | offen | Lastenheft nennt `harness/conventions.md` (spec/lastenheft.md:5) und `spezifikation.md`/`architecture.md` in der Regel-Prosa (:734); §8 listet weiter Punkte, die das Lastenheft inzwischen selbst entscheidet. Nach Skill-Wortlaut HIGH-Muster, im Folgereview als LOW geführt; Architect-Entscheidung steht aus |
| F-43 | offen | Versionierungsregel ohne Beleg unverändert (harness/conventions.md:186-191) |
| F-44, F-45 | INFO bleibt | unverändert |

## Findings

Neue Findings aus diesem Lauf (Output-Schema des Reviewer-Skills; Nummern setzen die des Folgereviews fort).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-46 | MEDIUM | Die Architektur nennt den neuen Interaktionstyp nur im Fließtext (Antworttypen, Matcher). Outbound-Port „je Session eine Query ausführen", Domain-Tabelle (Interaction = Request + Responses; Query = SQL-Text), Invariante „jede Interaktion besitzt einen gültigen Request", beide Sequenzdiagramme und beide Zustandsautomaten kennen weder Client-Ereignisse außer `Query` noch eine mehrteilige Interaktion; Anfrage-Typen (`Parse`, `Bind`, …) fehlen, nur Antwort-Typen sind genannt. | LH-FA-18; AGENTS.md §3.4 (Sicht trägt die Aussage, nicht nur ein Satz) | spec/architecture.md:196-197, 211-217, 238-242, 262-300, 366-385 | nein | Sicht deckt neuen Interaktionstyp nur im Fließtext |
| F-47 | MEDIUM | LH-FA-18.a zeichnet Client- und Server-Ereignisse „in Ankunftsreihenfolge" auf, der Replay gibt Server-Ereignisse „bis zum nächsten Client-Ereignis, das er noch nicht empfangen hat" frei. Wo ein Server-Ereignis relativ zu einem Client-Ereignis steht, hängt bei Pipelining vom Zeitverhalten der Aufzeichnung ab; ein Recording, das dieselbe Anwendung erzeugt, kann dadurch verschieden ausfallen (LH-QA-01) und im Replay andere Freigabepunkte haben. | LH-QA-01; LH-FA-18; [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | spec/spezifikation.md:519-523, 529-533 | ja — späterer Integrationstest (zwei Aufzeichnungen desselben Ablaufs) | Ereignisreihenfolge zeitabhängig |
| F-48 | MEDIUM | SPEC-041 zeigt `parse`, `bind`, `execute`, `sync` und fünf Server-Ereignisse; für `describe`, `close`, `flush`, `parameter_description`, `no_data`, `portal_suspended`, `empty_query_response`, `error_response`, `notice_response`, `close_complete` und NULL-Parameter gibt es keine Form. Die Matching-Feldliste nennt für `Describe`/`Close` weder die Zielart (Statement oder Portal) noch für `Flush` ein Feld; `result_formats` steht im Beispiel, nicht in der Liste. | LH-FA-18; LH-QA-06 | spec/spezifikation.md:526-530, 607-642 | nein — Schema-Tests fehlen noch | Format-Beispiel ohne Ereignisvollständigkeit |
| F-49 | MEDIUM | Der Abschluss einer Extended-Interaktion („`ReadyForQuery` nach `Sync`") ist nur im Normalfall beschrieben. Offen: Abbruch ohne `Sync` oder nur mit `Flush` (LH-FA-02.b, LH-FA-13.a kennen nur die Simple-Lesart), Client-Nachrichten nach `Sync` vor dem `ReadyForQuery` (Pipelining der nächsten Interaktion) und deren Zuordnung. Berührt „unvollständige Interaktion wird nicht übernommen" aus F-31. | LH-FA-02, LH-FA-13, LH-FA-18 | spec/spezifikation.md:91-103, 515-518; :412-416 | ja — Spec-Lesung | Fehlerbedingung ohne Definition |
| F-50 | MEDIUM | SPEC-041 führt eine neue Interaktionsart (`type: extended`) ein, SPEC-001 und SPEC-010 lassen die Version bei `1`; ob ein Leser von Version 1 ohne Kenntnis von `extended` ablehnt (`PGR-E3002`/`PGR-E3003`) oder nicht, ist nicht gesagt. `slice-extended-query-modell` führt „braucht eine neue Version" als offenes Risiko. | LH-QA-06; SPEC-001 | spec/spezifikation.md:564-575, 607-645, 678; docs/plan/planning/open/slice-extended-query-modell.md:98 | ja — Spec-Lesung | Versionierung einer Formatänderung offen |
| F-51 | MEDIUM | Lastenheft setzt drei Schwellen: SOLL „sofern ohne unverhältnismäßigen Aufwand realisierbar" (:99-100) neben „Fertig ist das Produkt erst, wenn auch die SOLL-Anforderungen umgesetzt sind" (:101-102, :16); §7 nennt „abnahmefähig" mit Szenarien 1 bis 7, die keine der SOLL-Anforderungen (LH-FA-14, -16, -17) prüfen. SPEC-040 führt daneben zwei Fertig-Begriffe („Das Produkt ist fertig …", „v1 ist technisch fertig …"). Vertragswortlaut ohne Auflösung, welche Schwelle gewinnt. | Source Precedence (Lastenheft vor Spec) | spec/lastenheft.md:16, 99-103, 637-639; spec/spezifikation.md:914-933 | nein | Vertragsschwelle widersprüchlich |
| F-52 | MEDIUM | `slice-replay-semantik-mismatch` fordert in der DoD „Der Prozess endet nach einem Mismatch mit Exit-Code 5", das Ziel sagt „meldet beim Beenden", die Spec legt fest: Verbindung endet, Prozess läuft weiter (LH-FA-13.b). Der Exit-Code entsteht erst nach kontrolliertem Herunterfahren, das `slice-v1-abschluss-betrieb` in der späteren Welle liefert; `slice-replay-semantik-meldungscodes` testet Klassen (`SPEC-026`, `SPEC-028`, `PGR-E2002`), deren Auslöser in `welle-v1-abschluss` entstehen. Closure-Trigger „Abnahmeszenario 4 end-to-end" der früheren Welle hängt damit von einer späteren ab. | LH-FA-13; Plan-Reihenfolge | docs/plan/planning/open/slice-replay-semantik-mismatch.md:32, 48, 99; docs/plan/planning/open/slice-replay-semantik-meldungscodes.md:47-48; docs/plan/planning/open/slice-v1-abschluss-betrieb.md:32; docs/plan/planning/welle-replay-semantik.md:38-39 | ja — Plan-Lesung | Slice-DoD verlangt Zustand einer späteren Welle |
| F-53 | MEDIUM | LH-FA-10 für Extended-Nachrichten ist zweifach zugewiesen: `slice-extended-query-replay` (DoD 2, „wie ein Simple-Query-Mismatch gemeldet") und `slice-replay-semantik-mismatch` (Ziel „… oder Extended-Query-Nachricht"). Der Mismatch-Pfad mit Diagnose, Index und Nachrichtentyp (LH-FA-18.a) entsteht erst in der späteren Welle; der frühere Slice kann ihn nicht an einem Referenzverhalten messen. | LH-FA-10; LH-FA-18; Plan-Reihenfolge | docs/plan/planning/open/slice-extended-query-replay.md:32, 47; docs/plan/planning/open/slice-replay-semantik-mismatch.md:32; docs/plan/planning/welle-extended-query.md:54 | ja — Plan-Lesung | Zuständigkeit für dieselbe Anforderung doppelt |
| F-54 | MEDIUM | `welle-replay-semantik` verspricht Abnahmeszenario 4 und 6 „für Simple und Extended Query" (Ziel). Kein Slice trägt Extended-Fehlerreplay (Verwerfen bis `Sync`, LH-FA-18.a „Fehler"): `slice-replay-semantik-fehlerreplay` ist auf Simple-Antworten geschnitten, `slice-extended-query-replay` hat keinen Fehlerfall in der DoD. | LH-FA-11; LH-FA-18; welle-replay-semantik Ziel | docs/plan/planning/welle-replay-semantik.md:21; docs/plan/planning/open/slice-replay-semantik-fehlerreplay.md:32, 47-48; docs/plan/planning/open/slice-extended-query-replay.md:46-47 | ja — Plan-Lesung | Wellen-Ziel ohne tragenden Slice |
| F-55 | MEDIUM | Das Handbuch behauptet Dinge, die die Spezifikation nicht festlegt: (a) „Laden Sie das Binary … herunter" und ein Image `pgwire-recorder:latest` ohne definierten Verteilungsweg; `slice-v1-abschluss-container` schließt die Registry-Veröffentlichung aus, `releasing.md` benennt keinen Veröffentlichungsort. (b) „beim Wiedergeben nimmt es jede Anmeldung an": LH-FA-05.b legt nur „für Testclients verwendbarer Ablauf" fest, die emulierten Nachrichten sind offen. | Handbuch-Standard §13 (Funktionen vorhanden/belegt); LH-FA-05.b | docs/user/benutzerhandbuch.md:51-71, 230-231; spec/spezifikation.md:205-216; docs/plan/planning/open/slice-v1-abschluss-container.md:36 | nein | Handbuch behauptet über die Spec hinaus |
| F-56 | LOW | Reste „nur Simple" nach der Verschiebung: Überschrift „Unterstützter Umfang: Simple Query Protocol" und „Für den Kern von v1 sind relevant: `Query`, `Terminate`" (LH-FA-05.a); LH-FA-06.a nennt nur „Query", „Anfrageart", „SQL-Text"; LH-FA-12.a verweist auf „geordnete Response-Liste (SPEC-002)", für Extended gilt `exchange` (SPEC-041); LH-FA-05.d nur Simple; LH-FA-10.a nennt „weitere Queries" ohne Extended-Pendant. Kein Widerspruch zu LH-FA-18, aber Lesepfad führt in die alte Lesart. | Maintainability | spec/spezifikation.md:157-168, 271-278, 384-385, 236-238, 335-338 | nein | Simple-only-Reste nach Scope-Erweiterung |
| F-57 | LOW | Handshake-Zuständigkeit überlappt: `slice-v1-abschluss-protokollrand` (DoD 2) liefert „Handshake, den ein verbreiteter Go-Client ohne Datenbank akzeptiert", Meilenstein M2 verlangt aber schon in `welle-extended-query` Abnahmeszenario 7 mit Standardtreiber im Replay; `slice-extended-query-replay` nennt „Handshake und Nachrichten". | Plan-Reihenfolge; M2 | docs/plan/planning/open/slice-v1-abschluss-protokollrand.md:47; docs/plan/planning/open/slice-extended-query-replay.md:46, 68; docs/plan/planning/in-progress/roadmap.md:45 | ja — Plan-Lesung | Zuständigkeit für dieselbe Anforderung doppelt |
| F-58 | LOW | `slice-extended-query-modell` plant Spec-Zeilen „`LH-FA-18.a` (neu)" und eine ADR „neu"; beides liegt in HEAD vor (LH-FA-18.a, SPEC-041, [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) Proposed). Liefergegenstand und DoD 1 beschreiben bereits geleistete Arbeit; offen bleiben Domain-Modell und `Accepted`. | Plan-Konsistenz | docs/plan/planning/open/slice-extended-query-modell.md:30-32, 46, 61-62, 15 | ja — Plan-Lesung | Slice beschreibt bereits geliefertes Ergebnis |
| F-59 | LOW | `slice-extended-query-record` nennt als Ziel `Parse`, `Bind`, `Describe`, `Execute`, `Sync`; LH-FA-18.a legt zusätzlich `Close` und `Flush` fest. | LH-FA-18 | docs/plan/planning/open/slice-extended-query-record.md:32; spec/spezifikation.md:507-509 | nein | Slice-Ziel schmaler als Spec |
| F-60 | LOW | `slice-v1-abschluss-container`: Risiko „`docs/user/` existiert noch nicht" (:98) und §3 „Betriebsdokumentation unter `docs/user/` — neu" (:66) sind überholt; `docs/user/benutzerhandbuch.md` und `docs/maintainer/releasing.md` liegen vor. Welcher Teil der Betriebsdokumentation der Slice noch liefert, ist offen. | Plan-Konsistenz | docs/plan/planning/open/slice-v1-abschluss-container.md:66, 98 | nein | Slice-Plan nennt überholten Ist-Stand |
| F-61 | LOW | Die Roadmap ordnet um (neue Welle `welle-extended-query` zwischen Skelett und Replay-Semantik, Start-Trigger von `welle-replay-semantik` von „walking-skeleton done" auf „extended-query done" geändert, M2 neu belegt, M3 neu); die Tabelle *Historische Trigger-Verschiebungen* bleibt leer. | `v6.13.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Drift-Log) | docs/plan/planning/in-progress/roadmap.md:42-47, 56-67 | nein | Umplanung ohne Drift-Log-Eintrag |
| F-62 | LOW | Gegen den Standard fehlen im Handbuch Angaben zu Autor/Team und Gültigkeitsbereich (§9); „Mit einem Datenbanktreiber arbeiten" ist als Aufgabe ohne Voraussetzung, Vorgehen und Ergebnis angelegt (§5). | `docs/user/benutzerhandbuch-standard.md` §5, §9 | docs/user/benutzerhandbuch.md:3-5, 161-174 | nein | Handbuch-Gerüst weicht vom Standard ab |
| F-63 | LOW | Handbuch gibt Hinweise ohne Spec-Beleg oder mit zu schmaler Lesart: Treiber „feste Namen verwenden oder keine Anweisungen zwischenspeichern" (Spec nennt nur die Grenze), und „Anfrage muss Zeichen für Zeichen der aufgezeichneten entsprechen", obwohl bei Extended auch Parameterwerte (bytegenau), Parametertypen, Format-Codes, Portalnamen und `max_rows` übereinstimmen müssen. | LH-FA-18.a | docs/user/benutzerhandbuch.md:155-158, 169-171, 286-288; spec/spezifikation.md:526-535 | nein | Handbuch-Aussage schmaler/weiter als Spec |
| F-64 | LOW | `releasing.md` §4 beschreibt im Indikativ „Ein Release entsteht durch einen Git-Tag" und zeigt `git tag`/`git push`, während §1 sagt, es gebe keinen Release-Mechanismus; §4 Nr. 2 „Welle, die den Release trägt" nennt nicht, bei welchem Meilenstein (M2 oder M3) v1 getaggt wird; §5 spricht von „veröffentlichten Tags und Artefakten" ohne Veröffentlichungsort. | AGENTS.md §3.7 (beschreibt, was da ist) | docs/maintainer/releasing.md:14-19, 56-72 | nein | Prozessbeschreibung ohne vorhandenen Mechanismus |
| F-65 | LOW | LH-FA-18.a verbirgt Parameterwerte in der Mismatch-Diagnose außer bei `debug`; LH-FA-10.a verlangt die empfangene Query im Klartext, sie kann dieselben Werte als Literale tragen. Die Schutzwirkung gilt nur für Extended. | SPEC-033; LH-FA-10.a | spec/spezifikation.md:543-547, 335-340 | nein | Schutzregel nur für einen Protokollpfad |
| F-66 | INFO | Die Vorlage sagt „Geplante Wellen bekommen noch keine Datei"; die drei Welle-Dateien liegen auf ausdrücklichen Wunsch des Auftraggebers vor, die Roadmap führt sie unter *Offene Wellen* und *Nächste Wellen* nennt „Keine". Keine Abweichung zu beheben; der Zustand „Welle-Datei ohne Slice in `in-progress/`" ist Normalfall (`regelwerk/modul-06-roadmap.md`). | `v6.13.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte | docs/plan/planning/welle-extended-query.md:3-8; docs/plan/planning/in-progress/roadmap.md:19-35 | nein | Auftraggeber-Entscheidung zur Vorlage |
| F-67 | INFO | Das Handbuch beschreibt Verhalten eines Produkts, das es noch nicht gibt (Kopf: „Software-Version: noch nicht veröffentlicht"). Die Standard-Prüfung „alle beschriebenen Funktionen vorhanden" (§13) ist erst mit Code beantwortbar; Zuständigkeit Verifier/Validator. | Handbuch-Standard §13 | docs/user/benutzerhandbuch.md:3-5 | nein | Handbuch vor Produkt |
| F-68 | INFO | Vorbestehend, nicht Teil des Diffs: die Roadmap trägt keinen Marker „Nichts in Arbeit", obwohl `in-progress/` keinen Slice enthält; `welle-walking-skeleton` nennt „Entscheidungen zu den offenen Spec-Punkten … fallen vor der jeweiligen Folge-Welle an", die Spec trifft sie inzwischen. | `v6.13.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte | docs/plan/planning/in-progress/roadmap.md:19-24; docs/plan/planning/welle-walking-skeleton.md:78 | nein | Roadmap-Marker fehlt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| spec/lastenheft.md — Reste „v2" oder „nicht Teil von v1" für Extended Query (Suche `v2`, `Extended`, `nicht Teil`, Glossar, §5, Zielversion) | geprüft, ohne Befund; Out-of-Scope-Verweise, Glossar (:631-633), Globaler Funktionsumfang (:104-106) und §5 sind angepasst |
| spec/lastenheft.md — Hard Rule 3.4/Decken-Regel für die neuen Stellen (LH-FA-18, Szenario 7) | geprüft, ohne Befund; kein Verweis auf Spec, ADR, Slice oder Welle (bestehender Rest: Folgereview F-41) |
| spec/lastenheft.md — LH-FA-18 gegen LH-FA-05, LH-FA-04, Szenario 7, LH-FA-10 | geprüft, ohne Befund; Boundary (Position je Ausführung) und Negative (LH-FA-10) stimmen mit LH-FA-18.a überein |
| spec/spezifikation.md — LH-FA-05.e, SPEC-039, SPEC-036, SPEC-038 gegen LH-FA-18.a | geprüft, ohne Befund; Extended in der Randtabelle „unterstützt", aus SPEC-039 gestrichen, Kompatibilitätstest in Standard- und Simple-Modus |
| spec/spezifikation.md — LH-FA-18.a gegen SPEC-041, SPEC-001/002 | Befunde F-47 bis F-50 |
| spec/spezifikation.md — SPEC-040 Nr. 8 bis 10 gegen LH-FA-14, -16, -17 (SOLL) und LH-FA-18 | geprüft; Nr. 9 und 10 decken die SOLL-Anforderungen und Szenario 7, Befund F-51 betrifft die Vertragsschwelle |
| spec/spezifikation.md — ID-Zählung SPEC-001 bis SPEC-041 | geprüft, ohne Befund (SPEC-041 je einmal definiert, Anker löst auf) |
| spec/architecture.md — Hard Rule 3.4 (Sprach-Artefakte, ADR/Slice/Welle/Commit) im Diff | geprüft, ohne Befund; neue Zeilen nennen Protokollnachrichten, keinen Code, keine Bibliothek |
| spec/architecture.md — Abdeckung des Interaktionstyps | Befund F-46 |
| docs/plan/adr/0012-extended-query-gruppen.md | geprüft, ohne Befund: Status `Proposed`; drei Alternativen (A, B, C) mit Pro/Contra; `Schärft` zeigt auf `LH-FA-18.a` und `SPEC-041`, beide Anker lösen auf (`make docs-check` 0 Befunde); Index-Eintrag ergänzt; keine `Accepted`-ADR überschrieben |
| docs/plan/planning/welle-*.md (drei neue) | geprüft gegen Vorlage: Slice-Tabelle, beobachtbarer Start-Trigger (Vorgänger-Welle `done`), Closure-Trigger mit „Mehr" gegenüber Slice-DoDs, Out-of-Scope, Ruheort-Zeiger vorhanden; Bezug-Links lösen auf. Befunde nur F-52 bis F-54 (Inhalt) |
| docs/plan/planning/open/slice-*.md (zehn neue) — §1 bis §8 | geprüft: ≤ 3 Liefer-Punkte (zwei bis drei), je Slice ein Risiko mit Ausgang „offen bis Closure", Start-Trigger und beide Rückführungen benannt, §8 mit beiden Vorgelagert-Schritten und GF-Modus; keine Suppression, keine ADR-Verletzung. Befunde nur F-52 bis F-54, F-57 bis F-60 |
| docs/plan/planning/in-progress/roadmap.md — Meilensteine M1 bis M3, Graph, „Blockiert"/„Wird blockiert von" in allen vier Welle-Dateien | geprüft, ohne Befund bei Konsistenz: Kette W1 → W2 → W3 → W4 in Graph, Triggern und Abschnitt 5 aller Wellen deckungsgleich; Status-Zellen tragen „offen" ohne Chronik; Drift-Log siehe F-61 |
| docs/plan/planning/welle-walking-skeleton.md (Änderung) | geprüft, ohne Befund; „Blockiert" und Out-of-Scope auf `welle-extended-query` umgestellt |
| docs/user/benutzerhandbuch.md — interne Kennungen (`SPEC-`, `LH-`, `ARC-`, `ADR-`, `slice-`, `welle-`, `MR-`) | geprüft, ohne Befund; nur Anwender-Codes `PGR-…` |
| docs/user/benutzerhandbuch.md — Chronik/Ist-Zustand | geprüft, ohne Befund; Änderungshistorie nennt nur „noch keine veröffentlichte Version" |
| docs/user/benutzerhandbuch.md — Optionen, Umgebungsvariablen, Standardwerte, Log-Level, Exit-Codes 0 bis 6, Meldungscodes `PGR-E…`/`PGR-W…`, `version`/`--help`/`-h`, TLS-Hinweis, Mehr-Session-Recordings, `--force`/`PGR-E2002` gegen Spec | geprüft, ohne Befund; deckt sich mit LH-FA-01.a, -02.a, -03.a/b, -17.a, SPEC-013 bis SPEC-019, SPEC-034 und LH-FA-05.c/e |
| docs/maintainer/releasing.md — genannte Targets/Workflows | geprüft, ohne Befund: einziges Target `make gates` existiert; Workflows werden nicht genannt; §1 sagt ausdrücklich, dass kein Mechanismus existiert; passt zu AGENTS.md §4 (kein nicht vorhandenes Target) |
| AGENTS.md, harness/, .harness/ im Diff | nicht berührt |
| `make docs-check` | 55 Dateien, 0 Befunde (fängt keines der Findings oben) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 10 |
| LOW | 10 |
| INFO | 3 |

(Zusätzlich offen aus dem Folgereview: F-11, F-17, F-30, F-36, F-39, F-41 (Skill-Wortlaut HIGH), F-43; INFO F-24, F-25, F-44, F-45.)

**Finding-Klassen dieses Laufs:** Sicht deckt neuen Interaktionstyp nur im Fließtext · Ereignisreihenfolge zeitabhängig · Format-Beispiel ohne Ereignisvollständigkeit · Fehlerbedingung ohne Definition · Versionierung einer Formatänderung offen · Vertragsschwelle widersprüchlich · Slice-DoD verlangt Zustand einer späteren Welle · Zuständigkeit für dieselbe Anforderung doppelt · Wellen-Ziel ohne tragenden Slice · Handbuch behauptet über die Spec hinaus · Simple-only-Reste nach Scope-Erweiterung · Slice beschreibt bereits geliefertes Ergebnis · Slice-Plan nennt überholten Ist-Stand · Umplanung ohne Drift-Log-Eintrag · Prozessbeschreibung ohne vorhandenen Mechanismus

Hinweis zum Zähler: „Fehlerbedingung ohne Definition" (F-31, F-49) und „Zuständigkeit … doppelt" (F-53, F-57) treten zum zweiten Mal auf; „Sicht trägt Sprach-Artefakte" trat in diesem Lauf nicht auf.

## Verdikt

**Merge-blockierend:** ja — für die MEDIUM-Findings F-46 bis F-54, die vor `Accepted` des Lastenhefts beziehungsweise vor Eröffnung von `welle-extended-query` und `welle-replay-semantik` entschieden sein müssen (Spec-Lücken im Kernformat, Plan-Reihenfolge, Wellen-Ziel ohne Slice). F-55 blockiert die Freigabe des Handbuchs. Kein HIGH; eine Spec-Konsistenz gegen die Verschiebung von Extended Query (Reste „v2"/„nicht Teil von v1") besteht, [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) erfüllt die Form.

**Übergabe:** Findings gehen an den Implementer; F-47, F-48, F-50, F-51 und die offenen Folgereview-Punkte F-30, F-36, F-41, F-43 an den Architect beziehungsweise den Auftraggeber (Vertragsrang); F-52 bis F-54 und F-57 an den Planner (Rückkante Review → Plan). Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
