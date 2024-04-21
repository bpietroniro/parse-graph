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
RETURNS TABLE(node_path varchar[]) AS $$
BEGIN
  RETURN QUERY
  WITH RECURSIVE paths AS (
    SELECT NULL::varchar AS from_node, start_node AS to_node, ARRAY[start_node] AS node_path FROM edges e WHERE e.graph_id = graph_id_in
    UNION    
    SELECT e.from_node, e.to_node, paths.node_path || e.to_node::varchar
    FROM edges e JOIN paths ON e.from_node = paths.to_node AND e.graph_id = graph_id_in WHERE NOT paths.is_cycle
  ) CYCLE to_node SET is_cycle TO TRUE DEFAULT FALSE USING cycle_path
  SELECT paths.node_path FROM paths WHERE NOT paths.is_cycle AND paths.to_node = end_node;
END;
$$ LANGUAGE plpgsql;

-- maybe
-- CREATE UNLOGGED TABLE path_cache (
--   from_node varchar,
--   to_node varchar,
--   graph_id varchar,
--   node_path varchar[],
--   inserted_at timestamp,
--   PRIMARY KEY (from_node, to_node, graph_id),
--   FOREIGN KEY (graph_id) REFERENCES graphs(id),
--   FOREIGN KEY (from_node, graph_id) REFERENCES nodes(id, graph_id),
--   FOREIGN KEY (to_node, graph_id) REFERENCES nodes(id, graph_id)
-- );