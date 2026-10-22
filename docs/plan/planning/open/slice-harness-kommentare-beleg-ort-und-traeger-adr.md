# Slice slice-harness-kommentare-beleg-ort-und-traeger-adr: Die Kommentare von fünf Werkzeug-Dateien sagen, was da ist: der Ort des Belegs eines Nachweises und die Herkunft der Kennungen des Bootstrap-Werkzeugs in den Träger-Dateien

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — **Empfehlung des Planners vom 2026-10-22** (Auftrag des Koordinators), nicht Entscheidung des Nutzers: Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist (`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Geprüft und nicht gewählt: `welle-dev-container` (zwei der fünf Dateien gehören dem Container-Werkzeug, die übrigen dem Release und dem Träger; ein Slice, der nur Kommentare ändert, dehnte ihren Closure-Trigger) und `welle-harness-regeln` (Regeln, hier steht keine). **Stelle:** vor `slice-harness-rollenlauf-zeitgrenze` (Empfehlung, kein Start-Trigger): Dessen Gegenprobe darf Code von `container-einrichten-gegenprobe.sh` übernehmen, und ein veralteter Kommentar würde mit übernommen. Entschieden hat der Nutzer die **Sache** am 2026-10-22: dass die Kommentare ohne Adresse einen Slice bekommen („ja“ auf die Empfehlungen des Koordinators, so übermittelt). Angelegt in `open/`, weil er nicht priorisiert ist.

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice ändert Kommentare von Werkzeugen des Harness. Regelstellen: `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist; die Herkunft steht in **einem** auflösbaren Feld) und §3.11 (ein Kommentar sagt nur zu, was ein Test oder Gate prüft). Bindung an das Register: `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` (1×; die Sendung des Architect nennt die Kommentare mit „§7“ als Fundstelle ohne Adresse), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (91×, `verkörpert`). Keine aktive ADR.

**Berührte Spec-Stellen:** [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18 (gelesen: Ort des Belegs eines Nachweises; die Kommentare zeigen dorthin). Der Slice ändert keinen Vertrag. Genannt, ohne berührt zu werden: [LH-QA-02](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) und [LH-QA-04](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) (in diesem Repo andere Anforderungen als die gleichnamigen Verweise des Bootstrap-Werkzeugs) und [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (Verweise in `container-einrichten-gegenprobe.sh`, die stimmen).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-22; nachgeführt am 2026-10-22 (Entscheidung zu den sechs Meldungen, §1).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Kommentare von **fünf** Dateien sagen, was da ist (`AGENTS.md` §3.7). (1) **Der Ort des Belegs eines Nachweises:** `harness/mk/homebrew-nachweis.mk` (Zeilen 5 und 19), `tools/harness/homebrew-nachweis.sh` (26) und `tools/dev/harness/container-einrichten-gegenprobe.sh` (16) sagen „der Beleg steht in §7 des Slice-Plans“. Seit der Entscheidung des Nutzers vom 2026-10-22 liegt der Beleg eines Nachweises als eigene Datei unter `docs/reviews/`, und §7 nennt Pfad, Summe, Exit-Zeile und Kopf ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6), [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18). Der Kommentar nennt den Ort wie der Vertrag und zeigt auf ihn. (2) **Die Herkunft der Kennungen des Bootstrap-Werkzeugs:** `harness/mk/traeger.mk` und `tools/harness/traeger-fetch.sh` tragen in Kommentaren Verweise auf ADR 0058 und 0059, `LH-QA-02` bis `LH-QA-04` und `MR-007`. Sie meinen Entscheidungen und Anforderungen von `ai-harness-init`, dem Werkzeug, das die Dateien emittiert; in diesem Repo ist [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) der Zeitwächter, eine ADR 0059 gibt es nicht, [LH-QA-02](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung) bis [LH-QA-04](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) sind andere Anforderungen (geringe Eingriffe, Portabilität, Automatisierbarkeit), und `MR-007` gibt es nicht. Ein Leser löst die Kennung hier auf und liest etwas anderes. Der Kommentar nennt die Herkunft: den Namen des Werkzeugs vor der Kennung (Entwurf, §3), oder er verzichtet auf die Kennung.

**Funde über die Liste der Koordination hinaus.** Die Architect-Liste nannte für den Träger ADR 0058 und 0059 in `traeger.mk` (4 Zeilen) und `traeger-fetch.sh` (7 Zeilen). Die Suche des Planners am 2026-10-22 (`grep -n "ADR-005[89]\|LH-QA-0[234]\|MR-007"`) fand mehr: Kommentarzeilen in `traeger.mk` 7, 13, 14, 15, 18, 20, 23 und in `traeger-fetch.sh` 3, 4, 7, 8, 10, 12, 19, 21, 24, 33, 50, 87, 121; dazu **sechs Zeilen, die keine Kommentare sind** (§1, *Ausdrücklich NICHT*). Maßgeblich ist die Suche zu Beginn des Slice; die Zeilennummern sind der Stand vom 2026-10-22.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **sechs Meldungen** von `tools/harness/traeger-fetch.sh` (Zeilen 43, 59, 67, 96, 101 und 130: `echo … >&2` mit ADR 0058 oder 0059, `LH-QA-02` oder `LH-QA-04` im Text) — Klasse 2 (Bestand bleibt bewusst stehen): Eine Meldung ist die Ausgabe eines Werkzeugs, kein Kommentar; sie zu ändern ist eine Änderung des Verhaltens (Text auf stderr) einer Datei, die der Bootstrap emittiert und die keine Gegenprobe prüft. Dieser Slice ändert **keinen Code** (DoD). **Entschieden (Entscheidung des Nutzers vom 2026-10-22, wörtlich „ja nächster slice“ auf die Vorlage mit den Empfehlungen des Koordinators, so übermittelt):** Die Meldungen bleiben stehen, bis der Bootstrap sie ändert (§6, Randform 1). Grund: Die Datei ist eine Kopie des Werkzeugs, kein Eingriff der Repo-Pflege. Eine Adresse gibt es nicht und braucht es nicht.
- Die **Träger-Dateien an der Quelle** (`ai-harness-init`, Kurs-Repo) — Klasse 3 (anderer Vorgang): Dieser Slice ändert die Kopie in diesem Repo; was der Bootstrap schreibt, entscheidet das Werkzeug.
- **Code und Verhalten** der fünf Dateien und die **Kennungen, die stimmen** (`SPEC-060`, `SPEC-067`, `SPEC-069`, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) und [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) in `container-einrichten-gegenprobe.sh`, `homebrew-nachweis.sh` und `homebrew-nachweis.mk`) — Klasse 2: Sie zeigen auf Dateien dieses Repos und lösen auf.
- Der **Ort des Belegs in den Anweisungen der Rollen und in `harness/README.md`** — Klasse 1: `slice-harness-regeln-implementer-uebergabe` (Anweisungen) und `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (die Zeile `make homebrew-nachweis`); beide nehmen die Sendung in ihrem §1 an. Grund: andere Dateien, andere Schicht.
- **Kommentare in anderen Dateien** — Bestand bleibt bewusst stehen: Die Suche nach ADR 0058 und 0059 in `tools`, `harness` (außer den Verträgen), `.claude` und `Makefile` fand außer den zwei Träger-Dateien nichts.
- Eine **Prüfung** (Sensor), die fremde Kennungen in Kommentaren meldet — anderer Vorgang (Werkzeug statt Kommentar): mit Vertrag unter `harness/vertraege/` und Gegenprobe; eine Adresse gibt es nicht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Ort des Belegs):** Die Kommentare in `harness/mk/homebrew-nachweis.mk` (Zeilen 5 und 19), `tools/harness/homebrew-nachweis.sh` (26) und `tools/dev/harness/container-einrichten-gegenprobe.sh` (16) nennen den Ort des Belegs eines Nachweises nach dem Vertrag (eine Datei unter `docs/reviews/`, §7 nennt Pfad und Summe) mit dem Zeiger auf [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) beziehungsweise [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18; kein Kommentar sagt mehr zu, als der Vertrag sagt (`AGENTS.md` §3.11), namentlich nicht, dass `make archive-welle` die Datei archiviert (nicht zugesagt).
- [ ] **Liefer-Punkt 2 (Herkunft der Kennungen des Werkzeugs):** In `harness/mk/traeger.mk` und `tools/harness/traeger-fetch.sh` nennt jede Kommentarzeile, die ADR 0058, ADR 0059, eine Anforderung des Werkzeugs (`LH-QA` mit Nummer 02 bis 04) oder die Kennung `MR` mit Nummer 007 als Verweis trägt, den Namen des Werkzeugs `ai-harness-init` vor der Kennung oder verzichtet auf die Kennung. Beleg in §7: die Suche aus §1 vor und nach der Änderung; nach ihr trifft sie in Kommentarzeilen nur Zeilen mit dem Namen des Werkzeugs.
- [ ] **Kein Code:** Der Diff ändert allein Kommentarzeilen (Beleg in §7: `git diff -U0 <Basis>..HEAD -- harness tools` mit allen Änderungszeilen, die nicht mit `#` nach Leerraum beginnen — das Ergebnis ist leer); die Meldungen von `traeger-fetch.sh` und jede Zeile Code bleiben byte-gleich.
- [ ] `make abdeckung-check` grün (der Kopf von `container-einrichten-gegenprobe.sh` trägt Abdeckungs-Deklarationen, und Zeile 16 liegt im Kopf), `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (der Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern). Eine Gegenprobe fährt dieser Slice nicht: Es ändern sich Kommentare, und keine Gegenprobe prüft einen (Mutationen sind Textersatz von Code-Zeilen); der Bericht nennt das.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: bei `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` der Vermerk, dass die Kommentare ihren Ort nennen (kein `verkörpert`: der Zähler steht bei 1×); für die Klasse „Kennung eines emittierenden Werkzeugs im Kommentar meint etwas anderes als dieselbe Kennung im Repo“ (§8: kein Eintrag trägt sie) ein neuer Eintrag mit Beleg dieses Slice, sofern der Planner bei der Closure die Klasse bestätigt; sonst „keine Beobachtung angefallen“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/homebrew-nachweis.mk` (Zeilen 5, 19) | update (Kommentar) | Liefer-Punkt 1: Ort des Belegs |
| `tools/harness/homebrew-nachweis.sh` (Zeile 26) | update (Kommentar) | Liefer-Punkt 1 |
| `tools/dev/harness/container-einrichten-gegenprobe.sh` (Zeile 16, im Kopf) | update (Kommentar) | Liefer-Punkt 1 |
| `harness/mk/traeger.mk` (sieben Kommentarzeilen) | update (Kommentar) | Liefer-Punkt 2: Herkunft der Kennungen; Datei emittiert von `ai-harness-init` |
| `tools/harness/traeger-fetch.sh` (dreizehn Kommentarzeilen) | update (Kommentar) | Liefer-Punkt 2; die sechs Meldungen bleiben (§1, entschieden am 2026-10-22) |
| `docs/plan/planning/observations/BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen/` | update bei Closure | Vermerk (kein Liefer-Punkt) |

Tests und Gates: keine neuen. Der Slice liefert Kommentare (`AGENTS.md` §3.10 verlangt Mutationen für Verträge von Code und Gates); die Prüfung ist der Diff (DoD, *Kein Code*).

**Entwurf der Form** (Empfehlung des Planners; der Implementer entscheidet keine Form selbst, er gibt eine andere an den Planner zurück): *Herkunft:* „(`ai-harness-init`, ADR 0058 Festlegung 1)“ statt „(ADR 0058 Festlegung 1)“ ohne Namen; für die Anforderungen des Werkzeugs „(`ai-harness-init`: Anforderung zur Portabilität)“ statt `LH-QA-03`, oder die Kennung entfällt, wo der Satz ohne sie trägt. *Ort des Belegs:* „der Beleg liegt als Datei unter `docs/reviews/`, §7 des Slice-Plans nennt Pfad und Summe (`SPEC-067` Punkt 20 (6))“; in der Gegenprobe „`SPEC-069` Punkt 18“.

**Reihenfolge:** (1) die Suche aus §1 am Anfang des Slice, Treffer in §7; (2) Liefer-Punkt 1; (3) Liefer-Punkt 2; (4) die Suche am Ende und der Diff gegen *Kein Code*. Jeder Commit nennt die Kennung dieses Slice.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Kein anderer Slice ist Start-Trigger. Reihenfolge: **vor** `slice-harness-rollenlauf-zeitgrenze` (**Empfehlung des Planners**, siehe Kopf). WIP-Limit 1.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`; wiederholt dieser Plan ihn, gilt der von dort): Summe aus Einfügungen und Löschungen über die Pfade `harness` und `tools`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Marke 200, Meldeschwelle 160.** Schätzung des Planners (2026-10-22, in dieser Metrik): Etwa 26 geänderte Kommentarzeilen (Ort des Belegs 6, Herkunft 20), je Zeile eine Einfügung und eine Löschung, mit Umbruch bis zu 100; **52 bis 100, das sind 26 bis 50 % der Marke**. Der Faktor Tests zu Produktcode entfällt: kein Code, keine Gegenprobe.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 200, oder die Suche am Anfang findet Kennungen in einer sechsten Datei.
- `in-progress` → `open` (blockiert): Der Nutzer nimmt die Entscheidung vom 2026-10-22 zurück und entscheidet, dass die Meldungen von `traeger-fetch.sh` mitzuändern sind (Verhalten einer emittierten Datei): Das ist ein anderer Schnitt (Code, Ausgabetext), und der Planner plant ihn neu.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, die Suche aus §1 am Ende in §7 mit dem Ergebnis, der Diff nach *Kein Code* leer in §7, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel entfällt, wenn nichts verkörpert wird; dann ist der Eintrag gezählt).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber Änderungen an Dateien, deren Ränder der Planner sonst im Einzelfall entscheidet):

1. **Die Meldungen von `traeger-fetch.sh`** (sechs Zeilen mit ADR 0058 oder 0059, `LH-QA-02` oder `LH-QA-04` in der Ausgabe) — **Entschieden (Entscheidung des Nutzers vom 2026-10-22, wörtlich „ja nächster slice“ auf die Vorlage mit den Empfehlungen des Koordinators, so übermittelt):** Sie bleiben stehen, bis der Bootstrap sie ändert. Eine Änderung änderte Text auf stderr einer vom Bootstrap emittierten Datei, ohne dass eine Gegenprobe ihn prüft; bleiben sie, führt eine Meldung weiter Kennungen, die hier anderes meinen (akzeptiert). Grund: Die Datei ist eine Kopie des Werkzeugs, kein Eingriff der Repo-Pflege.
2. **Der Träger wird vom Bootstrap emittiert.** Ein erneuter Lauf von `ai-harness-init` schreibt `harness/mk/traeger.mk` und `tools/harness/traeger-fetch.sh` neu, und der Kommentar mit der Herkunft **fehlt dann wieder**. Akzeptiertes Negativ; Grund: Die Dateien gehören dem Werkzeug, und dieses Repo sperrt den Bootstrap nicht; die Änderung dieses Slice hält bis zum nächsten Lauf. Gegenmittel: keines im Repo; der Befund gehört in das Register und in einen Änderungswunsch an das Werkzeug (nicht Gegenstand dieses Slice).
3. **Form der Herkunft:** Name des Werkzeugs vor der Kennung (Entwurf in §3); der Implementer ändert die Form nicht selbst.
4. **Die Suche ist nicht abschließend**, wenn eine Kennung in einer anderen Schreibweise steht (zum Beispiel „ADR 58“); die Suche der DoD ist der Maßstab, kein Beweis der Vollständigkeit.
5. **Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht:** nicht berührt (`AGENTS.md` §3.12); keine Zusage ändert sich.

**Risiken:**

- Die Suche am Anfang findet mehr Kennungen als §1 nennt — **Ausgang:** entfallen, wenn die Treffer in §7 stehen und der Diff in Marke und Schicht bleibt; sonst Rückführung nach `next` (§4).
- Ein neuer Kommentar sagt mehr zu, als der Vertrag oder ein Test hält (`AGENTS.md` §3.11) — **Ausgang:** entfallen, wenn jeder neue Satz auf den Vertrag zeigt, statt ihn zu wiederholen.
- Ein erneuter Bootstrap schreibt die Träger-Dateien zurück (Randform 2) — **Ausgang:** weiter offen → der neue Eintrag der Klasse „fremde Kennung im emittierten Kommentar“ (DoD, Register), mit dem Wiederauftreten als Beleg.
- Die Meldungen tragen weiter fremde Kennungen (Randform 1) — **Ausgang:** weiter offen → derselbe Eintrag; die Entscheidung des Nutzers (bleiben, 2026-10-22) steht in §1 und in Randform 1.
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 52 bis 100 gegen Marke 200) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: die Suche aus §1 vor und nach der Änderung mit ihrem Ergebnis, der Diff nach *Kein Code* (Befehl und leeres Ergebnis), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) mit dem Ergebnis der Schätzung gegen die Messung, der Ausgang des Registers und die Ausgänge der Risiken.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` (1×, unter der Schwelle) — der Auslöser für Liefer-Punkt 1; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (91×, `verkörpert` in `AGENTS.md` §3.11) — ein Kommentar sagt nur zu, was ein Test oder Gate prüft; `BEO-REPO/konjunktiv-im-kommentar-bei-kopplung` (1×) — Kommentare über Zustand statt Chronik. **Kein vorhandener Eintrag trägt die Klasse von Liefer-Punkt 2** (die Kennung eines emittierenden Werkzeugs im Kommentar meint etwas anderes als dieselbe Kennung in diesem Repo); das ist die Antwort der Sichtung und der Grund für den neuen Eintrag in der DoD. Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt):** Liefer-Punkte 2 (Ort des Belegs | Herkunft der Kennungen des Werkzeugs). Schichten 1: Werkzeuge und Fragmente (`tools/harness/`, `harness/mk/`); die Planung zählt keine Schicht, kein Produkt-Code. Schätzung 52 bis 100 in der Metrik von §4 (Marke 200, Meldeschwelle 160); der Faktor Tests zu Produktcode entfällt (kein Code, keine Gegenprobe). Der Slice trägt keine Sendung in einen anderen Nehmer ein und nimmt die Sendung des Architect zum Ort des Belegs für die Kommentare an (Register `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen`).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (2) und Schichten (1) unverändert.
