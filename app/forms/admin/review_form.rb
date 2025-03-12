class Admin::ReviewForm < Review
  include ActiveFormModel

  permit :user_id,
    :pinned,
    :state_event,
    :language_id,
    :locale,
    :body,
    :first_name,
    :last_name
end
