SELECT
  word.Lemma
  , word.Pos
FROM
  word
WHERE
  word.Lang = 'jpn'
  AND word.Pos in ('a', 'v', 'n');