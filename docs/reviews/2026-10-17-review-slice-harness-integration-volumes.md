# Review-Report: slice-harness-integration-volumes — 2026-10-17

**Review-Art:** Code — gegen Plan, `SPEC-057`, `SPEC-055` Punkt 9, `SPEC-052` Punkt 13 und die Hard Rules (`AGENTS.md` §3.7 bis §3.14).

**Gegenstand:** `slice-harness-integration-volumes`, Diff `de908591..c9400ec4`

**Skill:** `.harness/skills/reviewer.md` @ c9400ec4
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

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

- Slice-Plan `slice-harness-integration-volumes` (§1 bis §8, §7 Belege des Implementers)
- `SPEC-057` (neu), `SPEC-055` Punkt 9, `SPEC-052` Punkte 4 und 13
- [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) (Integrationstests im Multistage-`Dockerfile`)
- keine `LH-*`-Kennung (Plan: kein Lastenheft-Bezug)
- `AGENTS.md` (Hard Rules §3.7 bis §3.14)

**Selbst gefahren** (Kopien im Scratchpad, nie im Arbeitsbaum; Docker real, nur Objekte mit Präfix `pgr-ivg-`):
`make integration-volumes-gegenprobe`-Skript unverändert grün (59 s, danach keine Container und Volumes des Präfixes);
Mutationen im Code der Kopie, die die Gegenprobe rot sehen musste: `rm -fv`→`rm -f` im Funktionsfile (rot: `p2-grundfall`, `p2-ohne-v`),
Prüfung nach Namen entfernt (rot: `p3-angelegt`, `p4-rest-rot`, `p5-sigint`, `p5-sigterm`), Exit-Code der Prüfung verworfen (rot: `p3-angelegt`, `p4-rest-rot`),
Namen nach dem Entfernen lesen (wie Prüfung entfernt), Meldung bei fehlendem Container (rot: `p3-kein-container`), nur laufende Container lesen (rot: `p3-angelegt`, `p4-rest-rot`),
`local code=0` im Runner (rot: `p5-*`); `run-container.sh` mit `$pg` ohne `-v` im Aufräumen bzw. vor dem Wiedergeben: `container-gegenprobe.sh` rot (`p9-*-volume`, `p7-pg-vor-replay`);
Runner-Ende `exit "$ende"`→`exit "$code"`: **grün** (F-1). `make test-integration`: grün, Volumes nach Namen vor und nach dem Lauf: keines neu (zwei Einträge verschwanden, `pgr-ver-…` eines fremden Laufs).
`make docs-check` 0 Befunde; `make gates-schnell-gegenprobe` grün.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „Rest macht den Lauf rot, Exit 1, auch bei grünen Tests“ hat auf Runner-Ebene keine Mutation: `exit "$ende"` durch `exit "$code"` ersetzt in einer Kopie, die Gegenprobe bleibt grün. `p4-rest-rot` prüft nur `integration_ende`, die `p5-*`-Fälle enden ohnehin mit 143/130, ein grüner Testlauf mit Rest ist nicht abgedeckt und nicht als akzeptiertes Negativ genannt (Punkt 5 nennt nur roten Lauf und Fehler des Aufbaus). Scenario: der Runner verwirft den Rest-Exit, ein Leck bleibt im Gate grün. | `SPEC-057` Punkt 4; `AGENTS.md` §3.10 | `tools/test/run-integration-tests.sh` · `exit "$ende"` | ja — Mutation in Kopie, `make integration-volumes-gegenprobe` | Zusage ohne Mutation an der Verdrahtung des Neuen Vertrags |
| F-2 | MEDIUM | Der Vertrag sagt zu, die Gegenprobe verlange „genau diese Menge“ roter Fälle (mehr machen sie rot); `nur` in `mut` beschneidet die Fälle, sodass sieben Mutanten tatsächlich mehr Fälle färben als verlangt. Gemessen mit `nur` aufgehoben: `p2-ohne-tests` zusätzlich `p3-angelegt`; `p2-ohne-pg` zusätzlich `p3-angelegt`, `p5-sigint`, `p5-sigterm`; `p1-nur-tests`, `p4-nicht-entfernen`, `p4-ohne-namen` zusätzlich `p5-sigint`, `p5-sigterm`; `p4-alles-entfernen` (Zeile `p4-fremd` der Tabelle) zusätzlich `p5-*`; `p4-immer-1` (Zeile `p4-vorrang`) zusätzlich `p5-*`. Die Tabelle in Punkt 7 und §7 des Plans nennen für diese nur die erste Menge. | `AGENTS.md` §3.11; `SPEC-057` Punkt 7 | `harness/vertraege/SPEC-057-integration-aufraeumen.md` · „die Gegenprobe verlangt genau diese Menge“ | ja — `nur=""` in `mut`, `GP_ENTDECKE=1` | Vertragsaussage über Mengen roter Fälle weiter als die Prüfung |
| F-3 | MEDIUM | Aussagen über die Mutanten `p5-*` sind falsch. `harness/README.md` nennt unter den Mutationen „ohne Falle für `SIGTERM` oder `SIGINT`“, die rot wird; beide sind laut `SPEC-057` Punkt 7 und Gegenprobe äquivalent (kein roter Fall), rot werden nur die Mutanten „Falle beendet ohne Aufräumen“. Im Kopf der Gegenprobe steht zusätzlich „SPEC-057 Punkt 7 nennt ihn als roten Mutanten“, was der Vertrag nicht mehr sagt. | `AGENTS.md` §3.11 | `harness/README.md` · „ohne Falle für `SIGTERM` oder `SIGINT`“; `tools/test/integration-volumes-gegenprobe.sh` · „nennt ihn als roten Mutanten“ | nein — kein Gate prüft Prosa | Zusage im Kommentar/Gate-Index weiter als Prüfung |
| F-4 | LOW | Die Zusage „Meldung nennt Namen, kein Inhalt des Volumes“ hat keine Mutation für den zweiten Teil: `p3-angelegt` prüft nur, dass die Namen in der Meldung stehen, nicht, dass nichts anderes erscheint. | `SPEC-057` Punkt 4; `AGENTS.md` §3.10 | `tools/test/integration-volumes-gegenprobe.sh` · `in_err "${probe[0]}"` | ja — Mutation, die weitere Ausgabe ergänzt | Zusage ohne Mutation an der Verdrahtung des Neuen Vertrags |
| F-5 | LOW | Plan folgt dem Diff nicht an zwei Stellen: §3 hat keine Zeile für `harness/vertraege/SPEC-055-container-und-binaries.md`, obwohl der Diff sie ändert (Absatz zu `$rec`, `$rep`, `$img`, `version`, `psql`; Kopf und §3 sagen „nicht im Diff“); die Liste der weiteren `docker run` in §1 *Bestand* nennt `d-check.mk` und `a-check.mk` im Wurzelverzeichnis nicht (nur §7 holt es nach, §6 Prüfung des Architect Nr. 5 beschließt „§3 bleibt“). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-integration-volumes.md` · „Weitere `docker run` in `tools/` und `harness/mk/`“ | nein — Urteil von Review | Plan folgt der Korrektur nicht |
| F-6 | INFO | Der Befund „Falle für `INT` schützt ein Vordergrundkommando (`docker exec`, Exit 2)“ steht im Vertrag als Zusage-Begründung, ohne dass ein Fall ihn trägt; der Vertrag nennt es ausdrücklich als akzeptiertes Negativ, §7 des Plans hat die Messung. Konsistent. Zweites SIGTERM/SIGINT, Docker-Version, äquivalente Mutanten `p1-ohne-typ`/`p5-term-entfaellt`/`p5-int-entfaellt`: in Vertrag, Plan und Gegenprobe gleich und gemessen wahr (die Gegenprobe führt sie ohne roten Fall). Verantwortlich: Verifier für DoD. | `SPEC-057` Punkte 5, 7, 8 | `harness/vertraege/SPEC-057-integration-aufraeumen.md` · „akzeptiertes Negativ“ | ja — `make integration-volumes-gegenprobe` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/test/integration-aufraeumen.sh` (Funktionen) | geprüft, ohne Befund; Mutationen für Lesen, Zustand `created`, Entfernen, Prüfen, Vorrang rot |
| `tools/test/run-integration-tests.sh` | geprüft; Befund F-1 (Verdrahtung des Exit-Codes), sonst ohne Befund |
| `tools/harness/run-container.sh`, `tools/harness/container-gegenprobe.sh` | geprüft, ohne Befund; `$pg` ohne `-v` an beiden Stellen einzeln rot gesehen (`p9-*-volume`, `p7-pg-vor-replay`) |
| `tools/test/postgres-versionen-gegenprobe.sh`, `tools/harness/gates-schnell-gegenprobe.sh` | geprüft, ohne Befund |
| `harness/mk/integration.mk`, `harness/mk/gates-schnell.mk` | geprüft, ohne Befund; Ziel hängt an `GATE_CHECKS` und steht in `GATE_SCHNELL_AUSNAHME`; `make gates-schnell-gegenprobe` grün |
| `make test-integration` (Volumes) | geprüft, ohne Befund; nach Namen kein neues Volume, nach dem Lauf keine Container des Präfixes |
| `harness/README.md` | geprüft; Zeilen vorhanden und mit Bindung; Befund F-3 in einer Aussage |
| `harness/vertraege/SPEC-057`, `SPEC-055` | geprüft; Befunde F-2, F-4; akzeptierte Negative konsistent mit Plan und Gegenprobe |
| Core-Reinheit, ADR-Immutabilität, Suppression (`//nolint`), Move/Inhalt | geprüft, ohne Befund (kein Produkt-Code, keine ADR, kein Move im Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zusage ohne Mutation an der Verdrahtung des neuen Vertrags · Vertragsaussage über Mengen roter Fälle weiter als die Prüfung · Zusage im Kommentar/Gate-Index weiter als Prüfung · Plan folgt der Korrektur nicht

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM; keine HIGH. Die Gegenprobe erkennt alle selbst gefahrenen Fehlbedienungen des Aufräumens (`-v`, Prüfung, Exit-Code der Prüfung, Lesen vor dem Entfernen, Signale); die Findings betreffen die Verdrahtung des Exit-Codes im Runner und Aussagen, die weiter reichen als die Prüfung.

**Übergabe:** Findings gehen an den Implementer (F-2 ggf. über den Architect: Vertragstext oder `nur`); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
