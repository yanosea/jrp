UPDATE
  jrp
SET
  IsFavorite = 0
  , UpdatedAt = datetime(CURRENT_TIMESTAMP, '+9 hours')
WHERE
  jrp.IsFavorite = 1;
