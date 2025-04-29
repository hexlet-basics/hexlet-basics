TRUNCATE TABLE user_survey_pivots RESTART IDENTITY CASCADE;

INSERT INTO user_survey_pivots (
  user_id,
  study_plan_item_id,
  coding_experience_item_id,
  goal_item_id,
  created_at,
  updated_at
)
SELECT
  u.id,
  MAX(CASE WHEN s.slug = 'study-plan' THEN sa.survey_item_id END) AS study_plan_item_id,
  MAX(CASE WHEN s.slug = 'coding-experience' THEN sa.survey_item_id END) AS coding_experience_item_id,
  MAX(CASE WHEN s.slug = 'goal' THEN sa.survey_item_id END) AS goal_item_id,
  now(),
  now()
FROM users u
LEFT JOIN survey_answers sa ON sa.user_id = u.id AND sa.state = 'fulfilled'
LEFT JOIN surveys s ON s.id = sa.survey_id
GROUP BY u.id
HAVING COUNT(sa.id) > 0;
