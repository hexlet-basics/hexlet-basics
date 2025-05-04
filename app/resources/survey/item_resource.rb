class Survey::ItemResource < ApplicationResource
  typelize_from Survey::Item

  attributes :id, :survey_id, :value, :state, :order, :slug
end
