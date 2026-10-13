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
      geteilt, Reihe der Kennungen über beide Orte als Regel, nicht als Zahl (Randform 6), Einleitung von §11, Historie der Verträge) und je eine Datei für
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
| `harness/README.md` | update | §Source precedence: Rang 2 trägt beide Orte (Randform 7); §Sensors: Links der Zeilen von `make kopf-check` und `make abdeckung-check` auf den neuen Pfad (zu `SPEC-050` führt §Sensors keine Zeile; es entsteht keine) |
| `AGENTS.md` | update | §2 (Rang 2) und §3.12 (Ort der Randform-Entscheidung für Harness-Werkzeuge); die Herkunft in der Form `seit <Kennung>` setzt der Implementer bei der Closure |
| `harness/sensors/kopf-check.md`, `harness/sensors/abdeckung-check.md` | update | Links auf den Vertrag (MR-002 Adaption 8: die Lesart bleibt dort) |
| `tools/harness/kopf-check.sh`, `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung.sh`, `tools/test/abdeckung-gegenprobe.sh` | update (nur Kommentare) | „Vertrag: spec/spezifikation.md SPEC-0xx“ nennt den neuen Pfad |
| `docs/plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md` | update (Link) | Referenz-Nachzug nach [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), keine Überschreibung der Entscheidung |
| `docs/plan/planning/open/slice-harness-commit-struktur-id.md`, `docs/plan/planning/next/slice-harness-abdeckung-gate.md` | update | Kopf `Berührte Spec-Stellen`, §1, §2 (DoD, „Historien-Zeile in §12“), §3, §6: der Vertrag und seine Historie stehen in `harness/vertraege/`; beide schreiben ihn später fort (`AGENTS.md` §3.9; Stand von `commit-struktur-id`: Randform 11). `abdeckung-gate` nennt den Ort auch im Fließtext („Abschnitt für Harness-Werkzeuge der Spezifikation“, §3, §6) |
| `docs/reviews/2026-10-06-verifikation-slice-harness-vertraege-spezifikation.md` | update (Links) | nur die Pfade der Links; der Befund bleibt |
| `docs/plan/planning/in-progress/roadmap.md` | update | Drift-Log und Stand erst bei der Closure |

Ansatz: Der Text der Verträge wandert per Ausschneiden (Zeilenbereich am Stand `231089e`), nicht per Neuschrift; der Vergleich in §7 belegt es.
Es gibt keinen neuen Vertrag und keinen Test: Die Zusage „wortgleich“ prüft der Vergleich, nicht ein Gate; README und Handbuch sagen nichts weiter zu
(`AGENTS.md` §3.11). Die Verweise der Spezifikation auf die Verträge sind nur noch die Zeilen der *Historie*, und sie wandern mit.

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gate-kette-stop-hook-freigabe` liegt in `done/` (Stand `231089e`); WIP-Limit 1. Die Randformen aus §6 sind am 2026-10-13 vom Architect entschieden; der Implementer liest §6 vor dem ersten Commit an Spezifikation oder `harness/` und gibt eine Form, die dort fehlt,
an den Architect zurück (wellenlos heißt nicht ohne Architect, `BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

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

- Der Text eines Vertrags nennt selbst den alten Ort (in `SPEC-050` steht ein Beispiel-Anker `spezifikation.md#spec-049…` als Text, kein Link) — wortgleich, entschieden in Randform 3 — **Ausgang:** offen bis Closure.
- Zwischenstand mit zwei Orten: Wer `SPEC-049` oder `SPEC-052` sucht, findet sie bis zu den Folge-Slices noch in §11; das `README.md` des Ordners sagt, welche Verträge der Ordner führt, und nur das (`AGENTS.md` §3.7: Zustand, keine Chronik) — **Ausgang:** offen bis Closure.
- Verweise ohne Link (Pfadangaben in Kommentaren, Fließtext in Plänen) fängt kein Gate; der Beleg ist ein `grep`, der nur diese Form kennt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** offen bis Closure.
- Die Schnitt-Grenze: Eine Hälfte des Umzugs (Einleitung von §11, Reihe der Kennungen) liegt von Anfang an im Ordner, obwohl zwei Verträge noch fehlen; läge ein Fehler dort, träfe er die beiden Folge-Slices (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, 2×, offen) — **Ausgang:** offen bis Closure.

**Prüfung des Architect vor dem Code** (2026-10-13, Stand `795f1c0`, Spezifikation unverändert seit `231089e`): Bezug zu [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) und [MR-002](../../../../harness/conventions.md#mr-002) bestätigt (Ort, Ablage, Rang, Gate-Klasse `aussen`). Zeilenbereiche am Stand `231089e` nachgemessen: `SPEC-047` 2880 bis 2953 (74), `SPEC-048` 2954 bis 3011 (58), `SPEC-050` 3165 bis 3215 (51), zusammen 183. Größe bestätigt, in der Metrik von `AGENTS.md` §3.14: 183 Zeilen Vertragstext doppelt (neu und gelöscht) 366, dazu `README.md` des Ordners rund 45 (Einleitung, Rang, Reihe, Tabelle der Verträge, zwei Zeilen Historie), die zwei gelöschten Historie-Zeilen, rund 15 Dateien Nachzug mit je ein bis drei geänderten Zeilen (doppelt gezählt rund 60), `AGENTS.md` und `harness/README.md` rund 20: **480 bis 560**, innerhalb der Schätzung des Planners (450 bis 600), unter der Marke 750. Zählung 3 Liefer-Punkte und 2 Schichten bestätigt. Kein Schnitt, keine Folge-ADR.

**Randformen** (`AGENTS.md` §3.12), **entschieden vom Architect am 2026-10-13** vor dem ersten Code-Commit. Ort: dieser Plan, weil `harness/vertraege/README.md` erst mit Liefer-Punkt 1 entsteht und der Umzug keinen neuen Vertrag liefert; das `README.md` trägt das Ergebnis als Zustand (Verträge, Rang, Reihe, Historie), nicht die Entscheidung. Die Entscheidungen 1 bis 9 gelten für alle drei Umzug-Slices.

1. *Dateiname.* Wie vorgeschlagen, ASCII ohne Umlaut: `SPEC-047-kopf-der-plaene.md`, `SPEC-048-abdeckung.md`, `SPEC-050-struktur-kennungen-im-commit.md` (hier), `SPEC-049-lint-profil.md` (`slice-harness-vertraege-umzug-lint`), `SPEC-052-gate-kette.md` (`slice-harness-vertraege-umzug-gate-kette`).
2. *Überschrift.* Ebene `#` als Titel der Datei, Wortlaut unverändert (`# SPEC-047 — Kopf der Pläne (\`kopf-check\`)`). Anker bleiben: Belegt am 2026-10-13 mit einer Probe-Datei (nicht committet) unter `harness/vertraege/`: ein Link auf den Slug einer `#`-Überschrift löst im Doku-Gate auf, ein falscher Slug meldet `anchor-missing`.
3. *Text, der den Ort nennt.* `SPEC-050` bleibt wortgleich, auch der Beispiel-Anker `spezifikation.md#spec-049--lint-profil-lint` in Punkt 3: Er zeigt die Schreibweise eines Ankers (klein, darum keine Struktur-Kennung), er verweist nicht. Belegt am 2026-10-13: derselbe Code-Span in einer Datei unter `harness/vertraege/` bleibt im Doku-Gate ohne Befund. `SPEC-047` und `SPEC-048` nennen den Ort nicht (Zeilenbereiche gelesen). Der Beleg `MR-002` Adaption 5 in §7 ist ein `grep -n "harness/\|spezifikation.md" harness/vertraege/SPEC-047-kopf-der-plaene.md` ohne Treffer.
4. *Einleitung von §11.* Das `README.md` führt sie wortgleich bis auf die vom Ort erzwungenen Stellen: „Dieser Abschnitt“ wird „Dieser Ordner“, „Keine Kennung dieses Abschnitts“ wird „Keine Kennung dieses Ordners“, der Link heißt `../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes`. Die Spezifikation behält ihre Einleitung bis `slice-harness-vertraege-umzug-gate-kette`; die Dopplung bis dahin ist entschieden und akzeptiert (zwei Slices lang), `slice-harness-vertraege-umzug-gate-kette` löscht die Fassung der Spezifikation.
5. *Historie.* Eine Zeile zieht um, wenn jede Kennung, die sie nennt, ein umgezogener Vertrag ist (zusätzlich `LH-QA-07` als Decke). Für diesen Slice genau zwei Zeilen am Stand `231089e`: Zeile 3524 (2026-10-06, Abschnitt angelegt, `SPEC-047`, `SPEC-048`) und Zeile 3525 (2026-10-08, Commit-Träger, `SPEC-050`). Gemischte Zeilen gibt es unter den 16 Zeilen der drei Slices nicht (je Zeile geprüft); die Zeile zu `SPEC-038` (Warten in Tests) und die zu `SPEC-034` zwischen den Lint-Zeilen bleiben, weil sie Verträge des Produkts nennen. Der Wortlaut bleibt, auch die Daten (zwei Zeilen tragen 2026-10-14). Im `README.md` steht die Historie als Tabelle `Datum | Änderung` unter `## Historie`; Datum aufsteigend, am selben Tag steht die Zeile „Abschnitt angelegt“ zuerst, die übrigen in der Reihenfolge der Spezifikation am Stand `231089e`.
6. *Reihe.* Die höchste vergebene Nummer über beide Orte ist `SPEC-052` (`SPEC-051` steht in der Tabelle der Defaults der Spezifikation, nicht als Überschrift und bleibt Vertrag des Produkts); die nächste freie ist **053**. Das `README.md` nennt sie als Regel („höchste vergebene Nummer in `spec/spezifikation.md` und in `harness/vertraege/`, plus eins“, `MR-002` Adaption 3) und nicht als Zahl, die mit jedem Vertrag veraltet; die Zahl 053 steht in den Plänen, die sie brauchen.
7. *Rang ohne Gate.* `harness/README.md` §Source precedence: Rang 2 nennt beide Orte (`spec/spezifikation.md` für das Produkt, `harness/vertraege/` für die Werkzeuge), die Charakter-Spalte bleibt „technisch fortschreibbar“, es bleiben neun Ränge; der Absatz „Die Ränge 1–3 sind die drei Spec-Straten“ sagt dazu, dass Rang 2 den Ordner außerhalb der Straten mitführt (Verweis auf `MR-002`). `AGENTS.md` §2 Punkt 2 nennt beide Orte. Beide Stellen sagen nichts, was ein Gate prüfte (`MR-002` Adaption 4 und 5).
8. *Stop-Hook.* Bestätigt: `harness/vertraege/` steht nicht auf der Positivliste, jede Änderung verlangt `make gates`.
9. *Zeitdokumente.* Bestätigt: Links, die das Doku-Gate bricht, werden nachgezogen; Fließtext in Reviews und `done/` bleibt.
10. *Verweise ohne Link in Skripten.* Ein `grep -rn "spezifikation.md SPEC-0" --exclude-dir=.git --exclude-dir=.harness .` sucht die Form „Vertrag: spec/spezifikation.md SPEC-0xx“ in `tools/`, `.claude/` und `.golangci.yml`; für die drei Verträge dieses Slice sind es `tools/harness/kopf-check.sh`, `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung.sh`, `tools/test/abdeckung-gegenprobe.sh`. Der Beleg in §7 nennt diesen `grep` mit Ausgabe (nur Kommentare, kein Link-Sensor; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`).
11. *Stand von `slice-harness-commit-struktur-id`.* Der Plan liegt in `open/` (Rückführung vom 2026-10-08, Priorität), nicht geliefert; „geliefert“ in seinem §3 meint die Hälfte des Architect (`SPEC-050` und [ADR-0035](../../adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)). Träger `.githooks/commit-msg` und Gegenprobe fehlen im Bestand. Der Umzug nimmt ihm nichts ab und gibt ihm nichts: Er zieht in Kopf, §1, §2 (DoD, „Historien-Zeile in §12“), §3 und §6 nur die Ortsangabe nach; die Historie-Zeile steht danach in `harness/vertraege/README.md`.

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
