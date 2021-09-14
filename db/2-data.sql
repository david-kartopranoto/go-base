INSERT INTO app_user (email, password, username, created_at)
select 'dummy@dummy.com', 'dummy', 'dummy', NOW();