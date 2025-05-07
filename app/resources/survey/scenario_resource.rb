class Survey::ScenarioResource < ApplicationResource
  typelize_from Survey::Scenario

  attributes :id

  meta do
    {}
  end
end
