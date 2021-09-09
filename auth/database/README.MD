# Users Database Schema
- The email column has the type citext (case-insensitive text). 
This type stores text data exactly as it is inputted — without changing
the case in any way — but comparisons against the data are always case-insensitive… 
including lookups on associated indexes.
We’ve also got a UNIQUE constraint on the email column. Combined with the citext type, 
this means that no two rows in the database can have the same email value — even if they
have different cases. This essentially enforces a database-level business rule that no two 
users should exist with the same email address.


- The password column has the type bytea (binary string). In this column we’ll store a one-way hash of the user’s password generated using bcrypt —
not the plaintext password itself.