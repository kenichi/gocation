CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;

CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    point GEOGRAPHY,
    polygon GEOMETRY,
    longitude NUMERIC NOT NULL,
    latitude NUMERIC NOT NULL,
    accuracy REAL,
    speed REAL,
    bearing REAL,
    timestamp TIMESTAMP,
    altitude REAL,
    vertical_accuracy REAL,
    battery INT,
    charging BOOLEAN,
    properties JSON);

CREATE INDEX locations_points_gist ON locations USING gist(point);
CREATE INDEX locations_polygon_gist ON locations USING gist(polygon);
