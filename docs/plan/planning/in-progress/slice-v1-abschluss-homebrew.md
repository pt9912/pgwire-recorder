# Slice slice-v1-abschluss-homebrew: Homebrew-Bereitstellung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-19.a` · `SPEC-042` · `SPEC-067` (Punkt 1c und 19, [`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md)) · `SPEC-035` · `SPEC-055` (Punkt 3, Namen der Binaries) · `LH-FA-01`

**Entscheidung:** [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) (Accepted, Entscheidung des Nutzers vom 2026-10-20) · [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Namen der Binaries ohne Version; Accepted, Entscheidung des Nutzers vom 2026-10-20)

**Verantwortlich:** pt9912 (Zweig impl-brew)
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Bei jedem stabilen Tag erzeugt ein eigener Job der Release-Pipeline (`homebrew`, nach dem GitHub-Release) die Formel `Formula/pgwire-recorder.rb` aus einer Vorlage und der Datei `SHA256SUMS` des Release und legt sie im Tap `pt9912/homebrew-pgwire-recorder` ab; eine Vorabversion ändert den Tap nicht. Belegt ist das gegen ein lokales Probe-Tap (ein Bare-Repository im Temp-Verzeichnis der Gegenprobe), ohne Netz und ohne Geheimnis. Form der Formel, Ablage, Reihenfolge der Gründe, Zugang und Fälle entscheidet [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 1c und 19; warum so, [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md).

**Übernommen aus `slice-erster-release-homebrew-nachweis`** (dort §1, Abgrenzung): das Erzeugen der Formel; das ist DoD-Punkt 1, jener Slice wendet das Verfahren mit dem echten Tap an (der Job `homebrew` am nächsten stabilen Tag; ein Nachtrag von Hand nur für einen Tag mit dem neuen Schema, Punkt 19 (7); für `v0.1.0` entfällt er, Entscheidung des Nutzers vom 2026-10-20).

**Bindung aus `slice-v1-abschluss-binary-namen`** (`AGENTS.md` §3.13; Geber: dieser Slice liegt hinter ihm): Die Binaries heißen `pgwire-recorder-<os>-<arch>` ohne Version ([ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md), `SPEC-055` Punkt 3). Vorlage, Prüfung der Summen und Fälle (`SPEC-067` Punkt 19 (3), (4), (9)) rechnen mit diesen Namen; der Fall `p19-summen-altes-schema` (nur Namen des alten Schemas, wie sie `v0.1.0` trägt: Exit 1 `summen`) und die Mutation „der Name in der Adresse trägt die Version“ gehören zu Punkt 1 der DoD.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Formel für `v0.1.0` — entfällt (Entscheidung des Nutzers vom 2026-10-20): Seine Binaries tragen das alte Schema, und ein neuer Tag mit dem neuen Schema liefert die erste Formel; der Schritt lehnt das alte Schema ab (Klasse 3, anderer Vorgang; die Pflicht, den Tag zu setzen, liegt beim Nutzer, `slice-v1-abschluss-release-v0-1-1` bereitet ihn vor, `slice-erster-release-homebrew-nachweis` wartet auf ihn).
- Die Installation aus dem echten Tap, der Abschnitt zu Homebrew im Benutzerhandbuch (§2 *Installation*, Befehle) und die Zeile in `README.md` — `slice-erster-release-homebrew-nachweis` (Klasse 1; die Sendung steht dort in §1 und DoD, mit der Kennung dieses Slice). Grund: Das Handbuch nennt nur Befehle, die einmal gelaufen sind (`AGENTS.md` §3.11, `SPEC-067` Punkt 18), und der Lauf mit dem echten Tap ist dessen Gegenstand; die Formel ist kein Verhalten des Binaries. Dieser Slice ändert deshalb weder Handbuch noch README.
- Das Anlegen und die Rechte des Geheimnisses `HOMEBREW_TAP_GITHUB_TOKEN` und jede Handlung am echten Tap — Handlung des Betreibers (Klasse 3, anderer Vorgang). Der Slice liest, schreibt und gibt den Wert nie aus und pusht nie in den echten Tap.
- Ein Probe-Tag — er wäre eine Vorabversion (kann den Tap nicht ändern) und ein weiterer Tag mit Veröffentlichung ([ADR-0051](../../../plan/adr/0051-veroeffentlichung-auf-github-actions.md)); an seine Stelle tritt das lokale Probe-Tap. Das weicht von der Reihenfolge in `welle-v1-nachtrag` §5 ab, die den Probe-Tag nannte; die Reihenfolge ist nachgezogen.
- Aufnahme in das Standard-Repository von Homebrew — Out-of-Scope von LH-FA-19.
- Windows — Homebrew trägt Windows nicht; dort gilt das Binary direkt (LH-FA-01).
- Docker/OCI-Image — `slice-v1-abschluss-container`.
- Änderung des Lastenhefts — `LH-FA-19` trägt Happy Path, Boundary und Negative schon; die Einzelregeln stehen in der Spezifikation.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Der Schritt `homebrew` von `tools/harness/release.sh` mit der Vorlage `tools/harness/homebrew-formula.rb.tmpl` und der Job `homebrew` in `.github/workflows/release.yml` (Text nach `SPEC-067` Punkt 1c) erzeugen die Formel aus `SHA256SUMS` eines Fixture-Release und legen sie in einem Probe-Tap ab (lokales Bare-Repository, kein Netz, kein Geheimnis); die Prüfsummen und Plattformen stimmen, und jede Randform aus `SPEC-067` Punkt 19 (9) hat ihren Fall und ihre Mutation in `make release-gegenprobe` (Zusage · Mutation · roter Fall im Bericht, `AGENTS.md` §3.10); die Zeile von `make release-gegenprobe` in `harness/README.md` §Sensors nennt die neuen Fälle.
- [ ] Eine Vorabversion verändert den Tap nicht (`p19-vorabversion`, mit Mutation).
- [ ] `docs/maintainer/releasing.md` beschreibt das Ausgeführte: §3 (kein Probe-Tag für die Formel, das Probe-Tap), §4 (das Geheimnis als Vorbedingung des Betreibers, mit dem Recht, das der Token mindestens braucht), §6 (der Job), §7, §11 (die Zeile). Handbuch und `README.md` bleiben unverändert (§1).
- [ ] `make gates` grün.
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
| `tools/harness/release.sh` (Schritt `homebrew`) | update | Eingaben, Summen, Formel, Tap, Token über `GIT_ASKPASS` (`SPEC-067` Punkt 19) |
| `tools/harness/homebrew-formula.rb.tmpl` | neu | Vorlage mit Platzhaltern (Text in Punkt 19 (4)) |
| `.github/workflows/release.yml` | update | Job `homebrew` nach Punkt 1c; das Geheimnis nur im Schritt `Homebrew` |
| `tools/harness/release-gegenprobe.sh` | update | Fälle `p19-*` mit Probe-Tap, einer Attrappe von `gh` für `SHA256SUMS` und einem Protokoll vor dem echten `git`; `p1-*` um die Zeilen von 1c; Liste der Schritte mit Außenwirkung in `p2-*` |
| `harness/README.md` (§Sensors, Zeile `release-gegenprobe`) | update | nennt die neuen Fälle |
| `docs/maintainer/releasing.md` | update | §1, §2, §3, §4, §6, §7, §11 und die Änderungshistorie (Ist-Zustand des Ausgeführten; §1 und §2 sagten, die Formel entstehe nicht in der Pipeline) |

Nicht in dieser Liste: der Tap `pt9912/homebrew-pgwire-recorder` (Handlung des Betreibers, `slice-erster-release-homebrew-nachweis`), Handbuch und `README.md` (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-binary-namen` ist `done` (die Namen der Binaries, auf die Vorlage und Fälle rechnen; Reihenfolge in `welle-v1-nachtrag` §5), `slice-v1-abschluss-container` ist `done` (Release-Artefakte liegen vor), `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`; Entscheidung des Nutzers vom 2026-10-11), und [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) ist `Accepted` (die Annahme ist die Entscheidung des Nutzers; bis dahin sind Art und Rechte des Tap-Tokens eine Empfehlung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 0256ddff..<letzter Commit des Slice> -- tools .github harness/README.md docs/maintainer` (Tests eingeschlossen, `AGENTS.md` §3.14) **über 1300** liegt. Schätzung des Architect in dieser Metrik: Produktcode rund 250 (Schritt etwa 140, Vorlage 33, Job 24, `releasing.md` etwa 50), die Gegenprobe nach dem Faktor von `release-gegenprobe.sh` zu `release.sh` (rund das Vierfache) 700 bis 800, zusammen 900 bis 1200 — am oberen Rand einer Review-Sitzung. Vorgemerkter Schnitt: (a) Schritt, Vorlage und Fälle `p19-*` (liefert den Nachtrag von Hand), (b) Job in `release.yml`, `p1-*`, `releasing.md`; (b) setzt (a) voraus und bleibt einzeln lieferbar, weil (a) ohne Job über den Nachtrag schon eine Formel liefert.
- `in-progress` → `next`: Der Nutzer entscheidet gegen die Empfehlung in [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) (anderer Zugang als ein Token über HTTPS, etwa Deploy-Key oder App) — zurück zur Zerlegung.
- `in-progress` → `open`: Eine Eingabe des Schritts lässt sich nicht ohne Netz oder echtes Geheimnis fassen (etwa der Wortlaut von `gh release download`) und `SPEC-067` Punkt 19 kann sie nicht als Grenze festhalten — Carveout.

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

**Randformen des Vertrags** (`AGENTS.md` §3.12; alle vor dem Code entschieden, jede an der Stelle in `SPEC-067` bzw. [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md)). Das Vorbild ist das Verfahren von u-boot (`scripts/homebrew-formula.rb.tmpl`, `homebrew-formula-fill.sh`, `tap-nachzug.sh`, dort ADR 0016 von u-boot); Abweichungen stehen mit Grund in der Tabelle.

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Form der Formel | Ruby-Datei `Formula/pgwire-recorder.rb` aus einer Vorlage mit Platzhaltern wie bei u-boot; `on_macos`/`on_linux` × `on_arm`/`on_intel`, vier Plattformen, kein Windows; `url` je Plattform am Release des Tags, `sha256` aus `SHA256SUMS`; `test do` ruft `version` und erwartet `pgwire-recorder <Version>`. Abweichung von u-boot: eine Zeile `version` steht in der Formel, weil der Name `pgwire-recorder-<os>-<arch>` keine Version trägt ([ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md)) und die Erkennung aus dem Pfad der URL nicht gemessen ist; `brew audit --strict` kann sie als überflüssig melden, ein Befund im Nachweis, kein Fehler der Installation | `SPEC-067` Punkt 19 (4); `LH-FA-19.a` |
| (b) Wer erzeugt, wohin | Schritt `homebrew` in `release.sh` (nutzt Eingaben, Aufräumen, Abbruch des Skripts mit) und Vorlage als Datei; Ablage `Formula/pgwire-recorder.rb` im Tap, ein Commit, Push auf `main`. Abweichung von u-boot: kein Fill-Skript, kein Anhang am Release (SPEC-067 Punkt 9 verlangt genau sieben Dateien) | Punkt 19 (1), (4), (5) |
| (c) Vorabversion | Exit 0 ohne Aufruf von `gh` und `git`, vor der Prüfung der Token | Punkt 19 (2) |
| (d) Wiederholung, Fehler, Reihenfolge | Job nach `GitHub-Release` (`needs: release`), damit die URLs der Formel bestehen; Tap schon aktuell: Exit 0 ohne Commit; höhere Version im Tap: Exit 0, unverändert; gleiche Version, anderer Text: neuer Commit; Push abgewiesen: Exit 1, Tap unverändert; der Release besteht schon: die Summen kommen aus ihm, der Schritt legt keinen an; ein roter Job lässt Release und Registries stehen, Wiederanlauf nur dieses Jobs oder Nachtrag von Hand, kein neuer Tag; Exit 2 für Eingaben (Ereignis, Tag-Form, Token), Exit 1 für Summen und Tap, Reihenfolge der Gründe gebunden. Abweichung von u-boot: ein fehlendes Geheimnis ist rot statt Warnung und Überspringen, weil der Job hinter dem Release nichts gefährdet und ein Überspringen den Tap unbemerkt zurückließe | Punkt 19 (2), (3), (5); Punkt 11 |
| (e) Probe-Tap und Belegtes | lokales Bare-Repository (leer und mit Commits), echtes `git`, Attrappe von `gh`; belegt sind Text der Formel, Ablage, Randformen des Taps, Token-Weg; nicht belegt sind Rechte des Tokens, Annahme durch Homebrew, Binary auf macOS, Erreichbarkeit der URLs | Punkt 19 (9); Punkt 14, 15 |
| (f) Aufrufarten, Abbruch, Plattform | `release.sh` läuft nie unter `make`, `-n -q -t -i` und `MAKEFLAGS` haben keinen Fall (Grund wie `SPEC-067` Punkt 13: kein Hook und kein Stand liest das Ergebnis); `SIGINT` und `SIGTERM` räumen auf (Fälle `p19-sigint`, `p19-sigterm`); Linux geprüft, macOS und Windows als Ausführungsort Grenze, die Formel gilt für macOS und Linux und ist dort nicht ausgeführt | Punkt 19 (8), (9); Punkt 13 |
| (g) Zugang | Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` (vom Nutzer gesetzt, kein anderer Name), nur im Schritt `Homebrew` des Jobs `homebrew`, `permissions: {}` oben bleibt, Job mit `contents: read`; Token über `GIT_ASKPASS`, nie als Argument, URL, Datei oder Ausgabe (Abweichung von u-boot: dort steht er in der Klon-URL, hier gilt SPEC-067 Punkt 10). **Empfehlung:** mindestens *Contents: Read and write* auf `pt9912/homebrew-pgwire-recorder` und sonst nichts; der Nutzer bestätigt, dass sein Token dazu passt | Punkt 19 (6); Punkt 10; [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) |
| Nachtrag von Hand | derselbe Schritt von Hand mit den Token in der Umgebung, Handlung des Betreibers, für einen Tag mit dem neuen Schema; für `v0.1.0` nie (Exit 1 `summen`, alte Namen), Fall `p19-summen-altes-schema` | Punkt 19 (3), (7); [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) |

**Sendung an `slice-erster-release-homebrew-nachweis`** (`AGENTS.md` §3.13): Handbuch- und README-Teil, die Formel des ersten stabilen Tags mit dem Schema ohne Version im echten Tap (durch den Job `homebrew`) und die drei Grenzen unten stehen dort in §1, DoD und §6 mit der Kennung dieses Slice.

- Das Verhalten von Homebrew gegenüber Drittanbieter-Taps zeigt sich erst mit dem echten Tap — **Ausgang:** weiter offen, übergeben an `slice-erster-release-homebrew-nachweis` (`welle-v1-nachtrag`, dort §6).
- Zugangsdaten für den Tap-Push: das Geheimnis ist gesetzt (Nutzer, 2026-10-20), Art und Rechte sind unbestätigt, und das Schreibrecht zeigt allein der Push in den echten Tap — **Ausgang:** weiter offen, übergeben an `slice-erster-release-homebrew-nachweis` (dort §6, Risiko Rechte des Tokens).
- Dass Homebrew die Formel annimmt (Syntax, `version`-Zeile, Installation einer Datei ohne Archiv, `test do`) und das Binary auf macOS und Linux läuft, belegt die Gegenprobe nicht — **Ausgang:** weiter offen, übergeben an `slice-erster-release-homebrew-nachweis` (dort §6); eine Korrektur der Vorlage geht ohne neuen Tag über den Nachtrag.
- Der Wortlaut von `gh release download` auf dem Runner ist im Dev-Container nicht gemessen (kein `gh`) — **Ausgang:** weiter offen, übergeben an `slice-erster-release-homebrew-nachweis` (dort §6; der Nachtrag fährt den Schritt mit dem echten `gh`).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt am 2026-10-20 beim Eintragen der Randformen (Architect, `AGENTS.md` §3.13; die Schichtzählung aus `slice-doku-ist-stand` gilt): **Liefer-Punkte: 3** (1: Schritt, Vorlage, Job und Gegenprobe; 2: Vorabversion ändert den Tap nicht; 3: `releasing.md`). **Schichten: zwei** — Harness-Werkzeug (Skript, Vorlage, Workflow, Gegenprobe, Zeile in `harness/README.md`) und Wartungs-Doku (`docs/maintainer/`). Der Handbuch- und README-Teil, der bisher im ersten Punkt stand, ist an `slice-erster-release-homebrew-nachweis` abgegeben. Nachgezählt am 2026-10-20 mit der Bindung aus `slice-v1-abschluss-binary-namen` (Architect, `AGENTS.md` §3.13): Der Fall `p19-summen-altes-schema` und eine Mutation gehören zum ersten Punkt, kein neuer Liefer-Punkt; **Liefer-Punkte bleiben 3, Schichten zwei**; die Schätzung in §4 (900 bis 1200) wächst um unter 40 Zeilen, die Marke bleibt. Größe: Schätzung und vorgemerkter Schnitt in §4; Punkt 1 trägt fast den ganzen Umfang.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
