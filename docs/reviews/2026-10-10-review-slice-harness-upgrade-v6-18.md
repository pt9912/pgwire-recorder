# Review-Report: slice-harness-upgrade-v6-18 — 2026-10-10

**Review-Art:** Code-Review (Harness- und Dokumentations-Diff, kein Produkt-Code). Geprüft wird gegen Plan §1, §3 und §6 von `slice-harness-upgrade-v6-18`, gegen die Entscheidungen des Architect in §6, gegen `AGENTS.md` §3.3, §3.9, §3.11, §3.12 und §3.13 und gegen den Reviewer-Skill. Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** Diff `5e25195..HEAD`: `954fc1b` (Architect, Randformen in §6), `420f955` (Implementer, Baumtausch `v6.16.0` → `v6.18.0`, Symlinks, `AGENTS.md` §1, `harness/conventions.md`), `35503be` (Belege und Abgleich in §7), `0e61bed` (Folgelauf-Form im Skill und in `implement-slice`), `ddb288d` (Gate-Stand in §7). Rahmen: [MR-000](../../harness/conventions.md#mr-000--baseline-aussage).

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-harness-upgrade-v6-18` am Stand `ddb288d`, ganz gelesen, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen; die Zusagen in §7 habe ich nicht übernommen, sondern nachgemessen.
- `AGENTS.md` §3; Rahmen [MR-000](../../harness/conventions.md#mr-000--baseline-aussage); die aufgelöste [MR-001](../../harness/conventions.md#mr-001).
- `v6.18.0` · `regelwerk/modul-10-review-harness.md` §Harness-Einordnung, `templates/docs/reviews/review-report.template.md`, `regelwerk/modul-05-planning-harness.md`, `regelwerk/grundlagen-harness-dateien.md`.
- Vorherige Findings am Gegenstand: Report zu `slice-doku-ist-stand` (F-570 bis F-575). Die Nummern dieses Laufs beginnen bei F-576.

**Ausgeführte Läufe** (Scratch außerhalb des Repos, Präfix `rev-upg-`):

- **Asset.** `lab-regelwerk.zip` des Tags `v6.18.0` per `gh release download` frisch geladen; sha256 `18e5443b…b8d9` gegen die `SHA256SUMS` des Releases: gleich. `diff -r` des entpackten `regelwerk/` und `templates/` gegen `.harness/baseline/v6.18.0/`: leer. Vendored `SHA256SUMS`: 54 Zeilen, Form `<sha256>  <pfad>`, `sha256sum -c` ohne Abweichung, 54 Dateien im Baum.
- **Gate und Gegenprobe.** `make baseline-verify`: `v6.18.0 OK — 54 Dateien`. In einer Kopie des Arbeitsbaums eine Zeile an `modul-00-einfuehrung.md` angehängt: `make baseline-verify` meldet die Abweichung und endet mit Fehler.
- **Reste.** `grep v6.16.0` über das ganze Repo: Treffer nur in `done/`-Plänen, Register-Belegen, `docs/reviews/`, der Zeile `MR-001`, drei Herkunftsnennungen in `open/slice-harness-gate-index-werkzeug-teil` und einer Drift-Log-Zeile der Roadmap (Chronik). `find -xtype l` leer, beide Symlinks unter `.claude/rules/` lösen auf `v6.18.0` auf. Pins in `AGENTS.md` §1 und `harness/conventions.md` (Stand, Datum, Asset-URL, Stand-Zeile „Kurs-Welle 161“, Messzeile) nennen `v6.18.0`.
- **Delta.** Alter Baum aus `5e25195` gegen neuen: 31 Dateien verschieden (26 Regelwerk, 5 Vorlagen), 23 unverändert; 22 Regelwerk-Dateien und 2 Vorlagen nur Versionszeile — die Zählung in §7 stimmt. Inhaltlich: `grundlagen-harness-dateien`, `modul-13`, `modul-10`, `.d-check.yml`, `review-report.template`, `harness/README.template` (gelesen).
- **Zusage gegen Bestand.** Kein Eintrag im Diff außerhalb von Baum, Symlinks, `AGENTS.md`, `harness/conventions.md`, Skill, `implement-slice` und dem Slice-Plan; kein `.go`-File, keine Datei unter `docs/reviews/`, `done/`, `observations/` oder `harness/conventions/done/`. `make gates` am Stand `ddb288d` zusätzlich gefahren (Ergebnis siehe Verdikt).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-576 | INFO | Der Dateiname der Reports weicht in der Rollen-Infix-Form (`<datum>-review-<slice-Kennung>.md`) von der Vorlagen-Form `<YYYY-MM-DD>-<slice-Kennung>.md` ab. Der Wortlaut von `modul-10` („die volle Slice-Kennung im Dateinamen“) verbietet den Infix nicht, die Entscheidung des Implementers („Infix bleibt, nur `-r2` nachgezogen“) trägt also; die Abweichung ist im Skill sichtbar und in §7 begründet, hat aber bis zur Closure keine Adresse im Register. | `v6.18.0` · `regelwerk/modul-10-review-harness.md` §Harness-Einordnung (Modul 10); `v6.18.0` · `templates/docs/reviews/review-report.template.md` | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-18.md` · „**Akzeptierte Abweichung:** der Infix bleibt“ | nein — kein Sensor liest den Namen; erst ein künftiges Modul `reviews` mit `match: name` würde es prüfen | Abweichung von der Vorlagen-Form nur im Plan-§7 begründet |
| F-577 | INFO | Zwei Folgelauf-Formen stehen nebeneinander: der Bestand führt zehn Dateien mit `folge-review-…` bzw. Suffix `-folge`, Skill und `implement-slice` nennen jetzt `-r2`, `-r3`. Die Vorlage sagt nur „etwa mit Suffix `-r2`“; `-r3` ist die Fortschreibung. Der Diff ändert den Bestand nicht, wie in §1 zugesagt. | Maintainability | `.harness/skills/reviewer.md` · „hängt `-r2`, `-r3` an die Kennung“ | nein | Zwei Namensformen für Folgeläufe im Verzeichnis |
| F-578 | INFO | `modul-10` führt seit `v6.18.0` den Absatz „Ein Deckungs-Sensor prüft nur, was er als Zusage erkennt“ und die Vorlage nennt die Zuordnung Report → Slice über die Kennung. Dieses Repo hat keinen Sensor, der Reports den Slices zuordnet; §7 benennt den Absatz als Lesestoff für den noch nicht angelegten Slice zum Modul `reviews`. Die Beobachtung hat bis zur Closure keine Adresse (kein Slice, noch kein Register-Eintrag). | `v6.18.0` · `regelwerk/modul-10-review-harness.md` §Harness-Einordnung (Modul 10) | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-18.md` · „Lesestoff für den geplanten Slice“ | nein | Ausgang für Befund ohne Adresse noch offen |
| F-579 | INFO | `420f955` tauscht den Baum und ändert im selben Commit Dateiinhalte (Versionszeilen, `AGENTS.md`, `conventions.md`). `AGENTS.md` §3.3 verlangt für Move plus Inhaltsänderung zwei Commits. Der Architect hat den einen Commit ausdrücklich entschieden (§6, damit `make baseline-verify` nie zwei Tag-Verzeichnisse sieht); die Rename-Erkennung hält (`git show -M`: alle 54 Pfade als Rename, niedrigster Wert `SHA256SUMS` 65 %, Regelwerk ab 96 %). Kein Verstoß, festgehalten als Messung. | `AGENTS.md` §3.3 | `AGENTS.md` · „Der Move-Commit bleibt rein“ | ja — `git show -M --summary 420f955` | Move und Inhalt im Commit, vom Architect entschieden |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.harness/baseline/v6.18.0/` gegen Release-Asset | geprüft, ohne Befund (Byte-gleich, `SHA256SUMS` stimmt, Gegenprobe rot) |
| `.harness/baseline/v6.16.0/` | geprüft, ohne Befund (entfernt, `make baseline-verify` sieht ein Tag-Verzeichnis) |
| `.claude/rules/` | geprüft, ohne Befund (Symlinks lösen auf, keine kaputten Symlinks im Repo) |
| `AGENTS.md`, `harness/conventions.md` Pins und Messzeile | geprüft, ohne Befund |
| Regelblock-Tabelle `harness/conventions.md` | geprüft, ohne Befund: alle 25 `modul-*`/`grundlagen-*`-Dateien des neuen Baums stehen in der Tabelle, keine neue, umbenannte oder entfallene Datei; Zellen `modul-10`, `modul-13`, `modul-15` und `grundlagen-harness-dateien` stimmen mit dem Delta überein (das Delta betrifft Gate-Index in Teilen, Modul `reviews`, Deckungs-Sensor — keine Zelle widerspricht) |
| Verkörperte Regeln `AGENTS.md` §3 gegen `modul-05` §Closure und Risiko-Ausgänge, `modul-09` Hard Rules, `modul-10` Harness-Einordnung, `grundlagen-harness-dateien` §Kommentar | geprüft, ohne Befund (Delta berührt keine dieser Stellen außer dem Dateinamen; kein Widerspruch zu §3.2 bis §3.13) |
| `.claude/commands/implement-slice.md`, `.harness/skills/reviewer.md` (`0e61bed`) | geprüft, ohne Befund: auf die Folgelauf-Form und den Dateinamen beschränkt, jede Aussage trägt der Wortlaut von `modul-10` und der Vorlage (§3.11); der Infix wird nicht als Vorlagen-Aussage ausgegeben |
| Freshness-Audit-Aussage in §7 | geprüft, ohne Befund: Release-Liste am Tag des Laufs, neuester Tag `v6.18.0` = Pin; `conventions.md` §Der mitgelieferte Baum altert still nennt keinen Tag und bleibt wahr |
| `AGENTS.md` §3.9 Plan folgt der Korrektur | geprüft, ohne Befund: §1, §3, §6 und Kopf stimmen mit dem Diff überein (kein Folge-Slice angelegt, kein Roadmap-Eintrag nötig; Drift-Log-Zeile vom Planner steht) |
| `AGENTS.md` §3.13 Adresse und Nachzählen | geprüft, ohne Befund: der Diff nennt keinen anderen Slice neu als Adresse; zwei Liefer-Punkte, zwei Schichten (Harness, Doku), Maß ≤ 3 und ≤ 2 eingehalten |
| Eingefrorene Dokumente (`docs/reviews/`, `done/`, Register-Belege, `harness/conventions/done/`) | geprüft, ohne Befund (im Diff nicht angefasst) |
| `.d-check.yml`, `tools/`, `.githooks/commit-msg`, `*.mk`, Go-Code | geprüft, ohne Befund (nicht im Diff; kein Sensor liest den Report-Namen) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Abweichung von der Vorlagen-Form nur im Plan-§7 begründet · Zwei Namensformen für Folgeläufe im Verzeichnis · Ausgang für Befund ohne Adresse noch offen · Move und Inhalt im Commit, vom Architect entschieden

## Verdikt

**Merge-blockierend:** nein — keine HIGH- und keine MEDIUM-Findings.

**Übergabe:** Findings gehen an den Implementer. F-576 und F-578 brauchen zur Closure je einen Ausgang (Register-Eintrag, kein Zähler gesetzt; der Slice zum Modul `reviews` ist nicht angelegt, also keine Zuweisung). F-577 und F-579 sind ohne Aktion. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt). Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
