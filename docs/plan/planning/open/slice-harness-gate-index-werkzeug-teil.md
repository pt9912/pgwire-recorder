# Slice slice-harness-gate-index-werkzeug-teil: Gate-Index um den Teil des Werkzeugs teilen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Angelegt von
`slice-harness-upgrade-v6-16` als Ausgang eines Befunds aus dem Abgleich
v6.13.0 → v6.16.0; Start in §4.

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage: „keine inhaltlichen Adaptionen“ — der Slice hält sie wahr, sobald ein Werkzeug seinen Teil des Gate-Index schreibt). Keine Anforderung des Lastenhefts im Scope: Der Slice ändert den Harness-Einstieg, nicht das Produkt.

**Berührte Spec-Stellen:** —

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

**Ziel:** Schreibt ai-harness-init einen eigenen Teil des Gate-Index unter
`harness/mk/` (Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt, *Ein Index, mehrere Eigentümer*, seit
`v6.16.0`), dann führt `harness/README.md` §Sensors jedes Target genau einmal: die
Zeilen der Werkzeug-Targets stehen nur noch im Teil des Werkzeugs, eine Zeile unter den
Tabellen verlinkt ihn und trägt, was das Repo über dessen Targets entscheidet, und
`AGENTS.md` §4 nennt den Teil.

**Herkunft:** Befund des Abgleichs in `slice-harness-upgrade-v6-16`: Der Stand
`v6.16.0` legt den Teil eines Werkzeugs in dessen Verantwortung; die Fragmente unter
`harness/mk/` stammen teils von ai-harness-init, teils vom Repo, und
`harness/README.md` §Sensors führt heute die Targets beider. Solange das Werkzeug
keinen Teil schreibt, ist das die eine Datei, die die Regel für diesen Fall zulässt.
Entscheidung des Nutzers vom 2026-10-06: Der Befund bleibt ein Slice und wird kein
Eintrag im Beobachtungs-Register, auch nach Review F-438 (die Lesart trägt der
Baseline-Text nicht ausdrücklich, und der Start hängt an einem externen Lauf ohne
Datum); bis zum Start ruht die `MR-000`-Aussage auf dieser Lesart.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der erneute Bootstrap mit ai-harness-init selbst — ein anderer Vorgang: Er
  schreibt die tool-eigenen Fragmente und den Teil des Werkzeugs; dieser Slice zieht
  nur die Dateien des Repos nach, die danach doppelt führten.
- Ein Sensor auf Disjunktheit der Teile — ein anderer Vorgang: Die Baseline nennt
  die fehlende Prüfung als Grenze; ob das Repo sie baut, ist eine eigene Entscheidung
  mit eigener ADR.
- Die Targets der repo-eigenen Fragmente (etwa `make kopf-check`, `make lint`) —
  Bestand bleibt bewusst stehen: Sie gehören dem Repo und bleiben in
  `harness/README.md` §Sensors.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Kein Target steht zugleich in `harness/README.md` §Sensors und im Teil des
      Werkzeugs; jede Zeile, die das Repo über ein Werkzeug-Target entschieden hat
      (Bindung des Repos, Carveout, Sensor-Datei), steht in der Zeile, die den Teil
      verlinkt.
- [ ] `AGENTS.md` §4 nennt den Teil des Werkzeugs, den §Sensors verlinkt.
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
| `harness/README.md` §Sensors | update | Zeilen der Werkzeug-Targets entfernen, Link-Zeile auf den Teil des Werkzeugs mit den Entscheidungen des Repos |
| `AGENTS.md` §4 | update | ein Satz: Targets aus Werkzeug-Fragmenten stehen in dem Teil, den §Sensors verlinkt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Ein Lauf von ai-harness-init hat einen Teil des
Gate-Index unter `harness/mk/` geschrieben, und sein Diff liegt zum Commit vor; und
welle-erster-release ist `done` (Entscheidung des Nutzers vom 2026-10-08, Wellen vor
Harness). Tritt der Lauf früher ein, wartet der Slice bis nach den Wellen; bis dahin
führt `harness/README.md` §Sensors die Targets weiter. Gegenüber der Reihe der übrigen
Harness-Slices (§4 von `slice-harness-commit-struktur-id`) steht er außerhalb.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Teil des Werkzeugs
  führt Targets, die das Repo nicht kennt, oder lässt welche weg, die das Repo als
  Werkzeug-Targets führt — dann ist erst die Zuordnung zu klären.
- `in-progress` → `open` (blockiert — Carveout?): Der Teil des Werkzeugs
  widerspricht einer Entscheidung des Repos über ein Target, die nicht in die
  Link-Zeile passt.

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

**Randformen:** offen — der Slice legt keinen neuen Vertrag an; welche Zeile was trägt,
entscheidet der Architect nach dem Start gegen den Teil, den das Werkzeug tatsächlich
schreibt.

- **Doppelnennung bleibt still** — die Vereinigung der Teile zählt ein doppelt
  genanntes Target einmal; kein Gate dieses Repos sieht sie. — **Ausgang:** — (bei
  Closure)

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

**Vorgelagert — offene Beobachtungen sichten:** beim Start nachzuholen; der Stand
bei der Anlage ist der von `slice-harness-upgrade-v6-16` §8.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
