# Architecture

## Database

This project uses _two_ database engines. The reason of this choice is that we want
the graph traversal speed of ArangoDB, but the schema enforcement of PostgreSQL (and all the other goodies of
this DBMS).

### ArangoDB

ArangoDB is used to stored relationships that are likely to have a graph depth >= 2.  
Such relationships are for example:

- Follows / Followers relationships

### PostgreSQL

PostgreSQL contains everything else that doesn't need to be interacted with a graph traversal.
For example, the followers relationship of a user will be stored in ArangoDB, whereas the user display name will
be stored in PostgreSQL.