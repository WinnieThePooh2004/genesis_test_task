create database nbu_rates;

\connect nbu_rates;

create table subscription (
  id int generated always as IDENTITY primary key,
  email varchar(255) not null unique
);