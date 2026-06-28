# typed: true

class ApplicationResource
  include Alba::Resource
  include Typelizer::DSL

  meta do
    {}
  end

  def to_model
    object
  end
end
