UPDATE
  jrp
SET
  IsFavorite = 1
  -- add 9 hours to the current time to get JST
  , UpdatedAt = datetime(CURRENT_TIMESTAMP, '+9 hours')
WHERE
  jrp.IsFavorite = 0
  AND jrp.ID IN (%s);
