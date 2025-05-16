class SurveyResource < ApplicationResource
  typelize_from Survey

  attributes :id, :slug, :question, :description, :locale, :state, :created_at

  # has_many :items, resource: Survey::ItemResource

  meta do
    {}
  end
end
