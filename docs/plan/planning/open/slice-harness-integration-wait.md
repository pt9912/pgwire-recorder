# Slice slice-harness-integration-wait: Ein `Wait` je Prozess im Integrations-Testgeschirr

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-07: nach `slice-harness-lint` und vor `slice-harness-abdeckung-gate`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert nur das Testgeschirr, keine Zusage
des Produkts und keinen Nachweis; die Abdeckungs-Deklarationen der Tests bleiben
unverändert. Bindung an Entscheidungen:
[ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)
(Integrationstests in der Stufe `integration` gegen ein gepinntes PostgreSQL-Image),
[ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckungs-Deklaration
je Test, hier nur unverändert gehalten).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Im Integrations-Testgeschirr unter `test/integration` ruft für jeden
gestarteten Prozess genau eine Stelle `Wait` auf dem `exec.Cmd` auf; `stop`, die
Goroutinen der Tests und der `t.Cleanup` aus `startProzess` lesen ihr Ergebnis, statt
selbst zu warten. Ein Test, dessen Prozess nach dem Signal nicht endet, wird damit
nach seiner eigenen Frist rot, statt bis zum Zeitlimit von `go test` zu hängen.

**Übernimmt:** aus `slice-v1-abschluss-betrieb` das Risiko *Testgeschirr wartet zweimal
auf denselben Prozess* (Verifikation V-88 zu `slice-harness-blackbox-einstieg`), nicht
dessen Gegenstand; dort steht der Ausgang *übernommen von*
`slice-harness-integration-wait`.

**Befund** (`docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-einstieg.md`,
V-88; Review-Report `docs/reviews/2026-10-07-review-slice-harness-blackbox-einstieg.md`,
F-458): Die Goroutine `beendet <- rec.cmd.Wait()` in
`TestE2ERecordExtendedSigtermBeimPipelining` (`extended_e2e_test.go`) und der
`t.Cleanup` aus `startProzess` (`Kill`, dann `cmd.Wait()`, `record_e2e_test.go`) rufen
`Wait` auf demselben `exec.Cmd` nebeneinander. Den Kanal der Kopier-Goroutinen liest
nur einer der beiden, der andere wartet in `awaitGoroutines` für immer. Unter der
Mutation I2 der Verifikation (`ClientMessage` ohne die Sperre `herunterfahren`) hing der
Test an `1e4381b` in 2 von 5 Läufen bis zum Zeitlimit, an `f28a2df` in 1 von 5. Dasselbe
Muster steht in `stop` (`record_e2e_test.go`, `done <- r.cmd.Wait()`) und in
`TestE2EReplayExtendedSigtermMittenInFolge` (`extended_replay_e2e_test.go`,
`done <- rep.cmd.Wait()`). Im Grünfall wirkt es nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code — Schicht-Abgrenzung: Der Slice ändert nur Testcode unter
  `test/integration`. Die Mutation I2 wird nur in einer Kopie des Arbeitsbaums gefahren
  und nie committet.
- Neue Tests, geänderte Erwartungen, Fristen oder Abdeckungs-Deklarationen — Bestand
  bleibt bewusst stehen: Testliste und Deklarationen sind die Messlatte, an der sich
  zeigt, dass nur das Warten umgebaut ist.
- Die Tests mit Signal und Frist je Modus (`--shutdown-timeout`) — übernimmt
  `slice-v1-abschluss-betrieb`; er schreibt sie über die Helfer, die dieser Slice
  umbaut.
- Mutationen in den Integrationstests als Gate — anderer Vorgang; `slice-harness-mutation`
  schließt sie in seinem §1 aus.
- Ein Zeitlimit für `go test` in `make test-integration` — anderer Vorgang am Gate:
  Ein kürzeres Limit verkürzte das Hängen, beseitigte es nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Ein `Wait` je Prozess: Unter `test/integration` ruft für jeden gestarteten
      Prozess genau eine Stelle `Wait` auf; `stop`, die Goroutinen in
      `TestE2ERecordExtendedSigtermBeimPipelining` und
      `TestE2EReplayExtendedSigtermMittenInFolge` und der `t.Cleanup` aus
      `startProzess` lesen ihr Ergebnis. Beleg in §7: die Fundstellen von `.Wait()`
      unter `test/integration` vor und nach dem Umbau (Suche mit Pfad und Zeile).
- [ ] Rot statt hängend: Unter der Mutation I2 der Verifikation (`ClientMessage` ohne
      die Sperre `herunterfahren`, nur in einer Kopie des Arbeitsbaums) wird
      `TestE2ERecordExtendedSigtermBeimPipelining` in jedem von mindestens zehn Läufen
      rot mit der Meldung des Tests („Recorder endet nicht binnen 5 s nach SIGTERM“) und
      endet weit vor einem Zeitlimit von `go test`, das ein Hängen sichtbar macht; am
      Stand vor dem Umbau zeigt dieselbe Messung mindestens einen hängenden Lauf. Beleg
      in §7: je Stand Quellstand, Zahl der Läufe, rote und hängende Läufe, Dauer je Lauf,
      gesetztes Zeitlimit.
- [ ] Testliste und Deklarationen unverändert: Die Liste der Tests unter
      `test/integration` (Namen, Zahl) ist vor und nach dem Umbau gleich, keine Zeile
      `Abdeckung:` ist geändert, und `make abdeckung-check` ist grün ohne neu
      geschriebene Tabellen. Beleg in §7: beide Listen oder ihr Vergleich und der Diff
      ohne Deklarationszeile.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/record_e2e_test.go` | refactor | `startProzess` startet nach `Start` genau eine Goroutine mit `Wait` und legt ihr Ergebnis in einen Kanal, der nach dem Ende geschlossen wird; `recorder` trägt ihn. `stop` und der `t.Cleanup` (`Kill`, dann Lesen) lesen ihn, statt `Wait` aufzurufen |
| `test/integration/extended_e2e_test.go` | refactor | `TestE2ERecordExtendedSigtermBeimPipelining` liest das Ende über den Kanal statt `beendet <- rec.cmd.Wait()`; Erwartungen und Fristen gleich |
| `test/integration/extended_replay_e2e_test.go` | refactor | `TestE2EReplayExtendedSigtermMittenInFolge` ebenso statt `done <- rep.cmd.Wait()` |

- Wer den Exit-Code liest (`stop`, die beiden Tests), liest `ProcessState` erst nach dem
  Ende, das der Kanal meldet; danach schreibt niemand mehr daran.
- Der Implementer prüft vor dem Code mit einer Suche nach `.Wait()` und `cmd.Wait`
  unter `test/integration`, ob weitere Stellen dasselbe Muster tragen, und nennt sie
  hier.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-commit-struktur-id` liegt in `done/`
(WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05, 2026-10-06 und
2026-10-07: `slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint`,
`slice-harness-commit-struktur-id`, dieser Slice, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-harness-mutation`. Technisch hängt dieser Slice nur an
`slice-harness-lint`: Ab ihm ist `make lint` Gate auch für `test/integration`. Er steht
vor `slice-v1-abschluss-betrieb`, der die Tests mit Signal und Frist über dieselben
Helfer schreibt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Suche aus §3 findet das
  Muster in so vielen weiteren Tests, dass der Umbau nicht in einer Review-Sitzung
  prüfbar ist; dann trägt dieser Slice den Helfer und die drei genannten Stellen, die
  übrigen ein eigener Slice vor `slice-v1-abschluss-betrieb`.
- `in-progress` → `open` (blockiert — Carveout?): Am Stand vor dem Umbau zeigt die
  Messung aus DoD-Punkt 2 in keinem Lauf ein Hängen, und auch mehr Läufe oder eine
  andere Mutation, unter der der Prozess nach dem Signal nicht endet, führen nicht dazu;
  dann belegt die Messung den Fix nicht, und der Architect entscheidet einen anderen
  Nachweis.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `test-integration` und `abdeckung-check`, die
Messung aus DoD-Punkt 2 in §7 mit null hängenden Läufen nach dem Umbau, Closure-Notiz
mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen:** Der Slice liefert keinen neuen Vertrag (`AGENTS.md` §3.10, §3.12): Er
ändert keine Zusage des Produkts und keines Gates, nur das Warten im Testgeschirr. Die
Fälle, die der Umbau tragen muss, stehen darum als Risiken unten; entscheidet der Umbau
etwas, das hier nicht steht, gibt der Implementer es dem Architect zurück.

**Risiken:**

- **Messung unterscheidet nicht** — das Hängen trat an `1e4381b` in 2 von 5 Läufen auf;
  bei wenigen Läufen kann der Stand vor dem Umbau zufällig nicht hängen, und dann
  belegt der Stand nach dem Umbau nichts. Mindestens zehn Läufe je Stand, das
  Zeitlimit so gesetzt, dass ein Hängen als Abbruch erscheint (Rückführung in §4). —
  **Ausgang:** — (bei Closure)
- **Mutant kommt im Build-Kontext nicht an** — die Mutation I2 liegt im Produkt-Code
  und läuft über die Stufe `integration` des `Dockerfile`; BuildKit überträgt eine Datei
  gleicher Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`,
  1×). Je Lauf eine frische Kopie unter eigenem Pfad, wie in der Verifikation. —
  **Ausgang:** — (bei Closure)
- **Erste Phase verdeckt die zweite** — der Integrations-Runner endet nach einer roten
  ersten Phase (`BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar`, 1×); die
  Messung läuft darum gezielt für den einen Test (`-run`), nicht über
  `make test-integration`. — **Ausgang:** — (bei Closure)
- **`t.Cleanup` nach fehlgeschlagenem Start oder nach `Fatalf`** — der Cleanup läuft
  auch, wenn `stop` den Prozess schon getötet hat oder der Test vor dem Lauschen
  abbricht; er darf dann weder ein zweites Mal warten noch auf einen Kanal warten, der
  nie geschlossen wird. — **Ausgang:** — (bei Closure)
- **Neuer Fund am Testgeschirr** — F-458 nennt aus dem Bestand die Goroutine in
  `TestE2ERecordExtendedSigtermBeimPipelining`, die `pc` beim `Close` noch benutzt; der
  Umbau ändert daran nichts und darf es nicht verdecken. — **Ausgang:** — (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Belege zur DoD (Implementer):** <Fundstellen von `.Wait()` vor und nach dem Umbau;
  Messung unter I2 je Stand (Quellstand, Läufe, rot, hängend, Dauer, Zeitlimit);
  Testliste vor und nach dem Umbau und der Diff ohne Deklarationszeile>
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `cceb17f` gesichtet, dazu der
Nachtrag von F-336 im selben Commit wie dieser Plan (Zähler = Dateien unter
`evidence/`). V-88 selbst steht nicht im Register; seine Adresse ist dieser Slice.
Treffer:

- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — die Messung unter I2 baut
  einen Mutanten über den Build-Kontext; ein Risiko in §6.
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` (1×) — der Runner endet nach
  der ersten roten Phase; ein Risiko in §6.
- `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (3×, verkörpert) — darum stehen
  die Belege der DoD in §7, nicht im Bericht.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (16×, verkörpert in `AGENTS.md`
  §3.11) — ein neuer Kommentar am Helfer sagt nur zu, was die Messung zeigt.
- `BEO-REPO/plan-folgt-korrektur-nicht` (15×, verkörpert in §3.9) — findet die Suche
  aus §3 weitere Stellen, folgen §1 und §3 im selben Commit.

Keiner der Einträge erreicht mit diesem Slice die Schwelle 3×; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
