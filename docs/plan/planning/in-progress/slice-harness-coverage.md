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
Nachweisart Gate). Die ADR des Gates ist [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md)
(Proposed, Empfehlung des Architect vom 2026-10-17, bis der Nutzer entschieden hat).

**Berührte Spec-Stellen:** `SPEC-056` (Schwelle mit Randformen; vom Architect am 2026-10-17 angelegt, vor dem Code); `SPEC-052` (Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest, und Allowlist des Build-Kontexts mit dem Prüfschritt; Architect, am 2026-10-17 nachgezogen; Punkt 2 und 4 gelten für den Namen `…-gegenprobe` ohne Änderung)

**Verantwortlich:** pt9912

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
Messung des heutigen Stands und steht nach Entscheidung des Nutzers vom 2026-10-06 in `SPEC-056` Punkt 1, als
Konstante mit Einheit und Begründung (Empfehlung des Architect: 94 Prozent der Anweisungen, gemessen 94,45 %
bei Quellstand `69aaf3f7`); [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md) hält Mechanismus, Messung mit
Quellstand und die Regel, nach der die Zahl folgt, und setzt `Schärft:` auf die
Kennung — nicht in diesem Plan. Gemessen und gezählt wird nach `SPEC-056` Punkt 2 bis 5: nur die
Unit-Tests der Stufe `test`, ein Gesamtwert nach Anweisungen, alle Pakete mit Anweisungen, ein Stand genau an der Schwelle ist grün.
Die Gegenprobe zeigt, dass das Gate rot werden kann (`AGENTS.md` §3.10), und trägt die
Abdeckungs-Deklaration von Teil 2 von LH-QA-07. Eine
spätere Senkung der Schwelle braucht eine neue ADR (`AGENTS.md` §3.6).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: eine Coverage-Schwelle vor
dem nächsten großen Slice. `harness/README.md` §Sensors führt Testabdeckung heute
unter „Nicht behauptet … (keine Schwelle)“. Das Schwester-Repo ai-harness-init hat
keine Schwelle; es ist Kontext, kein Vorbild.

**Übernommen aus `slice-harness-mutation`** (dort §1, Abgrenzung): die Anweisungs-Abdeckung
als Gate; das ist das Ziel oben.

**Sendung an die schnelle Kette** (`AGENTS.md` §3.13; Geber dieser Slice): Das Ziel der
Gegenprobe (Name nach §6 *Benennung*) endet auf `-gegenprobe` und hängt an `GATE_CHECKS`. Es
steht im selben Commit in `GATE_SCHNELL_AUSNAHME` (`harness/mk/gates-schnell.mk`; `SPEC-052`
Punkt 4, sonst ist `gates-schnell-gegenprobe` rot), und `SPEC-052` Punkt 13 trägt den Beleg, dass
die Gegenprobe `docs/user/` und `README.md` des Arbeitsbaums nicht liest (sie baut aus einer
Kopie unter eigenem Temp-Pfad, §6; ein `docker build` sieht nur den Build-Kontext der Allowlist
`.dockerignore`). Die Stufe `test` selbst ist kein Gate der Ausnahmeliste und bleibt in der
schnellen Kette (`SPEC-056` Punkt 8): Die Schwelle steht in `SPEC-056` und liegt nicht im Build-Kontext, das
Ziel `make test` liest sie auf dem Host; in die Allowlist von `.dockerignore` kommt nur der Prüfschritt
`tools/harness/coverage.sh` (`SPEC-052` Punkt 13 nennt ihn). Die Zeile `GATE_SCHNELL_AUSNAHME` nennt
`coverage-gegenprobe`; die Liste `LANGSAM` in `tools/harness/gates-schnell-gegenprobe.sh` führt es
ebenfalls (sonst ist `gates-schnell-gegenprobe` rot, `p4-eintrag`).

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
- Eine Schwelle je Paket und neue Tests für `bootstrap/tlsproxy` (49,1 %), `cmd/pgwire-recorder`,
  `bootstrap/offsetlese` und `testpki` (0 %) — Bestand bleibt bewusst stehen (Klasse 2): Eine
  Schwelle je Paket verlangte diese Tests vor dem Gate (Rückführung in §4); die Pakete zählen mit
  ihrem Wert im Gesamtwert und stehen in der Ausgabe (`SPEC-056` Punkt 3, 4 und 7).
- Die Integrationstests in der Messung (`GOCOVERDIR`) und `-coverpkg` — ein anderer Vorgang (Klasse 3):
  `LH-QA-07` verlangt die Unit-Tests, und ein Weg aus dem Container der Integrationstests ist eine eigene
  Entscheidung (`SPEC-056` Punkt 2, [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md) Option C und D).
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

- [ ] Gate mit Schwelle: Die Stufe `test` erzeugt das Profil, prüft es mit
      `tools/harness/coverage.sh` gegen die Schwelle aus `SPEC-056` Punkt 1
      (nach den Randformen aus §6; `make test` liest sie und gibt sie als Build-Argument ohne Vorgabewert) und endet rot
      darunter; `gofmt`, `go vet` und die Unit-Tests der Stufe bleiben unverändert. Die
      Messung des Stands, aus der die Schwelle folgt, steht mit Quellstand in der ADR;
      der Bestand ist grün.
- [ ] Gegenprobe `coverage-gegenprobe` als eigenes Ziel an `GATE_CHECKS` und in `GATE_SCHNELL_AUSNAHME`: ein Stand unter der Schwelle wird
      abgelehnt, einer genau an der Schwelle angenommen; dazu je Punkt und je Element der Aufzählungen von
      `SPEC-056` Punkt 4, 6 und 10 ein Fall mit einer Mutation, die nur ihn bricht (Paket ohne Tests, Profil,
      Schwelle, Rundung, Ausgabe, Stufe). Je Zusage ist die Mutation gesehen (`AGENTS.md` §3.10). Die
      Gegenprobe ist der Nachweis von Teil 2 (Anweisungsabdeckung) von LH-QA-07 und
      trägt nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) dessen Abdeckungs-Deklaration (Nachweisart Gate, Teil 2
      von 3); mit Teil 1 und 3 aus `slice-harness-abdeckung-gate` ist LH-QA-07 in den
      Abdeckungstabellen vollständig, `abdeckung-vollstaendig.md` und damit die RTM
      führen sie, und das Abnahmeszenario 17 ist nachweisbar (Tabellen mit
      `make abdeckung` neu geschrieben, `make abdeckung-check` grün).
- [ ] Doku: `harness/README.md` §Sensors führt die Schwelle beim Vertrag von
      `make test` (mit `SPEC-056`) und die Gegenprobe mit Bindung an [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md); die Zeile „Nicht
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
| `docs/plan/adr/0044-schwelle-der-anweisungsabdeckung.md`, `docs/plan/adr/README.md` | neu / update (Architect, erledigt) | ADR des Gates: Messung mit Quellstand, Regel der Zahl, `Schärft:` auf `SPEC-056`; Index-Zeilen; Status wird mit der Entscheidung des Nutzers `Accepted` |
| `harness/vertraege/SPEC-056-testabdeckung-schwelle.md`, `harness/vertraege/README.md` | neu / update (Architect, erledigt) | Schwelle als Konstante mit Einheit und Begründung und die Randformen aus §6 (Entscheidung des Nutzers vom 2026-10-06) |
| `Dockerfile` | update | Stufe `test`: `ARG` ohne Vorgabewert; `gofmt`, `go vet` und `go test -covermode=set -coverprofile` (ohne `-coverpkg`) bleiben eine Schicht, der Prüfschritt ist die nächste (die Schicht der Prüfung läuft bei geänderter Schwelle neu, die Tests nicht); netzlos |
| `.dockerignore` | update | Allowlist: `!tools/harness/coverage.sh` |
| `tools/harness/coverage.sh` | neu | Prüfschritt nach `SPEC-056` Punkt 3 bis 7 (POSIX-Shell und `awk` des Images der Stufe, keine bash) |
| `harness/mk/build.mk`, `harness/mk/coverage.mk` | update / neu | `test` liest die Schwelle aus `SPEC-056` (`$(shell sed)` auf die Zeile `**Schwelle:**` am Zeilenanfang, Muster `^\*\*Schwelle:\*\*` ohne führende Leerzeichen, ein Einzeiler als Rezept, Punkt 9); `coverage.mk` trägt das Ziel der Gegenprobe an `GATE_CHECKS` |
| `tools/harness/coverage-gegenprobe.sh` | neu | Fälle und Mutationen nach `SPEC-056` Punkt 10 (Profile ohne Docker; Kopien unter eigenem Temp-Pfad mit Docker); im Kopf die Abdeckungs-Deklaration von LH-QA-07, Teil 2 von 3 |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update | `coverage-gegenprobe` in `GATE_SCHNELL_AUSNAHME` (`SPEC-052` Punkt 4) und in der Liste `LANGSAM` der Gegenprobe der Kette (deren Kopie nimmt das echte Fragment und ein Modell von `GATE_CHECKS`), im Commit des Ziels |
| `harness/vertraege/SPEC-052-gate-kette.md` | update (Architect, erledigt) | Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest; Allowlist mit dem Prüfschritt |
| `docs/user/abdeckung-*.md` | update | mit `make abdeckung` neu geschrieben: LH-QA-07 vollständig |
| `harness/README.md` | update | §Sensors: Vertrag von `make test` um die Schwelle ergänzt, Gegenprobe als Zeile; „Nicht behauptet“ angepasst |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug`, die vier
Umstellungs-Slices (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`), die beiden
Bereinigungs-Slices (`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`),
`slice-harness-lint`, `slice-harness-integration-wait`, `slice-harness-meldungskatalog-gate`,
die Slices von welle-v1-abschluss und `slice-harness-abdeckung-gate` liegen in `done/`
(WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13, vor die Release-Slices gezogen), dann die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-integration-volumes`,
`slice-harness-abdeckung-gate` und `slice-harness-coverage`, dann die Slices von welle-erster-release und welle-v1-nachtrag, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-lint-warten-ohne-frist`,
`slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Dieser Slice ist der letzte
vor M3: Mit ihm wird Abnahmeszenario 17 nachweisbar, nach der Entscheidung vom
2026-10-08 vor welle-erster-release, damit M3 vor M4 erreicht ist. Nach ihm folgen
welle-erster-release und welle-v1-nachtrag, dann `slice-harness-commit-struktur-id`,
`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving` und
`slice-harness-mutation`; keiner der vier Harness-Slices trägt zu M3 bei, und keiner ist
Vorbedingung dieses Slice. Vom Abdeckungs-Gate hängt dieser Slice
technisch ab: Ohne es kennt `tools/test/abdeckung.sh` die Deklaration von Teil 2
nicht. Grund der übrigen Reihenfolge, keine technische Abhängigkeit: Lint-Bereinigung und
Black-Box-Umstellung ändern Code und Tests, an denen die Schwelle gemessen wird; eine
Messung davor wäre veraltet, bevor das Gate greift. Erster Schritt nach
dem Start, vor jedem Code-Commit: Der Architect misst den Stand (gesamt und je Paket,
mit Quellstand), entscheidet die Randformen aus §6 und schreibt die ADR mit der
Schwelle (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Größe** (`AGENTS.md` §3.14): Metrik ist die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- Dockerfile .dockerignore harness/mk
tools/harness tools/test harness/vertraege harness/README.md docs/user docs/plan/adr`, Tests
eingeschlossen. Basis: der Commit des Anspruchs (`next` → `in-progress`), `9ee6387d`, beim Anspruch
einzutragen. Schätzung des Architect vom 2026-10-17 (ersetzt den Arbeitswert 450 bis 750; Produktcode rund 130 Zeilen:
Prüfschritt 70, `Dockerfile`, Fragmente, `.dockerignore` 40, Handbuch der Doku 20; Gegenprobe mit dem Doppelten, rund 300; `SPEC-056` rund 110 Zeilen,
[ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md) rund 60 und `SPEC-052`, Indizes rund 20 aus den Commits des Architect, die in die Metrik fallen): 600 bis 800, Mitte 700.
**Meldeschwelle 750:** Wer im Lauf darüber liegt oder liegen wird, meldet es dem Planner, bevor er
weiterbaut; die Schätzung reicht bis 800, der Planner kennt das. **Marke 900:** darüber gilt die Rückführung unten, ohne Ermessen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe der Metrik liegt über 900,
  oder die Entscheidung über die
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

**Randformen des Vertrags** (`AGENTS.md` §3.12) — alle **entschieden** vom Architect am 2026-10-17,
vor dem ersten Code-Commit, in `SPEC-056` (Punkt in Klammern); [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md)
trägt Entscheidung und Gründe. Die Schwelle selbst ist eine Empfehlung, bis der Nutzer die ADR annimmt. Was dort nicht steht, entscheidet der
Implementer nicht, er gibt es zurück.

- **Die Zahl** (Punkt 1) — gemessen 94,45 % (4201 von 4448 Anweisungen, Quellstand `69aaf3f7`), Schwelle auf die ganze Zahl
  abgerundet: 94. Abstand 19 Anweisungen. Eine Anhebung folgt derselben Regel.
- **Stelle der Zeile `**Schwelle:**`** (Punkt 1) — entschieden am 2026-10-17: am echten Zeilenanfang, ohne Einrückung;
  `build.mk` liest mit `^\*\*Schwelle:\*\*`, die Gegenprobe (`^ *` an drei Stellen) zieht nach, ein Fall „eingerückte Zeile
  zählt nicht“.
- **Reihenfolge mehrerer Gründe** (Punkt 6) — entschieden am 2026-10-17: (a), (b), (c), (g), (h), dann (d), (e), (f); (d) und
  (e) nach Zeilenfolge. Fälle je benachbartem Paar und (d)/(e) vertauscht, je mit einer Mutation, die nur die Reihenfolge ändert.
- **Gesamt oder je Paket** (Punkt 3, 7) — ein Gesamtwert nach Anweisungen; keine Schwelle je Paket (sie verlangte neue Tests,
  §4); die Ausgabe nennt jedes Paket.
- **Integrationstests über `GOCOVERDIR`** (Punkt 2) — zählen nicht; nur die Unit-Tests der Stufe `test`.
- **`-coverpkg`** (Punkt 2, 10) — nicht gesetzt; ein Paket zählt nur seine eigenen Tests; Fall mit Mutation `-coverpkg=./...`.
- **`cmd/` und Pakete ohne Anweisungen** (Punkt 4) — `cmd/` zählt mit 0 %; ein Paket ohne Anweisung hat keinen Block und
  steht nicht in der Ausgabe; Blöcke mit Anzahl 0 zählen nichts.
- **Paket ohne Tests** (Punkt 4) — zählt mit 0 % im Gesamtwert, kein eigener Befund, auch ein Hilfspaket für Tests.
- **Generierter Code** (Punkt 4) — keine Regel; der Prüfschritt liest nur das Profil.
- **Vergleich und Rundung** (Punkt 5) — `>=` in ganzen Zahlen (`ausgeführt * 100 >= Schwelle * gesamt`), genau an der
  Schwelle grün, keine Rundung vor dem Vergleich, Anzeige abgerundet auf zwei Stellen.
- **Werkzeug und Ort der Schwelle** (Punkt 1, 8) — Bordmittel (`go test -coverprofile`, ein POSIX-Skript mit `awk`); die Schwelle steht
  in `SPEC-056`, `make test` liest sie auf dem Host und reicht sie als Build-Argument ohne Vorgabewert; sie liegt nicht im
  Build-Kontext, nur `tools/harness/coverage.sh` kommt in die Allowlist. Die Stufe `test` bleibt in der schnellen Kette.
- **`covermode`** (Punkt 2, 6) — `set`; ein anderer Modus im Profil ist abgelehnt.
- **Ausgabe** (Punkt 7) — je Paket eine Zeile bei grün und rot (aufsteigend nach Importpfad), zuletzt der Gesamtwert mit
  Schwelle, bei rot mit Differenz; die Ausgabe erscheint nur, wenn die Schicht läuft.
- **Eingabe des Prüfschritts** (Punkt 6) — Profil fehlt, leer, Kopfzeile, Zeilenform, doppelter Block, keine Anweisung, Schwelle
  fehlt oder ungültig: je Form rot.
- **Ausnahmeliste und Beleg der schnellen Kette** (Punkt 11, `SPEC-052` Punkt 13) — `coverage-gegenprobe` steht in der
  Ausnahmeliste; der Beleg steht in Punkt 13; die Allowlist nennt den Prüfschritt.
- **Aufrufarten des Treibers** (Punkt 9) — Rezept von `test` ein Einzeiler ohne `$(MAKE)`; `-n`, `-q`, `-t` und `MAKEFLAGS`
  je ein Fall, `-i` wie bei jedem Gate.
- **Benennung** — Ziel `coverage-gegenprobe`, Skripte `tools/harness/coverage.sh` und `coverage-gegenprobe.sh`, Fragment
  `harness/mk/coverage.mk`; „Coverage“ nennt die Anweisungsabdeckung, „Abdeckung“ bleibt die je Anforderung. Die Zeile in
  `harness/README.md` §Sensors sagt „Anweisungsabdeckung der Unit-Tests“.

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

### Belege des Implementers

Stand der Läufe: Quelle `7f1f6ead` plus die Commits dieses Slice bis zur Übergabe (Basis der Größenmetrik `9ee6387d`).

**Messung** (`make test`, Stufe `test`, `COVERAGE_SCHWELLE=94`): Gesamtwert 94,44 % (4201 von 4448 Anweisungen), angezeigt abgerundet; `SPEC-056` Punkt 1 nennt 94,45 % (gerundet, derselbe Bruch 4201/4448 = 94,4469 %). Je Paket wie in der ADR (`bootstrap/tlsproxy` 49,09 %, `cmd/pgwire-recorder`, `bootstrap/offsetlese` und `testpki` 0,00 %). Der Bestand ist grün, Abstand 19 Anweisungen.

**Mutationen.** Jede Mutation wurde in einer Kopie gefahren und färbt die genannten Fälle von `make coverage-gegenprobe` rot (Fälle mit Docker auch, die ohne Docker nicht). Weitere rote Fälle derselben Mutation sind nicht aufgeführt. Treiber: Textersatz in einer Kopie des Baums, danach `touch`, je Mutant frischer Pfad.

| Zusage (SPEC-056) | Mutation | roter Fall |
|---|---|---|
| 3 Gewicht nach Anweisungen, nicht nach Blöcken | `gesamt += 1` statt `+= $2` | `p3-gewicht-gruen`, `p3-gewicht-rot` |
| 3 nicht Mittel der Pakete | Vergleich auf den Mittelwert der Paketwerte | `p3-gewicht-gruen`, `p3-gewicht-rot` |
| 3 ausgeführt zählt die Anweisungen des Blocks, nicht den Zähler | `+= $3` statt `+= $2` | `p3-zaehler-ist-nicht-anzahl` |
| 4 (a) `main` ohne Test zählt | Blöcke unter `/cmd/` ausgenommen | `p4-main-ohne-test` |
| 4 (b) Hilfspaket ohne Test zählt | Pfad `testpki` ausgenommen | `p4-hilfspaket-ohne-test` |
| 4 (c) Block mit Anzahl 0 zählt nichts | Anzahl 0 zählt als 1 | `p4-anzahl-null` |
| 4 ein Paket ohne Test zählt nur seine Tests (Punkt 2) | `-coverpkg=./...` im `Dockerfile` | `p2-paket-ohne-test` |
| 5 genau an der Schwelle grün | `>` statt `>=` | `p5-genau` |
| 5 keine Rundung vor dem Vergleich | Vergleich auf `int(Wert + 0,5)` | `p5-eins-darunter` |
| 5 Anzeige abgerundet | `int(Wert + 0,5)` in anzeige | `p5-abgerundet` |
| 6 (a) Profil fehlt / nicht lesbar | Prüfung entfernt; nur ohne `-f`; nur ohne `-r` | `p6-a-fehlt`, `p6-a-ohne-argumente`; `p6-a-verzeichnis`; `p6-a-nicht-lesbar` |
| 6 (b) Profil leer | `-s` entfernt | `p6-b-leer` |
| 6 (c) `mode: set` | Prüfung entfernt; nur `mode:`; nur Präfix; erste Zeile mit `mode` statt Zeile 1 | `p6-c-*`; `p6-c-count`, `p6-c-atomic`; `p6-c-zusatz`; `p6-c-leere-erste-zeile` |
| 6 (d) Form je Element | je Element der Aufzählung ein Mutant, der nur es lockert: Datei leer, Bereich leer, ohne `:`, Anzahl Text/Dezimal/negativ, Zähler Text/Dezimal/negativ, ein/zwei/vier Felder, ohne `^`, `$` gelockert, doppeltes Leerzeichen, Tab als Trenner, Leerzeile übersprungen | `p6-d-<Element>` je ein Fall |
| 6 (d) jede Zeile geprüft | nur Zeile 2; erste Zeile übersprungen; letzte übersprungen | `p6-d-pos-letzte`; `p6-d-pos-erste`; `p6-d-pos-letzte` |
| 6 Reihenfolge (a) vor (b) | Prüfung (b) vor (a) | `p6-r-a-b` |
| 6 Reihenfolge (b) vor (c) | Prüfung (c) vor (b) | `p6-r-b-c` |
| 6 Reihenfolge (c) vor (g) | Prüfungen (g), (h) vor (c) | `p6-r-c-g` |
| 6 Reihenfolge (g) vor (h) | Prüfung (h) vor (g) | `p6-r-g-h` |
| 6 Reihenfolge (h) vor (d) | Zeilenprüfung (d) vor den Prüfungen der Schwelle | `p6-r-h-d` |
| 6 Reihenfolge (d) vor (e) in einer Zeile | (e) vor (d) | `p6-r-d-e-zeile` |
| 6 (d) und (e) nach Zeilenfolge, nicht (d) zuerst | (e) bis zum Ende zurückgestellt, vor (f) gemeldet | `p6-r-d-e-folge` |
| 6 Reihenfolge (e) vor (f) | (e) bis zum Ende zurückgestellt, nach (f) gemeldet | `p6-r-e-f` |
| 6 (e) Block = Datei und Bereich, zweimal | Prüfung entfernt; Schlüssel nur Datei; nur Bereich; nur benachbart; ganze Zeile | `p6-e-doppelt`; `p6-e-anderer-bereich`; `p6-e-andere-datei`; `p6-e-getrennt`; `p6-e-anderer-zaehler`, `p6-e-andere-anzahl` |
| 6 (f) Summe 0 | Prüfung entfernt; nur Kopfzeile | `p6-f-nur-kopf`, `p6-f-nur-null-bloecke`; `p6-f-nur-null-bloecke` |
| 6 (g) Schwelle fehlt/leer | `-n`-Prüfung entfernt | `p6-g-fehlt`, `p6-g-leer` |
| 6 (h) Schwelle 1 bis 100 | je akzeptierter Wert ein Mutant (`0`, `101`, `94.5`, `094`, `94 94`, `-1`, `+94`, `abc`, `" 94"`, `"94 "`, `1e2`, `00`, `1000`, Zeilenumbruch); `100`, einstellig, zweistellig entfernt; alles akzeptiert | `p6-h-[<Wert>]`, `p6-h-zeilenumbruch`; `p6-h-grenze-100`; `p6-h-grenze-1`/`-9`; `p6-h-grenze-10` |
| 6 Ablehnung nennt die Form und schreibt nichts auf stdout | Meldung auf stdout (Shell, awk) | `p6-a-fehlt`, `p6-g-fehlt`, `p6-d-*`, `p6-e-doppelt` |
| 7 Liste aufsteigend nach Importpfad | ohne `sort`; `sort -r` | `p7-gruen-ausgabe`, `p7-rot-ausgabe` |
| 7 Gesamtzeile zuletzt | `close(sortiert)` entfernt (awk des Images) | `p7-gesamtzeile-zuletzt-in-der-stufe` |
| 7 Differenz nur bei rot | immer; nie; um 0,01 falsch | `p7-gruen-ausgabe`; `p7-rot-ausgabe`, `p5-eins-darunter`; `p7-rot-ausgabe`, `p5-eins-darunter` |
| 7 Liste bei grün und rot | keine Liste bei rot; keine bei grün | `p7-rot-ausgabe`, `p4-main-ohne-test`; `p7-gruen-ausgabe`, `p4-anzahl-null` |
| 7 Exit 1 nur bei rot | Exit immer 0; immer 1 | `p5-eins-darunter`, `p7-rot-ausgabe`; `p5-genau`, `p7-gruen-ausgabe` |
| 1 Wert aus der Datei, nicht fest | Zahl 94 im Rezept; Zahl im `Dockerfile` statt Argument | `p1-wert-aus-datei`; `p1-schwelle-kommt-an`, `p1-zeile-fehlt-stufe` |
| 1 Schwelle kommt in der Stufe an (Mutation des Vertrags) | `ARG COVERAGE_SCHWELLE=1` und kein `--build-arg` | `p1-schwelle-kommt-an`, `p1-wert-aus-datei` |
| 1 zwei Zeilen ergeben `94 94`; eine fehlende einen leeren Wert | `sed ...;q`; Vorgabe 94 bei leerem Wert; Argument ohne Anführungszeichen | `p1-zeile-doppelt`; `p1-zeile-fehlt`, `p1-zeile-fehlt-stufe`; `p1-zeile-doppelt` |
| 1 nur eine Zeile am Zeilenanfang (eine eingerückte Zeile zählt nicht) | `sed` ohne Anker `^`; mit `^ *` (Einrückung zugelassen); mit `-` zugelassen | `p1-zeile-mitten-im-text`; `p1-zeile-eingerueckt`; `p1-zeile-mitten-im-text` |
| 8 `gofmt` bricht die Stufe | `gofmt` entfernt | `p8-format` |
| 8 `go vet` über den Integrationstests | `go vet ./...` ohne `-tags integration` | `p8-vet-integration` |
| 8 ein roter Test bricht die Stufe | `go test ... || true` | `p8-test-rot` |
| 8 der Prüfschritt bricht die Stufe | `|| true` hinter dem Prüfschritt; Prüfschritt entfernt | `p1-schwelle-kommt-an`, `p1-zeile-fehlt-stufe`, `p2-paket-ohne-test` |
| 8 Allowlist nennt den Prüfschritt | Zeile aus `.dockerignore` entfernt | `p8-grundlauf` und alle Fälle der Stufe |
| 6/8 `mode: set` | `-covermode=atomic` | `p8-grundlauf` |
| 9 Rezept ohne `$(MAKE)` | `+` vor dem Rezept; `# $(MAKE)` im Rezept | `p9-flag-n`, `p9-flag-q`, `p9-flag-t`, `p9-makeflags-n`, `-q`, `-t` (alle sechs) |
| 11 Anschluss an `GATE_CHECKS` und Ausnahmeliste | Zeile aus `coverage.mk`; Eintrag aus `GATE_SCHNELL_AUSNAHME` | `p11-gate-checks`; `p11-ausnahmeliste` |

**Grüne Mutanten, eingeordnet.** (1) Vergleich auf den angezeigten (abgerundeten) Wert statt auf das Verhältnis: äquivalent, denn Anzeige unter der Schwelle ⇔ Verhältnis unter der Schwelle (Beweis: `floor(x) < n` für ganzes `n` ⇔ `x < n`); die Spezifikation sagt es in Punkt 5. (2) `close(sortiert)` entfernt: auf dem awk des Hosts (gawk) äquivalent, es schließt Pipes vor dem Leeren von stdout; auf dem awk des Images rot (`p7-gesamtzeile-zuletzt-in-der-stufe`). (3) `exit 0` statt `exit 1` in der Fehlerfunktion von awk: äquivalent für (d) und (e), weil `END` bei gesetztem Merker mit 1 endet; für (f) rot (`p6-f-nur-kopf`, `p6-f-nur-null-bloecke`). (4) `-covermode=set` ohne Angabe: `go test -coverprofile` wählt `set`, solange `-race` fehlt; nicht gefahren. Offen, ohne Fall (Kopf der Gegenprobe, Abschnitt OFFEN): Punkt 4 (d), Punkt 7 Ausgabe nur bei laufender Schicht (Cache), Punkt 8 netzlos und nichts in den Arbeitsbaum, Punkt 9 `-i`, Aufräumen der Kopien und Image-Tags.

**Läufe.** `make test` grün (94,44 % bei Schwelle 94); `make coverage-gegenprobe` grün; `make gates-schnell-gegenprobe` grün (Liste `LANGSAM` um `coverage-gegenprobe` ergänzt); `make abdeckung`, danach `make abdeckung-check`, `make docs-check`, `make kopf-check` grün; `make gates` einmal am Endstand (siehe Bericht).

**Größe** (`AGENTS.md` §3.14, Metrik §4): wird mit der Übergabe aus `git diff --numstat 9ee6387d..<letzter Commit>` genannt; die Meldeschwelle 750 war vor der Übergabe überschritten und ist dem Planner gemeldet.
 **Beantwortet am 2026-10-17 vom Architect:** beide in `SPEC-056` entschieden (Punkt 1 am echten Zeilenanfang, Punkt 6 Reihenfolge); §6 oben nennt sie.
**Rückgaben an den Architect** (Randform · Frage): `SPEC-056` Punkt 1 sagt „Zeilenanfang“, die Zeile steht in der Datei mit drei Leerzeichen als Fortsetzung eines Listenpunkts; `make test` liest sie mit führenden Leerzeichen (`^ *`), sonst fände es keine Zeile. Frage: gilt das so, oder wird die Zeile in der Spezifikation anders gesetzt. Ferner: nennt Punkt 6 keine Reihenfolge mehrerer Gründe in einer Eingabe (der Prüfschritt meldet in der Folge (a), (b), (c), (g), (h), dann (d), (e), (f)); jeder Fall der Gegenprobe hat genau einen Grund.

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
  Randformen ist der Ordner `harness/vertraege/`, angelegt von
  `slice-harness-vertraege-umzug-ordner`; die ADR trägt Entscheidung und Gründe.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×),
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Nachgezählt** mit der Sendung aus `slice-harness-vertraege-umzug-gate-kette` (Ort neuer
Kennungen: `harness/vertraege/`, nächste freie Nummer damals 053, heute 056 nach `SPEC-053` bis `SPEC-055`; `AGENTS.md` §3.13): drei
Liefer-Punkte (§2: Gate mit Schwelle, Gegenprobe, Doku, unverändert) und zwei Schichten,
Harness (Vertragsdatei, `Dockerfile`, `harness/mk/`, `tools/`, `harness/README.md`) und
Nutzer-Doku (`docs/user/abdeckung-*.md`); die ADR ist Planung. Die Vertragsdatei lag vorher
in der Spezifikation; die Sendung erhöht die Zahl der Schichten nicht.

**Nachgezählt** nach der Entscheidung der Randformen (Stand `9f56eeec`): Liefer-Punkte 3 (Gate mit Schwelle, Gegenprobe, Doku),
Schichten 2 (Harness mit `Dockerfile`, `.dockerignore`, `harness/mk/`, `tools/harness/`; Nutzer-Doku `docs/user/abdeckung-*.md`);
`Dockerfile`, Prüfschritt und Fragmente liegen in der Harness-Schicht, `SPEC-056`, [ADR-0044](../../adr/0044-schwelle-der-anweisungsabdeckung.md) und die Pläne sind Planung. Die Größe
steht in §4 (600 bis 800); die Meldeschwelle 750 liegt in der Schätzung.

**Nachgezählt** mit der Sendung an `gates-schnell.mk` und `SPEC-052` (§1; Stand `91a1bd02`,
`AGENTS.md` §3.13): Liefer-Punkte 3 (die Sendung wächst in Punkt 2, Gegenprobe, hinein), Schichten
2 (Harness mit `harness/mk/gates-schnell.mk` und `harness/vertraege/SPEC-052-gate-kette.md`;
Nutzer-Doku `docs/user/abdeckung-*.md`); `Dockerfile` und `tools/harness` liegen in der
Harness-Schicht, die Pläne und ADR sind Planung.

**Nachgezählt** mit dem Nachzug nach der Entscheidung des Architect vom 2026-10-17 (Zeile am Zeilenanfang, Reihenfolge der Gründe): Liefer-Punkte 3 (Gate mit Schwelle, Gegenprobe, Doku, unverändert; die neuen Fälle wachsen in Punkt 2 hinein), Schichten 2 (Harness mit `harness/mk/build.mk` und `tools/harness/coverage-gegenprobe.sh`; Nutzer-Doku `docs/user/abdeckung-*.md`); Größe der Metrik aus §4 siehe §7.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
