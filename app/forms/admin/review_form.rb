# frozen_string_literal: true

class Admin::ReviewForm < Review
  include ActiveFormModel

  permit :user_id,
         :state_event,
         :language_id,
         :locale,
         :body,
         :first_name,
         :last_name
end
