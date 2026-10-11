**Vorgang:** slice-v1-abschluss-einspielen-tls
**Fund:** Der Code-Commit `147afea` entschied vier Randformen, die §6 nicht nannte (Senden des `SSLRequest` als `PGR-E4002`, Ort der Testhilfen, Typ der Zertifikate, Hook-Signatur im Leser), und gab sie erst danach zurück; der Architect bestätigte sie in `c84df26` (Review F-605, MEDIUM). Der Ablauf war der der Regel, nur die Reihenfolge war falsch.
