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

**Ziel:** Bei jedem stabilen Tag erzeugt ein eigener Job der Release-Pipeline (`homebrew`, nach dem GitHub-Release) die Formel `Formula/pgwire-recorder.rb` aus einer Vorlage und der Datei `SHA256SUMS` des Release und legt sie im Tap `pt9912/homebrew-pgwire-recorder` ab; eine Vorabversion ändert den Tap nicht. Belegt ist das gegen ein lokales Probe-Tap (ein Bare-Repository im Temp-Verzeichnis der Gegenprobe), ohne Netz und ohne Geheimnis. Eine fehlende, leere oder nicht lesbare Vorlage endet mit Exit 1 und dem Stichwort `vorlage`, nach den Summen und vor dem Klon des Taps. Form der Formel, Ablage, Reihenfolge der Gründe, Zugang und Fälle entscheidet [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 1c und 19; warum so, [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md).

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

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Der Schritt `homebrew` von `tools/harness/release.sh` mit der Vorlage `tools/harness/homebrew-formula.rb.tmpl` und der Job `homebrew` in `.github/workflows/release.yml` (Text nach `SPEC-067` Punkt 1c) erzeugen die Formel aus `SHA256SUMS` eines Fixture-Release und legen sie in einem Probe-Tap ab (lokales Bare-Repository, kein Netz, kein Geheimnis); die Prüfsummen und Plattformen stimmen, die Zustände der Vorlage (fehlt, leer, Verzeichnis) enden mit Exit 1 `vorlage`, und jede Randform aus `SPEC-067` Punkt 19 (9) hat ihren Fall und ihre Mutation in `make release-gegenprobe` (Zusage · Mutation · roter Fall im Bericht, `AGENTS.md` §3.10); die Zeile von `make release-gegenprobe` in `harness/README.md` §Sensors nennt die neuen Fälle.
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
| `tools/harness/release.sh` (Schritt `homebrew`) | update | Eingaben, Summen, Zustände der Vorlage (Stichwort `vorlage`, nach den Summen, vor dem Klon), Formel, Tap, Token über `GIT_ASKPASS` (`SPEC-067` Punkt 19); sechs Anker-Kommentare ohne Leser entfernt |
| `tools/harness/homebrew-formula.rb.tmpl` | neu | Vorlage mit Platzhaltern (Text in Punkt 19 (4)) |
| `.github/workflows/release.yml` | update | Job `homebrew` nach Punkt 1c; das Geheimnis nur im Schritt `Homebrew` |
| `tools/harness/release-gegenprobe.sh` | update | Fälle `p19-*` (darunter `p19-vorlage-fehlt`, `-leer`, `-verzeichnis`, `-nach-summen`, `-vor-tap` und vier Mutationen der Prüfung der Vorlage, dazu `m-p19-token-in-url`) mit Probe-Tap, einer Attrappe von `gh` für `SHA256SUMS` und einem Protokoll vor dem echten `git`; `p1-*` um die Zeilen von 1c; Liste der Schritte mit Außenwirkung in `p2-*` |
| `harness/README.md` (§Sensors, Zeile `release-gegenprobe`) | update | nennt die neuen Fälle, auch die Zustände der Vorlage |
| `docs/maintainer/releasing.md` | update | §1, §2, §3, §4, §6, §7, §11 und die Änderungshistorie (Ist-Zustand des Ausgeführten; §1 und §2 sagten, die Formel entstehe nicht in der Pipeline) |

Nicht in dieser Liste: der Tap `pt9912/homebrew-pgwire-recorder` (Handlung des Betreibers, `slice-erster-release-homebrew-nachweis`), Handbuch und `README.md` (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-binary-namen` ist `done` (die Namen der Binaries, auf die Vorlage und Fälle rechnen; Reihenfolge in `welle-v1-nachtrag` §5), `slice-v1-abschluss-container` ist `done` (Release-Artefakte liegen vor), `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`; Entscheidung des Nutzers vom 2026-10-11), und [ADR-0053](../../../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) ist `Accepted` (die Annahme ist die Entscheidung des Nutzers; bis dahin sind Art und Rechte des Tap-Tokens eine Empfehlung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 8dbe6a14..<letzter Commit des Slice> -- tools .github harness/README.md docs/maintainer` (Basis: der Move nach `in-progress/`, ohne die Dateien von `slice-v1-abschluss-binary-namen`; Tests eingeschlossen, `AGENTS.md` §3.14) **über 1300** liegt. Schätzung des Architect in dieser Metrik: Produktcode rund 250 (Schritt etwa 140, Vorlage 33, Job 24, `releasing.md` etwa 50), die Gegenprobe nach dem Faktor von `release-gegenprobe.sh` zu `release.sh` (rund das Vierfache) 700 bis 800, zusammen 900 bis 1200 — am oberen Rand einer Review-Sitzung. Vorgemerkter Schnitt: (a) Schritt, Vorlage und Fälle `p19-*` (liefert den Nachtrag von Hand), (b) Job in `release.yml`, `p1-*`, `releasing.md`; (b) setzt (a) voraus und bleibt einzeln lieferbar, weil (a) ohne Job über den Nachtrag schon eine Formel liefert.
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
| (h) Vorlage fehlt, leer oder nicht lesbar (Rückgabe des Reviews F-1, Architect 2026-10-20) | zulässig: lesbare Datei mit mindestens einem Byte; sonst Exit 1, Stichwort `vorlage`, nach den Summen und vor dem Klon, der Tap bleibt unberührt; Reihenfolge `summen-laden`, `summen`, `vorlage`, dann Tap; Fälle `p19-vorlage-fehlt`, `p19-vorlage-leer`, `p19-vorlage-verzeichnis`, `p19-vorlage-nach-summen`, `p19-vorlage-vor-tap` und vier Mutationen (Prüfung entfällt, Stichwort, vor den Summen, nach dem Klon); der Kopfkommentar von `release.sh` nennt `vorlage` unter den Stichworten | Punkt 19 (4), (9) |
| (i) Anker-Kommentare ohne Leser (Review F-2) | Die sechs Anker `P19aktuell`, `P19liste`, `P19summen-aufruf`, `P19summen-zahl`, `P19url`, `P19version` entfallen aus `release.sh` (§3.7: die Mutationen ersetzen Text, nicht Anker; ein Anker ohne Leser ist keine Kopplung); keine neuen Fälle | — |
| (j) Beispiel-Mutation „Token in der URL“ (Review F-3) | eigene Mutation `m-p19-token-in-url` in `release-gegenprobe.sh`, rot `p19-token-nicht-im-argument` und `p19-token-nicht-in-url-und-config` (vom Review gesehen); zeigt der Lauf eine andere Menge, geht sie an den Architect | Punkt 19 (9) |
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

### Belege des Implementers

Nacharbeit zum Review `slice-v1-abschluss-homebrew` (Report vom 2026-10-20, Kopf `3043aed9`), auf den Entscheidungen des Architect in `SPEC-067` Punkt 19 (4), (6), (9) und §6 (h) bis (j). Umgebung: Dev-Container `pgr-dev2`, Worktree auf Zweig `impl-brew2`, die Läufe der Gegenprobe mit `GP_MUT=<Muster>`, `GP_ENTDECKE=1`, `GP_PARALLEL=4`, jeder Lauf unter drei Minuten. Kein voller Lauf von `make release-gegenprobe` (Entscheidung des Nutzers vom 2026-10-20).

**Befunde, Behebung, Beleg**

| Befund | Behebung | Beleg |
|---|---|---|
| F-1 (MEDIUM) Prüfung der Vorlage ohne Fall und Mutation | Stichwort `vorlage` im Kopfkommentar von `release.sh`; fünf Fälle `p19-vorlage-fehlt`, `-leer`, `-verzeichnis`, `-nach-summen`, `-vor-tap`; vier Mutationen (siehe Tabelle) | unveränderte Kopie: Gruppe der Fälle `p19-*` grün (30 s, Lauf mit Wegwerfkopie der Gegenprobe, die nur diese Gruppe fährt, danach gelöscht); Mutationen: Tabelle unten, Lauf `m-p19-vorlage-(pruefung|stichwort|vor-summen|nach-klon)` 34 s, genau die benannten Fälle rot |
| F-2 (LOW) sechs Anker ohne Leser | die Kommentare `P19aktuell`, `P19liste`, `P19summen-aufruf`, `P19summen-zahl`, `P19url`, `P19version` aus `release.sh` entfernt, Verhalten unverändert | vorher `git show 01850212:tools/harness/release.sh \| grep -c 'P19aktuell\|P19liste\|P19summen-aufruf\|P19summen-zahl\|P19url\|P19version'` = 6, nachher 0; die Mutationen, die diese Zeilen ersetzen, laufen unverändert rot (Tabelle) |
| F-3 (LOW) Beispiel-Mutation „Token in der URL“ fehlte | `m-p19-token-in-url` (die URL des Klons trägt `x-access-token:<Token>@`, für den Probe-Tap als `file://…@localhost<Pfad>`) | rot genau `p19-token-nicht-im-argument` und `p19-token-nicht-in-url-und-config`, die Menge des Architect |
| F-5 (INFO) Belege fehlten | dieser Abschnitt | Tabelle unten |
| F-7 (INFO) Basis der Größenmarke | §4 nennt `8dbe6a14` | `git diff --numstat 8dbe6a14..cdfb2ad7 -- tools .github harness/README.md docs/maintainer`: 860 Einfügungen, 49 Löschungen, Summe 909, Marke (1300) nicht erreicht |

Nebenfolge von F-1: Fünf bestehende Mutationen färbten die neuen Fälle mit und trugen sie nun in der Liste der erwarteten roten Fälle: `m-p19-summen-form-63`, `m-p19-summen-form-entfernt`, `m-p19-summen-liste-ohne-linux-amd64` (je `p19-vorlage-nach-summen`: ohne die Prüfung der Summe endet der Fall mit `vorlage` statt `summen`), `m-p19-tap-vor-summen-geklont` (alle fünf Fälle: der Klon steht dann vor der Prüfung der Vorlage), `m-p19-formula-nicht-angelegt` (die fünf Fälle aus der Liste der Fälle, die den Tap erreichen müssen, ausgenommen: sie enden vor der Formel). Der erste Lauf zeigte die fünf Abweichungen rot, nach der Anpassung ist der Lauf `m-p19-(summen-form-63|summen-form-entfernt|tap-vor-summen|formula-nicht|summen-liste-ohne-linux-amd64)` grün (62 s).

**Läufe** (jeder `GP_MUT=<Muster> GP_ENTDECKE=1 GP_PARALLEL=4 bash tools/harness/release-gegenprobe.sh`, Stand `cdfb2ad7` bis auf die Zeile der Anpassung oben; Exit 1 in zwei Läufen vor der Anpassung der fünf Listen)

| Muster (nach `^m-p19-`, außer der letzten Zeile) | Mutationen | Dauer | Ausgang |
|---|---|---|---|
| `(ereignis\|ref-nicht\|tag-form\|vorabversion\|reihenfolge-)` | 12 | 98 s | grün |
| `(gh-token-(entfernt\|leer)\|tap-token\|geheimnis\|braucht\|summen-(ohne\|tag\|dir\|laden\|zahl))` | 13 | 148 s | grün |
| `(summen-(name\|fremde\|liste\|form\|altes)\|tap-vor)` | 14 | 129 s | 4 rot (die genannten Listen), danach grün |
| `(version-\|formel-ohne\|vorlage-)` | 15 | 126 s | grün |
| `(formula\|andere\|gleiche-version-andere\|verlangt\|ruft-)` | 7 | 118 s | 1 rot (`formula-nicht-angelegt`), danach grün |
| `(leer-ohne\|zweig\|pfad\|formel-(zeilen\|stabil\|keine)\|sort\|hoeher\|aktuell\|gleiche-version-immer)` | 14 | 125 s | grün |
| `(klon\|autor\|committer\|nachricht\|zwei\|url-vorgabe)` | 11 | 94 s | grün |
| `push` | 6 | 63 s | grün |
| `(askpass\|prompt\|credential\|token\|gh-token-im\|set-x)` | 14 | 130 s | grün |
| `(ende\|int-\|term-)` | 7 | 49 s | grün |
| `^m-p1-homebrew` | 4 | 3 s | grün |
| `ok` (unveränderte Kopie, nur Gruppe `p19-*`) | 1 | 30 s | grün |

Die Mutationen jeder Zeile der Tabelle unten sind damit gefahren; jede Zeile nennt die roten Fälle aus der Ausgabe `GP_ENTDECKE`.

**Zusage · Mutation · roter Test**

| Zusage (`SPEC-067` Punkt 19) | Mutation | roter Test |
|---|---|---|
| Punkt 1c: `RELEASE_TAP_TOKEN` steht nur im Schritt, nicht im Job | `m-p1-homebrew-geheimnis-job-erlaubt` | `p1-homebrew-geheimnis-nur-im-schritt` |
| Punkt 1c: das Tap-Geheimnis steht nicht im Job `release` | `m-p1-homebrew-geheimnis-release-job-erlaubt` | `p1-homebrew-geheimnis-release-job` |
| Punkt 1c: das Tap-Geheimnis steht in keinem Schritt des Jobs `release` | `m-p1-homebrew-geheimnis-release-schritt-erlaubt` | `p1-homebrew-geheimnis-release-schritt` |
| Punkt 1c: der Job `homebrew` hat kein `packages: write` | `m-p1-homebrew-registry-rechte-erlaubt` | `p1-homebrew-keine-registry-rechte` |
| (5) gleiche Version, gleicher Text: kein Commit, kein Push | `m-p19-aktuell-immer-neu` | `p19-tap-gleich-aktuell` |
| (5) fremde Dateien im Tap bleiben | `m-p19-andere-dateien-geloescht` | `p19-tap-andere-dateien-bleiben`, `p19-tap-nicht-leer` |
| (5) andere Formeln im Tap bleiben | `m-p19-andere-formeln-geloescht` | `p19-tap-andere-dateien-bleiben`, `p19-tap-formula-vorhanden` |
| (6) Benutzername `x-access-token` | `m-p19-askpass-benutzer-anders` | `p19-token-askpass-benutzer` |
| (6) der Token geht über `GIT_ASKPASS` an `git` | `m-p19-askpass-nicht-gesetzt` | `p19-token-askpass-benutzer`, `p19-token-askpass-gesetzt`, `p19-token-askpass-passwort` |
| (6) Passwort ist der Token | `m-p19-askpass-passwort-leer` | `p19-token-askpass-passwort` |
| (5) Commit-Identität: Mail des Autors | `m-p19-autor-mail` | `p19-commit-identitaet` |
| (5) Commit-Identität: Name des Autors | `m-p19-autor-name` | `p19-commit-identitaet` |
| Punkt 13: der Schritt braucht `RUNNER_TEMP` nicht und liest keinen Zustand | `m-p19-braucht-runner-temp` | `p19-kein-zustand`, `p19-ohne-runner-temp` |
| (5) Commit-Identität: Mail des Committers | `m-p19-committer-mail` | `p19-commit-identitaet` |
| (5) Commit-Identität: Name des Committers | `m-p19-committer-name` | `p19-commit-identitaet` |
| (6) kein Credential-Helper | `m-p19-credential-helper-nicht-geleert` | `p19-credential-helper-leer` |
| Punkt 13 (Abbruch): Aufräumen bei Ende | `m-p19-ende-falle-fehlt` | `p19-ende` |
| (2a) Ereignis und Ref: nur ein Push eines Tags `v…` (Exit 2 `ereignis`), Prüfung entfällt | `m-p19-ereignis-entfernt` | `p19-ohne-push`, `p19-ref-kein-tag`, `p19-reihenfolge-gruende-a-b` |
| (2a) das Ereignis `push` wird verlangt | `m-p19-ereignis-nicht-geprueft` | `p19-ohne-push` |
| (5) Prüfung der Formel im Tap entfällt | `m-p19-formel-keine-pruefung` | `p19-tap-formel-ohne-version`, `p19-tap-formel-vorabversion`, `p19-tap-formel-zwei-versionen` |
| (4) LF am Ende der Formel | `m-p19-formel-ohne-lf-am-ende` | `p19-formel-text`, `p19-tap-gleich-aktuell` |
| (5) die Version im Tap ist eine stabile Version | `m-p19-formel-stabil-entfernt` | `p19-tap-formel-vorabversion` |
| (5) die Formel im Tap hat genau eine Zeile `version` | `m-p19-formel-zeilen-mindestens-eine` | `p19-tap-formel-zwei-versionen` |
| (5) die Formel wird unter `Formula/pgwire-recorder.rb` angelegt (alle Fälle, die sie lesen, rot) | `m-p19-formula-nicht-angelegt` | 37 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (2e) die Meldung nennt den Namen `HOMEBREW_TAP_GITHUB_TOKEN` | `m-p19-geheimnis-ohne-namen` | `p19-tap-token-fehlt`, `p19-tap-token-leer` |
| (2d) `GH_TOKEN` wird verlangt (Exit 2 `token`) | `m-p19-gh-token-entfernt` | `p19-gh-token-fehlt`, `p19-gh-token-leer`, `p19-reihenfolge-gruende-d-e` |
| (6) `GH_TOKEN` steht nicht in der Ausgabe | `m-p19-gh-token-im-log` | `p19-token-nicht-im-log` |
| (2d) ein leeres `GH_TOKEN` gilt als fehlend | `m-p19-gh-token-leer-gilt-als-gesetzt` | `p19-gh-token-leer` |
| (5) gleiche Version, andere Summe: neuer Commit | `m-p19-gleiche-version-andere-summe-gesperrt` | `p19-tap-gleich-andere-summe` |
| (5) gleiche Version, anderer Text: neuer Commit | `m-p19-gleiche-version-immer-aktuell` | `p19-tap-gleich-andere-summe`, `p19-tap-gleich-andere-vorlage` |
| (5) eine höhere Version im Tap bleibt unverändert | `m-p19-hoeher-ueberschrieben` | `p19-tap-hoeher`, `p19-tap-semver-numerisch` |
| Punkt 13 (Abbruch): `SIGINT`: Aufräumen, Exit 130 | `m-p19-int-falle-fehlt` | `p19-sigint` |
| Punkt 13 (Abbruch): `SIGINT` beendet die Nachkommen | `m-p19-int-kinder-bleiben` | `p19-sigint` |
| Punkt 13 (Abbruch): `SIGINT`: Verzeichnis erst nach den Nachkommen entfernt | `m-p19-int-verzeichnis-zu-frueh` | `p19-sigint` |
| (5) Klon scheitert: Exit 1 `tap-klon` | `m-p19-klon-fehler-ignoriert` | `p19-klon-scheitert` |
| (5) Fehlerausgabe von `git` mit Präfix `release: git: ` | `m-p19-klon-meldung-entfernt` | `p19-klon-scheitert` |
| (5) leerer Tap: Zweig `main` | `m-p19-leer-ohne-main` | `p19-tap-leer-head-anders` |
| (5) Nachricht ohne weiteren Text | `m-p19-nachricht-mit-text` | `p19-commit-nachricht` |
| (5) Nachricht `pgwire-recorder <Version>` ohne `v` | `m-p19-nachricht-mit-v` | `p19-commit-nachricht` |
| (5) `Formula/pgwire-recorder.rb` ist eine reguläre Datei | `m-p19-pfad-datei-nicht-geprueft` | `p19-tap-pfad-formel-symlink`, `p19-tap-pfad-formel-verzeichnis` |
| (5) die Formel im Tap ist kein Symlink | `m-p19-pfad-datei-symlink-erlaubt` | `p19-tap-pfad-formel-symlink` |
| (5) `Formula` ist ein Verzeichnis | `m-p19-pfad-verzeichnis-nicht-geprueft` | `p19-tap-pfad-datei`, `p19-tap-pfad-formula-symlink` |
| (5) `Formula` ist kein Symlink | `m-p19-pfad-verzeichnis-symlink-erlaubt` | `p19-tap-pfad-formula-symlink` |
| (6) keine Abfrage am Terminal (`GIT_TERMINAL_PROMPT=0`) | `m-p19-prompt-erlaubt` | `p19-kein-prompt` |
| (5) Push auf `main` | `m-p19-push-auf-anderen-zweig` | `p19-push-ohne-force` |
| (5) abgewiesener Push endet mit Exit 1 `tap-push` | `m-p19-push-fehler-ignoriert` | `p19-push-abgewiesen`, `p19-token-nicht-im-log` |
| (5) kein `--force` | `m-p19-push-force` | `p19-push-ohne-force` |
| (5) kein `--force-with-lease` | `m-p19-push-force-with-lease` | `p19-push-ohne-force`, `p19-push-ohne-lease` |
| (5) Fehlerausgabe von `git` beim abgewiesenen Push | `m-p19-push-meldung-entfernt` | `p19-push-abgewiesen` |
| (5) kein zweiter Push nach Abweisung | `m-p19-push-wiederholt` | `p19-push-abgewiesen` |
| (2a) der Ref `refs/tags/v…` wird verlangt | `m-p19-ref-nicht-geprueft` | `p19-ref-kein-tag`, `p19-reihenfolge-gruende-a-b` |
| (2) Reihenfolge der Gründe, Paar (a) vor (b) | `m-p19-reihenfolge-a-b-vertauscht` | `p19-ref-kein-tag`, `p19-reihenfolge-gruende-a-b` |
| (2) Reihenfolge der Gründe, Paar (b) vor (c) | `m-p19-reihenfolge-b-c-vertauscht` | `p19-reihenfolge-gruende-b-c` |
| (2) Reihenfolge der Gründe, Paar (c) vor (d) | `m-p19-reihenfolge-c-d-vertauscht` | `p19-reihenfolge-gruende-c-d`, `p19-vorabversion` |
| (2) Reihenfolge der Gründe, Paar (d) vor (e) | `m-p19-reihenfolge-d-e-vertauscht` | `p19-reihenfolge-gruende-d-e` |
| (1) der Schritt liest den Release nur (kein `gh release upload`) | `m-p19-ruft-gh-upload` | `p19-nur-lesen`, `p19-summen-aufruf` |
| (1) der Schritt ruft `skopeo` nicht | `m-p19-ruft-skopeo` | `p19-nur-lesen` |
| (6) kein `set -x` | `m-p19-set-x` | `p19-gh-token-fehlt`, `p19-reihenfolge-gruende-d-e`, `p19-token-nicht-im-log` |
| (5) Versionen vergleichen numerisch (`0.10.0` über `0.9.0`) | `m-p19-sort-lexikografisch` | `p19-tap-semver-numerisch` |
| (5) Richtung des Vergleichs | `m-p19-sort-umgekehrt` | `p19-tap-hoeher`, `p19-tap-niedriger`, `p19-tap-semver-numerisch` |
| (3) das alte Schema `pgwire-recorder_<Version>_<os>_<arch>` wird nicht angenommen | `m-p19-summen-altes-schema-angenommen` | `p19-summen-altes-schema` |
| (3) das Verzeichnis für `SHA256SUMS` ist neu und leer | `m-p19-summen-dir-nicht-neu` | `p19-summen-aufruf` |
| (3) Summe hat genau 64 Zeichen (nicht 63) | `m-p19-summen-form-63` | `p19-summen-form-zu-kurz`, `p19-summen-vor-tap`, `p19-vorlage-nach-summen` |
| (3) Summe hat genau 64 Zeichen (nicht 65) | `m-p19-summen-form-65` | `p19-summen-form-zu-lang` |
| (3) Form der Summe wird geprüft | `m-p19-summen-form-entfernt` | `p19-summen-form-grossbuchstaben`, `p19-summen-form-kein-hex`, `p19-summen-form-zu-kurz`, `p19-summen-form-zu-lang`, `p19-summen-vor-tap`, `p19-vorlage-nach-summen` |
| (3) Summe nur aus `0-9a-f` (keine Großbuchstaben) | `m-p19-summen-form-grossbuchstaben` | `p19-summen-form-grossbuchstaben` |
| (3) Summe nur aus `0-9a-f` (kein anderes Zeichen) | `m-p19-summen-form-kein-hex` | `p19-summen-form-kein-hex` |
| (3) weitere Namen in `SHA256SUMS` werden gelesen, nicht verwendet, kein Fehler | `m-p19-summen-fremde-zeile-abgelehnt` | `p19-summen-weitere-namen-ignoriert` |
| (3) Exit ungleich 0 von `gh` endet mit Exit 1 `summen-laden` | `m-p19-summen-laden-fehler-ignoriert` | `p19-summen-laden-scheitert` |
| (3) Windows steht nicht in der Liste und ist kein Fehler | `m-p19-summen-liste-mit-windows` | `p19-summen-windows-ignoriert` |
| (3) alle vier Plattformen werden gelesen, eine fehlende ist Exit 1 `summen` | `m-p19-summen-liste-ohne-darwin-amd64` | `p19-formel-kein-platzhalter`, `p19-formel-summe-darwin-amd64`, `p19-formel-text`, `p19-summen-fehlt-darwin-amd64`, `p19-tap-gleich-aktuell` |
| (3) alle vier Plattformen werden gelesen, eine fehlende ist Exit 1 `summen` | `m-p19-summen-liste-ohne-darwin-arm64` | `p19-formel-kein-platzhalter`, `p19-formel-summe-darwin-arm64`, `p19-formel-text`, `p19-summen-fehlt-darwin-arm64`, `p19-tap-gleich-aktuell` |
| (3) alle vier Plattformen werden gelesen, eine fehlende ist Exit 1 `summen` | `m-p19-summen-liste-ohne-linux-amd64` | 13 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (3) alle vier Plattformen werden gelesen, eine fehlende ist Exit 1 `summen` | `m-p19-summen-liste-ohne-linux-arm64` | `p19-formel-kein-platzhalter`, `p19-formel-summe-linux-arm64`, `p19-formel-text`, `p19-summen-fehlt-linux-arm64`, `p19-tap-gleich-aktuell` |
| (3) das zweite Feld ist gleich dem Namen (kein Präfix) | `m-p19-summen-name-praefix` | `p19-summen-name-genau` |
| (3) der Aufruf von `gh` trägt `--dir <Temp-Verzeichnis>` (ohne: der Schritt kommt nie zum Tap, alle Fälle bis dahin rot) | `m-p19-summen-ohne-dir` | 73 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (3) der Aufruf von `gh` trägt `--pattern SHA256SUMS` | `m-p19-summen-ohne-pattern` | `p19-summen-aufruf` |
| (3) der Aufruf von `gh` trägt `--repo $RELEASE_REPOSITORY` | `m-p19-summen-ohne-repo` | `p19-summen-aufruf` |
| (3) der Aufruf von `gh` nennt den Tag mit `v` | `m-p19-summen-tag-ohne-v` | `p19-summen-aufruf` |
| (3) genau eine Zeile je Name (nicht mindestens eine) | `m-p19-summen-zahl-mindestens-eins` | `p19-summen-doppelt` |
| (2b) Form des Tags (Exit 2 `tag-form`) | `m-p19-tag-form-entfernt` | `p19-reihenfolge-gruende-b-c`, `p19-tag-form` |
| (2e) `RELEASE_TAP_TOKEN` wird verlangt (Exit 2 `geheimnis`) | `m-p19-tap-token-entfernt` | `p19-tap-token-fehlt`, `p19-tap-token-leer` |
| (2e) ein leerer `RELEASE_TAP_TOKEN` gilt als fehlend | `m-p19-tap-token-leer-gilt-als-gesetzt` | `p19-tap-token-leer` |
| (3) `summen-laden`, `summen` und `vorlage` enden, bevor der Tap geklont wird | `m-p19-tap-vor-summen-geklont` | 7 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| Punkt 13 (Abbruch): `SIGTERM`: Aufräumen, Exit 143 | `m-p19-term-falle-fehlt` | `p19-sigterm` |
| Punkt 13 (Abbruch): `SIGTERM` beendet die Nachkommen | `m-p19-term-kinder-bleiben` | `p19-sigterm` |
| Punkt 13 (Abbruch): `SIGTERM`: Verzeichnis erst nach den Nachkommen entfernt | `m-p19-term-verzeichnis-zu-frueh` | `p19-sigterm` |
| (6) der Token steht in keinem Argument von `gh` | `m-p19-token-im-argument-gh` | `p19-token-nicht-im-argument` |
| (6) der Token steht in keinem Argument von `git` | `m-p19-token-im-argument-git` | `p19-token-nicht-im-argument` |
| (6) der Token steht nicht in der Ausgabe bei Erfolg | `m-p19-token-im-log-bei-erfolg` | `p19-token-nicht-im-log` |
| (6) der Token steht nicht in der Ausgabe bei Fehler | `m-p19-token-im-log-bei-fehler` | `p19-token-nicht-im-log` |
| (6) der Token steht nicht in der Konfiguration des Klons | `m-p19-token-in-config` | `p19-token-nicht-in-url-und-config` |
| (6) der Token steht in keiner Datei | `m-p19-token-in-datei` | `p19-token-nicht-in-url-und-config` |
| (6) der Token steht nicht in der URL des Klons (Beispiel-Mutation Punkt 19 (9)) | `m-p19-token-in-url` | `p19-token-nicht-im-argument`, `p19-token-nicht-in-url-und-config` |
| (5) ohne `RELEASE_TAP_URL` klont `git` die Vorgabe-URL | `m-p19-url-vorgabe-falsch` | `p19-tap-url-vorgabe` |
| Punkt 13: der Schritt verlangt keine Marke eines Vorgängers (alle Fälle außer den bis dahin endenden rot) | `m-p19-verlangt-marke` | 83 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (4) die Version in der Formel ohne `v` | `m-p19-version-mit-v` | 7 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (4) die Version steht an jeder Stelle des Platzhalters | `m-p19-version-nur-erste` | 7 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (2c) Vorabversion lässt den Tap unverändert (Erkennung entfällt) | `m-p19-vorabversion-nicht-erkannt` | `p19-reihenfolge-gruende-c-d`, `p19-vorabversion` |
| (2c) jede Vorabversion, nicht nur `-rc` | `m-p19-vorabversion-nur-rc` | `p19-vorabversion` |
| (2c) kein Aufruf von `gh` bei einer Vorabversion | `m-p19-vorabversion-ruft-gh` | `p19-reihenfolge-gruende-c-d`, `p19-vorabversion` |
| (2c) kein Aufruf von `git` bei einer Vorabversion | `m-p19-vorabversion-ruft-git` | `p19-reihenfolge-gruende-c-d`, `p19-vorabversion` |
| (4) Text der Vorlage, kein Windows | `m-p19-vorlage-mit-windows` | `p19-formel-kein-windows`, `p19-formel-text`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Reihenfolge: `vorlage` vor dem Klon des Taps | `m-p19-vorlage-nach-klon` | `p19-vorlage-fehlt`, `p19-vorlage-leer`, `p19-vorlage-verzeichnis`, `p19-vorlage-vor-tap` |
| (4) Name in der Adresse ohne Version (Punkt 19 (4)), je Plattform | `m-p19-vorlage-name-mit-version-darwin-amd64` | `p19-formel-text`, `p19-formel-url-darwin-amd64`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Name in der Adresse ohne Version (Punkt 19 (4)), je Plattform | `m-p19-vorlage-name-mit-version-darwin-arm64` | `p19-formel-text`, `p19-formel-url-darwin-arm64`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Name in der Adresse ohne Version (Punkt 19 (4)), je Plattform | `m-p19-vorlage-name-mit-version-linux-amd64` | `p19-formel-text`, `p19-formel-url-linux-amd64`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Name in der Adresse ohne Version (Punkt 19 (4)), je Plattform | `m-p19-vorlage-name-mit-version-linux-arm64` | `p19-formel-text`, `p19-formel-url-linux-arm64`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) fehlende, leere oder nicht lesbare Vorlage: Exit 1 `vorlage` (die Prüfung unterscheidet die Zustände nicht, akzeptiertes Negativ nach Punkt 19 (9)) | `m-p19-vorlage-pruefung-entfaellt` | `p19-vorlage-fehlt`, `p19-vorlage-leer`, `p19-vorlage-verzeichnis`, `p19-vorlage-vor-tap` |
| (4) Text der Vorlage, Summe je Plattform | `m-p19-vorlage-sha-vertauscht` | `p19-formel-summe-darwin-arm64`, `p19-formel-text`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) das Stichwort der Vorlage-Prüfung ist `vorlage` | `m-p19-vorlage-stichwort-tap-formel` | `p19-vorlage-fehlt`, `p19-vorlage-leer`, `p19-vorlage-verzeichnis`, `p19-vorlage-vor-tap` |
| (4) Text der Vorlage (Zeile `homepage`) | `m-p19-vorlage-text-geaendert` | `p19-formel-text`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Text der Vorlage (Zeile `version`) | `m-p19-vorlage-version-mit-v` | `p19-formel-text`, `p19-formel-version`, `p19-tap-gleich-aktuell`, `p19-vorlage-text` |
| (4) Reihenfolge: `summen` vor `vorlage` | `m-p19-vorlage-vor-summen` | `p19-vorlage-nach-summen` |
| (5) genau ein Commit je Nachtrag | `m-p19-zwei-commits` | 8 Fälle, alle außer den aus der Liste der Mutation ausgenommenen |
| (5) Tap auf einem anderen Zweig als `main`: Exit 1 `tap-zweig` | `m-p19-zweig-nicht-geprueft` | `p19-tap-zweig-anders` |

Eine Mutation mit vielen roten Fällen (`summen-ohne-dir`, `formula-nicht-angelegt`, `verlangt-marke`) bricht den Schritt vor dem Tap; ihr Beleg ist die Liste der Fälle, die grün bleiben müssen (`ausser_h` in der Gegenprobe), und die Zusage des jeweiligen Elements trägt der Fall in der Spalte der benannten Mutation darüber.

**Grenzen und Einordnung**

- Nicht gefahren: der volle Lauf von `make release-gegenprobe` (Nutzerentscheid), die übrigen Gruppen (`p1-*` außer `m-p1-homebrew`, `p2-*` bis `p16-*`) und die Mutationen der Gruppen `p1` bis `p16`; die Änderungen dieses Laufs berühren sie nicht (Zeilen im Schritt `homebrew`, in den Fällen der Gruppe `p19-*` und in einer Zeile von `harness/README.md`).
- `m-p19-vorlage-pruefung-entfaellt` färbt die drei Zustände gemeinsam: Eine Prüfung unterscheidet sie nicht (ein Satz Code, ein Zustand je Fall ohne eigene Stelle); das ist das akzeptierte Negativ aus `SPEC-067` Punkt 19 (9). Ein Zustand „nicht lesbar“ hat keinen eigenen Fall (als root liest jede Datei); das Verzeichnis ist der Zustand, in dem `cat` scheitert, und liegt an derselben Stelle.
- `m-p19-token-in-url` trägt den Token nur in die URL des Probe-Taps (Pfad `/…`), nicht in die Vorgabe-URL mit `https://`: der Fall `p19-tap-url-vorgabe` prüft die Vorgabe über `url.<…>.insteadOf` und bliebe sonst rot. Die Zusage „der Token steht nicht in der URL“ belegt das für den Klon des Probe-Taps; die Form `https://x-access-token:<Token>@github.com/…` ist nicht mutiert.
- Kein grüner Mutant in den gefahrenen Läufen außer den akzeptierten Negativen; keiner wurde als äquivalent eingeordnet.
- `make gates` am Endstand der Dateien außer diesem Abschnitt (Kopf `a0951261` bis `72b6e839` plus die Korrektur der Kennung): grün, 303 s. Davor zwei rote Läufe: der erste endete in `docs-check` (die Kennung `ADR-0054` stand in der Tabelle ohne Link, behoben), der zweite in `test-integration` (`TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild`, `DROP DATABASE … timeout: context deadline exceeded`, Last 10 bis 12 im Container; die Änderungen berühren keinen Go-Code und keinen Test des Binaries, das ist die Klasse `BEO-REPO/test-faellt-unter-cpu-last-aus`); der dritte Lauf war grün.

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
