class Admin::LanguageForm < Language
  include ActiveFormModel

  permit :progress,
    :learn_as,
    :slug,
    :cover,
    :openai_assistant_id,
    :hexlet_program_landing_page

  skip_if_empty :cover
end
