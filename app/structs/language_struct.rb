# typed: strict
# frozen_string_literal: true

class LanguageStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :progress, T.nilable(String)
  const :learn_as, T.nilable(String)
  const :slug, T.nilable(String)
  const :openai_assistant_id, T.nilable(String)
  const :hexlet_program_landing_page, T.nilable(String)
end
