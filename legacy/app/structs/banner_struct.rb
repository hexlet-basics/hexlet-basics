# typed: strict
# frozen_string_literal: true

class BannerStruct < T::Struct
  include ApplicationParamsStruct::Base

  const :locale, T.nilable(String)
  const :body, T.nilable(String)
  const :url, T.nilable(String)
  const :background, T.nilable(String)
  const :state, T.nilable(String)
  const :starts_at, T.nilable(String)
  const :finishes_at, T.nilable(String)
end
