DELETE
FROM
  jrp
WHERE
  jrp.IsFavorite = 0
  AND jrp.ID IN (%s);
