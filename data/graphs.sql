CREATE TABLE graphs (
  id varchar PRIMARY KEY,
  graph_name varchar NOT NULL
);

-- Holds data for graph node inputs.
-- Primary key indexed on the combination of unique ID and corresponding graph's unique ID.
-- graphs:nodes => one to many
CREATE TABLE nodes (
  id varchar,
  node_name varchar,
  graph_id varchar,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(id) ON DELETE CASCADE
);

-- Holds data for graph edge inputs.
-- Primary key indexed on the combination of unique ID and corresponding graph's unique ID.
-- graphs:edges => one to many
-- edges.from_node:nodes => many to one
-- edges.to_node:nodes => many to one
CREATE TABLE edges (
  id varchar,
  graph_id varchar,
  from_node varchar,
  to_node varchar,
  cost numeric,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(id) ON DELETE CASCADE,
  FOREIGN KEY (from_node, graph_id) REFERENCES nodes(id, graph_id) ON DELETE CASCADE,
  FOREIGN KEY (to_node, graph_id) REFERENCES nodes(id, graph_id) ON DELETE CASCADE,
  UNIQUE (from_node, to_node, graph_id)
);

-- Function that finds cycles in a given graph
CREATE OR REPLACE FUNCTION find_graph_cycles(graph_id_in varchar)
RETURNS TABLE(node_path varchar[]) AS $$
BEGIN
  RETURN QUERY

  WITH RECURSIVE paths AS (
    SELECT e.from_node, e.to_node, ARRAY[e.to_node::varchar] AS node_path FROM edges e WHERE e.graph_id = graph_id_in

    UNION    

    SELECT e.from_node, e.to_node, paths.node_path || e.to_node::varchar
    FROM edges e JOIN paths ON e.from_node = paths.to_node AND e.graph_id = graph_id_in WHERE NOT paths.is_cycle
  ) CYCLE to_node SET is_cycle TO TRUE DEFAULT FALSE USING cycle_path

  SELECT paths.node_path FROM paths WHERE paths.is_cycle AND paths.to_node = paths.node_path[1];
END;
$$ LANGUAGE plpgsql;

-- Function that finds all paths from A to B in a given graph 
CREATE OR REPLACE FUNCTION find_all_paths(graph_id_in varchar, start_node varchar, end_node varchar)
RETURNS TABLE(node_path varchar[], total_cost numeric) AS $$
BEGIN
  RETURN QUERY

  WITH RECURSIVE paths AS (
    SELECT
      NULL::varchar AS from_node,
      start_node AS to_node,
      ARRAY[start_node] AS node_path,
      0::numeric AS total_cost
    FROM edges e
    WHERE e.graph_id = graph_id_in

    UNION    

    SELECT e.from_node, e.to_node, paths.node_path || e.to_node::varchar, paths.total_cost + e.cost
    FROM edges e
    JOIN paths ON e.from_node = paths.to_node AND e.graph_id = graph_id_in
    WHERE NOT paths.is_cycle
  ) CYCLE to_node SET is_cycle TO TRUE DEFAULT FALSE USING cycle_path

  SELECT paths.node_path, paths.total_cost FROM paths WHERE NOT paths.is_cycle AND paths.to_node = end_node;
END;
$$ LANGUAGE plpgsql;
