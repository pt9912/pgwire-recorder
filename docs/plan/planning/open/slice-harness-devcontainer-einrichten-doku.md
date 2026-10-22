# Slice slice-harness-devcontainer-einrichten-doku: `docs/developer/devcontainer.md` §6 und §9 nennen `make container-einrichten` statt der Einzelbefehle

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (dieselbe Welle wie der Geber; `AGENTS.md` §3.11 verlangt den Doku-Folge-Slice **in derselben Welle**, direkt hinter dem Slice, der das Verhalten liefert). Angelegt am 2026-10-21 durch den Schnitt von `slice-harness-devcontainer-einrichten` (Entscheidung des Nutzers vom 2026-10-21, dritter von drei Punkten: den Doku-Teil vorab abspalten, weil die Schätzung des Gebers mit ihm bei 815 bis 1210 über der Marke 900 lag und er eine zweite Schicht trug).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als das Dokument die Host-Voraussetzungen des Repos beschreibt (`ssh`, Entscheidung des Nutzers vom 2026-10-17; §1); kein Anspruch des Produkts. Register: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (3×, Ausgang `geplant` mit diesem Slice seit 2026-10-22). Geber: `slice-harness-devcontainer-einrichten` (Liefer-Punkt 3 *Dokument*, §1), `slice-harness-devcontainer-u-boot-0-10` (`done/`; Dokument-Hälfte der Sendung „Arbeit im Container“ und die gemessenen Läufe) und `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (`open/`; Abschnitt über den Gegenproben-Container). Entscheidung: [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) (Accepted 2026-10-21) — das Dokument beschreibt das Ziel, das sie festlegt. §4 rechnet in der berichtigten Größen-Metrik: [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Status `Proposed`; `slice-harness-groessen-metrik` trägt sie in `AGENTS.md` §3.14 ein, bis dahin gilt sie aus der ADR).

**Berührte Spec-Stellen:** [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) (das Dokument beschreibt das Werkzeug, das diese Kennung festlegt; es legt nichts fest und sagt nichts zu, was dort nicht steht) · [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) (§3 bis §5 des Dokuments, vom Geber `slice-harness-devcontainer-einrichten-geruest` geschrieben, bleiben unverändert) · Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/developer/devcontainer.md` beschreibt den Ist-Zustand: §6 und §9 nennen `make container-einrichten` als den Weg und führen die Einzelbefehle nur noch, soweit das Dokument sie als Erklärung braucht; dazu steht dort der Ablauf eines Implementer-Laufs im Container-Klon (Arbeit im Klon per `ssh`, `make gates` dort, Ergebnis per `git fetch container` am Host). Jede Aussage ist durch einen Lauf belegt oder steht in §13 (Ungeprüftes). Dazu steht dort der Weg, wenn sich der Hostschlüssel eines Containers nach einem Neubau ändert (Entscheidung des Nutzers vom 2026-10-22, §1).

**Sendung von `slice-harness-devcontainer-einrichten` (Liefer-Punkt 3 *Dokument* des Gebers, eingetragen vom Planner am 2026-10-21, `AGENTS.md` §3.13; Entscheidung des Nutzers vom selben Tag):** Dieser Slice **nimmt an**: §6 und §9 des Dokuments nennen das Ziel statt der neun Einzelbefehle, §13 trägt, was ungeprüft bleibt. Der Geber ändert am Dokument nichts mehr (dessen §1, Ausschluss der Klasse 1). Grund des Schnitts: Mit dem Dokument lag die Schätzung des Gebers bei 815 bis 1210 über der Marke 900 und er trug zwei Schichten (Harness-Werkzeug **und** Wartungs-Doku); `AGENTS.md` §3.11 nennt für genau diesen Fall den Doku-Folge-Slice `<Kennung>-doku` direkt dahinter, in derselben Welle.

**Sendung „Arbeit im Container“, Dokument-Hälfte (von `slice-harness-devcontainer-u-boot-0-10` über `slice-harness-devcontainer-einrichten`, eingetragen vom Planner am 2026-10-21):** Dieser Slice **nimmt an**: den Ablauf eines Implementer-Laufs im Container-Klon im Dokument — Arbeit im Klon, `make gates` dort mit der Dauer, die `slice-harness-devcontainer-u-boot-0-10` gemessen hat (zwei Läufe zu je 1701 s, höchstens 5,42 GiB), Ergebnis per `git fetch container`. Die Zahlen werden **nicht neu gemessen**; sie stammen aus jenem Slice und stehen mit seiner Kennung da. Die Einrichtungs-Hälfte derselben Sendung (das Ziel) liegt beim Geber `slice-harness-devcontainer-einrichten`; die Adresse bleibt einen Hop länger, in `slice-harness-devcontainer-u-boot-0-10` (`done/`) wird kein Zeiger umgeschrieben.

**Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (eingetragen vom Planner am 2026-10-21, zuvor beim Geber; `AGENTS.md` §3.13):** Dieser Slice **nimmt an**: Der Abschnitt über den Gegenproben-Container in `docs/developer/devcontainer.md` verweist auf §3 bis §6 und §9 und nennt deren Einzelbefehle; wer §6 und §9 mit `make container-einrichten` umschreibt, erhält diesen Abschnitt gültig oder meldet es dem Planner. Die Kennung steht in der Closure-Notiz. Kein eigener Liefer-Punkt, keine eigene Schicht (§8).

**Sendung `Hostschlüssel nach Neubau` (Ausgang `geplant` von `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`, eingetragen vom Planner im Commit der Zuweisung, 2026-10-22, `AGENTS.md` §3.13; Entscheidung des Nutzers vom 2026-10-22, Variante A, Wortlaut „ja“ auf die Vorlage mit den Empfehlungen des Koordinators, so übermittelt):** **Keine Einstellung in den Werkzeugen; `PGR_SSH_OPTS` mit einer eigenen `known_hosts`-Datei ist der dokumentierte Weg.** Das Werkzeug setzt keine Einstellung zum Hostschlüssel ([`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 20 (5)); wer nach einem Neubau des Images eine Warnung von `ssh` bekommt oder für einen neuen Container keinen Eintrag hat, ruft es mit `PGR_SSH_OPTS` auf, das eine eigene Datei nennt (`-o UserKnownHostsFile=<Datei>`), und lässt `known_hosts` des Nutzers unberührt. Dieser Slice **nimmt an**: Das Dokument sagt das an der Stelle, die heute „das Dokument empfiehlt dazu keine Einstellung (§13)“ sagt (§5 des Dokuments, ein Satz), und führt es in §13 mit dem, was nicht gefahren ist. **Die Stelle liegt in §5, und §5 ist oben ausgeschlossen** (`slice-harness-devcontainer-einrichten-start`, `done/`): Der Planner fasst den Ausschluss für **genau diesen einen Satz** enger (unten, *Ausdrücklich NICHT*), weil der Satz mit der Entscheidung falsch würde und zwei Aussagen zum selben Schlüssel in einem Dokument widersprächen. Der Handbuchsatz liegt in der Schicht dieses Slice (Wartungs-Doku, `docs/developer/`), kein Liefer-Punkt und keine Schicht kommen hinzu (§8). **Was die Entscheidung nicht sagt:** Sie empfiehlt `StrictHostKeyChecking=accept-new` nicht; der Nachweis des Gebers fuhr es mit einer eigenen Datei, und das Dokument nennt es als das, was lief, nicht als Empfehlung. Der **Zielort** (der Satz im Dokument) entsteht erst in diesem Slice; deshalb ist der Ausgang des Eintrags `geplant`, nicht `verkörpert`; der Anker wäre `seit slice-harness-devcontainer-einrichten-doku`.

**Das Dokument beschreibt den Ist-Zustand** (`AGENTS.md` §3.11, Entscheidung des Nutzers vom 2026-10-10 zu `slice-doku-ist-stand`): keine Chronik („noch nicht“, „kommt mit“), kein Zielstand, kein Verweis auf Slices, Wellen, Reviews oder die Roadmap. Ein Satz über eine plattformabhängige Zusage nennt allein Linux (Geber §6 *Plattform*, `SPEC-069` Punkt 19). Was gilt, aber nicht gemessen ist, steht in §13, nicht im Fließtext.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Code, kein Fragment, keine Gegenprobe, kein Vertrag** (`tools/`, `harness/mk/`, `harness/vertraege/`, `harness/README.md`) — Klasse 1: Das Werkzeug `container-einrichten` und seine Gegenprobe samt Zeilen im Gate-Index liefert `slice-harness-devcontainer-einrichten` (dessen §1, §2 und §3; sein §4 *Start* nennt diesen Slice direkt dahinter). Grund: Dieser Slice beschreibt, was dort geliefert ist, und beschreibt nichts, was es nicht gibt (`AGENTS.md` §3.11). Damit ist beim Review sofort prüfbar, dass der Diff nur `docs/developer/` berührt.
- **§3 bis §5 des Dokuments** (Gerüst, Start, Stopp und Neustart des Containers) — Klasse 1: `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` (beide `done/`) haben sie geschrieben. Grund: Sie beschreiben `make container-start`, dessen Vertrag dieser Slice nicht berührt; wer sie anfasst, ändert die Beschreibung eines fremden Werkzeugs. **Ausgenommen** ist der eine Satz in §5, der zum Hostschlüssel sagt, das Dokument empfehle keine Einstellung (Sendung `Hostschlüssel nach Neubau`, oben); der Rest von §3 bis §5 bleibt unverändert.
- **§7 des Dokuments** (Zugang von Claude Code im Container) — Klasse 1: `slice-harness-devcontainer-einrichten-zugang` (`open/`, dessen §1 nimmt diese Sendung an; Entscheidung des Nutzers vom 2026-10-18: zurzeit nicht gebraucht). Grund: Das Ziel löst keinen Zugang auf, und ein Abschnitt über einen Weg, den es nicht gibt, wäre ein Zielstand.
- **Keine Messung von `make gates` im Container und keine neue Zeitangabe** — Bestand bleibt bewusst stehen (Klasse 2): Die Zahlen stammen aus `slice-harness-devcontainer-u-boot-0-10` (`done/`) und stehen mit seiner Kennung da. Grund: Eine eigene Messung wäre ein zweiter, teurer Lauf ohne neue Aussage; weicht die Zahl im Klon von `pgr-dev2` ab, ist das ein Befund an den Planner (§6).
- **Eine Einstellung in den Werkzeugen** (`StrictHostKeyChecking` oder `UserKnownHostsFile` als Vorgabe von `make container-einrichten` oder `make container-start`) — Entscheidung des Nutzers vom 2026-10-22 (Variante A): keine; das Werkzeug setzt keine Einstellung (`SPEC-069` Punkt 20 (5)). Grund: Eine Vorgabe lockerte die Prüfung des Hostschlüssels für jeden Aufruf; der Aufrufer wählt sie je Aufruf über `PGR_SSH_OPTS`. Dieser Slice beschreibt den Weg und ändert kein Werkzeug (Klasse 4, kein Code).
- **Kein Eintrag in `harness/README.md`** — Klasse 1 (Geber, dessen §3) und Klasse 4 (Schicht): Der Gate-Index ist Teil des Werkzeugs, nicht der Wartungs-Doku. Grund: Zwei Slices an derselben Tabelle erzeugen einen Konflikt statt eines Index.
- **Kein Produkt-Code und keine Nutzer-Doku** (`cmd/`, `internal/`, `docs/user/`) — Schicht-Abgrenzung (Klasse 4). Grund: Das Werkzeug richtet eine Entwicklungs-Sitzung ein; es ändert nichts, was ein Benutzer des Produkts sieht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkt (einer):

- [ ] **Dokument** `docs/developer/devcontainer.md` §6 und §9 nennen `make container-einrichten` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht (Sendung von `slice-harness-devcontainer-einrichten`, §1); dazu der Ablauf eines Implementer-Laufs im Container-Klon (Dokument-Hälfte der Sendung „Arbeit im Container“, §1): Arbeit im Klon per `ssh`, `make gates` dort mit der Dauer, die `slice-harness-devcontainer-u-boot-0-10` gemessen hat, Ergebnis per `git fetch container`. §3 bis §5 gehören `slice-harness-devcontainer-einrichten-geruest`, §7 `slice-harness-devcontainer-einrichten-zugang` (§1). Der Abschnitt über den Gegenproben-Container bleibt gültig oder der Befund geht an den Planner (Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`, §1). Jede Aussage ist durch einen Lauf belegt oder steht in §13; keine Chronik und kein Zielstand (`AGENTS.md` §3.11). **Hostschlüssel (Sendung, §1):** Ein Satz in §5, der „das Dokument empfiehlt dazu keine Einstellung“ ersetzt, und ein Eintrag in §13 nennen den Weg (`PGR_SSH_OPTS` mit eigener `known_hosts`-Datei) und sagen, dass die Werkzeuge keine Einstellung setzen; `StrictHostKeyChecking=accept-new` steht als das, was der Nachweis fuhr, nicht als Empfehlung.

Gates und Closure:

- [ ] `make docs-check` grün.
- [ ] `make gates` grün. Der Stop-Hook gibt `docs/developer/` über seine Positivliste frei (seit `slice-harness-stop-hook-positivliste-developer`); berührt der Slice wirklich nur diesen Pfad, genügt der Stand der schnellen Kette, und der volle Lauf des Gebers trägt weiter. Der Bericht nennt, welcher der beiden Fälle vorlag.
- [ ] **Die Befehle des Dokuments sind gefahren, nicht gelesen** (`AGENTS.md` §3.11): Jeder Aufruf, den §6 und §9 nennen, ist in dieser Form gegen einen laufenden Container gelaufen, und der Bericht nennt je Aufruf das Ergebnis. Ein Aufruf, der nicht lief, steht in §13 als ungeprüft. Dazu der Aufruf mit `PGR_SSH_OPTS` und eigener `known_hosts`-Datei: gegen einen Container ohne Eintrag gefahren; der Fall des geänderten Schlüssels nach einem Neubau ohne Cache ist eine Neubau-Probe, und was der Slice nicht fährt, steht mit der Kennung der Probe des Gebers (`slice-harness-devcontainer-u-boot-0-10`) in §13 als nicht neu gefahren.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/developer/devcontainer.md` §6 | update | `make container-einrichten` statt der Einzelbefehle; die Eingaben des Ziels (`PGR_DEV`, `PGR_SSH_PORT`, `PGR_REMOTE`, `PGR_ZWEIG`, `PGR_SSH_OPTS`) und seine Exit-Codes nach [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkte 2 und 14; `printf %q` entfällt, weil die Identität über den Standardeingang geht (Punkt 9) |
| `docs/developer/devcontainer.md` §9 | update | Ablauf eines Implementer-Laufs im Container-Klon: Zweig, Arbeit per `ssh`, `make gates` dort mit der Dauer aus `slice-harness-devcontainer-u-boot-0-10`, Ergebnis per `git fetch container` |
| `docs/developer/devcontainer.md` §13 | update | was ungeprüft bleibt: jeder Aufruf, der nicht gefahren wurde, und die Grenzen aus [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 20, die das Dokument berührt |
| `docs/developer/devcontainer.md` §5 (ein Satz) | update | „das Dokument empfiehlt dazu keine Einstellung“ wird durch den Weg ersetzt (Variante A, Sendung in §1); der übrige Abschnitt bleibt (§1) |
| Abschnitt über den Gegenproben-Container im selben Dokument | prüfen | Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (§1): bleibt nach dem Umbau von §6 und §9 gültig, sonst Befund an den Planner |

- **Reihenfolge der Arbeit:** erst die Aufrufe gegen einen laufenden Container fahren, dann schreiben. Ein Satz entsteht aus einem Lauf, nicht ein Lauf aus einem Satz (`AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`).
- **Vorbehalte der Vorfassung prüfen:** §11 des Dokuments trägt Sicherheits- und Isolationsaussagen („Behandle den Container nicht als Sandbox“). Wer §6 und §9 kürzt, vergleicht die Vorbehalte der Vorfassung mit dem Diff; eine Aussage, die ohne Entscheidung wegfällt, ist ein Befund (§6).
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-einrichten` liegt in `done/` — der Slice beschreibt das gelieferte Ziel, und vor dessen Lieferung gäbe es nichts zu beschreiben (`AGENTS.md` §3.11: der Doku-Folge-Slice steht **direkt dahinter**, in derselben Welle). Dazu: ein laufender Dev-Container, gegen den die Aufrufen des Dokuments gefahren werden können (DoD), und WIP-Limit frei. Den Anspruch trägt der Nutzer; `Verantwortlich:` setzt der Übergang `open` → `next`. Die Welle `welle-dev-container` schließt erst mit diesem Slice (deren §3 und §5).

**Größe.** Metrik in der berichtigten Fassung ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md), Status `Proposed`; `slice-harness-groessen-metrik` trägt sie in `AGENTS.md` §3.14 ein, bis dahin gilt sie aus der ADR):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- docs/developer \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Der einzige Pfad ist `docs/developer`; Plan, Register und Reports liegen unter den zwei Ausschlüssen und zählen nicht. Messpunkt ist der committete Stand zu jeder Zeit der Arbeit, nicht erst der letzte Commit.

**Schätzung (Planner, 2026-10-21; übernommen aus der Berichtigung des Architect vom 2026-10-21 zum Geber, §4 dort, Zeile *Dokument*): 70 bis 120; nachgezählt am 2026-10-22 mit dem Satz zum Hostschlüssel (ein Satz in §5, ein Eintrag in §13, höchstens 10 Zeilen): 75 bis 130** — §6 umgeschrieben und gekürzt (Löschungen zählen), §9 um den Ablauf des Implementer-Laufs ergänzt, §13 um das Ungeprüfte, Hinweis auf den Nachweis. Faktor Tests zu Produktcode entfällt: kein Code, kein Test, keine Gegenprobe (§1); der Beleg ist ein gefahrener Aufruf, kein geschriebener Test. **Marke 250**, **Meldeschwelle 180**. Die Mitte der Schätzung (102 nach der Nachzählung) liegt bei 41 % der Marke.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik oben liegt **über 250** (Meldeschwelle 180: anhalten und dem Planner melden, bevor weitergearbeitet wird), oder der Umbau verlangt mehr als einen Liefer-Punkt oder eine zweite Schicht (etwa eine Zeile in `harness/README.md` oder eine Änderung am Werkzeug, §1). Dann schneidet der Planner, Vorschlag: §6 (das Ziel statt der Einzelbefehle) zuerst, §9 mit dem Ablauf des Implementer-Laufs als zweiter Slice.
- `in-progress` → `open` (blockiert): Ein Aufruf, den das Dokument nennen soll, lässt sich nicht fahren (kein Container verfügbar, das Ziel verhält sich anders als `SPEC-069` sagt) — dann geht der Befund als Rückgabe an den Architect oder an den Planner, und das Dokument bleibt stehen, statt eine ungeprüfte Aussage zu tragen (`AGENTS.md` §3.11).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Zwei beobachtbare Kriterien: (1) Jeder Aufruf, den §6 und §9 nennen, ist in dieser Form gegen einen laufenden Container gelaufen, und der Bericht nennt je Aufruf das Ergebnis; was nicht lief, steht in §13 als ungeprüft. (2) Der Diff berührt allein `docs/developer/devcontainer.md` — `git diff --name-only` nennt keine andere Datei außer dem Plan, dem Register und dem Review-Report. Ein Mutationsprotokoll (`AGENTS.md` §3.10) entfällt: Der Slice liefert keinen Vertrag und keine Zusage des Codes; der Beleg einer Dokument-Aussage ist ihr Lauf.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert **keinen neuen Vertrag**: kein Werkzeug, kein Gate, keine Zusage mit Randformen, also kein Fall nach `AGENTS.md` §3.12. Die Randformen des Ziels stehen in [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) und im §6 des Gebers; dieser Slice entscheidet keine und beschreibt nur, was dort entschieden ist. Trifft er auf eine Form, die der Vertrag nicht nennt, hält er an und gibt sie dem Architect zurück, statt sie im Dokument zu entscheiden.

Die drei Dokument-Risiken aus dem Geber liegen seit dem 2026-10-21 hier (Schnitt, §1):

- **Das Dokument veraltet still gegen das Ziel** — kein Sensor prüft, ob §6 und §9 die Schritte richtig nennen ([`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 20 (9)), und eine Aussage reicht leicht weiter als der Lauf, der sie belegt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** *weiter offen* → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, oder *eingetreten* mit einem Beleg dieses Slice, wenn ein Satz des gelieferten Stands über seinen Lauf hinausreicht.
- **Der Umbau verliert Vorbehalte der Vorfassung** (Sicherheits- und Isolationsaussagen aus §11, „Behandle den Container nicht als Sandbox“; §3). — **Ausgang (Closure):** *entfallen* mit dem Vergleich der Vorbehalte der Vorfassung gegen den Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Die Sendung „Arbeit im Container“ wächst den Dokument-Punkt** (`BEO-REPO/slice-waechst-durch-uebernahmen`): §9 trägt neben dem Einrichtungsweg den ganzen Ablauf eines Implementer-Laufs. — **Ausgang (Closure):** *entfallen* mit der gemessenen Summe nach der Metrik aus §4 unter der Marke 250, oder *eingetreten* → Schnitt nach §4 (§6 zuerst, §9 als zweiter Slice), Kennung vom Planner.

Weitere Risiken:

- **Der Satz zum Hostschlüssel sagt mehr zu, als ein Lauf belegt** (etwa eine Empfehlung von `accept-new`, die die Entscheidung des Nutzers nicht trägt). — **Ausgang (Closure):** *entfallen*, wenn der Satz nur nennt, was lief, und die Entscheidung vom 2026-10-22 mit ihrem Wortlaut („keine Einstellung in den Werkzeugen, `PGR_SSH_OPTS` mit eigener `known_hosts`“) wiedergibt; sonst *eingetreten* → Befund an den Planner.
- **Der enger gefasste Ausschluss für §5 wächst zu mehr als einem Satz.** — **Ausgang (Closure):** *entfallen*, wenn `git diff` im Review in §3 bis §5 genau die eine Stelle zeigt; sonst *eingetreten* → der Planner entscheidet, ob ein weiterer Satz die Sendung trägt oder zurückgenommen wird.
- **Die Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` trägt nicht**: Der Abschnitt über den Gegenproben-Container verweist auf §6 und §9 und nennt deren Einzelbefehle; nach dem Umbau zeigt er auf Text, der nicht mehr dasteht. Der Nehmer dieser Bindung ist dieser Slice (§1); jener Slice liegt in `open/` und ist nicht geschlossen. — **Ausgang (Closure):** *entfallen*, wenn der Abschnitt im gelieferten Stand gültig ist und die Kennung in der Closure-Notiz steht; sonst *eingetreten* → Befund an den Planner, der entscheidet, ob der Abschnitt hier oder in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` nachgezogen wird.
- **Der Geber liefert etwas anderes, als der Vertrag heute sagt** (andere Eingaben, andere Exit-Codes, ein Zwischenstand nach einem Schnitt an der Naht aus dessen §4): Dann beschreibt dieses Dokument ein Ziel, das es nicht gibt. — **Ausgang (Closure):** *entfallen*, weil der Start den Geber in `done/` verlangt (§4) und das Dokument gegen den **gelieferten** Stand geschrieben und gefahren wird; weicht der gelieferte Stand von [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) ab, ist das ein Befund an den Architect (nicht an dieses Dokument).
- **Schneidet der Geber an der Naht seines §4** (Teil 1: Eingaben, Werkzeuge, Zustand; Teil 2: die handelnden Schritte), beschreibt dieser Slice nach Teil 1 ein Ziel, das noch nicht einrichtet, sondern nur prüft. — **Ausgang (Closure):** *entfallen*, wenn der Geber ungeteilt in `done/` liegt; *eingetreten* → der Planner entscheidet, ob dieser Slice hinter Teil 2 wandert oder ob Teil 1 seinen Zwischenstand selbst dokumentiert (Kennung vom Planner).
- **`make gates` im Klon von `pgr-dev2` weicht von der Messung ab** (`slice-harness-devcontainer-u-boot-0-10`: zwei Läufe zu je 1701 s, höchstens 5,42 GiB; anderes Gerüst, andere Last). Das Dokument nennt die Zahl mit der Kennung ihres Slice und nicht als eigene Messung (§1). — **Ausgang (Closure):** *entfallen* mit dem Satz im Dokument, der die Herkunft der Zahl nennt; ein abweichender Lauf ist *eingetreten* → Befund an den Planner und ein Eintrag im Register.
- **Die Plattformgrenze verschwindet beim Kürzen** (`BEO-REPO/verhalten-nur-unter-linux-geprueft`; `AGENTS.md` §3.12 *Die Plattform gehört dazu*): Geprüft ist ein Linux-Host gegen einen Linux-Container, und ein Satz über eine plattformabhängige Zusage nennt allein Linux ([`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 19). — **Ausgang (Closure):** *entfallen*, wenn der gelieferte Stand die Grenze an der Zusage nennt; sonst *weiter offen* → `BEO-REPO/verhalten-nur-unter-linux-geprueft`.
- **Der Stop-Hook verlangt den vollen Lauf**, obwohl der Diff nur `docs/developer/` berührt (Positivliste seit `slice-harness-stop-hook-positivliste-developer`). — **Ausgang (Closure):** *entfallen* mit der Angabe im Bericht, welcher der beiden Fälle vorlag (Freigabe über die Positivliste mit grüner schneller Kette, oder voller Lauf); ein unerwartetes `block` ist *eingetreten* → Befund an den Planner.

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-21, gemergter Stand). Treffer für ein Dokument über ein Werkzeug: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Stand verkörpert in `AGENTS.md` §3.11; zwei der drei Dokument-Risiken oben hängen daran), `BEO-REPO/slice-waechst-durch-uebernahmen` (Stand verkörpert in `AGENTS.md` §3.13; dieser Slice **ist** die Antwort darauf — der Dokument-Punkt wächst nicht länger im Geber), `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, `BEO-REPO/verhalten-nur-unter-linux-geprueft` (Stand verkörpert, §6), `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (§1, nicht Gegenstand), `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` (§2, der Pfad steht seit `slice-harness-stop-hook-positivliste-developer` auf der Positivliste), `BEO-REPO/gemeinsamer-klon-im-container-wechselt-zweig` (§9 des Dokuments beschreibt je Slice einen Zweig). Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Vertrag und Entscheidung gehen dem Dokument voraus — [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) und [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) liegen, und das Werkzeug liegt vor diesem Slice).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-21, im Commit, der diesen Plan anlegt und die Sendung aus `slice-harness-devcontainer-einrichten` einträgt):** **Liefer-Punkt 1** — das Dokument (§6, §9, §13). Mit ihm kamen drei Sendungen, die **keinen** eigenen Liefer-Punkt tragen: die Dokument-Hälfte von „Arbeit im Container“ (sie liegt **in** §9 und bringt keine zweite Datei), die Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (ein Abschnitt desselben Dokuments, zu prüfen, nicht zu schreiben) und die drei Dokument-Risiken (Risiken zählen nach `modul-05-planning-harness.md` §Ziel-Form: Slice nicht mit). **Schichten 1** — Wartungs-Doku `docs/developer/` (Teilung nach `slice-harness-d-check-v0-85` §8: Nutzer- und Wartungs-Doku ist **eine** Schicht; die Planung ist keine; Produkt-Code, Werkzeug, Gegenprobe, Gate-Index und Verträge sind ausgeschlossen, §1). **Schätzung 70 bis 120** in der Metrik von §4, Mitte 95 bei **Marke 250** und Meldeschwelle 180 (38 % der Marke). Damit unter drei Liefer-Punkten und unter zwei Schichten — mit Abstand, und das ist der Zweck des Schnitts: Der Geber stand mit diesem Punkt bei zwei Schichten und 815 bis 1210 über seiner Marke 900 (dessen §4 und §8).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Sendung `Hostschlüssel nach Neubau` einträgt): Der Satz in §5 und der Eintrag in §13 liegen in demselben Dokument und in demselben Liefer-Punkt (*Dokument*); die Sendung legt keinen Liefer-Punkt und keine Schicht an. Liefer-Punkte bleiben **1**, Schichten bleiben **1** (Wartungs-Doku `docs/developer/`; der Handbuchsatz liegt in dieser Schicht, anders als `AGENTS.md` §3.1, das in `slice-harness-rollenlauf-zeitgrenze-regel` steht). Schätzung 75 bis 130 in der Metrik von §4 (Marke 250, Meldeschwelle 180), Mitte 102, 41 % der Marke. Der Ausschluss von §3 bis §5 ist um den einen Satz enger gefasst (§1).
