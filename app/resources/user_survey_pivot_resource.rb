class UserSurveyPivotResource < ApplicationResource
  typelize_from UserSurveyPivot

  attributes :id, :user_id, :goal_item_id, :coding_experience_item_id, :created_at

  typelize :string, nullable: true
  attribute :goal do |obj|
    obj.goal_item&.value
  end

  typelize :string, nullable: true
  attribute :coding_experience do |obj|
    obj.coding_experience_item&.value
  end
end
