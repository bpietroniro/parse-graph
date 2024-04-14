CREATE TABLE graphs (
  id varchar PRIMARY KEY,
  graph_name NOT NULL
);

CREATE TABLE nodes (
  id varchar,
  node_name varchar,
  graph_id varchar,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(graph_id)
);

CREATE TABLE edges (
  id varchar,
  edge_name varchar,
  graph_id varchar,
  from_node varchar,
  to_node varchar,
  cost numeric,
  PRIMARY KEY (id, graph_id),
  FOREIGN KEY (graph_id) REFERENCES graphs(graph_id),
  FOREIGN KEY (from_node, graph_id) REFERENCES nodes(id, graph_id),
  FOREIGN KEY (to_node, graph_id) REFERENCES nodes(id, graph_id)
);