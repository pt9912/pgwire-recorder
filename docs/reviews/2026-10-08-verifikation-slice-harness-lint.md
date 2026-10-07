# Verifikation: slice-harness-lint — 2026-10-08

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft ([Review](2026-10-08-review-slice-harness-lint.md), F-471 bis F-479). F-477 bis F-479 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-lint.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `fdbcfaf~1..29bcb16`. Architect: `fdbcfaf`, `7a32d57`, `53c70fc`. Implementer: `221845d`, `3bcf000`, `29bcb16`. Review: `c3c8ead`.

**Eingang:**

- Plan ganz, besonders §2 (DoD), §3, §6 (*Prüfung vor dem Code*, *Randform-Rückgaben aus `221845d`*, *Randform aus dem Review, F-475*, *Risiken*) und §7 (*Belege des Implementers*, *Nachtrag zu den Rückgaben*, *Nachtrag zum Review*)
- `spec/spezifikation.md` `SPEC-049` ganz und die Historie; [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md)
- `AGENTS.md` §3.2, §3.9 bis §3.12
- Am Stand `29bcb16`: `harness/mk/lint.mk`, `tools/harness/lint.sh`, `tools/harness/lint-gegenprobe.sh` ganz, `.golangci.yml` ganz, `harness/sensors/lint.md`, die Vorlage `.harness/baseline/v6.16.0/templates/harness/sensors/gate.template.md`, `harness/README.md` §Sensors, `Makefile`, `tools/harness/record-gates.sh`, der Kommentar an `Meldungen` in `internal/hexagon/model/fehler.go`, die Stufe `lint` im `Dockerfile`
- Der Review-Report ganz, den Bericht des Implementers nicht

Bei meinem Start war der Arbeitsbaum sauber, HEAD `29bcb16`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

**Wie die Proben gebaut wurden.** Jede Kopie entstand mit `cp -R` aus dem Arbeitsbaum, ohne `-p`, in einem eigenen Unterverzeichnis des Scratchpads. Die Kopie für das rote Gate behielt `.git`, alle anderen nicht. Jede Mutation lief in einer frischen Kopie. Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft, prüft, dass die Datei sich geändert hat, und ruft danach `touch` auf die Datei. In jeder Mutanten-Kopie lief `make lint-gegenprobe`, das selbst noch einmal kopiert. Die Läufe liefen nacheinander, nie zwei Gegenproben gleichzeitig. Danach habe ich meine Kopien gelöscht, nur diese benannten Pfade. Docker habe ich nicht aufgeräumt. `make gates`, `make lint-gegenprobe` und `make docs-check` liefen im Arbeitsbaum, `git status --porcelain` war vor und nach `make lint-gegenprobe` gleich.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Gate `make lint`. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `harness/mk/lint.mk` hängt `lint` an `GATE_CHECKS` | `GATE_CHECKS += lint lint-gegenprobe`, ohne Ordnungskante; `record-gates: $(GATE_CHECKS)` im `Makefile` | bestätigt |
| `make gates` wird bei rotem Lint rot | Kopie **mit `.git`**, `//nolint` am Ende von `func Meldungen(…) {`. `make -k gates`: Exit 2. Rot sind genau `[harness/mk/lint.mk:15: lint] Fehler 1` mit `lint: internal/hexagon/model/fehler.go:175: Direktive nolint` und `[harness/mk/lint.mk:18: lint-gegenprobe] Fehler 1` (Folgerot über `p0-grundlauf`, wie die Sensor-Datei es nach F-478 sagt). `commit-msg-gegenprobe: gruen`, also ist `hook-gegenprobe` mit `.git` sauber. `make: Das Ziel „gates“ wurde wegen Fehlern nicht neugemacht.` Der Nachweis `.harness/state/gates-passed.diffsha` ist vor und nach dem Lauf gleich, `record-gates` lief nicht. Ohne `-k` endet `make gates` ebenfalls mit Exit 2. | bestätigt |
| `make gates` mit dem Gate grün | im Arbeitsbaum an `29bcb16`: Exit 0 (Abschnitt 4) | bestätigt |
| Bestand grün, ohne `//nolint` | frische Kopie mit einer zusätzlichen Datei `internal/frisch.txt`, damit der Build nicht aus dem Cache nimmt: `make lint` Exit 0, `0 issues.`, keine `lint:`-Zeile, die Stufe lief neu (9,0 s). **Damit ist die Vorbedingung der folgenden Slices erfüllt: `make lint` ist modulweit 0.** | bestätigt |
| `.golangci.yml` nur mit dauerhaften Ausnahmen nach Punkt 8, je mit `Why:`, keine ungenutzte | gelesen: 16 Regeln. Für Testdateien: Komplexität mit `funlen`, `noctx` und `unparam`, `revive` `unused-parameter`, `revive` `unused-receiver`. Dazu `ST1005`, zehn benannte Globale je mit Datei und Namen und `forbidigo` für `os.Std(out\|err)` unter `^(cmd\|test)/`. Jede trägt einen Block `# Why:`. Keine ungenutzt, das zeigt `p0-grundlauf` mit `warn-unused` und ohne `lint:`-Zeile. `git diff fdbcfaf~1..29bcb16 -- .golangci.yml` ist leer. | bestätigt |
| `testpackage` für jeden Pfad scharf | `skip-regexp: (^\|/)export_test\.go$`, sonst keine Ausnahme. V10 (`skip-regexp: ^nichts$`) wird rot, M2 (ohne Anker) ebenfalls (Abschnitt 2). | bestätigt |

### Punkt 2: Gegenprobe `make lint-gegenprobe`. **Nicht ganz bestätigt — V-96.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| an `GATE_CHECKS` | wie Punkt 1; `p10-gate-checks` liest die Datenbasis von `make` | bestätigt |
| Grundlauf: unveränderte Kopie mit Ausgang 0 und ohne `lint:`-Zeile | `p0-grundlauf` prüft beides. Mein Lauf: grün. | bestätigt |
| jeder rote Fall prüft seine eigene Zeile | gelesen: `befund` prüft Pfad, Zeile über einen Anker und den Linter, `lint_zeile` und `form` prüfen den Wortlaut, `allein` genau eine Zeile. Nur am Ausgang hängen `p9-laden` (nach der Grenze), `p10-gates-rot` und `p0-sammel`. Wie im Review (F-479). | bestätigt |
| Fall für Punkt 10: `make gates` endet bei rotem `make lint` mit Fehlerstatus | V1 (`-$(DOCKER_BUILD)` im Rezept von `lint`, Fehler ignoriert): **rot**, `p10-gates-rot` und `p9-lint-zeile-allein`. Die Gate-Kopie trägt nur `lint.mk`, ein Nachweis ohne `.git` kann den Fall also nicht zufällig rot machen. | bestätigt |
| Fälle aus der DoD-Aufzählung | Jeder steht im Skript: Linter je Gruppe (sogar je Linter, 17 + 8 + `revive` + `contextcheck` + `testpackage` = 28), `//nolint`, `gomodguard_v2`, Build-Tag, ungenutzte Regel, Regel ohne `Why:`, Profil fehlt, Schema, Form (Einzug, Flussform, Einzug 8 als Liste, Einzug 2 nur die erste Zeile), `p9-laden`, Feldfolge, `// x //nolint` und `// siehe nolint`, `lint:`-Zeile allein, `p7-whitebox-*` je Paketgruppe, `internal_test.go`, `Test…` und Variable in der Brücke | bestätigt, bis auf V-97 |
| je Zusage, die eine Mutation fangen kann, ein Fall | **zwei Zusagen ohne unterscheidenden Fall**, keine steht unter OFFEN: V7 (Punkt 6, `/` zwischen Kommentarzeichen und `nolint`) und V11 (Punkt 8, Doppelpunkt in `# Why:`). Beide grün (V-96). | **nicht bestätigt** |
| ausdrücklich benannte Zusagen aus `slice-harness-lint-werkzeug` | `^` in Pfaden: `p1-pfad-anker-*`. `dupl` an der Grenze: `p4-dupl-*` (M3 rot, siehe unten). `relative-path-mode`, `--network=none`, `-mod=readonly`, `GOTOOLCHAIN=local`, `-c`, Pin und Plattform sowie `Felder nicht erkannt` stehen unter OFFEN, je mit Grund. | bestätigt |
| keine Abdeckungs-Deklaration, `make abdeckung-check` grün | `grep -i abdeckung tools/harness/lint-gegenprobe.sh` ohne Treffer; `abdeckung-check` im Gate-Lauf grün | bestätigt |

### Punkt 3: Doku. **Bestätigt, mit V-98 und V-99.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `AGENTS.md` §3.2 nennt den Träger | Direktive `nolint` bricht `make lint`, Ausnahmen nur in `.golangci.yml` mit `# Why:`, Falsch/Richtig aus diesem Repo, Link auf die Sensor-Datei. Keine Platzhalter mehr. | bestätigt |
| `harness/README.md` | `make lint` steht unter §Sensors und verlinkt `sensors/lint.md`. Der Vertrag nennt `SPEC-049`, die Bindung ist [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) · seit slice-harness-lint. `make lint-gegenprobe` ist neu, mit Vertrag in einem Satz und Bindung. Die Zeile bei den Werkzeugen ist entfernt, „Nicht behauptet“ nennt nur Testabdeckung. | bestätigt |
| `harness/sensors/lint.md` nach der Vorlage | Abschnitte wie in `gate.template.md` (Vertrag, Grenze, Ausgabe und Ausgänge, Sperren, Bindung), Template-Hinweis und „Regeln dieser Datei“ entfernt. Vertrag als Link auf `SPEC-049`. Die Grenze gibt die der Spezifikation wieder, die Datei entscheidet nichts neu. Seit `29bcb16` steht dort auch das Folgerot der Gegenprobe (F-478). | bestätigt |
| kein Kommentar nennt `make lint` ein Werkzeug ohne Gate | `lint.mk` (Kopf und Hilfetext „Gate“), `Dockerfile` („Teil der Gate-Kette über `make lint`“). Im Arbeitsbaum steht „Werkzeug, kein Gate“ zu `lint` nirgends mehr. | bestätigt |
| Kopf von `lint.sh` zeigt auf die Sensor-Datei und ordnet jedem Punkt seine Fälle zu | Verweis auf `harness/sensors/lint.md`; `GEPRÜFT DURCH` hat Zeilen (1) bis (10) und Grenze. Gegen die echten Fälle gehalten: Jedes genannte Präfix und jeder Name existiert. Punkt 2 sagt „kein eigener Fall; offen“, was zum OFFEN der Gegenprobe passt (F-474 erledigt). Die Fälle zu `testpackage` aus Punkt 5 stehen nur unter (7) (V-98). | bestätigt, V-98 |
| Kommentar an `Meldungen` nennt die Tiefensuche | Der Diff ändert nur den Kommentar. Die Zusage prüft der Test `Tiefensuche` in `fehler_test.go` (Review F-479). | bestätigt |
| `lint.sh` ohne Verhaltensänderung | `git diff fdbcfaf~1..29bcb16 -- tools/harness/lint.sh` ohne eine Zeile außerhalb von Kommentaren; der Diff am `Dockerfile` ändert nur den Kommentar der Stufe `lint` | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt** (Abschnitt 4).

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-471 bis F-479) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, alle sechs Risiken in §6 tragen „— (bei Closure)“. Vorschläge in Abschnitt 5. |

---

## 2. Review-Befunde F-471 bis F-475: nachgefahren

| Befund | Probe | Ergebnis |
|---|---|---|
| F-471 (M1: Prüfung nach Punkt 6 ohne Testdateien) | M1 wie im Review | **rot**, `p6-testdatei` |
| F-472 (M2: `skip-regexp: export_test\.go$`) | M2 wie im Review | **rot**, `p7-endet-auf-export-test` |
| F-473 (M3: `dupl` `threshold: 149`) | M3 wie im Review | **rot**, `p4-dupl-gruen` mit „unerwarteter Befund (dupl) … `func doppeltKurzA(`“. Die Gegenrichtung (151) steht rot in §7. |
| F-474 (Kommentar-Zusagen) | gelesen | erledigt: `lint.sh` (2) „kein eigener Fall; offen“. Die Gegenprobe sagt „ohne Netz“ und „Arbeitsbaum unberührt“ nicht mehr zu, OFFEN nennt `--network=none` und „nichts in den Arbeitsbaum“. `lint.mk` verweist auf OFFEN. Rest im Hilfetext: V-99. |
| F-475 (`LINT_GEGENPROBE_PARALLEL`) | Skript direkt mit `0`, leer, `-1`, `abc`, ` 3`, `03`, `3 `, `DOCKER_BUILD=false`, `timeout 10` | jeder Wert: Exit 2, einzige Ausgabe `lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '<wert>'`, kein Build. Über `make … LINT_GEGENPROBE_PARALLEL=` kommt die Variable als gesetzt und leer an (Probe mit einem Mini-Makefile), wie §6 es voraussetzt. V13 (Muster `^[0-9]+$`) wird **rot**, `einstellung-0`. Die Fälle entsprechen §6 *Randform aus dem Review*. |
| F-476 (Plan folgt Punkt 8 nicht) | gelesen | erledigt: Kopf, §1 (Schicht-Abgrenzung) und §3 (Zeile `spec/spezifikation.md`) nennen die Änderung an Punkt 8 nach dem ersten Code-Commit mit `7a32d57`. |
| F-478 (Folgerot, Abbruch) | rotes Gate oben; gelesen | Die Sensor-Datei nennt das Folgerot. `aufraeumen` beendet die Hintergrund-Läufe und schreibt `ROT — Abbruch`, steht als OFFEN im Kopf. Nicht nachgefahren. |

---

## 3. Vollständigkeit der Zusagen von `SPEC-049`

Jede Zusage aus Kopf, Punkt 1 bis 10 und Grenze habe ich gegen einen Fall oder einen Eintrag unter OFFEN gehalten. Wo ich mir nicht sicher war, ob ein Fall wirklich unterscheidet, habe ich eine eigene Mutation gefahren.

| ID | Zusage | Mutation | Ergebnis |
|---|---|---|---|
| V1 | Punkt 10: `make gates` rot bei rotem `make lint` | `lint.mk`: `-$(DOCKER_BUILD) --target lint .` | **rot**: `p10-gates-rot`, `p9-lint-zeile-allein` |
| V2 | Punkt 4: `maintidx` unter 20 | `under: 21` | **rot**: `p4-maintidx-gruen` |
| V3 | Punkt 4: `nestif` 5 | `min-complexity: 6` | **rot**: `p4-nestif-rot` |
| V4 | Punkt 4: `interfacebloat` 10 | `max: 11` | **rot**: `p4-interfacebloat-rot` |
| V5 | Punkt 5: `errcheck` nur die fünf Funktionen | `- (io.Writer).Write` ergänzt | **rot**: `p5-errcheck-rot` |
| V6 | Punkt 5: `gomodguard_v2` nur zwei Module | `- module: golang.org/x/sync` ergänzt | **rot**: `p3-gomodguard` |
| V7 | Punkt 6: nach `//` oder `/*` auch `/` vor `nolint` | `lint.sh`: `[ ${tab}/]*` → `[ ${tab}]*` | **grün** (V-96) |
| V8 | Punkt 9: `uniq-by-line: false` | Zeile entfernt | **rot**: `p9-uniq-by-line-a`, `-b` und viele weitere |
| V9 | Punkt 5: `forbidigo` auch für `os.Stderr` | `pattern: ^os\.Stdout$` | **rot**: `p5-forbidigo-_ = os.Stderr` |
| V10 | Punkt 5: `testpackage` überspringt `export_test.go` | `skip-regexp: ^nichts$` | **rot**: `p0-grundlauf`, `p7-bruecke-kein-testpackage` und Folgefälle |
| V11 | Punkt 8: erste Zeile beginnt mit `# Why:` | `lint.sh`: `/^[ \t]*# Why:/` → `/^[ \t]*# Why/` | **grün** (V-96) |
| V12 | Punkt 5: `ireturn` für jedes Paket unter `ports/` | `ports(/[^.]+)?\.` → `ports\.` | **rot**: `p5-ireturn-frei-Ports` und Folgefälle |
| V13 | F-475: nur positive Zahlen | `^[1-9][0-9]*$` → `^[0-9]+$` | **rot**: `einstellung-0` |
| M1 bis M3 | Abschnitt 2 | | **rot** |

**Gegenprobe zu den beiden grünen Mutanten.** Auf dem Host mit dem Muster aus `lint.sh` gezeigt: `a++ // /nolint` und `a++ /*/nolint */` sind im Original ein Befund, im Mutanten V7 keiner. `///nolint` (`p6-drei-striche`) unterscheidet nicht, weil die letzten zwei Striche schon `//nolint` sind. Für V11: Ein Profil mit der ersten Kommentarzeile `# Why kein Doppelpunkt.` über einer Regel meldet im Original `lint: .golangci.yml:7: Regel ohne Kommentarblock "# Why:" unmittelbar darüber`, im Mutanten nicht. Beide Mutanten ändern also das Verhalten, und ein Fall kann sie fangen.

**Ohne Fall und ohne Eintrag unter OFFEN, aber kein Mutant des Werkzeugs (V-100):** Punkt 8 sagt „Zulässig sind nur diese Regeln“ und „Eine Regel, die einen Befund nur deshalb ausblendet, weil der Bestand ihn trägt, gibt es nicht.“ Das Werkzeug prüft beides nicht. Eine neue Regel mit `# Why:`, die einen neuen Befund ausblendet, lässt das Gate grün. Die Grenze von `SPEC-049` nennt das nicht, und der Kopf der Gegenprobe nennt es nicht unter OFFEN. Heute hält das Profil die Zusage ein (Punkt 1 oben).

Alle übrigen Zusagen haben einen Fall oder stehen unter OFFEN. Den Kopf der Spezifikation („schreibt nichts in den Arbeitsbaum“) habe ich einmal gemessen: `git status` gleich vor und nach `make lint-gegenprobe`. Nach allen Läufen der Gegenprobe lag kein Arbeitsverzeichnis von ihr unter `/tmp`.

---

## 4. Gate-Ergebnis und Laufzeit

Jeder Lauf allein, ohne anderen Docker-Lauf daneben, auf 20 Kernen:

| Lauf | Ergebnis | Zeit | §7 |
|---|---|---|---|
| `make gates` im Arbeitsbaum an `29bcb16` | Exit 0; die Stufe `lint` kam `CACHED`, `lint-gegenprobe: gruen` | 2 min 27 s | 2 min 49 s |
| `make lint-gegenprobe` im Arbeitsbaum, direkt danach | Exit 0, `lint-gegenprobe: gruen` | 1 min 24 s | 1 min 44 s |
| `make lint` in einer frischen Kopie | Exit 0, `0 issues.` | 10,4 s (Schritt der Stufe 9,0 s) | 8,6 s |
| `make -k gates` in der roten Kopie mit `.git` | Exit 2 (Punkt 1) | 3 min 19 s | — |

Die Gegenprobe ist mit rund 60 % des Gate-Laufs der längste Teil von `make gates`. Das Risiko *Laufzeit* in §6 ist damit eine Beobachtung mit Zahl.

**Commit-Messages:** Alle sieben nennen `slice-harness-lint` und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), `221845d` zusätzlich [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes). Keine nennt eine `SPEC-` oder `ARC-`Kennung.

---

## 5. Plan gegen Code

- **Kopf:** `Bezug` und `Berührte Spec-Stellen` stimmen mit dem Diff überein, auch mit `7a32d57` (F-476). `kopf-check` im Gate-Lauf grün.
- **§1:** Das Ziel ist erreicht, das Gate hängt an der Kette (Punkt 1). Die Abgrenzung hält: Profil, Stufe und `lint.sh` sind im Verhalten unverändert. Produkt-Code ändert sich nur im Kommentar an `Meldungen`. Keine Abdeckungs-Deklaration, kein Shell-Lint.
- **§2:** Die DoD verlangt noch „die Zeile der ungenutzten Regel … ohne erkannte Felder“ als Fall. §6 hat das in `7a32d57` als offen ohne Fall entschieden (V-97).
- **§3:** Jede Zeile hat ihre Änderung im Diff, und umgekehrt hat jede Änderung eine Zeile. Die Gegenprobe tut, was die Zeile sagt: `mktemp -d`, `cp -R` ohne `-p`, `touch`, bis zu sechs Läufe gleichzeitig, der Grundlauf zuerst, eine Gate-Kopie nur mit `lint.mk`.
- **§4:** Keine Rückführung ausgelöst. Zur Größe stimme ich F-477 zu: drei Liefer-Punkte, Harness und Doku. Das Skript ist lang (1177 Zeilen), aber gleichförmig: Fall und Erwartung je Zeile. V-96 ist eine Lücke in der Tabelle, keine Frage der Größe.
- **§6 *Randformen*:** Die Entscheidungen aus *Prüfung vor dem Code*, *Randform-Rückgaben aus `221845d`* und *Randform aus dem Review, F-475* hält der Code ein (`p8-leerraum`, `p8-strich-allein*`, `p8-form-strich-allein-8`, `einstellung-*`). Der Diff entscheidet keine neue Randform.
- **§6 *Risiken*, Vorschläge für die Ausgänge** (setzen muss sie die Closure):
  - *Bestand:* **eingetreten** → die sieben Folge-Slices aus §4 *Start* (in `done/`). Heute modulweit `0 issues.` (Punkt 1).
  - *Bestand wächst bis zum Gate nach:* **entfallen**. Begründung: `make lint` frisch an `29bcb16` ohne Befund, und ab jetzt prüft das Gate.
  - *Gegenprobe sieht den Mutanten nicht:* **entfallen**. Begründung: `cp -R` ohne `-p` mit `touch`; 14 von 16 Mutanten (V1 bis V13, M1 bis M3) rot, die zwei grünen sind echte Lücken der Fälle (V-96), kein Cache-Effekt.
  - *Konfiguration wirkt anders, als sie gelesen wird:* Urteil der Closure. Je Liste und je Ausnahme gibt es einen Fall (V5, V6, V9, V10, V12 rot). F-472 (Anker ohne Fall) und V-100 (zulässige Regeln ungeprüft) sind Prüf-Lücken, keine Konfiguration, die anders wirkt. Mein Vorschlag: **entfallen** mit dieser Begründung, also kein dritter Beleg für `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`.
  - *Bestehende Prüfung fällt weg:* **entfallen**. Begründung: Die Stufe `test` ist unverändert, der Diff am `Dockerfile` ändert nur den Kommentar der Stufe `lint`.
  - *Laufzeit von `make gates`:* **weiter offen** → Register, mit den Zahlen aus Abschnitt 4, oder **entfallen**, wenn der Planner 2½ Minuten für tragbar hält. Eine Lockerung ist es in keinem Fall.
- **§7 *Belege*:** Vollständig für das, was die DoD verlangt: Läufe, Fälle je Punkt, Mutationstabellen mit rotem Fall, grüne Mutanten eingeordnet, Nachträge zu den Rückgaben und zum Review. Die nachgemessenen Zahlen stimmen (Fälle unter Punkt 6: 13, unter Punkt 7: 14). Meine Zeiten sind kürzer als die in §7, wohl wegen des warmen Caches.
- **§8:** Für die Closure: V-96 ist dieselbe Klasse wie F-471 bis F-473 (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, Abgrenzung `BEO-REPO/liste-gruener-mutanten-unvollstaendig`). In diesem Slice wurde die Klasse damit in zwei Runden getroffen. V-99 gehört zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.

---

## 6. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-96 | MEDIUM (DoD Punkt 2; `AGENTS.md` §3.10; dieselbe Klasse wie F-471 bis F-473) | Zwei Zusagen von `SPEC-049` haben keinen Fall, der sie unterscheidet, und stehen nicht unter OFFEN. (a) Punkt 6: „gefolgt von beliebig vielen Leerzeichen, Tabs oder `/`“. Ohne `/` in der Klasse (V7) bleibt die Gegenprobe grün, weil `p6-drei-striche` (`///nolint`) auch ohne die Klasse trifft. Unterscheiden würde `// /nolint` oder `/*/nolint */`. (b) Punkt 8: „erste Zeile mit `# Why:`“. Ohne den Doppelpunkt im Muster (V11) bleibt die Gegenprobe grün. Unterscheiden würde eine Regel, deren erste Kommentarzeile `# Why kein Doppelpunkt.` ist; im Original ist das die Zeile `Regel ohne Kommentarblock "# Why:"`. Die DoD verlangt „je Zusage aus `SPEC-049`, die eine Mutation fangen kann, ein Fall“. **Vor dem Haken an Punkt 2:** je ein roter Fall in `sammel`, Adresse Implementer. | `tools/harness/lint.sh` · `nolint_muster="(//\|/\\*)[ ${tab}/]*`; · `zeile[i + 1] !~ /^[ \t]*# Why:/`; `tools/harness/lint-gegenprobe.sh` · `nolint_rot p6-drei-striche` | V7, V11 grün; Host-Probe mit Original und Mutant |
| V-97 | LOW (`AGENTS.md` §3.9) | Die DoD (Punkt 2) zählt „die Zeile der ungenutzten Regel mit Feldfolge **und ohne erkannte Felder**“ zu den Fällen. §6 *Randform-Rückgaben aus `221845d`* hat in `7a32d57` entschieden: `Felder nicht erkannt` ist offen im Kopf, ohne Fall. So steht es auch im Skript. Wörtlich lässt sich die DoD damit nicht abhaken. Adresse Architect bzw. Planner: den DoD-Satz auf „mit Feldfolge; ohne erkannte Felder offen im Kopf (§6)“ ziehen. | `docs/plan/planning/in-progress/slice-harness-lint.md` · „die Zeile der ungenutzten Regel mit Feldfolge und ohne erkannte Felder“ | Lesen |
| V-98 | LOW (DoD Punkt 3, Zuordnung) | Die Fälle aus F-471 und F-472 fehlen in den Listen der Köpfe. Der Kopf der Gegenprobe führt unter (6) `p6-testdatei` nicht und unter (5) und (7) `p7-endet-auf-export-test` nicht. Unter (5) nennt er für `testpackage` nur `p7-internal-test` und `p7-bruecke-kein-testpackage`. `lint.sh` (5) „p5-*, p1-cmd“ nennt die Fälle zur Einstellung von `testpackage` nicht, sie stehen nur über `p7-*` unter (7). Die Fälle laufen, nur die Zuordnung in den Köpfen ist unvollständig. Adresse Implementer, nur Kommentare. | `tools/harness/lint-gegenprobe.sh` · „p6-cmd, p6-test; kein Befund“; · „für testpackage p7-internal-test und p7-bruecke-kein-testpackage“; `tools/harness/lint.sh` · „(5) Einstellungen — p5-*, p1-cmd“ | Lesen gegen die Fälle |
| V-99 | LOW (`AGENTS.md` §3.11; Rest von F-474) | Der Hilfetext von `make lint` sagt „netzlos ausser deps“, die Vertragszelle in `harness/README.md` sagt „netzlos außer der Download-Stufe“. Die Gegenprobe führt `--network=none` als OFFEN, und §7 nennt den Mutanten äquivalent. F-474 hat diese Zusage in `lint.sh` und im Kopf der Gegenprobe enger gefasst, im Hilfetext und in der Zelle nicht. Mildernd: Die Zusage steht wörtlich im `Dockerfile` (`RUN --network=none`), und die Zellen von `make build` und `make test` sagen dasselbe ohne Test. Adresse Implementer: enger fassen oder als Eigenschaft der Stufe kennzeichnen. Ob „netzlos“ in einer Vertragszelle als Zusage zählt, kann der Architect einmal für alle Docker-Gates entscheiden. | `harness/mk/lint.mk` · „(Gate; netzlos ausser deps)“; `harness/README.md` · Zeile `make lint` | Lesen; §7 *Grüne Mutanten* |
| V-100 | LOW (Grenze unvollständig, Adresse Architect) | `SPEC-049` Punkt 8 sagt zu, dass nur die genannten Regeln zulässig sind und keine nur Bestand ausblendet. Das Werkzeug prüft das nicht: Eine neue Regel mit `# Why:`, die einen neuen Befund ausblendet, lässt das Gate grün. Die Grenze nennt es nicht, und OFFEN im Kopf der Gegenprobe auch nicht. Heute hält das Profil die Zusage ein (16 Regeln, gelesen). Vorschlag: ein Satz in der Grenze („ob eine Regel zu den zulässigen nach Punkt 8 gehört, prüft das Werkzeug nicht; Urteil des Review“) und derselbe Satz unter OFFEN bzw. bei den Grenzen der Gegenprobe. | `spec/spezifikation.md` `SPEC-049` · „Zulässig sind nur diese Regeln“; `tools/harness/lint-gegenprobe.sh` · „Ob ein `Why:` zutrifft“ | Lesen |
| V-101 | Hinweis (Planner) | Laufzeit: `make lint-gegenprobe` allein 1 min 24 s, `make gates` allein 2 min 27 s, bei warmem Cache. Die Gegenprobe ist der größte Posten. Für den Ausgang des Risikos *Laufzeit* (Abschnitt 5). | — | Abschnitt 4 |

Keine Abweichung von §6 *Randformen*. Kein Befund gegen DoD-Punkt 1, 3 (außer der Zuordnung V-98) und 4.

---

## 7. Urteil

**Fast fertig, ein Liefer-Punkt ist noch nicht ganz bestätigt.**

- **Gate (Punkt 1): bestätigt und selbst gezeigt.** `lint` und `lint-gegenprobe` hängen an `GATE_CHECKS`. In einer Kopie mit `.git` macht ein `//nolint` `make gates` rot (Exit 2). Rot sind nur `lint` mit der Zeile der Direktive und die Gegenprobe als Folgerot, `hook-gegenprobe` ist grün, und `record-gates` lief nicht. **`make lint` ist modulweit 0** (frischer Lauf, `0 issues.`), das Profil trägt nur dauerhafte Ausnahmen, und `testpackage` ist für jeden Pfad scharf. Die Vorbedingung der folgenden Slices ist erfüllt.
- **Gegenprobe (Punkt 2): Grundlauf grün, Verdrahtung gezeigt, F-471 bis F-473 und F-475 behoben** (M1, M2, M3 und V13 rot). Von 13 eigenen Mutationen sind 11 rot. **Zwei grüne zeigen Zusagen ohne Fall (V-96).** Das ist die Klasse, die die DoD ausschließen will. Punkt 2 ist abzuhaken, wenn die zwei Fälle da sind.
- **Doku (Punkt 3): bestätigt.** Sensor-Datei nach der Vorlage, `harness/README.md`, `AGENTS.md` §3.2 mit Träger, Kommentare zum Gate; `GEPRÜFT DURCH` stimmt mit den Fällen überein, unvollständig nur bei den zwei neuen Fällen (V-98). Rest von F-474 im Hilfetext: V-99.
- **`make gates` (Punkt 4): grün.**

**Vor der Closure offen:** V-96 (Implementer, zwei Fälle, blockiert den Haken an Punkt 2). V-97 (DoD-Satz, Architect bzw. Planner). V-98 und V-99 (Kommentare, Implementer). V-100 (Satz in der Grenze, Architect). Dazu die Closure-Pflichten mit den Vorschlägen aus Abschnitt 5.
