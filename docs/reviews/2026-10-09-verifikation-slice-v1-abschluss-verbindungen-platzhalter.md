# Verifikation: slice-v1-abschluss-verbindungen-platzhalter — 2026-10-09

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-v1-abschluss-verbindungen-platzhalter.md`, F-519 bis F-529). F-527 und F-529 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungen-platzhalter.md` im geschnittenen Zuschnitt (*Laden*), Kopf und §1 bis §8, gegen den Gesamt-Diff `539ebee..eaf8164`. Architect: `c297f9b`, `5eb485a`, `7a4dc1c`, `c4e7063`. Schnitt: `5c09d04`. Implementer: `f0180b6`, `d7e1833`, `f9b9584`, `a3bcad0`, `f4152c3`, `0c26212`, `eaf8164`. Review: `e87fe04`.

**Abgrenzung:** Das *Benutzen* der Verbindung (Name und `host:port` bei `--upstream` und `PGWIRE_RECORDER_UPSTREAM`, `sslmode=require` bei `record`, Einsetzen, `PGR-E2005`, Port nach dem Einsetzen, Verbinden, Handbuch) liegt bei `slice-v1-abschluss-upstream-verbinden`. Dass es fehlt, ist kein Befund.

**Eingang:**

- der Plan ganz, besonders §1, §2 (DoD), §3, §4, §6 (*Randformen*, A1 bis A11, L1 bis L6, *Gegengelesen*, *Befunde des Reviews*, *Grenze der Tests*, *Risiken*) und §7 (*Belege des Implementers*, *Nacharbeit zum Review*)
- `spec/spezifikation.md` `LH-FA-17.a` ganz, mit *Benannte Verbindungen*, *Geheimnisse*, *Fehler* und *Anzeige*
- am Stand `eaf8164`: `internal/adapters/driving/cli/verbindung.go` ganz, der Diff von `datei.go`, `export_test.go`, `leser_test.go`, `test/integration/konfiguration_e2e_test.go`; von den Tests `verbindung_test.go` (Zerlegung, `TestDateiDollar`, `TestDateiPlatzhalterAusserhalb`, `TestDateiUpstream`, `TestDateiUpstreamUngueltig`); die Wertemengen in `cli.go`
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- der Nehmer `slice-v1-abschluss-upstream-verbinden` (§1 *Übernimmt*, §6 *Randformen* ab *A3, Zusammensetzen*, *Risiken*)

Beim Start war der Arbeitsbaum sauber, HEAD `eaf8164`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

**So wurden die Proben gebaut.** Alles lief in einem eigenen Unterverzeichnis des Scratchpads, nie im Repo.

- **Binary:** eine frische Kopie per `git archive eaf8164 | tar -x`, ohne `.git` und ohne `cp -p`. Aus ihr wurde die Stufe `build` des `Dockerfile` mit eigenem Tag gebaut und das statische Binary herauskopiert. Jede Probe schrieb eine Datei und rief `config show --config <datei>` und `record --config <datei>` je mit einer Frist von 10 s auf, ohne Server und ohne `--listen`: Eine gültige Datei endet bei `record` deshalb mit `PGR-E2001 … Pflichtoption --listen fehlt`, also nach dem Laden der ganzen Datei und vor jedem Verbindungsaufbau. Jede Ausgabe auf `stderr` habe ich auf `GEHEIM` geprüft, bei Exit ungleich 0 auch `stdout`.
- **Mutanten:** je Mutant eine neue Kopie per `git archive eaf8164 go.mod go.sum cmd internal | tar -x`. Die Ersetzung machte ein Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Danach lief `go test -count=1 ./internal/adapters/driving/cli/` im Image der Stufe `deps`, mit der Kopie als Bind-Mount und ohne Netz; danach wurde nur das Verzeichnis des Mutanten gelöscht. Der Grundlauf ohne Mutation war grün. Die E2E-Mutanten liefen als Image der Stufe `integration` aus einer mutierten Kopie mit `-test.run '^TestE2EConfigShowVerbindungUngueltig$'`; `PGR_OHNE_POSTGRES=1` überspringt nur das Warten auf PostgreSQL, der Test verbindet sich nicht. Die ungeänderte Kopie lief zum Vergleich genauso.
- **Danach:** Kopien und eigene Images sind entfernt. `make gates` lief im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Zerlegung und Grammatik der URL nach A1 bis A6, L2, L3, L5 und F-521; Name einer Verbindung nach V-121 und L1 (`LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Zerlegung von links, Benutzerteil am letzten `@`, Passwort am ersten `:`, Datenbank bis `?` auch mit `/` (A1) | `TestVerbindungZerlegung` gelesen (`a@b@h/a/b`). Binary: `postgresql://h/a/b?sslmode=disable` gültig | bestätigt |
| Teil mit Benutzer, Host und Port endet an `/`, `?` und `#` (A1, F-519) | Binary: `postgresql://h?x/db` und `postgresql://h#x/db` je `PGR-E2004 … Datenbank fehlt`. **V-F519a/b/c** (endet nur an `/`; nicht an `?`; nicht an `#`): alle **rot**, `TestVerbindungUngueltig` | bestätigt |
| Leerer Benutzer `PGR-E2004`, leeres Passwort `PGR-E2006` (A2) | Binary: `postgresql://@h/db` ergibt `PGR-E2004`, `postgresql://u:@h/db` `PGR-E2006` | bestätigt |
| IPv6 in Klammern, hinter `]` nur `:` mit Port oder das Ende; `[` ohne `]` und `[]` ungültig (A3) | Binary gültig: `[::1]`, `u@[::1]:5`, `[fe80::1%25eth0]:5432`, `[::ffff:1.2.3.4]`; ungültig: `[abc`, `[]`, `::1` ohne Klammern. **V-A3-hinter**: **rot**, `TestVerbindungUngueltig` | bestätigt |
| Inhalt des Hosts (F-521) | Binary, je `PGR-E2004` an `connections.v` ohne `GEHEIM`: `[abc]`, `[GEHEIM]:5`, `[a b]`, `[1.2.3.4]`, `[${H}]`, `a%3Ab`, `GEHEIM%3Ab`, `h%2Fx`, `a%20b`, `a%25b`, `a%C2%A0b` (NBSP als Unicode-Leerraum), `a b${N}`. Gültig: `db-${N}.example`, `${H}`. **V-F521-ipv6**, **-ipv4** (`IsValid` statt `Is6`), **-platzh**, **-doppelpunkt**, **-leerraum**, **-vor-dekod**: alle **rot** (`TestVerbindungUngueltig`, bei IPv4 und Leerraum auch `TestDateiUpstreamUngueltig`). Zum Inhalt einer Zone siehe V-122 | bestätigt |
| Port: Default `5432`, Ziffern 1 bis 65535 wie geschrieben, führende Nullen, `:` ohne Port ungültig (Rückgabe 6, A4) | Binary: `h:0`, `h:65536`, `h:` je `PGR-E2004`, `h:05432` gültig. **V-Port-max** (Grenze 65536): **rot**, `TestDateiUpstreamUngueltig`, `TestVerbindungUngueltig` | bestätigt |
| Port mit Platzhalter: wörtliche Zeichen Ziffern (L5) | Binary: `h:5${P}` gültig, `h:x${P}` und `h:$$${P}` `PGR-E2004`. **V-L5**: **rot**, `TestVerbindungUngueltig` | bestätigt |
| Prozent-Dekodierung, UTF-8, keine Steuerzeichen (Rückgabe 10, L3, A5) | Binary: `%FF`, `%C2%80`, `%7F` je `PGR-E2004`; `%C3%A4` und `%c3%a4` gültig; `sslmode=%72equire` gültig. **V-L3-utf8**: **rot**, `TestVerbindungUngueltig` | bestätigt |
| Steuerzeichen im geschriebenen Text, auch der Tabulator (L2) | Binary: Tabulator in der Datenbank `PGR-E2004`. **V-L2-tab** (Tabulator ausgenommen): **rot**, `TestVerbindungUngueltig` | bestätigt |
| `sslmode` nur `disable` und `require`, höchstens einmal, Platzhalter an jeder Stelle ungültig (A6, F-520) | Binary: `prefer`, doppeltes `sslmode`, `re${X}quire`, `ssl${X}mode=disable` je `PGR-E2004`. **V-F520-M5**, **-M6**, **-M43** (Anker `^` entfernt), **V-ssl-doppelt**: alle **rot**, `TestVerbindungUngueltig` | bestätigt |
| Fragment ungültig, auch in einer nicht benutzten Verbindung | Binary: `postgresql://h/db#GEHEIM` ergibt `PGR-E2004` ohne `GEHEIM`; jede Probe stand in einer Verbindung, die kein Schlüssel nennt | bestätigt |
| Form der Platzhalter vor dem Rest des Teils (F-523) | Binary: `postgresql://${1}/db` `PGR-E2004`. **V-F523-klammer-zuerst** (`[ ohne ]` vor den Platzhaltern): **rot**, `TestVerbindungUngueltig` | bestätigt |
| Name ohne `:`, `@` und `$`, Stelle nur `connections` (V-121, L1) | Binary, je `PGR-E2004 … connections: Name einer Verbindung mit :, @ oder $` ohne den Namen: `postgresql://u:GEHEIM@h/db`, `a:GEHEIM`, `GEHEIM@b`, `a$GEHEIM`, `a$$GEHEIM`, `${GEHEIM}`. Gültig: `staging`, `a.b-c`, `mit Leerraum`. **V-L1**, **V-V121-at**, **V-V121-dp**, **V-V121-stelle** (Name in der Stelle): alle **rot**, `TestVerbindungName` | bestätigt |

### Punkt 2: `$$` und `${` außerhalb einer URL, Kommandozeile und Umgebung ohne beides, Schlüssel `upstream` gegen Namen oder `host:port` (L4, F-521), `config show` (`LH-FA-17`). **Bestätigt.** Zum Beleg in §7 siehe V-124.

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `${` außerhalb einer URL `PGR-E2004`, gleich ob gültig geformt, von links gelesen (A10) | Binary: `output: "${GEHEIM}"`, `output: "a$$${GEHEIM}"`, `log_level: "${1}"`, `replay.input: "${GEHEIM}"` (Abschnitt eines anderen Kommandos, auch bei `record`) je `PGR-E2004 … Platzhalter außerhalb einer URL` ohne `GEHEIM`. `output: "a$${X}"` und `output: "a$b"` gültig. **V-A10**: **rot**, `TestDateiPlatzhalterAusserhalb` | bestätigt |
| Die Wertemenge prüft den Text nach `$$` | Binary: `log_level: "in$$fo"` ist `PGR-E2004`. Zur Unterscheidbarkeit siehe Abschnitt 2, F-527 | bestätigt |
| Kommandozeile und Umgebung ohne `$$` (A11) | `TestDateiDollar` gelesen: `--output=a$${X}` bleibt `a$${X}`, `PGWIRE_RECORDER_OUTPUT=${X}$$` bleibt wie gesetzt. §7 DoD 2, Zeilen 5 und 6 | bestätigt |
| `upstream` nennt eine gültige Verbindung, auch eine nach ihm, oder hat die Form `host:port` (A8, A7) | Binary gültig: `staging` vor und nach `connections:`, `[::1]:5`, `h:5`, `[fe80::1%eth0]:5`. Ungültig, je `PGR-E2004 … record.upstream` ohne `GEHEIM`: `GEHEIM@h:5`, `GEHEIM h:5`, `h\x01:5`, `[GEHEIM]:5`, `a/b:5`, `[abc]:5`, `a%3Ab:5`, `h%41:1`, `:5432`, `h`, `Staging`, `${X}`, `GEHEIM`, `h:$$5`, `a$$b` gegen eine Verbindung `ab`. **V-A8-aus**, **V-A8-nachher** (keine Namen gesammelt), **V-F521-at-hostport**, **-ctrl-hostport**, **-ipv6-hostport**: alle **rot**, `TestDateiUpstreamUngueltig` (bei V-A8-nachher auch `TestDateiUpstream`) | bestätigt |
| Ein ungültiger Name zählt nicht, der erste Fehler steht in der Datei zuerst (L4) | Binary: `upstream: "a@b"` vor `connections:` mit `"a@b"` ergibt `record.upstream`; umgekehrt `connections:`. **V-L4**: **rot**, `TestDateiUpstreamUngueltig` | bestätigt |
| Auch bei `config show` und einem anderen Kommando | Jede Probe oben lief mit `config show` und `record` mit demselben Code; `config show` zeigte bei Exit 2 nichts. E2E-Mutant **V-E2E-upstream** (Prüfung von `upstream` entfernt): **rot**, `TestE2EConfigShowVerbindungUngueltig` (Exit 0 mit Anzeige); ungeändert grün | bestätigt |
| `config show` zeigt Platzhalter unaufgelöst und `$$` wie geschrieben | Binary: `postgresql://u:${PW}@db-${N}.example:${P}/db?sslmode=require` und `output: "a$$b"` erscheinen wörtlich, Exit 0. **V-ZEIGEN-dollar** (Knoten auf den Text nach `$$` gesetzt): **rot**, `TestDateiDollar`. Unaufgelöst ist strukturell: In diesem Stand setzt nichts ein | bestätigt |

### Punkt 3: Klartext-Passwort `PGR-E2006` in jeder Verbindung, auch bei `config show` (F-522, L6) (`LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Passwortteil, der nicht genau ein `${VAR}` ist | Binary, je `PGR-E2006` bei `config show` und `record`, ohne `GEHEIM`, `stdout` leer: `u:GEHEIM@h`, `u:@h`, `u:$${PW}@h`, `u:${PW}x@h`, `u:${1}@h`. **V-Klartext-aus**, **V-Klartext-anker** (`$` am Ende entfernt), **V-E2006-code**: alle **rot**, `TestVerbindungKlartext` | bestätigt |
| Parameter `password`, auch ohne `=`, leer und dekodiert (F-522, A5) | Binary, je `PGR-E2006`: `?password`, `?password=`, `?password=GEHEIM`, `?pass%77ord=GEHEIM`. `?PASSWORD=GEHEIM` und `?pass${X}word=GEHEIM` sind `PGR-E2004`. **V-F522-gleich-zuerst** (`=` vor dem Namen), **V-F522-pw-mit-gleich** (`password` nur mit `=`): beide **rot**, `TestVerbindungKlartext` | bestätigt |
| Passwort und Wert von `password` nicht dekodiert (L6) | Binary: `u:%ZZ@h` und `?password=%ZZ` je `PGR-E2006`. **V-L6** (Passwort vor der Prüfung dekodiert): **rot**, `TestVerbindungKlartext` | bestätigt |
| Reihenfolge innerhalb der URL | Binary: `postgresql://u:GEHEIM@[x/db?a=b` ergibt `PGR-E2006` (Passwort vor Host); `postgres://:GEHEIM@/` ergibt das Schema. Das akzeptierte Negativ aus §6 hält: `u:GEHEIM?w@h/db` ergibt `PGR-E2004` (Port), ohne `GEHEIM` | bestätigt |
| Auch bei `config show`, Exit 2, keine Ausgabe (Integration) | Eigener `make gates`: `TestE2EConfigShowVerbindungUngueltig` PASS. E2E-Mutant **V-E2E-e2006** (Code `PGR-E2004`): **rot**; ungeändert grün | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf.** Beleg in §7 nur für den Stand vor der Nacharbeit, siehe V-123.

Eigener Lauf im Arbeitsbaum am Stand `eaf8164`: Exit 0. Darin `d-check: 397 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `--- PASS: TestE2EConfigShowVerbindungUngueltig`, `run-integration-tests: gruen`, `abdeckung-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `baseline-verify: v6.16.0 OK`.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`e87fe04`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; §7 trägt `<…>`, alle vier Risiken in §6 tragen „offen bis Closure“ bzw. den Ausgang *eingetreten* mit Folge-Slice. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Das zweite Risiko (Größe) trat nicht ein, der Diff war nach F-529 in einer Sitzung prüfbar; das dritte (Ablage) trägt, siehe Abschnitt 3; das vierte schließen DoD-Punkt 1 und 3, am Binary bestätigt. In die Register-Fortschreibung gehören die Klassen aus der Summary-Zeile des Reviews und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, über die §8 die Closure entscheiden lässt.

---

## 2. Review-Befunde F-519 bis F-529: nachgefahren

| Befund | Prüfung | Urteil |
|---|---|---|
| F-519 Grenzen aus A1 | V-F519a, V-F519b, V-F519c je **rot**, `TestVerbindungUngueltig` | erledigt |
| F-520 Platzhalter neben Text | V-F520-M5, V-F520-M6, V-F520-M43 je **rot**; M4 ist durch F-522 entfallen (der Name `pass${X}word` scheitert vor dem Vergleich), V-F522-pw-mit-gleich ist rot | erledigt |
| F-521 Inhalt des Hosts | `LH-FA-17.a` entscheidet ihn seit `c4e7063`; Binary und Mutanten siehe Punkt 1 und 2. Offen bleibt der Inhalt einer Zone (V-122) | erledigt bis auf V-122 |
| F-522 Reihenfolge je Parameter | siehe Punkt 3 | erledigt |
| F-523 Reihenfolge | V-F523-klammer-zuerst **rot**; die Abdeckungs-Deklaration nennt die Reihenfolge als geprüfte Paare | erledigt |
| F-524 Teilzusagen | Binary: `%c3%a4` gültig, `%7F` und `upstream: "a[b:1"`-verwandte Fälle ungültig; die Mutanten dazu stehen in §7 *Nacharbeit* | erledigt |
| F-525, F-526 | Der Kommentar an `zerlegung` nennt nur `stelle` und `v`; §6 nennt `verbindung_test.go` | erledigt |
| F-527 grüne Mutanten | **V-GRUEN1** (Wertemenge prüft den Text vor `$$`) und **V-GRUEN2** (`$$` auf den Namen angewandt) laufen **grün**, wie §7 sagt. Die berichtigte Begründung hält: `artText` prüft nur auf leer, und `$$` ergibt nie einen leeren Text; `wahrheitswert`, `dauer` und `stufe` lehnen jeden Text mit `$` in beiden Formen ab, und ihre Fehlertexte sind feste Sätze, `dauer` verzweigt nur hinter der Regex, die ein `$` nie trifft. `upstream` ist `text` und vergleicht danach mit dem Text nach `$$`. Zu V-GRUEN2: Jeder Name mit `$` ergibt nach `$$` wieder einen Text mit `$`, sonst scheitert `${`; ohne `$` bleibt er gleich, leer wird er nie, und die übrigen Prüfungen von `nameFehler` sehen dieselben Zeichen. Beide sind äquivalent; das Register-Muster `BEO-REPO/gruener-mutant-faelschlich-aequivalent` liegt nicht vor | trägt |
| F-528 Ablage | siehe Abschnitt 3 | trägt |
| F-529 Größe | 1066 Zeilen Code und Tests im Diff `539ebee..eaf8164` (ohne Plan und Abdeckung) gegen geschätzte 650 bis 750; der Reviewer prüfte in einer Sitzung. Für die Closure: Das zweite Risiko trat nicht ein | Hinweis für die Closure |

---

## 3. Ablage für den Folge-Slice (F-528)

`slice-v1-abschluss-upstream-verbinden` erwartet in §1 und §6 (*Rekursion*, *A3*, *F-528*, *A12*, Risiko zwei), dass das Laden die URL vor dem Einsetzen in Teile zerlegt, je Teil die Platzhalter in der Reihenfolge der URL trägt und einen Host mit `:` nur aus Klammern oder einer Variable liefert. Am Code (`verbindung`, `stueck`, `teil`) und an `TestVerbindungZerlegung`:

- Benutzer, Passwort, Host, Port und Datenbank sind je eine Folge von Stücken, wörtlich oder Platzhalter. Die Felder stehen in der Reihenfolge der URL. Parameter tragen keinen Platzhalter (F-520), `password` ist immer Klartext. Damit liegen alle Platzhalter einer gültigen URL in diesen fünf Teilen, und die erste nicht gesetzte Variable in der Reihenfolge der URL (A12) ist ablesbar.
- Host und Port sind getrennt; `record` kann die Teile, die es ignoriert, ohne Auswertung übergehen.
- Ein Host ohne Klammern enthält nach der Dekodierung kein `:` (Binary: `a%3Ab` und `GEHEIM%3Ab` ungültig, V-F521-doppelpunkt und V-F521-vor-dekod rot). Die Regel aus F-528, „Klammern genau dann, wenn der Host ein `:` enthält“, ist damit eindeutig.
- Der Port steht wie geschrieben (`05432`), ohne Angabe als `5432`; `sslmode` ist abgelegt; `mitBenutzer` und `mitPasswort` sagen, ob die URL den Teil schreibt.

Die Ablage trägt, was der Folge-Slice erwartet; das dritte Risiko in §6 trat nicht ein. Der Inhalt einer Zone erreicht den Folge-Slice ungeprüft (V-122).

---

## 4. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`, `a-check-negativ: gruen`.
- `verbindung.go` importiert neben der Standardbibliothek (`net/netip`, `regexp`, `strconv`, `strings`, `unicode`, `unicode/utf8`) nur `internal/hexagon/model` für Fehler und Codes; keine YAML-Bibliothek. Kein Typ der Bibliothek verlässt den Adapter.
- Schicht-Abgrenzung (§1): `git diff --stat 539ebee..eaf8164 -- internal test` ändert nur `internal/adapters/driving/cli/` und `test/integration/konfiguration_e2e_test.go`. Kern, Code-Tabelle, PGWire-Adapter, Driven-Adapter, Bootstrap, `cmd/`, `go.mod` und `go.sum` sind unberührt. `spec/spezifikation.md` ändern nur die Architect-Commits `7a4dc1c` und `c4e7063`, wie §3 sagt.

---

## 5. Plan gegen Code

- **§1:** Geliefert ist genau das *Laden*. Kein Code löst einen Namen bei `--upstream` auf, setzt eine Variable ein oder erzeugt `PGR-E2005`; das Handbuch ist nicht berührt. Die Konstante `PGR-E2006` hat jetzt einen Erzeuger (`klartext`).
- **§3:** Jede geänderte Datei steht dort: der CLI-Adapter mit Unit-Tests (`verbindung.go` neu, `datei.go`, `export_test.go`, `leser_test.go`), `test/integration` mit `config show` für Klartext-Passwort, ungültige URL und ungültiges `upstream`, die Abdeckungstabellen. `make abdeckung-check` ist grün.
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-17.a`. Ausnahme ist die Zone in `host:port` (V-122): §6 *F-521* sagt dort „dazu kein Steuerzeichen“, die Spezifikation lässt die Zone in Klammern offen, der Code nimmt Steuerzeichen und Leerraum in der Zone an.
- **§7:** Die Tabellen decken alle drei Liefer-Punkte in der Form *Zusage · Mutation · roter Test*. Jede Zeile, die ich nachgefahren habe, deckt sich: 39 eigene Unit-Mutanten, 37 rot, 2 grün (äquivalent, wie §7 sagt), dazu 2 E2E-Mutanten, beide rot. Es fehlen eine Zeile zu `config show` und `$$` (V-124) und der Gate-Lauf am Stand der Nacharbeit (V-123).

---

## 6. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-122 | LOW | Der Inhalt einer Zone ist nicht entschieden, und der Code nimmt in `host:port` an, was §6 ausschließt. §6 *F-521* sagt für `host:port`: „dieselben Regeln, Zone hinter `%`; dazu kein Steuerzeichen“. `LH-FA-17.a` sagt: „in eckigen Klammern eine IPv6-Adresse, auch mit einer Zone hinter `%`, sonst kein Steuerzeichen, kein Leerraum …“, und „sonst“ gilt nur ohne Klammern. `ipv6` übergibt den Text an `netip.ParseAddr`, und diese Funktion nimmt jede nicht leere Zone an. Binary, je mit Exit 0 bei `config show`: `upstream: "[fe80::1%G\x01]:5"` (Steuerzeichen in der Zone), `"[fe80::1%GEHEIM h]:5"` (Leerraum), `"[fe80::1%a/b]:5"`. In der URL fangen L2 und L3 Steuerzeichen ab, `[fe80::1%25a%20b]` und `[fe80::1%25a%2Fb]` laden aber mit Leerraum bzw. `/` in der abgelegten Zone. Die Zeile „kein Steuerzeichen in `host:port`“ in §7 *Nacharbeit* prüft nur einen Host ohne Klammern (`h\x01…:5`). Eine Meldung nennt den Wert nicht, deshalb LOW. *Failure-Szenario:* Der Folge-Slice benutzt dieselbe Prüfung für `--upstream` und `PGWIRE_RECORDER_UPSTREAM` (§6 dort, *F-521*), erbt die offene Zone und setzt sie beim Zusammensetzen unverändert in die Adresse für den Verbindungsaufbau. | `internal/adapters/driving/cli/verbindung.go` · `a, err := netip.ParseAddr(text)` und `return ok && ipv6(host) && portForm(port)` | Architect: Inhalt der Zone in URL und `host:port` entscheiden (`AGENTS.md` §3.12) und §6 *F-521* mit `LH-FA-17.a` in Deckung bringen; danach Implementer: Test mit Mutation (§3.10), sonst die Zeile in §7 enger fassen (§3.11) |
| V-123 | LOW | Für den Code-Stand nach der Nacharbeit steht in §7 kein Beleg für `make gates`. Die Tabelle *Läufe* belegt `make gates` für `a3bcad0`. Die Nacharbeit `0c26212` ändert `verbindung.go` und die Tests, und §7 *Nacharbeit* sagt dazu nur: „`make gates` nach dem Commit dieses Abschnitts (Bericht des Implementers)“. Das verweist auf den Bericht, nicht auf ein Ergebnis im Plan. Ebenso fehlt `make test-integration` an diesem Stand. Mein eigener Lauf am Stand `eaf8164` ist grün (Punkt 4); der DoD-Punkt ist damit bestätigt, der Beleg im Plan fehlt trotzdem (seit slice-harness-blackbox-kern). | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungen-platzhalter.md` §7 *Nacharbeit zum Review* · „`make gates` nach dem Commit dieses Abschnitts (Bericht des Implementers)“ | Implementer: Ergebnis des Gate-Laufs am Stand `eaf8164` in §7 nachtragen |
| V-124 | INFO | Die Zusage aus DoD-Punkt 2 „`config show` zeigt `$$` wie geschrieben“ hat in §7 keine Zeile *Zusage · Mutation · roter Test*. Die DoD verlangt sie für jede Zusage. Der Test existiert (`TestDateiDollar`, Vergleich der ganzen Ausgabe), und mein Mutant V-ZEIGEN-dollar ist rot. Die Zusage ist damit gehalten, nur ihr Beleg fehlt im Plan. „Platzhalter unaufgelöst“ ist in diesem Stand strukturell wahr, weil nichts einsetzt; eine Mutation dazu gibt es erst mit dem Folge-Slice. | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungen-platzhalter.md` §7, Tabelle *DoD-Punkt 2* | Implementer: Zeile nachtragen |

---

## 7. Gefahrene Mutanten

Alle am Stand `eaf8164`, je in einer frischen Kopie, `go test -count=1 ./internal/adapters/driving/cli/`. Ausnahme sind die beiden E2E-Mutanten.

| Mutant | Änderung | Ergebnis |
|---|---|---|
| V-GRUND | keine | grün |
| V-F519a | Teil mit Benutzer, Host und Port endet nur an `/` | rot: `TestVerbindungUngueltig` |
| V-F519b | endet nicht an `?` | rot: `TestVerbindungUngueltig` |
| V-F519c | endet nicht an `#` | rot: `TestVerbindungUngueltig` |
| V-F520-M5 | `!w.platzhalter() &&` bei `sslmode` entfernt | rot: `TestVerbindungUngueltig` |
| V-F520-M6 | Prüfung auf Platzhalter im Namen eines Parameters entfernt | rot: `TestVerbindungUngueltig` |
| V-F520-M43 | `^` in `platzhalterName` entfernt | rot: `TestVerbindungUngueltig` |
| V-F521-ipv6 | IPv6-Prüfung in Klammern entfernt | rot: `TestVerbindungUngueltig` |
| V-F521-ipv4 | `a.IsValid()` statt `a.Is6()` | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungUngueltig` |
| V-F521-platzh | Platzhalter in Klammern zugelassen | rot: `TestVerbindungUngueltig` |
| V-F521-doppelpunkt | `:` aus den unzulässigen Zeichen entfernt | rot: `TestVerbindungUngueltig` |
| V-F521-leerraum | Leerraum zugelassen | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungUngueltig` |
| V-F521-vor-dekod | `:` im dekodierten Host übergangen | rot: `TestVerbindungUngueltig` |
| V-F521-at-hostport | `@` in `host:port` zugelassen | rot: `TestDateiUpstreamUngueltig` |
| V-F521-ctrl-hostport | Steuerzeichen in `host:port` zugelassen | rot: `TestDateiUpstreamUngueltig` |
| V-F521-ipv6-hostport | in Klammern nur nicht leer | rot: `TestDateiUpstreamUngueltig` |
| V-F522-gleich-zuerst | `=` vor dem Namen geprüft | rot: `TestVerbindungKlartext` |
| V-F522-pw-mit-gleich | `password` nur mit `=` als Klartext | rot: `TestVerbindungKlartext` |
| V-F523-klammer-zuerst | `[ ohne ]` vor der Form der Platzhalter | rot: `TestVerbindungUngueltig` |
| V-L1 | `$` im Namen erlaubt | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungName` |
| V-V121-at | `@` im Namen erlaubt | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungName` |
| V-V121-dp | `:` im Namen erlaubt | rot: `TestVerbindungName` |
| V-V121-stelle | Name in der Stelle der Meldung | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungName` |
| V-L4 | jeder Schlüssel unter `connections` zählt als Name | rot: `TestDateiUpstreamUngueltig` |
| V-A8-aus | Prüfung von `upstream` entfernt | rot: `TestDateiUpstreamUngueltig` |
| V-A8-nachher | keine Namen gesammelt | rot: `TestDateiUpstream`, `TestDateiUpstreamUngueltig` |
| V-A10 | `${` außerhalb einer URL als Text | rot: `TestDateiPlatzhalterAusserhalb` |
| V-A3-hinter | Text hinter `]` nicht geprüft | rot: `TestVerbindungUngueltig` |
| V-Port-max | Grenze 65536 | rot: `TestDateiUpstreamUngueltig`, `TestVerbindungUngueltig` |
| V-L5 | Ziffernprüfung bei Port mit Platzhalter entfernt | rot: `TestVerbindungUngueltig` |
| V-ssl-doppelt | doppeltes `sslmode` angenommen | rot: `TestVerbindungUngueltig` |
| V-L2-tab | Tabulator im Text der URL ausgenommen | rot: `TestVerbindungUngueltig` |
| V-L3-utf8 | UTF-8-Prüfung nach der Dekodierung entfernt | rot: `TestVerbindungUngueltig` |
| V-Klartext-aus | Klartext-Prüfung des Passwortteils entfernt | rot: `TestVerbindungKlartext` |
| V-Klartext-anker | `$` am Ende von `genauEinPlatzhalter` entfernt | rot: `TestVerbindungKlartext` |
| V-L6 | Passwort vor der Klartext-Prüfung dekodiert | rot: `TestVerbindungKlartext` |
| V-E2006-code | `PGR-E2004` statt `PGR-E2006` | rot: `TestVerbindungKlartext` |
| V-ZEIGEN-dollar | `config show` zeigt den Text nach `$$` | rot: `TestDateiDollar` |
| V-GRUEN1 | Wertemenge prüft den Text vor `$$` | grün, äquivalent (Abschnitt 2, F-527) |
| V-GRUEN2 | `$$` auf den Namen einer Verbindung angewandt | grün, äquivalent (Abschnitt 2, F-527) |
| V-E2E-e2006 | `PGR-E2004` statt `PGR-E2006`, Image der Stufe `integration` | rot: `TestE2EConfigShowVerbindungUngueltig`; ungeändert grün |
| V-E2E-upstream | Prüfung von `upstream` entfernt, Image der Stufe `integration` | rot: `TestE2EConfigShowVerbindungUngueltig` (Exit 0 mit Anzeige); ungeändert grün |

Jeder rote Mutant scheitert aus dem richtigen Grund; die gelesenen Fehlerzeilen nennen den Fall der mutierten Zusage. Der erste Ansatz für V-F521-ipv4 übersetzte nicht (Variable ohne Verwendung) und ist neu formuliert rot gesehen.

---

## Verdikt

**DoD-Liefer-Punkte 1 bis 3 und `make gates`: bestätigt.** Das Binary aus einer frischen Kopie hält jede Stichprobe aus `LH-FA-17.a` bei `config show` und `record`: Klartext-Passwörter in allen drei Formen (URL, `?password`, `?password=`) sind `PGR-E2006`. Hosts wie `[abc]`, `[GEHEIM]:5`, `a%3Ab` und `GEHEIM@h:5` sind `PGR-E2004`. `$$` und Platzhalter gelten in Namen und Werten, wie die Spezifikation sagt. `[::1]`, `[fe80::1%25eth0]` und `db-${N}.example` sind gültig. Keine der 113 Proben nennt `GEHEIM` in einer Meldung. Das Hexagon bleibt rein, `make a-check` ist grün. Die Review-Befunde F-519 bis F-526 sind erledigt. Die berichtigte Begründung der zwei grünen Mutanten (F-527) hält, und die Ablage trägt, was der Folge-Slice erwartet (F-528).

**Vor der Closure:** V-122 an den Architect, weil der Inhalt einer Zone eine Randform ist (`AGENTS.md` §3.12); er betrifft auch die Prüfung, die `slice-v1-abschluss-upstream-verbinden` für Option und Umgebung benutzt. V-123 und V-124 an den Implementer, nur Plan, ohne Code. Danach stehen die Closure-Pflichten in §7 aus.

**Summary:** 0 HIGH · 0 MEDIUM · 2 LOW · 1 INFO (V-122: In Klammern nimmt `host:port` eine Zone mit Steuerzeichen oder Leerraum an, die URL eine mit dekodiertem Leerraum oder `/`, entgegen §6 *F-521*, und die Spezifikation lässt die Zone offen; V-123: §7 trägt keinen Gate-Beleg für den Stand nach der Nacharbeit, mein eigener Lauf ist grün; V-124: Die Zusage „`config show` zeigt `$$` wie geschrieben“ hat keine Zeile in §7, mein Mutant ist rot).
