# Overview

## What This Does
This project encompasses a program that:
- Parses XML files containing representations of directed graphs
- Validates XML and graph-specific syntax and semantics
- Stores graph data in a normalized SQL schema in PostgreSQL
- Finds cycles in a given graph using SQL
- Accepts JSON inputs to query:
  - all paths between two nodes in a given graph
  - the cheapest path between two nodes in a given graph
- Fulfills these queries and outputs JSON-formatted answers

# Setup Instructions

The project uses Golang version 1.20 and PostgreSQL version 15.6.

## Installation

1. Unzip the project into a directory of your choosing.

2. Navigate to the root project directory (`parse-graph`).

3. Ensure that Go version 1.20 or later is [installed](https://go.dev/dl/) and currently in use in your environment. Run `go install`.

4. Ensure that PostgreSQL version 15.6 is [installed](https://www.postgresql.org/download/). To ensure that it's running, you may need the following command:
```bash
# On Linux
sudo service postgresql start

# On macOS with Homebrew
brew services start postgresql
```

## Database

1. Ensure you're logged in to PostgreSQL as a privileged user, and create a new database (e.g. `graph_data`). You can do this from the command line:
```bash
createdb graph_data
```
2. Set up the database schema and functions by reading in the commands from `data/graphs.sql`:
```bash
psql -d graph_data -f ./data/graphs.sql
```

## Environment

The application connects to the DB with the help of environment variables and the `github.com/joho/godotenv` package. In the project root, create a `.env` file, and populate it like so:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=<your_username>
DB_PW=<your_password>
DB_NAME=graph_data
```
Change the variables as needed for your local PostgreSQL configuration and choice of database name. If your chosen PostgreSQL user doesn't need a password, you can leave `DB_PW` empty.

## Run the Program

If all else went as planned, you should be able to run the program with `go run .` from the project root directory. (Or create an executable using `go build` if you'd like.)

### Command Line Arguments

To parse an XML file containing a graph representation, provide the filepath as an argument:

```bash
go run . test-inputs/xml/intersecting_cycles.xml
```

To parse a JSON file containing queries on a graph, provide the filepath as an argument, followed by the graph's ID as another argument:

```bash
go run . test-inputs/json/sample_query_1.json
```

Some examples (not exhaustive, but varied) are provided in the `test-inputs` folder.

# Implementation Notes

## XML Parsing

This project utilizes `github.com/beevik/etree` to parse XML inputs. `etree` is built on top of the standard library's `encoding/xml` package, and provides an easily traversable document-oriented model.

I began the project using the native `encoding/xml` package, a standard struct-based solution. Some sources warn that it doesn't perform the best for large or especially complex inputs, but large-scale graphs weren't a concern for this project, so it seemed suitable.

Once it came time to handle input validation, I reexamined other options just to see if any might be a little more ergonomic. The document-tree data model of the `etree` package made validating nested tags quite intuitive—akin to accessing objects in the browser DOM—and to my eyes, more readable than my original struct-based approach.

## SQL Schema

Here is an ERD to demonstrate the SQL schema for graphs:

![image](./schema.png)

The comments in `graphs.sql` provide further information on the relationships between these tables.

## JSON Parsing

I was able to implement the required query parsing using Golang's native `encoding/json` package without much trouble. I did briefly look into some alternatives, some of which are built on top of `encoding/json`, but to keep things simple I decided to stick with this go-to option, at least for now.

That said, I'm not quite satisfied with the structs I've set up for JSON purposes in `models/queries.go`. A little repetitive/clunky. I may not be using the `json` package in the most intuitive way possible, or it may be that different package could save the trouble of having to fuss with structs in this way. I plan to revisit this. Any ideas welcome!

## Path Algorithms

My first instinct was to implement Dijkstra's algorithm on the application level. This would efficiently find cheapest paths, and with a memoization component, finding all paths could be a convenient byproduct. The only database query needed would be to load the graph data initially.

However, after I'd gotten the `find_graph_cycles` SQL procedure working, I realized this procedure could be modified to return all paths between two nodes. This isn't the most efficient or database-friendly approach, but honestly, I was so happy to have managed the first bit with SQL that I wanted to do it again. 🤓 

### All Paths from A to B

The task of finding all paths between two given nodes is accomplished in the `find_all_paths` PL/pgSQL function defined in `graphs.sql`.

The function returns the results of a recursive SQL query that uses a CTE (Common Table Expression) to keep track of paths between nodes in the graph. With the use of PostgreSQL's `CYCLE` syntax, it avoids falling into infinite cyclic loops. It also keeps track of which paths contain cycles with the boolean `is_cycle` column, making it easy to filter out cyclic paths for the results.

To make finding cheapest paths easier, this function also tracks and returns the total cost of each path.

### Cheapest Path from A to B

If we've already calculated all paths between two nodes, finding the cheapest path is easy: just choose the path with the minimum cost.

Since our inputs may contain both "all paths" and "cheapest" queries for the same pair of nodes, I decided to implement an ephemeral cache of sorts within the `handleJSON` function. For each query, we first check local memory to see if we've already loaded the results of `find_all_paths` for the given start and end nodes. If so, we don't need to bother the database again.


# Thoughts for Further Improvement

## Efficiency

My current approach is not very database-friendly. In the worst case, it executes a recursive function in the database for every query in the input list. For the current purposes, optimization was not a priority, but to scale well with increased input length, frequency, or graph complexity, we'd want to do things differently.

Here are two potential alternatives:

### Execute Queries As a Single Batch
Instead of handling each query one-by-one, we could prepare one large batch query:
1. Iterate once through the query list, keeping track of start-end node pairs in memory.
2. Instead of calling `find_all_paths` on one node pair at a time, call it once, returning results for all pairs.
3. Aggregate the results by the value of the start node, or by node pairs.

### Move Pathfinding to Application

Dijkstra's algorithm would work well here.  This would move the computational work to the application level. It might not be as efficient as database functions under the hood, but if we needed to scale horizontally down the line, it would be easier to replicate the application than to replicate the database.

I've included an in-progress implementation in `paths.go`. It's not tested yet, and the application doesn't use it currently. Out of curiosity, I plan to finish it up and compare its performance with the database-centric approach.

## Caching

Another great way to increase efficiency would be to use a proper caching mechanism. I spent some time brainstorming and researching a few ways to do this, but for the sake of time decided not to implement them (at least not yet):

1. Use Redis as a cache to store serialized JSON outputs for repeat queries
2. Set up another PostgreSQL table to store paths as they are generated (this could get real big real fast, though)
3. Use PostgreSQL statement caching 

## Project Organization

I think there are several possible ways to organize this project (I've already evolved through a couple of them), and I'm not totally convinced by how it's laid out currently. For example I'd like to:

- possibly separate XML and JSON handling into two separate modules
- clean up the `handleJSON` function currently in `main`; break it into smaller functions
- figure out a struct schema for JSON marshalling that is less clunky than what I currently have in `queries.go`

## Tests

Thus far I've relied on manual, iterative testing while developing this project. Given more time, I'd like to develop a more robust suite of tests, and set up test automation.

## Feedback Welcome!

I'd be grateful for any questions, comments, criticism, and ideas about this project!
