SELECT
  jrp.ID
  , jrp.Phrase
  , jrp.Prefix
  , jrp.Suffix
  , jrp.IsFavorite
  , jrp.CreatedAt
  , jrp.UpdatedAt
FROM
  jrp
WHERE
  (%s)
ORDER BY
  jrp.ID ASC;
