# frozen_string_literal: true

class Admin::LanguageForm < Language
  include ActiveFormModel

  permit :progress,
    :learn_as,
    :slug,
    :cover,
    :openai_assistant_id

  skip_if_empty :cover
end
