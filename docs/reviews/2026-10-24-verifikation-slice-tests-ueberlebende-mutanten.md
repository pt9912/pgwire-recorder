# Verifikation slice-tests-ueberlebende-mutanten

**Stand:** main `48228da4`, Basis `83eb9993`. **Rolle:** Verifier, 2026-10-24. Der Review-Bericht wurde vor dem Urteil nicht gelesen.
**Urteil: bestanden.** Keine Befunde der Schwere hoch oder mittel; zwei niedrige, ein Hinweis.

## 1. Mutationen nachgefahren (Modul 11, bewusstes Brechen)

Frische Kopie von `git archive 48228da4` unter eigenem Pfad (Scratchpad, ohne `.git`), je Mutation eine Kopie; Tests per Bind-Mount im gepinnten `golang:1.27.3-alpine@sha256:b8c9…` mit `go test -count=1`, nicht gegen den Build-Kontext. Die unveränderte Kopie ist grün (alle drei Pakete `ok`). „Ohne neuen Test“ heißt `-skip` auf den neuen Test.

| Mutation | mit neuem Test | ohne neuen Test |
|---|---|---|
| M3: Zeile `bestaetigt[m.Type]--` in `extendedNachspielen` entfernt (`replay.go:439`) | rot: `TestReplayExtendedAbgelehntesParseLaesstStehen`, Diagnose `… Anweisung erwartet "SELECT 2", empfangen "SELECT 2"` | `ok` (grün) |
| P8: Zeile `TableOID:` in `spalten` entfernt (`upstream.go:326`) | rot: `TestQuerySpaltenMetadaten` und `TestReceiveSpaltenMetadaten`, `TableOID:0x0` | `ok` |
| P2: Zeile `ColumnNumber:` in `spalten` entfernt (`upstream.go:327`) | rot: beide Tests, `ColumnNumber:0x0` (TableOID 0x4000, 0x4267 erhalten) | `ok` |
| Y2a: `geprueftFromDTO(1, ii, id)` statt `sd.ID` (`yaml.go:719`) | rot: `TestUnmarshalFehlerNenntSession/Nummer` und `/Form` (`Session 1: Interaktion 2 trägt die Nummer 3`; `Session 1, Interaktion 2: endet nicht mit ready_for_query`), dazu der Bestandstest | rot: nur `TestLesenMeldungOhneWertUndZeile` (`sqlite_lesen_inhalt_test.go:436`, `Session 1, Interaktion 1`) |
| Y2n: Konstante `1` statt `sessionID` in „trägt die Nummer“ (`yaml.go:743`) | rot: nur `…/Nummer`, `…/Form` grün | `ok` (grün) |
| Y2f: Konstante `1` statt `sessionID` in `Session %d, Interaktion %d` (`yaml.go:750`) | rot: nur `…/Form` | rot: nur `TestLesenMeldungOhneWertUndZeile` |

Ergebnis: Jede Zeile von §7 stimmt, auch die Abweichung bei Y2 (zwei von drei Mutanten fängt schon ein Bestandstest im `sqlite`-Paket; nur Y2n ist ohne den neuen Test grün). Der neue Test fängt alle drei an der `yaml`-Schnittstelle. M3, P8 und P2 überleben ohne den neuen Test, wie §1 und §7 sagen. Die Tests sind aus dem richtigen Grund rot (Diagnose nennt das falsche SQL; Spalte trägt 0; Meldung nennt Session 1).

## 2. Sensoren (am Zweig von main `48228da4`)

- `make abdeckung-check`: Exit 0.
- `make docs-check`: 1447 Dateien, 0 Befunde (vor Ablage dieses Berichts).
- `make kopf-check`: Exit 0.
- Volles `make gates`: nicht gefahren (Auftrag); §7 belegt Exit 0 am Endstand `32fe4711`, danach änderten nur `798aa4c8` und `48228da4` (Plan, Kommentar, Report-Link); Code-Diff gegen `83eb9993` ist nach `git diff` nur Test, ein Kommentar und Spezifikation.

## 3. DoD

- Liefer-Punkt 1 (M3): Test in `services_test` über `NewReplayService` (`replayMit`) und `ClientMessage`; vergleicht `PGR-E5001` (`CodeReplayMismatch`) und `SELECT 1` / nicht `SELECT 2`; rot ohne die Zeile, sonst grün. Kommentar von `extendedNachspielen` erweitert (siehe F-1).
- Liefer-Punkt 2 (P8, P2): zwei Tests in `postgres_test` über `Upstream.Open`, `Query` bzw. `Send` und `Receive`, Simple- und Extended-Weg, Vergleich der ganzen Spalte (`reflect.DeepEqual`), OID und Nummer ungleich 0 und verschieden, dazu eine Spalte ohne Tabelle. Rot je Mutation, belegt.
- Liefer-Punkt 3 (Y2): `TestUnmarshalFehlerNenntSession` über `recording.Unmarshal`, Fälle Nummer und Form in Session 2, vergleicht `CodeRecordingBroken` und `Session 2`. Rot für beide Fehlerarten bei Session-1-Mutanten.
- Offen und erwartet vor der Closure: DoD-Häkchen, Closure-Notiz (§7 trägt noch Platzhalter), Lerneintrag, Register, Ausgänge der Risiken aus §6 (alle noch „bei Closure“); F-3 aus dem Review liegt laut §7 bei Closure und Nutzer. Das ist kein Befund an der Lieferung, aber Bedingung für `done/`.

## 4. Entscheidungs-Konformität

- Spalten-Metadaten (`LH-FA-02.a`, Extended über `LH-FA-18.a`): beide Wege geprüft. Hinweis H-1.
- Anweisung in der Diagnose (`LH-FA-18.a`): Test hält das benannte Statement; die Grenze beim unbenannten Statement steht in der Spezifikation und wird nicht getestet, wie in §6 zugesagt.
- Ort in der Meldung (`LH-FA-22.a`): beide Fälle (Nummer, Reihenfolge der Gründe: Nummer vor Form) vorhanden; Mutant `si+1` bleibt gleichwertig, nicht Gegenstand.
- Meldungstext in Tests (`SPEC-038`): die drei Tests vergleichen Code und zugesagte Textteile, nicht den Wortlaut der Einordnung.
- `SPEC-049` Punkt 7: Pakete `services_test`, `postgres_test`, `recording_test`; nur exportierte Schnittstelle (`NewReplayService`, `Upstream`, `Unmarshal`), keine Export-Test-Brücke erweitert; Produkt-Code nur im Kommentar geändert (`replay.go`, 3 Zeilen).

## 5. Plan gegen Diff

Der Diff gegen `83eb9993` deckt §3 vollständig: `spec/spezifikation.md`, `replay_extended_test.go`, `replay.go` (nur Kommentar), `upstream_test.go`, `upstream_extended_test.go`, `yaml_test.go`. Dazu `docs/user/abdeckung-unit.md` (fünf Zeilen, in §3 und §8 genannt und entschieden) und Planungsdateien (Plan, Report, ein Zeigernachzug in `slice-lint-bestand-kern-driven.md`). Keine ungenannte Code-Änderung, kein Eingriff in PGWire-Adapter, Bootstrap oder `test/integration`.

Abdeckung: fünf Deklarationen (`LH-FA-02`/Happy, `LH-FA-18`/Happy, `LH-FA-07`/Negative, `LH-FA-10`/Negative, `LH-FA-18`/Negative) entsprechen fünf neuen Zeilen in `abdeckung-unit.md` nach Anforderung, Art, Test und Pfad; `abdeckung-check` grün.

Marke 1500 (`AGENTS.md` §3.14), Befehl mit den Pfaden aus §4 (`internal/hexagon/services`, `…/driven/postgres`, `…/driven/recording`, `spec`, `docs/plan/adr`) und beiden Ausschlüssen, Basis `83eb9993`: **191 Zeilen**, wie in §7. Weit unter der Marke.

Liefer-Punkte 3, Schichten 2 (Kern, Driven-Adapter); die Entscheidung zur erzeugten Tabelle steht in §8 und ist als Einzelfall gekennzeichnet.

## 6. Befunde

| Nr. | Schwere | Befund |
|---|---|---|
| F-1 | niedrig | Der neue Satz im Kommentar von `extendedNachspielen` („ändert ein weiteres parse auf dasselbe Statement in der Interaktion dieses nicht“) sagt „in der Interaktion“ allgemein zu; der Test prüft ein abgelehntes `parse` in einer späteren Gruppe mit Antwort `42P05` bei n = 1. Ein verworfenes `parse` in derselben Gruppe prüft er nicht. Eng genug nach §3.11 für die Zusage von `LH-FA-18.a` (der Mutant bricht genau diesen Fall), aber der Satz reicht ein Stück weiter als der Fall. Kein Mutant bleibt grün; ein „(zweite Gruppe, abgelehnt)“ oder ein zweiter Fall würde es schließen. |
| F-2 | niedrig | Die Spezifikation (`LH-FA-02.a`, Spalten-Metadaten) sagt „zeichnet unverändert auf“ zu; kein Test dieses Slice prüft die Aufzeichnung der zwei Felder (§6 nennt das offen: kein Liefer-Punkt, kein Mutant). Die Aussage steht daher ohne Test in der Spezifikation; Rückweg zum Client liegt bei `slice-tests-ueberlebende-mutanten-driving` (G4). Der Plan hält beides ehrlich fest; nur Beobachtung für die Closure. |
| H-1 | Hinweis | §7 und §3 des Plans sind vor der Closure noch unvollständig (Platzhalter, Häkchen, Risiko-Ausgänge); beides gehört zur Closure, nicht zur Lieferung. |

Kein Befund zu DoD-Belegen: Die in §7 verlangte Mutationstabelle fehlt nicht, und jede Zeile ließ sich nachfahren.
