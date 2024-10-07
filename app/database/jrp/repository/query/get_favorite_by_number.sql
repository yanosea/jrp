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
  jrp.IsFavorite = 1
ORDER BY
  ID DESC
LIMIT ?;
