ALTER TABLE ads DROP COLUMN updated_at;
ALTER TABLE ads DROP COLUMN created_at;
ALTER TABLE ads ADD id INT(11) NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (id);
alter table ads ADD board VARCHAR(255) default '';

UPDATE ads SET uri = REPLACE( uri , "public" , "st" );

ALTER TABLE antispam CHANGE name ip VARCHAR(255);
ALTER TABLE antispam DROP COLUMN updated_at;
ALTER TABLE antispam DROP COLUMN created_at;

ALTER TABLE bans DROP COLUMN updated_at;
ALTER TABLE bans DROP COLUMN created_at;
ALTER TABLE bans ADD COLUMN ip VARCHAR(255);

ALTER TABLE mods DROP COLUMN updated_at;
ALTER TABLE mods DROP COLUMN created_at;

ALTER TABLE users DROP COLUMN updated_at;
ALTER TABLE users DROP COLUMN created_at;


