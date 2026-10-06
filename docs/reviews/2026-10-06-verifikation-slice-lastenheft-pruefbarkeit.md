# Verifikation: slice-lastenheft-pruefbarkeit — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-lastenheft-pruefbarkeit.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `27b33fe..b098fb3`. Darin: `a238108` ([ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) Proposed, Randformen), `daef83b` ([`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) und Abnahmeszenario 17, ADR Accepted), `dccfa8e` (`docs/plan/planning/open/slice-harness-abdeckung-gate.md`, Reihenfolge, Drift-Log), `9d85fa4` (Köpfe und Nachweis-DoDs der Folge-Slices), `9da7185` (Review-Report), die Nacharbeit `b9eae96` und `db0c6e4`, die Adaption [MR-001](../../harness/conventions.md#mr-001) von `c2b8d32` bis zur Auflösung in `6ca20bc`, `566782c` und `3dfef35`, und `8e98bc4`, `5809494` (neuer Slice `docs/plan/planning/open/slice-harness-vertraege-spezifikation.md`, Ort der Randformen in den Folge-Slices) sowie `b098fb3` (Randform *Spezifikation* im Plan).

**Eingang:**

- die DoD-Liefer-Punkte, die Entscheidungen des Nutzers im Plan (§1 *Herkunft*, §6 *Randformen*) und die Commit-Messages
- der [Review-Report](2026-10-06-review-slice-lastenheft-pruefbarkeit.md) (F-406 bis F-416, Stand `9d85fa4`); die Nacharbeit ab `b9eae96` hat kein eigenes Review gesehen
- `spec/lastenheft.md` §4 und §7, [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md), [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md), ADR-Index, `docs/plan/planning/in-progress/roadmap.md`
- die Pläne `docs/plan/planning/open/slice-harness-lint.md`, `docs/plan/planning/open/slice-harness-coverage.md`, `docs/plan/planning/open/slice-harness-mutation.md`, `docs/plan/planning/open/slice-harness-blackbox-kern.md`, `docs/plan/planning/open/slice-harness-blackbox-driven.md`, `docs/plan/planning/open/slice-harness-blackbox-pgwire.md`, `docs/plan/planning/open/slice-harness-blackbox-einstieg.md`, `docs/plan/planning/open/slice-harness-abdeckung-gate.md`, `docs/plan/planning/open/slice-harness-vertraege-spezifikation.md`
- `harness/conventions/done/MR-001-spezifikations-ort-werkzeugvertraege.md`, `AGENTS.md` §3.9 und §3.12, `.d-check.yml`, `tools/harness/kopf-check.sh` (Kopf)

Bei meinem Start war der Arbeitsbaum sauber, HEAD `b098fb3`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg.

**Wie die Proben gebaut wurden.** Der Slice liefert keinen Code; „bewusst brechen“ heißt hier: zeigen, dass die Sensoren, auf die sich DoD und Plan berufen, die Zusage halten. Jede Probe lief in einer frischen Kopie aus `git archive b098fb3` unter `verify-lh/` im Scratchpad, mit einer `sed`-Ersetzung an genau einer Stelle; die Wirkung der Ersetzung habe ich je Probe nachgelesen. Die unveränderte Kopie ist grün (`make kopf-check` Exit 0, `make docs-check` 0 Befunde). Die Kopien habe ich danach gelöscht.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — Lastenheft. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| §4 führt die Anforderung unter der nächsten freien QA-Kennung, Überschrift `### LH-QA-07 — Prüfbarkeit des Quellcodes`, nach LH-QA-06 und vor LH-RB-01 | `spec/lastenheft.md:820`, `:829`, `:843` | bestätigt |
| Wortlaut Byte für Byte gleich mit Plan §6, „muss“/„müssen“ | Zitat aus Plan §6 (ohne `> `-Präfix) gegen `spec/lastenheft.md:831`–`:841` per `diff`: leer. Im Lastenheft-Text steht kein „soll“/„sollen“ | bestätigt |
| Priorität MUSS allein über den Wortlaut, kein Feld `Priorität` | Feldliste von LH-QA-07: nur **Anforderung** und **Messmethode** | bestätigt |
| §7 führt Abnahmeszenario 17 im vom Nutzer bestätigten Text | Zitat aus Plan §6 *Abnahmeszenario* gegen `spec/lastenheft.md:1019`–`:1025` per `diff`: leer; Überschrift `### Abnahmeszenario 17 — Prüfbarkeit des Quellcodes` nach Szenario 16, `Bezug: LH-QA-07.` blank wie in 1 bis 16 | bestätigt |
| Referenz-Richtung: Szenario 17 und LH-QA-07 nennen kein Werkzeug, Gate, keine ADR, keinen Slice | Lesend: kein Treffer auf Werkzeug- oder Artefaktnamen. Proben PR1 bis PR3 (Abschnitt 3) rot | bestätigt; Grenze siehe PR4 |
| Version und Historie nach §6 *Version* (Draft 0.1.0, kein Sprung, keine Historie-Zeile) | `git diff 27b33fe..b098fb3 -- spec/lastenheft.md` hat genau zwei Hunks (§4, §7); Kopf und §8 unberührt | bestätigt |

### Punkt 2 — Roadmap und Köpfe. **Bestätigt**, mit einer Anmerkung zur Reichweite des Sensors (V-64).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Trigger von M3 nennt Abnahmeszenario 17 | `docs/plan/planning/in-progress/roadmap.md:46` „Abnahmeszenarien 1 bis 10 und 12 bis 17“; Drift-Log-Zeile `:102` | bestätigt |
| `Bezug` von Lint, Coverage und den vier Umstellungs-Slices verlinkt [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) | je Zeile 19 der sechs Pläne, Anker `#lh-qa-07--prüfbarkeit-des-quellcodes` | bestätigt |
| Platzhalter-Sätze entfernt | `grep` auf „hier nach“, „trägt ihre Kennung“, „neue(n) Anforderung“ in den sechs Plänen: leer; die entfernten Sätze stehen im Diff | bestätigt |
| `make kopf-check` und `make docs-check` grün | selbst gefahren: Exit 0 bzw. 250 Dateien, 0 Befunde; ebenso in `make gates` | bestätigt |
| Der Sensor hält die Zusage | Proben PK1 bis PK7 und PD1, PD3 rot (Abschnitt 3) | bestätigt; die Link-Form selbst hält kein Sensor (PD2, V-64) |

### Punkt 3 — Nachweis-Weg. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Entschieden vor dem Lastenheft-Commit | `a238108` (Randformen, ADR Proposed) liegt vor `daef83b` | bestätigt |
| [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) ergänzt [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), keine Ablösung | ADR §Entscheidung; Index-Zeile unter „Ergänzende ADRs“ (`docs/plan/adr/README.md:66`); [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) im Bereich unverändert | bestätigt |
| `Accepted` mit dem Lastenheft-Commit, Kennung im `Bezug` | `daef83b` ändert ADR (Status, `Bezug` mit Link auf LH-QA-07, Geschichte, Schlusssatz), Index und Lastenheft gemeinsam; danach kein Commit an der ADR (`git log 27b33fe..b098fb3 -- docs/plan/adr/0033-*`: nur `a238108`, `daef83b`) | bestätigt |
| `slice-harness-abdeckung-gate` liegt in `open/` | `docs/plan/planning/open/slice-harness-abdeckung-gate.md`, `Verantwortlich: —` | bestätigt |
| Lint, Blackbox-Einstieg und Coverage nennen in der DoD ihren Teil des Nachweises | Lint: Teil 1, mit den Umstellungs-Slices Teil 3, ohne Deklaration; Einstieg: Teil 3 vollständig, Deklaration in abdeckung-gate; Coverage: Teil 2 von 3 deklariert, danach vollständig | bestätigt |

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| `make gates` grün | bestätigt, selbst gefahren auf `b098fb3` (Abschnitt 5) |
| Review-Report liegt vor | bestätigt: `9da7185`, F-406 bis F-416 |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen, wie erwartet vor der Closure. §7 trägt noch die Platzhalter; beide Risiken aus §6 stehen mit „Ausgang: — (bei Closure)“. Das Register `BEO-REPO/harness-lesart-ohne-entscheidungsort` hat 1 Datei unter `evidence/`; der Beleg für F-406 (zweites Auftreten) folgt laut §8 mit der Closure. |

---

## 2. Review-Findings

| Finding | Status | Beleg |
|---|---|---|
| F-406 (MEDIUM) Ort der Randformen | behoben. Ort ist der Abschnitt für Harness-Werkzeuge der Spezifikation (Antwort des Kurs-Repos), angelegt von `slice-harness-vertraege-spezifikation`; [MR-001](../../harness/conventions.md#mr-001) eingeführt und aufgelöst, Eintrag in `conventions/done/`, Index unter *Aufgelöste Adaptionen*, *Aktive* wieder leer; `AGENTS.md` §3.12 nennt die Spezifikation auch für Harness-Werkzeuge; Plan §6 *Ort*, §6 *Spezifikation*, §8 und abdeckung-gate §6 *Ort* zeigen dorthin; die übrigen Folge-Slices ziehen in `5809494` nach. Rest: der Satz in [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md), der den Skriptkopf als Ort nennt, steht noch (V-65). | `b9eae96`, `db0c6e4`, `6ca20bc`, `566782c`, `3dfef35`, `8e98bc4`, `5809494`, `b098fb3` |
| F-407 (MEDIUM) zweite Bedingung von Messmethode 3 | adressiert im Plan. Lint-DoD und §3 verlangen je einen Rot-Fall „Zugriff an der Brücke vorbei“ und „Test in der Brückendatei“; Einstieg-DoD und das Risiko in abdeckung-gate nennen beide Bedingungen. Ob die Fälle tragen, zeigt erst der Code. | `b9eae96` |
| F-408 (MEDIUM) Randformen-Liste | behoben: §6 von abdeckung-gate nennt die drei fehlenden (geteilte Messung an Go-Tests, Deklaration außerhalb des Musters, Zahlform). | `b9eae96` |
| F-409 (LOW) Liefer-Punkte | behoben: drei Liefer-Punkte (Werkzeug und Vertrag, Gegenprobe, erste Deklarationen); die Doku-Zeile ist in die ersten beiden eingegangen. | `b9eae96` |
| F-410 (LOW) Kopf-Kennung ohne Gegenstück | behoben: kern, driven, pgwire nennen LH-QA-07 in §1 und §2; Probe PK3 bis PK5 zeigt, dass `make kopf-check` den Kopf jetzt hält. | `b9eae96`; Abschnitt 3 |
| F-411 (LOW) Ausweg im Coverage-Start | behoben: „ausdrücklich zurückgestellt“ gestrichen. | `b9eae96` |
| F-412 (LOW) Herkunft | behoben: Team für „nicht unbelegt“, Architect für den Vorschlag, Nutzer für Annahme der ADR und Platz von abdeckung-gate, gleichlautend in Plan und abdeckung-gate. | `b9eae96` |
| F-413 (LOW) „Messmethode“ doppelt belegt in der ADR | nicht adressiert, kein Ausgang notiert (V-65). | — |
| F-414 (LOW) Ende der Reihe | behoben (Kopf und §4 nennen `slice-harness-mutation`); durch den neuen Slice ist §4 aber erneut unvollständig (V-62). | `b9eae96` |
| F-415, F-416 (INFO) | keine Handlung nötig; F-415 ist Frage an Validator und Nutzer. | — |

---

## 3. Proben

| Probe | Mutation (je in eigener Kopie) | Erwartet | Ergebnis |
|---|---|---|---|
| PK1–PK7 | Link auf LH-QA-07 im `Bezug` durch `—` ersetzt, je einzeln in lint, coverage, blackbox-kern, blackbox-driven, blackbox-pgwire, blackbox-einstieg, abdeckung-gate | `make kopf-check` rot | rot in allen sieben: „§1: LH-QA-07 fehlt im Kopf“ und „§2: …“, Skript Exit 1 (make 2) |
| PK8, PK9 | dasselbe im eigenen Plan und in `slice-harness-vertraege-spezifikation` | — (dokumentiert) | grün. §1 und §2 beider Pläne nennen die Kennung nicht; der Sensor prüft nur die Richtung §1/§2 → Kopf ([ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md)). Siehe V-66. |
| PD1 | Anker im Lint-Kopf `#lh-qa-07--pruefbarkeit-des-quellcodes` (ohne Umlaut) | `make docs-check` rot | rot, 1 Befund `anchor-missing` |
| PD2 | Link im Coverage-Kopf durch blankes `` `LH-QA-07` `` ersetzt | — (Reichweite) | grün, 0 Befunde: `ids` in `.d-check.yml` verlangt den Link nur für ADR-Kennungen (V-64) |
| PD3 | Überschrift im Lastenheft zu „Pruefbarkeit“ geändert | `make docs-check` rot | rot, 14 Befunde `anchor-missing` (alle Verweise auf den Anker) |
| PR1 | Szenario 17: „geprüft von slice-harness-lint“ angehängt | rot | rot, `matrix-forbidden` spec-straten → slice |
| PR2 | Szenario 17: die Kennung der ADR 0033 blank, ohne Link, im `Bezug` | rot | rot, `id-unlinked` (verlinkt wäre es `matrix-forbidden` spec-straten → adr) |
| PR3 | Szenario 17: Link auf `tools/harness/kopf-check-gegenprobe.sh` | rot | rot, `matrix-forbidden` spec-straten → aussen |
| PR4 | Szenario 17: „(make lint, golangci-lint)“ als blanker Text | — (Reichweite) | grün. Werkzeugnamen als Text hält kein Sensor; dass keiner dasteht, ist nur gelesen (Punkt 1) |

---

## 4. Plan gegen Stand

| Abschnitt | Ergebnis |
|---|---|
| Kopf | folgt: `Bezug` mit LH-QA-07, [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) Accepted, [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); `Berührte Spec-Stellen` `lastenheft.md §4` · `§7` stimmt mit dem Diff; `Welle:` nennt Coverage als letzten vor M3 und Mutation als Ende der Reihe. |
| §1 | folgt im Ziel. Abgrenzung nennt den neuen Folge-Slice nicht, und die Schicht-Abgrenzung nennt nicht alles, was geändert wurde (V-62, V-63); eine Begründung in (3) beruft sich auf eine Sensor-Eigenschaft, die es nicht gibt (V-64). |
| §3 | folgt für Lastenheft, Folge-Pläne, abdeckung-gate, Mutation, Roadmap, ADR. `AGENTS.md` fehlt (V-63). |
| §6 | folgt: Randformen entschieden, *Ort* und *Spezifikation* nach der Kurs-Antwort nachgezogen (`db0c6e4`, `b098fb3`); die Rückführung aus §4 ist begründet nicht ausgelöst. |
| §8 | folgt: Treffer `BEO-REPO/harness-lesart-ohne-entscheidungsort` mit Ort Spezifikation; Register-Stand `ce50a10` als Datum der Sichtung. |
| Reihe und Roadmap | Roadmap-Drift-Log `:104` und `:105` ergeben zusammen: Lastenheft, vertraege-spezifikation, Lint, vier Umstellungs-Slices, abdeckung-gate, Coverage, Mutation. Die Start-Trigger bilden diese Kette lückenlos ab: vertraege-spezifikation ← dieser Slice; Lint ← vertraege-spezifikation; kern ← Lint; driven ← kern; pgwire ← driven; einstieg ← pgwire; abdeckung-gate ← einstieg; Coverage ← alle davor; Mutation ← Coverage. Die Aufzählungen in den Start-Absätzen sind verschieden lang; nur der eigene Plan und die Kurzformen lassen vertraege-spezifikation aus (V-62). |

---

## 5. Gate-Lauf

`make gates` auf `b098fb3`, selbst gefahren: Exit 0, Laufzeit 60 s. Grün gemeldet: `abdeckung-gegenprobe`, `a-check` (0 Befunde), `a-check-negativ`, `baseline-verify` (v6.13.0, 54 Dateien), `build`, `test` (gofmt leer, vet, Unit-Tests), `docs-check` (250 Dateien, 0 Befunde), `hook-gegenprobe`, `test-integration`, `kopf-check-gegenprobe`; `kopf-check` und `abdeckung-check` ohne Ausgabe, Exit 0; der Nachweis `record-gates` lief zuletzt. Danach war der Arbeitsbaum sauber. Die Commit-Messages im Bereich nennen keine Struktur-ID (`SPEC-*`, `ARC-*`).

---

## 6. Befunde

| ID | Befund | Ort | Empfehlung |
|---|---|---|---|
| V-62 | **§4 *Start* und §1 folgen dem neuen Slice nicht.** §4 zählt die Reihe als „Lastenheft, Lint, die vier Umstellungs-Slices, `slice-harness-abdeckung-gate`, Coverage, `slice-harness-mutation`“; seit `8e98bc4` steht `slice-harness-vertraege-spezifikation` zwischen diesem Slice und Lint (Roadmap `:105`, Lint §4, vertraege-spezifikation §4). §1 *Ausdrücklich NICHT* nennt die Werkzeug-Verträge in der Spezifikation nicht als Ausschluss mit Folge-Slice-Kennung; §6 *Ort dieser Randformen* und §8 sagen „ein eigener Folge-Slice“ ohne Kennung, nur §6 *Spezifikation* nennt sie. `AGENTS.md` §3.9, `BEO-REPO/plan-folgt-korrektur-nicht`. Kurzformen in den Start-Absätzen von blackbox-kern, -driven, -pgwire, -einstieg, Coverage und Mutation lassen den Slice ebenfalls aus; sie widersprechen der Kette nicht, weil Lint auf ihn wartet. | `docs/plan/planning/in-progress/slice-lastenheft-pruefbarkeit.md` §1, §4 *Start*, §6 *Ort dieser Randformen*, §8 | §4 um den Slice ergänzen; in §1 einen Ausschluss „Verträge der Gates in der Spezifikation — übernimmt `slice-harness-vertraege-spezifikation`“; in §6 und §8 die Kennung statt „eigener Folge-Slice“. Vor der Closure. |
| V-63 | **§1 und §3 nennen die Änderung an `AGENTS.md` §3.12 nicht.** `c2b8d32` trägt die Kennung dieses Slice und ändert `AGENTS.md`, `harness/conventions.md` und legt den MR-Eintrag an; die Endfassung von §3.12 kam mit `6ca20bc` unter [MR-001](../../harness/conventions.md#mr-001). Die Schicht-Abgrenzung in §1 zählt auf, was der Slice ändert, und `AGENTS.md` steht nicht darin. | `docs/plan/planning/in-progress/slice-lastenheft-pruefbarkeit.md` §1 *Produkt-Code und Harness-Werkzeuge*, §3 | Entweder §3 um `AGENTS.md` §3.12 ergänzen oder in §7 festhalten, dass Einführung und Auflösung von [MR-001](../../harness/conventions.md#mr-001) ein eigener Vorgang aus F-406 waren. |
| V-64 | **Die Begründung in §1 (3) beruft sich auf eine Sensor-Eigenschaft, die es nicht gibt.** „Das Doku-Gate verlangt für jede Lastenheft-Kennung einen auflösenden Link“: `.d-check.yml` führt unter `ids` nur das Muster für ADR-Kennungen. Probe PD2 (Link durch blanke Kennung ersetzt) bleibt grün. Die Reihenfolge-Begründung trägt trotzdem, weil ein Link auf einen fehlenden Anker rot wird (PD1, PD3). Für DoD 2 heißt das: Dass die Köpfe *verlinken*, ist gelesen bestätigt, nicht von einem Sensor gehalten. Der Satz stammt aus `20f9044`, also vor diesem Diff. `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. | `docs/plan/planning/in-progress/slice-lastenheft-pruefbarkeit.md` §1 *Warum ein eigener Slice* (3) | Satz auf das berichtigen, was gilt: Ein Link auf die Kennung löst erst auf, wenn die Überschrift im Lastenheft steht. Kein neuer Sensor nötig. |
| V-65 | **[ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) trägt zwei Textstellen ohne notierten Ausgang.** (a) „Form der Deklaration, Pfad-Schreibweise und Fehlformen stehen im Kopf des Abdeckungs-Skripts“ widerspricht dem Ort nach der Kurs-Antwort. Die Spezifikation steht im Rang über der ADR, der Widerspruch löst sich also zu ihren Gunsten auf. §6 *Lesart-Rangfolge* von `slice-harness-vertraege-spezifikation` nennt aber nur [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md), und §6 *Schärft-Bezug* nur das `Schärft: —`. (b) F-413 („Messmethode“ in zwei Bedeutungen, Zuschreibung von „genau ein Pfad *Messung*“ an [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)) ist in keinem Plan aufgegriffen. Beides ist nach §3.5 nicht im ADR-Text zu beheben. | `docs/plan/adr/0033-gate-nachweise-in-der-abdeckung.md:31`, `:19`–`:21`, `:28`; `docs/plan/planning/open/slice-harness-vertraege-spezifikation.md` §6 | Bei der Closure je einen Ausgang notieren: (a) als Randform in §6 von vertraege-spezifikation (der Abschnitt gilt, der Satz der ADR ist überholt), (b) als *entfallen* mit Begründung oder als Begriffsklärung im neuen Spezifikations-Abschnitt. |
| V-66 | **Hinweis: Der Kopf dieses Plans und der von `slice-harness-vertraege-spezifikation` werden von keinem Sensor gehalten.** Die Kennung LH-QA-07 aus dem `Bezug` zu entfernen, lässt `make kopf-check` grün (PK8, PK9), weil §1 und §2 beider Pläne die Kennung nicht nennen. Das ist die Grenze aus [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md), dieselbe Klasse wie F-410, und kein Verstoß. | `docs/plan/planning/in-progress/slice-lastenheft-pruefbarkeit.md` §1, §2; `docs/plan/planning/open/slice-harness-vertraege-spezifikation.md` §1, §2 | Nur vermerkt. Wer den Schutz will, nennt die Kennung in §2; für diesen Plan lohnt das vor der Closure kaum. |

---

**Summary:** 3/3 Liefer-Punkte bestätigt (DoD 2 mit der Reichweiten-Anmerkung V-64) · `make gates` grün auf `b098fb3` · Proben: 7 von 7 Kopf-Proben rot, 2 von 2 Anker-Proben rot, 3 von 3 Richtungs-Proben rot; grün blieben wie erwartet die Reichweiten-Proben PK8, PK9, PD2, PR4 · Review-Findings: F-406 bis F-412 und F-414 behoben oder im Plan adressiert, F-413 offen · Befunde V-62 bis V-65 (LOW) vor der Closure, V-66 vermerkt.
