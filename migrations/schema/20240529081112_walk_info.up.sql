CREATE OR REPLACE VIEW walk_info AS
SELECT 
  w.id as walk_id,
  w.start_time,
  w.finish_time,
  w.state,
  owner.id as owner_id,
  owner.name as owner_name,
  owner.email as owner_email,
  walker.id as walker_id,
  walker.name as walker_name,
  walker.email as walker_email,
  p.id as pet_id,
  p.name as pet_name,
  p.age as pet_age,
  p.additional_info as pet_additional_info,
  w.state as walk_state
FROM walks as w, users as owner, users as walker, pets as p
WHERE 
  w.owner_id = owner.id AND
  w.walker_id = walker.id AND
  w.pet_id = p.id;
