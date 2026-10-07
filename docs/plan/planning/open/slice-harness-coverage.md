# Slice slice-harness-coverage: Testabdeckung des Codes als Gate mit Schwelle

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei, das `slice-lastenheft-pruefbarkeit` anlegt und das mit
dem letzten Slice der Reihe nachweisbar wird. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 2). Bindung an Entscheidungen:
[ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Tests in der
Stufe `test` des Multistage-`Dockerfile`, netzlos),
[ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abgrenzung:
Abdeckung je Anforderung und Pfad), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Teil 2 der geteilten Messung,
Nachweisart Gate). Die ADR des neuen Gates schreibt der Architect
im Slice, vor dem Code; ihre Nummer vergibt der ADR-Index.

**Berührte Spec-Stellen:** `spezifikation.md §11` (*Harness-Werkzeuge*: neue Kennung des Gates mit Schwelle und Randformen, vergeben von diesem Slice)

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Stufe `test` des `Dockerfile` misst die Anweisungs-Abdeckung der
Unit-Tests (`go test -coverprofile`) und bricht ab, wenn sie unter der Schwelle
liegt; das Gate ist damit Teil von `make gates`. Die Schwelle ist eine Zahl aus der
Messung des heutigen Stands und steht nach Entscheidung des Nutzers vom 2026-10-06 im
Abschnitt für Harness-Werkzeuge der Spezifikation, in der Kennung des Gates, als
Konstante mit Einheit und Begründung; die ADR des Gates hält Mechanismus, Messung mit
Quellstand und die Regel, nach der die Zahl folgt, und setzt `Schärft:` auf die
Kennung — nicht in diesem Plan.
Die Gegenprobe zeigt, dass das Gate rot werden kann (`AGENTS.md` §3.10), und trägt die
Abdeckungs-Deklaration von Teil 2 von LH-QA-07. Eine
spätere Senkung der Schwelle braucht eine neue ADR (`AGENTS.md` §3.6).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: eine Coverage-Schwelle vor
dem nächsten großen Slice. `harness/README.md` §Sensors führt Testabdeckung heute
unter „Nicht behauptet … (keine Schwelle)“. Das Schwester-Repo ai-harness-init hat
keine Schwelle; es ist Kontext, kein Vorbild.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Abdeckung je Anforderung und Pfad (`make abdeckung-check`, Tabellen unter
  `docs/user/abdeckung-*.md`) — ein anderer Vorgang mit eigener Entscheidung
  ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)): Sie sagt, welche
  Anforderung durch welchen Test belegt ist, die Anweisungs-Abdeckung, welcher Code
  unter Tests läuft. Keine der beiden ersetzt die andere; Ziel- und Dateinamen dieses
  Slice tragen darum nicht das Wort „Abdeckung“ allein (§6 *Benennung*). Die
  Erweiterung des Abdeckungs-Skripts um die Nachweisart Gate und die Deklaration von
  Teil 1 und 3 liefert `slice-harness-abdeckung-gate` ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)); dieser Slice setzt
  nur die Deklaration von Teil 2 an seiner Gegenprobe.
- Neue Tests, um eine Zielzahl über dem gemessenen Stand zu erreichen — die Schwelle
  folgt der Messung; eine höhere Zielzahl ist eigene Arbeit mit eigenem Slice,
  aus der Messung heraus, nicht vorab.
- Ein Lint-Gate — übernimmt `slice-harness-lint`; die Umstellung der Tests auf
  Black-Box-Pakete übernehmen die vier Umstellungs-Slices aus §4.
- Eine Schwelle, die sich selbst anhebt (Ratsche) — ein anderer Vertrag mit
  Schreibzugriff auf die Konfiguration im Gate-Lauf; die Stufe `test` schreibt nichts
  in den Arbeitsbaum ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)).
  Eine Anhebung bleibt ein bewusster Commit.
- Ein Bericht als Datei (HTML, Profil im Arbeitsbaum) — ein Werkzeug, kein Gate;
  wird es gebraucht, ist es ein eigenes Werkzeug-Ziel.
- Produkt-Code, Lastenheft und die Spezifikation außerhalb ihres Abschnitts für
  Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert `Dockerfile`, Harness und
  ihre Doku; in der Spezifikation nur die Randformen seines Gates in jenem Abschnitt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Gate mit Schwelle: Die Stufe `test` erzeugt das Profil, prüft es gegen die
      Schwelle aus der Kennung des Gates im Abschnitt für Harness-Werkzeuge der
      Spezifikation (nach den Randformen aus §6) und endet rot
      darunter; `gofmt`, `go vet` und die Unit-Tests der Stufe bleiben unverändert. Die
      Messung des Stands, aus der die Schwelle folgt, steht mit Quellstand in der ADR;
      der Bestand ist grün.
- [ ] Gegenprobe als eigenes Ziel an `GATE_CHECKS`: ein Stand unter der Schwelle wird
      abgelehnt, einer genau an der Schwelle angenommen; dazu je entschiedener Randform
      aus §6, die eine Mutation fangen kann, ein Fall (Paket ohne Tests, ausgenommener
      Pfad, Rundung). Je Zusage ist die Mutation gesehen (`AGENTS.md` §3.10). Die
      Gegenprobe ist der Nachweis von Teil 2 (Anweisungsabdeckung) von LH-QA-07 und
      trägt nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) dessen Abdeckungs-Deklaration (Nachweisart Gate, Teil 2
      von 3); mit Teil 1 und 3 aus `slice-harness-abdeckung-gate` ist LH-QA-07 in den
      Abdeckungstabellen vollständig, `abdeckung-vollstaendig.md` und damit die RTM
      führen sie, und das Abnahmeszenario 17 ist nachweisbar (Tabellen mit
      `make abdeckung` neu geschrieben, `make abdeckung-check` grün).
- [ ] Doku: `harness/README.md` §Sensors führt die Schwelle beim Vertrag von
      `make test` und die Gegenprobe mit Bindung an die ADR; die Zeile „Nicht
      behauptet“ nennt keine Testabdeckung mehr und grenzt gegen die Abdeckung je
      Anforderung ab.
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
| `docs/plan/adr/<NNNN>-schwelle-testabdeckung.md`, `docs/plan/adr/README.md` | neu / update | ADR des Gates (Architect, vor dem Code): Messung mit Quellstand, Regel, nach der die Schwelle folgt, `Schärft:` auf die Kennung des Gates; Index-Zeile |
| `spec/spezifikation.md` | update | Abschnitt für Harness-Werkzeuge: Kennung des Gates mit der Schwelle als Konstante mit Einheit und Begründung und den Randformen aus §6 (Entscheidung des Nutzers vom 2026-10-06) |
| `Dockerfile` | update | Stufe `test`: Profil erzeugen und gegen die Schwelle prüfen, netzlos; die bestehenden Prüfungen der Stufe bleiben |
| Prüfung der Schwelle (Ort nach §6 *Werkzeug*) | neu | wertet das Profil aus, nennt je Paket und gesamt den Wert und die Schwelle |
| `harness/mk/build.mk` oder eigenes Fragment unter `harness/mk/` | update / neu | Schwelle als eine benannte Stelle (Build-Argument oder Konfiguration), Ziel der Gegenprobe an `GATE_CHECKS` |
| `tools/harness/<gegenprobe>.sh` | neu | Mutanten in einer Kopie unter eigenem Temp-Pfad: Code ohne Test hinzugefügt (rot), Stand genau an der Schwelle (grün), je Randform ein Fall; Name nach dem Muster `*-gegenprobe.sh`, im Kopf die Abdeckungs-Deklaration von LH-QA-07, Teil 2 von 3 |
| `docs/user/abdeckung-*.md` | update | mit `make abdeckung` neu geschrieben: LH-QA-07 vollständig |
| `harness/README.md` | update | §Sensors: Vertrag von `make test` um die Schwelle ergänzt, Gegenprobe als Zeile; „Nicht behauptet“ angepasst |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug`, die vier
Umstellungs-Slices (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`), die beiden
Bereinigungs-Slices (`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`),
`slice-harness-lint`, `slice-harness-commit-struktur-id`,
`slice-harness-integration-wait` und `slice-harness-abdeckung-gate` liegen in `done/`
(WIP-Limit 1). Dieser Slice ist der letzte der Reihe vor M3 nach Entscheidung
des Nutzers vom 2026-10-05, 2026-10-06 und 2026-10-07 (Lastenheft, Lint-Werkzeug,
Umstellung, Bereinigung, Lint-Gate, `slice-harness-commit-struktur-id`,
`slice-harness-integration-wait`, `slice-harness-abdeckung-gate`, Coverage). Nach ihm
folgen `slice-tests-ueberlebende-mutanten` und `slice-harness-mutation`; keiner der
beiden trägt zu M3 bei, und keiner ist Vorbedingung dieses Slice. Vom Abdeckungs-Gate hängt dieser Slice
technisch ab: Ohne es kennt `tools/test/abdeckung.sh` die Deklaration von Teil 2
nicht. Grund der übrigen Reihenfolge, keine technische Abhängigkeit: Lint-Bereinigung und
Black-Box-Umstellung ändern Code und Tests, an denen die Schwelle gemessen wird; eine
Messung davor wäre veraltet, bevor das Gate greift. Erster Schritt nach
dem Start, vor jedem Code-Commit: Der Architect misst den Stand (gesamt und je Paket,
mit Quellstand), entscheidet die Randformen aus §6 und schreibt die ADR mit der
Schwelle (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Entscheidung über die
  Schwelle (§6) verlangt neue Tests, bevor das Gate grün sein kann — etwa eine
  Schwelle je Paket, die Pakete unter dem gemessenen Stand nicht halten. Dann entstehen
  die Tests in eigenen Slices, und dieser liefert danach das Gate auf dem dann
  gemessenen Stand.
- `in-progress` → `open` (blockiert — Carveout?): Die Einbeziehung der
  Integrationstests (§6) ist entschieden, verlangt aber einen Weg aus dem Container
  der Integrationstests heraus, den
  [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nicht vorsieht
  (die Läufe schreiben nichts in den Arbeitsbaum); dann zuerst eine Entscheidung über
  diesen Weg.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit der Schwelle in der Stufe `test` und der
Gegenprobe, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — alle **offen**. Entschieden werden
sie vor dem ersten Code-Commit vom Architect und festgehalten im Abschnitt für Harness-Werkzeuge der Spezifikation (Technik-Stratum; angelegt von `slice-harness-vertraege-spezifikation`, ein bestehendes Werkzeug schreibt seine Kennung fort, ein neues bekommt die nächste freie Kennung, vergeben vom Slice, der es liefert); die ADR des
Gates trägt Entscheidung und Gründe und verweist mit `Schärft:` auf die Stelle. Eine
Entscheidung des Nutzers wird dort festgehalten. Was dort nicht steht, entscheidet der
Implementer nicht, er gibt es zurück.

- **Die Zahl** — erst nach der Messung des heutigen Stands; die Messung steht mit
  Quellstand in der ADR, die Zahl in der Kennung des Gates in der Spezifikation
  (Entscheidung des Nutzers vom 2026-10-06). Offen: Schwelle gleich dem gemessenen
  Wert, abgerundet auf welche Stelle, oder
  mit Abstand darunter (Spielraum für Refactorings, aber weniger Biss).
- **Gesamt oder je Paket** — eine Schwelle über alle Pakete (ein gut getestetes Paket
  verdeckt ein schwaches), je Paket (strenger; Pakete wie `bootstrap` oder `cli` mit
  wenig Logik können die Zahl der Kernpakete nicht halten) oder beides mit eigenen
  Zahlen. Gewichtung gesamt: nach Anweisungen (wie `go tool cover -func`) oder Mittel
  der Pakete.
- **Integrationstests über `GOCOVERDIR`** — `test/integration` startet das Binary
  als eigenen Prozess gegen PostgreSQL (`make test-integration`); mitgezählt würden
  sie nur mit einem `-cover`-Binary, `GOCOVERDIR` im Container und dem Zusammenführen
  der Profile über zwei Läufe hinweg. Offen: zählen sie mit, oder misst das Gate nur
  die Unit-Tests der Stufe `test` (und sagt das).
- **`-coverpkg`** — ohne die Option zählt ein Paket nur die eigenen Tests; mit
  `-coverpkg=./...` zählen auch Tests anderer Pakete (die Services-Tests decken
  Teile von `model`). Offen, welche Lesart die Zahl hat.
- **`cmd/` und Pakete ohne Anweisungen** — `cmd/pgwire-recorder` ist `main` ohne
  Test; die Port-Pakete und die `doc.go`-Pakete tragen keine Anweisung. Offen:
  ausgenommen (mit `Why:` an einer Stelle) oder mitgezählt; ob ein Paket ohne
  Anweisung 0 % oder „keine Aussage“ ist.
- **Paket ohne Tests** — neueres Go führt ein Paket mit Code, aber ohne Testdatei mit
  0 % im Profil. Offen: zählt es (dann zieht jedes neue Paket ohne Test die Zahl) oder
  ist es ein eigener Befund (Paket ohne Tests, rot) oder erlaubt.
- **Generierter Code** — der Bestand hat keinen (`Code generated` kommt nicht vor).
  Offen: Regel für später (Ausnahme nach Kopfzeile) oder ausdrücklich keine, bis es
  ihn gibt.
- **Vergleich und Rundung** — `≥ Schwelle` oder `>`, auf wie viele Nachkommastellen;
  der Stand genau an der Schwelle ist grün oder rot.
- **Werkzeug** — Auswertung mit Bordmitteln (`go tool cover -func` und ein kurzer
  Prüfschritt in der Stufe) oder ein Fremdwerkzeug (neues Modul oder Image, Pin,
  Netz in `deps`). Offen, ebenso, wo die Schwelle steht (Build-Argument aus dem
  Makefile, Konfigurationsdatei im Build-Kontext — dann `.dockerignore`).
- **`covermode`** — `set` (Default) oder `atomic`/`count`; mit `-race` wäre `atomic`
  Pflicht. Offen, ob die Zahl nur „mindestens einmal ausgeführt“ meint.
- **Ausgabe** — grün: eine Zeile gesamt mit Wert und Schwelle, je Paket der Wert;
  rot: dieselbe Ausgabe und die Differenz, Exit ≠ 0, die Stufe bricht ab. Offen: ob
  die Paket-Liste auch bei grün erscheint, Sortierung, und ob die Ausgabe ein Paket
  unter einer Paket-Schwelle beim Namen nennt.
- **Benennung** — „Abdeckung“ ist in diesem Repo die Abdeckung je Anforderung
  ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), `make abdeckung`,
  `make abdeckung-check`). Offen: Name des Gegenproben-Ziels und Wortlaut in der
  Sensors-Zeile, die beides nicht verwechselbar machen (etwa
  „Testabdeckung (Anweisungen)“ und ein Ziel `coverage-gegenprobe`).

**Risiken:**

- **Die Schwelle misst den falschen Stand** — gemessen vor einer Bestands-Bereinigung
  (`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`) oder vor der
  Black-Box-Umstellung (Black-Box-Tests
  erreichen unexportierte Pfade seltener), wäre die Zahl beim Einschalten schon
  überholt; darum die Reihenfolge in §4. — **Ausgang:** — (bei Closure)
- **Die Zahl wird zum Ziel** — Tests, die Zeilen ausführen, ohne etwas zu prüfen,
  heben die Zahl, nicht die Prüfung. Die Schwelle ersetzt nicht §3.10 (je Zusage eine
  Mutation); Review bleibt Urteil darüber. — **Ausgang:** — (bei Closure)
- **Bestehende Prüfung fällt weg** — der Umbau der Stufe `test` verliert `gofmt`,
  `go vet -tags integration` oder einen Testlauf
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Die Gegenprobe prüft, dass die
  Stufe bei einem Formatfehler weiter rot wird. — **Ausgang:** — (bei Closure)
- **Gegenprobe sieht den Mutanten nicht** — BuildKit überträgt eine Datei gleicher
  Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×).
  Die Gegenprobe baut aus einer Kopie unter eigenem Temp-Pfad. — **Ausgang:** — (bei
  Closure)
- **Konfiguration wirkt anders, als sie gelesen wird** — ein Ausschluss-Muster für
  `cmd/` greift nicht, oder die Schwelle kommt als Build-Argument nicht in der Stufe an
  und ein Default gilt (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 1×).
  Je Ausschluss und für die Schwelle ein Fall der Gegenprobe. — **Ausgang:** — (bei
  Closure)
- **Laufzeit** — `-coverprofile` verlängert die Unit-Tests leicht und hebt den
  Build-Cache der Stufe `test` bei jeder Code-Änderung ohnehin auf; kein neuer Engpass
  erwartet. — **Ausgang:** — (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `8b399ba` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (8×, verkörpert in `AGENTS.md`
  §3.10) — trifft das Gate unmittelbar: Ein Schwellen-Gate, das nur grün gesehen
  wurde, belegt nicht, dass es unter der Schwelle rot wird; daher der zweite
  Liefer-Punkt mit dem Fall genau an der Schwelle.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (9×, verkörpert in §3.11) — die
  Sensors-Zeile sagt nur zu, was gemessen wird (etwa „Unit-Tests“, wenn die
  Integrationstests nicht mitzählen).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (6×, verkörpert in §3.12) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — darum stehen die
  Randformen in §6 offen und werden vor dem Code entschieden. Mit einem dritten
  Auftreten in diesem Slice erreichte der zweite Eintrag die Schwelle.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — dieser Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code ausdrücklich.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×) — Entscheidungsort der
  Randformen ist der Abschnitt für Harness-Werkzeuge der Spezifikation, angelegt von
  `slice-harness-vertraege-spezifikation`; die ADR trägt Entscheidung und Gründe.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×),
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
