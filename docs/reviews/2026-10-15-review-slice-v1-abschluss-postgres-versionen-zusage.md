# Review-Report: slice-v1-abschluss-postgres-versionen-zusage — 2026-10-15

**Review-Art:** Code — Diff gegen Plan, Entscheidungen des Architect (§6) und Hard Rules (Maintainability); die DoD prüft der Verifier.

**Gegenstand:** Slice `slice-v1-abschluss-postgres-versionen-zusage`, Diff `510e7ab..4b214bc` (Move-Commit nach `in-progress/` bis HEAD; Commits `7d4534f` Architect, `9458483` Spezifikation, Lastenheft, Handbuch, README, `4b214bc` Plan).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0 · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-postgres-versionen-zusage`, §1 bis §8 (§6: Entscheidungen des Architect vom 2026-10-15; Entscheidung des Nutzers vom 2026-10-15, Frage B: SCRAM-SHA-256 für PostgreSQL 14 bis 18, Klartext-Passwort und MD5 nur PostgreSQL 17), dazu §7 *Belege des Implementers*
- [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md)
- [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [LH-RB-02](../../spec/lastenheft.md#lh-rb-02--einsatzumgebung); `LH-FA-05.e`, `SPEC-029`, `SPEC-038` der Spezifikation
- `AGENTS.md` (Hard Rules, hier §3.7, §3.9, §3.11, §3.12, §3.13, §3.14)

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `make docs-check` ist am HEAD rot: `d-check` meldet `id-unlinked` für die Kennung [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md) in der Zeile *Gates* von §7, die Commit `4b214bc` einträgt. Dieselbe Zeile belegt „`make gates` exit 0“; der Stand, den sie belegt, ist nicht der HEAD. | Maintainability (Gate rot am Übergabestand); `AGENTS.md` §6 Schritt 8 | `docs/plan/planning/in-progress/slice-v1-abschluss-postgres-versionen-zusage.md` · "(nach Korrektur des Links" an der Kennung der Versionsmatrix-ADR in §6 dieses Plans (Zeile 205) | ja — `make docs-check` meldet genau diese Zeile (selbst gefahren: 769 Dateien, 1 Befund) | Kennung in Beleg-Prosa ungelinkt, Gate nach dem Lauf rot |
| F-2 | LOW | §7 *Prüfungen der Texte* sagt zu `grep -n "harness/vertraege\|make " spec/*.md`: „keine Treffer“. Das Muster trifft `spec/spezifikation.md` (`make a-check`, `make a-check-negativ`); beide Stellen sind Bestand, nicht im Diff. Für „kein Link auf `harness/vertraege/`“ und „kein `test-postgres`/`test-integration`“ ist die Aussage wahr. | `AGENTS.md` §3.11 (Beleg sagt nur, was die Prüfung zeigt) | Plan · "keine Treffer. Quantoren-Durchsicht" (§7) | ja — der genannte `grep` | Prüfbefund im Plan nennt Trefferstand falsch |
| F-3 | LOW | §6 sagt nach der Umsetzung noch „Die Stellen sind noch nicht geschrieben; der Implementer schreibt sie mit diesen Vorgaben“. Der Diff hat sie geschrieben (`9458483`). | `AGENTS.md` §3.9 (Plan folgt dem Diff), §3.7 | Plan · "Die Stellen sind noch nicht geschrieben" (§6) | nein | Plan-Prosa hinter dem Diff |
| F-4 | LOW | Die Zeile *Nachweis* von `LH-RB-02` beginnt mit „Abnahmeszenarien in §7 laufen unter diesen Bedingungen“; „diesen Bedingungen“ schließt jetzt die zugesagten Hauptversionen ein. Gelesen so, laufen alle Abnahmeszenarien je Hauptversion; Abnahmeszenario 18 und der dritte Satz von Abnahmeszenario 16 laufen nur gegen 17 (`SPEC-029` Punkt 2). Der Wortlaut ist der des Architect; der Plan nennt den Satz „die automatisierten Abnahmeszenarien laufen gegen jede von ihnen grün“ als weiterreichend als das Gate. | `AGENTS.md` §3.11 | `spec/lastenheft.md` · "Abnahmeszenarien in §7 laufen unter diesen Bedingungen; die" | nein — die Lesart prüft kein Gate | Zusage im Lastenheft liest weiter als das Gate |
| F-5 | INFO | Der Beleg für „`psql` aus libpq 17“ im Handbuch ist ein Skript im Scratchpad des Laufs, nicht im Repo. Selbst gefahren (Skript des Implementers, `postgres:17-alpine` per Digest, `make build`): `psql` 17.11, vier Anmeldungen (scram-sha-256, md5, password, trust) liefern ihre Zeile, 4 `sql:`-Einträge, 0 Treffer für Passwörter. Die Angabe stimmt; ein Verifier findet den Lauf im Repo nicht. Rolle: Verifier. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · "`psql` aus libpq 17; andere Versionen der" | nein — kein Repo-Ziel fährt `psql` gegen den Recorder | Beleg nur im Scratchpad |
| F-6 | INFO | `spec/lastenheft.md` nennt mit „die die Spezifikation nennt“ erstmals das Spec-Stratum beim Namen. Das ist der Wortlaut des Architect und nennt keine Datei, keine ADR, keinen Slice; `make docs-check` meldet dazu nichts. Keine Aktion. | Hard Rule `AGENTS.md` §3.4 / Reviewer-Skill (Spec-Stratum nennt Artefakt) | `spec/lastenheft.md` · "die die Spezifikation nennt." | ja — `make docs-check` | Stratum-Name im Lastenheft |
| F-7 | INFO | Bestand außerhalb des Diffs: `spec/spezifikation.md` nennt `make a-check` und `make a-check-negativ` (Abschnitt Zeit-Port). Der Architect-Beschluss „kein `make`-Ziel in der Spezifikation“ gilt für diesen Slice; der Diff fügt keines hinzu. Kein Befund gegen den Diff. | Maintainability | `spec/spezifikation.md` · "Das Architektur-Gate (`make a-check`)" | ja — `grep -n "make " spec/*.md` | Bestand widerspricht Architect-Prinzip |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/spezifikation.md` (`SPEC-029` Tabelle und Absatz, `LH-FA-05.e`, `SPEC-038`, Historie): Liste 14 bis 18, Referenz 17, Vorabversion und Nebenversion, „unterstützt“ Punkt 1 und 2 gegen Architect-Beschluss | geprüft, ohne Befund (Punkt 1 deckt sich mit `POSTGRES_TESTS_PHASE1/2` im Fragment der Matrix; Punkt 2 „nur dort“ stimmt: Matrix fährt weder `md5`/`password` noch Abnahmeszenario 18 oder 16) |
| `spec/` ↔ `harness/vertraege/`: kein Link, kein Verweis, kein Gate-Name, kein `make`-Ziel, keine Testnamen im Diff | geprüft, ohne Befund (`grep vertraege\|SPEC-053\|ADR-00 spec/*.md` leer; Treffer `make` nur Bestand, F-7) |
| `spec/lastenheft.md` (`LH-FA-05`/`LH-RB-02`-Zeile): Draft, Version 0.1.0 unverändert, keine Historie-Zeile, Aussage entspricht dem Wortlaut des Architect | geprüft, ohne Befund (Einschränkung F-4) |
| `docs/user/benutzerhandbuch.md`: Ist-Zustand, keine Verweise auf Spezifikation, ADR, Slices, Reviews; Version 0.9 und Stand 15.10.2026 passen; vier Stellen plus *Voraussetzungen* tragen SCRAM-SHA-256 14 bis 18, Klartext und MD5 17 | geprüft, ohne Befund |
| „geprüft mit PostgreSQL 17“ ohne Verfahren in Handbuch und README | geprüft, ohne Befund (keine Fundstelle) |
| `README.md`: Satz zur Anmeldung nennt je Verfahren dieselben Versionen wie das Handbuch | geprüft, ohne Befund |
| Angabe „pgx 5.11“: `go.mod` führt `github.com/jackc/pgx/v5 v5.11.0`; `TestE2ERecordAnmeldung` nutzt `pgconn` daraus | geprüft, ohne Befund |
| Belege selbst gefahren, nacheinander: `make test-postgres-versionen` (14, 15, 16, 18 grün, 22 s) und `make test-integration` (exit 0, 127 PASS-Zeilen auf oberster Ebene, 0 FAIL, `TestE2ERecordAnmeldung/{scram-sha-256,md5,password}` PASS); Container und Netze danach entfernt | geprüft, ohne Befund |
| Größe nach Metrik §4: `git diff --numstat 510e7ab..HEAD -- spec docs/user README.md` = 3 + 24 + 6 + 23 = 56; Meldeschwelle 300 und Marke 400 nicht berührt | geprüft, ohne Befund (behauptete 56 stimmt) |
| Randformen (§3.12): alle 13 Entscheidungen des Architect stehen im Commit `7d4534f` vor dem Spezifikations-Commit; `4b214bc` ändert in §6 nur Link und Zusatz, keine neue Randform | geprüft, ohne Befund |
| Adressen und Nachzählung (§3.13): keine neue Adresse im Diff; §8 führt die Zählung (drei Liefer-Punkte, zwei Schichten); Register-Eintrag `BEO-REPO/postgres-versionsliste-altert` liegt vor | geprüft, ohne Befund |
| Kopf des Plans (`Berührte Spec-Stellen`) und `make kopf-check` | geprüft, ohne Befund |
| Core-Reinheit, Produkt-Code, Gates, ADRs: im Diff unberührt | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kennung in Beleg-Prosa ungelinkt, Gate nach dem Lauf rot · Prüfbefund im Plan nennt Trefferstand falsch · Plan-Prosa hinter dem Diff · Zusage im Lastenheft liest weiter als das Gate · Beleg nur im Scratchpad · Stratum-Name im Lastenheft · Bestand widerspricht Architect-Prinzip

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): `make docs-check` ist am HEAD rot, obwohl §7 `make gates` grün belegt. Die Änderung liegt auf der Positivliste des Stop-Hooks (`docs/plan/planning/`), der Hook gibt sie ohne neuen vollen Lauf frei; das Gate fängt es erst im nächsten Lauf. F-2 bis F-4 sind LOW, F-5 bis F-7 INFO.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
