# Slice slice-harness-vertraege-spezifikation: Verträge der Harness-Werkzeuge in der Spezifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach der Antwort des Kurs-Repos
ai-harness-course auf den Change Request „Spezifikations-Ort für Verträge der
Harness-Werkzeuge“ (2026-10-06): nach `slice-lastenheft-pruefbarkeit` und vor
`slice-harness-lint` (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) (die Gates laufen ohne Eingaben), [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Nachweis in den Abdeckungstabellen). Bindung an Entscheidungen: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) (Kopf-Sensor), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckung je Anforderung und Pfad), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Gate-Nachweise und geteilte Messung).

**Berührte Spec-Stellen:** `spezifikation.md` (neuer Abschnitt für Harness-Werkzeuge; Ort und Kennungen nach §6) · `spezifikation.md §11`

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Verträge der zwei Harness-Werkzeuge, die heute ohne Spezifikationsstelle
leben, stehen in einem Abschnitt der Spezifikation für Harness-Werkzeuge
(Technik-Stratum, fortschreibbar) mit eigenen Kennungen: `kopf-check` mit den neun
Punkten von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und den sieben
Lesarten, die der Architect beim Liefern bestätigt hat (heute im ZUSAGE-Kopf von
`tools/harness/kopf-check.sh`; Herkunft: §6 des archivierten
`slice-harness-kopf-sensor` in `docs/plan/planning/done/welle-replay-semantik/archiv.zip`),
und `abdeckung` mit [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und
dem Kopf von `tools/test/abdeckung.sh`. Wie ein Lauf zu lesen ist (Vertrag, Grenze,
Ausgabe, Ausgänge, Sperren), steht je Gate in einer Sensor-Datei
`harness/sensors/<target>.md` nach der Baseline-Vorlage, die den Vertrag verlinkt und
nichts neu entscheidet; `harness/README.md` §Sensors zeigt darauf. Die Skriptköpfe
behalten nur, womit das Werkzeug selbst geprüft ist. Danach entscheiden die
Folge-Slices ihre Randformen in diesem Abschnitt.

**Herkunft:** `BEO-REPO/harness-lesart-ohne-entscheidungsort` (Review F-380 und
Verifikation V-42 in `slice-harness-kopf-sensor`, Review F-406 in
`slice-lastenheft-pruefbarkeit`). Das Kurs-Repo hat den Change Request am 2026-10-06
beantwortet: Der Ort ist das Technik-Stratum, die Gate-ADR verweist mit `Schärft:`
darauf, Randformen nach `Accepted` werden dort ohne Folge-ADR fortgeschrieben; die
Lesart eines Laufs steht in der Sensor-Datei, die Prüfung des Werkzeugs bei ihm. Der
Ausgang „urteilsgebunden“ ist abgelehnt: Wo kein Sensor möglich ist, steht eine
Prosa-Begründung dafür. Der Adaptions-Eintrag
[MR-001](../../../../harness/conventions.md#mr-001), der den Ort bis dahin in den
Skriptkopf legte, ist damit aufgelöst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Verträge und Sensor-Dateien der übrigen Gates (`make a-check`, `make
  a-check-negativ`, `make hook-gegenprobe`, `make baseline-verify`, `make build`,
  `make test`, `make test-integration`, `make docs-check`) — Bestand bleibt bewusst
  stehen: Ihre Zeile in `harness/README.md` §Sensors trägt den Vertrag in einem Satz,
  und eine Sensor-Datei entsteht nach der Baseline-Vorlage erst, wenn ein Target mehr
  braucht. Ein Gate, dessen Randform später entschieden werden muss, bekommt seine
  Stelle in dem Slice, der sie entscheidet.
- Die Erweiterung des Abdeckungs-Vertrags nach
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Nachweisart Gate,
  geteilte Messung) — ein Folge-Slice übernimmt sie: `slice-harness-abdeckung-gate`
  baut sie und schreibt ihre Randformen in diesen Abschnitt (Empfehlung des Planners;
  entschieden wird es mit §6 *Umfang des Abdeckungs-Vertrags*).
- Die Verträge der neuen Gates (Lint, Coverage, Mutation, Black-Box-Stufen) — ein
  Folge-Slice übernimmt sie: `slice-harness-lint`, `slice-harness-coverage`,
  `slice-harness-mutation`; die vier Umstellungs-Slices folgen den Entscheidungen von
  `slice-harness-lint`. Dieser Slice legt nur den Abschnitt an, in den sie schreiben.
- Verhaltensänderung an `kopf-check` oder `abdeckung` — ein anderer Vorgang: Der Slice
  verlegt Verträge, er ändert sie nicht; eine Abweichung zwischen Skript und
  Spezifikationstext, die beim Übertragen auffällt, wird als Randform an den Architect
  gegeben, nicht im Skript berichtigt. Auch die offenen Randformen aus
  `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (nicht lesbares Lifecycle-Verzeichnis,
  Symlink als Plan) entscheidet dieser Slice nicht.
- Lastenheft, Sicht und Produkt-Code — Schicht-Abgrenzung: Der Slice ändert die
  Spezifikation (neuer Abschnitt, Historie), `harness/` (Sensor-Dateien, §Sensors) und
  die Kommentarköpfe zweier Skripte unter `tools/`; kein ausführbarer Code.
- Die angenommenen ADRs — [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md),
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) bleiben unverändert
  (`AGENTS.md` §3.5); wie ihr `Schärft: —` zum neuen Abschnitt steht, ist §6
  *Schärft-Bezug*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Spezifikation: Der Abschnitt für Harness-Werkzeuge liegt am Ort aus §6 und
      führt je Werkzeug eine Kennung nach §6; `kopf-check` mit den neun Punkten von
      [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und den Lesarten
      (1) bis (5), (a) und (b), `abdeckung` mit dem Vertrag aus
      [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und dem heutigen
      Skriptkopf (Deklarationsform, Pfade, Nachweisart nach Ort, erzeugte Tabellen,
      `--check`, Ausgänge). Jede Zusage des Abschnitts ist eine Zusage, die das Skript
      heute einhält: Zu jeder nennt der Skriptkopf den Fall der Gegenprobe
      (`AGENTS.md` §3.11), eine Zusage ohne Fall steht als Grenze im Abschnitt oder
      geht als Randform an den Architect. §11 trägt eine Historie-Zeile.
- [ ] Sensor-Dateien `harness/sensors/kopf-check.md` und
      `harness/sensors/abdeckung-check.md` (Satz nach §6 *Welche Targets*) per `cp` aus
      der vendored Vorlage, je mit Vertrag (Link auf den Abschnitt), Grenze, Ausgabe
      und Ausgängen, Sperren und Bindung; sie entscheiden nichts, was der Abschnitt
      nicht sagt. Die Target-Zellen in `harness/README.md` §Sensors verlinken sie.
- [ ] Skriptköpfe von `tools/harness/kopf-check.sh` und `tools/test/abdeckung.sh`
      tragen statt der Vertragspunkte den Verweis auf den Abschnitt und je Punkt nur
      noch, womit das Werkzeug geprüft ist (Fall der Gegenprobe); kein Vertragspunkt
      geht dabei verloren (Abgleich Punkt für Punkt im Bericht). Ausführbarer Code
      unverändert: `git diff` der beiden Skripte berührt nur Kommentarzeilen,
      `make kopf-check-gegenprobe` und `make abdeckung-gegenprobe` grün.
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
| `spec/spezifikation.md` | update | Abschnitt für Harness-Werkzeuge (Ort, Kennungen und Bezug nach §6) mit den Verträgen von `kopf-check` und `abdeckung`; Historie-Zeile in §11 |
| `harness/sensors/kopf-check.md`, `harness/sensors/abdeckung-check.md` | neu | per `cp` aus `.harness/baseline/v6.13.0/templates/harness/sensors/gate.template.md`; Lesart eines Laufs, Vertrag als Link auf den Abschnitt |
| `harness/README.md` | update | §Sensors: Target-Zellen von `make kopf-check` und `make abdeckung-check` verlinken die Sensor-Dateien; die Vertragszeile bleibt ein Satz |
| `tools/harness/kopf-check.sh`, `tools/test/abdeckung.sh` | update (nur Kommentare) | ZUSAGE- bzw. Beschreibungskopf → Verweis auf den Abschnitt plus Zuordnung Vertragspunkt → Fall der Gegenprobe |
| `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung-gegenprobe.sh` | geprüft | Fälle unverändert; ihr Kopf nennt den Vertrag, falls er heute die ADR nennt (Entscheidung beim Implementer, keine neue Zusage) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-lastenheft-pruefbarkeit` liegt in `done/`
(WIP-Limit 1). Reihenfolge der wellenlosen Reihe: `slice-lastenheft-pruefbarkeit`,
dieser Slice, `slice-harness-lint`, die vier Umstellungs-Slices,
`slice-harness-abdeckung-gate`, `slice-harness-coverage`, `slice-harness-mutation`.
Erster Schritt nach dem Start, vor jedem Commit an Spezifikation oder Skript: Der
Architect entscheidet die Randformen aus §6 (Ort, Kennungen, Schärft-Bezug,
Referenz-Richtung, Umfang); wellenlos heißt nicht ohne Architect
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Beim Übertragen zeigt sich,
  dass Skript und ADR-Wortlaut an mehr als zwei Stellen auseinanderliegen, oder der
  Architect entscheidet, den Abdeckungs-Vertrag samt
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) schon hier zu
  schreiben; dann trägt dieser Slice `kopf-check` und den Abschnitt, `abdeckung` geht
  an einen eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): Der Architect findet keinen
  Schärft-Bezug, der ohne Änderung einer angenommenen ADR trägt, und eine ersetzende
  ADR je Gate ist nötig; oder die Klasse `spec-straten` des Doku-Gates lässt den
  Abschnitt in der gewählten Form nicht zu und eine Änderung an `.d-check.yml`
  braucht eine eigene Entscheidung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Closure-Notiz mit Lerneintrag; dazu der Ausgang von
`BEO-REPO/harness-lesart-ohne-entscheidungsort` im Register.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — alle **offen**, soweit nicht als
entschieden markiert. Entschieden werden sie nach dem Start, vor dem ersten Commit an
Spezifikation oder Skript, vom Architect; eine Entscheidung des Nutzers wird
festgehalten. Wo eine Entscheidung steht, sagt die Randform selbst: Die zum Abschnitt
gehören dorthin, die zum Vorgehen dieses Slice hierher. Was hier nicht steht,
entscheidet der Implementer nicht, er gibt es zurück.

- **Ort: Abschnitt oder eigene Datei** — Kurs-Antwort: ein eigener Abschnitt in
  `spec/spezifikation.md` oder eine ihr zugeordnete Datei. Für den Abschnitt spricht:
  Die Klasse `spec-straten` in `.d-check.yml` nennt genau drei Dateien; eine eigene
  Datei fiele in die Klasse `aussen`, und `spec/spezifikation.md` dürfte sie nicht
  verlinken, bis `.d-check.yml` und die Source Precedence in `harness/README.md`
  nachgezogen sind. Offen: welcher Abschnitt — ein neuer zwischen §9
  *Testanforderungen* und §10, ein neuer vor §11 *Historie*, oder ein Unterabschnitt
  von §9 neben `SPEC-038` (die Prüfungen des Quellstands sind Testanforderungen im
  weiteren Sinn); bei einem neuen Abschnitt mit Nummer verschiebt sich die Nummer der
  folgenden, und Verweise auf `spezifikation.md §10`/`§11` sind nachzuziehen
  (Empfehlung des Planners: neuer §10 *Harness-Werkzeuge* vor *Nicht zugesichert in
  v1*, Verweise per `grep` nachziehen).
- **Kennungen** — Kurs-Antwort: eigene `SPEC`-Kennungen. Offen: fortlaufend im
  Zählraum der Datei (die höchste vergebene ist heute `SPEC-046`), eine Kennung je
  Werkzeug oder je Vertragsteil (etwa Gegenstand, Kennungen, Kopf, Ausgabe bei
  `kopf-check`). Bei eigener Datei: eigener Zählraum oder fortlaufend über beide
  (Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe, „fortlaufend je
  Datei“; zwei `SPEC-001` wären mehrdeutig). Die Kennungen vergibt dieser Slice;
  Folge-Slices nennen sie erst, wenn sie stehen (`make kopf-check`).
- **Bezug nach oben** — „präzisieren ja, erweitern nie“ (Kopf der Spezifikation): Der
  Abschnitt braucht eine Lastenheft-Anforderung, die er präzisiert. Kandidaten:
  `LH-QA-04` (Automatisierbarkeit; Bezug von [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)) für beide,
  `LH-QA-07` für die Abdeckung und die Gates der Prüfbarkeit. Offen, ob jede Kennung
  ihren Bezug nennt oder der Abschnitt einmal.
- **Schärft-Bezug der angenommenen ADRs** — [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) tragen
  `Schärft: —` und sind `Accepted`; der Schärft-Eintrag ist Inhalt (`AGENTS.md` §3.5).
  Kandidaten: (a) die ADRs bleiben unverändert, die Verbindung trägt die Sensor-Datei
  (Bindung: ADR, Vertrag: Link auf den Abschnitt) und der Kopf dieses Slice; die
  Spezifikation nennt keine ADR (Klasse `spec-straten`, Regel gegen `adr`). (b) Je ADR
  eine ersetzende ADR mit `Supersedes` und gesetztem `Schärft:` — drei ADRs nur für ein
  Feld, die Entscheidung selbst bliebe gleich. (c) Eine Rückrichtung über die
  Spezifikation, die die Verbindung ohne ADR-Änderung trägt (Vorschlag aus der
  Auftragslage; Form offen, sie darf die Regel `spec-straten → adr` nicht verletzen).
  Neue Gate-ADRs (Lint, Coverage, Mutation) setzen `Schärft:` von Anfang an.
- **Lesart-Rangfolge bei Doppelung** — [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) führt die neun Punkte als Regeldetail
  (vor der Kurs-Antwort so entschieden), der Abschnitt führt sie künftig ebenfalls. Die
  Spezifikation steht im Rang über der ADR (`harness/README.md` §Source precedence):
  Weicht der Abschnitt ab, gilt er. Offen: ob der Abschnitt den Wortlaut der ADR
  übernimmt oder in eigener Fassung schreibt; eine eigene Fassung darf inhaltlich
  nichts ändern (§1, *Verhaltensänderung*). Dieselbe Doppelung trägt
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Accepted, bleibt
  unverändert): Ihr Satz „Form der Deklaration, Pfad-Schreibweise und Fehlformen stehen
  im Kopf des Abdeckungs-Skripts“ nennt einen Ort, den die Kurs-Antwort ersetzt. Der
  Abschnitt übersteuert ihn: Form, Pfad-Schreibweise und Fehlformen stehen dort, weil
  die Spezifikation im Rang über der ADR steht, und der Skriptkopf trägt nur, womit das
  Werkzeug geprüft ist (Verifikation V-65 von `slice-lastenheft-pruefbarkeit`).
- **Was der Abschnitt nennen darf** — `spec-straten → aussen` ist verboten: Der
  Abschnitt verlinkt weder Skript noch Sensor-Datei noch Makefile. Offen: ob er
  Werkzeug und Pfade als Text nennt (etwa `make kopf-check`,
  `docs/plan/planning/open/`) — der geprüfte Gegenstand von `kopf-check` sind
  Repo-Pfade, ohne sie ist der Vertrag nicht sagbar.
- **Umfang des Abdeckungs-Vertrags** — der Vertrag, den `tools/test/abdeckung.sh`
  heute einhält ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)),
  oder schon die Zusagen aus
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), die
  `slice-harness-abdeckung-gate` erst baut. Empfehlung des Planners: nur der heutige
  Vertrag; die Erweiterung schreibt `slice-harness-abdeckung-gate` mit ihren Randformen
  in den Abschnitt, denn eine Zusage ohne Gegenprobe-Fall verletzte `AGENTS.md` §3.11.
- **Welche Targets eine Sensor-Datei bekommen** — nach der Vorlage nur, wo ein Satz
  nicht reicht. Kandidaten: `make kopf-check` (Befund-Zeile, Exit 0/1/2, Grenze) und
  `make abdeckung-check` (`--check`, Exit 0/1). Offen: ob die Gegenproben
  (`make kopf-check-gegenprobe`, `make abdeckung-gegenprobe`) und das Werkzeug
  `make abdeckung` (schreibt, prüft nicht) eine eigene Datei bekommen, in der Datei des
  Gates mitstehen oder bei ihrer Zeile bleiben.
- **Was im Skriptkopf bleibt** — Kurs-Antwort: womit das Werkzeug geprüft ist. Offen:
  Form der Zuordnung (je Vertragspunkt der Name des Gegenprobe-Falls, wie heute
  `nr2-punkt-endet-wort`), und ob die GRENZE-Zeile von `kopf-check.sh` als Grenze des
  Vertrags in den Abschnitt wandert oder als Grenze der Prüfung im Kopf bleibt.

**Risiken:**

- **Übertragung verliert oder verschiebt eine Zusage** — der Vertrag steht heute an
  drei Stellen (ADR, archivierter §6, Skriptkopf), die nicht wortgleich sind; beim
  Zusammenführen kann eine Lesart wegfallen oder schärfer werden. Gegenmittel: Abgleich
  Punkt für Punkt in der DoD, Abweichung als Randform an den Architect. — **Ausgang:** —
  (bei Closure)
- **Folge-Slices planen gegen einen Abschnitt, den es noch nicht gibt** —
  `slice-harness-lint`, `slice-harness-abdeckung-gate`, `slice-harness-coverage`,
  `slice-harness-mutation` und die vier Umstellungs-Slices nennen den Abschnitt als
  Ort ihrer Randformen; verschiebt dieser Slice Ort oder Form, sind ihre §6 nachzuziehen
  (`AGENTS.md` §3.9). — **Ausgang:** — (bei Closure)
- **Schärft-Bezug bleibt Prosa** — wählt der Architect (a), führt keine ADR einen
  maschinenlesbaren Bezug auf den Abschnitt; die Verbindung hängt an der Sensor-Datei.
  — **Ausgang:** — (bei Closure)
- **Nummernverschiebung** — ein neuer nummerierter Abschnitt vor *Historie* verschiebt
  §-Verweise in Plänen und Archiv; eingefrorene Zeitdokumente werden nicht
  nachgezogen. — **Ausgang:** — (bei Closure)

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
- **Drei Paarungen:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `3dfef35` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×; das zweite Auftreten, Review
  F-406, belegt die Closure von `slice-lastenheft-pruefbarkeit`) — dieser Slice ist die
  Antwort: Er legt den Ort an. Sein Ausgang im Register kommt mit der Closure.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×; V-40, V-43 aus
  `slice-harness-kopf-sensor`) — trifft den Vertrag von `kopf-check`; die zwei offenen
  Randformen werden übertragen, nicht entschieden (§1).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (11×, verkörpert in `AGENTS.md`
  §3.11) — jede Zusage des Abschnitts braucht ihren Gegenprobe-Fall (DoD, erster
  Liefer-Punkt).
- `BEO-REPO/plan-folgt-korrektur-nicht` (10×, verkörpert in §3.9, Sensor
  `make kopf-check`) — die Folge-Slices nennen den Abschnitt; Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (8×, verkörpert in §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) und
  `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — darum stehen die offenen
  Randformen in §6 und §4 *Start* nennt den Architect-Schritt; ein drittes Auftreten
  hier brächte `randform-im-code-entschieden-dann-zurueckgegeben` auf die Schwelle.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
