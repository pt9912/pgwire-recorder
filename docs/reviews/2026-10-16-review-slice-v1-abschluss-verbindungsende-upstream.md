# Review-Report: slice-v1-abschluss-verbindungsende-upstream — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41, Ende des Servers vom echten Upstream aufzeichnen), Diff `c5089444..5e13e1a6` (Basis `c5089444` laut Plan §4; Pfade `internal/adapters/driven/postgres internal/hexagon/ports/driven test/integration docs/user`, dazu `spec/spezifikation.md`, Plan und Roadmap).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
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

- Slice-Plan `slice-v1-abschluss-verbindungsende-upstream` (§1, §3, §4, §6, §7 Belege des Implementers, §8); Geber `slice-v1-abschluss-verbindungsende-record`, `slice-v1-abschluss-postgres-versionen-zusage`, Nehmer `slice-v1-abschluss-verbindungsende-upstream-extended`
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), `SPEC-029` Punkt 2, `SPEC-034`
- `LH-FA-25` (`LH-FA-25.a` *Aufzeichnen*, *Stelle*, *Leerlauf* mit der Berichtigung vom 2026-10-16, *Kein Ende*, *Senden nach dem Ende*), `LH-FA-02` (`LH-FA-02.b` *Gelesen*), Abnahmeszenario 18
- `AGENTS.md` (Hard Rules, §3.7 bis §3.14), Hexagon-Richtung (`make a-check`)

**Eigene Läufe des Reviewers** (Unit-Mutanten in Kopien von `go.mod`, `go.sum`, `cmd`, `internal` im Image `deps`, `--network=none`, je Mutant genau eine Ersetzung mit erzwungenem Treffer, ein Lauf zur Zeit; E2E-Mutanten im Arbeitsbaum mit `make test-integration` und Teilmenge, danach `git checkout`; Arbeitsbaum am Ende sauber):

- Aus der Tabelle in §7, rot wie behauptet: S1, S2, S4, S5, S7 (Schweregrad je Wert, `V` vor `S`, `fatal` kein Ende, leerer Schweregrad kein Ende), M1 (Merker bei `ReadyForQuery` leeren), D1, D3 (`Send` nach dem Ende, ERROR gegen FATAL), E1, Q4, R4 (`connection_closed` fehlt in `Query` und `Receive`), Q5, R5 (Lesepuffer hinter dem Ende, nur `TestEndeLiestNichtWeiter`, je nur der eigene Weg). Rot an den dort genannten Tests, die Zeilenmengen stichprobenhaft verglichen.
- Eigene Mutanten: N1 (nicht lesbare Antwort stuft über den Merker ein): rot an allen drei Zeilen von `TestNichtLesbareServerantwort`; N2 (nicht lesbare Antwort nennt den SQLSTATE des Merkers, der schärfere Mutant): rot allein an der Zeile „unbekannter Typ nach ErrorResponse“. **X1 grün** (siehe F-1).
- E2E: EM6 (Verbindungsende ohne Fehlerantwort als Ende werten) rot allein an `TestE2ERecordNetzwerkabbruch`; EM1 (SQLite-Version nicht setzen) rot allein an `TestE2ERecordEndeTerminate/rec.sqlite`.
- `make a-check`: 0 Befunde. `make abdeckung-check`: grün. Die Mutanten EM2 bis EM5, EM7 bis EM12 und RM1/RM2 nicht neu gefahren (Zeit; EM5 entspricht Q4, EM12 ist als grün angegeben).
- Messung nachgerechnet: `git diff --numstat c5089444..5e13e1a6 -- internal/adapters/driven/postgres internal/hexagon/ports/driven test/integration docs/user` gibt 818 Einfügungen und 57 Löschungen, Summe **875** (stimmt; der Commit `5e13e1a6` ändert in diesen Pfaden nichts). Aufschlüsselung stimmt (Produktcode 51+14=65, Port 13+2=15, Adapter-Tests 325+39+39=403, E2E 370, Tabellen 22). Nach dem Adapter-Teil (`37e78f76`): **483**, nicht 475 (F-3).

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „es wird nichts mehr gelesen“ gilt auch für `Send` (`LH-FA-25.a` *Senden nach dem Ende*); §7 nennt den Mutanten „`Send` liest“ äquivalent mit der Invariante „`Send` enthält keinen Lesezugriff“, die der Mutant gerade bricht. Der Mutant X1 (`Send` liest nach gelesener `FATAL` mit 30 ms Frist eine Nachricht und verwirft sie) bleibt grün: `TestSendNachDemEnde` (Schreiben gelingt und scheitert) und alle anderen Tests; die Gegenargumentation „D4/D5 fängt ihn“ trifft nur `Send` vor der gelesenen Fehlerantwort. Ein lesendes `Send` ohne Frist hinge im Test ohne Zeitgrenze (`TestSendNachDemEnde` nutzt `fertigBinnen` nicht, anders als `endeLauf`). | `AGENTS.md` §3.10, `AGENTS.md` §3.11; `LH-FA-25.a` *Senden nach dem Ende*; Register `BEO-REPO/gruener-mutant-faelschlich-aequivalent` | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-upstream.md` · „äquivalent. Invariante: `Send` enthält keinen Lesezugriff“ | ja — X1 (nach `Flush` ohne Fehler: `if s.endeGelesen() { SetReadDeadline(30 ms); fe.Receive() }`) im Image `deps`, `go test ./internal/adapters/driven/postgres/` grün | Äquivalenz-Aussage über Invariante, die der Mutant bricht (grüner Mutant als äquivalent) |
| F-2 | LOW | Der Port-Kommentar sagt, nach `connection_closed` rufe der Aufrufer nur noch `Close`, und im selben Absatz, `Send` liefere danach keinen Fehler; der Use Case ruft `upstream.Send` in `ClientMessage` nach der Prüfung von `serverEnde` ohne gehaltenes `mu`, so dass ein `Send` neben dem lesenden `Receive` nach `connection_closed` laufen kann. Plan §6 („kein Aufrufer, kein Test“) und der Satz „was Query, Send und Receive dann liefern, ist nicht zugesagt“ widersprechen der Zusage für `Send` zwei Sätze davor; die Aufrufer-Verpflichtung hat in diesem Slice keinen Test. | `AGENTS.md` §3.11 | `internal/hexagon/ports/driven/upstream.go` · „ruft der Aufrufer nur noch Close“ | nein — kein Gate-Lauf prüft eine Aufrufer-Verpflichtung gegen den Use Case | Kommentar sagt Aufrufer-Verpflichtung zu, die der Aufrufer nicht einhält |
| F-3 | LOW | Die Messung des Adapter-Teils in §7 nennt 475; `git diff --numstat c5089444..37e78f76` über die Pfade des Plans gibt 428 Einfügungen und 55 Löschungen, 483. Die Meldeschwelle 650 ist davon nicht berührt, die Gesamtzahl 875 stimmt. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-upstream.md` · „vor dem ersten Commit unter `test/integration`) **475**“ | ja — `git diff --numstat c5089444..37e78f76 -- <Pfade aus Plan §4>` | Zahl in der Messung weicht von der Metrik ab |
| F-4 | LOW | Die Zeile von `upstream.go` in §3 nannte nach den Code-Commits `37e78f76` und `51193b48` noch `verbindungsende` und nicht die neuen Funktionen `beendetVerbindung` und `endeGelesen`; sie folgt erst im letzten Commit `5e13e1a6`, der nur den Plan ändert. Der Plan hat dem Code damit zwei Commits lang nicht im selben Commit gefolgt. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-upstream.md` · „`lies` und neu `beendetVerbindung`, `endeGelesen`“ | ja — `git show 37e78f76 -- <Plan>` zeigt nur die Zeile der Tests | Plan folgt der Korrektur nicht im selben Commit |
| F-5 | INFO | Die Regel „Ende = `FATAL` oder `PANIC`, gelesen aus `V`, ohne `V` aus `S`“ steht jetzt an drei Stellen: `Response.beendetVerbindung` im Modell, `fehlerantwort` in `play.go` und `beendetVerbindung` im Upstream-Adapter. Die Tests des Adapters pinnen seine Kopie (`TestEndeSchweregrad`, zehn Zeilen); eine Änderung der Regel an einer Stelle färbt die anderen nicht. Hinweis für Planner und Architect, keine Aktion in diesem Slice. | Maintainability | `internal/adapters/driven/postgres/upstream.go` · „func beendetVerbindung(fehler map[string]string) bool“ | nein | Dreifache Definition derselben Schweregrad-Regel |
| F-6 | INFO | Die Probe der SQLite-Version liest Rohbytes („version2“ vorhanden, „version1“ nicht). Sie trägt aus dem Grund, den §7 nennt (Überschreiben an Ort und Stelle; EM1 selbst rot gesehen); sie hängt am Speicherlayout von SQLite (freie Seiten, Journal), würde „version20“ annehmen und nennt `meta.version` nicht über einen Leser. | Maintainability | `test/integration/ende_e2e_test.go` · `strings.Contains(text, "version2")` | ja — EM1 gefahren, rot | Probe über Speicherlayout statt über Leser |
| F-7 | INFO | Der Formvergleich gegen die handgeschriebenen Fixtures (`gleicheForm`) läuft nur für YAML; die SQLite-Aufzeichnung wird in Phase 1 nur auf `57P01` und die Version geprüft, `V`, `S`, `M` der SQLite-Fehlerantwort erst über `endeSicht` im Replay der Phase 2 (`TestE2EOhnePostgresEnde`). Für den Verifier: ob DoD 3 für SQLite damit gedeckt ist. | `LH-FA-25.a` | `test/integration/ende_e2e_test.go` · „gleicheForm(t, text, endeFatalYAML)“ | ja — `make test-integration` mit beiden Phasen | Formvergleich nur für ein Format |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/postgres/upstream.go` (Schweregrad, Merker, `Query`, `Send`, `Receive`, `lies`) | geprüft, ohne Befund; 13 Mutanten aus §7 und die eigenen Mutanten N1/N2 rot, `lies` hat genau zwei Aufrufer (`Query`, `Receive`), `Open` und `einspielen.go` lesen über `fe.Receive` und nicht über `lies` |
| Einordnung der sieben umgedeuteten Bestandstests (`upstream_test.go`) | geprüft, ohne Befund; `TestNichtLesbareServerantwort` „unbekannter Typ nach ErrorResponse“ ist mit `ERROR` und `57P01` aussagekräftig: Mutant N2 färbt genau diese Zeile (kein grün-ohne-Aussage) |
| `internal/adapters/driven/postgres/upstream_ende_test.go` | geprüft, ohne Befund; Wege, Fälle und Kommentare passen zu den Tests; Server bleibt offen (2 Sekunden Frist) |
| `internal/hexagon/ports/driven/upstream.go` | siehe F-2; Core-Reinheit gewahrt, nur Kommentare, `pgproto3` nur im Adapter (`make a-check`: 0 Befunde) |
| Leerlauf ohne nachfolgende Anfrage (`TestE2ERecordEndeLeerlauf/ohne_Anfrage_danach`) | geprüft, ohne Befund; die Begründung trägt: Invariante (`lies` nur von `Query` und `Receive`, `upstream.go`) und Ort (`services/record.go`) stimmen; ein Nichtgeschehen mit 300 ms Wartezeit, ohne Bezug zu F-1 |
| `kappProxy` und `TestE2ERecordNetzwerkabbruch` (EM12) | geprüft, ohne Befund; der Kommentar nennt `gekappt` als Grenze („keinen Test, der seinen Wegfall fängt“) und sagt nicht mehr zu; ein durchgeschlagener Frame würde den Test rot, nicht grün färben; EM6 rot allein hier |
| Formvergleich (`gleicheForm`, `endeFelder`, `nur`) | geprüft, ohne Befund außer F-7; liest die Zeile unmittelbar vor der ersten `connection_closed` strukturiert (Fluss- und Blockschreibweise, Fehlschlag bei fremdem Vorgänger), vier Felder je Fixture verlangt, kein Teilstring-Vergleich über den ganzen Output; Trennung von Leitungsprüfung (`pruefeBeendeAntworten`: `V`, `S`, `C`) und Formvergleich (YAML-Schreiber) wie in §7 Beobachtung 2 beschrieben |
| Fixtures und Unterdrückung | geprüft, ohne Befund; keine Fixture berührt, kein `//nolint` im Diff, kein altes Wort („bis einschließlich ReadyForQuery“ steht in `ports/driving/record.go` für den Use Case und bleibt) |
| `spec/spezifikation.md` `SPEC-029` Punkt 2 | geprüft, ohne Befund; der Satz zu Szenario 18 steht im Commit des Tests (`51193b48`), der Numerus-Satz davor steht seit `179ab424` |
| `docs/user/abdeckung-*.md` | geprüft, ohne Befund; `make abdeckung-check` grün |
| Randformen, die der Code still entscheidet (`AGENTS.md` §3.12) | geprüft, ohne Befund; gesucht in `Query`, `Send`, `Receive`, `lies` (leeres `V`, beide Felder leer, `Send` mit gelingendem Schreiben, Notices vor dem Ende, Lesepuffer): alle von §6 genannt |
| Plan §1, §4, §6, §8 gegen den Diff | geprüft, außer F-3 und F-4 ohne Befund; Adresse `slice-v1-abschluss-verbindungsende-upstream-extended` nimmt an (DoD 1 und 2, §1 des Nehmers), kein Zuwachs der Zählung |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Äquivalenz-Aussage über Invariante, die der Mutant bricht (grüner Mutant als äquivalent) · Kommentar sagt Aufrufer-Verpflichtung zu, die der Aufrufer nicht einhält · Zahl in der Messung weicht von der Metrik ab · Plan folgt der Korrektur nicht im selben Commit · Dreifache Definition derselben Schweregrad-Regel · Probe über Speicherlayout statt über Leser · Formvergleich nur für ein Format

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) blockiert: die Zusage „`Send` liest nichts nach dem Ende“ ist ohne Test und ihre Äquivalenz-Begründung trägt nicht; ob der Implementer einen Fall nachliefert oder die Zusage enger fasst (Spezifikation `LH-FA-25.a`, Port-Kommentar), entscheidet er mit dem Architect. F-2 bis F-4 (LOW) blockieren nicht; F-2 gehört mit F-1 zur Randform „`Send` neben `Receive` nach dem Ende“ und geht an den Architect.

**Übergabe:** Findings gehen an den Implementer (F-1, F-3, F-4) und den Architect (F-2, F-5); F-6 und F-7 an den Verifier. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler (F-1 betrifft `BEO-REPO/gruener-mutant-faelschlich-aequivalent`, F-4 `BEO-REPO/plan-folgt-korrektur-nicht`). Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
