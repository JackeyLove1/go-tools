CREATE DATABASE IF NOT EXISTS test2;
use test2;

DROP TABLE IF EXISTS student;
DROP TABLE IF EXISTS grades;

CREATE TABLE student
(
    id   INT(11)      NOT NULL UNIQUE AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    age  INT(11),
    PRIMARY KEY (id)
);

INSERT INTO student (name, age)
VALUES ('John Doe', 20),
       ('Jane Smith', 22),
       ('Bob Johnson', 21);

INSERT INTO student (name, age)
VALUES ('Jacky', 88),
       ('Cat', 100),
       ('Dog', 60);

CREATE TABLE grades
(
    id         INT         NOT NULL UNIQUE AUTO_INCREMENT,
    student_id INT(11)     NOT NULL,
    course     VARCHAR(30) NOT NULL,
    grade      VARCHAR(5)  NOT NULL,
    PRIMARY KEY (id),
    foreign key (student_id) references student (id)
);

INSERT INTO grades (student_id, course, grade)
VALUES (1, 'Math', 'A'),
       (1, 'English', 'B'),
       (2, 'Math', 'A'),
       (3, 'Science', 'A');

INSERT INTO grades (student_id, course, grade)
VALUES (3, 'Math', 'C'),
       (3, 'Branch', 'A'),
       (2, 'All', 'D'),
       (1, 'Science', 'B');

# INNER JOIN
SELECT * from student as s INNER JOIN grades as g ON s.id = g.student_id;

# LEFT JOIN
SELECT * from student as s LEFT JOIN grades as g ON s.id = g.student_id;

# RIGHT JOIN
SELECT * from student as s RIGHT JOIN grades as g ON s.id = g.student_id;

