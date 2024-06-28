INSERT INTO user_cms (id, created_at, updated_at, deleted_at, username, password, flag, level, token) VALUES (1, 1717578152, 1719474702, null, 'adminplanetban', '88741ccf1a3251cc435c23a20e4e3b10', false, 'Superadmin', '0ecc592dface49925711487b7159b330');

create type platform_type as enum ('Tiktok', 'Website', 'Instagram', 'Facebook');

create type role_type as enum ('Report', 'Admin', 'Superadmin');

create type source_type as enum ('Ads', 'Kol');