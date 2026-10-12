# Slice slice-harness-lint-gegenprobe-last: Die Gegenprobe des Lint-Gates fährt standardmäßig drei Läufe gleichzeitig

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeordnet am 2026-10-12 (Entscheidung des
Nutzers): als Nächstes, direkt hinter `slice-v1-abschluss-zeitangaben-doku` und vor
`slice-v1-abschluss-antwortvergleich` (WIP-Limit 1); Reihenfolge in §4 *Start* und im
Drift-Log der Roadmap.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert nur den Standard eines Werkzeugs des
Harness, keine Zusage des Produkts. Bindung an Entscheidungen:
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Lint-Gate mit Gegenprobe;
die ADR nennt die Zahl der gleichzeitigen Läufe nicht).

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make lint-gegenprobe` fährt ohne gesetztes `LINT_GEGENPROBE_PARALLEL` höchstens
drei Läufe der Stufe `lint` gleichzeitig (bisher sechs): `tools/harness/lint-gegenprobe.sh`
setzt den Standard auf 3 und nennt ihn in seinem Kopf, ohne die Zusage, ein Lauf von
`make gates` prüfe ihn, und ohne ihn als Obergrenze der Last von `make gates` zu lesen
(Entscheidung des Nutzers vom 2026-10-12, §6).

**Anlass** (Entscheidung des Nutzers vom 2026-10-12, den Standard auf 3 zu senken): Last auf
dem Rechner (20 Kerne) und ein OOM-Kill am 2026-10-12 um 06:05. Der Kernel-Log nennt als
getroffenen Prozess ein `bash` mit 17 GB RSS.

**Messung** (Planner-Auftrag, 2026-10-12, `make lint-gegenprobe` allein, ohne Änderung am
Repo; RAM-Basis vor dem Lauf rund 14,5 GB belegt, Swap 6,6 GB):

| Parallelität | Laufzeit | Load (Mittel / Maximum) | RAM belegt (Maximum) |
|---|---|---|---|
| 6 (bisheriger Standard) | 150 s | 57 / 75 | nicht erfasst |
| 3 | 167 s | 27,5 / 31,9 | 18,7 GB |
| 2 | 188 s | 23,3 / 27,3 | 17,9 GB |

**Die Ursache des OOM-Kills ist nicht belegt.** In der Messung wurde kein Einzelprozess in der
Größe von 17 GB gesehen, und für Parallelität 6 liegt keine RAM-Zahl derselben Messreihe vor.
Der Slice senkt die Last der Gegenprobe; dass er den OOM behebt, sagt er nicht zu (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Fall im Skript, der den Standard prüft — Bestand bleibt bewusst stehen: Entscheidung des
  Nutzers vom 2026-10-12 (Weg (2) der Randform *Standard ohne Fall*, §6); der Kopf nennt die
  Zahl und sagt nicht zu, dass ein Lauf sie prüft (`AGENTS.md` §3.11).
- Gegenproben nur bei geänderten Eingaben fahren oder `make gates` aufteilen — Lockerung
  beziehungsweise Umbau eines Gates, nach `AGENTS.md` §3.6 eine ADR und eine Entscheidung des
  Architect; es ist ein anderer Vorgang als eine Zahl im Werkzeug. Es gibt dafür noch keinen
  Folge-Slice und keine Kennung (Punkt an Architect und Nutzer, siehe Bericht des Planners).
- Eine Änderung des Stop-Hooks (`.claude/hooks/stop-require-gates.sh`) — Schicht-Abgrenzung:
  Der Slice ändert ein Werkzeug des Gates, nicht den Nachweis, der `make gates` verlangt; auch
  dafür gibt es noch keinen Folge-Slice und keine Kennung.
- Die Parallelität anderer Gegenproben — Bestand bleibt bewusst stehen: Außer
  `tools/harness/lint-gegenprobe.sh` führt kein Skript unter `tools/` und kein Fragment unter
  `harness/mk/` eine eigene Obergrenze (Suche nach `wait -n`, `xargs -P`, `-j`, 2026-10-12);
  die Parallelität von `make gates` selbst (die Checks laufen unter `make -j`, Kopf des
  `Makefile`) gehört zum ersten Punkt dieser Liste.
- Eine Zusage, der Slice behebe den OOM-Kill, in Ziel, DoD, Kopf oder Kommentar —
  `AGENTS.md` §3.11: Die Ursache ist nicht belegt, und kein Test prüft eine solche Zusage.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Standard: Ohne gesetztes `LINT_GEGENPROBE_PARALLEL` ist die Höchstzahl gleichzeitiger
      Läufe in `tools/harness/lint-gegenprobe.sh` 3 (Zeile `parallel=`); der Kopf des Skripts
      nennt 3 als Standard und sagt weder zu, dass ein Lauf von `make gates` ihn prüft, noch
      etwas über die Wirkung auf Last, Speicher oder einen OOM-Kill zu (`AGENTS.md` §3.11).
      Der Satz „den Default fährt jeder Lauf von `make gates`“ (Kopf, Zeilen 18–19) ist
      gestrichen, die 6 im Kopf (Zeile 14) ersetzt, und der Kopf liest 3 als Höchstzahl je
      Lauf des Skripts, nicht als Obergrenze der Last von `make gates` (§6).
      Kein Fall im Skript prüft den Standard (Entscheidung des Nutzers, §6). Ein gültiger gesetzter Wert gilt unverändert.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-harness-lint-gegenprobe-last.md` in
      `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` mit der Messung dieses Slice
      (`make lint-gegenprobe` allein, Laufzeit, Load im Mittel und Maximum, RAM belegt im
      Maximum und Basis vor dem Lauf, Parallelität 3 als Standard, und `make gates` als
      Ganzes, falls der Auftraggeber es fährt); der Zähler folgt aus den Dateien (2×). Die
      Messung ist Beleg, kein Liefer-Punkt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/lint-gegenprobe.sh` | update | `parallel=6` (Zeile 126) wird `parallel=3`; der Kopf (Zeile 13 ff.) nennt 3 als Standard (Zeile 14) und streicht den Satz „den Default fährt jeder Lauf von `make gates`“ (Zeilen 18–19), ohne die Zusage, ein Lauf prüfe ihn, und ohne 3 als Obergrenze der Last von `make gates` zu lesen; kein neuer Fall (§6) |

- Die Zahl 6 steht nach Suche vom 2026-10-12 (`grep -rn` nach `LINT_GEGENPROBE` und
  `PARALLEL` über `tools/`, `harness/`, `spec/`, `docs/`, `README.md`, `AGENTS.md`, `.claude/`)
  im Skript (Kopf und Zeile 126) und als „bis zu sechs gleichzeitig“ im Beleg
  `evidence/slice-harness-lint.md` unter `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`
  (Befund von damals, bleibt). Die Spezifikation `SPEC-049`, [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) und
  `harness/sensors/lint.md` nennen weder die Zahl noch die Einstellung; deshalb steht im Kopf
  des Plans keine Spec-Stelle, und es gibt keine Doku nachzuziehen. Der Implementer wiederholt
  die Suche vor dem Commit; findet er eine weitere Fundstelle in einem lebenden Dokument, zieht
  er sie im selben Commit nach.
- Kein Produkt-Code, kein Makefile-Fragment, keine Tests unter `test/`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-doku` liegt in `done/`
(WIP-Limit 1); dieser Slice steht direkt dahinter und vor `slice-v1-abschluss-antwortvergleich`
(Schritt 24) und den übrigen offenen Slices von welle-v1-abschluss, danach folgen sie in der
Reihenfolge des §5 der Welle, dann `slice-harness-integration-volumes` und die übrigen
(Drift-Log der Roadmap vom 2026-10-12). Technisch hängt er an nichts. Die Reihenfolge ist
Entscheidung des Nutzers vom 2026-10-12. Der Architect hat die Randformen aus §6 am 2026-10-12
geprüft (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`) und alle als entschieden
bestätigt; der Plan ist vor dem ersten Code-Commit bereit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Nutzer nimmt die Entscheidung
  zu *Standard ohne Fall* (§6) zurück, der Fall verlangt eine Änderung der
  Gegenprobe-Struktur, und der Diff ist nicht in einer Review-Sitzung prüfbar.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer verlangt statt der Senkung eine
  Aufteilung oder bedingte Läufe von `make gates` (erster Punkt von *Ausdrücklich NICHT*); dann
  entscheidet der Architect mit einer ADR, bevor der Slice weitergeht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Messung als Evidence im
Register, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — der Slice ändert den Standard einer vorhandenen Einstellung
und liefert keinen neuen Vertrag; die Randformen stehen trotzdem hier, weil die Einstellung
einen Vertrag im Kopf des Skripts hat. **Vom Architect am 2026-10-12 geprüft, alle entschieden.**
Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

- *Ungültiger Wert* (`0`, leer, `-1`, `abc`, ` 3`) — **entschieden:** im Kopf des Skripts
  (Zeile 13 ff.: `^[1-9][0-9]*$`, Ausgang 2 mit fester Zeile auf stderr), geprüft durch die
  Fälle `einstellung-<wert>`; unverändert.
- *Gültiger gesetzter Wert* (etwa `6`) — **entschieden:** gilt unverändert und überstimmt den
  Standard; wer die frühere Parallelität will, setzt sie. Kein Fall ändert sich. Ein Aufrufer,
  der die Einstellung setzt, ist im Repo nicht zu finden (Suche in §3). Das gilt auch für eine in
  der Umgebung des Aufrufers gesetzte gültige `LINT_GEGENPROBE_PARALLEL` (Shell-Export): sie
  überstimmt den neuen Standard still, denn `harness/mk/lint.mk` reicht die Umgebung unverändert
  durch. Entschieden (Architect, 2026-10-12): gilt unverändert, der Nutzer prüft seine
  Umgebung; kein Code.
- *Nicht gesetzt (Standard)* — **entschieden:** Entscheidung des Nutzers vom 2026-10-12, Weg (2):
  Kein Fall im Skript prüft den Standard. Der Kopf des Skripts nennt die Zahl 3 ohne die
  Zusage, ein Lauf von `make gates` prüfe sie (`AGENTS.md` §3.11). Der Implementer streicht im
  Kopf den Satz „den Default fährt jeder Lauf von `make gates`“ (Zeilen 18–19, am Ende des
  Absatzes zu `einstellung-<wert>`) und ersetzt die 6 im Kopf (Zeile 14) durch 3; es genügt
  nicht, nur die Zahl zu tauschen. **Akzeptiertes Negativ:** Wer `parallel=3` auf einen anderen
  Wert ändert, bricht kein Gate. Grund: Der Standard ist eine Höchstzahl mit Einstellung zum
  Überstimmen; Last oder Laufzeit sind keine Zusage des Werkzeugs. Die Entscheidung des Nutzers
  (Weg 2, 2026-10-12) steht in diesem Absatz und in der Closure-Notiz; die Spezifikation trägt
  sie nicht (nächster Punkt).
- *Weniger als drei Fälle oder Kerne* — **entschieden:** die Obergrenze ist eine Höchstzahl
  (`jobs -rp | wc -l` gegen `parallel`); weniger Fälle laufen alle, mehr warten. Unverändert,
  kein Fall nötig.
- *Ort der Zahl und Fundstelle der Entscheidung zu Standard ohne Fall* — **entschieden**
  (Architect, 2026-10-12): keine Zeile in `SPEC-049`. Der Standard ist die Höchstzahl eines
  Werkzeug-Parameters und keine Randform des Lint-Vertrags; ein Spec-Satz „höchstens drei
  Läufe“ wäre selbst eine Zusage ohne Prüfung (`AGENTS.md` §3.11). Fundstelle der Entscheidung
  ist dieser Abschnitt §6. Die Spezifikation ändert sich nicht, der Kopf des Plans bleibt `—`;
  es gibt keine Folge-ADR und keine Lockerung nach `AGENTS.md` §3.6.
- *Lesart des Standards 3* — **entschieden:** 3 ist eine Höchstzahl je Lauf des Skripts und
  keine Obergrenze der Last von `make gates`: Weitere Checks laufen unter `make` zusätzlich, und
  die Messung (§1) gilt für die Gegenprobe allein. Anforderung an den Wortlaut: Der Kopf des
  Skripts darf den Standard nicht anders lesen lassen, nennt also keine Last-, Speicher- oder
  Zeitgrenze für `make gates`.
- *Wirkung auf Last und Speicher* — **keine Randform des Skripts, sondern Messung:** siehe die
  Risiken; das Skript sagt über sie nichts zu.

**Mutationen** (`AGENTS.md` §3.10): Der Slice sagt nicht zu, dass der Standard 3 *geprüft*
wird; die Mutation `parallel=3` wird `parallel=6` bleibt deshalb grün (Weg (2), akzeptiertes
Negativ, kein roter Fall). Die Zusage „ein gesetzter Wert gilt unverändert“ ist
nicht neu; ihre Fälle sind die bestehenden `einstellung-<wert>` für die ungültigen Werte, für
gültige gibt es keinen.

**Risiken:**

- Die Senkung behebt den OOM-Kill möglicherweise nicht — **Ausgang:** weiter offen:
  → `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` (neu bei der Closure angelegt, falls
  bis dahin kein anderer Slice die Ursache untersucht). Belegt ist: ein `bash` mit 17 GB RSS
  wurde am 2026-10-12 um 06:05 getötet (Kernel-Log). Nicht belegt ist, dass dieser Prozess zur
  Gegenprobe gehörte, dass seine Größe mit der Zahl gleichzeitiger Läufe wächst, und wie sich
  der RAM bei Parallelität 6 verhielt (nicht mit derselben Messreihe erfasst). Gemessen wurde
  die Gegenprobe allein, nicht `make gates` als Ganzes, in dem weitere Checks gleichzeitig
  laufen; bei Parallelität 3 und 2 war der RAM im Maximum 18,7 GB und 17,9 GB bei einer Basis
  von rund 14,5 GB (Swap 6,6 GB), also kein Einzelprozess in der Größe des Kills. Zeigt der
  nächste Lauf von `make gates` wieder einen OOM-Kill, ist das Risiko eingetreten und braucht
  eine Untersuchung (Folge-Slice mit Kennung, vom Planner anzulegen); der Slice nimmt sie nicht
  auf.
- Die Gegenprobe läuft länger (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, 1×): 150 s bei
  6 gegen 167 s bei 3 (rund 11 % mehr, allein gemessen) — **Ausgang:** weiter offen:
  → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, Beleg
  `evidence/slice-harness-lint-gegenprobe-last.md` (Zähler 2×) mit den Zahlen dieses Slice;
  keine Lockerung von Fällen, die Laufzeit steht mit Zahl im Beleg.
- Die Messung gilt für einen Rechner mit 20 Kernen und dem Stand vom 2026-10-12; auf anderen
  Rechnern oder unter CI können die Zahlen anders ausfallen — **Ausgang:** entfallen: Der
  Standard ist eine Höchstzahl mit einer Einstellung zum Überstimmen; der Slice sagt keine
  Laufzeit oder Last für andere Rechner zu, und die Zahlen stehen nur im Beleg.

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

### Belege des Implementers

Änderung: Commit `10dc4c8` (`tools/harness/lint-gegenprobe.sh`: `parallel=3`, Kopf nennt 3 als
Höchstzahl je Lauf des Skripts, keine Obergrenze der Last von `make gates`; der Satz „den Default
fährt jeder Lauf von `make gates`“ ist gestrichen). Die Suche aus §3 (`LINT_GEGENPROBE`,
`PARALLEL` über `tools/`, `harness/`, `spec/`, `docs/user`, `README.md`, `AGENTS.md`, `.claude/`,
`Makefile`) fand danach nur das Skript selbst; keine weitere Fundstelle.

**Messung nachher** (`make lint-gegenprobe` allein, Standard 3, Sampler alle 2 s mit
`/proc/loadavg` (1-Minuten-Wert) und `/proc/meminfo` (MemTotal minus MemAvailable); Ausgang 0,
`gruen`):

| Lauf | Laufzeit | Load (Mittel / Maximum) | RAM belegt (Maximum) |
|---|---|---|---|
| 1 | 166 s | 27,8 / 38,0 | nicht erfasst (Sampler las `free` falsch) |
| 2 | 167 s | 31,0 / 37,7 | 19,3 GB (Basis vor dem Lauf rund 15,6 GB; Load-Start 22,2) |

Vorwerte aus §1: Parallelität 6: 150 s, Load 57 / 75; 3: 167 s, 27,5 / 31,9, 18,7 GB;
2: 188 s, 23,3 / 27,3, 17,9 GB. Der Lauf mit dem neuen Standard deckt sich mit der Messung des
Planners für Parallelität 3 (Laufzeit gleich; Load und RAM im Maximum etwas höher, der Host war
bei Start schon mit Load 22 belegt). Eine Messung der Parallelität 6 in diesem Lauf fehlt.

**Mutation:** `parallel=3` wird `parallel=6` bleibt grün — akzeptiertes Negativ (§6, Weg 2),
nicht gefahren, da kein Fall den Standard prüft und der Kopf es nicht zusagt. Die Fälle
`einstellung-<wert>` (ungültige Werte) liefen im Lauf mit und sind grün.

**Nicht belegt:** dass die Senkung den OOM-Kill behebt; die Ursache ist ungeklärt (Risiko 1 in §6).
Gemessen ist die Gegenprobe allein, nicht `make gates` als Ganzes.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/`
gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand `f60cb55`, 2026-10-12). Treffer:

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (1×, offen) — berührt: Der Slice ändert die
  Last und die Laufzeit der größten Gegenprobe in `make gates` (Risiko in §6). Er liefert bei
  seiner Closure die Messung als Beleg `evidence/slice-harness-lint-gegenprobe-last.md`; der
  Zähler steht dann bei 2×, unter der Schwelle 3×.
- `BEO-REPO/test-faellt-unter-cpu-last-aus` (1×, offen) — verwandt über das Thema Last auf dem
  Rechner, aber nicht der Gegenstand: Der Slice ändert kein Testverhalten, und er sagt nicht
  zu, dass die Senkung Ausfälle unter Last verhindert; kein Beleg wird eingetragen.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos;
  der Auslöser steht in §4 *Start*.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — der Slice ändert keine
  Konfiguration eines Gates; die Randform *Standard ohne Fall* (§6) ist dieselbe Frage in
  klein (eine Zahl im Kopf liest sich als geprüft).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — kein
  neuer Vertrag, aber der Kopf des Skripts ist eine Zusage: Mutation und Randformen in §6,
  keine Wirkungszusage im Kopf.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 1 (der Standard mit Kopf; kein Fall für den
Standard, Entscheidung des Nutzers vom 2026-10-12, §6). Die Messung als Beleg zählt nicht (Register). Schichten: 1
(Harness-Werkzeug `tools/harness`); keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein
Bootstrap. Der Slice trägt keine Sendung in einen anderen Slice ein (kein Folge-Slice
genannt, keine Bindung), also gibt es keinen Nehmer nachzuzählen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
