# Slice slice-harness-slice-mv-exakter-treffer: `slice-mv` findet die Quelle am exakten Namen, und eine Gegenprobe hält es

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
nächsten Welle-Closure. Angelegt am 2026-10-15 nach Hinweis des Nutzers vom selben Tag
(ai-harness-init `v0.7.0` löst das Präfix-Problem von `slice-mv`); die Stelle in der
Reihenfolge ist Empfehlung des Planners (§4 *Start*); seit dem 2026-10-16 hinter dem Release (Entscheidung des Nutzers).

**Bezug:** kein Anspruch des Produkts — der Slice ändert ein Werkzeug des Harness. Auslöser:
[slice-harness-upgrade-v6-16](../done/welle-v1-abschluss/slice-harness-upgrade-v6-16.md) (dokumentierter Auslöser
der Übernahme) und der Beleg in
[der Verifikation von slice-harness-gate-kette-stop-hook-freigabe](../../../reviews/2026-10-13-verifikation-slice-harness-gate-kette-stop-hook-freigabe.md)
(Zeile 49); die Moves von Schritt 35b und 36a gingen von Hand, weil das Muster `-sqlite` gegen
`-sqlite-schreiber` den Präfix mehrdeutig macht. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make slice-mv SLICE=slice-<Kennung> TO=…` bewegt einen Slice auch dann, wenn eine zweite
Kennung die erste als Präfix trägt (`slice-v1-abschluss-verbindungsende-sqlite` gegen
`…-sqlite-schreiber`), und eine Gegenprobe im Gate-Lauf hält die Zusagen der Quellensuche. Dazu übernimmt
der Slice `tools/harness/slice-mv.sh` aus ai-harness-init `v0.7.0` (neue Funktion `quelle_finden`;
Fragment `harness/mk/slice-mv.mk` ist byte-gleich und bleibt): der exakte Dateiname (mit oder ohne
`.md`) gewinnt vor dem Präfix; ohne exakten Treffer gilt der Präfix an einer Bindestrich-Grenze
(`slice-a` trifft `slice-a-b.md`, nicht `slice-ax.md`); zwei Präfix-Treffer und derselbe Name in zwei
Verzeichnissen enden mit Exit 2, ebenso kein Treffer. Die übernommene Fassung verweist in einem
Kommentar auf eine Datei `test/slice-mv.bats`, die dieses Repo nicht führt; der Kommentar zeigt nach der
Übernahme auf die Gegenprobe dieses Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die übrigen Dateien aus `v0.7.0` (`selbstpruefung.sh`, `stop-require-gates.sh`, `.githooks/commit-msg`,
  `.d-check.yml`, `harness/conventions.md`, `harness/README.md` außer der Zeile für die neue Gegenprobe,
  `implement-slice.md`) — Bestand bleibt bewusst stehen: Dieses Repo führt eigene Fassungen davon; eine
  Übernahme wäre ein anderer Vorgang (Abgleich mit dem Upstream) mit eigenem Befund je Datei.
- Eine Änderung am Fragment `harness/mk/slice-mv.mk` — Bestand bleibt bewusst stehen: Es ist byte-gleich
  zu `v0.7.0` (Beleg: `diff` in §7).
- Eine Änderung an der Verweis-Mechanik von `slice-mv.sh` (Ersetzungs-Funktionen, Ausnahmeliste,
  Commit-Reihenfolge) — Bestand bleibt bewusst stehen: `v0.7.0` ändert sie nicht, und der Slice
  übernimmt die Fassung, statt sie zu entwickeln. Die Gegenprobe prüft die Quellensuche und den Pfad
  durch `main()`, soweit er die Quellensuche braucht, nicht die Verweis-Ersetzung.
- Eine Aufnahme von `make slice-mv` in die schnelle Kette oder in `GATE_CHECKS` als Gate — anderer
  Vorgang: `slice-mv` ist ein Werkzeug, kein Gate (`harness/README.md` §Sensors); nur seine Gegenprobe
  ist ein Gate.
- Die Handlung des Freshness-Audits gegen die Release-Liste des Kurs-Repos — anderer Vorgang: Dieser
  Slice nimmt eine Datei eines Werkzeug-Repos (ai-harness-init) auf, nicht einen neuen Stand der
  vendored Baseline.
- Produkt-Code, Benutzerhandbuch, Lastenheft, Spezifikation — Schicht-Abgrenzung: Der Slice ändert ein
  Skript, ein Make-Fragment, eine Zeile im Gate-Index und Planung.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas
mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Übernahme):** `tools/harness/slice-mv.sh` trägt die Fassung von ai-harness-init
      `v0.7.0` (Funktion `quelle_finden`, Hilfetext, Aufruf in `main()`); der Kommentar, der auf
      `test/slice-mv.bats` verweist, nennt die Gegenprobe dieses Slice. Beleg in §7: `diff` gegen `v0.7.0`
      zeigt nur diese Kommentar-Zeilen; `diff` des Fragments `harness/mk/slice-mv.mk` gegen `v0.7.0` leer.
      Ein Lauf von `make slice-mv` mit `SLICE=slice-v1-abschluss-verbindungsende-sqlite` und
      `…-sqlite-schreiber` zeigt, dass der exakte Name den Präfix schlägt (in einem Wegwerf-Klon, nicht im
      Arbeitsbaum).
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `tools/harness/slice-mv-gegenprobe.sh` mit dem Fragment
      `harness/mk/slice-mv-gegenprobe.mk` (Ziel `slice-mv-gegenprobe`, hängt an `GATE_CHECKS`) prüft in
      einem Wegwerf-Repo (`mktemp -d`) die Fälle: exakter Treffer vor Präfix, Bindestrich-Grenze, zwei
      Präfix-Treffer (Exit 2), kein Treffer (Exit 2), derselbe Name in zwei Verzeichnissen (Exit 2),
      dazu die Randformen aus §6 nach der Entscheidung des Architect; je Zusage eine Mutation der
      Datei in einer Kopie, die nur ihre Fälle rot färbt (`AGENTS.md` §3.10, Bericht: Zusage · Mutation ·
      roter Test); die Zeile in `harness/README.md` §Sensors führt das Ziel; die Ausnahmeliste der
      schnellen Kette und `harness/sensors/slice-mv.md` nach der Entscheidung des Architect (§6).
- [ ] Mutation je Zusage der Gegenprobe selbst (`AGENTS.md` §3.10), vom Implementer rot gesehen; der Bericht
      nennt je Zusage: Zusage · Mutation · roter Test.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/slice-mv.sh` | update (Übernahme `v0.7.0`) | Funktion `quelle_finden`, Hilfetext, Aufruf in `main()`: +49/−15 Zeilen; Kommentar zur Prüfdatei auf die Gegenprobe umgebogen |
| `tools/harness/slice-mv-gegenprobe.sh` | neu | Gegenprobe im Stil von `kopf-check-gegenprobe.sh`: Wegwerf-Repo, Fälle je Zusage, Mutationen in Kopien |
| `harness/mk/slice-mv-gegenprobe.mk` | neu | Ziel `slice-mv-gegenprobe`, `GATE_CHECKS += slice-mv-gegenprobe` |
| `harness/mk/gates-schnell.mk` | update (nach Entscheidung des Architect) | Ausnahmeliste der schnellen Kette, falls die Gegenprobe dort steht (§6) |
| `harness/README.md` §Sensors | update | Zeile für `make slice-mv-gegenprobe` |
| `harness/sensors/slice-mv.md` | neu (nach Entscheidung des Architect) | Beschreibung nur, wenn der Architect sie für nötig hält (§6) |

Tests und Gates: eine neue Gegenprobe (`slice-mv-gegenprobe`). Reihenfolge: Übernahme zuerst (ein Commit,
nur die übernommene Datei plus Kommentar), dann die Gegenprobe gegen sie, dann Gate-Verdrahtung und
README-Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice steht in der wellenlosen Reihe der fünf Harness-Slices an erster Stelle, hinter
[welle-v1-nachtrag](../welle-v1-nachtrag.md) (in `done/`; Voraussetzung der Stelle, nicht der Technik) und vor
`slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`,
`slice-harness-schaetzung-planner-faktor` und `slice-harness-gegenprobe-matrix-p7-flanke`. Er stand bis zum 2026-10-16
unmittelbar hinter `slice-v1-abschluss-verbindungsende-doku` (Schritt 42 von welle-v1-abschluss, Empfehlung des Planners vom 2026-10-15)
und ist seitdem hinter das Release verschoben (Entscheidung des Nutzers vom 2026-10-16, Drift-Log der Roadmap). Er blockiert kein
Gate und keinen Schritt; bis dahin, also bis nach dem Release, gehen Moves von Hand, wenn das Präfix mehrdeutig ist. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die
  Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness .claude`,
  wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist und beim Start eingetragen wird; **über 500** geht der Slice
  zurück zur Zerlegung. **Meldeschwelle 400:** Erreicht die Summe 400, hält der Implementer an und meldet dem Planner,
  bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 280 bis 370,
  Mitte 325, 65 % der Marke. Aufbau: Skript 64 (gemessen: +49/−15 gegen die Fassung des Repos), Fragment und README-Zeile rund 15,
  Gegenprobe 170 bis 250, optional `harness/sensors/slice-mv.md` rund 35. Faktor Tests zu Produktcode: 5 auf die neu geschriebene Funktion
  `quelle_finden` (rund 42 Zeilen) — Beleg: Tabelle von Randformen mit einer Mutation je Zusage (`AGENTS.md` §3.10); als
  Mindestwert für diese Art Plan angesetzt, ein gemessener Vorgänger mit einer Gegenprobe über einem Shell-Werkzeug liegt nicht vor
  („keine Messung“); der Architect rechnet für Tabellen- und Gegenprobe-Code zusätzlich mit rund dem Doppelten
  des Produktcodes. Mit der Marke 400 läge die Mitte bei 81 % der Marke; gesetzt ist 500, und geschnitten wird nur, wenn die Schätzung die Marke überschreitet (Praxis des Planners seit 2026-10-16); die Meldeschwelle bleibt.
- `in-progress` → `next`: Der Architect findet, dass die Übernahme mehr als die Funktion `quelle_finden`, den Hilfetext und den
  Aufruf verlangt (etwa eine Änderung am Fragment oder an der Verweis-Mechanik) — zurück zur Zerlegung (§1).
- `in-progress` → `open`: `v0.7.0` weicht von der Fassung ab, die der Planner gelesen hat, oder die Übernahme bricht die Gegenprobe
  eines anderen Gates (etwa `lint`, wenn es auf `tools/harness/*.sh` wirkt, §6) und lässt sich nicht ohne Carveout lösen —
  der Slice wartet auf die Entscheidung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, `make gates` grün, Closure-Notiz mit Lerneintrag geschrieben (Form: neuer Sensor
`make slice-mv-gegenprobe`; die Quellensuche von `slice-mv` hat zugesagt und geprüft zugleich).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Vom Architect vor dem ersten Code-Commit zu klären (offen):**

- *Wirkt `make lint` oder `make docs-check` auf `tools/harness/*.sh`?* Wenn ja, trägt die übernommene Fassung Befunde (Quoting, `local` in
  Schleifen, Kommentarstil nach `AGENTS.md` §3.7); Antwort mit Beleg (Lauf oder Konfiguration), nicht aus dem Gedächtnis.
- *Gehört `slice-mv-gegenprobe` in `GATE_CHECKS` und in die Ausnahmeliste der schnellen Kette?* Der Plan setzt `GATE_CHECKS` an (wie die anderen
  Gegenproben); die Ausnahmeliste nennt „jedes Gate mit `-gegenprobe` im Namen“ (`harness/mk/gates-schnell.mk`), und
  `gates-schnell-gegenprobe` hält sie gegen `GATE_CHECKS` — die Gegenprobe gehört dann dazu. Der Architect bestätigt, und ob die Gegenprobe
  der Gegenprobe der schnellen Kette einen Fall oder eine Änderung braucht.
- *Braucht der Slice `harness/sensors/slice-mv.md`?* Die Zeile in `harness/README.md` reicht, wenn kein Vertrag mit Punkten zu führen ist; der
  Architect entscheidet, und mit Datei zählt der Plan neu (§8).

**Randformen der Gegenprobe** (`AGENTS.md` §3.12; Ort der Entscheidung: der Kopf von `tools/harness/slice-mv-gegenprobe.sh`, weil das Werkzeug
keinen Vertrag unter `harness/vertraege/` führt — der Architect bestätigt den Ort oder legt einen Vertrag an):

- *Glob-Metazeichen im Argument* (`slice-*`, `slice-?`, `[`): Das Argument ist ein Name, kein Muster; `quelle_finden` stellt den Namen in
  Anführungszeichen, ein Treffer über die Metazeichen wäre ein Fehler. Fall und Mutation (Entscheidung offen).
- *Leerzeichen im Argument:* kein Treffer, Exit 2; die Meldung nennt den Namen unverändert. Entscheidung offen.
- *Argument mit `.md`:* `slice-a.md` und `slice-a` meinen dasselbe; `slice-a.md.md` wird nur um eine Endung gekürzt. Fall offen.
- *Leeres Argument und Argument `.md`:* Das leere Argument bricht in `main()` mit der Hilfe ab; `.md` ergäbe den Namen `` und damit
  den Präfix `-`. Entscheidung offen, ob das ein Fall ist.
- *Pfad-Trennzeichen im Argument* (`../x`, `next/slice-a`): `"$1/$d/$name.md"` liest einen Pfad aus dem Argument; ob das erlaubt ist oder abgelehnt
  wird, steht weder in `v0.7.0` noch im Plan. Der Architect entscheidet; ändert die Entscheidung `slice-mv.sh` über `v0.7.0` hinaus,
  ist das eine Abweichung vom Upstream (siehe Risiken) und eine Rückführung nach §4, wenn sie über die Funktion hinausgeht.
- *Aufruf-Voraussetzung:* Die Quellensuche liest den Stand der Verzeichnisse in einem Aufruf; zwei gleichzeitige Aufrufe von `slice-mv` in einem
  Arbeitsbaum sind nicht zugesagt (das Skript verlangt einen sauberen Baum und committet). Der Architect bindet die Zusage daran oder nennt die
  Entscheidung.
- *Aufrufarten des Treibers* (`AGENTS.md` §3.12): Die Gegenprobe ist ein Gate, dessen Exit-Code der Stop-Hook und `make gates` glauben. `make -n`,
  `-q`, `-t`, `-i` und `MAKEFLAGS` aus der Umgebung: je Flag eine Entscheidung wie bei den anderen Gegenproben (Vorbild:
  `kopf-check-gegenprobe`, Punkt 9), ein Fall bei rotem Skript mit Exit ungleich 0 über `make`.
- *Umgebung des Wegwerf-Repos:* `git commit` braucht eine Identität, ein globaler Hook darf nicht greifen (`core.hooksPath`), das Skript wird in
  das Wegwerf-Repo kopiert, weil es `cd "$(dirname "$0")/../.."` ausführt. Entscheidung des Architect: Setzen der Identität in der
  Gegenprobe, nicht im Wegwerf-Repo geerbt.
- *Stelle einer Nachricht in einer Folge, Verhalten je Serverversion:* nicht berührt.

**Risiken:**

- Die übernommene Fassung weicht in einem Kommentar von `v0.7.0` ab (Verweis auf die Gegenprobe statt `test/slice-mv.bats`); ein erneuter Lauf von
  `ai-harness-init` schreibt das Werkzeug kanonisch neu (Kopf des Skripts) und macht die Abweichung rückgängig — **Ausgang:** offen bis Closure (erwartet:
  weiter offen → neuer Register-Eintrag, wenn der Bootstrap die Abweichung überschreibt; der Zeiger auf eine fehlende Datei kehrt dann zurück).
- Die Gegenprobe prüft die Zusagen der Quellensuche, nicht die Verweis-Mechanik; ein Befund in der Mechanik von `v0.7.0` bleibt unbemerkt — **Ausgang:** offen bis Closure
  (erwartet: weiter offen, Abgrenzung in §1).
- Die Gegenprobe (Wegwerf-Repo, `git`, `make`) verlängert den vollen Lauf (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`) — **Ausgang:** offen bis Closure
  (Messung der Dauer in §7).
- Die Schätzung liegt unter der Messung (Faktor 5 ohne gemessenen Vorgänger; `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`) — **Ausgang:** offen bis Closure;
  Marke und Meldeschwelle stehen in §4.
- Der Slice wächst über Marke oder Meldeschwelle von §4 — **Ausgang:** offen bis Closure.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-15). Treffer in der Sub-Area `REPO`: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` wird berührt (neue Gegenprobe im vollen Lauf, §6); `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/meldeschwelle-ohne-artefakt` stehen im Plan als Schätzung mit Faktor und Meldeschwelle in §4, ohne dass der Slice sie schließt. Mit diesem Slice erreicht keine Beobachtung neu eine Schwelle; kein weiterer Treffer.

Nachgezählt (Planner, 2026-10-15, `AGENTS.md` §3.13): Der Slice ist neu und trägt keine Sendung in einen anderen Slice ein; er nimmt keine an. Liefer-Punkte: 2 (Übernahme des Skripts; Gegenprobe mit Gate-Verdrahtung und README-Zeile); Mutationen und Zeile im Gate-Index sind Teile von Punkt 2. Schichten: zwei, das Werkzeug (`tools/harness`) und die Gate-Verdrahtung (`harness/mk`, `harness/README.md`); Planung zählt keine Schicht, kein Produkt-Code, keine Nutzer-Doku. Fällt `harness/sensors/slice-mv.md` an, bleibt es die Gate-Verdrahtung (Dokument desselben Gates). Geschätzt in der Metrik von §4: 280 bis 370 (Mitte 325) gegen Marke 500, 65 % der Marke und unter 70 %; der Plan bleibt ungeteilt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

Nachgezählt beim Verschieben hinter das Release (Planner, 2026-10-16, `AGENTS.md` §3.13): Start und Stelle geändert, Gegenstand und DoD unverändert; der Slice trägt keine Sendung ein und nimmt keine an: Liefer-Punkte weiter 2, Schichten weiter zwei. Kein Slice vor dem Release nennt ihn als Voraussetzung; die Moves der Schritte 39b bis 44 und der Release-Slices gehen von Hand, wo das Präfix mehrdeutig ist.
