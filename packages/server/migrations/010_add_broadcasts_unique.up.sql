DELETE FROM broadcasts b
USING broadcasts b2
WHERE b.channel_id = b2.channel_id
  AND b.url = b2.url
  AND b.id > b2.id;

ALTER TABLE broadcasts
ADD CONSTRAINT broadcasts_channel_id_url_unique UNIQUE (channel_id, url);
