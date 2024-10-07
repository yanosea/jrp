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
ORDER BY
  jrp.ID DESC
LIMIT ?;
