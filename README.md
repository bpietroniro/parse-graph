# Overview

# Setup Instructions

## Installation

## Database

## Environment

# Implementation Notes

## XML Parsing

This project utilizes the `github.com/beevik/etree` Golang package to parse XML inputs. This lightweight package is built on top of the standard library's `encoding/xml` package, and provides an easily traversable document-oriented model.

I began the project first by using the native `encoding/xml` package. This is a standard struct-based solution. Some sources warn that it doesn't perform the best for large or especially complex inputs, but graph scale wasn't a primary concern for this project, so it seemed a suitable choice.

However, once it came time to handle input validation logic, I reexamined several options just to see if any might be a little more ergonomic. The document-tree data model of the `etree` package made writing methodical checks of nested tags quite intuitive—akin to accessing objects in the browser DOM—and to my eyes, more readable than my original struct-based approach.

## SQL Schema
Here is an ERD to demonstrate how I set up the SQL schema:

![image](./schema.png)

(also see comments in `graphs.sql`)

## JSON Parsing

## Path Algorithms

### All Paths from A to B

### Cheapest Path from A to B

# Thoughts for Further Improvement

## Caching

## Project Layout