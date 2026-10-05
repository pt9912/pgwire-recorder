# Verifikation: slice-harness-randformen-vor-code — 2026-10-05

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-randformen-vor-code.md` (Kopf, §1 bis §4, §6) gegen den Gesamt-Diff `d2670ec..fed4387` bei HEAD `fed4387`. Commits: `8fc60d7` (Umsetzung), `6cb0377` (Review-Report), `fed4387` (Nacharbeit F-364 bis F-370). Reiner Prozess-Slice: `AGENTS.md` §3.12, `.claude/commands/implement-slice.md`, `.claude/commands/plan-welle.md`, `.claude/agents/architect.md`, Plan. Netto-Diff an `docs/plan/planning/observations/` ist leer.

**Eingang:**

- die DoD-Liefer-Punkte 1 bis 3 und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-harness-randformen-vor-code.md) (F-364 bis F-371)
- als Prüffall F-350 aus dem [Review zu slice-extended-query-lebendpruefung](2026-10-05-review-slice-extended-query-lebendpruefung.md), dazu der archivierte Slice-Plan aus `docs/plan/planning/done/welle-extended-query/archiv.zip` und die Plan-Stände aus `git show 6d49f23`, `f906c04`, `7ffd029`, `0841891`

Der Arbeitsbaum war bei meinem Start sauber. Die Behauptungen des Implementers zählen nicht als Beleg; `make gates` habe ich selbst gefahren. Einen Test, den man brechen könnte, gibt es nicht. An seine Stelle tritt das Durchspielen an F-350 (Abschnitt 3).

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — Planner und Architect. **Bestätigt.**

- `plan-welle.md:69`–`:71` (Schritt 6, Slices bereitstellen): „Liefert der Slice einen neuen Vertrag, nennt §6 dessen Randformen“, mit Zeiger auf §3.12 und „entschieden werden sie vor dem Code, mit dem Architect“. Die Regel wird nicht neu formuliert.
- `architect.md:23`–`:26`: prüft die Liste in §6 vor dem Code, entscheidet jede Randform aus §6 oder aus einer Rückgabe des Implementers an dem Ort, den §3.12 nennt, oder legt sie dem Nutzer vor. Fehlende meldet er dem Planner.
- Ort laut `AGENTS.md:230`–`:232`: Spezifikation oder ADR; auch eine Entscheidung des Nutzers wird dort festgehalten. Das stimmt mit dem DoD-Wortlaut überein.

Rest: V-38 (Eingang und Ausgang des Rollen-Typs nennen die neue Kante nicht).

### Punkt 2 — Implementer. **Bestätigt.**

- Schritt 12 (`implement-slice.md:111`–`:113`): Ist bei neuem Vertrag eine Randform in §6 nicht entschieden, gilt die Rückgabe, bevor Code entsteht.
- Schritt 13 (`:114`–`:116`): Eine Randform, die §6 nicht nennt und die beim Planen oder Implementieren auffällt, geht denselben Weg.
- Randform-Rückgabe (`:129`–`:132`): Der Implementer entscheidet nicht, hält an und gibt die Randformen im Bericht als Liste *Randform · Frage* an den **Architect** zurück. Der Slice bleibt in `in-progress/`.

Damit ist auch die Lücke aus dem Review benannt: „Schritt 12/13“ ist jetzt tatsächlich in beiden Schritten geändert.

Rest: V-36 (die Aufzählung „weder im Code noch in §6“ ist enger als §3.12).

### Punkt 3 — Regel und Herkunfts-Anker. **Bestätigt.**

- `AGENTS.md:225` §3.12 mit Überschrift „(seit slice-harness-randformen-vor-code)“, Falsch/Richtig/Begründung, Retirement-Check. Die Form entspricht §3.9 bis §3.11.
- Anker an jedem Zielort, geprüft mit `grep`: `plan-welle.md:70`, `architect.md:23`, `implement-slice.md:113` (Schritt 12) und `:116` (Schritt 13).
- `state.md` von `BEO-REPO/spec-randform-erst-im-review-entschieden` steht auf `geplant` mit `Kennung: slice-harness-randformen-vor-code` und ist damit wieder auf dem Stand von `d2670ec`. Den Stand setzt die Closure, wie die DoD es verlangt.

### Punkt 4 — `make gates` grün. **Bestätigt** (eigener Lauf, Abschnitt 5).

### Punkt 5 — Review-Report liegt vor. **Bestätigt.** [Report](2026-10-05-review-slice-harness-randformen-vor-code.md) von `6cb0377`, in anderem Kontext erstellt (Skill `.harness/skills/reviewer.md`).

Die Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen) sind offen. Das ist vor der Closure richtig, und ich habe keine Häkchen gesetzt.

### Der Ablauf als Ganzes

| Frage | Antwort in den Texten | Stelle |
|---|---|---|
| Wer nennt? | der Planner in §6 (Welle: `plan-welle` Schritt 6); der Implementer, wenn er eine nicht genannte Randform findet | `plan-welle.md:69`; `implement-slice.md:114` |
| Wer entscheidet? | der Architect, oder er legt sie dem Nutzer vor | `architect.md:24`–`:26`; `AGENTS.md:232`–`:234` |
| Wo landet die Entscheidung? | in der Spezifikation oder einer ADR; §6 zeigt darauf | `AGENTS.md:227`–`:232`, Beispiel „Richtig“ `:238`–`:239` |
| Wer gibt zurück, an wen, womit? | der Implementer an den Architect, als Liste *Randform · Frage* im Bericht; der Architect meldet fehlende Randformen an den Planner (Artefakt nicht benannt, V-38) | `implement-slice.md:129`–`:131`; `architect.md:26` |
| Was passiert mit dem Slice? | er bleibt in `in-progress/`, kein Lifecycle-Rücksprung; weiter geht es erst nach der Entscheidung | `implement-slice.md:132` |
| Wer zieht §6 nach der Entscheidung nach? | nicht eigens gesagt; §3.9 trägt es für eine Spezifikations-Änderung („wer … die Spezifikation ändert, zieht … §6 nach“), für eine ADR nicht wörtlich (V-39) | `AGENTS.md:180`–`:182` |

Der Ablauf hängt zusammen. Bei einem Slice einer Welle trägt er von der Planung bis zur Rückgabe. Bei einem wellenlosen Slice beginnt er erst mit Schritt 12 und greift dort nur, wenn §6 die Randform schon offen nennt (V-37).

## 2. Review-Findings F-364 bis F-370

| Finding | Ergebnis | Beleg |
|---|---|---|
| F-364 (MEDIUM) | **behoben.** Die Rückgabe hat einen Empfänger (Architect), ein Artefakt (Liste *Randform · Frage* im Bericht) und eine Lifecycle-Aussage (bleibt in `in-progress/`). Sie steht als eigener Absatz neben den Plan-Defekt-Kanten und führt nicht mehr in die Schleife zurück zu 13. §3.12 nennt den Architect statt „den Plan“. Den Weg über §6 schließt „entscheidest du weder im Code noch in §6“; den Weg über die Spezifikation nur §3.12 selbst (V-36). | `implement-slice.md:129`–`:132`; `AGENTS.md:232`–`:234` |
| F-365 (LOW) | **behoben.** Alle vier Zielorte tragen den Anker; `state.md` listet keine Zielorte mehr. | Punkt 3 |
| F-366 (LOW) | **behoben.** Der Stand ist wieder `geplant`; DoD Punkt 3, §1 und §3 übergeben ihn an die Closure (Schritt 25). | `state.md:1`; Plan §1, §2, §3 |
| F-367 (LOW) | **teilweise behoben.** Schritt 12 ist jetzt ein Auslöser vor dem Code, auch für wellenlose Slices, und Risiko 3 sagt nicht mehr zu, als der Ablauf trägt. Kein Command-Schritt ruft bei einem wellenlosen Slice den Architect vor dem Code auf. Der Satz „Der Architect prüft die Liste vor dem Code“ in §3.12 gilt dort nur über die Rückgabe (V-37). | `implement-slice.md:111`; Plan §6 Risiko 3; `AGENTS.md:232` |
| F-368 (LOW) | **behoben im beanstandeten Satz.** Der Handelnde ist genannt („meldest du sie dem Planner“), und „dem Nutzer vorlegen“ deckt sich mit §3.12. Eingang und Ausgang des Rollen-Typs sind unverändert (V-38). | `architect.md:23`–`:26`, `:9`–`:10` |
| F-369 (LOW) | **behoben.** §1 nennt „Gate“ und bindet an die Liste aus §3.10. | Plan §1 Ziel |
| F-370 (INFO) | **aufgenommen.** Entscheidungsort ist nur noch Spezifikation oder ADR; §6 zeigt darauf, die Entscheidung des Nutzers wird dort festgehalten. Eine entschiedene Randform in §6 braucht bei der Closure weiterhin einen Ausgang. Dafür passt *entfallen* mit dem Ort als Begründung, wie es `slice-extended-query-lebendpruefung` für Risiko 2 tat. Kein Befund. | `AGENTS.md:230`–`:232` |

F-371 (INFO) verlangte keine Aktion.

## 3. Durchspielen an F-350 (`\v` je Serverversion)

Ich habe geprüft, wann die Randform im Plan von `slice-extended-query-lebendpruefung` stand und welche Rolle sie schrieb.

| Stand | Commit | §6 nennt `\v`/Serverversion? |
|---|---|---|
| Übernahme in `in-progress/` | `6d49f23` | nein |
| Architect: [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md), Plan nachgezogen | `f906c04` (10:00:06) | nein |
| Spezifikation `LH-FA-09.a` mit Leerraum „nach dem Scanner von PostgreSQL“, ohne Version | `7ffd029` (10:42:10) | **ja**, neues Risiko 5: „Ob ältere Server den vertikalen Tabulator (`\v`) ebenso als Leerraum lesen, ist nicht geprüft; die Spezifikation nennt keine Serverversion … **Ausgang:** offen bis Closure.“ |
| erster Code-Commit | `0841891` (10:42:36) | unverändert offen |

Der Implementer hat die Randform also beim Schreiben der Spezifikation selbst gefunden, 26 Sekunden vor dem Code-Commit, und sie offen in §6 eingetragen. Dann hat er weitercodiert. Das Review fand sie als HIGH (F-350), und der Nutzer entschied sie danach.

**Mit den neuen Texten:**

1. *Planung und Architect* (`plan-welle` Schritt 6, `architect.md`). Diese Stellen greifen hier nicht, weil §6 die Randform bei `f906c04` nicht nannte. Ob Planner oder Architect sie bemerkt hätten, ist Urteil. Das Beispiel „Verhalten je Serverversion“ in §3.12 macht es wahrscheinlicher, sichert es aber nicht.
2. *Schritt 12.* Greift ebenfalls nicht: Beim Eintritt in den Plan stand die Randform noch nicht in §6.
3. *Schritt 13 und Randform-Rückgabe.* **Diese Stelle greift.** Der Implementer findet beim Implementieren eine Randform, die §6 nicht nennt. Sie ist offen, und §3.12 führt „Verhalten je Serverversion“ wörtlich als Beispiel. Nach `implement-slice.md:129`–`:132` hätte er angehalten, statt §6 als Risiko fortzuschreiben und weiterzucodieren. Er hätte als Liste zurückgegeben: *`\v` als Leerraum · gilt das je Serverversion, und was gilt vor der Zuordnung?* Das wäre vor `0841891` beim Architect angekommen, der Slice wäre in `in-progress/` geblieben.

**Urteil:** Ja, die neuen Texte hätten F-350 vor dem Code an den Architect gebracht, und zwar über Schritt 13 und die Randform-Rückgabe, nicht über die Planungs-Hälfte.

**Wo der Auslöser schwach ist:** Er hängt daran, dass der Implementer die Frage selbst bemerkt und als offen erkennt. Hier stand sie im Plan bewusst unter „Risiken“ mit dem Standard-Ausgang „offen bis Closure“, also in derselben Form wie die drei entschiedenen Einträge darüber. Hätte er die Frage in `7ffd029` selbst beantwortet, etwa „`\v` gilt immer“, stünde in §6 ein Zeiger auf `LH-FA-09.a` als „entschieden in der Spezifikation“, und das vor dem ersten Code-Commit. Formal wäre §3.12 damit erfüllt. Den Implementer hält davon nur der allgemeine Satz in §3.12 ab („entscheidet der Implementer nicht“). Die Aufzählung im Command, „weder im Code noch in §6“, nennt die Spezifikation nicht, obwohl in diesem Repo der Implementer im Slice die Spezifikation schreibt (`7ffd029`, `4149080`). Das ist V-36.

## 4. Plan gegen Diff

- **Kopf:** `Bezug: —`, `Berührte Spec-Stellen: —`, `Verantwortlich: pt9912`. Stimmt, kein Spec-Stratum im Diff.
- **§1:** Ziel einschließlich „Gate“ und Rückgabe an den Architect; Abgrenzungen (vendored Baum, Sensor, Altbestand, Produkt/Spec/Lastenheft, übrige Rollen-Typen und Skill) treffen zu. `planner.md`, `implementer.md`, `reviewer.md` und `.harness/` sind im Diff unberührt.
- **§2:** Die Liefer-Punkte 1 bis 3 entsprechen dem Diff (Abschnitt 1).
- **§3:** Jede Zeile entspricht einer Änderung; `state.md` ist als „update bei Closure“ geführt und im Netto-Diff unverändert.
- **§4:** Keine Rückführungs-Bedingung ist eingetreten; der vendored Baum ist unberührt.
- **§6:** Risiko 3 ist auf den Stand von `fed4387` nachgezogen; Risiko 1 und 2 stimmen.

**Gegen bestehende Regeln:**

- **§3.9:** kein Widerspruch. §3.9 verlangt, §6 nachzuziehen, wenn die Spezifikation sich ändert; die Rückgabe verbietet, eine Randform in §6 zu *entscheiden*. Sie zu nennen bleibt erlaubt.
- **§3.10:** §3.12 übernimmt dessen Vertragsliste über „(§3.10)“, statt eine eigene zu führen.
- **§3.11:** Die Plan-Zeilen sagen nicht mehr zu, als der Ablauf trägt. Den Rest im Regeltext beschreibt V-37.
- **Rücksprung-Abschnitt:** kein Widerspruch. Die Randform-Rückgabe steht als eigener Absatz neben den Plan-Defekt-Kanten. Sie ist ausdrücklich kein Lifecycle-Rücksprung (Schritt 11), sondern ein Anhalten in `in-progress/`.
- **Command-Kopf** (`implement-slice.md:10`, „Findings · Folge-ADR · Carveout“): Mit der Liste *Randform · Frage* gibt es ein viertes Übergabe-Artefakt, das der Kopf nicht aufzählt (V-38).

## 5. Gate-Lauf

`make gates` auf `fed4387`, sauberer Arbeitsbaum, eigener Lauf: **Exit 0.**

- `abdeckung-gegenprobe`: gruen
- `a-check`: 0 Befunde. Der Hinweis auf 5 Integrationstest-Dateien außerhalb der Schichten ist bekannter Bestand.
- `a-check-negativ`: gruen
- `baseline-verify`: v6.13.0 OK, 54 Dateien
- `build` und `test`: Stufen `runtime` und `test` aus dem Docker-Cache, unveränderte Go-Quellen
- `d-check`: 190 Dateien, 0 Befunde
- `commit-msg-gegenprobe`: gruen
- `run-integration-tests`: gruen
- `abdeckung-check` und `record-gates` gaben nichts aus und haben den Lauf nicht abgebrochen

Nach dem Commit dieser Datei habe ich `make gates` erneut gefahren. Der erste Lauf meldete eine unverlinkte ADR-Kennung in diesem Bericht (`id-unlinked`). Ich habe sie verlinkt, und der Lauf auf dem Commit des Berichts ist grün.

## 6. Befunde

| ID | Kategorie | Befund | Pfad |
|---|---|---|---|
| V-36 | LOW | Die Randform-Rückgabe sagt „entscheidest du weder im Code noch in §6“. Die Spezifikation fehlt in dieser Aufzählung, obwohl in diesem Repo der Implementer im Slice die Spezifikation schreibt. Genau dort hat er in F-350 die Leerraum-Menge ohne Version festgelegt (`7ffd029`), vor dem ersten Code-Commit. Ein solcher Eintrag erfüllt den Buchstaben von §3.12 („entschieden … in der Spezifikation“). Nach Source Precedence gilt der allgemeinere Satz in §3.12, also kein Widerspruch. Der Command lädt aber zu der engeren Lesart ein. | `.claude/commands/implement-slice.md:129`–`:130`; `AGENTS.md:230`–`:234` |
| V-37 | LOW | F-367, Rest. Bei einem wellenlosen Slice gibt es vor dem Code keinen Schritt, der den Architect aufruft. Schritt 12 greift nur bei einer Randform, die §6 nennt und offen lässt, nicht bei einer fehlenden Liste. Damit gilt „Der Architect prüft die Liste vor dem Code“ dort nur über die Rückgabe. Im F-350-Fall hätte auch bei einem Slice einer Welle die Planungs-Hälfte nicht gegriffen (Abschnitt 3). Plan §6 Risiko 3 führt den Rest ehrlich und bis zur Closure offen. Eine Aktion in diesem Slice ist nicht nötig; der Rest braucht dort einen Ausgang. | `AGENTS.md:232`; `.claude/commands/implement-slice.md:111`; Plan §6 Risiko 3 |
| V-38 | LOW | F-368, Rest. Kopf und `description` von `architect.md` sind unverändert. Der Eingang nennt nur den Slice-Plan vom Planner, nicht die Rückgabe des Implementers. Der Ausgang nennt weder die Randform-Entscheidung (Spezifikation oder ADR) noch die Meldung einer fehlenden Randform an den Planner, und für diese Meldung ist kein Artefakt benannt. Ebenso zählt der Kopf von `implement-slice.md` die Liste *Randform · Frage* nicht zu den Übergabe-Artefakten. Die Kanten stehen im Fließtext; die Kopfzeilen, die Modul 8 als Schnittstelle liest, wissen von ihnen nichts. | `.claude/agents/architect.md:3`, `:9`–`:10`, `:26`; `.claude/commands/implement-slice.md:10` |
| V-39 | INFO | Wer §6 nach der Entscheidung auf den Ort zeigen lässt, sagt kein neuer Text. Bei einer Spezifikations-Änderung trägt es §3.9; bei einer ADR nennt §3.9 sie nicht. In der Praxis zog der Architect den Plan im ADR-Commit nach (`f906c04`). Kein Handlungsbedarf, solange der Brauch hält. | `AGENTS.md:180`–`:182` |

**Summary:** 0 HIGH · 0 MEDIUM · 3 LOW · 1 INFO. DoD-Liefer-Punkte 1 bis 3 bestätigt, `make gates` grün auf `fed4387`. F-364 bis F-366, F-368 im beanstandeten Satz und F-369 sind behoben, F-370 ist aufgenommen, F-367 teilweise. Am Fall F-350 greift der neue Ablauf über Schritt 13 und die Randform-Rückgabe vor dem Code, nicht über Planung oder Architect-Prüfung. Er hängt daran, dass der Implementer die Frage bemerkt. Merge-blockierend ist keiner der Befunde.
