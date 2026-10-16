# Verifikation slice-v1-abschluss-verbindungsende-einspielen-extended

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-16. **Spanne:** `08290f8` (Basis der Marke) bis `1a1a888` (HEAD): Umsetzung `5cab7bf`, Review-Report `bcab0c6`, Architect `51145c0`, Nacharbeit `bf4fa5b`, `8b673fd`, `12e46f6`, `1a1a888`. Plan: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-einspielen-extended.md`.

**Eigene Läufe statt Zusage.** Unit-Mutanten in Kopien (`go.mod`, `go.sum`, `internal/`, `cmd/` im Scratchpad, Image aus `docker build --target deps`, `--network=none`, `go test -count=1 ./internal/hexagon/services/`, ein Lauf zur Zeit, Arbeitsbaum unberührt). E2E-Mutanten im Arbeitsbaum, `PGR_TESTS_1="TestE2EPlayEndeExtendedSync TestE2EPlayEndeExtendedFlush" make test-integration` (Server 17), danach `git checkout`; Arbeitsbaum am Ende sauber (`git status` leer). Dazu `make gates-schnell` (Exit 0, Schichten aus dem Cache des grünen Stands), `make abdeckung-check` (Exit 0), Messung nachgerechnet.

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt, die Nacharbeit trägt, kein Befund über LOW. Ein LOW (V-1: Kommentare, die für die letzte Gruppe mit `Flush` zu weit reichen) und drei INFO. Kein Produktcode von mir.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1: Einstufung `FATAL`/`PANIC` in der Extended-Interaktion (erwartet / `PGR-E4003`; Sync, Flush, frühere Gruppe, nach `ERROR`, normale Antwort/`ERROR` auf letzte Flush-Gruppe; kein `Terminate`, keine weitere Gruppe; Fall je Signal) | bestätigt | `traegtEnde`, `endetMitFlush`, Rückkehr vor `bisBereit` in `play.go`; Tests `TestPlayEndeExtendedErwartet` (24), `…NichtErwartet` (7), `…NachError`, `…LetzteFlush` (6), `TestPlayExtendedFatal` (umgekehrt), `TestPlayExtendedSignal` (3 Fälle); Mutanten unten rot |
| 2: Binary gegen PostgreSQL, Sync und Flush, Exit 0, nächste Session, eine Datei je Fall | bestätigt | `TestE2EPlayEndeExtendedSync`/`…Flush`, Fixtures `play-ende-extended-sync.yaml`, `…-flush.yaml` (unterscheiden sich in genau einer Zeile: `sync` gegen `flush`); E1 bis E3 unten selbst gefahren |
| 3: je Zusage eine Mutation, Tabelle in §7 | bestätigt | Tabelle mit 20 Zeilen plus E1 bis E3 und die Tabelle zu F-2 vorhanden; ich habe die Zeilen 1 bis 20 als Unit-Mutanten und E1 bis E3 am Binary nachgebaut (Mutant 15 und 6 am Binary: nicht wiederholt, laut §7 grün); die Angaben „roter Test“ stimmen, Abweichungen keine |
| 4: `make gates` grün | durch den Nachweis getragen | `.harness/state/gates-passed.dateiliste` (14:54:19, Stand `12e46f6`); die Liste unterscheidet sich vom Arbeitsbaum an `1a1a888` nur in der Plan-Datei (Positivliste, mit `dateiliste_erzeugen` gegen den Stand gelesen: eine Zeile); `stop-require-gates.sh` antwortet `approve`; `make gates-schnell` Exit 0 |
| 5: Review durchgeführt | bestätigt | Report `docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-einspielen-extended.md` (0 HIGH, 1 MEDIUM, 2 LOW, 6 INFO); kein Self-Review |
| 6 bis 9: Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §2 ohne Häkchen, §7 trägt Belege und Ausgang je Finding, §6-Risiken „offen bis Closure“; Vorschlag unten |

## Plan-vs-Code-Diff

§3 gegen `08290f8..1a1a888`: `play.go` (update, 39 Zeilen Produktcode, Spanne 25 bis 45), `play_ende_extended_test.go` (neu), `play_extended_test.go` (update), `test/integration/play_ende_extended_e2e_test.go` und zwei Fixtures (neu), `abdeckung-unit.md`/`abdeckung-e2e.md` (update). **Port `einspielziel.go` und Einspiel-Adapter unverändert, Bootstrap und CLI-Adapter unberührt** (§4 Rückführungen nicht eingetreten). Dazu außerhalb der Metrik: Plan, Roadmap, Review-Report, `spec/spezifikation.md` (Architect), Nehmer `next/slice-v1-abschluss-verbindungsende-vergleich.md`. Nichts Fremdes im Diff.

**Endstand sauber.** `git ls-files` kennt weder `nicht/` noch `Signal)/`; `git diff --name-status 08290f8..1a1a888` führt nur die 13 erwarteten Pfade; `make abdeckung-check` Exit 0. Hinweis zur History: `8b673fd` (und `bf4fa5b`) tragen die Kopie des Quellbaums unter `Signal)/` beziehungsweise `nicht/` (Hunderte Dateien als Zugang), `12e46f6` entfernt sie; wer die History aufräumt, ist Sache des Nutzers, der Endstand ist davon frei (V-4).

**Messung nachgerechnet.** `git diff --numstat 08290f8..bf4fa5b -- internal test docs/user`: 499 Einfügungen, 37 Löschungen, **536** (Marke 1000: 54 %, Meldeschwelle 800 nicht erreicht). `play.go` 32 + 7 = **39**; Tests im Service 236 + 74 + 28 = 338; E2E 55 + 43 + 43 = 141; Tabellen 6 + 12 = 18. Bis `1a1a888` unverändert 536 (die Commits danach ändern keinen Pfad der Metrik). Bis `5cab7bf`: 447. Stimmt mit §7.

## Nacharbeit (Nehmer, Signal, Kommentare, Spezifikation)

* **Nehmer `-vergleich` (`next/`).** §1 („Übernommen aus `-einspielen-extended`“), DoD Punkt 1 (Fall mit `Flush`), §6 (Bindung an (g), Zeile *Extended*) und §8 (Nachzählung: Liefer-Punkte weiter 3, Schichten zwei, +40 bis 70 Zeilen) tragen die Bindung mit der Kennung des Gebers; die frühere falsche Verweisung auf eine „Zeile *Extended*“ der DoD steht nur noch als Berichtigung in §8. `grep` nach der Kennung findet sie in den vier Abschnitten. **Konsistent.**
* **`TestPlayExtendedSignal`** mit drei Fällen (Sync-Bestand, erstes Signal / `Flush`, zweites Signal / `Flush`). Mutant „Rückkehr vor `bisBereit` entfällt“ färbt `…Signal/erstes_Signal,_letzte_Gruppe_mit_Flush` (ein `naechste` mehr im Ablauf) und die fünf Fälle von `…LetzteFlush`; der Sync-Fall bleibt grün. Mutant „`us.Schliesse()` im zweiten Signal entfällt“ färbt `…Signal/zweites_Signal,_letzte_Gruppe_mit_Flush` und `TestPlayZweitesSignalNachFehlerantwort`; Grund: der Fake blockiert in der Antwort der Gruppe, ohne Schließen endet Play nicht (30 s), also der richtige Grund, aber die Prüfung ist die Zeitüberschreitung, nicht ein Wert. **Richtig rot.**
* **Kommentare von `extended` und `interaktion`** stimmen mit `LH-FA-20.a` überein (Sync liest bis `ReadyForQuery`; `Flush` endet mit den Antworten oder der Fehlerantwort; nach früherem `ERROR` ohne Lesen) und jeder Satz hat einen Test (`…LetzteFlush`, `…Signal`). **Aber** V-1 unten.
* **`LH-FA-20.a` *Gruppen*:** Der Widerspruch („nur noch auf das `ReadyForQuery` nach ihrem `Sync`“) trägt jetzt die Ausnahme *(ausgenommen die letzte Gruppe mit `Flush`, unten)*; die Wortlaut-Sätze (Warten endet mit der Zahl, ohne Frist, mit und ohne Option; Sendefehler nach `ERROR` ist `PGR-E4003` bis zur gelesenen Fehlerantwort) und die Grenzen F-7/F-8 (`FATAL` nach vollständiger Zahl wird nicht gelesen; weniger Antworten als die Zählregel: Warten ohne Frist, nur Referenz 17) stehen dort und in §6 *Grenzen des Lesens* und in den Risiken (§6 Punkte 4 und 5) übereinstimmend. Das Abbruchsignal (Spezifikation Zeile 1987) nennt die Flush-Gruppe. Kein Widerspruch gefunden.

## Bewusstes Brechen (selbst gefahren)

Basis ohne Mutant: Paket `ok`. „rot“ nennt die gefärbten Tests, alle anderen blieben grün.

| Zusage | Mutant (Kopie von `play.go`) | Rot |
|---|---|---|
| `traegtEnde` erkennt Ende bei `Flush` (d, g) | `&& !endetMitFlush(letzte Gruppe)` | `…Erwartet`, genau die 8 `Flush`-Fälle |
| … bei `Sync` (c) | `&& endetMitFlush(letzte)` | `…Erwartet` 16, `…NachError`, `TestPlayExtendedFatal` |
| … Stelle ist die letzte Gruppe | erste Gruppe gelesen; irgendeine statt letzter Nachricht (`server[0]`) | `…Erwartet` 24, `…NachError`, `TestPlayExtendedFatal` (beide Mutanten gleich) |
| einfache Anfrage bleibt (Geber) | `endetEinfach` → `false` | `TestPlayEndeErwartet` 8 |
| (d) `FATAL` in früherer Gruppe erwartet | `l.ende = false` bei Fehlerantwort in nicht letzter Gruppe | `…Erwartet/…/frühere Gruppe` 8, `…NachError`, `TestPlayExtendedFatal` |
| Konjunktion: keine weitere Gruppe | `errSessionEnde` überlesen | `…Erwartet/…/frühere Gruppe` 8, `TestPlayExtendedFatal` |
| Konjunktion: kein `Terminate` | `Schliesse` statt `SchliesseOhneTerminate` | 66 FAIL-Zeilen, u.a. `…Erwartet` 24, `…NichtErwartet` 7, `TestPlayExtendedFatal`, die Fälle des Gebers |
| (e) `FATAL` nach `ERROR` | `l.ende = false` nach `fehlerGesehen = true` | nur `TestPlayEndeExtendedNachError` |
| Glied Option | `l.ende = traegtEnde(in)` | `…NichtErwartet` 3 (ohne Option Sync, mit `--continue-on-error`, Flush), `TestPlayEndeNichtErwartet` 2 |
| Glied Ende an der Stelle | `l.ende = Option` | `…NichtErwartet` 4, `TestPlayExtendedFatal`, Geber 3, `TestPlayErwarteterFehler` |
| Ende ist `connection_closed`, nicht irgendeine Fehlerantwort | `&& mitFehlerantwort(in)` | `…NichtErwartet/ohne Ende an der Stelle, aufgezeichnet ERROR` (die eine Zeile), `TestPlayExtendedFatal`, Geber |
| gilt ohne / mit `--continue-on-error` | `&& ContinueOnError` / `&& !ContinueOnError` | `…Erwartet` je 12, `…NachError` (nur das erste), `TestPlayExtendedFatal` (nur das zweite) |
| `verglichen`/`abweichend` bleiben 0 | `Verglichen++` / `Abweichend++` im Zweig des Endes | `…Erwartet` 24 und `…NachError`, `TestPlayEndeErwartet` 8 (je Zähler) |
| (g) kein Warten auf `ReadyForQuery` (Zählstand `naechste`) | Rückkehr vor `bisBereit` entfällt | `…LetzteFlush` 5, `…Signal/erstes Signal, Flush` |
| (g) `ERROR` beendet das Zählen | `n = 0` entfällt | `…NachError`, `…LetzteFlush` 3, `TestPlayExtendedFortsetzung` 2, `…ErwarteterFehler` 4, zwei Vergleichs-Tests |
| (g) nach früherem `ERROR` kein Lesen | `fehlerGesehen` ignoriert | `…NachError`, `…LetzteFlush/früherer ERROR` (beide), `…Fortsetzung`, `…ErwarteterFehler`, `TestVergleichExtendedFehlerOhneWarten` |
| `endetMitFlush` Sync-Seite / Flush-Seite | Helfer prüft `Sync`; Helfer immer `true` | 122 beziehungsweise 101 FAIL-Zeilen, über Sync- und Flush-Fälle beider Wege |
| Rückkehr nur für die **letzte** Gruppe | `endetMitFlush(gruppen[0])` | 38 FAIL-Zeilen (u.a. `…Erwartet` 8, `…NichtErwartet` 4, `…Signal` 1) |
| Schwere: `V` vor `S`; `S` zählt ohne `V`; `PANIC` wie `FATAL`; `FATAL` zählt | `S` zuerst; `S` ignoriert; `PANIC` gestrichen; `FATAL` gestrichen | je Klasse 12 Fälle `…Erwartet` (PANIC/S/V) bzw. alle; `S`/`V`/`PANIC`: je 12; `FATAL` gestrichen: 12 `…Erwartet` + 7 `…NichtErwartet` + 5 Geber + weitere |
| zweites Signal schließt | `us.Schliesse()` im Rückruf entfällt | `…Signal/zweites Signal, Flush`, `TestPlayZweitesSignalNachFehlerantwort` (30 s) |
| E1 Binary, Sync | Ende nur bei `Flush` erkannt | `TestE2EPlayEndeExtendedSync` rot (Exit 4), `…Flush` grün |
| E2 Binary, Flush | Ende nur ohne `Flush` erkannt | `TestE2EPlayEndeExtendedFlush` rot (Exit 4), `…Sync` grün |
| E3 Binary, Option | `l.ende = false` | beide E2E-Tests rot (Exit 4) |

**Grün gebliebener Mutant, eingeordnet.** `n > 0` in `traegtEnde` entfällt: grün (exit 0); äquivalent für geladene Aufzeichnungen, weil `Group.validate`/`pruefeEnde` jede Server-Folge mit Ende als nicht leer garantieren (Review F-6, akzeptiertes Negativ des Architect, bestätigt). Mutant Rückkehr vor `bisBereit` am Binary: laut §7 grün (die `FATAL` kommt als erste Antwort), äquivalent für diesen Fall, die Zusage trägt der Service-Test; von mir nicht wiederholt (Server 17, `5cab7bf`-Lauf des Implementers).

**Konjunktions-Glieder:** Option, Ende an der Stelle, Ende ist `connection_closed`, kein `Terminate`, keine weitere Gruppe, `--continue-on-error` beide Richtungen: je ein Mutant, der nur sie bricht, alle rot. **`FATAL`/`PANIC` × `V`/`S`:** vier Fälle `fatalFaelle` je mit und ohne `--continue-on-error` (24 Fälle in `…Erwartet`), alle drei Schweregrad-Mutanten färben nur ihre Klasse. **Abdeckung:** `make abdeckung-check` grün, jede der neuen Tests trägt eine Deklaration (`LH-FA-20`, `LH-FA-25`, `LH-FA-18`), die Fixtures sind je Fall eine Datei.

## Findings

| # | Schwere | Befund | Zuständig |
|---|---|---|---|
| V-1 | LOW | **Kommentare reichen für die letzte Gruppe mit `Flush` weiter als der Test (`AGENTS.md` §3.11, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`).** `play.go:129` (Kommentar von `Play`): „eine laufende Interaktion [läuft zu Ende] bis zu ihrem ReadyForQuery“, die Spezifikation (Abbruchsignal) und `TestPlayExtendedSignal` sagen für die Flush-Gruppe „bis zur Antwort auf jede Client-Nachricht“; das ist genau der Fall von F-2, der Kommentar wurde in `bf4fa5b` nicht nachgezogen. Ebenso `play.go:247` (*Mit CompareResponses hält sie die Antworten bis zum ReadyForQuery*) und `:451` (`abschluss`: „bis zu ihrem ReadyForQuery gelesen“) sowie `play_vergleich.go:14` (Kommentar von `Zaehler`: `Verglichen` „bis zum ReadyForQuery oder bis zur Fehlerantwort“), die die Spezifikation (`LH-FA-24.a`, Zähler) für die Flush-Gruppe schon anders sagt. Der Code läuft dort so (Rückkehr vor `bisBereit`, `abschluss` ohne `ready_for_query`), aber kein Test dieses Slice prüft den Vergleichspfad (siehe V-3). | `:129`: Implementer dieses Slice, vor der Closure (ein Satz). `:247`, `:451`, `play_vergleich.go:14`: Nehmer `slice-v1-abschluss-verbindungsende-vergleich` (Bindung in §1/DoD/§6/§8 trägt den Pfad; dort mit dem Fall); der Nehmer nimmt an |
| V-2 | INFO | Der Mutant „`us.Schliesse()` im zweiten Signal entfällt“ färbt die Fälle über eine Zeitüberschreitung von 30 s, nicht über einen Wert (`TestPlayExtendedSignal` und `TestPlayZweitesSignalNachFehlerantwort`). Die Prüfung ist trotzdem eindeutig (ohne Schließen endet Play nicht), und der Mutant kostet im Fehlerfall 30 s je Fall. Kein Handlungsbedarf, nur Kosten des Fehlerfalls. | — |
| V-3 | INFO | Neue Zusagen von `LH-FA-24.a` *Ende der Aufzeichnung* (vierter Punkt) und *Zähler* (`verglichen` bis zur Antwort auf jede Client-Nachricht) sind in der Spezifikation, vom Architect vor dem Code gesetzt, aber ohne Test, bis der Nehmer liefert; der Einleitungssatz der Stelle („Die Folge des Servers wird bis zum `ReadyForQuery` verglichen, oder … bis zu ihr“, Zeile 2459) nennt die Ausnahme nicht, nur der Aufzählungspunkt darunter. Der Mutant „Rückkehr vor `bisBereit` entfällt“ färbt keinen Vergleichs-Test. Das ist die Bindung an den Nehmer (Risiko §6, zweiter Punkt), durch den Nehmer getragen. | Nehmer `slice-v1-abschluss-verbindungsende-vergleich` (Fall mit Vergleich und der Einleitungssatz); Architect für den Wortlaut |
| V-4 | INFO | `bf4fa5b` und `8b673fd` tragen Kopien des Quellbaums (`nicht/`, `Signal)/`, vom Mutanten-Skript mit Namen aus mehreren Wörtern), `12e46f6` entfernt sie. Endstand sauber (siehe oben); die Zwischen-Commits blähen die History und die Dateiliste von `git log --stat`. Ob die History bereinigt werden muss, ist Entscheidung des Nutzers; die Closure ist davon nicht abhängig. Der Lerneintrag kann es als Beobachtung fassen (Skript mit Namen als ein Argument, `git add -A` nur nach `git status`). | Nutzer |

## Vorbereitung der Closure (Vorschlag, nicht ausgeführt)

Risiko-Ausgänge aus §6 (jeder genau einer; Urteil über das Tragen beim Planner):

| Risiko | Vorschlag | Begründung |
|---|---|---|
| Diff über einer Review-Sitzung | **entfallen** | gemessen 536 (54 % der Marke), Meldeschwelle 800 nicht erreicht, Schnitt vor dem Anspruch erfolgt |
| Bindung an `-vergleich` | **eingetreten → `slice-v1-abschluss-verbindungsende-vergleich`** | der Nehmer liegt in `next/`, nimmt in §1, DoD, §6, §8 an; V-1 (Kommentare) und V-3 gehen dorthin |
| Port/Adapter tragen den Pfad der Gruppen nicht (zweite Rückführung) | **entfallen** | Rückführung nicht eingetreten; Port und Adapter unverändert, die Tests der Warteschlange liegen beim Geber |
| Randform (d), (e), (g) nur gegen Fake / `FATAL` nur gegen PostgreSQL | **weiter offen → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`** | (d)/(e)/(g) sind entschieden und getestet (entfallen als Randform), der Rest (normale Antwort auf `Flush` ohne `Sync` gegen den Server) ist das akzeptierte Negativ und bleibt Register-Zähler |
| Hängen oder ungelesene `FATAL` bei anderer Serverversion | **weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`** | Grenzen (1)/(2) in `LH-FA-20.a`, nur Referenz 17 |
| Handgeschriebene Aufzeichnung ohne Beleg gegen den Server | **weiter offen → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`** | Beleg liegt für das Einspielen, Aufzeichnen in `slice-v1-abschluss-verbindungsende-upstream-extended` |
| `57P01` nur gegen Referenz 17 | **weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`** | dieselbe Beobachtung wie die Grenze oben: ein Zählerschritt, kein zweiter |

Lerneintrag, Vorschlag: Kommentare und Plan-Zeile, die eine Zusage zu einer Ablaufgrenze nennen („bis zum ReadyForQuery“), sind beim Ändern der Grenze in **jedem** Pfad nachzuziehen, nicht nur im geänderten (V-1; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und `BEO-REPO/plan-folgt-korrektur-nicht` zitieren, nicht neu benennen). Beobachtung (aus V-4): Mutanten-Skripte nur mit einem Argument ohne Leerzeichen und Sonderzeichen aufrufen; `git add -A` erst nach `git status`.

## Offene Handlungen vor `done/`

1. V-1 für `play.go:129` im Slice nachziehen (Kommentar, ein Satz), dann Closure-Notiz, Register, Risiko-Ausgänge (Vorschlag oben), Häkchen, reiner `git mv` nach `done/`.
2. V-1 (Rest) und V-3 im Nehmer `-vergleich` als Teil seiner DoD (Fall mit `Flush` und Vergleich) halten.
