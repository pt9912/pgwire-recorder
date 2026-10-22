# Slice slice-harness-dev-kopie-und-pin: Die Skripte liegen zweimal vor — Quelle unter `tools/dev/`, Nutzer-Kopie gepinnt mit SHA256 (Entwurf, Fragen offen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-22). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist. Angelegt als Entwurf und als Adresse: `slice-harness-gegenproben-umzug-dev` (Schritt 1) schließt das Kopieren der Skripte nach `tools/dev/` und das Pinnen aus (§1 dort) und braucht einen Slice, der sie annimmt (`AGENTS.md` §3.13). Der Plan liegt in `open/`: Er ist nicht priorisiert, und die Fragen in §1 sind nicht entschieden.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert, wo Werkzeuge des Harness entwickelt und wo sie gelaufen werden, keinen Anspruch des Produkts. Entscheidungen des Nutzers vom 2026-10-22 (§1, wörtlich).

**Berührte Spec-Stellen:** — (offen). Ob ein Vertrag für die Freigabe und das Gate entsteht, entscheidet sich mit den Fragen in §1; ihn schreibt der Architect nach der Entscheidung des Nutzers, in `harness/vertraege/`.

**Verantwortlich:** — (bis zur Priorisierung)

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel (Entwurf):** Die Skripte des Harness liegen zweimal vor: die Quelle zum Entwickeln unter `tools/dev/`, die Kopie zum Nutzen im bisherigen Ordner (`tools/harness/`, `tools/test/`), gepinnt mit SHA256; ein Gate prüft die Nutzer-Kopie gegen die Summe, ein Weg bringt eine geprüfte dev-Fassung in die Nutzer-Kopie, und die Gegenproben prüfen die dev-Fassung. Wie, ist nicht entschieden; die offenen Fragen stehen unten, und dieser Plan ist bis zu ihrer Entscheidung ein Merkzettel und eine Adresse, kein lieferbarer Plan.

**Sendung aus `slice-harness-gegenproben-umzug-dev` (Geber, Schritt 1; `AGENTS.md` §3.13):** Das Kopieren der Skripte nach `tools/dev/` als Quelle und das Pinnen mit SHA (Schritt 2 nach der Reihenfolge des Nutzers) sowie die Umstellung dessen, was die Gegenproben als Skript unter Test lesen: Der Umzug lässt jede Gegenprobe ihr Skript weiter an seinem bisherigen Ort prüfen (rund 270 Zeilen in den 21 Gegenproben nennen `tools/…`, Zählung des Planners am Stand `1ad628a7`; die meisten betreffen das Skript unter Test). Dieser Slice nimmt die Sendung an; sein §1 schließt sie nicht aus.

**Entscheidungen des Nutzers vom 2026-10-22** (wörtlich, so vom Koordinator übermittelt):

1. „erst tools/dev/... anlegen und dann die script kopieren/gegeproben umziehen, pinnen (sha, ...)“ — die Reihenfolge; das Anlegen und der Umzug der Gegenproben sind `slice-harness-gegenproben-umzug-dev`.
2. „Aber für das Pinnen braucht man kein eigenes Repo“
3. „die scripte halten wir jetzt doppelt, einmal zum Nutzen im aktuellen Ordner und den dev-Ordner zum Entwickeln“
4. „auch für die anderen Punkte, wie agents, hooks, mk könnten wir mit "dev"-Ordner lösen“ — die **Richtung**; der Nutzer hat den Umfang (Frage f) und den Mechanismus (Frage g) nicht bestätigt.
5. „ok, dein Vorschlag“ auf folgenden Vorschlag: **Die Nutzer-Seite sind Dateien, gepinnt mit SHA256 (Variante A).** Ein lokal gebautes Image (Variante B: Digest des Basis-Images gepinnt, die Skripte per Summe gepinnt, die Summe der Dateien bleibt der Pin) nur dort, wo ein Skript Werkzeuge braucht, die ein Host nicht hat; ein veröffentlichtes Image mit Registry-Digest wie bei `d-check` und `a-check` (Variante C) **nicht**, bis ein zweites Projekt die Werkzeuge braucht. Hooks bleiben Dateien (ein `docker run` pro Werkzeug-Aufruf wäre zu langsam). Damit ist die Frage (h) entschieden.

**Begründung der Richtung (4)** (Vorschlag des Koordinators, nicht Entscheidung des Nutzers): Die Quelle eines Hooks oder einer Rolle darf sich nicht im laufenden Pfad ändern, weil ein fehlerhafter Hook (der Stop-Hook, der geplante Zeitwächter) jede Sitzung sperrt. **Am Repo geprüft (Planner, 2026-10-22, Stand `1ad628a7`):** `.claude/settings.json` startet Hooks als `bash "$CLAUDE_PROJECT_DIR"/.claude/hooks/<name>.sh` (`span-emit.sh`, `pretooluse-command-guard.sh`, `stop-require-gates.sh`); Rollen und Kommandos liegen in `.claude/agents/` und `.claude/commands/`, Regeln in `.claude/rules/` (Symlinks auf `AGENTS.md`, `harness/conventions.md`, `harness/README.md` und zwei Module der Baseline); `git grep -n 'tools/dev'` nennt am Stand `1ad628a7` nichts. Eine Datei unter `tools/dev/` lädt demnach nichts im Repo; eine dev-Kopie eines Hooks wird nie ausgeführt, solange keine Zeile in `.claude/settings.json` auf sie zeigt. Dass Claude Code selbst nur diese Pfade liest, belegt das Repo nicht; es belegt nur, dass es keinen Verweis auf `tools/dev/` gibt.

**Offene Fragen** (Fragen, keine Entscheidungen; die Empfehlung ist die des Planners, wo nichts anderes steht):

- **(a) Wer kopiert von dev nach Nutzen?** Empfehlung: ein Make-Ziel (Name wählt der Architect), das die Gegenprobe des Skripts fährt, bei grün die Nutzer-Kopie überschreibt und die Summe schreibt; kein Kopieren von Hand. Unterfrage: Für die zwei Gegenproben, die nur mit engster Auswahl laufen (`release-gegenprobe`, `test-container-gegenprobe`, bis `slice-harness-gegenproben-schneller` in `done/` liegt), verlangt die Freigabe den engsten Fall oder den vollen Lauf?
- **(b) Wo liegt die Summe?** Empfehlung: `SHA256SUMS` neben den Nutzer-Kopien, wie bei `.harness/baseline/<tag>/SHA256SUMS`, je Verzeichnis der Nutzer-Kopien.
- **(c) Welches Gate prüft Nutzer-Kopie gegen Summe?** Empfehlung: ein eigenes Skript mit Fragment, bash und netzlos, in `GATE_CHECKS` und in der schnellen Kette; nicht als Erweiterung von `tools/harness/baseline-verify.sh`, dessen Gegenstand die vendored Baseline ist.
- **(d) Müssen dev-Kopie und Nutzer-Kopie nach dem Freigeben gleich sein?** Empfehlung: gleich zum Zeitpunkt der Freigabe (die Gleichheit ist die Freigabe); danach darf dev vorauslaufen, und das Gate prüft allein, dass die Nutzer-Kopie ihrer Summe entspricht, nicht Gleichheit mit dev.
- **(e) Was ändert sich an den rund 270 Pfaden der Gegenproben auf das Skript unter Test?** Empfehlung: Jede Gegenprobe führt am Kopf eine Variable für den Ort des Skripts unter Test (neue Variablen mit dem Präfix `gp_`, `BEO-REPO/gegenprobe-variable-kollidiert-mit-skript`); die Temp-Bäume der Gegenproben bilden weiter das Layout der Nutzer-Seite nach, nur die Quelle der Kopie in den Baum wechselt auf dev.
- **(f) Umfang: Was wird doppelt gehalten?** Empfehlung des Koordinators (nicht Entscheidung des Nutzers): ja für die laufenden Skripte unter `tools/` (ohne Gegenproben), die Fragmente `harness/mk/` (ohne die der Gegenproben, die selbst Entwicklung sind und ganz nach `tools/dev/` ziehen könnten), `.claude/hooks/`, `.claude/agents/`, `.claude/commands/` und `.githooks/commit-msg`; **nicht** `AGENTS.md`, `harness/vertraege/`, `spec/`, ADRs und Pläne (Dokumente mit eigenem Entscheidungsweg und Rang in der Source Precedence). Empfehlung des Planners zur Reihenfolge: mit den Skripten unter `tools/` beginnen, dann die Hooks (dort ist der Nutzen der Trennung am größten), dann Rollen und Kommandos, zuletzt die Fragmente; je Bereich ein Slice (§8).
- **(g) Mechanismus:** Empfehlung des Koordinators (nicht Entscheidung des Nutzers): ein Verzeichnis `tools/dev/` mit Spiegelstruktur und einer Manifest-Datei (je Zeile dev-Pfad und Nutzer-Pfad) als einzigem Mechanismus für alle Bereiche, `SHA256SUMS` neben den Nutzer-Kopien, ein Gate, das beides prüft. Empfehlung des Planners: ja zum Manifest, weil `.claude/` und `harness/mk/` sich unter `tools/dev/` nicht spiegeln lassen, ohne Pfade zu erfinden; die Zuordnung steht dann an einer Stelle.
- **(i) Braucht die Trennung eine ADR?** Empfehlung: ja, vom Architect vor dem Vertrag: Ort der Quelle, Pin, Freigabe, und dass die Nutzer-Seite Dateien sind (Entscheidung 5); die Regeldetails stehen im Vertrag (`AGENTS.md` §3.8). Grund: Sie führt ein neues Gate ein und berührt die Rollen-Anweisungen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Umzug der Gegenproben nach `tools/dev/`** — Klasse 1, Geber `slice-harness-gegenproben-umzug-dev` (liegt vor diesem Slice). Grund: Er ist Schritt 1 der Reihenfolge des Nutzers (Entscheidung 1) und mechanisch; dieser Slice setzt ihn voraus.
- **Ein veröffentlichtes Image mit Registry-Digest (Variante C) für die Werkzeuge** — Klasse 3, entschieden (Entscheidung 5). Grund: Es lohnt erst, wenn ein zweites Projekt die Werkzeuge braucht.
- **Ein eigenes Repo für die Werkzeuge** — entschieden (Entscheidung 2): Das Pinnen braucht keines.
- **Dokumente mit eigenem Entscheidungsweg** (`AGENTS.md`, `harness/vertraege/`, `spec/`, ADRs, Pläne) — Klasse 2 nach der Empfehlung (f), nicht entschieden. Grund: Ihr Rang in der Source Precedence verlangt den eigenen Weg; eine Doppelhaltung ließe zwei Fassungen desselben Dokuments gelten.
- **Produkt-Code** (`cmd/`, `internal/`, `test/`) — Klasse 4. Grund: Die Trennung betrifft Werkzeuge des Harness, nicht das Produkt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Entwurf: Die Liefer-Punkte stehen erst fest, wenn der Nutzer (a) bis (g) und (i) aus §1 entschieden hat; der Architect füllt §3 danach, der Planner schneidet nach dem Umfang (f).

- [ ] **Sendung aus `slice-harness-gegenproben-umzug-dev` geliefert (Adresse, `AGENTS.md` §3.13):** Die Skripte liegen als Quelle unter `tools/dev/` und als Nutzer-Kopie im bisherigen Ordner; die Nutzer-Kopie ist mit SHA256 gepinnt; ein Gate in `GATE_CHECKS` prüft Kopie gegen Summe; ein Weg bringt eine dev-Fassung in die Nutzer-Kopie; die Gegenproben prüfen die dev-Fassung. Der Inhalt jedes Halbsatzes hängt an den Fragen (a) bis (g) und (i) und ist bis dahin nicht prüfbar.
- [ ] Die Fragen (a) bis (g) und (i) aus §1 sind vom Nutzer entschieden, die Entscheidung steht in einem Vertrag unter `harness/vertraege/` (Architect), und eine ADR liegt vor, wenn (i) bejaht ist.
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

Wird vom Architect nach der Entscheidung des Nutzers gefüllt (Fragen in §1). Pfad-Kandidaten nach der Empfehlung (f): `tools/dev/` (neu), `tools/harness/` und `tools/test/` (Nutzer-Kopien, Summe), ein Skript und Fragment für das Gate, `harness/mk/`, die 21 Gegenproben unter `tools/dev/` (Ort des Skripts unter Test), `harness/vertraege/` (Architect).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next`, noch nicht eingetreten): `slice-harness-gegenproben-umzug-dev` liegt in `done/`, der Nutzer hat (a) bis (g) und (i) aus §1 entschieden und priorisiert den Slice, der Architect hat Vertrag und ADR geschrieben, und der Planner hat nach dem Umfang (f) geschnitten. **Schätzung, Marke und Meldeschwelle** (`AGENTS.md` §3.14) trägt der Planner beim Schnitt ein; solange der Plan in `open/` liegt und die Fragen offen sind, gibt es keine.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:** werden beim Schnitt benannt (Bedingung je Teil-Slice). Vorab gilt: Kommt der Umfang „alle Bereiche aus (f)“, ist der Slice über drei Liefer-Punkten und wird nach Bereich geschnitten (§8).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen:** werden mit dem Vertrag des Architect genannt, nicht vorher (`AGENTS.md` §3.12); sicher sind die Aufrufarten von `make` für das Freigabe-Ziel, die Plattform des Pins (`sha256sum` unter Linux, macOS und Windows) und ein Abbruch mitten im Kopieren (halb geschriebene Nutzer-Kopie mit alter Summe).

- Der Plan wird gestartet, bevor der Nutzer (a) bis (g) und (i) entschieden hat, und der Implementer entscheidet sie im Code — **Ausgang:** offen bis Closure (entfallen: Die Entscheidungen stehen im Vertrag vor dem ersten Code-Commit; eingetreten: Rückgabe an den Architect nach `AGENTS.md` §3.12).
- Die Doppelhaltung driftet: Die Nutzer-Kopie wird im Nutzer-Pfad bearbeitet statt in dev, oder dev läuft der Summe davon, ohne dass es jemand merkt — **Ausgang:** offen bis Closure (Gegenmittel: das Gate aus (c); entfallen, wenn es eine geänderte Nutzer-Kopie ohne neue Summe rot meldet und eine Gegenprobe es zeigt).
- Der Umfang (f) wächst über drei Liefer-Punkte (`BEO-REPO/slice-waechst-durch-uebernahmen`) — **Ausgang:** offen bis Closure (entfallen: Der Slice ist nach Bereich geschnitten, §8).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`, Greenfield); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist für eine Sub-Area, die das ganze Repo deckt, nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (Stand `1ad628a7`). Treffer: `folge-slice-adresse-nimmt-nicht-an` (13×; dieser Plan ist die Adresse für die Sendung aus `slice-harness-gegenproben-umzug-dev` und nimmt sie in §1 und DoD an), `slice-waechst-durch-uebernahmen` (7×; Zählung unten), `gegenprobe-variable-kollidiert-mit-skript` (1×; Empfehlung (e): Präfix `gp_`), `gegenprobe-nicht-in-kette-regression-unbemerkt` (5×; ein Pin ohne Gate wäre dieselbe Lücke, deshalb (c)). Erreicht keiner mit diesem Slice 3×.

**Zählung** (Planner, 2026-10-22, `AGENTS.md` §3.13; **vorläufig**, der Umfang (f) ist offen): Mit dem Umfang „Skripte unter `tools/`“ (Empfehlung des Planners, erster Teil) liegen drei Liefer-Punkte nahe — (1) Quelle unter `tools/dev/` und Manifest, (2) Pin und Gate, (3) die Gegenproben prüfen die dev-Fassung — und zwei Schichten (Harness-Werkzeug; Fragmente und Gegenproben zählen zur selben). Mit „alle Bereiche aus (f)“ liegt der Slice über drei Punkten und über zwei Schichten (`.claude/` ist eine zweite Gruppe von Harness-Artefakten). **Vorgemerkter Schnitt:** nach Bereich — Skripte unter `tools/`, Hooks, Rollen und Kommandos, Fragmente von `harness/mk/` —, je Bereich ein Slice mit demselben Mechanismus. Die Sendung aus dem Umzug (die Umstellung des Skripts unter Test in den Gegenproben, rund 270 Zeilen) liegt im dritten Punkt und ändert die Zählung nicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Verträge und Konventionen des Harness gehen den Skripten voraus).
