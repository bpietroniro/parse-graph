# Overview

# Setup Instructions

## Installation

## Database

## Environment

# Implementation Notes

## XML Parsing

This project utilizes `github.com/beevik/etree` to parse XML inputs. `etree` is built on top of the standard library's `encoding/xml` package, and provides an easily traversable document-oriented model.

I began the project using the native `encoding/xml` package, a standard struct-based solution. Some sources warn that it doesn't perform the best for large or especially complex inputs, but large-scale graphs weren't a concern for this project, so it seemed suitable.

Once it came time to handle input validation, I reexamined other options just to see if any might be a little more ergonomic. The document-tree data model of the `etree` package made validating nested tags quite intuitive—akin to accessing objects in the browser DOM—and to my eyes, more readable than my original struct-based approach.

## SQL Schema

Here is an ERD to demonstrate how I set up the SQL schema:

![image](./schema.png)

(also see comments in `graphs.sql`)

## JSON Parsing

## Path Algorithms

Upon first glance at the problem, my first instinct was to implement Dijkstra's algorithm on the application level. This would efficiently find cheapest paths, and with a memoization component, finding all paths could be a convenient byproduct. The only database query needed would be to load the graph data initially.

However, after I'd gotten the `find_graph_cycles` SQL procedure working, I realized this procedure could be modified to return all paths between two nodes. This isn't the most efficient or database-friendly approach, but I'll be honest: I was so happy to have managed this in SQL that I wanted to do it again. 🤓 

### All Paths from A to B

The task of finding all paths between two given nodes is accomplished in the `find_all_paths` PL/pgSQL function defined in `graphs.sql`.

The function returns the results of a recursive SQL query that uses a CTE (Common Table Expression) to keep track of paths between nodes in the graph. With the use of PostgreSQL's `CYCLE` syntax, it avoids falling into infinite cyclic loops. It also keeps track of which paths contain cycles with the boolean `is_cycle` column, making it easy to filter out cyclic paths for the results.

To make finding cheapest paths easier, this function also tracks and returns the total cost of each path.

### Cheapest Path from A to B

If we've already calculated all paths between two nodes, finding the cheapest path is easy: just choose the path with the minimum cost.

Since our inputs may contain both "all paths" and "cheapest" queries for the same pair of nodes, I decided to implement an ephemeral cache of sorts within the `handleJSON` function. For each query, we first check local memory to see if we've already loaded the results of `find_all_paths` for the given start and end nodes. If so, we don't need to bother the database again.

# Thoughts for Further Improvement

## Efficiency

My current approach is not very database-friendly.

### Execute Queries As a Single Batch

### Move Pathfinding to Application

## Caching

## Project Layout