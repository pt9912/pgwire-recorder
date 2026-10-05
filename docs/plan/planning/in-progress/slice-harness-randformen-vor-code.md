# Slice slice-harness-randformen-vor-code: Randformen vor dem Code entscheiden

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Eingesammelt wird er von der
nächsten Welle-Closure.

**Bezug:** — (Harness-Arbeit; keine Produkt-Anforderung). Herkunft:
`BEO-REPO/spec-randform-erst-im-review-entschieden`, Lese-Schritt der Closure von
welle-extended-query.

**Berührte Spec-Stellen:** —

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

**Ziel:** Der Ablauf von Planner, Architect und Implementer hat einen Schritt
„Randformen vor dem Code entscheiden“: Ein Slice, der einen neuen Vertrag liefert
(Format, Leser, Protokollrand, Diagnose, Option, Gate — die Liste von `AGENTS.md`
§3.10), nennt in §6 dessen Randformen und wo jede entschieden ist; entschieden ist
sie in der Spezifikation oder einer ADR, auch als Entscheidung des Nutzers, bevor der
erste Code-Commit entsteht. Trifft der Implementer auf eine nicht genannte oder offene,
entscheidet er sie nicht still, sondern gibt sie dem Architect zurück.

**Herkunft:** In allen vier Slices von welle-extended-query entschied der Code eine
Randform still, und erst Review oder Verifikation fanden sie, Runde um Runde eine
weitere (`BEO-REPO/spec-randform-erst-im-review-entschieden`, 4×: fehlende und
`null`-Felder, Aliase und Merge-Keys im Format; Herunterfahren; „erwartete Query“
in der Extended-Diagnose; `\v` je Serverversion). Der Gegenbeleg steht in
slice-extended-query-lebendpruefung: Die Randformen, die §6 vor dem Code nannte und
eine ADR vor der ersten Code-Zeile entschied, öffnete kein Review wieder; die nicht
genannten fielen erst dort auf.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die vendored Slice-Vorlage unter `.harness/baseline/` — Bestand bleibt: Der Baum
  ist Kurs-Inhalt, byte-genau gehalten (`make baseline-verify`); der Schritt lebt in
  den repo-eigenen Rollen-Typen und Commands. Braucht er doch eine Vorlagen-Abweichung,
  ist das eine Adaption `MR-<NNN>` in `harness/conventions.md`, kein Edit am Baum.
- Ein Sensor, der die Randformen-Liste prüft — ob eine Randform fehlt, ist Urteil;
  maschinell prüfbar wäre nur, dass §6 eine Liste trägt, und das erzeugte
  Pflichterfüllung (v6.13.0 · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice, „Keine Mindestzahl, und kein Sensor darauf“).
- Randformen bestehender Verträge nachträglich sammeln — ein anderer Vorgang; die
  offenen gehören in die Slices, die die Verträge berühren.
- Produkt-Code, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Rollen-Typen, Commands und `AGENTS.md`. Den Stand seines Register-Eintrags setzt
  die Closure (`implement-slice.md` Schritt 25), nicht die Implementation.
- Die übrigen Rollen-Typen (`planner.md`, `implementer.md`) und die Reviewer-Skill —
  Bestand bleibt: Sie zeigen auf ihre Commands bzw. prüfen gegen `AGENTS.md` §3, und
  die Regel steht nur dort (Pointer statt zweiter Formulierung).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Planner und Architect: `.claude/commands/plan-welle.md` (Slice bereitstellen)
      verlangt die Randformen-Liste in §6 für jeden neuen Vertrag;
      `.claude/agents/architect.md` prüft sie vor dem Code und entscheidet jede an dem
      Ort, den `AGENTS.md` §3.12 nennt, oder legt sie dem Nutzer vor — bestätigt an
      `fed4387` (Verifikation, Punkt 1); `7d9b230` ergänzt Kopf und Eingang des
      Rollen-Typs um die Rückgabe (V-38), ohne Review, siehe §7.
- [x] Implementer: `.claude/commands/implement-slice.md` Schritt 12/13 und
      Randform-Rückgabe — eine nicht genannte oder offene Randform wird als Liste
      „Randform · Frage“ an den Architect zurückgegeben, nicht still entschieden —
      bestätigt an `fed4387` (Verifikation, Punkt 2); `7d9b230` nimmt die Spezifikation
      in die Aufzählung der ausgeschlossenen Entscheidungsorte auf (V-36), ohne Review.
- [x] Regel in `AGENTS.md` §3; sie und jeder Zielort (plan-welle Schritt 6,
      `architect.md`, implement-slice Schritt 12/13) tragen den Herkunfts-Anker
      `seit slice-harness-randformen-vor-code` (oder `seit welle-<Kennung>`, falls eine
      Welle ihn einsammelt). Den Register-Stand setzt die Closure — bestätigt an
      `fed4387` (Verifikation, Punkt 3); die Closure setzt den Stand und trägt den Anker
      auch an der Randform-Rückgabe nach (§7, Drei Paarungen).
- [x] `make gates` grün — an `fed4387` (Verifikation, Punkt 4) und nach `7d9b230`; die
      Closure ändert Planungsdokumente, Register und eine Anker-Zeile.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review bis
      `8fc60d7` (F-364 bis F-371); `fed4387` von der Verifikation geprüft, `7d9b230` ohne
      Review, siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — bei Closure
      geprüft (§7); die nächste Welle-Closure prüft sie erneut.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/commands/plan-welle.md` | update | Schritt 6 (Slices bereitstellen): Randformen-Liste in §6, Zeiger auf `AGENTS.md` §3.12 mit Anker |
| `.claude/agents/architect.md` | update | Kontext-Zuschnitt: Prüfung und Entscheidung der Liste und der Rückgaben vor dem Code, fehlende an den Planner, Zeiger auf §3.12 mit Anker |
| `.claude/commands/implement-slice.md` | update | Schritt 12 (Auslöser vor dem Code, auch für wellenlose Slices), Schritt 13 (beim Implementieren gefunden), Absatz „Randform-Rückgabe“ bei den Rücksprungkanten (Liste „Randform · Frage“ an den Architect, Slice bleibt in `in-progress/`); Herkunfts-Anker in beiden Schritten und im Absatz (der im Absatz mit der Closure) |
| `AGENTS.md` | update | §3.12 als Hard Rule mit Herkunfts-Anker `seit slice-harness-randformen-vor-code`; Prüfweg ist Architect vor dem Code und Review gegen §3 |
| `docs/plan/planning/observations/BEO-REPO/spec-randform-erst-im-review-entschieden/state.md` | update bei Closure | Ausgang verkörpert (Schritt 25, nicht in der Implementation) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): jederzeit; keine Abhängigkeit. Er wirkt erst auf
Slices, die nach ihm geplant werden; vor dem ersten Slice von welle-replay-semantik
geliefert, trifft er deren Pläne.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Schritt verlangt eine
  Abweichung von der vendored Vorlage (dann zuerst eine Adaption, getrennt).
- `in-progress` → `open` (blockiert — Carveout?): keine Bedingung absehbar; der
  Slice ändert nur repo-eigene Prozess-Texte.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Schritt bläht Pläne auf: Eine Liste möglicher Randformen ohne Entscheidung ist
  Text, keine Grenze. Gegenmittel ist die Pflicht, jede genannte zu entscheiden, nicht
  die Liste lang zu machen; der Architect hält sie nach seiner Regel „Lösungen, nicht
  neue Hürden“ klein — **Ausgang:** entfallen: Die Form, die das Risiko benennt, eine
  Liste ohne Entscheidung, kommt nicht mehr bis zum Code. §3.12 verlangt jede genannte
  Randform vor dem ersten Code-Commit in der Spezifikation oder einer ADR entschieden,
  und eine offene hält Schritt 12 von `implement-slice.md` als Rückgabe an den Architect
  an (Verifikation, Punkt 2). Eine längere Liste entschiedener Randformen ist kein
  Aufblähen, sondern Vertragsinhalt; jede kostet einen Satz am Entscheidungsort, und die
  Liste gilt nur für neue Verträge (§3.10). Dass der HIGH-Anker der Reviewer-Skill zu
  vorsorglich langen Listen drängen kann (F-371), steht in §7 beim Lerneintrag.
- Was „neuer Vertrag“ heißt, ist Urteil; zu eng gefasst, greift der Schritt bei einer
  Diagnose-Änderung nicht (die Klasse von slice-extended-query-replay). §3.12 bindet den
  Begriff an die Liste von §3.10, die Diagnose nennt — **Ausgang:** entfallen: §3.12
  übernimmt die Vertragsliste von §3.10 (Protokollrand, Format, Leser, Diagnose, Option,
  Gate), statt eine eigene zu führen; der Fall von slice-extended-query-replay, eine
  Diagnose, liegt darin. Die Lücke, die das Review im Plan fand (F-369, „Gate“ fehlte
  im Ziel), ist in `fed4387` geschlossen (Verifikation, Abschnitt 2). Ob ein Vertrag
  *neu* ist, bleibt Urteil des Planners; das ist keine Lücke der Regel, sondern derselbe
  Begriff, auf den §3.10 seit welle-extended-query baut.
- Ein wellenloser Slice wird nicht über `plan-welle` bereitgestellt; Schritt 6 erreicht
  ihn nicht. Den Auslöser trägt implement-slice Schritt 12: Ist eine Randform in §6
  nicht entschieden, geht sie vor dem Code als Rückgabe an den Architect. Eine
  Randform, die §6 gar nicht nennt und die der Implementer nicht bemerkt, findet erst
  das Review — **Ausgang:** weiter offen → Register
  `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×). Rest nach F-367 und V-37:
  Kein Command-Schritt bringt einen wellenlosen Slice vor dem Code zum Architect;
  Schritt 12 greift nur bei einer Randform, die §6 nennt und offen lässt, Schritt 13 nur,
  wenn der Implementer sie bemerkt. *Entfallen* wäre falsch, denn das Risiko kann
  eintreten: Am Fall F-350 hätte auch bei einem Slice einer Welle allein Schritt 13
  gegriffen (Verifikation, Abschnitt 3). *Eingetreten* ist es nicht, denn seit §3.12
  lief kein wellenloser Slice mit neuem Vertrag; V-37 ist ein Befund am Text, kein
  Auftreten. Ein Folge-Slice jetzt hieße einen Architect-Schritt für jeden wellenlosen
  Slice, auch ohne neuen Vertrag, ohne einen einzigen Beleg, dass er fehlt — eine Hürde
  vor der Evidenz. Der Zähler entscheidet: Findet ein Review bei einem wellenlosen Slice
  eine Randform, die §6 nicht nannte, ist das eine Evidenz hier und zugleich in
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (Retirement-Check von §3.12).

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

- **Was hat funktioniert:** Der Ablauf hat an allen drei Stellen einen Träger: `plan-welle` Schritt 6 verlangt die Randformen-Liste in §6, `architect.md` prüft und entscheidet sie vor dem Code oder legt sie dem Nutzer vor, und `implement-slice` hält bei einer offenen (Schritt 12) oder nicht genannten Randform (Schritt 13) an und gibt sie als Liste *Randform · Frage* an den Architect zurück. Entschieden wird nur in der Spezifikation oder einer ADR; §6 zeigt darauf. Die Verifikation hat die drei Liefer-Punkte an `fed4387` bestätigt und den Ablauf an F-350 (`\v` je Serverversion, slice-extended-query-lebendpruefung) durchgespielt: Er hätte die Randform vor dem ersten Code-Commit zum Architect gebracht, aber nur über Schritt 13 und die Randform-Rückgabe. Planungs-Hälfte und Schritt 12 hätten nicht gegriffen, weil §6 die Randform beim Eintritt nicht nannte (Verifikation, Abschnitt 3). Der Ablauf hängt damit an der Stelle, die auch wellenlose Slices erreicht, und dort daran, dass der Implementer die Frage bemerkt.
- **Was ging anders als geplant:** Die erste Fassung (`8fc60d7`) hatte für die Rückgabe weder Empfänger noch Übergabe-Artefakt; sie lief in die Plan-Defekt-Schleife des Implementers zurück, der die Randform dann selbst in §6 hätte entscheiden können (F-364, MEDIUM). Der Plan führte den Register-Stand als Liefer-Punkt, gegen Schritt 25 von `implement-slice.md` (F-366), drei Zielorte trugen den Anker nicht (F-365), und Risiko 3 sagte für wellenlose Slices eine Architect-Prüfung zu, die kein Schritt trägt (F-367). `fed4387` hat F-364 bis F-370 abgearbeitet, F-367 teilweise; den Rest führt Risiko 3. Die Verifikation fand drei Reste: V-36 (die Randform-Rückgabe nannte die Spezifikation nicht als ausgeschlossenen Entscheidungsort, obwohl der Implementer im F-350-Fall genau dort die Leerraum-Menge ohne Version festlegte) und V-38 (Kopf und Eingang von `architect.md`, Artefakt-Liste von `implement-slice.md` ohne die neue Kante) hat `7d9b230` nachgezogen; dieser Commit hat kein Review und keine Verifikation, nur `make gates` grün. V-37 ist Risiko 3. Zwei Prüfrunden und eine Nacharbeit statt einer Runde.
  - **V-39, ohne Aktion:** Wer §6 nach der Entscheidung auf den Ort zeigen lässt, sagt kein neuer Text. Bei einer Änderung der Spezifikation trägt es §3.9; bei einer ADR nennt §3.9 sie nicht wörtlich, und der Brauch ist, dass der Architect den Plan im ADR-Commit nachzieht (`f906c04`). Hält der Brauch nicht, zeigt es sich als `BEO-REPO/plan-folgt-korrektur-nicht`.
- **Steering-Loop-Eintrag:** Guide geschärft: Ein Slice, der einen neuen Vertrag liefert (§3.10), nennt in §6 dessen Randformen und wo jede entschieden ist; entschieden wird vor dem ersten Code-Commit in der Spezifikation oder einer ADR, und eine nicht genannte oder offene Randform gibt der Implementer an den Architect zurück, statt sie zu entscheiden
  — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-REPO/spec-randform-erst-im-review-entschieden` (slice-extended-query-modell, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 4×). Weitere Zielorte: `.claude/commands/plan-welle.md` Schritt 6, `.claude/agents/architect.md`, `.claude/commands/implement-slice.md` Schritte 12 und 13 und Randform-Rückgabe.
  - **Proportionalität am Review (F-371):** Die Reviewer-Skill bleibt unverändert; ihr HIGH-Anker „ADR-Verstoß (Layer, Tool, Hard Rule)“ fasst einen Verstoß gegen §3.12. Dem Wortlaut nach ist damit jede Randform, die ein Diff entscheidet, ohne dass §6 sie nennt, ein HIGH, gleich wie selten oder folgenlos sie ist. Das ist gewollt, weil die Regel genau diese stille Entscheidung verhindern soll; es kann aber Reviews dazu bringen, gegen den Wortlaut herabzustufen, oder Pläne dazu, vorsorglich lange Listen zu führen (§6, Risiko 1). Nicht eingetragen, weil kein Fehlermuster vorliegt. Wer §3.12 lockern oder die Skill nach Wirkung staffeln will, prüft zuerst im Retirement-Check, ob Reviews seit slice-harness-randformen-vor-code solche HIGH vergeben oder herabgestuft haben.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` — Stand `verkörpert`, Zielort und Herkunfts-Anker in `state.md`; Zähler unverändert **4×**. Keine neue Evidenz: Dieser Slice liefert keinen Vertrag im Sinn von §3.10 (Prozess-Texte, kein Format, Leser, Protokollrand, Diagnose, Option oder Gate), und keiner von F-364 bis F-369 ist eine Randform, die Code still entschied. Am nächsten liegt F-364, die fehlende Empfänger-Kante; sie ist eine Lücke im Ablauf-Text, die das Review fand, keine Randform eines Vertrags.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11 seit welle-extended-query) — Beleg `evidence/slice-harness-randformen-vor-code.md`: Risiko 3 im Plan und der Satz „Der Architect prüft die Liste vor dem Code“ in §3.12 sagten für wellenlose Slices eine Prüfung zu, die kein Schritt trägt (F-367, Rest V-37); Zuordnung laut Review-Summary. **7×**. Retirement-Check von §3.11: wieder aufgetreten, die Regel bleibt.
  - `BEO-REPO/randform-wellenlos-ohne-architect-vor-code/` — neu, aus Risiko 3 (weiter offen), Beleg `evidence/slice-harness-randformen-vor-code.md` (F-367, V-37), **1×**, `offen`.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` — kein Auftreten: F-369 (Plan-Ziel enger als die Regel) entstand im selben Commit wie die Regel, nicht nach einer Korrektur, und §3.9 greift laut Review im Wortlaut nicht; die Verifikation fand Plan und Diff deckungsgleich (Abschnitt 4).

  Einmalig und nicht eingetragen: F-365 (Zielorte ohne Anker), F-366 (Register-Stand vor der Closure gesetzt; der Plan selbst verlangte es), F-368 (Handlung ohne Handelnden im Rollen-Typ), F-369, F-370 und F-371 (kein Fehlermuster), V-36 und V-38 (behoben in `7d9b230`), V-39 (oben). Über der Schwelle stehen außer dem verkörperten Eintrag dieses Slice nur verkörperte Einträge (`negativtests-fehlen-bei-neuem-vertrag` 6×, `plan-folgt-korrektur-nicht` 6× mit Sensor geplant in `slice-harness-kopf-sensor`, `zusage-im-kommentar-weiter-als-pruefung` 7×); den Lese-Schritt für sie führt die nächste Welle-Closure.
- **Folge-Slices:** keine. Risiko 3 geht ins Register, nicht in einen Slice (Begründung in §6).
- **Risiken aus §6:** drei, jedes mit genau einem Ausgang — zwei entfallen (Risiko 1 und 2), eines weiter offen → `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (Risiko 3); siehe §6.
- **Drei Paarungen:** Anker — `liegt in` nennt `AGENTS.md §3.12`; `grep -n "seit slice-harness-randformen-vor-code"` findet ihn in der Überschrift von §3.12 und an jedem weiteren Zielort: `plan-welle.md` Schritt 6, `architect.md` (Abschnitt *Randformen eines neuen Vertrags*), `implement-slice.md` Schritt 12, Schritt 13 und Randform-Rückgabe (dort mit der Closure nachgetragen). Folge-Slice — keiner genannt; `slice-harness-kopf-sensor` oben ist Zustand eines anderen Eintrags und liegt als Datei im Planning-Lifecycle. Register — die vier genannten Kennungen bestehen als Verzeichnis mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-harness-randformen-vor-code.md` (bis `8fc60d7`; F-364 bis F-371), Verifikation `docs/reviews/2026-10-05-verifikation-slice-harness-randformen-vor-code.md` (bis `fed4387`; V-36 bis V-39; DoD 1 bis 5 bestätigt, `make gates` grün an `fed4387`), Nacharbeit `7d9b230` (V-36, V-38; `make gates` grün, kein Review). Validierung: n/a, der Slice ändert Prozess-Texte und liefert keinen End-Nutzer-Wert.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register am Stand der Closure von
welle-extended-query gesichtet. Treffer: `BEO-REPO/spec-randform-erst-im-review-entschieden`
(4×, geplant mit diesem Slice) ist der Gegenstand; `BEO-REPO/plan-folgt-korrektur-nicht`
(6×, verkörpert, Sensor geplant in slice-harness-kopf-sensor) trifft jede Plan-Änderung
dieses Slice.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
