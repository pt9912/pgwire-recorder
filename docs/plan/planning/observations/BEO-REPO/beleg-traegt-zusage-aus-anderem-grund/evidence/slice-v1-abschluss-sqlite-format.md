**Vorgang:** slice-v1-abschluss-sqlite-format
**Fund:** Der Mutant N1b (Primärschlüssel von `session` entfernt) macht 33 Tests rot, weil das Schema ohne Primärschlüssel den Fremdschlüssel ablehnt, nicht nur den benannten; ein Mutant, der allein die Eindeutigkeit aufhebt, macht genau `TestSQLiteDoppelteSessionNummer` rot (Review F-620, INFO). Damals als einmalig nicht eingetragen.
