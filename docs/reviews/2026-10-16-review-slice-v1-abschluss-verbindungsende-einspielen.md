# Review-Report: slice-v1-abschluss-verbindungsende-einspielen — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-einspielen` (Schritt 38a, einfache Anfrage), Diff `33e5eab..69b712e` (Basis `33e5eab`).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-einspielen` (§1, §3, §4, §6, §7 Belege des Implementers, §8) und der Nehmer `slice-v1-abschluss-verbindungsende-einspielen-extended`
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), `SPEC-034` (`PGR-E4003`)
- `LH-FA-20` (`LH-FA-20.a` *Interaktion*, *Ende einer Session*), `LH-FA-25` (`LH-FA-25.a`), `LH-FA-24`
- `AGENTS.md` (Hard Rules, §3.7 bis §3.14), Hexagon-Richtung (`make a-check`, vom Implementer gefahren, nicht neu gefahren)

**Eigene Läufe des Reviewers** (Mutanten in einer Kopie des Baums im Image `pgr-impl-pgwire:deps`, `--network=none`, Original unberührt; ein E2E-Mutant im Arbeitsbaum per `make test-integration`, danach `git checkout`, Baum sauber):

- Service, rot wie in §7 behauptet: `errSessionEnde` als `continue` (nur `TestVergleichFatalFortsetzung` rot), `S` oder `V` für `FATAL`, `ERROR` bei `l.ende` beendet die Session, Aufruf `SchliesseOhneTerminate` nur bei Vergleich oder Ende, Aufruf nach der Verzweigung (entfällt), `l.ende` auch mit `ContinueOnError`, `l.ende` für Extended, `l.ende` nur ohne aufgezeichnete Fehlerantwort.
- Adapter, rot: `Terminate` trotz `SchliesseOhneTerminate`, Schutz `erstes` entfernt, `Close` nur bei `Terminate`.
- E2E, rot: Rückgabe `nil` statt `errSessionEnde` beim erwarteten Ende (`TestE2EPlayEndeErwartet` FAIL).
- Adapter, **grün**: `close(s.ende)` nur bei `terminate` (F-1).
- Messung `665` nachgerechnet: Einfügungen 572, Löschungen 93; Produktcode 81, Tests und Testdaten 558, Tabellen 26; Faktor 6,9; E2E-Fall 65 (37 und 28). Stimmt.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Port sagt zu, `SchliesseOhneTerminate` beende den Sender; der Mutant, der `close(s.ende)` nur bei `terminate` ausführt, bleibt grün (Sender-Goroutine läuft nach dem Ende weiter), weil `TestEinspielGruppeSchliesseBeendetSender`, `…VerwirftWarteschlange` und `…SchliesseBeimSenden` nur `Schliesse` rufen. Der Port nennt „Anfrage sendet oder Gruppe sendet oder Naechste wartet“; getestet ist für `SchliesseOhneTerminate` nur der Fall „Anfrage sendet“. Zuständig: Implementer. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `internal/adapters/driven/postgres/einspielen.go` · `if erstes {` mit `close(s.ende)` in `schliesse`; `internal/hexagon/ports/driven/einspielziel.go` · „Schliesse und SchliesseOhneTerminate dürfen jederzeit laufen“ | ja — `make test` mit dem Mutanten bleibt grün | Aufzählung in Port-Zusage ohne Fall je Element |
| F-2 | LOW | Die Abdeckungs-Deklaration von `TestEinspielSchliesseSchreibtNichtBeimSenden` sagt „SchliesseOhneTerminate (Fehlerantwort FATAL während eine Gruppe sendet)“; der Test sendet eine Anfrage, keine Gruppe. Der Satz reicht weiter als die Prüfung. Zuständig: Implementer. | `AGENTS.md` §3.11 | `internal/adapters/driven/postgres/einspielen_test.go` · „SchliesseOhneTerminate (Fehlerantwort FATAL während eine Gruppe sendet)“; gleich in `docs/user/abdeckung-unit.md` | ja — `make abdeckung-check` prüft nur die Deckung der Zeilen, nicht den Satz | Kommentar sagt mehr zu, als der Test prüft |
| F-3 | LOW | Der Nehmer `slice-v1-abschluss-verbindungsende-einspielen-extended` §6 nennt den Fall „Sendefehler nach der gelesenen `FATAL`“ mit „Test im Adapter beim Geber vorhanden“; für `SchliesseOhneTerminate` mit Gruppen in der Warteschlange gibt es beim Geber keinen (F-1). Die Sendung trägt eine Bindung, die nicht stimmt. Zuständig: Planner (Nehmer), nach Behebung von F-1 Implementer. | `AGENTS.md` §3.13 (auch eine Bindung ist eine Sendung) | `docs/plan/planning/next/slice-v1-abschluss-verbindungsende-einspielen-extended.md` · „Test im Adapter beim Geber vorhanden“ | nein — Urteil | Bindung an Geber nennt nicht vorhandenen Test |
| F-4 | LOW | Der Test-Kommentar und die Abdeckungs-Zeile zu `TestPlayExtendedFatal` sagen „noch nicht erwartet … Zwischenstand bis zur Einstufung der Extended-Interaktion“ ohne auflösbare Kennung des Folge-Slice; der Satz ist eine Chronik-Wendung und nennt nicht, wer ihn umkehrt (`slice-v1-abschluss-verbindungsende-einspielen-extended`). Der Fall selbst ist getestet (Mutant `l.ende` für Extended rot). Zuständig: Implementer. | `AGENTS.md` §3.7 (Herkunft in ein auflösbares Feld) | `internal/hexagon/services/play_extended_test.go` · „Zwischenstand\n// bis zur Einstufung der Extended-Interaktion“; `docs/user/abdeckung-unit.md` · „Zwischenstand bis zur Einstufung“ | nein — Urteil | Zeitwort statt Kennung im Kommentar |
| F-5 | LOW | §3 des Plans nennt `play_laufsteuerung_test.go` und `spec/spezifikation.md` nicht; der Diff ändert beide (Kommentar von `TestPlayErwarteterFehler`, Spezifikation durch den Architect vor dem Code). §7 Beobachtung (3) nennt die Datei erst im Nachhinein. Zuständig: Implementer. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen.md` · „`TestRunPlayFehlerantwort` im Bootstrap bleibt unverändert“ (Zeile der Tests in §3) | nein — Urteil | Plan folgt dem Diff nicht |
| F-6 | INFO | Mutant (1) von §7 (`errSessionEnde` als `continue`): die Begründung trägt nur, solange die Aufzeichnung gültig ist (Ende letztes Ereignis der Session, `Session.ValidateEnd` beim Laden); der Service prüft es nicht selbst. Der Mutant ist für gültige Aufzeichnungen äquivalent, ein Test dafür wäre ohne Aussage. Kein Handlungsbedarf. Zuständig: —. | Maintainability | `internal/hexagon/services/play.go` · `if errors.Is(err, errSessionEnde) {` | nein | — |
| F-7 | INFO | Der „gestellte“ Mutant („Ende gilt nur an der aufgezeichneten Interaktion, früher“) ist im Aufbau nicht herstellbar: `interaktion` sieht nur die eigene Interaktion (`endetEinfach(in)`). Die Begründung trägt; die Fälle „FATAL in früherer Interaktion“ und „ohne Ende an der Stelle“ färben die erreichbaren Varianten (Mutant `l.ende = AllowRecordedErrors` rot). Zuständig: —. | Maintainability | `internal/hexagon/services/play.go` · `l.ende = s.optionen.AllowRecordedErrors && endetEinfach(in)` | nein | — |
| F-8 | INFO | Nicht in §6 genannt: Gewinnt beim ersten Signal das Schließen aus dem Abbruch (`us.Schliesse()` in `context.AfterFunc`) gegen `SchliesseOhneTerminate` aus `fehlerantwort`, geht nach einer gelesenen `FATAL` ein `Terminate` auf die Leitung. §6 nennt nur die umgekehrte Richtung des Rennens als folgenlos. Folgenlos für den Server (Verbindung ist beendet), aber die Spezifikation sagt „außer nach einem Ende des Servers“ ohne Vorbehalt. Zuständig: Architect. | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen.md` · „ein Rennen mit dem Schließen nach der Fehlerantwort ist folgenlos“ | nein | Rennen nur in einer Richtung entschieden |
| F-9 | INFO | Die Handbuch-Lücke (`docs/user/benutzerhandbuch.md`: „Ein Fehler mit dem Schweregrad `FATAL` beendet die Verbindung und bricht das Einspielen auch dann ab“ stimmt für aufgezeichnetes Ende mit `--allow-recorded-errors` nicht mehr) ist in §1 als akzeptiert geführt (Entscheidung des Nutzers vom 2026-10-12, Nehmer `slice-v1-abschluss-verbindungsende-doku`). Der Hilfetext der Option sagt über das Ende nichts zu. Zuständig: —. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „und bricht das Einspielen auch dann ab (`PGR-E4003`)“ | nein | Akzeptierte Lücke bis Doku-Slice |
| F-10 | INFO | Zwei Kommentare in `play.go` laufen über die übliche Umbruchbreite (`session`: „… kehrt erst nach dem Schließen zurück. Mit KeepTiming wartet sie …“; `interaktion`: „… Jede andere ist mit AllowRecordedErrors erwartet, wenn die“). Kein Konventions-Anker (`gofmt` und `make lint` rügen es nicht). Zuständig: —. | Maintainability | `internal/hexagon/services/play.go` · „schließt sie sofort und kehrt erst nach dem Schließen zurück. Mit KeepTiming wartet sie vor jeder Interaktion (warte).“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/play.go` (Einstufung, `lauf.ende`, `fehlerantwort`, `session`) | geprüft, ohne Befund über F-6 bis F-8 hinaus; Mutanten zu den Konjunktions-Gliedern (Option, Ende an der Stelle, `ContinueOnError`), V/S bei `FATAL` und `PANIC`, Aufruf mit und ohne Vergleich, `ERROR` bei Ende, Verbindungsverlust, Zähler (laut §7 rot, Stichprobe selbst rot gesehen) |
| `internal/hexagon/ports/driven/einspielziel.go` | geprüft, Befund F-1 (Zusage `SchliesseOhneTerminate` beendet den Sender) |
| `internal/adapters/driven/postgres/einspielen.go` | geprüft, Befund F-1; kein `Terminate` auf der Leitung, zweites Schließen, `Close` beim Schließen ohne `Terminate`: Mutanten rot |
| Tests Service (`play_ende_test.go`, `play_test.go`, `play_extended_test.go`, `play_vergleich_test.go`, `play_laufsteuerung_test.go`) | geprüft, Befund F-4 (Kommentar); Fixture je Fall in `endeRec` selbst gebaut, Fälle isoliert; Extended bleibt `PGR-E4003` getestet |
| Tests Adapter (`einspielen_ende_test.go`, `einspielen_test.go`) | geprüft, Befunde F-1 und F-2 |
| `test/integration` (`play_ende_e2e_test.go`, `testdata/play-ende-erwartet.yaml`) | geprüft, ohne Befund; Aufzeichnung gehört genau einem Test, Mutant am Binary rot; Grenze „nur Referenzversion 17“ in §6 geführt |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | geprüft, Befund F-2, F-4; die übrigen Zeilen sagen nicht mehr zu, als die Tests prüfen |
| `spec/spezifikation.md` (Architect-Commit vor dem Code) | geprüft, ohne Befund; Drift-Log-Zeile vorhanden, Stratum nennt keine ADR, keinen Slice |
| Veraltete Wörter in `internal/` und `cmd/` („erwartet wäre“, „sendet Terminate“, Hilfetext von `--allow-recorded-errors`) | geprüft, ohne Befund; der Wortlaut in `play_laufsteuerung_test.go` und `abdeckung-unit.md` ist enger gefasst |
| Hexagon-Richtung (Core-Reinheit, Imports der Diff-Dateien) | geprüft, ohne Befund; Diff berührt nur Kern (Service, Port) und Einspiel-Adapter, Bootstrap und CLI unberührt |
| Messung und Metrik (§4, §7) | geprüft, ohne Befund; 665 und Faktor 6,9 nachgerechnet |
| Nehmer `slice-v1-abschluss-verbindungsende-einspielen-extended` (§1, DoD, §8) | geprüft, Befund F-3; Sendung und Nachzählung mit Kennung des Gebers vorhanden |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft, ohne Befund; Zustandsfeld nennt Zustand ohne Chronik |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Aufzählung in Port-Zusage ohne Fall je Element · Kommentar sagt mehr zu, als der Test prüft · Bindung an Geber nennt nicht vorhandenen Test · Zeitwort statt Kennung im Kommentar · Plan folgt dem Diff nicht

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): grüner Mutant bei einer Zusage des neuen Port-Vertrags. LOW und INFO blockieren nicht.

**Übergabe:** Findings gehen an den Implementer (F-3 zusätzlich an den Planner, F-8 an den Architect); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
