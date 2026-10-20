# Slice slice-v1-abschluss-binary-namen: Release-Binaries ohne Version im Dateinamen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität)

**Berührte Spec-Stellen:** `SPEC-055` (Punkt 2 und 3, [`harness/vertraege/`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md)) · `SPEC-067` (Punkt 4 und 19 (3), (4), (7), (9), [`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md)) · `LH-FA-01` · `LH-FA-19.a` · `LH-QA-03`

**Entscheidung:** [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-20: keine Version im Dateinamen der Binaries, Schema `pgwire-recorder-<os>-<arch>`)

**Verantwortlich:** pt9912
**Autor:** Architect. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make binaries`, `make test-container` und die Pipeline (`release.sh`) schreiben, erwarten und prüfen die Binaries als `pgwire-recorder-<os>-<arch>` (unter Windows `.exe`) ohne Version im Namen; das Image behält den Tag mit Version. Schema, Dateiliste, Zeilen von `SHA256SUMS`, Eingaben und Fälle entscheidet [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 2 und 3 und [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 4; warum so, [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md). Der Slice liegt vor `slice-v1-abschluss-homebrew`, dessen Vorlage und Fälle die Namen brauchen (Bindung unten).

**Sendung an `slice-v1-abschluss-homebrew`** (`AGENTS.md` §3.13, Bindung): Die Vorlage der Formel, die Prüfung der Summen (`SPEC-067` Punkt 19 (3)) und ihre Fälle `p19-*` rechnen mit den Namen `pgwire-recorder-<os>-<arch>`, die dieser Slice einführt; der Fall `p19-summen-altes-schema` (Namen des alten Schemas, Exit 1 `summen`) liegt dort. Beides steht dort in §1 und §6 mit der Kennung dieses Slice und als Start-Trigger.

**Sendung an `slice-erster-release-homebrew-nachweis`** (Bindung): Der Nachtrag der Formel für `v0.1.0` entfällt; die erste Formel entsteht durch den Job `homebrew` mit dem nächsten stabilen Tag, dessen Binaries das neue Schema tragen. Das steht dort in §1, DoD, §4 und §6 mit der Kennung dieses Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Vorlage der Formel, der Schritt und der Job `homebrew` — `slice-v1-abschluss-homebrew` (Klasse 1; dort §1 und §6 nehmen die Namen an). Grund: ein eigener Liefergegenstand mit eigener Gegenprobe; dieser Slice ändert nur Namen und ihre Prüfung.
- Ein neuer Tag und ein neuer Release (der erste mit dem neuen Schema) — Handlung des Nutzers mit neuem Ja zum Tag (Klasse 3, anderer Vorgang). Den Tag setzt der Nutzer; der Slice, der ihn vorbereitet und belegt, ist `slice-v1-abschluss-release-v0-1-1` (Welle `welle-v1-nachtrag`, Reihenfolge 4), und `slice-erster-release-homebrew-nachweis` wartet auf ihn (dort §4).
- Der Abschnitt „Binary vom GitHub-Release“ im Handbuch und die Zeile zu `v0.1.0` in `README.md` — Bestand bleibt bewusst stehen (Klasse 2). Grund: Sie beschreiben, was `v0.1.0` trägt (Namen mit Version, einmal gelaufen); das bleibt wahr, bis ein neuer Tag erscheint. `slice-v1-abschluss-release-v0-1-1` schreibt sie nach `SPEC-067` Punkt 18 um (dort §1 und DoD 2 nehmen die Sendung an); dieser Slice ändert nur den Abschnitt zu `make binaries`.
- Der Release `v0.1.0` selbst — unverändert (Entscheidung des Nutzers vom 2026-10-20); ein Tag wird nie verändert (`SPEC-067` Punkt 11).
- Das Image, seine Tags und die Registries — Bestand bleibt (Klasse 2); die Entscheidung gilt nur den Binaries.
- Änderung des Lastenhefts — es legt kein Namensschema fest (`LH-FA-01`, `LH-QA-03` nennen Plattformen); eine Anhebung der Version nach `harness/conventions.md` entfällt.
- Produkt-Code (`cmd/`, `internal/`) und `Dockerfile` — Schicht-Abgrenzung (Klasse 4): `VERSION` geht weiter in `-X main.version`; der Name entsteht allein in `binaries.sh`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung): `tools/harness/binaries.sh`, `run-container.sh` und `release.sh` (Dateiliste des Schritts `binaries`, Summe je Plattform im Schritt `image`) schreiben und erwarten `pgwire-recorder-<os>-<arch>` ohne Version, `.exe` allein bei `windows`; `make test-container-gegenprobe` und `make release-gegenprobe` tragen die Fälle und Mutationen aus `SPEC-055` Punkt 3 (`p3-namen`, `p3-exe`, `p3-name-ohne-version`, `p3-name-trenner`, `p3-sums-zeile-<os>-<arch>`, `p3-dist-andere-version`, `p3-dist-altes-schema`; die Mutation `p3-name` kehrt sich um) und `SPEC-067` Punkt 4 (`p4-altes-schema`, `p4-namen-je-plattform`) und führen jeden alten Namen mit Version nach; je Zusage Mutation und roter Fall im Bericht (`AGENTS.md` §3.10); die Zeile von `make binaries` in `harness/README.md` §Sensors nennt das neue Schema.
- [ ] Benutzerhandbuch (Abschnitt zu `make binaries`: Namen, Beispiele mit `dev`, Aufruf des Binaries) und `docs/maintainer/releasing.md` (§2 *Was lokal läuft*, §3 *Versionierung*) beschreiben das Ausgeführte im Ist-Zustand ohne Chronik (`AGENTS.md` §3.11); der Abschnitt „Binary vom GitHub-Release“ und die Zeile zu `v0.1.0` in `README.md` bleiben (§1); die Beispiele des Handbuchs, die Namen nennen, liegen als Datei vor und laufen gegen das gebaute Binary, soweit sie es aufrufen (`dist/pgwire-recorder-linux-amd64 version`).
- [ ] `make gates` grün; `make test-container-gegenprobe` und `make release-gegenprobe` am Endstand (`AGENTS.md` §6).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/binaries.sh` | update | Name `pgwire-recorder-${os}-${arch}` (Marke `P3name`), Meldung bei `DIST` unverändert |
| `tools/harness/run-container.sh` | update | `binaer()` nennt den neuen Namen (Punkt 4 und 6 von `SPEC-055`) |
| `tools/harness/release.sh` | update | `erwartete_namen` ohne Version; Summe im Schritt `image` aus `pgwire-recorder-linux-<arch>` |
| `tools/harness/container-gegenprobe.sh` | update | `bname`, `namen_soll`, Muster der Summen-Zeile, Meldung zu Punkt 4; neue Fälle und Mutationen der Namen und von `DIST` |
| `tools/harness/release-gegenprobe.sh` | update | Attrappe von `make`, `NAMEN`, `ALLE7`, Fälle zu Punkt 5 und 9; Fälle `p4-altes-schema`, `p4-namen-je-plattform` und ihre Mutationen |
| `harness/README.md` (§Sensors, Zeile `make binaries`) | update | nennt das Schema |
| `docs/user/benutzerhandbuch.md` (Abschnitt zu `make binaries`) | update | Namen und Beispiele im Ist-Zustand |
| `docs/maintainer/releasing.md` | update | §2, §3 (Namen tragen keine Version mehr) |

Nicht in dieser Liste: `Dockerfile`, `harness/mk/container.mk` und `cmd/`, `internal/` (§1), der Abschnitt zu `v0.1.0` im Handbuch und `README.md` (§1), die Verträge `SPEC-055` und `SPEC-067` (liegen vor dem Code, Architect).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-20), und `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`). Der Slice steht in der Reihenfolge der Welle vor `slice-v1-abschluss-homebrew` und blockiert ihn (dessen Start-Trigger nennt ihn).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 3d927615..<letzter Commit des Slice> -- tools harness/README.md docs/user docs/maintainer README.md` (Tests eingeschlossen, `AGENTS.md` §3.14) **über 800** liegt. Schätzung des Architect in dieser Metrik: Produktcode rund 30 (die drei Skripte, je ersetzte Zeile zwei), Doku rund 80 (Handbuch, `releasing.md`, `harness/README.md`), die Gegenprobe rund 250 bis 400 (viele ersetzte Namen, dazu neue Fälle und Mutationen; Faktor von Tabellen- und Gegenprobe-Code zum Produktcode rund zwei und mehr), zusammen 360 bis 510. Kein vorgemerkter Schnitt: Der Slice hat zwei Liefer-Punkte, und die Obergrenze liegt unter der Marke.
- `in-progress` → `open`: Der Nutzer nimmt [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) nicht an oder wählt ein anderes Schema — zurück in die Planung; die Verträge `SPEC-055` und `SPEC-067` ziehen mit der Entscheidung nach (Architect).

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

**Randformen des Vertrags** (`AGENTS.md` §3.12; alle vor dem Code entschieden, jede an der genannten Stelle der Verträge):

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Schema | `pgwire-recorder-<os>-<arch>`, Bindestriche, klein; `.exe` allein bei `windows`; sechs Namen | `SPEC-055` Punkt 3 |
| (b) Woher die Version, wenn der Name sie nicht trägt | aus `VERSION` wie bisher (`-X main.version`); nach dem Lauf steht sie nur im Binary (`version`), eine Versionsdatei in `DIST` entsteht nicht; die Pipeline bindet sie durch den Aufruf `VERSION=<Version>` und für `linux/amd64` durch Summe im Image und Smoke | `SPEC-055` Punkt 2 und 3; `SPEC-067` Punkt 4 |
| (c) Zeilen von `SHA256SUMS` | `<Summe>  <Name>`, sechs Zeilen, C-Locale sortiert, LF; Namen nach (a) | `SPEC-055` Punkt 3 |
| (d) Mehrere Versionen im selben `DIST` | `DIST` muss fehlen oder leer sein; Dateien einer anderen `VERSION` oder des alten Schemas sind ein nicht leeres `DIST`, rot vor dem ersten `docker` | `SPEC-055` Punkt 3 (`p3-dist-andere-version`, `p3-dist-altes-schema`) |
| (e) Eingaben und Fehler | das Muster von `VERSION` bleibt (Wert geht weiter in `-X` und in die Befehlszeile); `dev` am Ziel `test-container` bleibt rot; keine Lockerung eines Gates | `SPEC-055` Punkt 2 |
| (f) Release-Pipeline | Dateiliste der sechs und `SHA256SUMS`, der alte Name ist „zu viel“ und „fehlt“; sieben Anhänge bleiben; Summe im Image aus `pgwire-recorder-linux-<arch>` | `SPEC-067` Punkt 4, 5 und 9 |
| (g) Altes Schema (`v0.1.0`) | wird nicht gelesen und nicht übersetzt; der Schritt `homebrew` endet bei diesen Namen mit Exit 1 `summen` (Fall in `slice-v1-abschluss-homebrew`) | `SPEC-067` Punkt 19 (3), (7) |
| (h) Plattform | Namen aus Kleinbuchstaben, Ziffern, `-`, `.`: gültig auf Linux, macOS, Windows; erzeugt und geprüft nur unter Linux im Docker, macOS und Windows Grenze | `SPEC-055` Punkt 3 und 13 (4) |
| (i) Aufrufarten, Abbruch | unverändert; die Rezepte und die Fälle zu `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`, `SIGINT`, `SIGTERM` bleiben, nur ihre Namen ziehen nach | `SPEC-055` Punkt 9 und 11 |

- Fundstellen des alten Schemas, die dieser Plan nicht nennt (`grep -rn 'pgwire-recorder_' --exclude-dir=.git --exclude-dir=baseline`, ohne `docs/reviews/` und `done/`; Stand beim Schreiben: die fünf Skripte, die zwei Verträge, `harness/README.md`, Handbuch, `releasing.md`, `README.md`, `slice-v1-abschluss-homebrew`) — **Ausgang:** offen bis Closure (`grep` am Endstand ohne Treffer außer den Zeilen, die das alte Schema als alt nennen).
- Der Abschnitt zu `v0.1.0` in Handbuch und `README.md` nennt Namen, die der nächste Tag ablöst; wer ihn umschreibt, ist `slice-v1-abschluss-release-v0-1-1` (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`) — **Ausgang:** weiter offen, übergeben an `slice-v1-abschluss-release-v0-1-1` (dort §1, DoD 2 und §8; `SPEC-067` Punkt 18).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `verhalten-nur-unter-linux-geprueft` (Plattform; hier als Randform (h) vor dem Code entschieden, `AGENTS.md` §3.12), `folge-slice-adresse-nimmt-nicht-an` (die zwei Sendungen in §1 stehen in den Nehmern, §3.13), `folge-slice-liefert-handbuch-teil-nicht` (der Abschnitt zu `v0.1.0` hat keinen Nehmer mit Kennung, §6 und Auftrag an den Planner). Der Stand der Zähler dieser drei Einträge ist beim Start zu lesen; erreicht einer mit diesem Slice drei, braucht er einen Folge-Slice.

**Zählung** (Architect, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: Namen in den drei Skripten, ihren Gegenproben und der Zeile von `harness/README.md`; 2: Handbuch und `releasing.md`). **Schichten: zwei** — Harness-Werkzeug (Skripte, Gegenproben, Zeile in `harness/README.md`) und Nutzer- und Wartungs-Doku (`docs/user/`, `docs/maintainer/`); Produkt-Code zählt nicht (§1). Die Verträge liegen vor dem Code beim Architect und zählen nicht als Liefer-Punkt. Größe: Schätzung und Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält Spezifikation und Entscheidungen vor dem Code; Verträge und [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) gehen den Skripten voraus).
