INSERT INTO `account` (username, password_hash, display_name) VALUES
('test-user1', 'hashed_password1', 'TestUser1'),
('test-user2', 'hashed_password2', 'TestUser2'),
('test-user3', 'hashed_password3', 'TestUser3'),
('test-user4', 'hashed_password4', 'TestUser4'),
('test-user5', 'hashed_password5', 'TestUser5'),
('test-user6', 'hashed_password6', 'TestUser6'),
('test-user7', 'hashed_password7', 'TestUser7'),
('test-user8', 'hashed_password8', 'TestUser8'),
('test-user9', 'hashed_password9', 'TestUser9');

INSERT INTO `status` (account_id, content) VALUES
(1, 'Test content for user 1'),
(2, 'Test content for user 2'),
(3, 'Test content for user 3'),
(4, 'Test content for user 4'),
(5, 'Test content for user 5'),
(6, 'Test content for user 6'),
(7, 'Test content for user 7'),
(8, 'Test content for user 8'),
(9, 'Test content for user 9');
