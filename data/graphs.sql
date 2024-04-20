CREATE TABLE graphs (
  id varchar PRIMARY KEY,
  graph_name varchar NOT NULL
);

CREATE TABLE nodes (
  id varchar,
  node_name varchar,
  graph_id varchar,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(id)
  ON DELETE CASCADE
);

CREATE TABLE edges (
  id varchar,
  graph_id varchar,
  from_node varchar,
  to_node varchar,
  cost numeric,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(id) ON DELETE CASCADE,
  FOREIGN KEY (from_node, graph_id) REFERENCES nodes(id, graph_id) ON DELETE CASCADE,
  FOREIGN KEY (to_node, graph_id) REFERENCES nodes(id, graph_id) ON DELETE CASCADE
);

-- Function that finds cycles in a graph
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

-- maybe
-- CREATE TABLE shortest_paths (
--   from_node varchar,
--   to_node varchar,
--   graph_id varchar,
--   node_path json,
--   PRIMARY KEY (from_node, to_node),
--   FOREIGN KEY (graph_id) REFERENCES graphs(id),
--   FOREIGN KEY (from_node, graph_id) REFERENCES nodes(id, graph_id),
--   FOREIGN KEY (to_node, graph_id) REFERENCES nodes(id, graph_id)
-- );