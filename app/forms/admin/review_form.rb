class Admin::ReviewForm < Review
  include ActiveFormModel

  permit :user_id,
    :pinned,
    :language_id,
    :locale,
    :state,
    :body,
    :first_name,
    :last_name
end
