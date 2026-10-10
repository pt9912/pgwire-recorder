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
von der nächsten Welle-Closure. Vorgezogen nach Entscheidung des Nutzers vom 2026-10-10: nach
`slice-v1-abschluss-einspielen-laufsteuerung` und vor `slice-v1-abschluss-einspielen-extended`
(WIP-Limit 1); welle-v1-abschluss geht danach in der Reihenfolge ihres §5 weiter.

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

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen-laufsteuerung` liegt in
`done/` (WIP-Limit 1). Erster Schritt nach dem Start, vor dem ersten Commit an Baum oder
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

**Randformen** (`AGENTS.md` §3.12) — **offen, der Architect entscheidet sie nach dem Start und
vor dem ersten Commit** (§4). Der Slice legt keinen Werkzeugvertrag an und ändert keinen; der
Ort der Entscheidung ist deshalb, wie in `slice-harness-upgrade-v6-16` (dort §6), dieser
Abschnitt. Die Entscheidungen jenes Slice sind der Vorschlag, den der Architect gegen den neuen
Stand prüft. Was hier nicht steht, entscheidet der Implementer nicht, er gibt es zurück
(`.claude/commands/implement-slice.md`, Randform-Rückgabe).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

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
