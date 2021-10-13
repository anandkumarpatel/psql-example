create table services (
  name varchar not null primary key,
  description varchar not null,
  versions varchar[] not null
);

-- example data

INSERT INTO services(name, description, versions)
VALUES ('first', 'cool thing', ARRAY[md5(RANDOM()::TEXT), md5(RANDOM()::TEXT)]);

INSERT INTO services(name, description, versions)
VALUES ('next', 'broken thing', ARRAY[md5(RANDOM()::TEXT)]);

INSERT INTO services(name, description, versions)
VALUES ('before', 'legacy', ARRAY[md5(RANDOM()::TEXT), md5(RANDOM()::TEXT)]);