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

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode (4): die Prüfumgebung des Projekts, deren Werkzeuge der Abschnitt festlegt; §6 *Bezug nach oben*). Bindung an Entscheidungen: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) (Kopf-Sensor), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckung je Anforderung und Pfad), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Gate-Nachweise und geteilte Messung).

**Berührte Spec-Stellen:** `spezifikation.md §11` (neu: *Harness-Werkzeuge*, `SPEC-047` für `kopf-check`, `SPEC-048` für `abdeckung`) · `spezifikation.md §12` (*Historie*, bisher §11)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Verträge der zwei Harness-Werkzeuge, die vor diesem Slice ohne
Spezifikationsstelle lebten, stehen im neuen §11 *Harness-Werkzeuge* der Spezifikation (Technik-Stratum,
fortschreibbar), eine Kennung je Werkzeug: `kopf-check` als `SPEC-047` mit den neun
Punkten von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und den sieben
Lesarten, die der Architect beim Liefern bestätigt hat (vor diesem Slice im ZUSAGE-Kopf von
`tools/harness/kopf-check.sh`; Herkunft: §6 des archivierten
`slice-harness-kopf-sensor` in `docs/plan/planning/done/welle-replay-semantik/archiv.zip`),
und `abdeckung` als `SPEC-048` mit [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und
dem damaligen Kopf von `tools/test/abdeckung.sh`; wo das Skript genauer ist, in
dessen Fassung (§6 *Lesarten beim Übertragen*). Wie ein Lauf zu lesen ist (Vertrag, Grenze,
Ausgabe, Ausgänge, Sperren), steht je Gate in einer Sensor-Datei
`harness/sensors/<target>.md` nach der Baseline-Vorlage, die den Vertrag verlinkt und
nichts neu entscheidet; `harness/README.md` §Sensors zeigt darauf. Die Skriptköpfe
behalten nur, womit das Werkzeug selbst geprüft ist; fehlt einem Vertragspunkt der
Fall, ergänzt ihn die Gegenprobe (§6). Danach entscheiden die
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
  baut sie und schreibt ihre Randformen in `SPEC-048` fort, ohne neue Kennung
  (entschieden in §6 *Umfang des Abdeckungs-Vertrags* und *Kennungen*).
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
  die Kommentarköpfe zweier Skripte unter `tools/` und, wo einem Vertragspunkt der
  Fall fehlt, deren Gegenproben (§6 *Vertragspunkt ohne Gegenprobe-Fall*); kein
  ausführbarer Code der Werkzeuge selbst.
- Die angenommenen ADRs — [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md),
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) bleiben unverändert
  (`AGENTS.md` §3.5); ihr `Schärft: —` bleibt stehen, die Verbindung zum Abschnitt
  trägt die Sensor-Datei (§6 *Schärft-Bezug*, Variante (a)).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Spezifikation: Der Abschnitt für Harness-Werkzeuge liegt als §11 nach *Nicht
      zugesichert in v1* und führt `SPEC-047` für `kopf-check` mit den neun Punkten von
      [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und den Lesarten
      (1) bis (5), (a) und (b), `SPEC-048` für `abdeckung` mit dem Vertrag aus
      [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und dem heutigen
      Skriptkopf (Deklarationsform, Pfade, Nachweisart nach Ort, erzeugte Tabellen,
      `--check`, Ausgänge). Jede Zusage des Abschnitts ist eine Zusage, die das Skript
      heute einhält: Zu jeder nennt der Skriptkopf den Fall der Gegenprobe
      (`AGENTS.md` §3.11); fehlt der Fall, ergänzt ihn die Gegenprobe (§6). §12
      *Historie* trägt eine Zeile.
- [ ] Sensor-Dateien `harness/sensors/kopf-check.md` und
      `harness/sensors/abdeckung-check.md` per `cp` aus
      der vendored Vorlage, je mit Vertrag (Link auf den Abschnitt), Grenze, Ausgabe
      und Ausgängen, Sperren und Bindung; sie entscheiden nichts, was der Abschnitt
      nicht sagt. Die Target-Zellen in `harness/README.md` §Sensors verlinken sie.
- [ ] Skriptköpfe von `tools/harness/kopf-check.sh` und `tools/test/abdeckung.sh`
      tragen statt der Vertragspunkte den Verweis auf den Abschnitt und je Punkt nur
      noch, womit das Werkzeug geprüft ist (Fall der Gegenprobe); kein Vertragspunkt
      geht dabei verloren (Abgleich Punkt für Punkt im Bericht). Ausführbarer Code
      unverändert: `git diff` der beiden Skripte berührt nur Kommentarzeilen,
      `make kopf-check-gegenprobe` und `make abdeckung-gegenprobe` grün; jeder
      ergänzte Fall ist gegen eine Mutation des Skripts rot gesehen (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` | update | neuer §11 *Harness-Werkzeuge* mit `SPEC-047` (`kopf-check`) und `SPEC-048` (`abdeckung`), Bezug und Form nach §6; *Historie* wird §12 und bekommt eine Zeile |
| `harness/sensors/kopf-check.md`, `harness/sensors/abdeckung-check.md` | neu | per `cp` aus `.harness/baseline/v6.13.0/templates/harness/sensors/gate.template.md`; Lesart eines Laufs, Vertrag als Link auf den Abschnitt |
| `harness/README.md` | update | §Sensors: Target-Zellen von `make kopf-check` und `make abdeckung-check` verlinken die Sensor-Dateien; die Vertragszelle ist ein Satz, die von `make kopf-check` schrumpft darauf |
| `tools/harness/kopf-check.sh`, `tools/test/abdeckung.sh` | update (nur Kommentare) | ZUSAGE- bzw. Beschreibungskopf → Verweis auf den Abschnitt plus Zuordnung Vertragspunkt → Fall der Gegenprobe |
| `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung-gegenprobe.sh` | update | bestehende Fälle behalten ihre Eingaben; ein Fall kommt hinzu, wo einem Vertragspunkt keiner gilt (§6 *Vertragspunkt ohne Gegenprobe-Fall*, Stand dort), je gegen eine Mutation rot gesehen; bei `abdeckung` ist die Erwartung abgelehnter Fälle geschärft (Exit genau 1, Datei und Zeile auf stderr) und die Prüfungen im gültigen Baum tragen Fall-Namen; ihr Kopf nennt den Vertrag statt der ADR |
| `docs/plan/planning/open/slice-harness-abdeckung-gate.md` | update | Kopf `Berührte Spec-Stellen` und §6 *Ort*: schreibt den Vertrag von `abdeckung` in `SPEC-048` fort und vergibt keine neue Kennung (§6 *Kennungen*; `AGENTS.md` §3.9) |
| `docs/plan/planning/open/slice-harness-coverage.md`, `docs/plan/planning/open/slice-harness-mutation.md` | update | je ein Satz mit der Entscheidung des Nutzers vom 2026-10-06 zum Ort der Schwelle (§1 bzw. §2, Liefer-Punkte unverändert; §6 *Schwellen künftiger Gates*) |

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

DoD vollständig, `make gates` grün, Closure-Notiz mit Lerneintrag. Der Register-Ausgang
von `BEO-REPO/harness-lesart-ohne-entscheidungsort` steht schon (`gestrichen` mit der
Closure von `slice-lastenheft-pruefbarkeit`); dieser Slice schuldet für ihn keinen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — **entschieden vom Architect am
2026-10-06**, nach dem Start und vor dem ersten Commit an Spezifikation, Skript oder
Gegenprobe. Die Entscheidungen zum Inhalt des Abschnitts überträgt der Implementer in
die Spezifikation, dort gelten sie (Rang über ADR und Plan); die zum Vorgehen dieses
Slice gelten hier. Was hier nicht steht, entscheidet der Implementer nicht, er gibt es
zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe).

- **Ort** — neuer Abschnitt `## 11. Harness-Werkzeuge` in `spec/spezifikation.md`,
  **nach** §10 *Nicht zugesichert in v1*; *Historie* wird §12. Keine zugeordnete
  Datei: Die Klasse `spec-straten` in `.d-check.yml` nennt genau drei Dateien, eine
  vierte fiele in `aussen` und verlangte eine Änderung an `.d-check.yml` (Pfade und
  `order`) und an der Source Precedence in `harness/README.md`; der Abschnitt
  braucht keine von beiden. Abweichung von der Empfehlung des Planners (neuer §10):
  Hinter §10 verschiebt sich nur die Nummer der Historie, `SPEC-039` behält §10, und
  der Produkt-Teil endet mit seiner Abgrenzung, bevor die Werkzeuge der
  Prüfumgebung kommen. Kein lebendes Dokument verweist heute auf `spezifikation.md`
  §11 oder den Anker der Historie (geprüft per `grep`); die Folge-Slices nennen den
  Abschnitt beim Namen *Abschnitt für Harness-Werkzeuge*, der Titel hält ihn.
- **Kennungen** — fortlaufend im Zählraum der Datei, **eine Kennung je Werkzeug**:
  `SPEC-047` für `kopf-check`, `SPEC-048` für `abdeckung`. Innerhalb einer Kennung
  sind die Vertragspunkte nummeriert, `(1)`, `(2)`, …; die Nummer ist die Adresse,
  auf die der Skriptkopf zeigt. Für `SPEC-047` sind es die Nummern 1 bis 9 von
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), damit die Fall-Präfixe
  `nr<n>-` der Gegenprobe gültig bleiben; die sieben bestätigten Lesarten werden in
  den Punkt eingearbeitet, den sie präzisieren, und bekommen keine eigene Nummer.
  Kennung je Vertragsteil ist verworfen: Sie bläht jeden Kopf eines Folge-Slice, der
  eine Randform eines Werkzeugs fortschreibt, und adressiert nichts, was die
  Punkt-Nummer nicht schon adressiert. Eine Randform eines bestehenden Werkzeugs
  schreibt der Folge-Slice in dessen Kennung fort (`slice-harness-abdeckung-gate` in
  `SPEC-048`, ohne neue Kennung); ein neues Werkzeug bekommt die nächste freie
  Kennung, vergeben vom Slice, der es liefert.
- **Bezug nach oben** — einmal, im Einleitungsabsatz des Abschnitts, nicht je
  Kennung: Der Abschnitt legt die Werkzeuge der *Prüfumgebung des Projekts* fest,
  die [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
  Messmethode (4) nennt, und sagt dazu ausdrücklich: Keine Kennung des Abschnitts ist
  eine Zusage des Produkts, und keine erweitert eine Anforderung des Lastenhefts.
  `LH-QA-04` ist **nicht** der Bezug: Es verlangt, dass Start, Betrieb und Beendigung
  *des Produkts* skriptgesteuert gehen; die Bezüge von [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) darauf
  bleiben als eingefrorene Lesart stehen, die Spezifikation übernimmt sie nicht. Ein
  Bezug je Kennung ist verworfen: `kopf-check` präzisiert keine Anforderung einzeln,
  der Bezug wäre erfunden. Akzeptiertes Negativ: Die Präzisierung ist eng — sie trägt
  nur, weil das Lastenheft die Prüfumgebung nennt, nicht ihre Werkzeuge.
- **Schärft-Bezug der angenommenen ADRs** — **(a)**: [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)
  bleiben unverändert, `Schärft: —` bleibt stehen (`AGENTS.md` §3.5). Die Verbindung
  ADR ↔ Vertrag trägt die Sensor-Datei (Vertrag: Link auf die Kennung, Bindung: die
  ADR) und die Index-Zeile in `harness/README.md` §Sensors. Gegen (b): drei
  ersetzende ADRs für ein Feld bei gleicher Entscheidung; dazu setzte jede Ablösung
  `superseded`, und `status: forbidden` in `.d-check.yml` färbte jedes lebende
  Dokument rot, das die alte ADR verlinkt (`AGENTS.md` §3.9, `harness/README.md`
  §Sensors, offene Slices) — Aufwand ohne neue Entscheidung. Gegen (c): Eine
  Rückrichtung aus der Spezifikation auf eine ADR hat keine zulässige Form — die
  Regel `{from: spec-straten, to: adr, allow: false}` gilt auch in der Historie, und
  das Baseline-Regelwerk (`grundlagen-referenz-richtung.md` §Spec-Straten) kennt
  Spec → ADR auch nicht als Quellen-Spalte; die einzige Rückrichtung ohne
  Regelverstoß liegt außerhalb der Straten, und das ist (a). Akzeptiertes Negativ:
  Für diese drei ADRs fehlt die maschinenlesbare Änderungskopplung ADR → Spec; wer
  eine von ihnen ablöst, setzt in der Nachfolgerin `Schärft:` auf die Kennung. Neue
  Gate-ADRs (Lint, Coverage, Mutation, Gate-Nachweise) setzen `Schärft:` von Anfang
  an.
- **Lesart-Rangfolge bei Doppelung** — Der Abschnitt schreibt in **eigener Fassung**,
  nicht im Wortlaut der ADR, und sagt nur, was das Skript heute tut; wo das Skript
  genauer ist als [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) (Feldmarke nur am Absatzanfang, Leerzeile mit Leerzeichen
  und Tabs, nicht lesbarer Plan als ein Befund), steht die genauere Fassung — das
  präzisiert die ADR, es ändert sie nicht. Ein echter Widerspruch zwischen ADR und
  Skript ist eine Randform-Rückgabe, kein Text-Entscheid. Bei Doppelung gilt der
  Abschnitt (Rang über der ADR). Der Satz in [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), Form der Deklaration,
  Pfad-Schreibweise und Fehlformen stünden im Kopf des Abdeckungs-Skripts, ist
  **übersteuert**, nicht abgelöst: Er nennt einen Ort, ein Regeldetail im Sinn von
  `AGENTS.md` §3.8, keine Entscheidung; die Entscheidung der ADR (Nachweisart Gate,
  geteilte Messung) bleibt unberührt, also braucht es keine Folge-ADR. Wer der ADR
  zum Skriptkopf folgt, findet dort als erste Zeile den Zeiger auf `SPEC-048`; das
  trägt den Übergang. Dasselbe gilt für den Satz im Kontext von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), den
  Vertrag nenne `harness/README.md` §Sensors.
- **Was der Abschnitt nennen darf** — Targets (`make kopf-check`,
  `make abdeckung-check`, `make abdeckung`) und geprüfte Repo-Pfade
  (`docs/plan/planning/open/`, `docs/user/abdeckung-*.md`, `test/integration/`) als
  Text in Code-Spans, **ohne Link**; keinen Skriptpfad, keine Sensor-Datei, kein
  Makefile, auch nicht als Text — die Schnittstelle ist das Target. Verboten bleiben
  Links nach `aussen` und jedes Token, das die Matrix fängt: die Zeichenfolge aus
  `slice` und Bindestrich in Kleinschrift (den Gegenstand von `kopf-check` schreibt
  der Abschnitt wie [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md): *Dateiname mit Präfix `slice` und Bindestrich*), die aus
  `welle` und Bindestrich, eine `MR`-Kennung und eine ADR-Kennung, auch in der
  Historie-Zeile. Links auf das Lastenheft sind erlaubt (aufwärts).
- **Umfang des Abdeckungs-Vertrags** — nur der heutige Vertrag: Deklarationsform,
  Pfade, Nachweisart nach Ort, die vier Tabellen, `--check`, Ausgänge, dazu was die
  Gegenprobe heute schon prüft (Maskierung von `|`, Dateirechte 0644). Die
  Erweiterung nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) schreibt `slice-harness-abdeckung-gate` mit ihren
  Randformen in `SPEC-048` fort; eine Zusage ohne Fall verletzte `AGENTS.md` §3.11.
- **Vertragspunkt ohne Gegenprobe-Fall** — Ergibt der Abgleich, dass ein Punkt, den
  das Skript heute einhält, keinen Fall hat (bei `abdeckung` absehbar: Nachweisart
  E2E nach Ort, Inhalt von `abdeckung-e2e.md` und `abdeckung-gesamt.md`, eine
  Anforderung mit fehlendem Pfad fehlt in `abdeckung-vollstaendig.md`, `LH-QA`/`LH-RB`
  nur mit `Messung`, Exit genau 1), **ergänzt der Implementer den Fall in der
  Gegenprobe** und sieht ihn gegen eine Mutation des Skripts rot (`AGENTS.md` §3.10);
  der ausführbare Code von `kopf-check.sh` und `abdeckung.sh` bleibt unverändert.
  Ein Punkt als „Grenze“ im Abschnitt ist dafür kein Ausweg: Die Grenze sagt, was
  das Gate nicht prüft, nicht, was am Werkzeug ungeprüft ist. Ist ein Fall nicht ohne
  Eingriff ins Werkzeug herstellbar, ist das eine Randform-Rückgabe.
  **Stand** (Implementer, 2026-10-06): ergänzt, je gegen eine Mutation in einer Kopie
  des Skripts rot gesehen, die alte Gegenprobe bleibt gegen dieselbe Mutation grün —
  bei `kopf-check` `nr3-bereich-zeilenumbruch-abschnitt` (Bereich über einen
  Zeilenumbruch in §1) und die Abbruchzeile in `nr8-ohne-ablage`; bei `abdeckung`
  `doppelt`, `dateiende`, `qa-happy`, `fa-messung`, `ueberschrift-ebene`,
  `ohne-deklaration-unit`, `nicht-gelesen`, `inhalt-e2e`, `inhalt-unit`,
  `inhalt-gesamt`, `inhalt-vollstaendig`, dazu Exit genau 1 mit Datei und Zeile auf
  stderr in jedem abgelehnten Fall und in `check-veraltet`. Kein Fall brauchte einen
  Eingriff ins Werkzeug.
- **Lesarten beim Übertragen** (Implementer, 2026-10-06; zur Bestätigung beim
  Architect, nicht hier entschieden) — der Abschnitt schreibt sie nach
  *Lesart-Rangfolge* in der genaueren Fassung des Skripts, je mit Fall:
  `TestE2E` als Name, nicht als Ort, „jeder E2E-Test“ in
  [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) heißt also jeder
  Test mit diesem Namen (`ohne-deklaration-unit`; `TestMain` unter
  `test/integration/` trägt keine); gelesen werden die Testdateien außer unter `.git/`
  und `.harness/` (`nicht-gelesen`); die Anforderung steht als Überschrift dritter
  Ebene im Lastenheft (`ueberschrift-ebene`); die Ausgänge 0, 1 und 2 sind die der
  Prüfung, `make` meldet jeden Ausgang ungleich 0 als Fehler (`nr9-scharf`, die
  Sensor-Datei sagt, woran die Ursachen auseinanderzuhalten sind). Nicht übertragen
  und weiter offen: Symlink als Plan und nicht lesbares Lifecycle-Verzeichnis
  (`BEO-REPO/gate-uebergeht-ablage-eintrag-still`); der Abschnitt sagt „Dateien“ wie
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), die Sensor-Datei nennt
  beides als offen.
- **Welche Targets eine Sensor-Datei bekommen** — `harness/sensors/kopf-check.md` und
  `harness/sensors/abdeckung-check.md`, je per `cp` aus der vendored Vorlage. Die
  Gegenproben bekommen keine: Was sie prüfen, ist die Frage, ob das Werkzeug richtig
  ist, und lebt beim Werkzeug; ihr Vertrag bleibt ein Satz in ihrer Zeile. `make
  abdeckung` bekommt keine: Es ist ein Werkzeug, seine Wirkung ist Teil von
  `SPEC-048`, und `abdeckung-check.md` nennt es als Weg aus dem Rot (*Ausgabe und
  Ausgänge*). Die Index-Zelle von `make kopf-check` schrumpft auf einen Satz, die
  Einzelheiten stehen in Abschnitt und Sensor-Datei. Eine *Sperre* nennt die Datei
  nur, wo das Skript eine benannte Abbruch-Meldung hat und ein Fall sie hält
  (`kopf-check`: fehlende Ablage, Exit 2); `abdeckung` hat keine und sagt das.
- **Was im Skriptkopf bleibt** — der Zweck in einem Satz, der Zeiger
  `Vertrag: spec/spezifikation.md SPEC-047` (bzw. `SPEC-048`) samt Sensor-Datei als
  Text, die Aufruf-Zeile, und die Zuordnung `GEPRÜFT DURCH <Gegenprobe>:` mit einer
  Zeile je Vertragspunkt, `(<n>) <Stichwort> — <Fall-Namen>`. Fall-Namen wie in der
  Gegenprobe; `nr<n>-*` steht für alle Fälle mit diesem Präfix und nur dort, wo sie
  alle zu Punkt `n` gehören; ein Fall ohne Namen (etwa Nr. 9, der `make`-Lauf im
  Temp-Baum) wird mit seinem Kommentar-Titel in der Gegenprobe genannt. Ausgabe- und
  Ausgang-Zeilen wandern in den Abschnitt. Die **GRENZE-Zeile** von `kopf-check.sh`
  wandert als *Grenze* in `SPEC-047` (was kein Befund ist, ist Vertrag), die
  Sensor-Datei liest sie unter *Grenze*, und der Skriptkopf ordnet ihr die grünen
  Fälle `neg-*` zu; eine Grenze der Prüfung selbst (etwa der übersprungene Fall
  unter root) bleibt im Kopf der Gegenprobe.
- **Schwellen künftiger Gates — Empfehlung für die Folge-Slices** (vom Nutzer am
  2026-10-06 angenommen: Schwellen künftiger Gates stehen in der Spezifikation; im
  Einzelnen entschieden wird sie in deren Architect-Schritt): Die **Zahl**
  (Coverage-Schwelle, Mutations-Schwelle) steht in der Kennung des Werkzeugs in
  diesem Abschnitt, als Konstante mit Einheit und einem Satz Begründung; die **ADR**
  entscheidet Mechanismus, Messart und die Regel, nach der die Zahl aus der Messung
  folgt, hält die **Messung** mit Quellstand im Kontext fest und setzt `Schärft:` auf
  die Kennung. Gründe: LH-QA-07 (2) verlangt eine *festgelegte* Schwelle, und das
  Festlegen ist Präzisieren, also Technik-Stratum; eine Zahl in der Accepted-ADR ist
  Regeldetail (`AGENTS.md` §3.8) und zwänge bei jeder Anhebung zu einer
  Ersetzungs-ADR. `AGENTS.md` §3.6 bleibt unberührt: Eine **Senkung** braucht
  weiterhin eine ADR, die die Kennung schärft; eine Anhebung ist Fortschreibung. Die
  Messung ist ein Zeitdokument und gehört in die ADR, nicht in den Abschnitt.
  Betroffen: §1 *Ziel* von `slice-harness-coverage` und §2 von
  `slice-harness-mutation`; hier ist je ein Satz mit der Entscheidung nachgezogen, die
  Liefer-Punkte bleiben unverändert und werden in deren Architect-Schritt
  angeglichen (DoD von `slice-harness-coverage` nennt noch „Schwelle aus der ADR des
  Gates“).

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
  Gewählt ist (a); akzeptiertes Negativ mit Grund in §6 *Schärft-Bezug*. — **Ausgang:** — (bei Closure)
- **Nummernverschiebung** — ein neuer nummerierter Abschnitt vor *Historie* verschiebt
  §-Verweise in Plänen und Archiv; eingefrorene Zeitdokumente werden nicht
  nachgezogen. Gewählt ist §11 nach *Nicht zugesichert*; es verschiebt sich nur die
  Historie (§6 *Ort*). — **Ausgang:** — (bei Closure)

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

- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (2×, `gestrichen` mit der Closure
  von `slice-lastenheft-pruefbarkeit`, nachgesehen am Stand `e491e0c`: die Ursache ist
  mit der Antwort des Kurs-Repos weggefallen) — dieser Slice legt den Ort an, der die
  Streichung trägt; ein Ausgang steht nicht mehr aus. Eine Randform, die trotzdem in
  Skriptkopf oder Plan entschieden wird, zählt bei den beiden Randform-Einträgen
  unten.
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
