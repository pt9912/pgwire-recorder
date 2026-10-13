# Slice slice-harness-vertraege-umzug-ordner: Ordner `harness/vertraege/` mit Rang 2 und die Verträge von `kopf-check`, `abdeckung` und `commit-msg`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom 2026-10-13 ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)): Der Umzug der Verträge der Harness-Werkzeuge geschieht direkt hinter `slice-harness-gate-kette-stop-hook-freigabe` und vor `slice-v1-abschluss-tls-client`, geschnitten in drei Slices nach Vertrag (WIP-Limit 1, Reihenfolge: `slice-harness-vertraege-umzug-ordner`, `slice-harness-vertraege-umzug-lint`, `slice-harness-vertraege-umzug-gate-kette`); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode (4): die Prüfumgebung des Projekts, deren Werkzeuge die Verträge festlegen). Bindung an Entscheidungen: [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) (Ort, Ablage, Rang, Gate-Klasse, Zeitpunkt) mit Adaption [MR-002](../../../../harness/conventions.md#mr-002); die Verträge bleiben von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) geschärft.

**Berührte Spec-Stellen:** `SPEC-047` · `SPEC-048` · `SPEC-050` (die Verträge von `kopf-check`, `abdeckung` und `commit-msg`: der Ort wechselt von `spezifikation.md §11` nach `harness/vertraege/`, der Text bleibt) · `spezifikation.md §11` (*Harness-Werkzeuge*, verliert diese drei Verträge) · `spezifikation.md §12` (*Historie*, verliert deren Zeilen)

**Verantwortlich:** pt9912

**Autor:** pt9912 (Planner). **Datum:** 2026-10-13.

---

## 1. Ziel und Abgrenzung


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Ordner `harness/vertraege/` besteht und trägt Rang 2 der Source Precedence neben `spec/spezifikation.md`; er führt die drei Verträge
`SPEC-047` (`kopf-check`), `SPEC-048` (`abdeckung`) und `SPEC-050` (`commit-msg`) je in einer Datei, im Text unverändert
gegenüber `spezifikation.md §11` am Stand `231089e`; die drei verschwinden samt ihren Zeilen der *Historie* aus der
Spezifikation, und jeder Verweis auf sie zeigt auf den neuen Pfad. Die Spezifikation verweist nicht auf `harness/`.

Warum diese drei zuerst: Sie sind die kleinen Verträge (zusammen 183 Zeilen), zwei von ihnen haben eine Sensor-Datei, die den
Vertrag verlinkt, und mit ihnen entsteht alles, was der Ordner einmal braucht (`README.md`, Rang, Reihe der Kennungen,
Ort der Historie). Die Verträge von Lint-Profil und Gate-Kette folgen in je einem eigenen Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Verträge des Lint-Profils und der Gate-Kette bleiben in `spezifikation.md §11`; sie ziehen in
  `slice-harness-vertraege-umzug-lint` und `slice-harness-vertraege-umzug-gate-kette` um (Klasse 1). Grund: Der Umzug aller fünf
  läge bei geschätzt 1700 (Summe aus Einfügungen und Löschungen) und über den drei Liefer-Punkten; ein Vertrag ist die kleinste
  Einheit, die einzeln lieferbar bleibt, und der Zwischenstand mit zwei Orten ist in `MR-002` vorgesehen. Die Einleitung von §11
  bleibt in der Spezifikation, bis der letzte Vertrag fort ist; `harness/vertraege/README.md` führt sie von Anfang an.
- Der Text der drei Verträge wird nicht umgeschrieben (Entscheidung des Nutzers: wortgleich verschieben). Berichtigungen, die der
  Ort erzwingt, benennt §6. Die Klarstellung zu Punkt 22 (b) der Gate-Kette (V-343) gehört zu deren Umzug und steht in
  `slice-harness-vertraege-umzug-gate-kette` (Klasse 1).
- `.d-check.yml` bleibt unverändert: Der Ordner steht in der Klasse `aussen` (Entscheidung des Nutzers, [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), Lesart (b)); eine Klasse oder
  ein Gate für „Rang 2“ entsteht nicht (Bestand bleibt bewusst stehen, Klasse 2).
- Zeitdokumente (Reviews unter `docs/reviews/`, Pläne in `done/`) werden nur dort geändert, wo ein Link sonst bricht; Ortsangaben im
  Fließtext bleiben (akzeptiertes Negativ der ADR, Klasse 2: Sie tragen den Stand ihrer Zeit).
- Pläne in `open/` und `next/`, die `spezifikation.md §11` als Ort *neuer* Kennungen nennen (`slice-harness-lh-links-pflicht`,
  `slice-harness-meldungskatalog-gate`, `slice-harness-coverage`, `slice-harness-mutation`), bleiben unberührt: Ihre Ortsangabe stimmt, bis
  §11 fällt; `slice-harness-vertraege-umzug-gate-kette` zieht sie nach (Klasse 1, dort in §1 und DoD).
- Kein Produkt-Code und kein Test: Kommentare in Skripten, die den alten Pfad nennen, ändern sich, Logik nicht (Klasse 4).
  Die Kommentare in `internal/` nennen die Kennung ohne Pfad und bleiben (Klasse 2: Die Kennung löst weiter auf).
- Kein Stub in der Spezifikation ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) schließt ihn aus: Die Spezifikation darf nicht auf `harness/` verweisen).

## 2. Definition of Done


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — der Ordner.** `harness/vertraege/README.md` (Gegenstand des Ordners, Rang 2 mit `spec/spezifikation.md`
      geteilt, Reihe der Kennungen über beide Orte mit der nächsten freien Nummer, Einleitung von §11, Historie der Verträge) und je eine Datei für
      `SPEC-047`, `SPEC-048`, `SPEC-050`; der Text ist gegenüber `spezifikation.md §11` am Stand `231089e` unverändert bis auf die
      Abweichungen, die §6 benennt. Beleg in §7: der Vergleich (Befehl und Ausgabe), nicht „im Bericht“.
- [ ] **Liefer-Punkt 2 — die Spezifikation.** `spec/spezifikation.md` führt die drei Verträge und die Zeilen ihrer *Historie* nicht
      mehr; §11 behält Einleitung und die beiden übrigen Verträge unverändert; die Spezifikation verweist nicht auf `harness/`;
      `make docs-check` grün (Klasse `spec-straten` gegen `aussen`).
- [ ] **Liefer-Punkt 3 — die Verweise.** Jeder Verweis auf die drei Verträge zeigt auf den neuen Ort: Links (Anker bleiben, der Pfad
      wechselt) in `harness/README.md`, den Sensor-Dateien, [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md), den Plänen in `open/` und `next/`
      und, wo sie sonst bricht, in Reviews und `done/`; `harness/README.md` §Source precedence (Rang 2 trägt beide Orte, neun Ränge)
      und §Sensors; `AGENTS.md` §2 (Rang 2) und §3.12 (Ort der Randform-Entscheidung); die Pfadangabe in den Kommentaren der Skripte.
      Beleg in §7: `grep -rn "spezifikation.md#spec-0\(47\|48\|50\)" --exclude-dir=.git --exclude-dir=.harness .` ohne Treffer, und
      der Text von `SPEC-047` nennt den Ort nicht (`MR-002` Adaption 5).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Handbuch und README unberührt: Das Produktverhalten ändert sich nicht (`AGENTS.md` §3.11).
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
| `harness/vertraege/README.md` | neu | Gegenstand, Rang 2 geteilt, Reihe über beide Orte (höchste vergebene Nummer, plus eins), Einleitung von §11, Historie der Verträge (Wortlaut der Zeilen aus §12) |
| `harness/vertraege/SPEC-047-kopf-der-plaene.md`, `SPEC-048-abdeckung.md`, `SPEC-050-struktur-kennungen-im-commit.md` | neu | Text aus §11 (Zeilen 2880 bis 2953, 2954 bis 3011, 3165 bis 3215 am Stand `231089e`); Dateinamen sind Vorschlag des Planners, der Architect entscheidet (§6) |
| `spec/spezifikation.md` | update | die drei Verträge aus §11 und die Zeilen der *Historie* mit `SPEC-047`, `SPEC-048`, `SPEC-050` entfernt; sonst nichts |
| `harness/README.md` | update | §Source precedence: Rang 2 trägt beide Orte; §Sensors: Links der Zeilen von `make kopf-check`, `make abdeckung-check` und der Zeile zu `SPEC-050` auf den neuen Pfad |
| `AGENTS.md` | update | §2 (Rang 2) und §3.12 (Ort der Randform-Entscheidung für Harness-Werkzeuge); die Herkunft in der Form `seit <Kennung>` setzt der Implementer bei der Closure |
| `harness/sensors/kopf-check.md`, `harness/sensors/abdeckung-check.md` | update | Links auf den Vertrag (MR-002 Adaption 8: die Lesart bleibt dort) |
| `tools/harness/kopf-check.sh`, `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung.sh`, `tools/test/abdeckung-gegenprobe.sh` | update (nur Kommentare) | „Vertrag: spec/spezifikation.md SPEC-0xx“ nennt den neuen Pfad |
| `docs/plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md` | update (Link) | Referenz-Nachzug nach [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), keine Überschreibung der Entscheidung |
| `docs/plan/planning/open/slice-harness-commit-struktur-id.md`, `docs/plan/planning/next/slice-harness-abdeckung-gate.md` | update | Kopf `Berührte Spec-Stellen`, §1, §3, §6: der Vertrag steht in `harness/vertraege/`; beide schreiben ihn später fort (`AGENTS.md` §3.9) |
| `docs/reviews/2026-10-06-verifikation-slice-harness-vertraege-spezifikation.md` | update (Links) | nur die Pfade der Links; der Befund bleibt |
| `docs/plan/planning/in-progress/roadmap.md` | update | Drift-Log und Stand erst bei der Closure |

Ansatz: Der Text der Verträge wandert per Ausschneiden (Zeilenbereich am Stand `231089e`), nicht per Neuschrift; der Vergleich in §7 belegt es.
Es gibt keinen neuen Vertrag und keinen Test: Die Zusage „wortgleich“ prüft der Vergleich, nicht ein Gate; README und Handbuch sagen nichts weiter zu
(`AGENTS.md` §3.11). Die Verweise der Spezifikation auf die Verträge sind nur noch die Zeilen der *Historie*, und sie wandern mit.

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gate-kette-stop-hook-freigabe` liegt in `done/` (Stand `231089e`); WIP-Limit 1. Erster Schritt nach dem
Start, vor jedem Commit an Spezifikation oder `harness/`: Der Architect entscheidet die Randformen aus §6 (Dateinamen, Überschriftsebene, Auswahl der
Historie-Zeilen, Einleitung, Reihe); wellenlos heißt nicht ohne Architect (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von
  `git diff --numstat <Basis>..<letzter Commit des Slice> -- spec harness AGENTS.md tools .claude docs/plan docs/reviews` liegt über 750
  (`<Basis>` ist der Commit des Übergangs `next → in-progress` auf dem Hauptzweig; Metrik nach `AGENTS.md` §3.14; Schätzung des Planners 450 bis 600: 183 Zeilen Vertragstext, doppelt gezählt, dazu README, Historie und Nachzug; der Faktor 2 für Tabellen- und
  Gegenprobe-Code entfällt, weil kein Code entsteht), oder der Architect entscheidet, den Text eines Vertrags beim Umzug umzuschreiben statt zu verschieben
  (dann prüft der Vergleich nichts mehr); in beiden Fällen schneidet der Planner `SPEC-050` ab.
- `in-progress` → `open` (blockiert — Carveout?): Das Doku-Gate verlangt für den Ordner eine Änderung an `.d-check.yml` (die Entscheidung des Nutzers lässt sie
  unverändert; dann entscheidet der Nutzer neu), oder ein Verweis der Spezifikation auf einen der drei Verträge lässt sich nur auflösen, indem Text der
  Spezifikation geändert wird, den dieser Slice nicht entscheidet.

## 5. Closure-Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (die Dateien unter `harness/vertraege/` liegen nicht auf der Positivliste des Stop-Hooks, der volle Lauf ist nötig),
Closure-Notiz mit Lerneintrag. Beobachtungs-Register: `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen) wird von diesem Slice berührt; ob eine
Evidence-Datei entsteht, entscheidet der Implementer an der Frage, ob der Ort noch einmal übersehen wurde (Urteil, kein Zähler).

## 6. Risiken und offene Punkte


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Text eines Vertrags nennt selbst den alten Ort (in `SPEC-050` steht ein Beispiel-Link `spezifikation.md#spec-049…` als Text, kein Link) — wortgleich heißt: bleibt, mit Entscheidung des Architect (§6) — **Ausgang:** offen bis Closure.
- Zwischenstand mit zwei Orten: Wer `SPEC-049` oder `SPEC-052` sucht, findet sie bis zu den Folge-Slices noch in §11; das `README.md` des Ordners sagt, welche Verträge der Ordner führt, und nur das (`AGENTS.md` §3.7: Zustand, keine Chronik) — **Ausgang:** offen bis Closure.
- Verweise ohne Link (Pfadangaben in Kommentaren, Fließtext in Plänen) fängt kein Gate; der Beleg ist ein `grep`, der nur diese Form kennt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** offen bis Closure.
- Die Schnitt-Grenze: Eine Hälfte des Umzugs (Einleitung von §11, Reihe der Kennungen) liegt von Anfang an im Ordner, obwohl zwei Verträge noch fehlen; läge ein Fehler dort, träfe er die beiden Folge-Slices (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, 2×, offen) — **Ausgang:** offen bis Closure.

**Randformen** (`AGENTS.md` §3.12; der Umzug liefert keinen neuen Vertrag, verschiebt aber Text, und jede Entscheidung dazu fällt vor dem ersten Commit; Ort: `harness/vertraege/README.md`, vom Architect):

1. *Dateiname.* Vorschlag: `SPEC-047-kopf-der-plaene.md`, `SPEC-048-abdeckung.md`, `SPEC-050-struktur-kennungen-im-commit.md` (ASCII, ohne Umlaut; `MR-002` Adaption 1 fordert `SPEC-<NNN>-<kurztitel>.md`).
2. *Überschrift.* Der Wortlaut bleibt (Anker bleiben); offen ist die Ebene: `###` wie in §11 oder `#` als Titel der Datei. Vorschlag: `#`.
3. *Text, der den Ort nennt.* `SPEC-050` nennt ein Beispiel `spezifikation.md#spec-049--lint-profil-lint` als Text; `SPEC-047` nennt den Ort nicht. Vorschlag: wortgleich lassen (ein Beispiel, kein Verweis), die Entscheidung steht in der Datei `README.md`.
4. *Einleitung von §11.* Sie nennt `LH-QA-07` als Link `lastenheft.md#…`; im `README.md` gilt `../../spec/lastenheft.md#…`. Die Einleitung bleibt in der Spezifikation bis zum letzten Umzug (zwei Orte, eine Dopplung, die der letzte Slice auflöst).
5. *Historie.* Eine Zeile aus §12 zieht um, wenn ihr Gegenstand ein Harness-Werkzeug ist; Zeilen mit gemischtem Gegenstand (Produkt und Werkzeug) bleiben. Für diesen Slice: die Zeilen vom 2026-10-06 (Abschnitt angelegt, `SPEC-047`, `SPEC-048`) und 2026-10-08 (Commit-Träger, `SPEC-050`). Der Wortlaut bleibt, die Reihenfolge nach Datum auch.
6. *Reihe.* Höchste vergebene Nummer über beide Orte ist 052; die nächste freie ist 053; `SPEC-051` (Schreibfrist der Fehlerantwort) bleibt in der Spezifikation.
7. *Rang ohne Gate.* `harness/README.md` und `AGENTS.md` sagen „Rang 2 trägt beide Orte“ und sagen nichts zu, was ein Gate prüfte (`MR-002` Adaption 4 und 5).
8. *Stop-Hook.* Dateien unter `harness/vertraege/` sind nicht auf der Positivliste (Punkt 12 der Gate-Kette): Jede Änderung verlangt den vollen Lauf; die Dateien von §7 zählen mit.
9. *Verweise in Zeitdokumenten.* Links, die brechen (Gate `anchors`), werden nachgezogen; Fließtext bleibt ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), akzeptiertes Negativ).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Die Belege aus §2 (Vergleich gegen den Stand `231089e`, `grep` der alten Pfade, Prüfung von `SPEC-047` auf den Ort) stehen hier, nicht im Bericht.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield); der Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` am Stand `231089e` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:
`werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen: der Ort der Werkzeug-Festlegung ist der Gegenstand; der Umzug ändert den Ort, nicht das Stratum) ·
`schnitt-laesst-haelfte-an-der-grenze` (2×, offen: dieser Plan ist ein Schnitt; die Größe ist bekannt und wird mit Marke und Schätzung in der Metrik von `AGENTS.md` §3.14 benannt, ein drittes Auftreten wäre die Schwelle) ·
`plan-folgt-korrektur-nicht` (41×, verkörpert in §3.9 und `make kopf-check`: der Nachzug der Pläne und der Kopf sind Auftrag dieses Plans) ·
`zusage-im-kommentar-weiter-als-pruefung` (51×, verkörpert in §3.11: „wortgleich“ und „Rang 2“ sind ohne Gate nur so weit zugesagt, wie der Vergleich in §7 reicht) ·
`spec-randform-erst-im-review-entschieden` (27×, verkörpert in §3.12: §6 nennt die Randformen) ·
`groessenprobe-ohne-messbaren-stand` und `schaetzung-in-zeilen-unterschaetzt-tests` (je 3×, verkörpert in §3.14: Marke mit Metrik in §4) ·
`slice-waechst-durch-uebernahmen` (3×, verkörpert in §3.13): Nachgezählt beim Anlegen, siehe unten.

**Nachgezählt** (`AGENTS.md` §3.13): Liefer-Punkte 3 (Ordner, Spezifikation, Verweise). Schichten 2: die Spezifikation (`spec/`) und das Harness (`harness/`, `AGENTS.md`, Skripte, Rollen-Dateien); Planung (ADR, Pläne, Reviews, Roadmap) ist keine Schicht, Produkt-Code und Handbuch bleiben unberührt.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
