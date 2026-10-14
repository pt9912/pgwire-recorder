**Stand:** offen

Kein Auftreten: `slice-v1-abschluss-tls-client-doku` — betrifft die Version des Clients (libpq 17), nicht des Servers; das Handbuch nennt sie seit der Nacharbeit zu V-1 im Satz. Der Zähler bleibt bei 1.

Kein Auftreten: `slice-v1-abschluss-anmeldung-doku` — das Handbuch nennt „PostgreSQL 17“ an allen Stellen zur Anmeldung und sagt, dass andere Versionen nicht geprüft sind (Gegenbeleg; Review F-3 fand die fehlende Version vor `1291bcd`). Die Matrix weiterer Versionen bleibt Gegenstand von `slice-v1-abschluss-postgres-versionen`, dessen Plan vier Schichten zählt; ein Eintrag der Handbuch-Sätze dort unterblieb. Der Zähler bleibt bei 1.

Kein Auftreten: `slice-v1-abschluss-postgres-versionen-matrix` — die Adresse liefert: Das Gate fährt die serverabhängige Teilmenge (Szenarien 1, 2 und 7, Fall SCRAM-SHA-256, Lebendprüfungs-Abgleich) gegen 14, 15, 16 und 18, die Grenze 17 stimmt für alle. Der Zähler bleibt bei 1, der Eintrag bleibt `offen`: Abnahmeszenario 18 und Szenario 16 Satz 3 laufen nur gegen 17 (Entscheidung des Nutzers vom 2026-10-12), Klartext und MD5 der Anmeldung nur gegen 17, und die Handbuch-Sätze zur Anmeldung zieht `slice-v1-abschluss-postgres-versionen-zusage` nach.
