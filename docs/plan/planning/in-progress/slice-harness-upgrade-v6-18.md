# Slice slice-harness-upgrade-v6-18: Baseline v6.18.0 nachziehen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er
von der nächsten Welle-Closure. Vorgezogen nach Entscheidung des Nutzers vom 2026-10-10: nach `slice-doku-ist-stand` und vor
`slice-v1-abschluss-einspielen-extended` (WIP-Limit 1); welle-v1-abschluss geht danach in der
Reihenfolge ihres §5 weiter.

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage: adoptierter Stand, Asset-Quelle, „keine inhaltlichen Adaptionen“ — der Slice hebt den Stand und prüft die Aussage gegen den neuen). Keine Anforderung des Lastenhefts im Scope: Der Slice ändert die Prüfumgebung, nicht das Produkt.

**Berührte Spec-Stellen:** — (die Werkzeugverträge im Abschnitt für Harness-Werkzeuge der
Spezifikation bleiben unverändert; verlangt der Abgleich eine Änderung, zieht sie den Kopf
nach, `AGENTS.md` §3.9, und ist nach §4 eine Rückführung)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Repo führt die vendored Baseline des Kurs-Stands `v6.18.0` (`make
baseline-verify` grün und nennt `v6.18.0`, alle lebenden Verweise und Pins auf den Stand
nachgezogen, `v6.16.0` entfernt); der Abgleich `v6.16.0` → `v6.18.0` gegen die verkörperten
Regeln dieses Repos ist gemacht, und jeder Befund daraus hat genau einen Ausgang.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-10. Der Freshness-Audit in
`slice-harness-d-check-v0-85` (dort §7) fand `v6.17.0`; der Nutzer verschob das Upgrade, weil
eine neuere Version angekündigt war. Die Release-Liste des Kurs-Repos (`pt9912/ai-harness-course`)
führt jetzt `v6.18.0` mit dem Asset `lab-regelwerk.zip` und einer `SHA256SUMS` zum Release. Das
inhaltliche Delta laut Vorprüfung des Nutzers, über beide Releases:

- die Disjunktheit geteilter Gate-Index-Teile (kein Target in zwei Teilen; seit `v6.17.0`) —
  dieses Repo führt keinen Gate-Index in Teilen (`slice-harness-d-check-v0-85`, dort §6);
- eine Anleitung zum d-check-Modul `reviews`;
- der Dateiname des Review-Reports mit der vollen Slice-Kennung.

Der Abgleich hält den ganzen Baum per `diff -r` gegen den Pin, nicht nur diese drei Punkte
(`BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen`, §8).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das d-check-Modul `reviews` aktivieren (`.d-check.yml`, Zeile `modules:`) — ein anderer
  Vorgang: Es verschärft das Doku-Gate und ist eine eigene Entscheidung; der Nutzer plant dafür
  einen eigenen Slice (Entscheidung vom 2026-10-10), angelegt ist er nicht. Dieser Slice prüft
  nur, ob die Anleitung eine verkörperte Regel berührt, und trägt die Regelblock-Zelle nach dem
  Zustand dieses Repos ein (*liegt bei, nicht verdrahtet*, solange das Modul aus ist).
- Der Pin von d-check (`d-check.mk`, `ghcr.io/pt9912/d-check:v0.86.1`) — Bestand bleibt bewusst
  stehen: Der Slice hebt nur die Baseline. Verlangt der neue Stand eine neuere d-check-Version,
  ist das ein Befund mit Adresse, kein stiller Bump hier.
- Produkt-Code, Tests des Produkts, `Dockerfile` und die Skripte unter `tools/` —
  Schicht-Abgrenzung: Der Slice ändert den vendored Baum unter `.harness/baseline/`, die
  Symlinks unter `.claude/rules/`, Harness-Dokumente (`harness/`, `AGENTS.md`, `.claude/`,
  `.harness/skills/`) und Planungsdokumente. Kein `.go`-File im Diff.
- Ein erneuter Bootstrap mit ai-harness-init (tool-eigene Fragmente unter `harness/mk/`,
  Skripte unter `tools/harness/`) — ein anderer Vorgang: Das Asset trägt nur Regelwerk und
  Vorlagen. `slice-harness-gate-index-werkzeug-teil` bleibt in `open/` mit seinem Trigger.
- Eingefrorene Zeitdokumente (Pläne und Stubs in `done/`, Archive, Review- und
  Verifikations-Reports unter `docs/reviews/`, `harness/conventions/done/`, Belege im Register)
  — Bestand bleibt bewusst stehen: Sie nennen den Stand ihrer Zeit; ein Review-Report behält
  seinen Dateinamen, auch wenn die neue Regel ihn anders bildet. Befunde in ihnen nach §6
  *Befunde in eingefrorenen Dokumenten*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Baseline `v6.18.0` vendored:** `.harness/baseline/v6.18.0/{regelwerk,templates}/` mit
      `SHA256SUMS` aus dem Release-Asset `lab-regelwerk.zip` des Tags, dessen sha256 vor dem
      Entpacken gegen die `SHA256SUMS` des Releases geprüft ist (Wert und Quelle in §7);
      `.harness/baseline/v6.16.0/` ist im selben Commit entfernt; `make baseline-verify` grün
      und nennt `v6.18.0`. `harness/conventions.md` §Baseline (Stand, Datum der Adoption) und
      §Adoptierte Konventions-Quellen (Asset-URL, Stand-Zeile aus `regelwerk/README.md`), die
      Asset-URL in `AGENTS.md` §1 und die Messzeile der Regelblock-Tabelle nennen `v6.18.0`;
      `grep -rn "v6\.16\.0"` über `AGENTS.md`, `README.md`, `harness/` (ohne
      `harness/conventions/done/` und ohne die Zeile `MR-001` in §Aufgelöste Adaptionen, die
      den Stand ihrer Auflösung nennt), `.claude/`, `.harness/skills/`, `tools/`, `Makefile`,
      `*.mk` und `spec/` ist leer, die beiden Symlinks unter `.claude/rules/` lösen auf
      `v6.18.0` auf, und `find . -path ./.git -prune -o -xtype l -print` ist leer.
- [ ] **Abgleich `v6.16.0` → `v6.18.0`:** Jeder Regelblock des neuen Baums steht in der Tabelle
      *Welche Regelblöcke des Baums hier einen Träger haben* in `harness/conventions.md` mit
      genau einem Wert; neue, umbenannte und entfallene Blöcke sind in §7 einzeln genannt. Das
      Delta der Regelwerk-Dateien und der Vorlagen (`diff -r` über den ganzen Baum) ist gegen
      die verkörperten Regeln (`AGENTS.md` §3, die Commands unter `.claude/commands/`, die
      Rollen-Typen unter `.claude/agents/`, die Skills unter `.harness/skills/`,
      `docs/plan/planning/observations/README.md`) und gegen die ausgefüllten Singletons
      (`AGENTS.md`, `harness/README.md`, `harness/conventions.md`, Spec-Straten) gehalten, der
      Adaptions-Durchgang ist gemacht (`MR-000`, die aufgelöste `MR-001`). Die drei Punkte der
      Vorprüfung (§1 *Herkunft*) stehen je mit Ausgang in §7, der Dateiname des Review-Reports
      gegen `.harness/skills/reviewer.md` und den Bestand unter `docs/reviews/`. Jeder Befund hat
      genau einen Ausgang: behoben (klein, Maß in §4), Folge-Slice mit Kennung in `open/` oder
      Eintrag im Beobachtungs-Register.
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
| `.harness/baseline/v6.18.0/{regelwerk,templates}/`, `SHA256SUMS` | neu | aus `lab-regelwerk.zip` des Tags `v6.18.0`, sha256 vor dem Entpacken geprüft; Netz nach §6 *Netzzugriff* |
| `.harness/baseline/v6.16.0/` | entfernt | `tools/harness/baseline-verify.sh` duldet ein Tag-Verzeichnis zur Zeit (§6 *Alte Baseline*) |
| `.claude/rules/modul-01-entwicklungszyklus.md`, `.claude/rules/modul-05-planning-harness.md` | update (Symlink-Ziel) | zeigen auf `v6.16.0`; im Commit des Tauschs auf `v6.18.0` (§6 *Was am alten Pfad hängt*) |
| `harness/conventions.md` | update | §Baseline, §Adoptierte Konventions-Quellen (Asset-URL, Stand-Zeile), Regelblock-Tabelle (Messzeile, Werte gegen den neuen Baum), `MR-000` nach dem Adaptions-Durchgang |
| `AGENTS.md` | update | §1 Asset-URL; §3 nur, wenn der Abgleich eine verkörperte Regel als gegenstandslos oder widersprüchlich findet und die Anpassung klein ist |
| `.claude/commands/*.md`, `.claude/agents/*.md`, `.harness/skills/*.md`, `docs/plan/planning/observations/README.md` | prüfen, update falls klein | Abgleich gegen das Delta, darunter der Dateiname des Review-Reports in `.harness/skills/reviewer.md` und `.claude/commands/implement-slice.md`; Verweise auf den Baum nennen `.harness/baseline/<tag>/…` und bleiben so |
| `harness/README.md`, `harness/sensors/*.md` | prüfen, update falls klein | Singleton-Abgleich gegen die Vorlagen des neuen Stands |
| `docs/plan/planning/open/` (Folge-Slices, falls Befunde) | neu | per `cp` aus `slice.template.md` des **neuen** Stands, je Befund, der nicht klein ist |
| `docs/plan/planning/in-progress/roadmap.md` | update (Drift-Log) | eine Zeile je angelegtem Folge-Slice |

**Ansatz:**

- Reihenfolge der Commits nach dem Vorbild `slice-harness-upgrade-v6-16` (dort §3 *Ansatz*):
  (1) Baum tauschen und Verweise nachziehen in einem Commit, damit `make baseline-verify` an
  keinem Commit zwei Tag-Verzeichnisse sieht; (2) Abgleich und Regelblock-Tabelle; (3) kleine
  Anpassungen und Folge-Slices. Der Vergleich der Bäume läuft vor (1) in einem Temp-Baum
  außerhalb des Repos (`diff -r` alt gegen neu), sein Ergebnis steht in §7.
- Der Abgleich geht nach Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit
  der vendored Baseline vor: Adaptions-Durchgang (auch aufgelöste Einträge), Form-Vergleich der
  Vorlagen (Singletons nacharbeiten, wiederkehrende Vorlagen append-only), Stichprobe gegen den
  Bestand.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-doku-ist-stand` liegt in `done/` (WIP-Limit 1). Erster Schritt nach dem Start, vor dem ersten Commit an Baum oder
Konventionen: Der Architect entscheidet die Randformen aus §6
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`); der Slice legt keinen neuen Vertrag
an, aber er ändert die Grundlage des Gates `make baseline-verify`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Abgleich verlangt mehr als kleine
  Anpassungen — Maß wie in `slice-harness-upgrade-v6-16` (dort §4): mehr als ein Satz je
  betroffener verkörperter Regel, eine Änderung an einem Werkzeugvertrag in der Spezifikation,
  eine neue oder ersetzende ADR, eine neue Adaption `MR-<NNN>`, ein Eingriff in ein Skript unter
  `tools/` oder eine Änderung an `.d-check.yml`. Dann trägt dieser Slice nur Baum und Verweise
  (Liefer-Punkt 1) und den Abgleich als Befundliste; die Umsetzung geht an Folge-Slices mit
  Kennung.
- `in-progress` → `open` (blockiert — Carveout?): Das Asset ist nicht über einen Weg
  beschaffbar, den §6 *Netzzugriff* zulässt, oder sein sha256 weicht von der `SHA256SUMS` des
  Releases ab; dann bleibt `v6.16.0`, und der Nutzer entscheidet.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (darin `make baseline-verify` mit `v6.18.0`), `grep` aus dem
ersten Liefer-Punkt leer, Closure-Notiz mit Lerneintrag; jeder Befund des Abgleichs hat seinen
Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — **entschieden vom Architect am 2026-10-10**, nach dem
Start und vor dem ersten Commit (§4). **Ort ist dieser Abschnitt**, nicht die Spezifikation und
keine ADR: Der Slice legt keinen Werkzeugvertrag an und ändert keinen (kein `make`-Ziel, kein
Skript unter `tools/`, `.d-check.yml` und `d-check.mk` unberührt, Kopf bleibt `—`); es ist ein
einmaliger Wartungsschritt nach Baseline-Regelwerk `modul-02-harness-bootstrap.md`
§Freshness-Audit der vendored Baseline, wie in `slice-harness-upgrade-v6-16` (dort §6). Die
Punkte mit „Vorschlag“ unten sind **bestätigt, wie formuliert**, soweit der Block *Entscheidungen
des Architect* nach der Liste sie nicht ändert oder ergänzt. Verlangt die Umsetzung doch eine
Änderung an einem Werkzeugvertrag, gilt §4 *Rückführungen*. Was hier nicht steht, entscheidet
der Implementer nicht, er gibt es zurück (`.claude/commands/implement-slice.md`,
Randform-Rückgabe).

- **Alte Baseline** — Vorschlag: `v6.16.0` fällt im selben Commit, der `v6.18.0` anlegt; das
  Delta zeigen `diff -r` im Temp-Baum vorher und `git diff -M` danach (dort *Alte Baseline*).
- **Was am alten Pfad hängt** — die beiden Symlinks unter `.claude/rules/` zeigen auf
  `v6.16.0`; ist eine der beiden Dateien im neuen Stand umbenannt oder entfallen: Rückgabe an
  den Architect. Ob weitere lebende Stellen den Pfad nennen, zeigt eine Probe in einer
  HEAD-Kopie ohne `.harness/baseline/v6.16.0/` mit `make docs-check` vor dem Tausch.
- **`SHA256SUMS`** — Soll-Wert des Assets ist die `SHA256SUMS` des Releases (Rang 1 in
  `slice-harness-upgrade-v6-16`, dort *Netzzugriff*); offen ist, ob die vendored `SHA256SUMS`
  wie bisher nach dem Entpacken neu erzeugt wird oder die des Assets übernommen wird, wenn sie
  dieselbe Form hat (`<sha256>  <pfad>` relativ zum Tag-Verzeichnis, über `regelwerk/` und
  `templates/`, nach `LC_ALL=C sort`).
- **Netzzugriff** — Vorschlag: einmaliger Wartungsschritt mit dem gepinnten Image der Stufe
  `deps`, Mount nur eines Temp-Verzeichnisses außerhalb des Repos, ins Repo per `cp` auf dem
  Host (dort *Netzzugriff*); kein `make`-Ziel, kein Gate mit Netz.
- **Zwei Releases in einem Sprung** — `v6.17.0` wird nicht vendored; der Abgleich misst
  `v6.16.0` gegen `v6.18.0`, und §7 nennt je Punkt, mit welchem Release er kam, soweit die
  Release-Liste es sagt.
- **Dateiname des Review-Reports** — gilt die neue Regel für Reports ab dem Upgrade, und
  weicht der Bestand dieses Repos (`<datum>-review-<slice-kennung>.md`) von ihr ab, ist der Weg
  zu entscheiden: Skill und Commands nachziehen (klein) oder Befund mit Adresse; eingefrorene
  Reports behalten ihren Namen (§1).
- **Anleitung zum Modul `reviews`** — die Regelblock-Zelle sagt den Zustand dieses Repos
  (Modul aus); berührt die Anleitung eine verkörperte Regel, etwa die Zitier-Form der Reports,
  ist das ein Befund des Abgleichs, keine Aktivierung (§1).
- **Befunde in eingefrorenen Dokumenten** — Vorschlag: Form-Befunde, die das Dokument nicht
  umdeuten (Link → Code-Span mit gleichem Text und unverändertem Pfad), je Datei eine Zeile in
  §7; alles, was Inhalt änderte, ist Halt und Rückgabe (dort *Befunde in eingefrorenen
  Dokumenten*).

**Entscheidungen des Architect (2026-10-10)** — Ergänzungen und Schärfungen zu den Punkten oben:

- **Reihenfolge der Schritte.** (0) Vorbereitung, nichts im Repo: Asset in ein Temp-Verzeichnis
  außerhalb des Repos laden (Weg wie *Netzzugriff*), sha256 gegen die `SHA256SUMS` des Releases
  halten (Halt bei Abweichung, §4), `unzip -l` auf absolute Pfade, `..`, Symlinks (Halt),
  entpacken, `diff -r` alt gegen neu, Probe in einer HEAD-Kopie ohne
  `.harness/baseline/v6.16.0/` mit `make docs-check`. (1) **Ein Commit:** Baum tauschen,
  Symlinks, `AGENTS.md` §1 (Asset-URL), `harness/conventions.md` §Baseline, §Adoptierte
  Konventions-Quellen (Asset-URL, Stand-Zeile) und Messzeile der Regelblock-Tabelle — so ist der
  `grep` aus Liefer-Punkt 1 an jedem Commit leer und `make baseline-verify` sieht nie zwei
  Tag-Verzeichnisse. (2) Abgleich: Regelblock-Tabelle (Werte), `MR-000`, Aufgelöste Adaptionen.
  (3) Kleine Anpassungen und Folge-Slices, je Befund ein Commit-Thema. Closure-Reihenfolge nach
  `AGENTS.md` §3.3.
- **Verweise auf `v6.16.0` in Zeitdokumenten** (Pläne in `done/`, `docs/reviews/`,
  Register-Belege, `harness/conventions/done/`): **nicht anfassen.** Einzige Ausnahme: ein
  Befund der Probe (d-check meldet `target-missing`/`anchor-missing`) wird nach *Befunde in
  eingefrorenen Dokumenten* behandelt (Link → Code-Span mit unverändertem Pfad, eine Zeile in
  §7). Lebende Pläne, die `v6.16.0` als Maß ihrer Aussage nennen (`open/slice-harness-gate-index-werkzeug-teil`,
  Herkunft „v6.13.0 → v6.16.0“) und die Roadmap-Drift-Log-Zeilen sind Chronik und bleiben; nennt
  der Plan `v6.16.0` als aktuellen Stand statt als Herkunft, ist es ein Befund des Abgleichs
  (Satz, im Maß von §4). Die Zeile `MR-001` in §Aufgelöste Adaptionen bleibt (Stand ihrer
  Auflösung).
- **`.d-check.yml` / `scan.ignore`.** Unverändert. `.harness/**` steht in `scan.ignore`, der neue
  Baum liegt unter demselben Pfad; d-check sieht ihn nicht. Jede Änderung an `.d-check.yml` ist
  eine Rückführung (§4); die Aktivierung des Moduls `reviews` ist Nicht-Ziel (§1).
- **`SHA256SUMS`.** Es zählt die `SHA256SUMS` des **Releases** (Soll-Wert für das Asset, Halt bei
  Abweichung). Die vendored `SHA256SUMS` unter `v6.18.0/` wird wie unter `v6.16.0` **nach dem
  Entpacken neu erzeugt** (`<sha256>  <pfad>` relativ zum Tag-Verzeichnis, nur `regelwerk/` und
  `templates/`, `LC_ALL=C sort`), weil das Skript genau diese Form liest und die Vollständigkeit
  gegen die Datei prüft. Hat das Asset eine Datei gleichen Namens in der Form, wird der entpackte
  Inhalt zusätzlich gegen sie geprüft und das Ergebnis in §7 genannt; weicht ihre Form ab, wird
  sie nicht übernommen. Die sha256 des Assets steht in §7 *Belege*, in keinem lebenden Dokument.
- **Regelblock-Tabelle.** Vorprüfung des Koordinators: keine neuen, umbenannten, entfallenen
  Blöcke. Der Implementer belegt das mit `ls` beider Bäume (Dateinamen und Abschnitte
  `modul-*.md`, `grundlagen-*.md`); jede Abweichung ist ein Eintrag in §7 und eine Zeile mit
  genau einem Wert in der Tabelle. Zellen von Blöcken ohne Delta bleiben unberührt; zu
  prüfen sind die Zellen der Blöcke, die das Delta berührt: `modul-10` (Dateiname des
  Reports), `modul-13` (Gate-Index, Disjunktheit), `modul-15` (Modul `reviews`; Zelle *Doku-Konsistenz-Drift*
  bleibt *liegt bei, nicht verdrahtet*, solange `modules:` das Modul nicht führt) und
  `grundlagen-harness-dateien.md` (Zelle „Teil des Gate-Index“ bleibt *kommt nicht mit*, Dauer
  unverändert, solange `slice-harness-gate-index-werkzeug-teil` offen ist). Die Messzeile nennt
  `v6.18.0`.
- **Freshness-Audit-Aussage.** Die Aussage „Sensor kommt nicht mit“ im Abschnitt *Der
  mitgelieferte Baum altert still* bleibt unverändert und wahr; sie nennt keinen Tag. In §7 steht
  als Audit-Beleg: Release-Liste gelesen am Tag des Laufs, neuester Tag `v6.18.0` = Pin.
  `v6.17.0` ist durch den Sprung abgedeckt (*Zwei Releases in einem Sprung*).
- **Stichprobe gegen den Bestand** (Modul 02: ein Abschnitt, rotierend, Komplementärmenge zum
  Delta). Letzter Lauf: `modul-05-planning-harness.md` §Offene Risiken (v6.16-Slice). Dieser
  Lauf: **`modul-07-carveouts.md` §Ziel-Form: Carveout** gegen `docs/plan/carveouts/`
  (Vorlage `carveout.template.md`) und den Bestand dort; fehlt der Abschnitt im Komplement (hat
  `diff -r` ihn berührt), nimmt der Implementer den nächsten ohne Delta in der Reihenfolge
  `modul-04-adrs.md` §Ziel-Form, `modul-03-spec.md` §Ziel-Form: Architektur-Sicht und nennt die
  Wahl in §7. Frage je Regel: im ausgefüllten Artefakt oder als deklarierte Abweichung? Ein
  einzelner Fund geht den Weg einer Diskrepanz (Folge-Slice oder Register), mehrere treffen die
  `MR-000`-Aussage (Rückführung §4, neue Adaption).
- **Dateiname des Review-Reports** (Vorprüfung nennt „tun wir bereits“ — **das trifft nur halb**).
  Bestand: `<datum>-<rolle>-<slice-Kennung>.md` (`-review-`, `-verifikation-`, `-mutationen-`,
  Folgeläufe `folge-review-N-…`); neue Regel laut Vorprüfung:
  `<YYYY-MM-DD>-<slice-Kennung>.md`, Folgeläufe Suffix `-r2`. Die volle Kennung führt der Bestand,
  Rollen-Infix und Folgelauf-Form weichen ab. **Entscheidung nur für das, was der Wortlaut des
  neuen Stands trägt:** Der Implementer liest die Regel im Temp-Baum (`modul-10`, `review-report.template.md`)
  und trägt die Folgelauf-Form `-r2` (Suffix) in `.harness/skills/reviewer.md` §Output und die
  Commands nach, **wenn** die Regel den Rollen-Infix nicht verbietet. Verbietet sie ihn, ist das
  eine **Nutzerentscheidung** (Optionen im Bericht), nicht Sache des Implementers. Eingefrorene
  Reports behalten ihren Namen. Kein Sensor liest den Namen (`.d-check.yml`, `.githooks/commit-msg`
  und `tools/` nennen kein Muster in `docs/reviews/`; `harness/mk/vorgaben.mk` nimmt den Ordner
  nur aus `slice-mv` aus) — das dritte Risiko ist damit auf den Wortlaut reduziert.
- **Schnitt nach `AGENTS.md` §3.13.** Zwei Liefer-Punkte (Baum und Verweise · Abgleich), Maß
  ≤ 3. Schichten: **Harness** (vendored Baum, Konventionen, Skills, Commands, Symlinks) und
  **Doku** (Planung, Register); kein Code, keine Spec — zwei, Maß ≤ 2. Wächst der Abgleich über
  das Maß aus §4, gilt die Rückführung. Keine Sendung an einen anderen Slice nötig; ein Befund,
  den der Abgleich an `slice-harness-gate-index-werkzeug-teil` oder an den geplanten
  Reviews-Modul-Slice weist, trägt der Implementer dort im selben Commit in §1 oder DoD mit der
  Kennung dieses Slice ein (§3.13); der Reviews-Slice ist **nicht angelegt**, also keine Adresse:
  der Punkt geht als Register-Eintrag, nicht als Zuweisung.

**Risiken:**

- Das Delta berührt eine verkörperte Regel, ohne dass der Abgleich es sieht
  (`BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen`, 1×) — Gegenmittel:
  `diff -r` über den ganzen Baum, je geänderter Abschnitt die verkörperte Stelle gesucht —
  **Ausgang:** offen bis Closure.
- Ein lebender Verweis auf `.harness/baseline/v6.16.0/` bleibt stehen und bricht mit dem
  Entfernen, oder bleibt in einem Code-Span unentdeckt — Gegenmittel: `grep` und `find` aus
  dem ersten Liefer-Punkt — **Ausgang:** offen bis Closure.
- Die Regel zum Dateinamen des Review-Reports passt nicht zum Commit-Träger oder zum
  Doku-Gate (ein Name, den `.d-check.yml` oder `.githooks/commit-msg` anders liest) —
  **Ausgang:** offen bis Closure.

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

Die Closure-Notiz wird bei Closure gefüllt (vor dem `git mv` nach `done/`); bis dahin
stehen hier die Belege des Implementers.

### Belege des Implementers

- *Asset:* `lab-regelwerk.zip` des Tags `v6.18.0`, sha256
  `18e5443b4fca32ce452ec5dcf7cf011d7e04b0d29741b777a8dceaa452dcb8d9`. Soll-Wert Rang (1):
  das Release-Asset `SHA256SUMS` (eine Zeile, nur das Zip) — `sha256sum -c` OK; Rang (2): das
  Feld `digest` des Assets in der Release-API nennt denselben Wert. Abruf einmalig im Image der
  Stufe `deps` (`--user`, Mount nur eines Temp-Verzeichnisses außerhalb des Repos); ins Repo mit
  `cp` auf dem Host. Vor dem Entpacken: 68 Einträge, oberste Ebene nur `regelwerk/` und
  `templates/`, kein absoluter Pfad, kein `..`, kein Symlink (54 × `-rw-r--r--`, 14 × `drwxr-xr-x`).
- *`SHA256SUMS` des Tag-Verzeichnisses:* nach dem Entpacken neu erzeugt (`<sha256>  <pfad>`,
  `regelwerk/` und `templates/`, `LC_ALL=C sort`), 54 Zeilen. Das Asset bringt keine eigene
  Prüfsummenliste mit; die Gegenprobe läuft deshalb gegen die entpackten Asset-Dateien im
  Temp-Baum (gleiche Erzeugung dort): identisch. `make baseline-verify` meldet
  `v6.18.0 OK — 54 Dateien`.
- *Audit:* Release-Liste des Kurs-Repos am 2026-10-10 gelesen (`v6.18.0`, `v6.17.0`, `v6.16.0`,
  `v6.15.0`, `v6.14.1`); neuester Tag `v6.18.0` = Pin. `v6.17.0` ist durch den Sprung abgedeckt.
- *Probe vor dem Tausch:* `make docs-check` in einer HEAD-Kopie ohne
  `.harness/baseline/v6.16.0/` (Temp-Baum außerhalb des Repos): 473 Dateien, 0 Befunde. Es
  gibt keinen Befund in eingefrorenen Dokumenten.
- *Dateimenge und Gliederung:* `diff -rq` alt gegen neu: dieselben 54 Dateien, keine neue,
  umbenannte oder entfallene; die Überschriften (`grep '^#'`, sortiert) aller Regelwerk-Dateien
  sind identisch. Es gibt keinen neuen, umbenannten oder entfallenen Regelblock.
- *Delta je Datei* (`diff -r`, 31 Dateien verschieden; 54 − 31 = 23 unverändert):
  - nur Versionszeile (Quelle-Kommentar bzw. Asset-URL `v6.16.0` → `v6.18.0`, je 2 Zeilen):
    `grundlagen-begriffe`, `-bootstrap`, `-durchsetzungsschicht`, `-klassifikation`,
    `-referenz-richtung`, `-source-precedence`, `-traceability`, `modul-00`, `-01`, `-02`, `-03`,
    `-04`, `-05`, `-06`, `-07`, `-08`, `-09`, `-11`, `-12`, `-14`, `-15`, `-16`,
    `templates/AGENTS.template.md`, `templates/harness/conventions.template.md`;
  - `regelwerk/README.md`: Versionszeilen, Stand-Zeile „Kurs-Welle 161 · 2026-10-10“;
  - inhaltlich: `grundlagen-harness-dateien.md` (Disjunktheit geteilter Gate-Index-Teile:
    wer die Prüfung hat, schaltet sie mit der Vereinigung ein; Release `v6.17.0`),
    `modul-13-quality-gates.md` (derselbe Satz: ein Target in zwei Teilen sieht der Sensor nur
    bei eigener Disjunktheitsprüfung), `modul-10-review-harness.md` (Absatz „Ein Deckungs-Sensor
    prüft nur, was er als Zusage erkennt“; Dateiname des Reports mit der vollen Slice-Kennung),
    `templates/.d-check.yml` (Kommentare zu `authority-disjoint`, d-check ≥ v0.83.0, und die
    Schlüssel des Moduls `reviews`: `match: name`, `require-promises`, `recursive`,
    `skip-pattern`, `skip-allows-empty`, d-check ≥ v0.85.0/v0.86.0),
    `templates/docs/reviews/review-report.template.md` (Dateiname `<YYYY-MM-DD>-<slice-Kennung>.md`,
    ein Review ohne Slice `<YYYY-MM-DD>-<diff-ref>.md`, „etwa mit Suffix `-r2`“),
    `templates/harness/README.template.md` (Anleitung zum Modul `reviews`, ≥ v0.86.0).
- *Die drei Punkte der Vorprüfung — Ausgang je Punkt:*
  - Disjunktheit geteilter Gate-Index-Teile: dieses Repo führt den Gate-Index in einer Datei
    (`harness/README.md` §Sensors); der Satz setzt ein Repo mit Teilen voraus. Kein Befund,
    keine Änderung; die Zelle *Teil des Gate-Index* bleibt *kommt nicht mit*, Dauer unverändert
    (`slice-harness-gate-index-werkzeug-teil` in `open/`).
  - Anleitung zum Modul `reviews`: berührt keine verkörperte Regel dieses Repos (kein Text in
    `AGENTS.md`, `harness/README.md`, den Commands oder dem Skill nennt das Modul); `.d-check.yml`
    bleibt unverändert, das Modul bleibt aus (§1). Kein Befund; die Zelle *Doku-Konsistenz-Drift*
    bleibt *liegt bei, nicht verdrahtet*. Der Absatz „Ein Deckungs-Sensor prüft nur, was er als
    Zusage erkennt“ ist Lesestoff für den geplanten Slice, der das Modul aktiviert (nicht angelegt, keine Adresse, §6).
  - Dateiname des Review-Reports: siehe unten.
- *Dateiname des Review-Reports* (Wortlaut `modul-10-review-harness.md` §Harness-Einordnung
  und `review-report.template.md`): „Abgelegt wird ein Report pro Lauf unter `docs/reviews/`, die
  volle Slice-Kennung im Dateinamen, Folgeläufe als neue Datei statt Überschreibung“; Vorlage:
  `<YYYY-MM-DD>-<slice-Kennung>.md`, „etwa mit Suffix `-r2`“. Der Wortlaut verbietet einen
  Rollen-Infix nicht und verlangt die volle Kennung im Namen; der Bestand (`<datum>-review-<slice-Kennung>.md`,
  `-verifikation-`, `-mutationen-`, `-validierung-`) trägt die volle Kennung. **Akzeptierte
  Abweichung:** der Infix bleibt, weil er die Rolle des Reports im Namen hält und die volle
  Kennung unberührt lässt; eingefrorene Reports behalten ihren Namen (§1). Nachgezogen ist nur die
  Folgelauf-Form `-r2` (Suffix an der Kennung statt Präfix `folge-review-N-`) in
  `.harness/skills/reviewer.md` §Output und `.claude/commands/implement-slice.md` Schritt 21.
  Kein Sensor liest den Namen (`.d-check.yml`, `.githooks/commit-msg`, `tools/`); das dritte
  Risiko aus §6 trifft damit nicht. Gegenstand dieser Abweichung ist ein künftiger Slice, der das
  Modul `reviews` aktiviert (dort `match: name` gegen die Namen prüfen); er ist nicht angelegt.
- *Regelblock-Tabelle:* alle Zellen gegen den neuen Baum gehalten, Messzeile auf `v6.18.0`.
  Geprüft und unverändert: `modul-10` (Skill trägt die Urteilsgrundlage), `modul-13` (Aggregator
  und Fragmente; die Disjunktheit betrifft Gate-Index in Teilen), `modul-15` (beide Zellen),
  `grundlagen-harness-dateien` (beide Zellen), `modul-02` (Freshness-Audit: weiter *kommt nicht
  mit*, die Aussage nennt keinen Tag). Zellen der übrigen Blöcke: kein Delta, unberührt. Es gibt
  keine Abweichung, die eine neue Zeile verlangt.
- *Adaptions-Durchgang:* `MR-000` (keine inhaltlichen Adaptionen) gilt unverändert; das Delta
  führt zu keiner Adaption `MR-<NNN>`. Die aufgelöste `MR-001` trägt ihre Auflösung („seit
  `v6.16.0`“) als Stand ihrer Zeit; das Delta berührt `grundlagen-referenz-richtung.md` nur in
  der Versionszeile.
- *Singletons gegen die Vorlagen des neuen Stands:* `AGENTS.md` (Vorlage: nur Asset-URL, nachgezogen),
  `harness/conventions.md` (Vorlage: nur Asset-URL, nachgezogen), `harness/README.md` (Vorlage:
  nur die Anleitung zum Modul `reviews`, nicht verkörpert, siehe oben).
- *Stichprobe gegen den Bestand* (Modul 02, rotierend; letzter Lauf `modul-05` §Offene Risiken):
  `modul-07-carveouts.md` §Ziel-Form: Carveout — im Komplement des Deltas (nur Versionszeile).
  Gehalten gegen `docs/plan/carveouts/` und `carveout.template.md`: Der Bestand ist leer (nur
  `.gitkeep`), `harness/README.md` nennt „Rote Gates: … bisher keiner“, in der Bindung-Spalte
  steht kein `CO-<NNN>`; die sechs Pflicht-Header-Felder der Regel stehen in der Vorlage
  (Status, Datum angelegt, Letzte Prüfung, Betroffenes Gate, Geltungsbereich, Folge-Slice).
  Kein Fund.
- *Verweisprobe:* `grep "v6\.16\.0"` über `AGENTS.md`,
  `README.md`, `harness/` (ohne `harness/conventions/done/` und die Zeile `MR-001`), `.claude/`,
  `.harness/skills/`, `tools/`, `Makefile`, `*.mk`, `spec/` ist leer;
  `find . -path ./.git -prune -o -xtype l -print` ist leer; die Symlinks unter `.claude/rules/`
  lösen auf `v6.18.0` auf. Verweise auf `v6.16.0` in `open/slice-harness-gate-index-werkzeug-teil.md`
  und im Drift-Log der Roadmap nennen Herkunft, nicht den Stand, und bleiben.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `0859a19` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen` (1×) — genau der Fall
  dieses Slice; Gegenmittel und Risiko in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — der Slice ist wellenlos; der
  Architect entscheidet die Randformen vor dem ersten Commit (§4).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×) — trifft nur, wenn der Abgleich
  eine Änderung an `.d-check.yml` verlangt; die ist hier eine Rückführung (§4), das Modul
  `reviews` bleibt aus (§1). Erreichte der Eintrag mit diesem Slice 3×, wäre er eine Lücke mit
  eigenem Folge-Slice.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (29×, `AGENTS.md` §3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (23×, §3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (6×, §3.13) — verkörpert; die Regelblock-Zellen
  sagen nur den Zustand dieses Repos, der Plan folgt jeder Korrektur, und jeder Folge-Slice
  nimmt seine Sendung an.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
