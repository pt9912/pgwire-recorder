# Slice slice-harness-stop-hook-positivliste-developer: `docs/developer/` auf die Positivliste des Stop-Hooks

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Auftrag des Nutzers („hook anpassen“,
2026-10-17, ausdrückliches Ja zur Regeländerung). Er hängt an keinem anderen Slice. Die Reihenfolge
gegenüber den übrigen Slices in `next/` legt der Nutzer fest (WIP-Limit 1); eingesammelt wird er von
der nächsten Welle-Closure.

**Bezug:** kein Anspruch des Produkts — der Slice ändert ein Werkzeug des Harness. Bindung an
Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Positivliste; „weitere
Pfade nur, wenn belegt ist, dass kein langsames Gate sie liest“; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt.

**Berührte Spec-Stellen:** `SPEC-052` (Punkte 12, 13 und 15 in `harness/vertraege/SPEC-052-gate-kette.md`:
die Liste, der Beleg, die Verzeichnis-Namen der Link-Regeln); die Sätze „liest nichts unter …“ in
`SPEC-053` (Punkt 13), `SPEC-055` (Punkt 10) und `SPEC-057` (Punkt 10), siehe §3

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Positivliste des Stop-Hooks (`inert` in `.claude/hooks/stop-require-gates.sh`) nimmt
`docs/developer/` auf, mit Vertrag (`SPEC-052`) und Gegenprobe, sodass eine Änderung nur unter diesem
Pfad nach grüner schneller Kette frei gibt, ohne vollen `make gates`.

Grund: `docs/developer/` ist Prosa, die kein Gate der Ausnahmeliste liest; `docs-check` liest sie und
läuft in der schnellen Kette (Beleg in `SPEC-052` Punkt 13, siehe §3). Die Beobachtung
`BEO-REPO/stop-hook-positivliste-ohne-docs-developer` steht bei 1×
(`slice-harness-devcontainer-doku`).

**Entscheidung des Nutzers (2026-10-17, „hook anpassen“, ausdrückliches Ja zur Regeländerung):** Die
Positivliste wird um `docs/developer/` erweitert. Nur dieser Pfad.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Weitere Pfade auf der Positivliste (`docs/maintainer/`, `harness/README.md`, `spec/`,
  `docs/plan/adr/`, `AGENTS.md`) — der Auftrag nennt nur `docs/developer/`; jeder Pfad braucht seinen
  eigenen Beleg nach `SPEC-052` Punkt 13 (`harness/README.md` ist dort ausdrücklich ausgeschlossen, bis
  jemand ihn belegt). `docs/maintainer/releasing.md` liegt heute auf keiner Liste; ob die Pfade aus
  `slice-v1-abschluss-container-doku` §2 („nur `docs/user`, `docs/maintainer`, …“) damit stimmten, ist
  Bestand und wird hier nicht berichtigt — das wäre ein Beschluss des Nutzers über einen weiteren Pfad.
- Eine Änderung am Verfahren des Hooks (Reihenfolge der Prüfungen, Stand der schnellen Kette, Dateiliste,
  Schleifenschutz, Symlink-Regeln) — das Verfahren bleibt; nur die Liste wächst. Wer die Regeln anfasst,
  hat den Plan geändert.
- Produktcode (`cmd/`, `internal/`, `test/`) — Schicht-Abgrenzung: der Slice ändert Harness-Werkzeug und
  seinen Vertrag.
- Der Inhalt von `docs/developer/devcontainer.md` — geliefert von `slice-harness-devcontainer-doku`
  (`done/`), das Dokument bleibt unverändert.
- Eine ADR: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) sieht „weitere Pfade“ bei Beleg vor und legt die Liste in die Spezifikation
  (`AGENTS.md` §3.8). Der Architect bestätigt vor dem Code, dass §3.6 hier keine neue ADR verlangt
  (Empfehlung des Planners: keine; wie bei der Aufnahme von `docs/user/`).
- Die Texte `.claude/commands/implement-slice.md`, `plan-welle.md`, `close-welle.md`,
  `.claude/agents/reviewer.md` und `harness/conventions.md` — sie nennen die Pfade der Liste nicht,
  sondern zeigen auf `inert` im Hook als Quelle (Stand 2026-10-17, per `grep` nach `Positivliste` und
  `docs/user`); der Implementer bestätigt das im Bericht und ändert sie nur, wenn der Befund anders ist.

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Hook und Gegenprobe:** `inert` und `spiegelt` in `.claude/hooks/stop-require-gates.sh` nennen
      `docs/developer/` (Kopf-Kommentar des Hooks nachgezogen); `tools/harness/stop-hook-gegenprobe.sh` hat
      die Fälle und Mutationen aus §3, `make stop-hook-gegenprobe` grün.
- [ ] **Vertrag:** `SPEC-052` Punkte 12, 13 und 15 nennen `docs/developer/` samt Grenzen und Beleg (§3);
      die Frage der Sätze in `SPEC-053`, `SPEC-055`, `SPEC-057` ist vom Architect entschieden und umgesetzt.
- [ ] **Gate-Index:** die Zeile `make stop-hook-gegenprobe` in `harness/README.md` nennt `docs/developer/`.

Pflichten der Closure:

- [ ] `make kopf-check` und `make docs-check` grün; `make gates` einmal am Ende grün (nicht pro Commit;
      Hook, Vertrag und Gegenprobe sind nicht inert, `SPEC-052` Punkt 12).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [ ] Je Zusage des geänderten Vertrags Zusage · Mutation · roter Test im Bericht (`AGENTS.md` §3.10, je
      Element der Aufzählung).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `evidence/slice-harness-stop-hook-positivliste-developer.md`
      in `BEO-REPO/stop-hook-positivliste-ohne-docs-developer/`; **kein Zähler wird gesetzt**.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/hooks/stop-require-gates.sh` | update | `inert`: `docs/developer/*) return 0 ;;` neben `docs/user/*`; `spiegelt`: `docs/developer` in der Zeile der Verzeichnisse der Liste; Kopf-Kommentar (Zeile 16) nennt den Pfad |
| `harness/vertraege/SPEC-052-gate-kette.md` | update | Punkt 12: `docs/developer/` in der Liste, Grenzfälle `docs/developer-alt/…`, `docs/developers/…`, Datei namens `docs/developer`, `docs/Developer/…` in der Aufzählung der nicht inerten Pfade. Punkt 13: Titel und Satz „Alle … sind ohne Ausnahme aufgenommen“ um den Pfad; Beleg für `docs/developer/` (siehe unten); Bedingung der Aufnahme, akzeptiertes Negativ. Punkt 15: die Verzeichnis-Namen, die „nicht inert“ bzw. „spiegeln“ sind (`docs/user`, `docs/user/`) um `docs/developer`, `docs/developer/` |
| `tools/harness/stop-hook-gegenprobe.sh` | update | Fälle, Mutationen, Fixture (Liste unten) |
| `harness/README.md` | update | Zeile `make stop-hook-gegenprobe`: Liste der Positivliste um `docs/developer/` |
| `harness/vertraege/SPEC-053-postgres-versionen.md` (Punkt 13), `SPEC-055-container-und-binaries.md` (Punkt 10), `SPEC-057-integration-aufraeumen.md` (Punkt 10) | update, falls der Architect es verlangt | Jeder nennt, was sein Gate nicht liest („nichts unter `docs/plan/planning/`, `docs/reviews/`, `docs/user/`“); Empfehlung: den Pfad zufügen, ein Satzteil je Datei, damit die Sätze der Liste folgen (`AGENTS.md` §3.9). Sie sind Nehmer der Liste von `SPEC-052` Punkt 13 |

**Beleg für `docs/developer/` (`SPEC-052` Punkt 13), Stand 2026-10-17:**

- Kein Gate der Ausnahmeliste liest den Pfad: `grep -rn "docs/developer"` über `tools/`, `harness/mk/`,
  `Makefile`, `Dockerfile`, `.dockerignore`, `.a-check.yml`, `.golangci.yml` und `harness/` findet keine
  Zeile, die ihn liest (Treffer allein in Plänen und Reports). Die Docker-Gates sehen nur den Build-Kontext;
  die Allowlist von `.dockerignore` (`*`, dann `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`,
  `.golangci.yml`, `tools/harness/lint.sh`) gibt `docs/` nicht frei. `a-check` liest nach `.a-check.yml`
  nur `**/*.go`. Die Gates mit `-gegenprobe` im Namen fahren in Temp-Bäumen und lesen den Pfad des
  Arbeitsbaums nicht (wie für `docs/user/` belegt). `abdeckung-check` liest `docs/user/abdeckung-*.md`,
  `meldungskatalog-check` das Handbuch unter `docs/user/`, `kopf-check` den Planungsbaum, keiner
  `docs/developer/`. Der Implementer wiederholt das `grep` am Stand seines ersten Commits und trägt den
  Befund in `SPEC-052` ein; eine Zeile, die den Pfad liest, hält ihn an und geht an den Architect.
- Gelesen wird der Pfad von `docs-check`: `.d-check.yml` hat `scan.roots: ["."]`, `scan.ignore` nennt
  `**/*.template.md`, `.tmp/**`, `.harness/**` und nicht `docs/developer/`; Links und Anker im Dokument
  prüft das Gate (Beleg: `slice-harness-devcontainer-doku` §3, `make docs-check` über das Dokument
  grün). `docs-check` steht in der schnellen Kette; die Aufnahme gilt, solange es dort steht und der
  Build-Kontext den Pfad nicht freigibt (beides Änderungen an Gates oder `.dockerignore`, nicht inert).
  Eine Prosa-Änderung unter `docs/developer/` mit totem Link färbt `docs-check` rot, und die schnelle
  Kette fängt sie.
- Ob die Schritte des Dokuments gegen den Container laufen, prüft kein Gate; das bleibt Pflicht von
  Implementer und Verifier (wie beim Handbuch, `SPEC-052` Punkt 13).

**Gegenprobe (Namen nach der Nummer des Punktes):**

- Fixture: `DATEIEN_FIXTURE` bekommt `docs/developer/d.md`; im Temp-Repo liegt ein Link
  `docs/developer/ln` → `d.md` (neben `docs/user/ln`).
- Fälle, die freigeben (`approve`): `p12-developer` (`docs/developer/d.md`), `p12-developer-tief`
  (`docs/developer/a/b/c.md`), `p15-link-inert-developer` (Ziel des Links `docs/developer/ln` geändert auf
  `x.md` im selben Verzeichnis).
- Fälle, die sperren (`block`): Grenze `p12-grenze-developer-alt` (`docs/developer-alt/n.md`),
  `p12-grenze-developers` (`docs/developers/n.md`), `p12-grenze-datei-developer` (Datei namens
  `docs/developer`, über `ersetze_durch_datei`), Schreibweise `p12-case-developer`
  (`docs/Developer/n.md`), Link `p15-ziel-aussen-developer` (Link `docs/developer/ln` mit Ziel
  `../../internal/x.go`), Ziel ist das Verzeichnis `p15-ziel-verzeichnis-developer` (Ziel `docs/developer`
  bzw. `docs/developer/`), Spiegelregel `p15-spiegel-developer` (Link außerhalb der Liste mit Ziel
  `docs/developer`, danach Änderung einer inerten Datei: sperrt), Art `X` `p12-art-x-developer`
  (Gitlink `docs/developer/sub`).
- Mutationen (je Element eine, die nur seinen Fall rot färbt): `m-ohne-developer` (Zeile `docs/developer/*`
  aus `inert` genommen → `p12-developer`, `p12-developer-tief`, `p15-link-inert-developer`),
  `m-developer-erste-ebene` (`docs/developer/*/*) ;;` davor → `p12-developer-tief`),
  `aufnahme 'docs/developer-*'` → `p12-grenze-developer-alt`, `aufnahme 'docs/developers/*'` →
  `p12-grenze-developers`, `aufnahme 'docs/developer'` → `p12-grenze-datei-developer`,
  `m-spiegel-ohne-developer` (nur `| docs/developer` aus der Zeile der Verzeichnisse in `spiegelt` →
  `p15-spiegel-developer`, `p15-ziel-verzeichnis-developer`), `nocasematch` (bestehende Mutation, der
  Fall `p12-case-developer` kommt in ihre rote Menge).
- Bestehende Mutationen bleiben gültig (die erwartete Menge der „nur sie“-Prüfung wird geprüft:
  `m-ohne-user` färbt `p12-user-tief`, `p13-docs-user`, `p9-freigabe`, weiterhin nicht die neuen Fälle).
- `Aufzählung zählt je Element` (§3.10): die Liste hat fünf Pfade (`planning`, `reviews`, `user`,
  `developer`, `README.md`); für `developer` gibt es Fall und Mutation wie für die anderen vier. Es gibt
  keine neue Ortsangabe.

## 4. Trigger

**Start** (`next` → `in-progress`): ohne Voraussetzung. Der Anspruch ist ein reiner `git mv` auf dem
Hauptzweig vor der Arbeit (WIP-Limit 1; zum Zeitpunkt der Anlage liegt kein Slice in `in-progress/`).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- .claude/hooks tools/harness harness`,
Tests und Vertrag eingeschlossen. Basis ist der Elternteil des Commits, der diese Datei nach
`in-progress/` bewegt. Schätzung des Planners 120 bis 220 (Produktcode des Hooks 5, Vertrag etwa 40,
Gegenprobe etwa 100 bis 150 mit dem Doppelten des Produktcodes für Tabellen- und Gegenprobe-Code; der
Architect rechnet nach). Meldeschwelle 300, Marke 400.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 400, oder
  die Arbeit verlangt mehr als zwei Schichten (§8), oder ein Gate der Ausnahmeliste liest `docs/developer/`
  (Beleg in §3 fällt) und braucht eine andere Behandlung als die Aufnahme.
- `in-progress` → `open` (blockiert — Carveout?): Ein anderer Slice ändert gleichzeitig Hook oder
  Gegenprobe (Zusammenführung der Listen), oder die schnelle Kette lässt sich für den Arbeitsbaum nicht
  grün fahren.

## 5. Closure-Trigger

DoD vollständig, `make kopf-check`, `make docs-check` und der volle `make gates` (einmal, am Ende; der
Stop-Hook verlangt ihn, weil Hook, Vertrag und Gegenprobe nicht inert sind) grün, Review- und
Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

**Erwarteter Gate-Aufwand.** Ein voller `make gates` am Ende, nicht pro Commit (bei
`slice-harness-devcontainer-doku` endeten zwei volle Läufe an `lint-gegenprobe` und
`test-postgres-versionen-gegenprobe` rot, der dritte grün; ein Wiederholungslauf ist einzuplanen). Davor
genügen `make stop-hook-gegenprobe` und `make gates-schnell`. Danach gibt eine Änderung nur unter `docs/developer/`
nach grüner schneller Kette frei.

## 6. Risiken und offene Punkte

Randformen des Vertrags (`AGENTS.md` §3.12), mit Fundort der Entscheidung; entschieden sind sie vor dem
ersten Code-Commit in `SPEC-052` (der Architect prüft die Liste):

- **Pfadgrenze:** `docs/developer-alt/…`, `docs/developers/…`, eine Datei namens `docs/developer`
  (auch als Verzeichnis-Name ohne Schrägstrich) sind nicht inert; der Eintrag ist der Präfix
  `docs/developer/` mit Schrägstrich. Fundort: `SPEC-052` Punkt 12. Das Verzeichnis `docs/developer`
  selbst steht nicht in einem Dateipfad der Liste.
- **Groß- und Kleinschreibung:** `docs/Developer/…` ist nicht inert, bytegenau und case-sensitiv.
  Fundort: `SPEC-052` Punkt 12.
- **Symlink unter `docs/developer/`:** inert nur, wenn Link-Pfad und Ziel inert sind; ein Ziel außerhalb
  der Liste (`../../internal/x.go`), ein absolutes, ein Ziel über die Wurzel hinaus und ein Ziel, das
  `docs/developer` oder `docs/developer/` selbst benennt, sind nicht inert. Fundort: `SPEC-052` Punkt 15.
  Ein Link außerhalb der Liste mit Ziel in `docs/developer/` oder auf das Verzeichnis spiegelt
  (jede Änderung wird nicht inert). Der Bestand an Links (`git ls-files -s | awk '$1==120000'`:
  `.claude/rules/*`, Ziele außerhalb der Liste) trägt keinen Link auf `docs/developer`; der Implementer
  wiederholt die Zählung am Stand seines Commits und trägt sie in Punkt 15 ein.
- **Zusammenspiel mit Punkt 15 (Wechsel der Art, Umbenennung):** Datei zu Link und zurück unter dem Pfad,
  Umbenennung zwischen `docs/developer/` und einem nicht inerten Pfad (ein gelöschter und ein
  hinzugekommener Pfad, beide müssen inert sein, Punkt 14). Fundort: `SPEC-052` Punkte 14 und 15; die
  bestehenden Fälle für `docs/user/` gelten für den Pfad mit, weil die Regeln pfadunabhängig sind; der Slice
  fügt keinen Fall dafür zu, außer den oben in §3 genannten (Entscheidung des Planners; der Architect
  bestätigt oder erweitert die Liste vor dem Code).
- **Art `X` (Gitlink) und unzulässige Pfade (`.`, `..`, leeres Segment, absolut):** nie inert, auch
  unter `docs/developer/`; Fundort: `SPEC-052` Punkte 12 und 20; der Fall `p12-art-x-developer`
  deckt den Gitlink, die unzulässigen Formen die bestehenden Fälle (die Prüfung steht vor `inert`).
- **Aufrufarten des Treibers:** nicht berührt. Der Slice ändert weder `make gates` noch
  `make gates-schnell`, weder Stand noch Exit-Code; die Aufrufarten `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS`
  entscheiden `SPEC-052` Punkte 5 bis 7 und die bestehenden Fälle.
- **Aufruf-Voraussetzung (Ordnung über mehrere Aufrufe):** nicht berührt; der Vertrag sagt dazu nichts zu.
- **Stelle einer Nachricht:** keine Nachricht im Vertrag.

Risiken:

- Ein künftiges Gate liest `docs/developer/` (zum Beispiel ein Schritt von `u-boot` oder ein
  Docker-Gate mit freigegebenem Build-Kontext) und die Aufnahme gibt zu früh frei — **Ausgang:** weiter
  offen → `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` im Register (Beobachtung steht bei 1×, der
  Slice fügt Beleg zu; akzeptiertes Negativ wie in `SPEC-052` Punkt 13: der Beleg wird beim Review der
  Änderung am Gate neu geführt).
- Die Sätze in `SPEC-053`, `SPEC-055`, `SPEC-057` ziehen die neue Liste nicht nach und widersprechen
  `SPEC-052` Punkt 13 — **Ausgang:** entfallen, sobald der Architect sie vor dem Code entscheidet (§3
  Empfehlung: nachziehen); bleibt die Entscheidung aus, ist das eine Rückgabe an den Architect, kein
  Wert für den Implementer.
- Die Zählung (§8) liegt mit zwei Schichten an der Grenze; eine Sendung an diesen Slice ist ohne
  Nachzählung und ohne Planner keine Adresse (`AGENTS.md` §3.13) — **Ausgang:** entfallen: Der Slice
  nimmt keine Sendung an; `grep` nach der Kennung dieses Slice in anderen Plänen ist leer (Stand der Anlage).

## 7. Closure-Notiz

(leer — wird bei der Closure geschrieben)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Zählung (Planner, 2026-10-17, `AGENTS.md` §3.13).** Liefer-Punkte: drei (Hook samt Gegenprobe;
Vertrag samt der Sätze der Nehmer-Verträge; Gate-Index-Zeile). Schichten: zwei, (1) Harness-Werkzeug
(`.claude/hooks`, `tools/harness`), (2) Vertrag und Harness-Dokumentation (`harness/vertraege/`,
`harness/README.md`; keine Nutzer- oder Wartungs-Doku, kein Produkt-Code). Planung zählt nicht. Der
Slice liegt an beiden Grenzen; weitere Sendungen (Texte der Commands, Wurzel-`README.md`, Handbuch) würden
eine dritte Schicht oder einen vierten Punkt bringen und gehen an den Planner. Dass Hook, Vertrag und
Gegenprobe eine Einheit sind, ist die Lesart des Auftrags; die Zählung ändert sie nicht.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist nur die Default-Sub-Area `*` (Kürzel `REPO`, GF) aus
`harness/conventions.md`; sie erfüllt die Schwelle als Default für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register unter `../observations/` durchgegangen für `REPO`:
Treffer ist `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` (1×, Stand offen, Beleg
`slice-harness-devcontainer-doku`); dieser Slice fügt seinen Beleg zu (2×); der Eintrag erreicht 3× nicht, also bekommt er keinen Ausgang
`geplant`/`verkörpert` zugewiesen, obwohl dieser Slice die Regel schon verkörpert. Daneben einschlägig: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10 je Element, in §3
berücksichtigt) und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Faktor 2 in der Schätzung, §4).
Der Zähler-Stand ist bei der Closure aus den Dateien zu lesen.

Alle berührten Sub-Areas GF.
