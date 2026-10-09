# Review-Report: slice-v1-abschluss-upstream-verbinden — 2026-10-09

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3 und §6, gegen die Spezifikation und gegen die Hard Rules (Modul 10). Die Spezifikation ist `LH-FA-17.a` mit *Benannte Verbindungen*, *Geheimnisse* und *Fehler*, dazu A7 bis A9, A12, F-521, F-528 und U1 bis U8 am Stand `9a04b12`. Hinzu kommen [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) und [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers. Die Wirkung einer Verbindung bei `play` liegt bei `slice-v1-abschluss-einspielen` und ist kein Befund.

**Gegenstand:** fünf Code-Commits und ein Plan-Commit.

- `fe849ac`: Prüfung von `--upstream` und `PGWIRE_RECORDER_UPSTREAM` nach der Zusammenführung, dazu die Hilfe.
- `694eb99`: benutzte Verbindung, `sslmode=require`, Einsetzen, `PGR-E2005`, Port, Zusammensetzen.
- `e35fb56`: zwei E2E-Tests.
- `8ba8caa`: Handbuch §5 und §7.
- `80cb896`: Testfall Leerraum.
- `74eb4b8`: §7 *Belege des Implementers*.

Neu sind `upstream.go`, `upstream_test.go`, `einsetzen_test.go` und `test/integration/verbindung_e2e_test.go`. Geändert sind `cli.go`, `datei.go`, `verbindung.go`, `export_test.go`, `leser_test.go`, `verbindung_test.go` und die Abdeckungstabellen. Der Diff `9a04b12..80cb896` hat 712 Zeilen hinzu und 24 entfernt.

**Skill:** `.harness/skills/reviewer.md` @ `32afc69`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-upstream-verbinden.md`: ganz gelesen am Stand `74eb4b8`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-17.a`: *Konfigurationsdatei*, *Benannte Verbindungen*, *Geheimnisse* und *Fehler*. `spec/` ist im Diff `9a04b12..74eb4b8` unverändert.
- [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md), [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md).
- `AGENTS.md` §3.3, §3.7 und §3.9 bis §3.13.
- `internal/bootstrap/bootstrap.go`: Übergabe an `postgres.Upstream`, Zeile beim Start.
- Vorheriger Report am selben Modul: `slice-v1-abschluss-verbindungen-platzhalter` (F-519 bis F-529). Die Nummern dieses Laufs beginnen bei F-530.

**Ausgeführte Läufe:**

- **Arbeitsweise.** Je Lauf eine frische Kopie per `git archive 74eb4b8 go.mod go.sum cmd internal | tar -x` unter dem Scratch-Pfad, ohne `cp -p`. `go test -count=1 ./internal/adapters/driving/cli/` lief im Image der Stufe `deps` aus dem `Dockerfile`, netzlos, mit einem Bind-Mount der Kopie. Im Repo habe ich nichts geändert. Gelöscht habe ich nur benannte Pfade: je Mutant sein Verzeichnis `upvmut-<X>`, dazu `base`, `probe_test.go`, `probe2_test.go`, `upv_mut.py` und `acheck.log`, außerdem das eigene Image `pgr-rev-upv:deps`.
- **Sonde.** Eine Testdatei im Paket `cli_test` mit 15 Fällen, Ergebnisse unter den Negativbefunden und in F-534.
  - `GEHEIM` steht in Benutzer, Passwort und Datenbank, und zwar für `sslmode=require`, für `PGR-E2005` und für den Port nach dem Einsetzen (`GEHEIMP`, `9GEHEIMP`). Dazu kommt eine gültige Adresse.
  - Der zusammengeführte Wert `host:port` liegt über einem Schlüssel `upstream`, der auf eine Verbindung mit `sslmode=require` zeigt. Geprüft mit `--upstream`, mit der Umgebungsvariable und mit beiden.
  - Eine fehlende Pflichtoption trifft auf einen ungültigen Wert der Umgebungsvariable.
  - Eingesetzter Host `a]b` und `[::1]`, eingesetzter Port mit 22 führenden Nullen und mit der arabischen Ziffer `٥`, dazu `replay` mit Datei.
- **Mutanten.** Je Mutant eine frische Kopie. Das Skript ersetzt genau eine Stelle und bricht bei einer anderen Trefferzahl ab.
  - A: Klammern bei mehr als einem `:` statt bei einem. rot (`TestUpstreamEinsetzen`)
  - C: Umgebung nur ohne Kommandozeile geprüft. rot (`TestUpstreamUngueltig`, `TestUpstreamReihenfolgeAmEnde`)
  - D: leere Variable eingesetzt. rot (drei Tests)
  - E: Port nach dem Einsetzen nur auf Ziffern geprüft, ohne Bereich. rot (`TestUpstreamPortNachEinsetzen`)
  - F: `sslmode=require` nicht geprüft. rot (`TestUpstreamReihenfolgeAmEnde`)
  - Q: `zuletzt` in der Schleife der Zusammenführung, vor den späteren Pflichtoptionen. rot (`TestUpstreamReihenfolge`)
  - T: fehlende Variable in Host oder Port durch `5` ersetzt. rot (zwei Tests)
  - R: Verbindung aus `g.cli` gesucht, wenn die Kommandozeile gesetzt ist. **grün**, äquivalent (Negativbefunde)
  - S: `datei.verbindung` vergleicht ohne Rücksicht auf Groß- und Kleinschreibung. **grün** (F-534). Eine zweite Kopie mit Mutant und Sonde zeigt eine Eingabe, die Original und Mutant unterscheidet.
- **Architektur-Gate:** `make a-check` Exit 0, `gesamt: 0 Befund(e)`.
- `make docs-check` vor dem Commit.
- Nicht gefahren: `make gates`, `make lint`, `make test-integration`, Integrations-Mutanten. Von den 38 Mutanten aus §7 habe ich nur die gefahren, die mit C, D, F und Q zusammenfallen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-530 | LOW | Laut Kommentar an `SetzeEin` ist ein Teil, den die URL nicht schreibt, `""`. Für den Port stimmt das nicht. `TestEinsetzenAlleTeile` erwartet für `postgresql://${PGR_T_H}/${PGR_T_DB}` den Port `"5432"`, weil das Laden ohne Port `5432` ablegt. Der Kommentar sagt also etwas zu, was der Code nicht tut. | `AGENTS.md` §3.7, §3.11 | `internal/adapters/driving/cli/export_test.go` · „ein Teil, den die URL nicht schreibt, ist ""“ (Z. 84) | ja (`TestEinsetzenAlleTeile`, Fall *ohne Benutzer*) | Kommentar sagt etwas zu, was der Code nicht tut |
| F-531 | LOW | Laut den Kommentaren an `option.zuletzt` und `lies` *prüft* `zuletzt` nach der Zusammenführung. Der einzige Träger, `upstreamRecord`, schreibt aber auch das Ergebnis: Er ersetzt `c.Record.Upstream` durch die zusammengesetzte Adresse. Wer am Leser ändert und nur die Kommentare von `option` und `lies` liest, hält `zuletzt` für eine reine Prüfung. Er übersieht damit, dass die Reihenfolge der `zuletzt`-Aufrufe auch bestimmt, welchen Wert ein späterer Aufruf sieht. | `AGENTS.md` §3.7 (Kopplung) | `internal/adapters/driving/cli/cli.go` · „zuletzt prüft die“ (Z. 194); · „Zuletzt prüft es je“ (Z. 330); `upstream.go` · `c.Record.Upstream = adresse` (Z. 38) | nein (Lesen) | Kommentar verschweigt eine Wirkung |
| F-532 | INFO | Laut `LH-FA-17.a` *Fehler* gilt am Ende die Reihenfolge Pflichtoptionen, Kombinationen, `--upstream`. Geliefert ist: Alle Pflichtoptionen stehen vor `zuletzt`, das ist mit Q belegt. Kombinationen prüft `record` heute nicht, denn `--tls-cert`, `--tls-key` und `--allow-plaintext` kennt der Leser nicht (grep ohne Treffer). Damit ist die Reihenfolge erfüllt. Die `zuletzt`-Aufrufe laufen in der Reihenfolge der Optionstabelle. Hängt eine spätere Kombinationsprüfung an einer Option, die in der Tabelle nach `upstream` steht, käme sie nach `--upstream`. Das ist ein Hinweis an den Architect, keine Adresse. | `LH-FA-17.a` *Fehler* | `internal/adapters/driving/cli/cli.go` · `if err := o.zuletzt(&cmd, stand[i], d); err != nil {` (Z. 375) | ja (Mutant Q) | — (Hinweis an den Architect) |
| F-533 | INFO | Das Handbuch §5 nennt `play` an drei Stellen. Neu formuliert sind das Passwort aus dem Platzhalter (*„Passwörter geben Sie als Platzhalter“*) und *„Beim Einspielen verbindet `require` verschlüsselt“*. Unverändert stehen geblieben ist *„gilt `PGWIRE_RECORDER_PASSWORD`“*. Keine dieser Zusagen liefert einer der beiden Slices, `play` gibt es noch nicht. Das Handbuch beschreibt `play` auch an anderen Stellen, und Plan §3 nennt für das Handbuch den *Zielstand*. Gegen den Plan ist das deshalb kein Befund. DoD-Punkt 3 verlangt dagegen *„wie geliefert“*. Wie weit diese Zusagen dazu passen, ist eine Frage an den Verifier. Was das Handbuch zu `record` zusagt, stimmt mit dem Code überein (Negativbefunde). | Plan §3, Zeile `docs/user/benutzerhandbuch.md`; Plan §6 *Risiken*, dritter Punkt | `docs/user/benutzerhandbuch.md` · „Beim Einspielen verbindet `require` verschlüsselt“ (Z. 650); · „Passwörter geben Sie als Platzhalter“ (Z. 654) | nein | — (Hinweis an den Verifier) |
| F-534 | MEDIUM | Die Zusage, dass die benutzte Verbindung die ist, deren Namen der zusammengeführte Wert genau nennt, hält beim Auswählen der Verbindung keine Mutation. Mutant S ersetzt in `datei.verbindung` den Vergleich `v.name == name` durch `strings.EqualFold`, und alle Tests bleiben grün. Die Zeile *„Name nur bei genauer Übereinstimmung“* in §7 belegt nur die Prüfung über `hatName`, nicht die Auswahl. Die Sonde unterscheidet die beiden: Die Datei führt `Staging: postgresql://gross/db` vor `staging: postgresql://klein/db`, beide Namen sind gültig. Mit `--upstream=staging` liefert das Original `klein:5432`, der Mutant `gross:5432`. *Failure-Szenario:* Ein Umbau der Suche, etwa über eine Abbildung mit kleingeschriebenen Schlüsseln, lässt `record` zum falschen Server verbinden. Es bleibt bei Exit 0, die Zeile beim Start nennt die falsche Adresse, und kein Gate wird rot. Die Klasse war in den vorigen Reports zweimal LOW und einmal MEDIUM (F-512, F-513, F-519). | `LH-FA-17.a` *Benannte Verbindungen* (*nur bei genauer Übereinstimmung, auch in Groß- und Kleinschreibung*; *Benutzt ist die Verbindung, deren Namen …*); `AGENTS.md` §3.10 | `internal/adapters/driving/cli/datei.go` · `if v.name == name {` | ja (Mutant S, Sonde, `go test`) | Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten |
| F-535 | INFO | Zwei Prüfungen von `TestUpstreamEinmalGelesen` sind unterschiedlich stark. Die erste prüft, dass die Adresse schon im Ergebnis von `Parse` eingesetzt ist (`erst:5432`). Ein Einsetzen erst im Bootstrap würde sie rot machen. Die zweite (`Setenv` danach) kann nicht rot werden, weil `RecordOptions.Upstream` ein Text ist. §6 U2 und §7 *Grenzen* nennen das offen. Die Abdeckungszeile sagt nur zu, was gilt. | Plan §6 U2; `AGENTS.md` §3.11 | `internal/adapters/driving/cli/einsetzen_test.go` · `func TestUpstreamEinmalGelesen` (Z. 221) | nein | — (Hinweis an den Verifier) |
| F-536 | INFO | Größe: Der Diff war bequem in einer Review-Sitzung prüfbar. Er hat 712 Zeilen bei einer Schätzung von 660 bis 870. Der Code umfasst 141 Zeilen in drei kleinen Funktionen, die Tests sind tabellengetrieben, und die Zusagen sind in §7 je Zeile belegt. Aus Sicht des Reviews ist das erste Risiko aus §6 nicht eingetreten. | Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice; Plan §4, §6 | Plan §6 · „**660 bis 870 Zeilen**“ | nein | — (Hinweis an den Verifier) |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Prüfreihenfolge am Ende (`LH-FA-17.a` *Fehler*, A9) | geprüft, ohne Befund. `zuletzt` läuft erst nach der Schleife der Zusammenführung, also nach allen Pflichtoptionen. Das hält auch für `--output`, das in der Tabelle nach `upstream` steht (Mutant Q rot). In `upstreamRecord` kommt zuerst die Kommandozeile, dann die Umgebungsvariable. `adresseRecord` prüft danach in dieser Reihenfolge: `sslmode=require`, dann das Einsetzen von Host und Port (`PGR-E2005`), dann die Form des Ports. Laut Sonde geht eine fehlende Pflichtoption einer ungültigen Umgebungsvariable vor, und diese geht `sslmode=require` vor. Kombinationen: F-532. |
| Option und Umgebung (A7, A8, F-521, U1) | geprüft. Befund F-534. Beide Werte werden geprüft, auch der, der nicht gilt (Mutant C rot). Sie laufen über dieselbe Regel `upstreamGueltig` wie der Schlüssel. Namen kommen nur aus `datei.namen`, und eine nil-Datei hat keine. Beide Meldungen sind konstant, nennen nur `--upstream` bzw. den Namen der Umgebungsvariable und keinen Wert. `TestUpstreamUngueltig` vergleicht sie genau, mit `GEHEIM` in jedem Fall. Laut Sonde wird eine Verbindung mit `sslmode=require`, auf die der Schlüssel `upstream` zeigt, nicht benutzt, wenn Option oder Umgebung einen anderen Wert setzen. Mutant R ist äquivalent: Die Kommandozeile geht immer vor, also ist `c.Record.Upstream` bei gesetzter Kommandozeile gleich `g.cli`. Die Mutanten C, D, F und Q fallen mit Zeilen der Tabellen in §7 zusammen und sind rot, wie §7 sagt. Auswahl der Verbindung: F-534. |
| Einsetzen (*Geheimnisse*, A12, Rückgabe 8, U8) | geprüft, ohne Befund. Das Einsetzen läuft in einem Durchlauf über die Stücke. Der Wert wird unverändert geschrieben: weder dekodiert noch erneut ausgewertet noch gekürzt, auch mit `@ : / %`, `${` und Leerraum (`TestUpstreamEinsetzen`). Ein leerer Wert gilt als nicht gesetzt (Mutant D rot), und gemeldet wird die erste fehlende Variable in der Reihenfolge der Teile und Stücke. Den Port prüft `portForm` nach dem Einsetzen. Das schließt Überlauf (`99999999999999999999`), Nicht-Ziffern, `٥`, `-5`, Leerraum und Steuerzeichen ein (Mutant E rot, Sonde). 22 führende Nullen sind gültig, wie A4 es verlangt. Einen eingesetzten Host prüft der Start nicht. |
| Klammerung (U5, F-528) | geprüft, ohne Befund. `zusammensetzen` setzt Klammern genau dann, wenn der Host ein `:` enthält (Mutant A rot, Fall `a:b`). Laut Sonde wird `a]b` ohne Klammern zusammengesetzt und das eingesetzte `[::1]` zu `[[::1]]:5432`. Beides folgt aus derselben Regel. Eine Zone bleibt als `%eth0` stehen. |
| U6, Meldungen und Zeile beim Start | geprüft, ohne Befund. In der Sonde steht `GEHEIM` in Benutzer, Passwort und Datenbank. `sslmode=require`, `PGR-E2005` und der Port nach dem Einsetzen nennen dann nur `connections.v` und, bei `PGR-E2005`, den Namen der Variable. Die gültige Adresse ist `h:7`. Die Zeile beim Start schreibt `bootstrap.go` aus `o.Upstream`. `TestE2ERecordVerbindungAusDatei` prüft den ganzen stderr gegen `staging`, `NUTZERX`, `GEHEIMPW` und `DBNAMEX`. Für `PGR-E4002` belegt `TestE2ERecordVerbindungIPv6` die Adresse. Benutzer, Passwort und Datenbank kann `postgres.Upstream` nicht nennen, denn es erhält nur `Address`. |
| Hexagon, [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) | geprüft, ohne Befund. `internal/bootstrap`, `internal/hexagon`, `internal/adapters/driven` und `internal/adapters/driving/pgwire` sind im Diff unverändert. `upstream.go` importiert nur `errors`, `os` und `model`. `make a-check` endet mit Exit 0 und `0 Befund(e)`. Den Hinweis auf zehn Dateien ohne Schicht kennt der Bestand schon für `test/integration/`, die neue E2E-Datei liegt dort. |
| [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) | geprüft, ohne Befund. Der neue Code nutzt keinen YAML-Typ. `PGR-E2005` kommt über `model.CodeConfigVariable`, `PGR-E2004` über `fehlerDatei`, `PGR-E2001` über `model.CodeUsage`. Die Code-Tabelle ist unberührt. |
| Handbuch §5 und §7 zu `record` | geprüft. Befund F-533. Die Zusagen zu `record` stimmen mit dem Code überein: Form der URL, Name nur aus der gewählten Datei, Prüfung jedes Werts, nur Host und Port, Klammern bei `:`, die Zeile beim Start, `sslmode=require` als `PGR-E2004`, einmaliges und unverändertes Einsetzen, leere Variable, `PGR-E2005` und der Port. Ebenso die Zeilen zu `PGR-E2005` und `PGR-E2006`. |
| Tests `upstream_test.go`, `einsetzen_test.go`, `leser_test.go`, `verbindung_test.go`, E2E | geprüft. Befunde F-530 und F-535. Die Fälle vergleichen die Meldungen genau. Damit ist „ohne Wert“ in jedem Fall mitgeprüft. Die Abdeckungs-Deklarationen sagen nicht mehr zu, als die Fälle prüfen. `basis` in `leser_test.go` ändert nur den Beispielwert. |
| Plan §1, §3, §6; `AGENTS.md` §3.9, §3.12, §3.13 | geprüft, ohne Befund. Kein Code-Commit ändert den Plan (`git diff 9a04b12..80cb896 -- docs/plan` ist leer). §3 sagt „keine Änderung erwartet“ für den Bootstrap, und so ist es geliefert. Der Code entscheidet keine Randform, die §6 nicht nennt. Die Fälle `a]b` und `[[::1]]` folgen aus U5, die 22 Nullen aus A4. Dass die Meldung zur Kommandozeile den Wert nicht nennt, ist innerhalb des „darf“ aus *Fehler* gewählt. Eine neue Adresse nennt der Diff nicht. |
| Hard Rule 3.3, Kommentare (3.7) | geprüft. Befunde F-530 und F-531. In den sechs Commits gibt es keinen Move. Die übrigen neuen Kommentare sind Zusagen oder Kopplungen. |
| Commit-Messages | geprüft, ohne Befund. Alle sechs nennen `slice-v1-abschluss-upstream-verbinden` und `LH-FA-17`, keine nennt eine `SPEC-` oder `ARC-`Kennung. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Summary:** 0 HIGH · 1 MEDIUM · 2 LOW · 4 INFO.

- F-530: Der Kommentar an `SetzeEin` sagt `""` für einen nicht geschriebenen Port; es ist `5432`.
- F-531: Die Kommentare an `zuletzt` nennen eine Prüfung und verschweigen, dass `zuletzt` die Adresse schreibt.
- F-532: Die Reihenfolge Pflichtoptionen, Kombinationen, `--upstream` ist erfüllt, weil `record` keine Kombination prüft.
- F-533: Das Handbuch nennt `play` im Zielstand; ob das zu DoD-Punkt 3 passt, entscheidet der Verifier.
- F-534: Die Auswahl der Verbindung nach dem genauen Namen hält keine Mutation (Mutant S grün, Sonde unterscheidet).
- F-535: `TestUpstreamEinmalGelesen` trägt nur in seiner ersten Hälfte.
- F-536: Der Diff war in einer Sitzung prüfbar.

Wiederkehrende Klassen: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-534), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-530).

**Finding-Klassen dieses Laufs:** Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten (vierter Lauf in Folge, deshalb MEDIUM) · Kommentar sagt etwas zu, was der Code nicht tut · Kommentar verschweigt eine Wirkung

## Verdikt

**Merge-blockierend:** nein für das gelieferte Verhalten. Die Schwerpunkte halten am Code und an den Tests: die Prüfreihenfolge am Ende, beide Werte von `--upstream` ohne Wert in der Meldung, das unveränderte Einsetzen mit der ersten fehlenden Variable, der Port nach dem Einsetzen, die Klammerung genau bei `:` und U6. Von neun eigenen Mutanten sind sieben rot, einer ist äquivalent (R), und einer bleibt grün (S). Der Code wählt heute richtig, wie die Sonde zeigt. Vor der Closure fehlt der Test, der die Auswahl nach dem genauen Namen hält (F-534). Hexagon und Bootstrap sind unverändert.

**Übergabe:** F-534, F-530 und F-531 gehen an den Implementer, F-532 an den Architect. F-533, F-535 und F-536 gehen an den Verifier. Dieser Report ersetzt keine Verifikation.
