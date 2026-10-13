# Slice slice-harness-vertraege-umzug-lint: Der Vertrag des Lint-Profils nach `harness/vertraege/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom 2026-10-13 ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)): Der Umzug der Verträge der Harness-Werkzeuge geschieht direkt hinter `slice-harness-gate-kette-stop-hook-freigabe` und vor `slice-v1-abschluss-tls-client`, geschnitten in drei Slices nach Vertrag (WIP-Limit 1, Reihenfolge: `slice-harness-vertraege-umzug-ordner`, `slice-harness-vertraege-umzug-lint`, `slice-harness-vertraege-umzug-gate-kette`); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode (4)). Bindung an Entscheidungen: [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) (Ort) mit Adaption [MR-002](../../../../harness/conventions.md#mr-002); der Vertrag bleibt von [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) geschärft.

**Berührte Spec-Stellen:** `SPEC-049` (Vertrag von `lint`: der Ort wechselt von `spezifikation.md §11` nach `harness/vertraege/`, der Text bleibt) · `SPEC-052` (nur die Ortsangabe in `AGENTS.md` §3.12 bleibt für ihn stehen; Umzug im Folge-Slice) · `spezifikation.md §11` (*Harness-Werkzeuge*, verliert diesen Vertrag) · `spezifikation.md §12` (*Historie*, verliert neun Zeilen)

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

**Ziel:** Der Vertrag des Lint-Profils, `SPEC-049` (152 Zeilen, zehn Punkte), steht als Datei in `harness/vertraege/`, im Text unverändert gegenüber
`spezifikation.md §11` am Stand `231089e`; er und die neun Zeilen seiner *Historie* sind aus der Spezifikation verschwunden, und jeder Verweis auf ihn zeigt
auf den neuen Pfad.

**Sendung aus `slice-harness-vertraege-umzug-ordner` (Zwischenstand mit zwei Orten).** `AGENTS.md` §3.12 nennt als Ort der Randform-Entscheidung `harness/vertraege/` und für die Verträge `SPEC-049` und `SPEC-052` `spec/spezifikation.md` §11. Dieser Slice nimmt `SPEC-049` an: Nach seinem Umzug nennt §3.12 `SPEC-049` nicht mehr als Ausnahme (der Ort ist `harness/vertraege/`); die Ausnahme für `SPEC-052` bleibt stehen, bis `slice-harness-vertraege-umzug-gate-kette` sie streicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Ordner, `README.md`, Rang und Reihe der Kennungen: Sie entstehen in `slice-harness-vertraege-umzug-ordner` und liegen vor diesem Slice (Voraussetzung: Er startet erst, wenn jener in `done/` liegt). Dieser Slice trägt nur die Historie-Zeilen und den Eintrag des Vertrags in das `README.md` nach.
- Die Ausnahme für `SPEC-052` in `AGENTS.md` §3.12 und der Wegfall des Zwischenstands: `slice-harness-vertraege-umzug-gate-kette` (Klasse 1; Sendung aus `slice-harness-vertraege-umzug-ordner`). Grund: Solange `SPEC-052` in §11 steht, ist die Ausnahme wahr.
- Der Vertrag der Gate-Kette und Einleitung und Rest von §11: `slice-harness-vertraege-umzug-gate-kette` (Klasse 1). Grund: Der Vertrag ist mit 262 Zeilen der größte und trägt eine Klarstellung (V-343), die einen eigenen Slice füllt.
- Der Text von `SPEC-049` wird nicht umgeschrieben (Entscheidung des Nutzers: wortgleich verschieben); der neue Punkt, den `slice-harness-lint-warten-ohne-frist` schreibt, entsteht dort, in der Datei am neuen Ort.
- Das Lint-Profil `.golangci.yml`, `tools/harness/lint.sh`, `tools/harness/lint-gegenprobe.sh`: nur Kommentare mit dem alten Pfad ändern sich, keine Regel, kein Fall, kein Mutant (Klasse 4: kein Produkt-Code, und die Gegenprobe prüft den Vertrag, nicht seinen Ort).
- Kommentare in `internal/**/export_test.go` nennen `SPEC-049 Punkt 7` ohne Pfad und bleiben (Klasse 2: Die Kennung löst weiter auf).
- Zeitdokumente: Links in Reviews und `done/` nur dort, wo sie sonst brechen; Fließtext bleibt ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), akzeptiertes Negativ).

## 2. Definition of Done


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 — die Datei.** `harness/vertraege/SPEC-049-lint-profil.md` trägt den Text von `SPEC-049` aus `spezifikation.md §11` am Stand `231089e` (Zeilen 3012 bis 3163),
      unverändert bis auf die Abweichungen, die `slice-harness-vertraege-umzug-ordner` für alle Verträge festlegt (Dateiname, Überschriftsebene); das `README.md` des Ordners führt den Vertrag
      und die neun Zeilen seiner Historie (Wortlaut aus §12). Beleg in §7: der Vergleich (Befehl und Ausgabe).
- [x] **Liefer-Punkt 2 — die Spezifikation.** `spec/spezifikation.md` führt `SPEC-049` und die neun Zeilen seiner *Historie* nicht mehr; sonst unverändert; kein Verweis der Spezifikation
      auf `harness/`; `make docs-check` grün.
- [x] **Liefer-Punkt 3 — die Verweise.** Jeder Verweis auf `SPEC-049` zeigt auf den neuen Ort: die Links in `harness/README.md`, `harness/sensors/lint.md`,
      [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), den Plänen in `open/`, den Plänen in `done/` und den Reviews, wo sie sonst bricht; die Pfadangaben in den Kommentaren von
      `.golangci.yml`, `tools/harness/lint.sh`, `tools/harness/lint-gegenprobe.sh`; die Ortsangabe in `slice-harness-lint-warten-ohne-frist` (Kopf, §1, §3, §6).
      `AGENTS.md` §3.12 (Sendung aus `slice-harness-vertraege-umzug-ordner`): nennt `SPEC-049` nicht mehr mit `spec/spezifikation.md` §11, die Ausnahme gilt nur noch für `SPEC-052`.
      Beleg in §7: `grep -rn "spezifikation.md#spec-049" --exclude-dir=.git --exclude-dir=.harness .` und `grep -rn "spezifikation.md SPEC-049" --exclude-dir=.git .` ohne Treffer.
- [x] `make gates` grün.
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
| `harness/vertraege/SPEC-049-lint-profil.md` | neu | Text aus §11 (Zeilen 3012 bis 3163 am Stand `231089e`) |
| `harness/vertraege/README.md` | update | Eintrag des Vertrags und neun Zeilen Historie (Wortlaut aus §12, Reihenfolge nach Datum) |
| `spec/spezifikation.md` | update | `SPEC-049` und neun Zeilen der *Historie* entfernt |
| `harness/README.md` | update | §Sensors: Links und Nennungen in den Zeilen von `make lint` und `make lint-gegenprobe` |
| `harness/sensors/lint.md` | update | drei Links auf den Vertrag |
| `AGENTS.md` | update | §3.12: Ausnahme für `SPEC-049` gestrichen, Ausnahme für `SPEC-052` bleibt (Sendung aus `slice-harness-vertraege-umzug-ordner`, Randform 7) |
| `.golangci.yml`, `tools/harness/lint.sh`, `tools/harness/lint-gegenprobe.sh` | update (nur Kommentare) | „Vertrag: spec/spezifikation.md SPEC-049“ nennt den neuen Pfad (`harness/mk/lint.mk` nennt keinen Pfad und bleibt, Randform 5) |
| `docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md` | update (Link) | Referenz-Nachzug nach [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), keine Überschreibung der Entscheidung |
| `docs/plan/planning/open/slice-harness-lint-warten-ohne-frist.md`, `slice-tests-ueberlebende-mutanten.md`, `slice-tests-ueberlebende-mutanten-driving.md` | update | Links und Ortsangabe des Vertrags (`AGENTS.md` §3.9) |
| `docs/plan/planning/done/slice-harness-blackbox-driven.md`, `-blackbox-einstieg.md`, `-blackbox-kern.md`, `-blackbox-pgwire.md`, `-lint.md`, `-lint-werkzeug.md`, `slice-lint-bestand-driving.md`, `slice-lint-bestand-kern-driven.md`; `docs/reviews/2026-10-06-review-slice-harness-lint-werkzeug.md`, `…-verifikation-slice-harness-blackbox-kern.md`, `…-verifikation-slice-harness-lint-werkzeug.md`, `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-driven.md` | update (Links) | nur die Pfade der Links, weil das Doku-Gate sonst bricht |

Ansatz: wie in `slice-harness-vertraege-umzug-ordner`: Der Text wandert per Ausschneiden, der Vergleich in §7 belegt es. Es entsteht kein Code, kein Test und kein Mutant; die
Gegenprobe des Lint-Gates (`make lint-gegenprobe`) läuft im Gate-Lauf unverändert und belegt, dass die Kommentaränderungen in `.golangci.yml` und den Skripten das Gate nicht verändern.
Die Zahl der Verweise ist hoch (18 Dateien außerhalb der Spezifikation mit Link, rund 22 Zeilen) und mechanisch; sie ändert die Größe nicht, weil sie unter Liefer-Punkt 3 eine Aufgabe zählt.

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-vertraege-umzug-ordner` liegt in `done/` (WIP-Limit 1). Die Randformen sind am 2026-10-13 vom Architect entschieden (§6 hier und im Ordner-Slice); der Implementer liest sie vor dem ersten Commit an Spezifikation oder
`harness/` und gibt eine fehlende Form an den Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von
  `git diff --numstat <Basis>..<letzter Commit des Slice> -- spec harness AGENTS.md tools .claude .golangci.yml docs/plan docs/reviews` liegt über 600
  (`<Basis>` ist der Commit des Übergangs `next → in-progress` auf dem Hauptzweig; Metrik nach `AGENTS.md` §3.14; Schätzung des Planners 350 bis 480: 152 Zeilen Vertragstext, doppelt gezählt, dazu Historie und rund 22 Zeilen Verweise; kein Code, der Faktor 2 aus §3.14 entfällt).
- `in-progress` → `open` (blockiert — Carveout?): Eine Änderung der Kommentare in `.golangci.yml` oder in den Skripten ändert die Antwort von `make lint` oder `make lint-gegenprobe`
  (dann ist der Kommentar ein Teil der Zusage und der Fall gehört dem Architect), oder `slice-harness-vertraege-umzug-ordner` schließt nicht.

## 5. Closure-Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (die Dateien unter `harness/vertraege/` liegen nicht auf der Positivliste des Stop-Hooks), Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein Link in einem Zeitdokument bricht erst im Doku-Gate und nicht beim Umzug des Vertrags; der Slice zählt die Links vorher (`grep`, 18 Dateien außerhalb der Spezifikation) — **Ausgang:** offen bis Closure.
- Pfadangaben in Kommentaren ohne Link fängt kein Gate (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** offen bis Closure.
- Ein Kommentar in `.golangci.yml` ist Teil der Datei, die `make lint` und seine Gegenprobe lesen; eine Änderung des Wortlauts darf deren Antwort nicht ändern — **Ausgang:** offen bis Closure.

**Prüfung des Architect vor dem Code** (2026-10-13, Stand `795f1c0`): Bezug zu [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Schärft `SPEC-049`), [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) und [MR-002](../../../../harness/conventions.md#mr-002) bestätigt. Nachgemessen am Stand `231089e`: `SPEC-049` 3012 bis 3163 (152 Zeilen), neun Historie-Zeilen 3514 bis 3522, jede nennt nur `SPEC-049`. Größe bestätigt, in der Metrik von `AGENTS.md` §3.14: 152 Zeilen doppelt 304, neun gelöschte Historie-Zeilen und rund zehn im `README.md` des Ordners, rund 22 Zeilen Verweise doppelt 44: **370 bis 430**, innerhalb der Schätzung des Planners (350 bis 480), unter der Marke 600. Zählung 3 Liefer-Punkte und 2 Schichten bestätigt. Kein Schnitt, keine Folge-ADR.

**Randformen** (`AGENTS.md` §3.12): Es gelten die Entscheidungen 1 bis 11 aus §6 von `slice-harness-vertraege-umzug-ordner` (Dateiname, Überschrift, Text mit Ortsangabe, Einleitung, Regel der Historie-Zeilen, Reihe, Rang, Stop-Hook, Zeitdokumente, Skript-Kommentare, Stand von `commit-struktur-id`). Zusätzlich für diesen Vertrag, **entschieden vom Architect am 2026-10-13**, Ort: dieser Plan:

1. *Dateiname.* `SPEC-049-lint-profil.md`, bestätigt.
2. *Neun Zeilen Historie.* Genau die Zeilen 3514 bis 3522 am Stand `231089e` (2026-10-06 bis 2026-10-08), die erste („Harness-Werkzeuge: Lint-Profil mit Linter, Schwellen, …“) eingeschlossen. Sie stehen im `README.md` nach der Regel der Einsortierung aus Randform 5 des Ordner-Slice (Datum aufsteigend; am selben Tag die Zeile „Abschnitt angelegt“ zuerst, die übrigen in der Reihenfolge der Spezifikation): die fünf Lint-Zeilen vom 2026-10-06 folgen also auf „Abschnitt angelegt“.
3. *Vertrag nennt eine Kennung des Produkts oder einen Pfad.* Entschieden: Der Text von `SPEC-049` nennt keinen Pfad der Spezifikation und kein `harness/` (Zeilen 3012 bis 3163 gelesen); ein Befund beim Vergleich ist eine Randform, die der Implementer zurückgibt.
4. *Plan `slice-harness-lint-warten-ohne-frist`.* Ortsangabe des Vertrags (Kopf, §1, §3, §6) ist `harness/vertraege/SPEC-049-lint-profil.md`, und die Historie des Vertrags steht in `harness/vertraege/README.md`; die Nennung `spezifikation.md` §12 im Kopf bleibt für die Historie der Spezifikation (die Zeile zu `SPEC-038`) und zieht mit `slice-harness-vertraege-umzug-gate-kette` auf §11 (Nummer der Historie nach dem Wegfall von §11). Sonst nichts (Bindung nach `AGENTS.md` §3.13: der Nehmer nennt den Ort mit der Kennung dieses Gebers).
5. *Kommentare.* `harness/mk/lint.mk` nennt `SPEC-049` ohne Pfad (Zeilen 5 und 17) und bleibt unverändert; der Pfad steht in `.golangci.yml` (Zeile 2), `tools/harness/lint.sh` (Zeile 5) und `tools/harness/lint-gegenprobe.sh` (Zeile 4). Der Kommentar in `.golangci.yml` bleibt **eine** Zeile, damit die Zeilen der Datei, die `make lint` und die Gegenprobe nennen, sich nicht verschieben; dass die Antwort der Gegenprobe gleich bleibt, belegt der Gate-Lauf.
6. *Weitere Pläne mit Link.* Neben `slice-tests-ueberlebende-mutanten` und `-driving` (Link auf `SPEC-049` Punkt 7) zählt der `grep` aus Liefer-Punkt 3 jede weitere Datei; was er findet und nicht in §3 steht, wird im selben Commit nachgezogen und in §7 genannt.
7. *Ortsangabe in `AGENTS.md` §3.12 (Sendung aus `slice-harness-vertraege-umzug-ordner`, `AGENTS.md` §3.13).* Zusage: Nach diesem Slice nennt §3.12 `spec/spezifikation.md` §11 nur noch für `SPEC-052`; `SPEC-049` steht unter `harness/vertraege/`. Beleg: `grep -n "SPEC-049" AGENTS.md` zeigt keine Zeile, die §11 als Ort nennt. Ein Rest, der `SPEC-049` weiter in §11 verortet, ist ein Befund gegen diesen Slice. Die Zeile ist Verdrahtung des Ortes im Harness, keine neue Schicht und kein neuer Liefer-Punkt (sie gehört zu Liefer-Punkt 3).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Die Belege aus §2 (Vergleich gegen den Stand `231089e`, `grep` der alten Pfade) stehen hier, nicht im Bericht.

### Belege des Implementers

Basis der Größenprobe: `f4b6d6d` (Move `next → in-progress`). Commits: `cdefb4c` (Vertrag, Spezifikation, Historie), `2b2a1cc` (Verweise, Kommentare, `AGENTS.md`, Plan `slice-harness-lint-warten-ohne-frist`).

**Läufe.** `make gates` auf `2b2a1cc`: Exit 0, `lint: 0 issues`, `lint-gegenprobe: gruen — Grundlauf ohne Befund, jeder Fall je Punkt von SPEC-049 wie erwartet`, `stop-hook-gegenprobe: gruen — 138 Fälle`. `make docs-check` auf demselben Stand: `714 Datei(en) geprüft, 0 Befund(e)`. Die Kommentare in `.golangci.yml`, `lint.sh` und `lint-gegenprobe.sh` ändern die Antwort der Gegenprobe nicht (derselbe Wortlaut des Ergebnisses wie vor dem Slice; `.golangci.yml` bleibt bei einer Kommentarzeile an Zeile 2).

**Liefer-Punkt 1 — wortgleich.** Vergleich gegen `231089e` (Zeilen 3012 bis 3163; Abweichung nur `#` statt `###` in der Überschrift und keine Leerzeile am Ende):

```
git show 231089e:spec/spezifikation.md | sed -n '3012,3163p' | sed '1s/^### /# /' | diff - harness/vertraege/SPEC-049-lint-profil.md
```

Ausgabe leer. Die neun Historie-Zeilen: `git show 231089e:spec/spezifikation.md | sed -n '3514,3522p'` gegen die neun Zeilen im `README.md` des Ordners (Einsortierung nach `Abschnitt angelegt`, vor der Zeile von `SPEC-050`): wortgleich. `grep -n "harness/\|spezifikation.md" harness/vertraege/SPEC-049-lint-profil.md` und `grep -n "](" …`: keine Ausgabe (Randform 3: der Text nennt keinen Pfad).

**Liefer-Punkt 2 — Spezifikation.** `git diff 231089e HEAD --stat -- spec`: `1 file changed, 347 deletions(-)`, keine Einfügung (185 aus `slice-harness-vertraege-umzug-ordner`, 162 hier: 152 Zeilen Vertrag samt Leerzeile davor, neun Historie-Zeilen). `grep -c "harness/vertraege" spec/spezifikation.md` ergibt 0. Die Zeile der Historie vom 2026-10-11 (`modernc.org/sqlite` in der Erlaubnisliste, nennt `SPEC-049` als eine von mehreren Änderungen) gehört nicht zu den neun Zeilen und bleibt (Randform 2); sie nennt keinen Pfad.

**Liefer-Punkt 3 — Verweise.** `grep -rn "spezifikation.md#spec-049" --exclude-dir=.git --exclude-dir=.harness .` trifft nur Text, keinen Link: `harness/vertraege/SPEC-050-struktur-kennungen-im-commit.md:25` (Beispiel-Anker, wortgleich nach Randform 3 des Ordner-Slice), `done/slice-harness-vertraege-umzug-ordner.md:161,172` (Zeitdokument, dieselbe Stelle) und diesen Plan (Zeilen 67, 94: die Greps und das Zitat des alten Kommentars). `grep -rn "spezifikation.md SPEC-049" --exclude-dir=.git .` trifft nur diesen Plan (Zeile 67, 94). Gezählt vorher: 17 Dateien mit Link außerhalb der Spezifikation (§6); nachgezogen sind 18 (`harness/README.md`, `harness/sensors/lint.md`, [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), drei Pläne in `open/`, acht in `done/`, vier Berichte in `docs/reviews/`); die Zählung des Plans lag um eine Datei darunter, alle 18 stehen in §3. Der Diff hat 25 Dateien mit Vertrag, Historie, Spezifikation, `AGENTS.md`, drei Kommentardateien und diesem Plan. Der `grep` aus Randform 6 fand keine Datei, die nicht in §3 steht. Plan `slice-harness-lint-warten-ohne-frist`: Kopf, §1, §3 (Zeile für den Vertrag und Zeile für `SPEC-038`/§12 getrennt), §6 nennen `harness/vertraege/SPEC-049-lint-profil.md`; §1 trägt `Ort seit slice-harness-vertraege-umzug-lint` (Bindung nach `AGENTS.md` §3.13). `grep -n "SPEC-049" AGENTS.md`: keine Zeile nennt `SPEC-049` mit §11 (§3.12 nennt `SPEC-052`).

**Zusage · Mutation · roter Lauf** (einzeln gefahren, danach `git checkout`):

| Zusage | Mutation | roter Lauf |
|---|---|---|
| der Text ist wortgleich | ein `X` vor Zeile 3 von `SPEC-049-lint-profil.md` | der Vergleich oben meldet `3c3` mit beiden Zeilen |
| jeder Link auf den Vertrag zeigt auf den neuen Ort | `harness/sensors/lint.md:8` zurück auf `../../spec/spezifikation.md#spec-049…` | `make docs-check`: 1 Befund, `anchor-missing` |
| die Spezifikation verweist nicht auf `harness/` | Link auf `../harness/vertraege/README.md` am Ende von `spec/spezifikation.md` | `make docs-check`: `matrix-forbidden — Referenz spec-straten → aussen ist nicht erlaubt` |
| die Spezifikation führt `SPEC-049` nicht mehr | Zeile `### SPEC-049 — Lint-Profil (`lint`)` am Ende von `spec/spezifikation.md` | `grep -c "^### SPEC-049" spec/spezifikation.md` meldet 1 statt 0 (der Beleg ist ein `grep`, kein Gate) |
| `AGENTS.md` §3.12 nennt `SPEC-049` nicht mehr als Ausnahme | „für die Verträge `SPEC-049` und“ wieder eingesetzt | `grep -n "SPEC-049" AGENTS.md` meldet Zeile 305 (ein `grep`, kein Gate) |
| die neun Historie-Zeilen stehen im `README.md` | die Zeile vom 2026-10-08 gelöscht | `grep -c` auf die neun Zeilen meldet 8 statt 9 |

Rang- und Ortsaussagen in Prosa haben kein Gate; die zwei `grep`-Zeilen oben sind der Beleg, nicht eine Sperre. Ein Kommentar mit dem alten Pfad (`lint.sh`, `lint-gegenprobe.sh`, `.golangci.yml`) fängt kein Gate (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); der Nachzug ist per `grep` belegt (Ausgabe der drei Dateien nennt `harness/vertraege/SPEC-049-lint-profil.md`).

**Größe** (Metrik `AGENTS.md` §3.14, `git diff --numstat f4b6d6d..2b2a1cc -- spec harness AGENTS.md tools .claude .golangci.yml docs/plan docs/reviews`, Summe aus Einfügungen und Löschungen): **393** (+198, −195) vor diesem Beleg-Commit; Marke 600 nicht erreicht, Schätzung des Architect 370 bis 430 getroffen. Eine Review-Runde steht noch aus (Schätzung ≈ 15 % darunter, Lerneintrag aus dem Ordner-Slice).

**Rückgaben an den Architect:** keine; keine Randform außerhalb von §6 aufgetreten.

**Beobachtung:** keine Evidence-Datei vom Implementer vorgeschlagen (Urteil der Closure).

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

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` am Stand `231089e` gesichtet. Treffer: `werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen: Ort der Festlegung) ·
`plan-folgt-korrektur-nicht` (41×, verkörpert: Ortsangaben in Plänen ziehen mit) · `zusage-im-kommentar-weiter-als-pruefung` (51×, verkörpert: Kommentare nennen den Pfad nur, soweit ein `grep` ihn belegt) ·
`schnitt-laesst-haelfte-an-der-grenze` (2×, offen: dieser Slice ist die mittlere Hälfte des Schnitts, Größe bekannt und kleiner als die beiden anderen) · `slice-waechst-durch-uebernahmen` (3×, verkörpert: Nachzählung unten).

**Nachgezählt** (`AGENTS.md` §3.13): Liefer-Punkte 3 (Datei, Spezifikation, Verweise). Schichten 2: Spezifikation (`spec/`) und Harness (`harness/`, Skripte, `.golangci.yml`); Planung (ADR, Pläne, Reviews) ist keine Schicht, Produkt-Code unberührt (die Kommentare in `internal/` bleiben). **Mit der Sendung aus `slice-harness-vertraege-umzug-ordner` nachgezählt:** `AGENTS.md` §3.12 gehört zu Liefer-Punkt 3 (Verweise) und liegt in der Schicht Harness; es bleiben 3 Liefer-Punkte und 2 Schichten. Größe: eine Zeile mehr, Schätzung 370 bis 440, Marke 600 unberührt.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
