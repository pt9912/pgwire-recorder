# Slice slice-harness-upgrade-v6-16: d-check v0.82.0 und Baseline v6.16.0 nachziehen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-06 (d-check v0.82.0 und Regelwerk v6.16.0 des Kurs-Repos ai-harness-course sind
erschienen): nach `slice-harness-lint-werkzeug` und vor `slice-harness-blackbox-kern`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage: adoptierter Stand, Asset-Quelle, „keine inhaltlichen Adaptionen“ — der Slice hebt den Stand und prüft die Aussage gegen den neuen). Keine Anforderung des Lastenhefts im Scope: Der Slice ändert die Prüfumgebung, nicht das Produkt.

**Berührte Spec-Stellen:** — (die Werkzeugverträge im Abschnitt für Harness-Werkzeuge der Spezifikation bleiben unverändert; verlangt der Abgleich eine Änderung, zieht sie den Kopf nach, `AGENTS.md` §3.9, und ist nach §4 eine Rückführung)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Repo prüft seine Doku mit d-check `v0.82.0` (Tag und Digest gepinnt,
`make docs-check` grün) und führt die vendored Baseline des Kurs-Stands `v6.16.0`
(`make baseline-verify` grün, alle lebenden Verweise auf den Stand nachgezogen); der
Abgleich v6.13.0 → v6.16.0 gegen die verkörperten Regeln dieses Repos ist gemacht, und
jeder Befund daraus hat eine Adresse — in diesem Slice behoben, wenn er klein ist, sonst
Folge-Slice oder Register-Eintrag.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-06; der Freshness-Audit nach
Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit der vendored Baseline
hat einen neueren Tag in der Release-Liste gefunden. Ein neuer Tag löst einen Review aus,
keinen stillen Auto-Bump — dieser Slice ist der Review.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Tests des Produkts, `Dockerfile` und die Go-Werkzeuge der Prüfumgebung —
  Schicht-Abgrenzung: Der Slice ändert Pins (`d-check.mk`), den vendored Baum unter
  `.harness/baseline/`, Harness-Dokumente (`harness/`, `AGENTS.md`, `.claude/`,
  `.harness/skills/`) und Planungsdokumente; wo d-check `v0.82.0` an `.d-check.yml` eine
  Änderung verlangt, auch diese. Kein `.go`-File im Diff.
- Die Umplanung der Lint-Reihe (`slice-harness-lint`, die vier Umstellungs-Slices, die
  zwei Bereinigungs-Slices, `slice-harness-abdeckung-gate`, `slice-harness-coverage`,
  `slice-harness-mutation`) — Bestand bleibt bewusst stehen: Ihre Reihenfolge ist
  Entscheidung des Nutzers vom 2026-10-06; dieser Slice schiebt sich nur vor
  `slice-harness-blackbox-kern`. Verlangt der Abgleich eine Änderung an einem dieser
  Pläne, ist das ein Befund mit Adresse (der betroffene Slice, `AGENTS.md` §3.9), keine
  stille Umplanung hier.
- Befunde der strengeren d-check-Version oder des Abgleichs, die mehr als eine kleine
  Anpassung verlangen (Maß in §4 *Rückführungen*) — ein Folge-Slice übernimmt sie, je
  Befund mit Kennung, angelegt in `open/` mit der Closure dieses Slice; der Slice
  bereinigt nur, was klein ist.
- Ein erneuter Bootstrap mit ai-harness-init (Neuschreiben der tool-eigenen Fragmente
  und Skripte wie `tools/harness/slice-mv.sh`, `tools/harness/baseline-verify.sh`,
  `harness/mk/*.mk`) — ein anderer Vorgang: Das Release-Asset `lab-regelwerk.zip` trägt
  nur Regelwerk und Vorlagen; den Präfix-Fehler von `make slice-mv` behebt
  ai-harness-init, dieser Slice prüft ihn nicht (§6 *Präfix-Fehler von `make slice-mv`*).
- Die Werkzeugverträge von `kopf-check`, `abdeckung` und `lint` im Abschnitt für
  Harness-Werkzeuge der Spezifikation und ihre Sensor-Dateien — Bestand bleibt bewusst
  stehen, solange der Abgleich keine Änderung verlangt; der Slice legt keinen neuen
  Vertrag an. Verlangt er eine, gilt §4 *Rückführungen*.
- Eingefrorene Zeitdokumente (Pläne und Stubs in `done/`, Archive, Review- und
  Verifikations-Reports unter `docs/reviews/`, aufgelöste Adaptionen in
  `harness/conventions/done/`) — Bestand bleibt bewusst stehen: Sie nennen den Stand
  ihrer Zeit, und eine wiederkehrende Form wird nach Baseline-Regelwerk
  `modul-02-harness-bootstrap.md` §Freshness-Audit nicht rückwirkend umgeschrieben. Wie
  mit einem neuen Befund in einem solchen Dokument umzugehen ist, entscheidet §6
  *Befunde in eingefrorenen Dokumenten*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **d-check `v0.82.0`:** `DCHECK_IMAGE` und `DCHECK_DIGEST` in `d-check.mk` nennen
      Tag `v0.82.0` und den Digest, den `docker pull` für diesen Tag meldet (Herkunft des
      Digest im Bericht: Kommando und Ausgabe); die Adaption des Fragments (Target
      `docs-check`, gepinnter Digest, Kopfkommentar) bleibt erhalten, die übrigen Targets
      folgen der neuen Version nach §6 *Neu erzeugen oder nur umpinnen*. `make docs-check`
      grün; jeder neue Befund der strengeren Version ist im Bericht gezählt und hat
      einen Ausgang — behoben in diesem Slice oder Folge-Slice mit Kennung (§6
      *Befunde in eingefrorenen Dokumenten*).
- [x] **Baseline `v6.16.0` vendored:** `.harness/baseline/v6.16.0/{regelwerk,templates}/`
      mit `SHA256SUMS` aus dem Release-Asset `lab-regelwerk.zip` des Tags, dessen sha256
      vor dem Entpacken geprüft ist (Wert und Quelle im Bericht); `make baseline-verify`
      grün und nennt `v6.16.0`. `harness/conventions.md` §Baseline (Stand, Datum) und
      §Adoptierte Konventions-Quellen, die Asset-URL in `AGENTS.md` §1 und die Messzeile
      der Regelblock-Tabelle nennen `v6.16.0`; `grep -rn "v6\.13\.0"` über `AGENTS.md`,
      `README.md`, `harness/` (ohne `harness/conventions/done/`), `.claude/`,
      `.harness/skills/`, `tools/`, `Makefile`, `*.mk` und `spec/` ist leer, und kein
      Plan in `open/`, `next/` oder `in-progress/` verweist auf einen Pfad unter
      `.harness/baseline/v6.13.0/` (Pläne dürfen den alten Stand als Wort nennen, dieser
      eingeschlossen); `find . -path ./.git -prune -o -xtype l -print` ist leer. Was mit
      dem alten Baum geschieht, nach §6 *Alte Baseline* und *Was am alten Pfad hängt*.
- [x] **Abgleich v6.13.0 → v6.16.0:** Jeder Regelblock des neuen Baums steht in der
      Tabelle *Welche Regelblöcke des Baums hier einen Träger haben* in
      `harness/conventions.md` mit genau einem Wert; neue, umbenannte und entfallene
      Blöcke sind im Bericht einzeln genannt. Das Delta der Regelwerk-Dateien und der
      Vorlagen ist gegen die verkörperten Regeln (`AGENTS.md` §3.1 bis §3.12, die
      Commands unter `.claude/commands/`, die Rollen-Typen unter `.claude/agents/`, die
      Skills unter `.harness/skills/`, `docs/plan/planning/observations/README.md`) und
      gegen die ausgefüllten Singletons (`AGENTS.md`, `harness/README.md`,
      `harness/conventions.md`, Spec-Straten) gehalten; der Adaptions-Durchgang ist
      gemacht, auch für die aufgelöste Adaption in `harness/conventions/done/` und für
      die `MR-000`-Aussage (§6 *Klarstellung zum Technik-Stratum*, *Sensor-Datei*). Jeder
      Befund hat genau einen Ausgang: behoben (klein), Folge-Slice mit Kennung in
      `open/` oder Eintrag im Beobachtungs-Register.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `d-check.mk` | update | `DCHECK_IMAGE` auf `v0.82.0`, `DCHECK_DIGEST` auf den Digest des Tags; Targets nach §6 *Neu erzeugen oder nur umpinnen* |
| `.d-check.yml` | update (nur falls nötig) | nur, wenn `v0.82.0` die Konfiguration anders liest oder ein Modul umbenennt; jede Änderung, die einen Befund wegnimmt statt ihn zu beheben, ist eine Lockerung (`AGENTS.md` §3.6) und braucht eine ADR — dann Folge-Slice |
| Doku mit neuen Befunden von `v0.82.0` | update | kleine Befunde in lebenden Dokumenten beheben; eingefrorene nach §6 *Befunde in eingefrorenen Dokumenten* |
| `.harness/baseline/v6.16.0/{regelwerk,templates}/`, `SHA256SUMS` | neu | aus `lab-regelwerk.zip` des Tags `v6.16.0`, sha256 vor dem Entpacken geprüft; Netz nach §6 *Netzzugriff* |
| `.harness/baseline/v6.13.0/` | entfernt (§6 *Alte Baseline*) | `tools/harness/baseline-verify.sh` duldet ein Tag-Verzeichnis zur Zeit |
| `.claude/rules/modul-01-entwicklungszyklus.md`, `.claude/rules/modul-05-planning-harness.md` | update (Symlink-Ziel) | zeigen auf `v6.13.0`; im Commit des Tauschs auf `v6.16.0` (§6 *Was am alten Pfad hängt*) |
| `harness/conventions/done/MR-001-spezifikations-ort-werkzeugvertraege.md` | update (Form-Reparatur) | Link auf `v6.13.0` wird Code-Span mit unverändertem Pfad (§6 *Was am alten Pfad hängt*) |
| `harness/conventions.md` | update | §Baseline (Stand `v6.16.0`, Datum der Adoption), §Adoptierte Konventions-Quellen (Asset-URL, Stand-Zeile aus `regelwerk/README.md`), Regelblock-Tabelle (Messzeile, neue/umbenannte/entfallene Blöcke), `MR-000` und Aufgelöste Adaptionen nach dem Adaptions-Durchgang |
| `AGENTS.md` | update | §1 Asset-URL; §3.12 Ort einer entschiedenen Randform (Review F-434); §3.x sonst nur, wenn der Abgleich eine verkörperte Regel als gegenstandslos oder widersprüchlich findet und die Anpassung klein ist |
| `.claude/commands/*.md`, `.claude/agents/*.md`, `.harness/skills/*.md`, `docs/plan/planning/observations/README.md` | prüfen, update falls klein | Abgleich gegen das Delta; Verweise auf den Baum nennen heute `.harness/baseline/<tag>/…` und bleiben so — ein fester Tag in einem lebenden Dokument ist ein Befund |
| `harness/README.md`, `harness/sensors/*.md` | prüfen, update falls klein | Singleton-Abgleich gegen `README.template.md` und `sensors/gate.template.md` des neuen Stands (§6 *Sensor-Datei*) |
| `docs/plan/planning/open/`, `next/` (Folge-Slices, falls Befunde) | neu | per `cp` aus `slice.template.md` des **neuen** Stands, je Befund, der nicht klein ist; angelegt: `slice-harness-gate-index-werkzeug-teil` (Teil des Gate-Index, der einem Werkzeug gehört) |
| `docs/plan/adr/README.md` | update (ein Satz) | Abschnitt *Konventionen* folgt `adr/README.template.md` des neuen Stands: die ADR eines Gates schärft dessen Spec-Stelle |
| `docs/plan/planning/in-progress/roadmap.md` | update (Drift-Log) | eine Zeile für den angelegten Folge-Slice, wie für jeden wellenlos angelegten Slice |

**Ansatz:**

- Reihenfolge der Commits: (1) d-check umpinnen und dessen Befunde bereinigen —
  unabhängig von der Baseline, eigener Commit; (2) Baseline tauschen und Verweise
  nachziehen in einem Commit, damit `make baseline-verify` an keinem Commit zwei
  Tag-Verzeichnisse sieht; (3) Abgleich und Regelblock-Tabelle; (4) kleine Anpassungen
  und Folge-Slices. Der Vergleich der Bäume läuft vor (2) in einem Temp-Baum außerhalb
  des Repos (`diff -r` alt gegen neu), sein Ergebnis steht im Bericht.
- Der Abgleich geht nach Baseline-Regelwerk `modul-02-harness-bootstrap.md`
  §Freshness-Audit der vendored Baseline vor: Adaptions-Durchgang (auch aufgelöste
  Einträge), Form-Vergleich der Vorlagen (Singletons nacharbeiten, wiederkehrende
  Vorlagen append-only), Stichprobe gegen den Bestand (ein Abschnitt ohne Delta,
  rotierend).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug` liegt in `done/`
(WIP-Limit 1). Reihenfolge der wellenlosen Reihe nach Entscheidung des Nutzers vom
2026-10-06: `slice-harness-lint-werkzeug`, dieser Slice, `slice-harness-blackbox-kern`,
danach unverändert die übrige Reihe aus §4 von `slice-harness-blackbox-kern`. Erster
Schritt nach dem Start, vor dem ersten Commit an Pin, Baum oder Konventionen: Der
Architect entscheidet die Randformen aus §6
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`); der Slice legt keinen neuen
Vertrag an, aber er ändert die Grundlage zweier Gates (`make docs-check`,
`make baseline-verify`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Abgleich verlangt mehr
  als kleine Anpassungen — Maß: mehr als ein Satz je betroffener verkörperter Regel,
  eine Änderung an einem Werkzeugvertrag in der Spezifikation, eine neue oder
  ersetzende ADR, eine neue Adaption `MR-<NNN>`, ein Eingriff in ein Skript unter
  `tools/` oder mehr als eine Handvoll neuer d-check-Befunde in lebenden Dokumenten.
  Dann trägt dieser Slice nur Pin, Baum und Verweise (Liefer-Punkte 1 und 2) und den
  Abgleich als Befundliste; die Umsetzung geht an Folge-Slices mit Kennung.
- `in-progress` → `open` (blockiert — Carveout?): `make docs-check` mit `v0.82.0` wird
  durch Befunde rot, die sich nur durch eine Lockerung der Konfiguration oder eine
  Änderung eingefrorener Dokumente beheben ließen, und der Architect entscheidet keine
  der Varianten aus §6 *Befunde in eingefrorenen Dokumenten*; oder das Asset bzw. der
  Image-Digest ist nicht über einen Weg beschaffbar, den §6 *Netzzugriff* zulässt. Ein
  roter `docs-check` landet nicht in `done/` ohne Carveout (Modul 7); Ausweg ohne
  Carveout ist, nur die Baseline zu heben und d-check auf `v0.79.0` zu lassen, mit
  Folge-Slice für den Bump.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (darin `make docs-check` mit `v0.82.0` und
`make baseline-verify` mit `v6.16.0`), `grep` aus dem zweiten Liefer-Punkt leer,
Closure-Notiz mit Lerneintrag; jeder Befund des Abgleichs hat seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — **entschieden vom Architect am 2026-10-06**, nach
dem Start und vor dem ersten Commit. **Ort ist dieser Abschnitt**, nicht die
Spezifikation und keine ADR: Der Slice legt keinen Werkzeugvertrag an und ändert keinen —
kein neues `make`-Ziel, kein Skript unter `tools/`, die Rezepte von `d-check.mk` und die
Prüfung von `make baseline-verify` bleiben gleich (Abschnitt für Harness-Werkzeuge der
Spezifikation unberührt, Kopf bleibt `—`). Was hier entschieden ist, ist ein einmaliger
Wartungsschritt nach Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit
(„Wartung, kein Feedback-Gate“); den adoptierten Stand trägt `MR-000` in
`harness/conventions.md`, die Belege trägt §7. Verlangt die Umsetzung doch eine Änderung
an einem Werkzeugvertrag, gilt §4 *Rückführungen*. Was hier nicht steht, entscheidet der
Implementer nicht, er gibt es zurück (`.claude/commands/implement-slice.md`,
Randform-Rückgabe).

- **Alte Baseline** — `tools/harness/baseline-verify.sh` lässt genau ein Tag-Verzeichnis
  zu und endet mit zwei rot; `modul-02-harness-bootstrap.md` §Freshness-Audit lässt alte
  und neue Form zum Vergleich nebeneinander liegen. **Entscheidung: (a)** — `v6.13.0`
  fällt im selben Commit, der `v6.16.0` anlegt. Die Wirkung der Regel (die alte Form liegt
  während des Reviews als Vergleichsgrundlage vor) bleibt erhalten, nur der Ort wechselt:
  vor dem Commit `diff -r` im Temp-Baum außerhalb des Repos, danach
  `git diff -M <vorher> <nachher> -- .harness/baseline/` (Rename-Erkennung zeigt das Delta
  je Datei) und `git show <vorher>:.harness/baseline/v6.13.0/…`. Keine Adaption, kein
  Eingriff ins Skript; der Widerspruch zwischen Regelwerk-Text und emittiertem Skript wird
  in §7 nur festgehalten, ohne Frage an das Kurs-Repo (Entscheidung des Nutzers vom
  2026-10-06). Der Tausch ist **kein**
  `git mv` mit Inhaltsänderung im Sinn von `AGENTS.md` §3.3, sondern das Ersetzen eines
  vendored Baums: Ein reiner Move-Commit vorab hinterließe ein Verzeichnis `v6.16.0` mit
  dem Inhalt von `v6.13.0`, das `make baseline-verify` grün meldete. Ein Commit;
  `git log --follow` über den Baum ist ein akzeptiertes Negativ, weil kein Lauf die
  Historie einer vendored Datei verfolgt und `git diff -M` das Delta zeigt.
- **Was am alten Pfad hängt** (vom Architect ergänzt; Probe am 2026-10-06: HEAD-Kopie
  ohne `.harness/baseline/v6.13.0/`, d-check `v0.79.0`). Zwei Stellen brechen mit dem
  Entfernen, beide im Commit des Tauschs nachgezogen:
  (1) `.claude/rules/modul-01-entwicklungszyklus.md` und
  `.claude/rules/modul-05-planning-harness.md` sind Symlinks auf
  `../../.harness/baseline/v6.13.0/regelwerk/…`; sie zeigen danach auf dieselbe Datei
  unter `v6.16.0` (ist eine umbenannt oder entfallen: Rückgabe an den Architect). Kein
  Gate sieht einen hängenden Symlink, und `grep -r` folgt ihm nicht; Prüfung im
  zweiten Liefer-Punkt zusätzlich: `find . -path ./.git -prune -o -xtype l -print` ist
  leer. Dass ein Symlink den Tag fest nennt, ist ein akzeptiertes Negativ — ein Symlink
  kann kein `<tag>` tragen, und die Prüfung fängt ihn beim nächsten Tausch.
  (2) `harness/conventions/done/MR-001-spezifikations-ort-werkzeugvertraege.md` Zeile 21
  verlinkt `…/v6.13.0/regelwerk/grundlagen-referenz-richtung.md#spec-straten-…` — der
  einzige d-check-Befund der Probe (`target-missing`). Ausgang nach *Befunde in
  eingefrorenen Dokumenten* (a): Der Link wird zum Code-Span mit **unverändertem** Pfad
  samt `v6.13.0` und Anker, der Text bleibt; **nicht** auf `v6.16.0` umbiegen, das
  verwiese die damalige Aussage auf einen anderen Text. Eine Zeile in §7. Code-Spans mit
  dem alten Pfad in `done/`-Plänen und `docs/reviews/` meldet d-check nicht (Probe) und
  bleiben stehen.
- **`make baseline-verify` und der Tag im Verzeichnisnamen** — **Entscheidung:**
  (1) `SHA256SUMS` wird nach dem Entpacken **neu erzeugt**, im selben Container, in der
  Form, die das Skript liest und die `v6.13.0` hat: `<sha256>  <pfad>` relativ zum
  Tag-Verzeichnis, über jeden Nicht-Verzeichnis-Eintrag unter `regelwerk/` und
  `templates/`, nach `LC_ALL=C sort` geordnet. Vendored werden nur `regelwerk/` und
  `templates/`; andere Einträge auf oberster Ebene des Assets nennt der Bericht. Bringt
  das Asset selbst eine Prüfsummenliste mit, wird der entpackte Inhalt zusätzlich gegen
  sie geprüft (Ergebnis im Bericht); vendored wird sie nur, wenn sie unter `regelwerk/`
  oder `templates/` liegt und damit ohnehin gelistet ist. Vor dem Entpacken listet
  `unzip -l` die Einträge; ein absoluter Pfad, ein `..` oder ein Symlink im Asset ist
  Halt und Rückgabe. `SHA256SUMS` belegt danach nur, dass der Baum sich seit dem
  Vendoring nicht bewegt hat — die Herkunft hängt am sha256 des Assets (unten).
  (2) Verzeichnisname exakt `v6.16.0`. (3) Der sha256 des Assets steht im Bericht und in
  §7 *Belege*, in keinem lebenden Dokument. Eine `sources`-Prüfung von d-check bleibt
  außerhalb (neuer Sensor).
- **Neu erzeugen oder nur umpinnen** — **Entscheidung: (a)**, mit festgestelltem
  Ergebnis: `d-check --print-mk` von `v0.82.0` (lokal vorhanden, `--network none`) führt
  dieselben dreizehn Ziele wie das Fragment heute, `doc-immutable` und `doc-commits` mit
  Rezept, Rezepte unverändert. Neu erzeugt unterscheidet sich die Datei vom heutigen
  `d-check.mk` nur im Kopf und in den Pin-Zeilen; wieder angewendet werden genau die vier
  Adaptionen: Kopfkommentar von ai-harness-init, `doc-check` → `docs-check`, Muster von
  `doc-help` `^docs?-` statt `^doc-`, `DCHECK_DIGEST` gepinnt. Der Kopfkommentar nennt
  die `doc-help`-Adaption mit, er sagt heute „advisory doc-*-Targets verbatim“ und
  verschweigt sie (`AGENTS.md` §3.11). Digest: `docker image inspect
  ghcr.io/pt9912/d-check:v0.82.0 --format '{{json .RepoDigests}}'` nach `docker pull`
  meldet am 2026-10-06 `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`;
  der Implementer hält ihn gegen den Digest der Release-Notes von d-check `v0.82.0`
  (der erzeugte Kopf verweist dorthin). Weichen beide ab: Halt und Rückgabe. Nennen die
  Notes keinen, gilt der `RepoDigests`-Wert, und §7 sagt das.
- **Befunde in eingefrorenen Dokumenten** — **Entscheidung: wie empfohlen.** (a) für
  Form-Befunde, die das Dokument nicht umdeuten (Link-Ziel, Anker, Link → Code-Span mit
  gleichem Text), je Datei eine Zeile in §7; (c) für alles, was Inhalt ändern würde —
  dann bleibt d-check auf `v0.79.0`, Folge-Slice für den Bump. (b) ist ausgeschlossen:
  kein `scan.ignore`, kein Marker. Festgestellt: `v0.82.0` meldet am HEAD-Stand
  0 Befunde bei 279 Dateien (Probe in einer HEAD-Kopie); der Bump allein (Commit 1) ist
  grün, und der einzige erwartete Befund kommt aus dem Tausch (oben, MR-001).
- **Netzzugriff** — **Entscheidung:** Einmaliger Wartungsschritt, kein `make`-Ziel, kein
  Gate mit Netz. Image: `golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414`,
  der Pin der Stufe `deps` im `Dockerfile` — kein neuer Pin; es führt `wget` mit
  `ssl_client` und CA-Bündel, `sha256sum`, `unzip`, `find`, `sort` (geprüft mit
  `--network none`). Lauf mit `--user "$(id -u):$(id -g)"`, Mount nur eines
  Temp-Verzeichnisses außerhalb des Repos; ins Repo kopiert wird danach mit `cp` auf dem
  Host. **Soll-Wert des sha256, in dieser Rangfolge:** (1) eine Prüfsumme, die das
  Kurs-Repo selbst zum Release veröffentlicht (eigenes Asset oder Release-Text);
  (2) sonst das Feld `digest` des Assets in
  `https://api.github.com/repos/pt9912/ai-harness-course/releases/tags/v6.16.0` — von
  GitHub beim Hochladen berechnet, über einen anderen Abruf als der Download; (3) gibt es
  keines von beiden, gibt es keinen Soll-Wert. Was das heißt, ehrlich: Bei (2) belegt die
  Prüfung, dass die geladenen Bytes die hochgeladenen sind (Transport, CDN), nicht, wer
  sie hochgeladen hat; bei (3) ist es Trust-on-first-use — der Wert wird nur
  festgehalten, damit ein späteres Abweichen sichtbar wird, und §7 sagt „kein Soll-Wert“
  statt „geprüft“. Beides ist ein akzeptiertes Negativ für diesen Slice: Der Baum ist
  Text, der nicht ausgeführt wird, und dieser Slice liest sein Delta vollständig gegen
  den Bestand. Weicht bei (1) oder (2) der Wert ab: Halt, nichts entpacken, Rückgabe.
- **Präfix-Fehler von `make slice-mv`** — **Entscheidung: nicht in diesem Slice**, auch
  keine Prüfung des Bootstrap-Stands. Grund: Entscheidung des Nutzers vom 2026-10-06, der
  Fix kommt in ai-harness-init; die Übernahme eines neuen Werkzeug-Stands ist nach §1 ein
  anderer Vorgang. Kein Folge-Slice, kein Register-Eintrag — akzeptiertes Negativ, weil
  der Fehler laut scheitert, nicht still: Wer `slice-harness-lint` bewegt, bekommt die
  Mehrdeutigkeit gemeldet und nimmt den Handweg wie in `bf00052`.
- **Klarstellung zum Technik-Stratum** — **Entscheidung:** Kein Nachfolge-Eintrag
  `MR-<NNN>` in beiden Fällen, denn `MR-001` ist schon aufgelöst; ein Rückbau-Eintrag
  hätte nichts zurückzubauen. Bringt `v6.16.0` die Klarstellung, ergänzt die Zelle
  *aufgelöst durch* der Zeile `MR-001` in `harness/conventions.md` §Aufgelöste Adaptionen
  den Stand und den Abschnitt, der es trägt; die Datei unter
  `harness/conventions/done/` bleibt bis auf die Form-Reparatur oben unberührt. Bringt
  sie es nicht, bleibt alles stehen: Der Grund der Auflösung ist die Antwort des
  Kurs-Repos, nicht ein Baseline-Stand, und das Repo folgt der Lesart, die das Kurs-Repo
  selbst als seine erklärt hat — die `MR-000`-Aussage bleibt wahr, keine neue Adaption,
  kein Register-Eintrag (`BEO-REPO/harness-lesart-ohne-entscheidungsort` bleibt
  gestrichen). §7 nennt dann die offene Frage, mit welchem Stand die
  Baseline die Antwort trägt; der nächste Freshness-Audit liest die Zeile `MR-001` im
  Adaptions-Durchgang ohnehin.
- **Sensor-Datei** — **Entscheidung: wie formuliert bestätigt.** Eine neue Pflicht-Sektion
  in `harness/sensors/gate.template.md` oder in `grundlagen-harness-dateien.md`
  §Einstiegspunkt zieht die Sensor-Dateien nach, wenn es je Datei beim Maß aus §4 bleibt,
  sonst Folge-Slice; eine geänderte Regel, *welche* Gates eine Sensor-Datei brauchen,
  geht als Nachzug an die Pläne der Lint-Reihe (`AGENTS.md` §3.9).

**Risiken:**

- **Strengere d-check-Version färbt viele Dateien rot** — `v0.82.0` kann neue Module oder
  strengere Regeln mitbringen; bei 240 und mehr gescannten Dateien kann die Zahl der
  Befunde den Slice sprengen. — **Ausgang:** entfallen: `v0.82.0` meldet am Stand vor
  dem Tausch 0 Befunde (Probe `m9` der Verifikation), nach dem Tausch genau den einen
  aus §6 *Was am alten Pfad hängt*, behoben nach (a); der Bump ist vollzogen, für ihn
  kann das Risiko nicht mehr eintreten.
- **Konfigurations-Drift in `.d-check.yml`** — eine neue Version liest einen Schlüssel
  anders oder ignoriert ihn still; ein grünes `docs-check` behauptete dann eine Prüfung,
  die nicht läuft (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`). Gegenmittel:
  je aktives Modul ein Probe-Befund in einem Temp-Baum, rot gesehen, vor und nach dem
  Bump. — **Ausgang:** entfallen: `.d-check.yml` ist unverändert (`git diff 53ea519
  481ef52 -- .d-check.yml` leer), und die Probe je aktives Modul (`links`, `anchors`,
  `ids`, `matrix`, `spans`) war unter `v0.79.0` und `v0.82.0` je einmal rot (§7
  *Belege*; `links`, `anchors`, `ids` von der Verifikation unabhängig rot gesehen,
  Mutation `m7`). Für eine spätere Version zählt
  `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` weiter.
- **Verkörperte Regel und neue Baseline widersprechen sich** — eine Regel aus
  `AGENTS.md` §3.9 bis §3.12, aus den Commands oder dem Register-README regelt der neue
  Stand anders; die `MR-000`-Aussage „keine inhaltlichen Adaptionen“ wäre dann falsch.
  — **Ausgang:** eingetreten, an den Folge-Slice `slice-harness-gate-index-werkzeug-teil`:
  `harness/README.md` §Sensors führt die Targets der Werkzeug-Fragmente, die `v6.16.0`
  einem Teil des Werkzeugs zuweist; die `MR-000`-Aussage ruht bis zu dessen Start auf
  der Lesart, die dort unter *Herkunft* steht (Entscheidung des Nutzers vom
  2026-10-06). Der zweite Widerspruch, `AGENTS.md` §3.12 mit der ADR als Ort einer
  Randform (Review F-434), ist in `481ef52` behoben und braucht keinen Ausgang mehr.
- **Lebender Verweis auf den alten Baum bleibt stehen** — ein Pfad unter
  `.harness/baseline/v6.13.0/` in einem lebenden Dokument bricht mit dem Entfernen
  (`anchor-missing` bzw. Link-Befund) oder bleibt in Code-Spans unentdeckt, die kein Gate
  liest. Gegenmittel: `grep` aus dem zweiten Liefer-Punkt. — **Ausgang:** entfallen:
  Der `grep` aus dem zweiten Liefer-Punkt ist leer, `find … -xtype l` ebenso, beide
  Symlinks unter `.claude/rules/` lösen auf `v6.16.0` auf, `MR-001` trägt den alten
  Pfad als Code-Span (Verifikation, Punkt 2 und Mutation `m5`); die Nennungen im
  eigenen Plan sind Wort, kein Verweis (V-78).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Bericht des Implementers:** Diese Sektion ist der Bericht des Implementers, den die
  DoD in Punkt 1 bis 3 nennt (Digest-Herkunft, Asset-sha256, Liste der Blöcke); einen
  eigenen Bericht außerhalb des Plans gab es nicht. Die Verifikation hat jede Aussage
  unter *Belege* selbst nachgefahren (V-79).
- **Was hat funktioniert:** Die Randformen standen vor dem ersten Commit an Pin, Baum oder
  Konventionen fest (`53ea519` vor `01f9933`), und die Commit-Folge (1) bis (4) aus §3
  *Ansatz* hielt: An keinem Stand lagen zwei Tag-Verzeichnisse (Verifikation, Punkt 2),
  `git show -M` erkannte den Tausch als Rename (Review F-437). Die Probe in einer
  HEAD-Kopie vor dem Tausch sagte den einzigen d-check-Befund voraus (`MR-001`,
  `target-missing`), und er kam genau so (Mutation `m8`). Keine Rückführung aus §4: kein
  Werkzeugvertrag, keine ADR, keine neue Adaption, kein Skript unter `tools/` im Diff.
- **Was ging anders als geplant:** Der Abgleich hielt das Delta von
  `grundlagen-referenz-richtung.md` §Spec-Straten gegen Spezifikation, Sensor-Dateien und
  Pläne, aber nicht gegen `AGENTS.md` §3.12, das die ADR weiter als Ort einer Randform
  zuließ; das Review fand es (F-434, MEDIUM), `481ef52` zog §3.12 nach, ebenso F-435,
  F-436 und F-439. Die Verifikation fand §7 an zwei Stellen unvollständig: eine falsche
  Begründung und ein Teilbefund ohne Ausgang bei Befund (2) (V-76), zwei Entscheidungen
  des Nutzers nicht festgehalten (V-77). Beides ist mit dieser Closure nachgezogen.
  V-78 (Urteil zu F-440): kein Befund, erledigt — die Nennungen von
  `.harness/baseline/v6.13.0/` im eigenen Plan sind Wort, kein Verweis. Entscheidungen
  des Nutzers vom 2026-10-06, hier festgehalten:
  - (a) `slice-harness-gate-index-werkzeug-teil` bleibt ein Slice und wird kein
    Register-Eintrag, auch nach F-438; eingetragen in dessen §1 *Herkunft* (V-77).
  - (b) Kein Change Request an das Kurs-Repo, weder zu `baseline-verify` und zwei
    Tag-Verzeichnissen noch zur Lesart unter `MR-000` (F-438).
  - (c) Die angenommenen Gate-ADRs bleiben unverändert (`AGENTS.md` §3.5); Begründung
    unter *Belege*, Befund (2).
  - (d) Die älteren Werkzeuge (Commit-Hook, `baseline-verify`, `a-check-negativ`,
    `hook-gegenprobe`) bekommen ihre Festlegungen in Spezifikation §11 erst, wenn für
    eines von ihnen eine Randform zu entscheiden ist.
  - **Summary-Zeilen:** Review `docs/reviews/2026-10-06-review-slice-harness-upgrade-v6-16.md`:
    „HIGH 0, MEDIUM 1, LOW 2, INFO 4. Finding-Klassen dieses Laufs: Verkörperte Regel nach
    Baseline-Sprung nicht nachgezogen · Regelblock-Zelle trägt zwei Werte · Teilersetzung
    lässt Satzrest stehen.“ Verifikation
    `docs/reviews/2026-10-06-verifikation-slice-harness-upgrade-v6-16.md`, Urteil: Liefer-Punkt
    1 und 2 erfüllt, Punkt 3 erfüllt bis auf V-76 und V-77 (mit dieser Closure
    nachgezogen); V-78 kein Befund, V-79 Hinweis; `make gates` grün am Stand `481ef52`.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke: Spezifikation §11 *Harness-Werkzeuge*
  führt Festlegungen nur für `kopf-check`, `abdeckung` und `lint` (`SPEC-047` bis
  `SPEC-049`). Für den Commit-Hook ([ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md),
  [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md)), `baseline-verify`,
  `a-check-negativ` und `hook-gegenprobe` stehen die Regeln in ADR, Skriptkopf oder
  `harness/README.md` — an Orten, die `v6.16.0` für Festlegungen eines Werkzeugs
  ausschließt. Die Lücke schließt sich je Werkzeug mit seiner ersten zu entscheidenden
  Randform (Entscheidung (d)); bis dahin zählt sie im Register unter
  `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`. Kein `liegt in`: Mit diesem
  Eintrag ist nichts verkörpert. Retirement-Check von `AGENTS.md` §3.9 bis §3.12: keine
  Regel gelockert oder entfernt; §3.12 ist enger geworden (F-434).
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `2fc5f77`.
  - `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen/` — **neu**,
    Beleg (F-434, Finding-Klasse des Reviews), **1×**. Erster Baseline-Sprung des Repos;
    kein älterer Eintrag trägt die Klasse.
  - `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum/` — **neu**, Beleg (V-76,
    Teilbefund von Befund (2)), **1×**. Nicht `harness-lesart-ohne-entscheidungsort`
    (gestrichen, bleibt 2×): Dessen Ursache, der fehlende Ort, ist weiterhin weg; hier
    liegt Bestand aus der Zeit davor am alten Ort, und eine neue Randform ginge nach
    §3.12 in die Spezifikation.
  - `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht/` — Beleg (V-79), **2×**. Die
    DoD verlangte Belege „im Bericht“; der Verifikation lag keiner vor, sie hat §7
    *Belege* selbst nachgefahren. Erreicht die Schwelle nicht; der nächste Beleg macht
    ihn zur Lücke, die einen Folge-Slice braucht.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11
    seit welle-extended-query) — Beleg (Kopf von `d-check.mk` sagte „verbatim“ über ein
    angepasstes `doc-help`, vom Architect in §6 gefunden, behoben in `01f9933`, Mutation
    `m6`), **15×**.
  - Ohne Beleg: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (bleibt 2×; die
    Probe je Modul war rot, `.d-check.yml` unverändert),
    `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×; keine Gate-Regel geändert),
    `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (bleibt 1×; §4 nannte den
    Architect, er entschied in `53ea519` vor dem ersten Commit),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 2×; keine
    Rückgabe), `BEO-REPO/plan-folgt-korrektur-nicht` (bleibt 13×; `481ef52` zog die
    Zeile `AGENTS.md` in §3 nach; V-76 und V-77 sind eine falsche Begründung und eine
    fehlende Entscheidung in §7, kein Plan auf altem Stand),
    `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` (bleibt 2×; der Präfix-Fehler
    von `make slice-mv` ist nach §6 nicht eingetragen).

  Einmalig und nicht eingetragen: F-435 (Regelblock-Zelle trägt zwei Werte) und F-436
  (Teilersetzung lässt Satzrest stehen), beide in `481ef52` behoben; F-437 und F-438
  (Hinweise an den Architect, F-438 durch Entscheidung (a) erledigt), F-439 (behoben),
  F-440 (V-78). Kein Eintrag erreicht mit diesem Slice die Schwelle 3× neu; über ihr
  stehen nur verkörperte Einträge, ihnen gibt diese Closure keinen Ausgang, den
  Lese-Schritt führt die nächste Welle-Closure.
- **Folge-Slices:** `slice-harness-gate-index-werkzeug-teil` (`open/`, Befund (3) und das
  Risiko *Verkörperte Regel und neue Baseline widersprechen sich*); als nächster in der
  wellenlosen Reihe `slice-harness-blackbox-kern`.
- **Risiken aus §6:** vier. *Verkörperte Regel und neue Baseline widersprechen sich*:
  **eingetreten**, an `slice-harness-gate-index-werkzeug-teil`. *Strengere
  d-check-Version färbt viele Dateien rot*, *Konfigurations-Drift in `.d-check.yml`* und
  *Lebender Verweis auf den alten Baum bleibt stehen*: **entfallen**, mit Begründung. Die
  Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker — kein `liegt in`, der Lerneintrag ist eine benannte
  Spec-Lücke, nichts ist verkörpert;
  `grep -rn "seit slice-harness-upgrade-v6-16"` über `AGENTS.md`, `.claude/` und
  `harness/` ist leer, wie es sein soll. Folge-Slice — `slice-harness-gate-index-werkzeug-teil`
  liegt in `open/`, ist nicht geschlossen, nennt in §1 *Herkunft* diesen Slice und
  Entscheidung (a) und schließt den Teil des Gate-Index nicht aus. Register — die vier
  Kennungen mit Beleg und die sechs ohne bestehen als Verzeichnis, jedes mit nicht leerem
  `evidence/`, der Beleg dieses Slice heißt `evidence/slice-harness-upgrade-v6-16.md`. Die
  nächste Welle-Closure prüft erneut.
- **Belege:**
  - *Asset:* `lab-regelwerk.zip` des Tags `v6.16.0`, sha256
    `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063`, geprüft gegen
    Rang (1) aus §6 *Netzzugriff*: das Asset `SHA256SUMS` desselben Releases, vom
    Kurs-Repo veröffentlicht (`sha256sum -c`: OK); das Feld `digest` des Assets in der
    Release-API (Rang 2) nennt denselben Wert. Abruf einmalig im Image der Stufe `deps`
    mit `--user`, Mount nur eines Temp-Verzeichnisses. Vor dem Entpacken: 68 Einträge,
    oberste Ebene nur `regelwerk/` und `templates/`, kein absoluter Pfad, kein `..`,
    kein Symlink (Unix-Modus je Eintrag aus dem Central Directory: 54 × `100644`,
    14 × `40755`); das Asset bringt keine eigene Prüfsummenliste mit. `SHA256SUMS` nach
    dem Entpacken neu erzeugt, 54 Zeilen; `make baseline-verify` meldet `v6.16.0 OK`.
  - *d-check:* `docker image inspect ghcr.io/pt9912/d-check:v0.82.0 --format
    '{{json .RepoDigests}}'` nach `docker pull` meldet
    `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`; die
    Release-Notes von d-check `v0.82.0` nennen denselben Digest-Pin. Neu erzeugtes
    Fragment gegen das alte: Unterschied nur in Kopf und Pin-Zeilen. Probe je aktives
    Modul (`links`, `anchors`, `ids`, `matrix`, `spans`) in einer HEAD-Kopie: unter
    `v0.79.0` und `v0.82.0` je ein Befund und Exit 1.
  - *`diff -r` der Bäume* (Temp-Baum, vor dem Tausch): dieselben 54 Dateien, keine
    neue, umbenannte oder entfallene; alle 26 Regelwerk-Dateien und 14 Vorlagen
    geändert. In 18 Regelwerk-Dateien nur die Quell-Zeile mit dem Tag; Inhalt in
    `grundlagen-begriffe.md`, `grundlagen-harness-dateien.md`,
    `grundlagen-referenz-richtung.md`, `modul-03-spec.md`, `modul-10-review-harness.md`,
    `modul-13-quality-gates.md`, `modul-15-observability.md` und in `README.md`
    (Stand-Zeile „Kurs-Welle 159 · 2026-10-06“).
  - *d-check-Befunde:* `v0.82.0` am Stand vor dem Tausch 0; nach dem Tausch 1
    (`target-missing` in
    `harness/conventions/done/MR-001-spezifikations-ort-werkzeugvertraege.md`),
    behoben nach §6 *Befunde in eingefrorenen Dokumenten* (a): Link zu Code-Span mit
    unverändertem Pfad samt `v6.13.0` und Anker, Text unverändert. Danach 0.
  - *Abgleich v6.13.0 → v6.16.0, Befunde und Ausgänge:* Regelblöcke: keiner neu,
    umbenannt oder entfallen; die Tabelle in `harness/conventions.md` führt alle 26
    Dateien; `grundlagen-harness-dateien.md` steht in zwei Zeilen je Abschnitt, der Teil des
    Gate-Index eines Werkzeugs mit *kommt nicht mit*, Grund und Dauer, die übrigen Werte
    bleiben. `MR-000` bleibt wahr (keine aktive Adaption; die
    aufgelöste `MR-001` ist nach §6 *Klarstellung zum Technik-Stratum* behandelt).
    (1) Reviewer ohne Stil-Finding, HIGH und MEDIUM nur mit Failure-Szenario, LOW mit
    Konventions-Anker, `pfad` als Kurzzitat (`modul-10-review-harness.md`, Vorlagen der
    beiden Skills) — behoben in `.harness/skills/reviewer.md` und
    `.harness/skills/closure-note-reviewer.md`. (2) Werkzeug-Festlegungen im
    Technik-Stratum (`grundlagen-referenz-richtung.md` §Spec-Straten, `modul-03-spec.md`,
    Vorlagen der Spezifikation, der ADR, des ADR-Index und der Sensor-Datei) — behoben:
    Zeile `MR-001` in §Aufgelöste Adaptionen nennt Stand und Abschnitt, der ADR-Index
    nennt die Regel in *Konventionen*, `AGENTS.md` §3.12 nennt die Spezifikation als Ort
    einer entschiedenen Randform und die ADR nur mit Entscheidung, Gründen und
    `Schärft:` (Review F-434; `.claude/agents/architect.md` und
    `.claude/commands/implement-slice.md` zeigen auf §3.12 und wiederholen den Ort nicht); Spezifikation §11, beide Sensor-Dateien und die
    Pläne der Lint-Reihe folgen ihr schon. Die angenommenen Gate-ADRs bleiben unverändert,
    weil der Nutzer es am 2026-10-06 entschieden hat (Entscheidung (c)) und eine
    angenommene ADR nicht überschrieben wird (`AGENTS.md` §3.5). Ihr Kopf sieht so aus:
    [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md),
    [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md),
    [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) und
    [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) tragen `Schärft: —`,
    [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) steht in
    ihrem `Bezug:`; für die zweite und die vierte gäbe es mit `SPEC-048` und `SPEC-047` eine
    Spec-Stelle, die sie nach der Regel im ADR-Index schärfen würden.
    [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) trägt `Schärft:
    SPEC-049` und folgt der Regel. Teilbefund: Für Commit-Hook, `baseline-verify`,
    `a-check-negativ` und `hook-gegenprobe` hat Spezifikation §11 keine Festlegung, ihre
    Regeln stehen in ADR, Skriptkopf oder `harness/README.md` — Ausgang: Eintrag im
    Beobachtungs-Register, `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`;
    Grund: Die Festlegungen kommen nach §11, wenn für das Werkzeug eine Randform zu
    entscheiden ist (Entscheidung (d)), und kein Slice dafür ist geplant. Derselbe
    Eintrag trägt das fehlende `Schärft:` dieser beiden (V-76). (3) Teil des Gate-Index,
    der einem Werkzeug gehört (`grundlagen-harness-dateien.md`, `modul-13-quality-gates.md`,
    Vorlagen von `AGENTS.md`, `harness/README.md`, `Makefile`, `.d-check.yml`) —
    Folge-Slice `slice-harness-gate-index-werkzeug-teil` in `open/`. (4) Form der
    Register-Kennung und Ablage `observations/` in den Vorlagen — kein Befund: lebende
    Harness-Dateien führen die neue Form schon, Pläne sind wiederkehrende Vorlagen
    (append-only). (5) Pflicht-Feld ohne Wert in `modul-15-observability.md` — kein
    Befund: `harness/erfassung-feldliste.md` (werkzeug-erzeugt) sagt „leer heißt
    unbekannt“ und nennt die Quelle. Stichprobe ohne Delta: `modul-05-planning-harness.md`
    §Offene Risiken gegen §6 dieses Plans und die offenen Pläne — ohne Befund.
- **Festgehalten, nicht weitergereicht:** (1) `modul-02-harness-bootstrap.md`
  §Freshness-Audit lässt alte und neue Form nebeneinander liegen („Das alte Verzeichnis
  fällt erst, wenn der Review durch ist“), das emittierte `tools/harness/baseline-verify.sh`
  endet bei zwei Tag-Verzeichnissen rot. Aufgelöst durch das Entfernen von `v6.13.0` im
  Commit, der `v6.16.0` anlegt (§6 *Alte Baseline*); kein Change Request an das Kurs-Repo
  (Entscheidung (b) des Nutzers vom 2026-10-06). (2) Die Antwort zum Technik-Stratum trägt die Baseline seit `v6.16.0` in
  `grundlagen-referenz-richtung.md` §Spec-Straten; keine offene Frage mehr.

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
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie. Der vendored Baum unter
`.harness/baseline/` ist keine eigene Sub-Area: Er ist Kurs-Inhalt, nicht
autoritativ und vom Doku-Gate ausgenommen (`harness/conventions.md` §Was der
mitgelieferte Baum ist).

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `8d9f0ee` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×, offen) — trifft den
  d-check-Bump unmittelbar: eine neue Version kann `.d-check.yml` anders lesen; Risiko
  in §6 mit Probe-Befund je Modul.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×, offen) — trifft eine Änderung an
  `.d-check.yml`, falls `v0.82.0` sie verlangt: eine Regel, die eine andere ersetzt
  statt ergänzt, nimmt still Prüfung weg; im Review zu sehen, gegen §6 *Befunde in
  eingefrorenen Dokumenten* (b).
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (2×, gestrichen) — Grund der
  Streichung ist die Antwort des Kurs-Repos; trägt `v6.16.0` sie nicht, fehlt der
  Streichung die Baseline (§6 *Klarstellung zum Technik-Stratum*). Ein drittes
  Auftreten wäre dann keine Notiz mehr.
- `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` (2×, behoben) — verwandt mit dem
  Präfix-Fehler von `make slice-mv` (beide: ein tool-eigenes Werkzeug passt nicht zur
  benannten Kennung des Repos), aber nicht dieselbe Klasse: dort lehnt der Hook ab, hier
  findet das Werkzeug den Plan mehrdeutig. Wird der Präfix-Fehler eingetragen (§6), dann
  als eigener Eintrag, nicht als dritter Beleg hier — sonst stünde er über der Schwelle
  ohne gleiche Ursache.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — darum stehen die
  Randformen in §6 offen und §4 *Start* nennt den Architect-Schritt; ein drittes
  Auftreten des zweiten Eintrags in diesem Slice erreichte die Schwelle.
- `BEO-REPO/plan-folgt-korrektur-nicht` (12×, verkörpert in `AGENTS.md` §3.9, Sensor
  `make kopf-check`) — ändert der Abgleich eine Regel, auf die offene Pläne bauen,
  ziehen deren Pläne nach; die Verweise auf `.harness/baseline/v6.13.0/` in offenen
  Plänen fängt der `grep` aus dem zweiten Liefer-Punkt.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (13×, verkörpert in §3.11) — die
  Zeile von `make baseline-verify` in `harness/README.md` §Sensors und der Kopf von
  `d-check.mk` sagen nur zu, was die neue Version prüft.

Keiner der offenen Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
