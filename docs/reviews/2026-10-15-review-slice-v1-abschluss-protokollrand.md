# Review-Report: slice-v1-abschluss-protokollrand — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-protokollrand`, Diff `5c4b3ae..44956ac` (Basis = Move-Commit; Architect-Commit `5101274` liegt darin; Commits des Implementers `9d56324`, `0e90cf0`, `f99192a`, `c7e0f97`, `44956ac`).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-protokollrand` (§1, §3, §4, §6 mit den Entscheidungen des Architect, §7 *Belege des Implementers*, §8)
- aktive ADRs: [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) (Hexagon), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) (`pgproto3`), [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) (Build und Test im Dockerfile), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Lint-Profil)
- berührte `LH-*`-IDs: `LH-FA-05` (`.b` *Auswahl im Replay*, `.e` *Randformen des Protokollrands*), `LH-FA-18`, `LH-FA-12.a`, `LH-FA-23.a`, `SPEC-045`
- `AGENTS.md` (Hard Rules, §3.7, §3.9 bis §3.14)

Eigene Nachprüfung (nicht aus dem Bericht übernommen): Metrik aus §4 nachgemessen; `make abdeckung-check` und `make kopf-check` ausgeführt (beide grün); Mutanten in Kopien des Arbeitsbaums per Bind-Mount im Image `pgr-deps` (`go test ./internal/adapters/driving/pgwire/`, ohne Netz), Arbeitsbaum danach unverändert: Basislauf grün; `pq-praefix-ohne-punkt`, Version 3.2, Version 65535.1, ohne Leser nach TLS, ohne Leser am Kanal rot (Aussagen der Mutationstabelle bestätigt); Reihenfolge Länge vor Code für Version 4.0, 65535.1 und 3.0 grün (F-1). Die E2E-Läufe habe ich nicht wiederholt (Container, Netz, PostgreSQL); die Mutationszeilen zu `test/integration` sind nicht nachgeprüft.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage *Prüfreihenfolge der ersten Nachricht* (erst Länge, dann Code; bei abgelehntem Code wird der Rest nicht gelesen) hat keinen Fall, der sie bricht: Ein Mutant, der die Längenprüfung für 4.0, für 65535.1 oder für 3.0 überspringt, lässt alle Unit-Tests grün (Szenario: erste Nachricht mit Länge 4 oder über `SPEC-045` und Code 4.0 antwortet `PGR-E6002` statt still mit `PGR-W3003` zu enden). Die Zeilen zu Version und Sonderanfrage haben durchweg Länge 8, keine mit Länge über 8 und fehlendem Rest; die Mutationstabelle in §7 führt die Zusage nicht. | `AGENTS.md` §3.10 (Aufzählung und Ortsangabe zählen je Element), `LH-FA-05.e` *Prüfreihenfolge der ersten Nachricht* | `internal/adapters/driving/pgwire/server_protokollrand_test.go` · `{"0.0", 0}, {"2.0", 2 << 16}, {"3.1", 3<<16 \| 1}` (Tabelle der Versionen, nur `rohKopf(8, v.code)`); `internal/adapters/driving/pgwire/server.go` · `case laenge < 8 \|\| laenge > maxStartLaenge:` | ja — `make test` mit dem Mutanten `case (laenge < 8 \|\| laenge > maxStartLaenge) && code != 4<<16:` bleibt grün (ebenso `65535<<16\|1` und `codeProtocol30`); Sonderanfrage-Variante wird vom Bestand (`TestTLSSonderanfrageMitFalscherLaenge`) gefangen | Zusage mit Ortsangabe ohne brechende Mutation |
| F-2 | MEDIUM | Die neue Zeile in `LH-FA-05.e` (*Eine Verbindung, die der Recorder vor dem Aufbau ablehnt oder beendet, verbraucht im `replay` keine Session*, mit Verweis auf `LH-FA-12.a` *Zählung der Verbindungen*) widerspricht dem Wortlaut dort: „Die n-te Verbindung ist die n-te, deren Startnachricht gelesen ist“ und „zugeordnet, sobald deren Startnachricht gelesen ist“. Ein Parameter `_pq_.` ist eine gelesene Startnachricht 3.0, der Code zählt sie nicht (`sessions: 0`); ein Leser der Zählregel erwartet die Session verbraucht und ordnet bei `--session-assignment connection` die falsche Session zu. | `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut`; `AGENTS.md` §3.9 | `spec/spezifikation.md` · `verbraucht im \`replay\`` (Zeile *Modus* in *Randformen des Protokollrands*) gegen `spec/spezifikation.md` · `*Zählung der Verbindungen.* Die n-te Verbindung ist die n-te, deren Startnachricht gelesen` | nein — kein Gate liest die Spezifikation auf Widerspruch; Urteil des Architect | Zählregel der Verbindungen im Wortlaut zweimal verschieden |
| F-3 | LOW | Die Rückverweise in der Spezifikation nennen eine Überschrift, die es nicht gibt: *Auswahl im Replay* verweist auf `LH-FA-05.e`, *Startnachricht*; der Punkt heißt *Startnachricht 3.0*. Die Deklaration von `TestE2ERecordErweiterungOhneUpstream` nennt ihn richtig. | Maintainability (`AGENTS.md` §3.7, Rang-Zeiger) | `spec/spezifikation.md` · `(\`LH-FA-05.e\`, *Startnachricht*).` | nein — `make docs-check` prüft Anker, keine kursiven Überschriftsnamen | Verweis auf nicht vorhandene Überschrift |
| F-4 | LOW | Der Kommentar von `startup` bricht mitten im Satz auf eine Kurzzeile um („…Eine erste Nachricht,“ und dann der Rest erst auf der nächsten Zeile); die Teilersetzung hat den Absatz nicht neu gesetzt. Der Inhalt bleibt zutreffend. | `AGENTS.md` §3.7 | `internal/adapters/driving/pgwire/server.go` · `Eine erste Nachricht,` | nein | Kommentar nach Teilersetzung nicht neu gesetzt |
| F-5 | LOW | Die Spalte „roter Test (nur dieser, außer wo genannt)“ der Mutationstabelle stimmt für die Versionen nicht: Der Mutant für 3.2 färbt auch `TestAndereProtokollversion` (Bestand); dasselbe gilt für 3.1, 4.0, 2.0 und 0, sie nennt die Tabelle nicht. Die Zusage der Tabelle (jeder Mutant färbt die eigene Zeile) trägt, der Wortlaut „nur dieser“ nicht. | `AGENTS.md` §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-protokollrand.md` · `roter Test (nur dieser, außer wo genannt)` | ja — Mutant `\|\| code == 3<<16\|2` in `startkopf`: rot in `TestProtokollrandStartnachricht` und `TestAndereProtokollversion` | Zusage in Belegen weiter als die Messung |
| F-6 | INFO | Grenzen der Tests sind ehrlich benannt: §6 *Grenzen der Tests* nennt Abbruch mitten in einer Nachricht (Gutfall; der Leser kennzeichnet jedes Ende gleich, ReadAtLeast vergleicht per Identität, der Wrapper bleibt erhalten), Benachrichtigung nach der letzten Interaktion (Gutfall), Reset-Zeilen (nur ohne TLS, belegen nur das stille Ende), `NotificationResponse` in einer Extended-Sequenz ohne eigenen Test und die Replay-Auswahl mit TLS nur im Unit-Test; §7 ordnet den grünen E2E-Mutanten `anfrage-NotificationResponse` (als `NoData`) mit gemessener Meldung ein und nennt den Unit-Test, der ihn färbt. Die Deklarationen der Tests sagen nicht mehr zu (`AGENTS.md` §3.11 gewahrt). Nicht in §6 gelistet, nur in §7: der E2E-Mutant der `NoData`-Abbildung. | `AGENTS.md` §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-protokollrand.md` · `Grenzen der Tests (Implementer, 2026-10-15; Messungen in §7)` | nein | — |
| F-7 | INFO | Für die Aufzählung 0 bis 7 Bytes mit `--record-empty-sessions` belegt das gebaute Binary nur 0, 1 und 7 (Unit alle acht, aber ohne die Option); ein Mutant der Option allein ist in der Tabelle nicht aufgeführt. Die Spezifikation sagt „auch mit `--record-empty-sessions` ohne Session“. Zuständig für die Einordnung: Architect/Planner. | `AGENTS.md` §3.10 | `test/integration/protokollrand_e2e_test.go` · `for _, n := range []int{0, 1, 7} {` | ja — Mutant, der im Record mit der Option bei weniger als 8 Bytes eine Session anlegt; E2E `TestE2EProtokollrandTCPProbe` (von mir nicht gefahren) | Aufzählung im E2E nur in Stichprobe |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| (1) Tabelle `LH-FA-05.e` × Modus: `internal/adapters/driving/pgwire/server_protokollrand_test.go` | geprüft; jede Zeile (Start, Versionen, Sonderanfragen, `CancelRequest` Länge 8, Folgen der Anfragen, Client-Nachrichten, Bytes 0 bis 7, gleichzeitig) in beiden Modi und mit/ohne TLS; Mutationen der Stichprobe rot, bis auf F-1 |
| (2) Grenzen und §3.11 | geprüft, ohne Befund (F-6 als Hinweis); keine Zusage weiter als ihr Test |
| (3) Prüfreihenfolge `_pq_.` vor Upstream-Aufbau und vor `PGR-E6003` | geprüft, ohne Befund: `annehmen` prüft erst `_pq_.`, dann TLS; `ErweiterungVorKlartext`, E2E `TestE2ERecordErweiterungOhneUpstream`; Präfix-Mutanten rot. Reihenfolge Länge vor Code und Versionen: F-1 |
| (3) Versionen 0.0, 2.0, 3.1, 3.2, 4.0, 65535.1 → `PGR-E6002` | geprüft, ohne Befund: je ein Fall, Mutanten für 3.2 und 65535.1 rot |
| (4) `internal/adapters/driving/pgwire/leser.go`, `tls.go` | geprüft, ohne Befund: Wrapper zustandslos (keine Race, kein Leck); Unterscheidung `io.EOF` des Stroms und abgelehnter Rumpf trägt auch nach `bufio` und `io.ReadAtLeast` (`streamFehler` bleibt über `errors.As`/`Unwrap` erhalten); beide Konstruktionen des Lesers (Kanal, nach TLS) per Mutant als nötig bestätigt; Anmeldung (`anmeldung.go`) nutzt `verbindungsende` unverändert |
| (5) Plan folgt dem Diff (§3.9), Größe, Zählung §8 | geprüft, ohne Befund: Metrik `git diff --numstat 5c4b3ae..44956ac -- internal test spec docs/user` ergibt 1004 + 22 = 1026 (Spezifikation 109, Produkt-Code 101, Tests 791, Tabellen 25), unter Marke 1400 und Meldeschwelle 1050; §8 zählt 2 Liefer-Punkte und 2 Schichten, der Diff berührt nur PGWire-Adapter (Code und Tests) und Upstream-Adapter (nur Tests) |
| (6) Scope: Handbuch/README, Upstream-Adapter, Hexagon | geprüft, ohne Befund: kein Handbuch und kein README im Diff, keine Änderung an Produkt-Code des Upstream-Adapters, Kern und Ports unberührt, kein `//nolint`, kein `pgproto3` in `test/integration` |
| (7) `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | geprüft, ohne Befund: `make abdeckung-check` Exit 0, `make kopf-check` Exit 0 |
| `spec/spezifikation.md` (Stratum-Reinheit) | geprüft, ohne Befund außer F-2 und F-3: keine Slice-, Wellen-, ADR- oder Bibliotheksnennung in den neuen Absätzen |
| `docs/plan/planning/in-progress/slice-v1-abschluss-protokollrand.md` | geprüft, ohne Befund außer F-5; Risiko-Ausgänge stehen noch offen (Closure-Pflicht, nicht Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage mit Ortsangabe ohne brechende Mutation · Zählregel der Verbindungen im Wortlaut zweimal verschieden · Verweis auf nicht vorhandene Überschrift · Kommentar nach Teilersetzung nicht neu gesetzt · Zusage in Belegen weiter als die Messung · Aufzählung im E2E nur in Stichprobe

## Verdikt

**Merge-blockierend:** ja — zwei MEDIUM (F-1: Mutation für Länge vor Code und „Rest nicht gelesen“ fehlt; F-2: Widerspruch der Spezifikation zur Zählregel in `LH-FA-12.a`). Kein HIGH: Der Code folgt der Spezifikation, und die geprüften Mutanten der Tabelle in §7 sind rot.

**Übergabe:** F-1 und F-3 bis F-5 an den Implementer; F-2 an den Architect (Spezifikation, Technik-Stratum); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
