# Slice slice-harness-mutation: Mutationstests als Gate

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
nächsten Welle-Closure. Angelegt bei der Closure von welle-replay-semantik nach
Entscheidung des Nutzers vom 2026-10-06; Reihenfolge in §4 *Start*.

**Bezug:** — (Harness-Arbeit). Bindung an Entscheidungen:
[ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Tests in den
Stufen des Multistage-`Dockerfile`, netzlos außer der Download-Stufe). Die ADR des
neuen Gates schreibt der Architect im Slice, vor dem Code; ihre Nummer vergibt der
ADR-Index.

**Berührte Spec-Stellen:** `spezifikation.md §11` (*Harness-Werkzeuge*: neue Kennung des Gates mit Schwelle und Randformen, vergeben von diesem Slice)

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

**Ziel:** Ein Werkzeug für Mutationstests (Kandidaten: gremlins, go-mutesting) läuft
Docker-only, per Digest gepinnt, als Gate in `make gates`: Es mutiert den Code und
meldet rot, wenn zu viele Mutanten die Tests überleben — über eine Schwelle für den
Bestand oder für die geänderten Pakete. Damit fängt eine Maschine, was heute nur
Review und Verifikation finden: eine Zusage, die kein Test hält.

**Herkunft:** Lese-Schritt der Closure von welle-replay-semantik. Zwei Einträge des
Beobachtungs-Registers sind über der Schwelle und als Prosa ausgeschöpft:
`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (11×, `AGENTS.md` §3.11) und
`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (10×, `AGENTS.md` §3.10). Beide
traten in jedem Slice der Welle wieder auf, obwohl der Implementer die Mutationen
selbst fährt; die Funde lagen in Mutationen, die er nicht gewählt hatte
(V-52, V-56, V-60). Der Nutzer hat am 2026-10-06 einen Sensor entschieden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Neue Tests für Mutanten, die heute überleben — die Schwelle folgt der Messung; das
  Schließen der Lücken ist eigene Arbeit mit eigenem Slice, aus der Liste der
  Überlebenden heraus, nicht vorab.
- Die Anweisungs-Abdeckung als Gate — übernimmt `slice-harness-coverage`; sie sagt,
  welcher Code unter Tests läuft, dieses Gate, ob die Tests ihn prüfen. Keine der
  beiden ersetzt die andere.
- Die Mutationsberichte von Implementer und Verifikation (`AGENTS.md` §3.10,
  `.claude/commands/implement-slice.md` Schritt 19, Modul 11) — Bestand bleibt bewusst
  stehen: Eine handgewählte Mutation je Zusage prüft den Vertrag des Slice, das
  Werkzeug mutiert nach Operatoren. Ob das Gate einen Teil der Handarbeit ersetzt,
  entscheidet der Retirement-Check nach seiner Einführung, nicht dieser Slice.
- Mutationen in den Integrationstests (`make test-integration`) — es wäre ein anderer
  Vorgang: Je Mutant ein Lauf mit PostgreSQL in einem eigenen Docker-Netz kostet
  Minuten. Ob die Integrationstests als Prüfer zählen, ist Randform *Umfang* in §6;
  eine Einbeziehung mit eigenem Lauf ist ein eigener Slice.
- Produkt-Code, Lastenheft und die Spezifikation außerhalb ihres Abschnitts für
  Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert `Dockerfile` oder ein
  eigenes Werkzeug-Image, Harness und ihre Doku; in der Spezifikation nur die
  Randformen seines Gates in jenem Abschnitt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Die Schwelle steht nach Entscheidung des Nutzers vom 2026-10-06 im Abschnitt für
Harness-Werkzeuge der Spezifikation, in der Kennung des Gates, als Konstante mit
Einheit und Begründung; die ADR des Gates hält Mechanismus, Messung mit Quellstand und
die Regel, nach der die Zahl folgt, und setzt `Schärft:` auf die Kennung.

- [ ] Gate: Das Werkzeug läuft netzlos aus einem per Digest gepinnten Image oder einer
      gepinnten Stufe, als Ziel an `GATE_CHECKS`, nach den Randformen am Ort aus §6;
      es endet rot, wenn der Anteil getöteter Mutanten unter der Schwelle
      liegt. Die Messung des Bestands, aus der die Schwelle folgt, steht mit Quellstand
      in der ADR; der Bestand ist grün.
- [ ] Gegenprobe als eigenes Ziel an `GATE_CHECKS`: ein Stand mit einem Test weniger
      (ein überlebender Mutant mehr) wird abgelehnt, einer genau an der Schwelle
      angenommen; dazu je entschiedener Randform aus §6, die eine Mutation fangen
      kann, ein Fall (äquivalenter Mutant in der Ausnahmeliste, Paketauswahl,
      Zeitgrenze). Die Gegenprobe baut aus einer Kopie unter eigenem Pfad
      (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`). Je Zusage ist die Mutation
      gesehen (`AGENTS.md` §3.10).
- [ ] Doku: `harness/README.md` §Sensors führt Gate und Gegenprobe mit Bindung an die
      ADR und sagt nur zu, was das Gate prüft (Paketumfang, Operatoren, Schwelle);
      `AGENTS.md` §3.10 und §3.11 nennen den Sensor für ihre maschinell prüfbare
      Hälfte, so wie §3.9 `make kopf-check` nennt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; die Zeile `**Sensor:**` in
      `state.md` von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und
      `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` wechselt von `geplant` auf
      `verkörpert` mit Zielort und Herkunfts-Anker.
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
| `docs/plan/adr/<NNNN>-mutationstests-als-gate.md`, `docs/plan/adr/README.md` | neu / update | ADR des Gates (Architect, vor dem Code): Werkzeug, Messung mit Quellstand, Regel, nach der die Schwelle folgt, `Schärft:` auf die Kennung des Gates; Index-Zeile |
| `spec/spezifikation.md` | update | Abschnitt für Harness-Werkzeuge: Kennung des Gates mit der Schwelle als Konstante mit Einheit und Begründung und den Randformen aus §6 (Entscheidung des Nutzers vom 2026-10-06) |
| `Dockerfile` oder eigenes Werkzeug-Image | update / neu | Stufe oder Image mit dem gepinnten Werkzeug; netzlos außer der Download-Stufe |
| `harness/mk/<mutation>.mk` | neu | Ziel des Gates und der Gegenprobe an `GATE_CHECKS`, Schwelle als eine benannte Stelle |
| `tools/harness/<gegenprobe>.sh` | neu | Mutanten in einer Kopie unter eigenem Pfad: ein Test entfernt (rot), Stand genau an der Schwelle (grün), je Randform ein Fall |
| `harness/README.md`, `AGENTS.md` | update | §Sensors: Gate und Gegenprobe; §3.10 und §3.11: Zeile *Sensor* |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-coverage` liegt in `done/` oder ist
ausdrücklich zurückgestellt (WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers
vom 2026-10-06: nach der Reihe Lastenheft, Lint-Werkzeug, Black-Box-Umstellung,
Bereinigung, Lint-Gate, `slice-harness-abdeckung-gate`, Coverage.
Grund der Reihenfolge, keine technische Abhängigkeit: Lint-Bereinigung und
Black-Box-Umstellung ändern die Tests, gegen die Mutanten laufen; eine Messung davor
wäre veraltet, bevor das Gate greift. Erster Schritt nach dem Start, vor jedem
Code-Commit: Der Architect wählt das Werkzeug, misst den Bestand (Mutanten gesamt,
getötet, überlebt, Laufzeit, je Paket, mit Quellstand), entscheidet die Randformen
aus §6 und schreibt die ADR (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung ergibt so viele
  überlebende Mutanten, dass eine Schwelle mit Biss erst nach neuen Tests grün wäre.
  Dann entstehen die Tests in eigenen Slices, und dieser liefert danach das Gate auf
  dem dann gemessenen Stand — oder das Gate startet nur für geänderte Pakete
  (Randform *Umfang*), mit Hochschalt-Trigger in der ADR.
- `in-progress` → `open` (blockiert — Carveout?): Keines der Kandidaten-Werkzeuge
  läuft netzlos und reproduzierbar mit der Go-Version des `Dockerfile`, oder die
  Laufzeit über den Bestand übersteigt jede Grenze, die der Nutzer für `make gates`
  hinnimmt. Dann zuerst eine Entscheidung über Werkzeug oder Lauf außerhalb von
  `make gates`.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit dem Mutations-Gate und seiner Gegenprobe,
Closure-Notiz mit Lerneintrag, die beiden Register-Einträge mit Sensor-Ausgang
`verkörpert`.

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

- **Werkzeug** — gremlins, go-mutesting oder ein anderes; Kriterien: gepflegt, läuft
  mit der Go-Version des `Dockerfile`, netzlos nach dem Download, Ausgabe
  maschinenlesbar, Operatoren wählbar. Offen, ebenso der Weg ins Image (Modul in
  `deps`, eigenes Image per Digest, Binary per Prüfsumme).
- **Schwelle** — die Zahl erst nach der Messung; die Messung steht mit Quellstand in
  der ADR, die Zahl in der Kennung des Gates in der Spezifikation (Entscheidung des
  Nutzers vom 2026-10-06). Offen:
  Anteil getöteter Mutanten gesamt, je Paket oder beides; gleich dem gemessenen Wert
  oder mit Abstand darunter; Rundung; der Stand genau an der Schwelle grün oder rot.
- **Umfang** — der ganze Bestand oder nur geänderte Pakete (gegen welchen
  Vergleichsstand: Hauptzweig, letzter Tag, Arbeitsbaum). Ob die Integrationstests als
  Prüfer zählen (§1). Ob `cmd/`, Port-Pakete und Testhilfen ausgenommen sind, mit
  `Why:` an einer Stelle.
- **Laufzeit** — Zeitgrenze je Mutant (ein Mutant, der eine Schleife endlos macht),
  Parallelität, Gesamtdauer in `make gates`. Offen: läuft das Gate in jedem
  `make gates` oder als eigenes Ziel mit eigenem Trigger, und was das für den
  Gate-Nachweis heißt.
- **Äquivalente Mutanten** — ein Mutant, den kein Test töten kann, weil er das
  Verhalten nicht ändert (in der Verifikation von `slice-replay-semantik-meldungscodes`
  etwa VE2). Offen: Ausnahmeliste mit Begründung je Eintrag (Ort, Form, Prüfung, dass
  ein Eintrag nicht veraltet), Schwelle unter 100 % als Puffer, oder beides.
- **Build-Kontext** — BuildKit überträgt eine Datei gleicher Größe und mtime nicht neu
  (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`). Offen: Mutiert das Werkzeug im
  Container (dann wirkt der Build-Kontext nicht) oder im Arbeitsbaum (dann nie), und
  wie die Gegenprobe ihre Mutanten in den Lauf bringt.
- **Operatoren** — welche Klassen (Bedingungsgrenzen, Negation, arithmetisch,
  entfernte Anweisungen, Rückgabewerte). Offen, ebenso, ob verneinende Zusagen
  (`.claude/commands/implement-slice.md` Schritt 19, verneinende Zusage) einen
  Operator haben.
- **Ausgabe** — grün: Gesamtzahl, getötet, überlebt, Schwelle; rot: dazu je
  Überlebendem Datei, Zeile, Operator. Offen: Sortierung, Kürzung langer Listen, ob die
  Liste auch bei grün erscheint.

**Risiken:**

- **Die Schwelle misst den falschen Stand** — gemessen vor Lint-Bereinigung oder
  Black-Box-Umstellung wäre die Zahl beim Einschalten überholt; darum die Reihenfolge
  in §4. — **Ausgang:** — (bei Closure)
- **Laufzeit macht `make gates` unbrauchbar** — dauert der Lauf so lange, dass er vor
  dem Handoff übersprungen wird, ist das Gate ein stilles Rot. — **Ausgang:** — (bei
  Closure)
- **Gegenprobe sieht den Mutanten nicht** — BuildKit überträgt eine geänderte Datei
  nicht (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×). — **Ausgang:** — (bei
  Closure)
- **Konfiguration wirkt anders, als sie gelesen wird** — eine Ausnahme oder die
  Paketauswahl greift nicht, oder die Schwelle kommt nicht an und ein Default gilt
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 1×). Je Ausnahme und für die
  Schwelle ein Fall der Gegenprobe. — **Ausgang:** — (bei Closure)
- **Die Zahl wird zum Ziel** — Tests, die einen Operator töten, ohne eine Zusage zu
  prüfen, heben die Zahl, nicht die Prüfung. Das Gate ersetzt §3.10 nicht; Review
  bleibt Urteil darüber. — **Ausgang:** — (bei Closure)

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `e05e4bd` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (11×) und
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (10×) — der Anlass dieses Slice
  (§1 *Herkunft*); Sensor-Ausgang `geplant` mit dieser Kennung.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — Randform *Build-Kontext*
  und ein Risiko in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — dieser Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code ausdrücklich.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×) — Entscheidungsort der
  Randformen ist der Abschnitt für Harness-Werkzeuge der Spezifikation, angelegt von
  `slice-harness-vertraege-spezifikation`; die ADR trägt Entscheidung und Gründe.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×) und
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — je ein Risiko oder ein Fall der
  Gegenprobe; das Gate kommt hinzu, keine bestehende Prüfung der Stufe `test` fällt
  weg.

Keiner der Einträge unter der Schwelle erreicht mit diesem Slice allein 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
