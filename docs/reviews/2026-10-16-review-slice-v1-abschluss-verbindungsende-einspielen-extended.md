# Review-Report: slice-v1-abschluss-verbindungsende-einspielen-extended — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-einspielen-extended` (Schritt 38b, Extended-Interaktion), Diff `08290f8..a927048` (Basis `08290f8` laut Plan §4; der Plan-Diff ist gegen den Stand nach `e92825a`, der Prüfung des Architect, gelesen).

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-einspielen-extended` (§1, §3, §4, §6, §7 Belege des Implementers, §8) und der Nehmer `slice-v1-abschluss-verbindungsende-vergleich`
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), `SPEC-034` (`PGR-E4003`)
- `LH-FA-20` (`LH-FA-20.a` *Interaktion*, *Gruppen*, *Abbruchsignal*; Spezifikationsstände aus `526a50f`, `5c8c592`, `e92825a`), `LH-FA-25` (`LH-FA-25.a`), `LH-FA-18`
- `AGENTS.md` (Hard Rules, §3.7 bis §3.14), Hexagon-Richtung (`make a-check`, vom Implementer gefahren, nicht neu gefahren)

**Eigene Läufe des Reviewers** (Unit-Mutanten in Kopien von `go.mod`, `go.sum`, `internal`, `cmd` im Image der Stufe `deps`, `--network=none`, Original unberührt; ein E2E-Mutant im Arbeitsbaum per `make test-integration`, danach `git checkout`, Baum sauber):

- Service, rot wie in §7 behauptet: Mutant 1 (Flush-Ende nicht erkannt: nur `…Erwartet`), 2 (nur Flush erkannt: `…Erwartet`, `NachError`, `TestPlayExtendedFatal`), 3 (erste statt letzte Gruppe: `…Erwartet`, `NachError`, `TestPlayExtendedFatal`), 4 (grob wie in §7: `…Erwartet`, `NachError`, `TestPlayExtendedFatal`), 6 (`Schliesse` statt `SchliesseOhneTerminate`: acht Tests, darunter alle Extended-Tests), 15 (Rückkehr vor `bisBereit` entfällt: `TestPlayEndeExtendedLetzteFlush`, fünf Fälle), 16 (`n = 0` entfällt: sechs Tests), 17-Variante (`fehlerGesehen` beim Lesen ignoriert: fünf Tests), `Eingespielt++` entfällt (zwanzig Tests).
- Weitere Mutanten, rot: Rückkehr nur ohne `fehlerGesehen`, Rückkehr nur mit `l.ende` (`LetzteFlush`: die Fälle ohne Option), Rückkehr nach der ersten statt der letzten Gruppe, `Sync` zählt als `Flush` im Helfer.
- Mutant 4 in der feineren Form (nur `FATAL` oder `PANIC` in einer nicht letzten Gruppe setzt `l.ende` zurück): nur `…Erwartet` und `TestPlayExtendedFatal` rot, `NachError` grün.
- Service, **grün** (kein Mutant der Tabelle in §7): `n > 0` in `traegtEnde` entfernt (A); „irgendeine Server-Nachricht der letzten Gruppe ist `connection_closed`“ statt „die letzte“ (B); Rückkehr vor `bisBereit` nur ohne Vergleich (`l.v == nil`, O).
- E2E: Mutant 1 (Flush-Ende nicht erkannt) färbt `TestE2EPlayEndeExtendedFlush` (Exit-Code 4), `TestE2EPlayEndeExtendedSync` bleibt grün, wie in §7 (E2).
- Messung `447` nachgerechnet: Einfügungen 431, Löschungen 16 (`git diff --numstat 08290f8..5cab7bf -- internal test docs/user`), dieselbe Summe bis `a927048`; Produktcode 34, Tests im Service 256, E2E mit Fixtures 141, Tabellen 16; Faktor 311 / 34 = 9,1, mit Fixtures 397 / 34 = 11,7; E2E-Fall 141 / 2 = 70,5. Stimmt.
- Port, Adapter, Model, CLI, Bootstrap: `git diff --stat 08290f8..a927048` über `internal/adapters`, `internal/hexagon/ports`, `internal/hexagon/model`, `cmd`, `internal/bootstrap` leer; unverändert wie geplant.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Bindung (g) an `slice-v1-abschluss-verbindungsende-vergleich` steht dort nur in §6 (Randformen, Zeile *Extended*) und §8, nicht in §1 und nicht in der DoD; `grep` nach `Flush` findet im Nehmer nur diese eine §6-Zeile. §8 des Nehmers nennt „die Zeile *Extended* von DoD Punkt 1“, die es in der DoD nicht gibt (ein Punkt, Klammer „Test, einfache Anfrage und Extended-Interaktion“). Der Mutant „Rückkehr vor `bisBereit` nur ohne Vergleich“ (O) bleibt in diesem Slice grün; der Fall im Pfad mit Vergleich ist damit allein auf eine Adresse gelegt, deren DoD ihn nicht führt. Failure-Szenario: Der Implementer des Nehmers liest §1 und DoD, liefert keinen Fall für die letzte `Flush`-Gruppe mit Vergleich, und die geänderte Funktion `extended` bleibt dort unbelegt. | `AGENTS.md` §3.13 (Bindung ist Sendung: §1 des Nehmers mit Kennung des Gebers, §6; Nachzählen in §8); Reviewer-Skill, MEDIUM *Adresse nimmt nicht an* | `docs/plan/planning/next/slice-v1-abschluss-verbindungsende-vergleich.md` · „**Bindung an `slice-v1-abschluss-verbindungsende-einspielen-extended` (g)**“ und „ein Fall im Service in der Zeile *Extended* von DoD Punkt 1“ | ja — `grep -n "Flush" <Nehmer>` trifft nur §6; ein Mutant O im Service bleibt grün | Bindung nur in §6 und §8 des Nehmers, nicht in §1 oder DoD |
| F-2 | LOW | Der neue Satz zum Abbruchsignal in `LH-FA-20.a` (*Abbruchsignal*: „die mit `Flush` endende letzte Gruppe bis zur Antwort auf jede ihrer Client-Nachrichten“) hat keinen Test: `TestPlayExtendedSignal` hat eine letzte Gruppe mit `Sync`, und §6 nennt das Signal während der letzten `Flush`-Gruppe nicht als Randform. | `AGENTS.md` §3.10 (Aufzählung und Ortsangabe zählen je Element) und §3.12 | `spec/spezifikation.md` · „die mit `Flush` endende letzte Gruppe bis zur“; `internal/hexagon/services/play_extended_test.go` · `func TestPlayExtendedSignal` | ja — ein Test mit Signal in der letzten `Flush`-Gruppe, Mutant: Warten auf `ReadyForQuery` nach dem Signal | Zusage der Spezifikation ohne Test und ohne Randform in §6 |
| F-3 | LOW | Der Kommentar von `extended` sagt, nach einer Fehlerantwort, die nicht abbricht, lese `play` „dann bis zum ReadyForQuery“, und die letzte `Flush`-Gruppe ende „mit den Antworten dieser Gruppe oder mit der Fehlerantwort“. Nach einem `ERROR` in einer früheren Gruppe liest `play` bei einer letzten `Flush`-Gruppe nichts mehr (Test `früherer ERROR`, Spezifikation *Gruppen*); beide Sätze nennen diesen Fall nicht und der erste widerspricht ihm. | `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist) | `internal/hexagon/services/play.go` · „liest dann bis zum ReadyForQuery. Endet die letzte Gruppe mit Flush (nur bei“ | nein — Lesen des Kommentars gegen `TestPlayEndeExtendedLetzteFlush/früherer ERROR` | Kommentar nennt den Sonderfall der Rückkehr ohne Lesen nicht |
| F-4 | INFO | Zu (2): Mutant 4 färbt `TestPlayEndeExtendedNachError` als Artefakt, nicht aus fehlender Trennschärfe. Die feinere Form (nur `FATAL`/`PANIC` in einer nicht letzten Gruppe) färbt `…Erwartet` und `TestPlayExtendedFatal` und lässt `NachError` grün; die Zusage (d) ist mit den acht Fällen „frühere Gruppe“ getragen, die Einordnung in §7 Beobachtung (2) stimmt. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen-extended.md` · „Mutant 4 färbt `NachError` mit“ | ja — Mutant in Kopie | — |
| F-5 | INFO | Zu (1): Die Einordnung der grünen Mutanten 15 und 6 am Binary trägt. Mutant 15 färbt im Service `normale Antwort mit Option`, `normale Antwort ohne Option`, `ERROR mit Option` und beide `früherer ERROR`-Fälle über den Zählstand von `naechste`; Mutant 6 färbt acht Tests über `war()`. Am Binary ist 15 äquivalent (die `FATAL` ist die erste Antwort der Gruppe), 6 nicht sichtbar. Das akzeptierte Negativ (normale Antwort auf die letzte `Flush`-Gruppe nicht am Binary) trägt mit Grund: Die Zählregel gegen PostgreSQL ist für `Execute` mit `Flush` in `TestE2EPlayExtendedFlushWarten` belegt; das Nicht-Warten belegt nur der Zählstand des Fakes, nicht ein Hängen, das der Fake nicht zeigt (§6 sagt das). | Maintainability | `internal/hexagon/services/play_ende_extended_test.go` · `func TestPlayEndeExtendedLetzteFlush` | ja — Mutanten 15 und 6 in Kopie | — |
| F-6 | INFO | Zwei Mutanten der Konjunktions-Glieder in `traegtEnde` bleiben grün: `n > 0` entfällt (A) und „irgendeine statt der letzten Server-Nachricht ist `connection_closed`“ (B). Beide sind für geladene Aufzeichnungen äquivalent (`Group.validate` verlangt eine nicht leere letzte Server-Folge, `pruefeEnde` lässt `connection_closed` nur zuletzt); der neue Zugriff `gruppen[len(gruppen)-1]` in `extended` hat für eine Extended-Interaktion ohne Gruppe keinen Wächter, auch dort trägt das Format. Kein Test belegt die Wächter. | Maintainability | `internal/hexagon/services/play.go` · „return n > 0 && server[n-1].Type == model.ResponseConnectionClosed“ | ja — Mutanten A und B in Kopie | — |
| F-7 | INFO | Zu (6): Die Frage, ob PostgreSQL vor der `FATAL` von `pg_terminate_backend(pg_backend_pid())` in `Execute` noch `DataRow` oder `CommandComplete` sendet, ist beim Schreiben der Fixtures aufgetreten und steht in §7, nicht in §6. Der Fall `Flush` am Binary wäre bei vorangehendem `CommandComplete` ebenfalls grün (Exit-Code 0, `3/0/0`, nächste Session): Er belegt die Einstufung nur, wenn die `FATAL` die erste Antwort ist; die Zählung, die dann die `FATAL` ungelesen lässt, ist als Grenze in `LH-FA-20.a` *Gruppen* gedeckt. | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen-extended.md` · „war offen, ob PostgreSQL vor der FATAL“ | nein — Verhalten des Servers, nur gegen Referenz 17 | — |
| F-8 | INFO | Zu (8): Das Hängen bei einer kleineren Antwortzahl des Servers als `erwarteteAntworten` (Rückkehr ohne `bisBereit`, kein Timeout) ist in §6 benannt (Absatz *Letzte Gruppe mit `Flush`, Lesen und Hängen*), steht aber nicht unter *Risiken* mit Ausgang und nicht als Grenze in der Spezifikation für den Anwender. Es bestand für Gruppen mit `Flush` davor schon. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen-extended.md` · „Ein Hängen gäbe es nur, wenn der Server weniger Antworten sendet“ | nein | — |
| F-9 | INFO | In `LH-FA-20.a` *Gruppen* steht „Nach einer `ErrorResponse` wartet `play` in dieser Interaktion nur noch auf das `ReadyForQuery` nach ihrem `Sync`“ und mehrere Sätze später, dass `play` bei einer letzten `Flush`-Gruppe nach einem früheren `ERROR` nichts mehr liest; der erste Satz nennt die Ausnahme nicht. Zuständig: Architect. | Maintainability | `spec/spezifikation.md` · „wartet `play` in dieser Interaktion nur noch auf das“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/play.go` (`interaktion`, `extended`, `endetMitFlush`, `traegtEnde`) | geprüft; Befunde F-3, F-6 |
| `internal/hexagon/services/play_ende_extended_test.go`, `play_extended_test.go` | geprüft, Mutanten der Tabelle in §7 stichprobenartig rot gesehen; Befunde F-4, F-5 |
| `test/integration` (`play_ende_extended_e2e_test.go`, zwei Fixtures) | geprüft, je Fall eine Datei und eine Datenbank (`play_ende_extended_sync`, `play_ende_extended_flush`), Mutant 1 rot nur im `Flush`-Fall; Befund F-7 |
| `internal/adapters`, `internal/hexagon/ports`, `internal/hexagon/model`, `cmd`, `internal/bootstrap` | geprüft, unverändert, ohne Befund |
| Hexagon-Richtung und Core-Reinheit (`internal/hexagon/services`) | geprüft, keine neuen Importe, ohne Befund |
| Veraltete Wörter in `internal/` und `test/` („außer bei einer einfachen Anfrage“, „Zwischenstand“, „ist hier nicht erwartet“, `kehrt den Fall um`) | geprüft, kein Treffer; `endetEinfach` bleibt als Helfer von `traegtEnde` und `replay.go` |
| `spec/spezifikation.md` (Schritt 3, *Gruppen*, *Abbruchsignal*, Drift-Log) | geprüft; Befunde F-2, F-9; keine ADR-, Slice- oder Wellen-Kennung im Stratum |
| `docs/user/abdeckung-unit.md`, `abdeckung-e2e.md` | geprüft, Deklarationen decken, was die Tests prüfen (Fälle mit `--continue-on-error` in `…NichtErwartet` nur in zwei Zeilen, die Aussage hängt am Code ohne Bedingung auf die Option) |
| Slice-Plan §1, §3, §4, §6, §7, §8 (Plan folgt Diff, `AGENTS.md` §3.9) | geprüft; §3 und §1 folgen dem Diff, Messung 447 stimmt; Befund F-1 zum Nehmer |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Bindung nur in §6 und §8 des Nehmers, nicht in §1 oder DoD (`BEO-REPO/folge-slice-adresse-nimmt-nicht-an`) · Zusage der Spezifikation ohne Test und ohne Randform in §6 · Kommentar nennt den Sonderfall nicht

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM). Der Befund liegt im Plan des Nehmers `slice-v1-abschluss-verbindungsende-vergleich`, nicht im Code dieses Slice; Produktcode, Tests und Fixtures tragen ihre Zusagen, alle stichprobenartig gefahrenen Mutanten der Tabelle sind rot, die Messung stimmt, Port und Adapter sind unverändert. F-2 und F-3 sind LOW.

**Übergabe:** F-1 an den Architect (Eintrag in §1 oder DoD des Nehmers mit Kennung des Gebers, Berichtigung des Verweises in §8, Nachzählen); F-2 an Implementer oder Architect (Test oder engere Fassung des Satzes); F-3 an den Implementer; F-9 an den Architect. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg**; DoD- und Spec-Konformität prüft der Verifier separat.
