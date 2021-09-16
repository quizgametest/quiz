DROP TABLE answers;
DROP TABLE questions;
DROP TABLE game;
DROP TABLE users;
CREATE TABLE questions (
  id smallserial not null check(id > 0) primary key,
  question varchar(255) not null);

CREATE TABLE answers (
  question_id smallint not null references questions (id),
  answer varchar(255) not null,
  is_correct boolean default false);

CREATE TABLE users (
  id smallserial not null check(id > 0) primary key,
  username varchar(255) not null);

CREATE TABLE game (
  id smallserial not null check(id > 0) primary key,
  user_id smallint not null references users (id),
  rightAnswer smallint not null,
  answered boolean default false);

insert into questions values (1, 'first question');
insert into answers values (1, 'answer-one-true', true);
insert into answers values (1, 'answer-two');
insert into answers values (1, 'answer-three');

insert into questions values (2, 'second question');
insert into answers values (2, 'answer-one');
insert into answers values (2, 'answer-two-true', true);
insert into answers values (2, 'answer-three');

insert into questions values (3, 'third question');
insert into answers values (3, 'answer-one');
insert into answers values (3, 'answer-two');
insert into answers values (3, 'answer-three-true', true);

insert into users (username) values ('Modest');
