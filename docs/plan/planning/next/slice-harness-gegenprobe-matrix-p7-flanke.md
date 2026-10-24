# Slice slice-harness-gegenprobe-matrix-p7-flanke: Eine Mutation der Gruppe `p7-sig` färbt nur Fälle ihrer Gruppe, oder die Voraussetzung dafür steht im Vertrag

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln (Entscheidung des Nutzers vom 2026-10-18; Neuschnitt der Wellen 2 und 3 am selben Tag, vorher welle-gegenproben-und-mutanten; bis dahin ohne Welle). Angelegt am 2026-10-15 bei der Closure von
`slice-harness-randform-stelle-und-lesart` (Register-Ausgang *geplant*: das dritte Auftreten
erreicht die Schwelle); die Stelle in der Reihenfolge war entschieden (Nutzer, 2026-10-15) und ist seit dem 2026-10-16 hinter dem Release (Entscheidung des Nutzers; §4 *Start*).

**Bezug:** kein Anspruch des Produkts — die Prüfung gilt dem Gate `make test-postgres-versionen-gegenprobe`, nicht einer Zusage
des Produkts. Bindung an das Register: `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`;
Bindung an Entscheidungen: [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) (Versionsmatrix der Integrationstests und ihre Gegenprobe).

**Berührte Spec-Stellen:** `SPEC-053` (Punkt 7, *Aufräumen*, mit Punkt 14 zum nachlaufenden Aufräumen nach `SIGTERM`, und die
Zusage der Gegenprobe, dass Mutationen nur die Fälle ihrer Gruppe rot färben; Vertrag in `harness/vertraege/SPEC-053-postgres-versionen.md`)

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Kopplung zwischen der Mutation `p7-int-treiber` und den TERM-Fällen der Gruppe `p7-sig` in
`tools/dev/test/postgres-versionen-gegenprobe.sh` ist in ihrer Ursache ermittelt, und die Gegenprobe der Versionsmatrix färbt unter
dieser Mutation nur die Fälle ihrer Gruppe (`p7-sigINT-gruppe`, `p7-sigINT-skript`) — oder `SPEC-053` nennt die Bedingung,
unter der die Zusage gilt, und die Gegenprobe prüft sie.

**Sendung der drei Geber (Gegenstand):** Das Register zählt dreimal dasselbe Rot im ersten vollen `make gates` eines Slice, der Gate und
Gegenprobe nicht berührt: `slice-v1-abschluss-postgres-versionen-zusage` (zusätzlich rot `p7-sigTERM-gruppe`),
`slice-v1-abschluss-protokollrand` (zusätzlich rot `p7-sigTERM-skript`, V-8) und `slice-harness-randform-stelle-und-lesart`
(zusätzlich rot `p7-sigTERM-skript`, Review F-4, Verify V-4). Einzellauf und zweiter voller Lauf waren jeweils grün, die Ursache ist nicht
untersucht. Dieser Slice nimmt alle drei an (`AGENTS.md` §3.13); die DoD führt sie.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung an `tools/test/run-postgres-versionen.sh` oder an der Zusage des Aufräumens nach `SIGINT` und `SIGTERM` (`SPEC-053` Punkt 7), wenn die Ursache dort liegt — anderer Vorgang: Das wäre Arbeit am Gegenstand der Gegenprobe, nicht an ihr. Liegt die Ursache im Runner, geht der Slice mit dem Befund zurück zum Planner (§4), statt den Runner mitzunehmen.
- Die Dauer der Gegenprobe und ihre Aufteilung — gehört `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`: Dort geht es um die Dauer, hier um den Ausgang. Eine Verlängerung der Wartezeit zwischen Fällen als Gegenmittel braucht dessen Beobachtung (Abwägung im Architect, §6).
- Die Versionsliste und die Fälle der Punkte 1 bis 6 und 8 bis 12 von `SPEC-053` — Bestand bleibt bewusst stehen: Keine Beobachtung nennt sie; die Flanke trifft allein die Gruppe `p7-sig`.
- Produkt-Code, Benutzerhandbuch, Lastenheft — Schicht-Abgrenzung: Der Slice ändert ein Test-Skript des Gates und dessen Vertrag.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** Die Ursache der Kopplung zwischen `p7-int-treiber` und den TERM-Fällen von `p7-sig` ist ermittelt und belegt (ein Lauf, der sie zeigt, nicht die Beschreibung des Musters), und `SPEC-053` trägt die Entscheidung dazu: entweder die Zusage „Mutationen färben nur die Fälle ihrer Gruppe“ unverändert, mit der Ursache behoben, oder die Bedingung, unter der sie gilt (Aufruf-Voraussetzung, `AGENTS.md` §3.12). Die Sendung der drei Geber `slice-v1-abschluss-postgres-versionen-zusage`, `slice-v1-abschluss-protokollrand` und `slice-harness-randform-stelle-und-lesart` ist damit angenommen.
- [ ] **Liefer-Punkt 2:** `tools/dev/test/postgres-versionen-gegenprobe.sh` färbt unter `p7-int-treiber` nur `p7-sigINT-gruppe` und `p7-sigINT-skript`; ein Fall mit Mutation hält die Entscheidung aus Punkt 1 (rot, wenn die Kopplung wiederkehrt), und eine Wiederholung der Gruppe `p7-sig` in der Zahl, die der Architect vor dem Code festlegt, bleibt grün.
- [ ] Mutation je Zusage des geänderten Vertrags (`AGENTS.md` §3.10), vom Implementer selbst rot gesehen; Bericht nennt je Zusage: Zusage · Mutation · roter Test.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` von `geplant` auf `verkörpert` (Zielort und Anker) oder, wenn die Flanke nicht reproduzierbar bleibt, mit Begründung; der Zähler folgt den Dateien in `evidence/`.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-053-postgres-versionen.md` | update | Entscheidung zur Ursache in Punkt 7 und 14; vor dem Code vom Architect geschrieben (`AGENTS.md` §3.12) |
| `tools/dev/test/postgres-versionen-gegenprobe.sh` | update | Gruppe `p7-sig` (Fälle und Mutationen `p7-int-treiber`, `p7-term-treiber`, `p7-term-runner`); Ansatz nach der Ursache, nicht vorab |

Welche Zeile in `harness/README.md` (Sensors, `make test-postgres-versionen-gegenprobe`) mitzieht, sagt der Architect nach der Ursache; trägt der Slice sie, zählt er neu (§8).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice steht in der wellenlosen Reihe der fünf Harness-Slices an letzter Stelle, hinter `slice-harness-schaetzung-planner-faktor` und hinter [welle-v1-nachtrag](../done/welle-v1-nachtrag/welle-v1-nachtrag.md) (in `done/`; Voraussetzung der Stelle, nicht der Technik). Er stand bis zum 2026-10-16 nach dem Verbindungsende-Block von welle-v1-abschluss (hinter `slice-v1-abschluss-verbindungsende-doku`, Schritt 42) und vor den wellenlosen Slices `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage` (Entscheidung des Nutzers vom 2026-10-15) und ist seitdem hinter das Release verschoben (Entscheidung des Nutzers vom 2026-10-16, Drift-Log der Roadmap); `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage` liegen vor dem Release. Er blockiert kein Gate: Der Stop-Hook gab in allen drei Auftreten nach dem zweiten vollen Lauf frei. Wiederholt sich das Rot vor seinem Start in einem weiteren Slice, bleibt es ein Auftreten im Zähler dieser Beobachtung. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools/test tools/dev/test harness/vertraege harness/README.md`, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist und beim Start eingetragen wird; **über 200** geht der Slice zurück zur Zerlegung. **Meldeschwelle 160:** Erreicht die Summe 160, hält der Implementer an und meldet dem Planner. Schätzung des Planners: 60 bis 120 (Mitte 90), unsicher, solange die Ursache fehlt; Marke 200 und Meldeschwelle 160 bleiben Schätzung, der Architect setzt sie neu, sobald die Ursache feststeht (Entscheidung des Nutzers vom 2026-10-15); die Gegenprobe trägt rund das Doppelte des Produktcodes, Tabellen- und Fall-Code eingeschlossen.
- `in-progress` → `open`: Die Flanke lässt sich nicht herstellen und die Ursache nicht belegen (kein Lauf zeigt sie in der Zahl von Wiederholungen, die der Architect festlegt) — der Slice wartet auf ein weiteres Auftreten, und die Beobachtung geht mit Begründung nach *offen* zurück, nicht nach *gestrichen*.
- `in-progress` → `open`: Die Ursache liegt im Runner `tools/test/run-postgres-versionen.sh` oder im Aufräumen nach `SIGTERM` (`SPEC-053` Punkt 14) — der Gegenstand gehört einem anderen Slice (§1); der Planner schneidet ihn.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, `make gates` grün, Closure-Notiz mit Lerneintrag geschrieben (Form: geschärfte Regel im Vertrag oder neuer Fall der Gegenprobe).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12; zu entscheiden vom Architect vor dem ersten Code-Commit, in `SPEC-053`):

- *Aufruf-Voraussetzung:* Die Fälle der Gruppe `p7-sig` laufen nacheinander in einer Gegenprobe; ob ein Fall beginnt, während der vorige noch aufräumt (`SPEC-053` Punkt 14: das Aufräumen folgt bis zu einige Sekunden nach `SIGTERM` an die Prozessgruppe), ist die Frage, die der Plan stellt, nicht beantwortet. Ob die Zusage an „nacheinander, der vorige Fall ist beendet“ gebunden wird, entscheidet der Architect.
- *Reihenfolge der Fälle:* Die Gruppe führt `INT:gruppe TERM:gruppe INT:skript TERM:skript`; ob die Kopplung an die Reihenfolge hängt, ist zu zeigen, nicht zu vermuten.
- *Aufrufarten des Treibers:* `make -n`, `-q`, `-t`, `-i` und `MAKEFLAGS` aus der Umgebung ändern an der Hängung der Gegenprobe an ihrer Liste (seit `slice-harness-gegenproben-aus-gates` nicht mehr `GATE_CHECKS`) nichts, solange der Slice das Rezept nicht anfasst; der Architect bestätigt es oder nennt die Entscheidung je Flag.
- *Wiederholungszahl:* Wie oft die Gruppe `p7-sig` laufen muss, damit „bleibt grün“ eine Aussage über eine zeitabhängige Flanke ist, legt der Architect fest; eine einmalige Grünmeldung ist es nicht (drei Einzelläufe waren grün, der volle Lauf war rot).
- *Stelle einer Nachricht in einer Folge:* nicht berührt; der Vertrag enthält keine Nachrichtenfolge.

**Risiken:**

- Die Flanke lässt sich nicht herstellen, die Ursache bleibt eine Beschreibung des Musters — **Ausgang:** offen bis Closure (Rückführung nach `open/`, §4).
- Die Ursache liegt im Runner oder im nachlaufenden Aufräumen und nicht in der Gegenprobe — **Ausgang:** offen bis Closure (Folge-Slice mit Kennung, §4).
- Eine Wartezeit zwischen den Fällen behebt das Rot, verlängert aber jeden vollen Lauf (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`) — **Ausgang:** offen bis Closure.
- Der Slice wächst über Marke oder Meldeschwelle von §4 — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-15). Treffer in der Sub-Area `REPO`: `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` ist der Gegenstand (3×, Schwelle erreicht, `geplant` mit dieser Kennung); `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` wird berührt, falls die Abhilfe Wartezeit kostet (§6, drittes Risiko). Kein weiterer Treffer.

Nachgezählt (Planner, 2026-10-15, `AGENTS.md` §3.13): Der Slice ist neu; die Sendung der drei Geber (dieselbe Flanke, dreimal) ist ein Gegenstand und kein zusätzlicher Liefer-Punkt. Liefer-Punkte: 2 (Ursache und Entscheidung im Vertrag; Gegenprobe). Schichten: zwei, das Gate-Werkzeug (`tools/test`) und der Vertrag (`harness/vertraege`); Planung zählt keine Schicht, kein Produkt-Code, keine Nutzer-Doku. Zieht die Zeile in `harness/README.md` mit, zählt der Architect neu; läge der Slice darüber, geht der Punkt an den Planner. Geschätzt in der Metrik von §4: 60 bis 120 (Mitte 90) gegen Marke 200, unter 70 % der Marke; der Plan bleibt ungeteilt. Der Nehmer trägt keine Sendung in einen anderen Slice ein.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

Nachgezählt beim Verschieben hinter das Release (Planner, 2026-10-16, `AGENTS.md` §3.13): Start und Stelle geändert, Gegenstand und DoD unverändert; der Slice trägt keine Sendung ein und nimmt keine an: Liefer-Punkte weiter 2, Schichten weiter zwei. Kein Slice vor dem Release nennt ihn als Voraussetzung; tritt das Rot der Gruppe `p7-sig` bis dahin wieder auf, bleibt es ein Auftreten im Zähler von `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`.

**Nachgezählt** (`AGENTS.md` §3.9, §3.13, Planner, 2026-10-18, im Commit, der den Kettenumbau `slice-harness-gegenproben-aus-gates` mit seinem Vertrag (`SPEC-060`) nachzieht): Die Gegenprobe der Versionsmatrix hängt seit `slice-harness-gegenproben-aus-gates` an der Container-Liste und nicht mehr an `GATE_CHECKS`; die Aufrufarten-Randform nennt die Liste. Liefer-Punkte und Schichten ändern sich nicht.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (2) und Schichten (2) unverändert. Die Marke von §4 misst mit `-- tools/test …`; weil die Gegenprobe jetzt unter `tools/dev/test/` liegt, nennt der Befehl dort zusätzlich `tools/dev/test`.
